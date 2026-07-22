package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const width, height, tile = 480, 720, 48
const saveKey = "ebiShowcaseQuest."

const (
	battleChoose = iota
	battlePlayerAction
	battleEnemyAction
)

// A small, explicit route graph makes the world read like a classic tile RPG:
// # is scenery and . is a walkable road. Every quest target is on one route.
var worldRoute = [...]string{
	"##########",
	"#######..#",
	"#######..#",
	"#####....#",
	"#####.####",
	"#####.####",
	"##.......#",
	"##.##.####",
	"##.##.####",
	"#..##.####",
	"#.....####",
	"##########",
}

//go:embed assets/quest-pearl-kingdom.png assets/quest-moon-arena.png assets/quest-party-atlas.png assets/quest-enemy-atlas.png
var questArtFS embed.FS

var (
	questArtOnce sync.Once
	questArt     map[string]*ebiten.Image
	questParty   [2]*ebiten.Image
	questEnemies [3]*ebiten.Image
	questFace14  *text.GoTextFace
	questFace16  *text.GoTextFace
	questFace20  *text.GoTextFace
)

type game struct {
	x, y, quest, hp, enemyHP, enemyMax, enemy, scene int
	turn, tick, shake, flash                         int
	action, actionTick, battlePhase, pendingChoice   int
	companion, clear                                 bool
	defend, effectApplied                            bool
	message                                          string
	audio                                            *audio.Context
	gate                                             audiolab.Gate
	pulse                                            *shaderlab.Pulse
	cam                                              cameralab.State
	badge                                            *ebiten.Image
}

func newGame() *game {
	loadQuestArt()
	b := ebiten.NewImage(20, 20)
	b.Fill(color.RGBA{255, 210, 80, 255})
	g := &game{x: 1, y: 10, hp: 60, message: "Meet Momo in the southwest village.", audio: audiolab.Context(), pulse: shaderlab.NewPulse(), cam: cameralab.State{ViewW: width, ViewH: height}, badge: b}
	g.load()
	return g
}
func (g *game) load() {
	q, ok := storageGet(saveKey + "quest")
	if !ok {
		return
	}
	x, _ := storageGet(saveKey + "x")
	y, _ := storageGet(saveKey + "y")
	hp, _ := storageGet(saveKey + "hp")
	g.quest, _ = strconv.Atoi(q)
	g.x, _ = strconv.Atoi(x)
	g.y, _ = strconv.Atoi(y)
	g.hp, _ = strconv.Atoi(hp)
	g.companion = g.quest > 0
	g.setMessage()
}
func (g *game) save() {
	storageSet(saveKey+"quest", strconv.Itoa(g.quest))
	storageSet(saveKey+"x", strconv.Itoa(g.x))
	storageSet(saveKey+"y", strconv.Itoa(g.y))
	storageSet(saveKey+"hp", strconv.Itoa(g.hp))
}
func (g *game) Update() error {
	g.tick++
	g.cam.Pos = cameralab.Vec{X: float64(g.x * tile), Y: float64(g.y * tile)}
	if g.shake > 0 {
		g.shake--
	}
	if g.flash > 0 {
		g.flash--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		for _, k := range []string{"quest", "x", "y", "hp"} {
			storageRemove(saveKey + k)
		}
		*g = *newGame()
		return nil
	}
	if g.clear {
		if any() {
			g.resetSave()
			*g = *newGame()
		}
		return nil
	}
	if g.scene == 1 {
		return g.battle()
	}
	dx, dy := 0, 0
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dx = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dx = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		dy = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		dy = 1
	}
	if x, y, ok := press(); ok {
		rx, ry := x-(g.x*tile+24), y-(74+g.y*tile+24)
		if abs(rx) > abs(ry) {
			if rx < 0 {
				dx = -1
			} else {
				dx = 1
			}
		} else {
			if ry < 0 {
				dy = -1
			} else {
				dy = 1
			}
		}
	}
	if dx != 0 || dy != 0 {
		nextX, nextY := g.x+dx, g.y+dy
		if !passableTile(nextX, nextY) {
			g.message = "The reef blocks that route. Follow the glowing road."
			return nil
		}
		g.x, g.y = nextX, nextY
		switch {
		case g.quest == 0 && g.x == 2 && g.y == 9:
			g.quest = 1
			g.companion = true
			g.message = "Momo joined! Find the crystal in the northeast."
		case g.quest == 1 && g.x == 8 && g.y == 2:
			g.startBattle(0, 36, "Crystal Slime guards the shard!")
		case g.quest == 2 && g.x == 8 && g.y == 1:
			g.startBattle(1, 68, "Tower Knight raises its shield!")
		case g.quest == 3 && g.x == 5 && g.y == 6:
			g.startBattle(2, 110, "Shadow Crab ambushes the road!")
		case g.quest == 4 && g.x == 1 && g.y == 10:
			g.quest = 5
			g.clear = true
			g.message = "The village is safe! Quest complete."
		}
		g.save()
	}
	return nil
}
func (g *game) startBattle(enemy, hp int, msg string) {
	g.scene = 1
	g.enemy = enemy
	g.enemyHP = hp
	g.enemyMax = hp
	g.turn = 0
	g.battlePhase = battleChoose
	g.pendingChoice = -1
	g.action = 0
	g.actionTick = 0
	g.effectApplied = false
	g.defend = false
	g.message = msg
}
func (g *game) resetSave() {
	for _, k := range []string{"quest", "x", "y", "hp"} {
		storageRemove(saveKey + k)
	}
}
func (g *game) battle() error {
	if g.battlePhase != battleChoose {
		g.advanceBattle()
		return nil
	}

	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		choice = 2
	}
	if x, y, ok := press(); ok && y > 530 {
		choice = min(2, x/(width/3))
	}
	if choice < 0 {
		return nil
	}
	g.pendingChoice = choice
	g.action = choice + 1
	g.actionTick = 36
	g.battlePhase = battlePlayerAction
	g.effectApplied = false
	g.message = []string{"Ebi rushes in...", "Ebi raises the shell guard...", "Momo begins a healing song..."}[choice]
	return nil
}

