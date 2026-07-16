package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const width, height = 480, 720

const (
	handCardWidth = 142
	handCardGap   = 10
)
const (
	phaseBattle = iota
	phaseReward
	phaseRoute
)

type card struct {
	name          string
	damage, block int
	color         color.RGBA
}

// These original, Apache-2.0 assets are generated specifically for Ebi
// Ascent. They are decoded when a run is created, never lazily from Draw.
//
//go:embed assets/deckbuilder-abyss-background.png
var backgroundPNG []byte

//go:embed assets/deckbuilder-abyss-king.png
var bossPNG []byte

//go:embed assets/deckbuilder-card-atlas.png
var cardAtlasPNG []byte

//go:embed assets/deckbuilder-enemy-atlas.png
var enemyAtlasPNG []byte

type enemy struct {
	name       string
	hp, attack int
}

var encounters = []enemy{{"REEF SCOUT", 32, 6}, {"SHELL MAGE", 42, 8}, {"IRON CRAB", 54, 10}, {"TIDE KNIGHT", 66, 12}, {"ABYSS KING", 88, 14}}
var rewards = []card{{"PEARL STRIKE", 12, 0, color.RGBA{219, 91, 76, 255}}, {"TIDAL SHELL", 0, 12, color.RGBA{70, 148, 226, 255}}, {"EBI COMET", 16, 0, color.RGBA{169, 88, 211, 255}}}

type game struct {
	deck, hand                                      []card
	hp, enemyHP, energy, block, floor, phase, route int
	message                                         string
	clear, over                                     bool
	turn, tick, flash, shake, fx, score, best       int
	sparks                                          []spark
	audio                                           *audio.Context
	gate                                            audiolab.Gate
	pulse                                           *shaderlab.Pulse
	cam                                             cameralab.State
	badge                                           *ebiten.Image
	background, boss                                *ebiten.Image
	cardArt                                         [3]*ebiten.Image
	enemyArt                                        [4]*ebiten.Image
}
type spark struct{ x, y, vx, vy, life float64 }

func newGame() *game {
	g := &game{hp: 46, phase: phaseBattle, message: "Play cards, then end the turn."}
	g.loadGeneratedArt()
	g.audio = audiolab.Context()
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: width / 2, Y: height / 2}, ViewW: width, ViewH: height}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{253, 200, 70, 255})
	// The opening hand deliberately shows all three visual verbs immediately:
	// red impact attacks, blue shell defense, and the violet-gold finisher.
	g.deck = []card{{"JAB", 7, 0, color.RGBA{215, 92, 77, 255}}, {"SHELL", 0, 8, color.RGBA{71, 147, 224, 255}}, rewards[2]}
	g.startBattle()
	return g
}

func mustDecodePNG(data []byte) *ebiten.Image {
	source, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(source)
}

