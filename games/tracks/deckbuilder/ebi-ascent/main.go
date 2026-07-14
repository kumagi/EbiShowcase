package main

import (
	"fmt"
	"image/color"
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
	phaseBattle = iota
	phaseReward
	phaseRoute
)

type card struct {
	name          string
	damage, block int
	color         color.RGBA
}
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
}
type spark struct{ x, y, vx, vy, life float64 }

func newGame() *game {
	g := &game{hp: 46, phase: phaseBattle, message: "Play cards, then end the turn."}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: width / 2, Y: height / 2}, ViewW: width, ViewH: height}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{253, 200, 70, 255})
	g.deck = []card{{"JAB", 7, 0, color.RGBA{215, 92, 77, 255}}, {"SHELL", 0, 8, color.RGBA{71, 147, 224, 255}}, {"JAB", 7, 0, color.RGBA{215, 92, 77, 255}}}
	g.startBattle()
	return g
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
				choice = min(len(g.hand)-1, x/(width/len(g.hand)))
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
func (g *game) drawBattle(s *ebiten.Image) {
	e := encounters[g.floor]
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 5
	}
	sprites := []string{"swarm", "slug", "boss-crab", "species-3", "king-crab"}
	sprite := sprites[g.floor]
	size := 110 + float64(g.floor)*10
	if g.flash > 0 {
		trackatlas.DrawTinted(s, sprite, 240+ox, 150, size, 1, 1, .3, 1)
	} else {
		trackatlas.DrawCentered(s, sprite, 240+ox, 150+math.Sin(float64(g.tick)*.13)*3, size)
	}
	intent := e.attack + []int{0, 4, -2}[g.turn%3]
	intentName := []string{"STRIKE", "HEAVY", "FEINT"}[g.turn%3]
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s HP %02d/%02d NEXT %s %d", e.name, max(0, g.enemyHP), e.hp, intentName, intent), 95, 75)
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENERGY %d/3   BLOCK %d", g.energy, g.block), 160, 340)
	ebitenutil.DebugPrintAt(s, g.message, 45, 390)
	if len(g.hand) > 0 {
		w := float32(width / len(g.hand))
		for i, c := range g.hand {
			x := float32(i)*w + 3
			y := float32(465)
			if i%2 == 0 {
				y -= float32(math.Sin(float64(g.tick)*.08)) * 3
			}
			vector.DrawFilledRect(s, x, y, w-6, 155, c.color, false)
			trackatlas.DrawCentered(s, cardSprite(c), float64(x+w/2-3), 505, 40)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s\nDMG %d\nBLK %d", i+1, c.name, c.damage, c.block), int(x)+8, 550)
		}
	}
	vector.DrawFilledRect(s, 55, 642, 370, 55, color.RGBA{240, 177, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, "END TURN [E]", 185, 665)
}
func (g *game) drawReward(s *ebiten.Image) {
	ebitenutil.DebugPrintAt(s, "CHOOSE ONE CARD REWARD", 145, 320)
	for i, c := range rewards {
		x := float32(i*160 + 5)
		vector.DrawFilledRect(s, x, 390, 150, 200, c.color, false)
		trackatlas.DrawCentered(s, cardSprite(c), float64(x+75), 440, 48)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s\n\nDMG %d\nBLOCK %d", i+1, c.name, c.damage, c.block), int(x)+10, 500)
	}
	ebitenutil.DebugPrintAt(s, g.message, 70, 630)
}
func (g *game) drawRoute(s *ebiten.Image) {
	ebitenutil.DebugPrintAt(s, "CHOOSE THE NEXT ROUTE", 145, 280)
	vector.DrawFilledRect(s, 20, 400, 210, 170, color.RGBA{66, 158, 129, 255}, false)
	vector.DrawFilledRect(s, 250, 400, 210, 170, color.RGBA{178, 113, 58, 255}, false)
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