func (g *game) advanceBattle() {
	if g.actionTick > 0 {
		g.actionTick--
	}
	if g.actionTick == 18 && !g.effectApplied {
		if g.battlePhase == battlePlayerAction {
			g.applyPlayerAction()
		} else {
			g.applyEnemyAction()
		}
		g.effectApplied = true
	}
	if g.actionTick > 0 {
		return
	}

	if g.battlePhase == battlePlayerAction {
		if g.enemyHP <= 0 {
			g.finishBattle()
			return
		}
		g.battlePhase = battleEnemyAction
		g.action = 4
		g.actionTick = 36
		g.effectApplied = false
		g.message = "Enemy turn: " + enemyIntentName(g.turn%3)
		return
	}

	if g.hp <= 0 {
		g.hp = 60
		g.scene = 0
		g.x, g.y = 2, 9
		g.message = "The party escaped and recovered in the village."
		g.save()
		return
	}
	g.turn++
	g.battlePhase = battleChoose
	g.pendingChoice = -1
	g.action = 0
	g.effectApplied = false
	g.message = "Your turn. Choose a command for the next intent."
	g.save()
}

func (g *game) applyPlayerAction() {
	switch g.pendingChoice {
	case 0:
		d := partyAttackDamage(g.enemy, g.turn, g.companion)
		g.enemyHP -= d
		g.play(720)
		g.flash = 7
		g.shake = 4
		g.message = fmt.Sprintf("Party attack: %d damage!", d)
	case 1:
		g.defend = true
		g.message = "Shell guard ready: the next hit is halved."
	case 2:
		before := g.hp
		g.hp = min(60, g.hp+15)
		g.play(540)
		g.message = fmt.Sprintf("Momo restores %d HP!", g.hp-before)
	}
}

func (g *game) applyEnemyAction() {
	intent := g.turn % 3
	damage := enemyAttackDamage(g.enemy, intent, g.defend)
	g.defend = false
	g.hp -= damage
	g.play(180)
	g.shake = 7
	g.message = fmt.Sprintf("%s: party takes %d damage.", enemyIntentName(intent), damage)
}

