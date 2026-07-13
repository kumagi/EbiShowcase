package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type card struct {
	name                string
	cost, damage, block int
	c                   color.RGBA
}

var library = []card{{"Strike", 1, 7, 0, color.RGBA{231, 87, 77, 255}}, {"Guard", 1, 0, 7, color.RGBA{72, 147, 232, 255}}, {"Heavy", 2, 15, 0, color.RGBA{160, 83, 210, 255}}}

type game struct {
	draw, hand, discard                        []card
	rng                                        *rand.Rand
	hp, enemyHP, energy, block, turn, shuffles int
	message                                    string
	clear, over                                bool
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(5103)), hp: 45, enemyHP: 100, turn: 1}
	for i := 0; i < 5; i++ {
		g.draw = append(g.draw, library[0], library[1])
	}
	g.draw = append(g.draw, library[2], library[2])
	g.shuffleDraw()
	g.startTurn()
	return g
}
func (g *game) shuffleDraw() {
	g.rng.Shuffle(len(g.draw), func(i, j int) { g.draw[i], g.draw[j] = g.draw[j], g.draw[i] })
}
func (g *game) startTurn() {
	g.energy = 3
	g.block = 0
	for len(g.hand) < 5 {
		if len(g.draw) == 0 {
			if len(g.discard) == 0 {
				break
			}
			g.draw = append(g.draw, g.discard...)
			// Keep the reusable backing array; only the cards' logical count is reset.
			g.discard = g.discard[:0]
			g.shuffleDraw()
			g.shuffles++
		}
		last := len(g.draw) - 1
		g.hand = append(g.hand, g.draw[last])
		g.draw = g.draw[:last]
	}
	g.message = "Draw five. Play cards or end the turn."
}
func (g *game) Update() error {
	if g.clear || g.over {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	choice := -1
	keys := []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5}
	for i, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			choice = i
		}
	}
	end := inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if x, y, ok := press(); ok {
		if y > 645 {
			end = true
		} else if y > 470 && len(g.hand) > 0 {
			choice = min(len(g.hand)-1, x/(width/len(g.hand)))
		}
	}
	if choice >= 0 && choice < len(g.hand) {
		c := g.hand[choice]
		if g.energy < c.cost {
			g.message = "Not enough energy for " + c.name
		} else {
			g.energy -= c.cost
			g.enemyHP -= c.damage
			g.block += c.block
			g.discard = append(g.discard, c)
			g.hand = append(g.hand[:choice], g.hand[choice+1:]...)
			g.message = fmt.Sprintf("Played %s. It moved to discard.", c.name)
			if g.enemyHP <= 0 {
				g.clear = true
			}
		}
	}
	if end && !g.clear {
		g.discard = append(g.discard, g.hand...)
		g.hand = nil
		damage := max(0, 9-g.block)
		g.hp -= damage
		g.turn++
		if g.hp <= 0 {
			g.over = true
		} else {
			g.startTurn()
			g.message = fmt.Sprintf("Enemy dealt %d. Drew a new hand.", damage)
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	vector.DrawFilledCircle(s, 240, 145, 58, color.RGBA{215, 103, 76, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ARCHIVE GOLEM HP %03d/100", max(0, g.enemyHP)), 160, 60)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TURN %d   HP %02d/45   ENERGY %d/3   BLOCK %d", g.turn, max(0, g.hp), g.energy, g.block), 105, 300)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DRAW %02d   HAND %d   DISCARD %02d   SHUFFLES %d", len(g.draw), len(g.hand), len(g.discard), g.shuffles), 90, 335)
	ebitenutil.DebugPrintAt(s, g.message, 55, 380)
	if len(g.hand) > 0 {
		w := float32(width / len(g.hand))
		for i, c := range g.hand {
			x := float32(i)*w + 3
			shade := c.c
			if g.energy < c.cost {
				shade = color.RGBA{75, 80, 96, 255}
			}
			vector.DrawFilledRect(s, x, 455, w-6, 170, shade, false)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s\nCOST %d\nDMG %d\nBLK %d", i+1, c.name, c.cost, c.damage, c.block), int(x)+10, 480)
		}
	}
	vector.DrawFilledRect(s, 55, 642, 370, 55, color.RGBA{240, 177, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, "END TURN [E]", 185, 665)
	if g.clear {
		overlay(s, "DECK CYCLE COMPLETE!\n\nTAP / SPACE TO RESTART")
	} else if g.over {
		overlay(s, "THE DECK FAILED!\n\nTAP / SPACE TO RETRY")
	}
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
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 115, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Draw Hand Discard — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