func (g *game) loadGeneratedArt() {
	g.background = mustDecodePNG(backgroundPNG)
	g.boss = mustDecodePNG(bossPNG)
	atlas := mustDecodePNG(cardAtlasPNG)
	panelW := atlas.Bounds().Dx() / len(g.cardArt)
	for i := range g.cardArt {
		panel := atlas.SubImage(image.Rect(i*panelW, 0, (i+1)*panelW, atlas.Bounds().Dy()))
		g.cardArt[i] = ebiten.NewImageFromImage(panel)
	}
	enemyAtlas := mustDecodePNG(enemyAtlasPNG)
	enemyW := enemyAtlas.Bounds().Dx() / len(g.enemyArt)
	for i := range g.enemyArt {
		panel := enemyAtlas.SubImage(image.Rect(i*enemyW, 0, (i+1)*enemyW, enemyAtlas.Bounds().Dy()))
		g.enemyArt[i] = ebiten.NewImageFromImage(panel)
	}
}
func (g *game) startBattle() {
	g.enemyHP = encounters[g.floor].hp
	g.energy = 3
	g.block = 0
	g.drawHand()
	g.phase = phaseBattle
}
func (g *game) drawHand() {
	g.hand = nil
	for i := 0; i < 5; i++ {
		g.hand = append(g.hand, g.deck[i%len(g.deck)])
	}
}
func (g *game) Update() error {
	g.tick++
	if g.flash > 0 {
		g.flash--
	}
	if g.shake > 0 {
		g.shake--
	}
	if g.fx > 0 {
		g.fx--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.clear || g.over {
		if restart() {
			best := g.best
			*g = *newGame()
			g.best = best
		}
		return nil
	}
	choice := -1
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5} {
		if inpututil.IsKeyJustPressed(k) {
			choice = i
		}
	}
	end := inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if x, y, ok := press(); ok {
		switch g.phase {
		case phaseBattle:
			if y > 640 {
				end = true
			} else if y > 470 {
				choice = handCardAt(x, len(g.hand))
			}
		case phaseReward:
			if y > 390 && y < 590 {
				choice = min(2, x/160)
			}
		case phaseRoute:
			if y > 400 {
				choice = min(1, x/240)
			}
		}
	}
	switch g.phase {
	case phaseBattle:
		g.updateBattle(choice, end)
	case phaseReward:
		if choice >= 0 && choice < 3 {
			g.deck = append(g.deck, rewards[choice])
			g.phase = phaseRoute
			g.message = "Reward added. Choose the next route."
		}
	case phaseRoute:
		if choice >= 0 && choice < 2 {
			g.route = choice
			if choice == 0 {
				g.hp = min(46, g.hp+9)
				g.message = "REST route restored 9 HP."
			} else {
				g.deck = append(g.deck, rewards[(g.floor+1)%3])
				g.message = "TREASURE route added a card."
			}
			g.floor++
			if g.floor >= len(encounters) {
				g.clear = true
				g.score = g.hp*20 + len(g.deck)*75
				if g.score > g.best {
					g.best = g.score
				}
			} else {
				g.startBattle()
			}
		}
	}
	return nil
}
func (g *game) updateBattle(choice int, end bool) {
	if choice >= 0 && choice < len(g.hand) && g.energy > 0 {
		g.play(660)
		c := g.hand[choice]
		g.energy--
		g.enemyHP -= c.damage
		g.block += c.block
		g.fx = 16
		if c.damage > 0 {
			g.flash = 7
			g.shake = 4
			g.burst(240, 165, 10)
		} else {
			g.burst(240, 360, 7)
		}
		g.hand = append(g.hand[:choice], g.hand[choice+1:]...)
		g.message = fmt.Sprintf("%s: damage %d, block %d.", c.name, c.damage, c.block)
		if g.enemyHP <= 0 {
			g.phase = phaseReward
			g.message = "Victory! Choose exactly one reward."
			return
		}
	}
	if end {
		g.play(220)
		e := encounters[g.floor]
		intent := e.attack + []int{0, 4, -2}[g.turn%3]
		taken := max(0, intent-g.block)
		g.hp -= taken
		g.block = 0
		g.energy = 3
		g.drawHand()
		g.turn++
		g.shake = 5
		g.message = fmt.Sprintf("%s dealt %d. New hand drawn.", e.name, taken)
		if g.hp <= 0 {
			g.over = true
		}
	}
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	p := g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .07))
	p.Play()
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * 6.283 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 25 + float64(i%9)})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	bgs := []color.RGBA{{18, 27, 45, 255}, {30, 43, 58, 255}, {45, 31, 55, 255}, {54, 39, 31, 255}, {45, 20, 32, 255}}
	s.Fill(bgs[min(g.floor, 4)])
	g.drawGeneratedBackground(s)
	vector.DrawFilledRect(s, 10, 10, 460, 48, color.RGBA{5, 11, 25, 220}, false)
	vector.StrokeRect(s, 10, 10, 460, 48, 2, color.RGBA{253, 200, 70, 130}, false)
	for i := 0; i < 5; i++ {
		c := color.RGBA{80, 91, 111, 180}
		if i <= g.floor {
			c = color.RGBA{253, 200, 70, 230}
		}
		vector.DrawFilledCircle(s, float32(182+i*30), 66, 5, c, true)
		if i < 4 {
			vector.StrokeLine(s, float32(187+i*30), 66, float32(207+i*30), 66, 2, c, false)
		}
	}
	g.drawHUD(s)
	g.drawEffectBadge(s)
	switch g.phase {
	case phaseBattle:
		g.drawBattle(s)
	case phaseReward:
		g.drawReward(s)
	case phaseRoute:
		g.drawRoute(s)
	}
	if g.clear {
		overlay(s, "ASCENT COMPLETE!\n\nTAP / SPACE FOR A NEW RUN")
	} else if g.over {
		overlay(s, "THE RUN ENDED!\n\nTAP / SPACE TO RETRY")
	}
}