func (g *game) finishBattle() {
	if g.enemyHP <= 0 {
		g.quest++
		g.scene = 0
		if g.enemy == 0 {
			g.x, g.y = 8, 2
			g.message = "Crystal recovered! Enter the dark tower."
		} else if g.enemy == 1 {
			g.x, g.y = 8, 2
			g.message = "Tower cleared! Cross the center road."
		} else {
			g.x, g.y = 5, 6
			g.message = "Shadow Crab defeated! Return to the village."
		}
		g.save()
		return
	}
}
func (g *game) play(hz float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, hz, .10)).Play()
}
func (g *game) setMessage() {
	switch g.quest {
	case 0:
		g.message = "Meet Momo in the southwest village."
	case 1:
		g.message = "Defeat the crystal guardian in the northeast."
	case 2:
		g.message = "Enter the dark tower beside the crystal."
	case 3:
		g.message = "Cross the center road. Something is waiting."
	case 4:
		g.message = "Return to the southwest village."
	case 5:
		g.message = "Quest complete."
	}
}
func (g *game) Draw(s *ebiten.Image) {
	if g.scene == 1 {
		g.drawBattle(s)
		return
	}
	drawQuestCover(s, questArt["world"], 0, 0, width, height)
	vector.DrawFilledRect(s, 0, 0, width, height, color.RGBA{3, 20, 35, 14}, false)
	drawWorldRoute(s)

	// The painted landmarks are the actual quest graph, not detached key art.
	drawMapLabel(s, "PEARL VILLAGE", 14, 548)
	drawMapLabel(s, "CRYSTAL TOWER", 328, 118)
	drawMapLabel(s, "OLD BRIDGE", 200, 405)

	// Momo is visibly waiting at the first destination, then follows Ebi.
	if !g.companion {
		drawQuestContain(s, questParty[1], 82, 497, 78, 84, false, 1, false)
		vector.StrokeCircle(s, 121, 540, 34+float32(math.Sin(float64(g.tick)*.1)*3), 2, color.RGBA{255, 211, 232, 205}, true)
	}
	if g.quest == 1 {
		drawQuestContain(s, questEnemies[0], 373, 145, 72, 72, false, .92, false)
	}
	if g.quest == 3 {
		drawQuestContain(s, questEnemies[2], 223, 337, 86, 70, false, .78, false)
	}

	px, py := float64(g.x*tile+24), float64(74+g.y*tile+24)
	stepBob := math.Abs(math.Sin(float64(g.tick)*.16)) * 2
	vector.DrawFilledCircle(s, float32(px), float32(py+19), 24, color.RGBA{2, 12, 24, 100}, true)
	drawQuestContain(s, questParty[0], px-38, py-52-stepBob, 76, 80, false, 1, false)
	if g.companion {
		drawQuestContain(s, questParty[1], px-58, py-23-stepBob, 48, 52, false, 1, false)
	}

	tx, ty := g.targetPosition()
	marker := 24 + float32(math.Sin(float64(g.tick)*.12)*4)
	vector.StrokeCircle(s, float32(tx), float32(ty), marker, 4, color.RGBA{255, 224, 105, 230}, true)
	vector.StrokeCircle(s, float32(tx), float32(ty), marker+7, 1, color.RGBA{117, 242, 255, 150}, true)
	for i := 0; i < 3; i++ {
		a := float64(g.tick)*.035 + float64(i)*math.Pi*2/3
		vector.DrawFilledCircle(s, float32(float64(tx)+math.Cos(a)*float64(marker+10)), float32(float64(ty)+math.Sin(a)*float64(marker+10)), 3, color.RGBA{255, 240, 159, 230}, true)
	}

	vector.DrawFilledRect(s, 10, 10, 460, 82, color.RGBA{4, 15, 31, 226}, true)
	vector.StrokeRect(s, 10, 10, 460, 82, 2, color.RGBA{107, 225, 238, 170}, true)
	drawCenteredQuestLabel(s, fmt.Sprintf("EBI QUEST   •   HP %02d/60   •   QUEST %d/5", g.hp, g.quest), 240, 22, questFace16, color.White)
	drawCenteredQuestLabel(s, g.message, 240, 52, questFace14, color.RGBA{255, 232, 177, 255})
	drawCenteredQuestLabel(s, "MOVE: ARROWS / WASD / TAP TOWARD A TILE", 240, 101, questFace14, color.RGBA{221, 241, 245, 255})
	drawCenteredQuestLabel(s, "AUTOSAVED   •   R: DELETE SAVE", 240, 692, questFace14, color.RGBA{221, 231, 238, 255})
	g.drawBadge(s)
	if g.clear {
		overlay(s, "EBI QUEST COMPLETE!\n\nTAP / SPACE TO PLAY AGAIN")
	}
}
func (g *game) drawBattle(s *ebiten.Image) {
	drawQuestCover(s, questArt["battle"], 0, 0, width, height)
	washes := []color.RGBA{{16, 102, 119, 16}, {81, 43, 126, 28}, {139, 24, 58, 35}}
	vector.DrawFilledRect(s, 0, 0, width, height, washes[g.enemy], false)
	ox := 0.0
	if g.shake > 0 {
		ox = float64((g.tick%3)-1) * 5
	}
	enemyW := []float64{150, 178, 235}[g.enemy]
	enemyH := []float64{145, 220, 205}[g.enemy]
	enemyX := 285.0 + ox - enemyW/2
	actionPulse := battleActionPulse(g.actionTick)
	if g.action == 4 {
		enemyX -= 48 * actionPulse
	}
	enemyY := []float64{148, 115, 128}[g.enemy] + math.Sin(float64(g.tick)*.08)*3
	vector.DrawFilledCircle(s, float32(285+ox), float32(enemyY+enemyH*.78), float32(enemyW*.34), color.RGBA{2, 10, 25, 100}, true)
	drawQuestContain(s, questEnemies[g.enemy], enemyX, enemyY, enemyW, enemyH, false, 1, g.flash > 0)

	// One party atlas supports all actions: Update records the chosen action,
	// while Draw changes pose, position and effects without changing the rules.
	ebiX, ebiY := 30.0, 330.0
	momoX, momoY := 132.0, 354.0
	if g.action == 1 {
		ebiX += 72 * actionPulse
		vector.StrokeLine(s, float32(ebiX+94), 387, 274, 270, 7, color.RGBA{115, 244, 255, 205}, true)
	}
	partyShake := 0.0
	if g.shake > 0 && g.flash == 0 {
		partyShake = -ox
	}
	drawQuestContain(s, questParty[0], ebiX+partyShake, ebiY, 145, 180, false, 1, false)
	if g.companion {
		drawQuestContain(s, questParty[1], momoX+partyShake, momoY, 122, 155, false, 1, false)
	}
	if g.action == 2 {
		for i := 0; i < 3; i++ {
			vector.StrokeCircle(s, 105, 412, float32(50+i*8), float32(5-i), color.RGBA{121, 209, 255, uint8(210 - i*45)}, true)
		}
	}
	if g.action == 3 && g.companion {
		for i := 0; i < 8; i++ {
			a := float64(i)*math.Pi/4 + float64(g.tick)*.05
			x := 192 + math.Cos(a)*float64(42+i%3*5)
			y := 405 + math.Sin(a)*float64(58+i%2*5)
			vector.DrawFilledCircle(s, float32(x), float32(y), 5, color.RGBA{255, 163, 224, 210}, true)
			vector.StrokeLine(s, float32(x+4), float32(y), float32(x+4), float32(y-12), 2, color.RGBA{255, 236, 170, 230}, true)
		}
	}

	names := []string{"CRYSTAL SLIME", "TOWER KNIGHT", "SHADOW CRAB"}
	intent := []string{"NORMAL ATTACK", "HEAVY ATTACK — GUARD!", "QUICK ATTACK"}[g.turn%3]
	vector.DrawFilledRect(s, 42, 12, 396, 105, color.RGBA{4, 13, 29, 226}, true)
	vector.StrokeRect(s, 42, 12, 396, 105, 2, color.RGBA{255, 220, 144, 170}, true)
	drawCenteredQuestLabel(s, fmt.Sprintf("%s   HP %02d/%02d", names[g.enemy], max(0, g.enemyHP), g.enemyMax), 240, 27, questFace16, color.White)
	vector.DrawFilledRect(s, 72, 57, 336, 12, color.RGBA{40, 42, 62, 255}, true)
	vector.DrawFilledRect(s, 75, 60, float32(330*max(0, g.enemyHP)/g.enemyMax), 6, color.RGBA{244, 83, 103, 255}, true)
	intentColor := color.RGBA{126, 242, 255, 255}
	if g.turn%3 == 1 {
		intentColor = color.RGBA{255, 154, 115, 255}
	}
	drawCenteredQuestLabel(s, "NEXT  •  "+intent, 240, 79, questFace14, intentColor)
	phaseLabel := []string{"CHOOSE A COMMAND", "PARTY ACTION", "ENEMY ACTION"}[g.battlePhase]
	drawCenteredQuestLabel(s, phaseLabel, 240, 98, questFace14, color.RGBA{255, 238, 181, 255})

	vector.DrawFilledRect(s, 20, 467, 440, 72, color.RGBA{4, 14, 29, 226}, true)
	vector.StrokeRect(s, 20, 467, 440, 72, 2, color.RGBA{112, 225, 238, 150}, true)
	drawQuestLabel(s, fmt.Sprintf("PARTY HP %02d/60", g.hp), 36, 481, questFace16, color.RGBA{160, 246, 218, 255})
	drawCenteredQuestLabel(s, g.message, 240, 510, questFace14, color.White)
	drawBattleButton(s, 0, 12, "1  ATTACK", color.RGBA{27, 127, 133, 245})
	drawBattleButton(s, 1, 167, "2  GUARD", color.RGBA{45, 77, 137, 245})
	drawBattleButton(s, 2, 322, "3  SONG", color.RGBA{133, 72, 119, 245})
	if g.battlePhase != battleChoose {
		vector.DrawFilledRect(s, 12, 550, 455, 88, color.RGBA{4, 12, 26, 145}, true)
		drawCenteredQuestLabel(s, "RESOLVING — WATCH EACH ACTION", 240, 581, questFace14, color.RGBA{255, 238, 181, 255})
	}
	drawCenteredQuestLabel(s, "TAP A COMMAND   •   ENEMY INTENT IS SHOWN ABOVE", 240, 654, questFace14, color.RGBA{231, 238, 245, 255})
	g.drawBadge(s)
}

