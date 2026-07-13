package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const width, height, tile = 480, 720, 48
const saveKey = "ebiShowcaseQuest."

type game struct {
	x, y, quest, hp, enemyHP, enemyMax, enemy, scene int
	turn, tick, shake, flash                         int
	companion, clear                                 bool
	defend                                           bool
	message                                          string
}

func newGame() *game {
	g := &game{x: 1, y: 10, hp: 60, message: "Meet Momo in the southwest village."}
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
		g.x = max(0, min(9, g.x+dx))
		g.y = max(0, min(11, g.y+dy))
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
	g.message = msg
}
func (g *game) resetSave() {
	for _, k := range []string{"quest", "x", "y", "hp"} {
		storageRemove(saveKey + k)
	}
}
func (g *game) battle() error {
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
	if choice == 0 {
		d := 10
		if g.companion {
			d += 5
		}
		if g.enemy == 1 && g.turn%3 == 0 {
			d /= 2
		}
		g.enemyHP -= d
		g.flash = 7
		g.shake = 4
		g.message = fmt.Sprintf("Ebi and Momo deal %d damage!", d)
	} else if choice == 1 {
		g.defend = true
		g.message = "Ebi braces for the next attack!"
	} else {
		g.hp = min(60, g.hp+15)
		g.message = "Momo's song restores 15 HP!"
	}
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
		return nil
	}
	damage := []int{7, 11, 15}[g.enemy]
	intent := g.turn % 3
	if intent == 1 {
		damage += 6
	}
	if intent == 2 {
		damage -= 3
	}
	if g.defend {
		damage = (damage + 1) / 2
		g.defend = false
	}
	g.hp -= damage
	g.turn++
	g.shake = 7
	if intent == 1 {
		g.message = fmt.Sprintf("Heavy attack! Party takes %d.", damage)
	} else if intent == 2 {
		g.message = fmt.Sprintf("Quick attack! Party takes %d.", damage)
	} else {
		g.message = fmt.Sprintf("Enemy attacks for %d.", damage)
	}
	if g.hp <= 0 {
		g.hp = 60
		g.scene = 0
		g.x, g.y = 2, 9
		g.message = "The party escaped and recovered in the village."
	}
	g.save()
	return nil
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
	s.Fill(color.RGBA{12, 28, 40, 255})
	oy := 74
	for y := 0; y < 12; y++ {
		for x := 0; x < 10; x++ {
			t := "tile-grass"
			if x > 5 {
				t = "tile-cobble"
			}
			trackatlas.Draw(s, t, float64(x*tile), float64(oy+y*tile), tile)
		}
	}
	trackatlas.DrawCentered(s, "npc", 2*tile+24, float64(oy+9*tile+24), 32)
	if g.quest == 1 {
		vector.DrawFilledCircle(s, 8*tile+24, float32(oy+2*tile+24), 11, color.RGBA{72, 205, 255, 255}, false)
	}
	vector.DrawFilledRect(s, 8*tile+7, float32(oy+tile+7), 34, 34, color.RGBA{55, 45, 75, 255}, false)
	if g.quest == 3 {
		vector.StrokeCircle(s, 5*tile+24, float32(oy+6*tile+24), 18, 3, color.RGBA{255, 90, 90, 220}, true)
	}
	hero.DrawCentered(s, float64(g.x*tile+24), float64(oy+g.y*tile+24), 34)
	if g.companion {
		trackatlas.DrawCentered(s, "ally", float64(g.x*tile+10), float64(oy+g.y*tile+38), 22)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI QUEST  HP %02d/60  QUEST %d/5", g.hp, g.quest), 120, 25)
	ebitenutil.DebugPrintAt(s, g.message, 60, 50)
	ebitenutil.DebugPrintAt(s, "AUTOSAVED IN THIS BROWSER   R: DELETE SAVE", 85, 690)
	if g.clear {
		overlay(s, "EBI QUEST COMPLETE!\n\nTAP / SPACE TO PLAY AGAIN")
	}
}
func (g *game) drawBattle(s *ebiten.Image) {
	colors := []color.RGBA{{20, 48, 58, 255}, {44, 36, 67, 255}, {52, 24, 35, 255}}
	s.Fill(colors[g.enemy])
	ox := 0.0
	if g.shake > 0 {
		ox = float64((g.tick%3)-1) * 5
	}
	sprites := []string{"species-1", "knight", "king-crab"}
	sprite := sprites[g.enemy]
	if g.enemy == 1 {
		sprite = "boss-crab"
	}
	size := []float64{100, 130, 160}[g.enemy]
	if g.flash > 0 {
		trackatlas.DrawTinted(s, sprite, 240+ox, 210, size, 1, 1, .35, 1)
	} else {
		trackatlas.DrawCentered(s, sprite, 240+ox, 210+math.Sin(float64(g.tick)*.12)*3, size)
	}
	names := []string{"CRYSTAL SLIME", "TOWER KNIGHT", "SHADOW CRAB"}
	intent := []string{"NORMAL ATTACK", "HEAVY ATTACK — GUARD!", "QUICK ATTACK"}[g.turn%3]
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s HP %02d/%02d", names[g.enemy], max(0, g.enemyHP), g.enemyMax), 145, 105)
	ebitenutil.DebugPrintAt(s, "NEXT: "+intent, 150, 350)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PARTY HP %02d/60", g.hp), 185, 415)
	ebitenutil.DebugPrintAt(s, g.message, 60, 480)
	vector.DrawFilledRect(s, 12, 550, 145, 80, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 167, 550, 145, 80, color.RGBA{80, 145, 225, 255}, false)
	vector.DrawFilledRect(s, 322, 550, 145, 80, color.RGBA{255, 190, 75, 255}, false)
	ebitenutil.DebugPrintAt(s, "1 ATTACK", 50, 585)
	ebitenutil.DebugPrintAt(s, "2 GUARD", 207, 585)
	ebitenutil.DebugPrintAt(s, "3 SONG", 362, 585)
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