func (g *game) drawGeneratedBackground(s *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// The embedded image already matches the logical canvas. A slight floor
	// tint differentiates the climb while preserving the authored lighting.
	tints := [...]struct{ r, gr, b float32 }{{.78, .90, 1}, {.72, 1, .93}, {.90, .76, 1}, {1, .83, .68}, {1, .68, .76}}
	t := tints[min(g.floor, len(tints)-1)]
	op.ColorScale.Scale(t.r, t.gr, t.b, .86)
	s.DrawImage(g.background, op)
	// Cards need a quiet landing zone but the illustrated stairs remain visible.
	vector.DrawFilledRect(s, 0, 350, width, 370, color.RGBA{3, 8, 20, 155}, false)
	vector.DrawFilledRect(s, 0, 0, width, 82, color.RGBA{3, 8, 20, 115}, false)
}
func (g *game) drawHUD(s *ebiten.Image) {
	label := fmt.Sprintf("EBI ASCENT  FLOOR %d/5  HP %02d/46  DECK %d  BEST %04d", g.floor+1, max(0, g.hp), len(g.deck), g.best)
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(38, 24)
		text.Draw(s, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(s, label, 28, 35)
}
func (g *game) drawEffectBadge(s *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.tick)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(width-38, 16)
	s.DrawImage(fx, op)
}
func cardSprite(c card) string {
	switch {
	case c.damage > 0 && c.block > 0:
		return "card-skill"
	case c.damage > 0:
		return "card-attack"
	default:
		return "card-block"
	}
}

func cardArtIndex(c card) int {
	if c.name == "EBI COMET" {
		return 2
	}
	if c.block > 0 {
		return 1
	}
	return 0
}

func handLayout(count int) (start, cardWidth, step float32) {
	if count <= 0 {
		return 0, 0, 0
	}
	cardWidth = handCardWidth
	step = cardWidth + handCardGap
	if count > 1 {
		step = min(step, (width-8-cardWidth)/float32(count-1))
	}
	total := cardWidth + step*float32(count-1)
	return (width - total) / 2, cardWidth, step
}

func handCardAt(x, count int) int {
	start, cardWidth, step := handLayout(count)
	// Later cards are drawn over earlier cards, so test them in reverse order
	// when a five-card hand overlaps.
	for i := count - 1; i >= 0; i-- {
		left := start + float32(i)*step
		if float32(x) >= left && float32(x) < left+cardWidth {
			return i
		}
	}
	return -1
}

func (g *game) drawCardArt(s *ebiten.Image, c card, x, y, w, h float64) {
	art := g.cardArt[cardArtIndex(c)]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(w/float64(art.Bounds().Dx()), h/float64(art.Bounds().Dy()))
	op.GeoM.Translate(x, y)
	s.DrawImage(art, op)
	// A bright lower edge connects the illustration to its rules text.
	vector.DrawFilledRect(s, float32(x), float32(y+h-3), float32(w), 3, c.color, false)
}

func (g *game) drawBoss(s *ebiten.Image, x, y, size float64, hit bool) {
	g.drawEnemyPortrait(s, g.boss, x, y, size, size, 1, hit)
}

func (g *game) drawEnemyPortrait(s *ebiten.Image, art *ebiten.Image, x, y, drawW, drawH, alpha float64, hit bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(drawW/float64(art.Bounds().Dx()), drawH/float64(art.Bounds().Dy()))
	op.GeoM.Translate(x-drawW/2, y-drawH/2)
	if hit {
		op.ColorScale.Scale(1.35, .42, .42, float32(alpha))
	} else {
		op.ColorScale.ScaleAlpha(float32(alpha))
	}
	s.DrawImage(art, op)
}