func (g *game) drawBadge(dst *ebiten.Image) {
	if !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if g.pulse.Draw(fx, g.badge, float32(g.tick)*.08) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(442, 16)
		dst.DrawImage(fx, op)
	}
}

func loadQuestArt() {
	questArtOnce.Do(func() {
		questArt = make(map[string]*ebiten.Image, 4)
		for key, filename := range map[string]string{
			"world":   "quest-pearl-kingdom.png",
			"battle":  "quest-moon-arena.png",
			"party":   "quest-party-atlas.png",
			"enemies": "quest-enemy-atlas.png",
		} {
			data, err := questArtFS.ReadFile("assets/" + filename)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			questArt[key] = ebiten.NewImageFromImage(decoded)
		}
		partyAtlas := questArt["party"]
		for i := range questParty {
			questParty[i] = partyAtlas.SubImage(image.Rect(i*512, 0, (i+1)*512, 512)).(*ebiten.Image)
		}
		enemyAtlas := questArt["enemies"]
		for i := range questEnemies {
			questEnemies[i] = enemyAtlas.SubImage(image.Rect(i*512, 0, (i+1)*512, 512)).(*ebiten.Image)
		}
		questFace14, _ = uilab.Face("en", 14)
		questFace16, _ = uilab.Face("en", 16)
		questFace20, _ = uilab.Face("en", 20)
	})
}

