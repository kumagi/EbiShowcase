package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
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

var encounters = []enemy{{"REEF SCOUT", 36, 7}, {"IRON CRAB", 52, 10}, {"ABYSS KING", 76, 14}}
var rewards = []card{{"PEARL STRIKE", 12, 0, color.RGBA{219, 91, 76, 255}}, {"TIDAL SHELL", 0, 12, color.RGBA{70, 148, 226, 255}}, {"EBI COMET", 16, 0, color.RGBA{169, 88, 211, 255}}}

type game struct {
	deck, hand                                      []card
	hp, enemyHP, energy, block, floor, phase, route int
	message                                         string
	clear, over                                     bool
}

func newGame() *game {
	g := &game{hp: 46, phase: phaseBattle, message: "Play cards, then end the turn."}
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
	if g.clear || g.over {
		if restart() {
			*g = *newGame()
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
			} else {
				g.startBattle()
			}
		}
	}
	return nil
}
func (g *game) updateBattle(choice int, end bool) {
	if choice >= 0 && choice < len(g.hand) && g.energy > 0 {
		c := g.hand[choice]
		g.energy--
		g.enemyHP -= c.damage
		g.block += c.block
		g.hand = append(g.hand[:choice], g.hand[choice+1:]...)
		g.message = fmt.Sprintf("%s: damage %d, block %d.", c.name, c.damage, c.block)
		if g.enemyHP <= 0 {
			g.phase = phaseReward
			g.message = "Victory! Choose exactly one reward."
			return
		}
	}
	if end {
		e := encounters[g.floor]
		taken := max(0, e.attack-g.block)
		g.hp -= taken
		g.block = 0
		g.energy = 3
		g.drawHand()
		g.message = fmt.Sprintf("%s dealt %d. New hand drawn.", e.name, taken)
		if g.hp <= 0 {
			g.over = true
		}
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI ASCENT   FLOOR %d/3   HP %02d/46   DECK %d", g.floor+1, max(0, g.hp), len(g.deck)), 85, 35)
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
	trackatlas.DrawCentered(s, "king-crab", 240, 150, 124)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  HP %02d/%02d  NEXT %d", e.name, max(0, g.enemyHP), e.hp, e.attack), 125, 75)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENERGY %d/3   BLOCK %d", g.energy, g.block), 160, 340)
	ebitenutil.DebugPrintAt(s, g.message, 45, 390)
	if len(g.hand) > 0 {
		w := float32(width / len(g.hand))
		for i, c := range g.hand {
			x := float32(i)*w + 3
			vector.DrawFilledRect(s, x, 465, w-6, 155, c.color, false)
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