func (g *game) drawBattle(s *ebiten.Image) {
	e := encounters[g.floor]
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 5
	}
	// Shadow and danger aura are shared by all five generated portraits.
	vector.DrawFilledCircle(s, float32(240+ox), 282, 66, color.RGBA{1, 5, 14, 105}, true)
	vector.DrawFilledCircle(s, float32(240+ox), 190, float32(94+math.Sin(float64(g.tick)*.09)*6), color.RGBA{255, 93, 108, 18}, true)
	if g.floor == len(encounters)-1 {
		g.drawBoss(s, 240+ox, 205+math.Sin(float64(g.tick)*.05)*2, 286, g.flash > 0)
	} else {
		widths := [...]float64{190, 194, 208, 184}
		heights := [...]float64{278, 282, 278, 290}
		g.drawEnemyPortrait(s, g.enemyArt[g.floor], 240+ox, 205+math.Sin(float64(g.tick)*.08)*3, widths[g.floor], heights[g.floor], 1, g.flash > 0)
	}
	intent := e.attack + []int{0, 4, -2}[g.turn%3]
	intentName := []string{"STRIKE", "HEAVY", "FEINT"}[g.turn%3]
	vector.DrawFilledRect(s, 70, 76, 340, 42, color.RGBA{5, 12, 27, 205}, false)
	vector.DrawFilledRect(s, 95, 107, 290, 7, color.RGBA{55, 55, 75, 255}, false)
	vector.DrawFilledRect(s, 95, 107, float32(290*max(0, g.enemyHP)/e.hp), 7, color.RGBA{240, 79, 100, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  HP %02d/%02d  NEXT %s %d", e.name, max(0, g.enemyHP), e.hp, intentName, intent), 95, 88)
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENERGY %d/3   BLOCK %d", g.energy, g.block), 160, 340)
	ebitenutil.DebugPrintAt(s, g.message, 45, 390)
	if len(g.hand) > 0 {
		start, w, step := handLayout(len(g.hand))
		for i, c := range g.hand {
			x := start + float32(i)*step
			y := float32(465)
			if i%2 == 0 {
				y -= float32(math.Sin(float64(g.tick)*.08)) * 3
			}
			vector.DrawFilledRect(s, x+4, y+7, w, 155, color.RGBA{1, 5, 14, 150}, false)
			vector.DrawFilledRect(s, x, y, w, 155, color.RGBA{11, 18, 34, 248}, false)
			vector.StrokeRect(s, x, y, w, 155, 3, c.color, false)
			g.drawCardArt(s, c, float64(x+6), float64(y+7), float64(w-12), 74)
			vector.DrawFilledCircle(s, x+16, y+16, 12, color.RGBA{35, 223, 235, 255}, true)
			ebitenutil.DebugPrintAt(s, "1", int(x)+13, int(y)+11)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s\nDMG %d  BLK %d", i+1, c.name, c.damage, c.block), int(x)+8, int(y)+91)
		}
	}
	vector.DrawFilledRect(s, 55, 642, 370, 55, color.RGBA{240, 177, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, "END TURN [E]", 185, 665)
}
func (g *game) drawReward(s *ebiten.Image) {
	vector.DrawFilledCircle(s, 240, 400, 210, color.RGBA{255, 211, 90, 13}, true)
	ebitenutil.DebugPrintAt(s, "CHOOSE ONE CARD REWARD", 145, 320)
	for i, c := range rewards {
		x := float32(i*160 + 5)
		vector.DrawFilledRect(s, x+6, 398, 150, 200, color.RGBA{1, 5, 14, 160}, false)
		vector.DrawFilledRect(s, x, 390, 150, 200, color.RGBA{11, 18, 34, 250}, false)
		vector.StrokeRect(s, x, 390, 150, 200, 4, c.color, false)
		g.drawCardArt(s, c, float64(x+8), 400, 134, 105)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s\n\nDMG %d  BLOCK %d", i+1, c.name, c.damage, c.block), int(x)+10, 520)
	}
	ebitenutil.DebugPrintAt(s, g.message, 70, 630)
}
func (g *game) drawRoute(s *ebiten.Image) {
	for i := 0; i < 9; i++ {
		vector.DrawFilledCircle(s, float32(35+i*55), float32(350+(i%2)*230), 20, color.RGBA{71, 168, 145, 28}, true)
	}
	ebitenutil.DebugPrintAt(s, "CHOOSE THE NEXT ROUTE", 145, 280)
	vector.DrawFilledRect(s, 20, 400, 210, 170, color.RGBA{66, 158, 129, 255}, false)
	vector.DrawFilledRect(s, 250, 400, 210, 170, color.RGBA{178, 113, 58, 255}, false)
	vector.StrokeRect(s, 20, 400, 210, 170, 4, color.RGBA{180, 255, 218, 150}, false)
	vector.StrokeRect(s, 250, 400, 210, 170, 4, color.RGBA{255, 224, 156, 160}, false)
	trackatlas.DrawCentered(s, "route-rest", 125, 445, 56)
	trackatlas.DrawCentered(s, "route-treasure", 355, 445, 56)
	ebitenutil.DebugPrintAt(s, "1  REST\n\nHEAL 9 HP", 75, 490)
	ebitenutil.DebugPrintAt(s, "2  TREASURE\n\nADD A CARD", 285, 490)
	ebitenutil.DebugPrintAt(s, g.message, 65, 620)
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
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, m string) {
	vector.DrawFilledRect(s, 45, 275, 390, 165, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, m, 110, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Ascent — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