func drawQuestCover(dst, img *ebiten.Image, x, y, w, h float64) {
	b := img.Bounds()
	scale := math.Max(w/float64(b.Dx()), h/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Min.X), -float64(b.Min.Y))
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-float64(b.Dx())*scale)/2, y+(h-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawQuestContain(dst, img *ebiten.Image, x, y, w, h float64, mirror bool, alpha float32, flash bool) {
	b := img.Bounds()
	scale := math.Min(w/float64(b.Dx()), h/float64(b.Dy()))
	dw, dh := float64(b.Dx())*scale, float64(b.Dy())*scale
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Min.X), -float64(b.Min.Y))
	if mirror {
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(x+(w+dw)/2, y+(h-dh)/2)
	} else {
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x+(w-dw)/2, y+(h-dh)/2)
	}
	if flash {
		op.ColorScale.Scale(1, .28, .28, alpha)
	} else {
		op.ColorScale.ScaleAlpha(alpha)
	}
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawQuestLabel(dst *ebiten.Image, label string, x, y float64, face *text.GoTextFace, c color.Color) {
	if face == nil {
		ebitenutil.DebugPrintAt(dst, label, int(x), int(y))
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, label, face, op)
}

func drawCenteredQuestLabel(dst *ebiten.Image, label string, centerX, y float64, face *text.GoTextFace, c color.Color) {
	if face == nil {
		ebitenutil.DebugPrintAt(dst, label, int(centerX)-len(label)*3, int(y))
		return
	}
	w, _ := text.Measure(label, face, 0)
	drawQuestLabel(dst, label, centerX-w/2, y, face, c)
}

func drawMapLabel(dst *ebiten.Image, label string, x, y float64) {
	w, _ := text.Measure(label, questFace14, 0)
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(w+16), 25, color.RGBA{4, 17, 31, 205}, true)
	vector.StrokeRect(dst, float32(x), float32(y), float32(w+16), 25, 1, color.RGBA{255, 224, 145, 155}, true)
	drawQuestLabel(dst, label, x+8, y+4, questFace14, color.RGBA{255, 242, 209, 255})
}

func drawWorldRoute(dst *ebiten.Image) {
	for y, row := range worldRoute {
		for x, cell := range row {
			if cell != '.' {
				continue
			}
			px, py := float32(x*tile), float32(74+y*tile)
			vector.DrawFilledRect(dst, px+4, py+4, float32(tile-8), float32(tile-8), color.RGBA{77, 222, 221, 38}, true)
			vector.StrokeRect(dst, px+5, py+5, float32(tile-10), float32(tile-10), 1, color.RGBA{188, 250, 224, 115}, true)
			vector.DrawFilledCircle(dst, px+float32(tile)/2, py+float32(tile)/2, 3, color.RGBA{255, 238, 153, 180}, true)
		}
	}
}

func drawBattleButton(dst *ebiten.Image, action int, x float64, label string, fill color.RGBA) {
	vector.DrawFilledRect(dst, float32(x), 550, 145, 88, fill, true)
	vector.StrokeRect(dst, float32(x), 550, 145, 88, 3, color.RGBA{255, 235, 190, 175}, true)
	if action == 0 {
		drawQuestContain(dst, questParty[0], x+6, 555, 52, 52, false, 1, false)
	} else if action == 1 {
		vector.StrokeCircle(dst, float32(x+32), 579, 20, 5, color.RGBA{164, 225, 255, 235}, true)
		vector.StrokeCircle(dst, float32(x+32), 579, 13, 2, color.RGBA{255, 244, 202, 220}, true)
	} else {
		drawQuestContain(dst, questParty[1], x+6, 555, 52, 52, false, 1, false)
	}
	drawQuestLabel(dst, label, x+52, 568, questFace14, color.White)
	drawQuestLabel(dst, []string{"15 DMG", "1/2 HURT", "+15 HP"}[action], x+52, 595, questFace14, color.RGBA{205, 243, 240, 255})
}

func (g *game) targetPosition() (int, int) {
	switch {
	case g.quest == 0:
		return 2*tile + 24, 74 + 9*tile + 24
	case g.quest == 1 || g.quest == 2:
		return 8*tile + 24, 74 + 2*tile + 24
	case g.quest == 3:
		return 5*tile + 24, 74 + 6*tile + 24
	default:
		return tile + 24, 74 + 10*tile + 24
	}
}

func partyAttackDamage(enemy, turn int, companion bool) int {
	damage := 10
	if companion {
		damage += 5
	}
	if enemy == 1 && turn%3 == 0 {
		damage /= 2
	}
	return damage
}

func enemyAttackDamage(enemy, intent int, defend bool) int {
	damage := []int{7, 11, 15}[enemy]
	if intent == 1 {
		damage += 6
	} else if intent == 2 {
		damage -= 3
	}
	if defend {
		damage = (damage + 1) / 2
	}
	return damage
}

func enemyIntentName(intent int) string {
	return []string{"NORMAL ATTACK", "HEAVY ATTACK", "QUICK ATTACK"}[intent]
}

func battleActionPulse(actionTick int) float64 {
	if actionTick <= 0 || actionTick >= 36 {
		return 0
	}
	return 1 - math.Abs(float64(actionTick-18))/18
}

func passableTile(x, y int) bool {
	return y >= 0 && y < len(worldRoute) && x >= 0 && x < len(worldRoute[y]) && worldRoute[y][x] == '.'
}

func press() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Quest — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
