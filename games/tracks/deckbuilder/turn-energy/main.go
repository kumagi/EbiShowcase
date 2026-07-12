package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720
const (
	playerPhase = iota
	enemyPhase
)

type card struct {
	name                string
	cost, damage, block int
	c                   color.RGBA
}

var hand = []card{{"Jab", 1, 7, 0, color.RGBA{230, 91, 78, 255}}, {"Big Wave", 2, 16, 0, color.RGBA{170, 85, 220, 255}}, {"Shell", 1, 0, 9, color.RGBA{74, 151, 235, 255}}}

type game struct {
	hp, enemyHP, energy, block, phase, timer, turn int
	message                                        string
	clear, over                                    bool
}

func newGame() *game {
	return &game{hp: 45, enemyHP: 75, energy: 3, phase: playerPhase, turn: 1, message: "PLAYER TURN: spend energy or end the turn."}
}
func (g *game) Update() error {
	if g.clear || g.over {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	if g.phase == enemyPhase {
		g.timer--
		if g.timer <= 0 {
			damage := max(0, 11-g.block)
			g.hp -= damage
			g.message = fmt.Sprintf("Enemy deals %d after block. New turn!", damage)
			g.block = 0
			g.energy = 3
			g.turn++
			g.phase = playerPhase
			if g.hp <= 0 {
				g.over = true
			}
		}
		return nil
	}
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		choice = 2
	}
	end := inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if x, y, ok := press(); ok {
		if y > 630 {
			end = true
		} else if y > 430 {
			choice = min(2, x/(width/3))
		}
	}
	if choice >= 0 {
		c := hand[choice]
		if g.energy < c.cost {
			g.message = "Not enough energy for " + c.name
		} else {
			g.energy -= c.cost
			g.enemyHP -= c.damage
			g.block += c.block
			g.message = fmt.Sprintf("Played %s: damage %d, block %d.", c.name, c.damage, c.block)
			if g.enemyHP <= 0 {
				g.clear = true
			}
		}
	}
	if end && !g.clear {
		g.phase = enemyPhase
		g.timer = 35
		g.message = "ENEMY TURN: preparing 11 damage."
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	vector.DrawFilledCircle(s, 240, 155, 62, color.RGBA{222, 104, 79, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("CLOCKWORK CRAB HP %02d/75", max(0, g.enemyHP)), 160, 70)
	vector.DrawFilledCircle(s, 240, 320, 32, color.RGBA{45, 225, 194, 255}, false)
	phase := "PLAYER"
	if g.phase == enemyPhase {
		phase = "ENEMY"
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s PHASE   TURN %d   HP %02d/45   ENERGY %d/3   BLOCK %d", phase, g.turn, max(0, g.hp), g.energy, g.block), 55, 370)
	ebitenutil.DebugPrintAt(s, g.message, 55, 410)
	for i, c := range hand {
		x := float32(10 + i*157)
		shade := c.c
		if g.energy < c.cost {
			shade = color.RGBA{75, 80, 96, 255}
		}
		vector.DrawFilledRect(s, x, 455, 147, 160, shade, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d  %s\n\nCOST %d\nDMG %d  BLK %d", i+1, c.name, c.cost, c.damage, c.block), int(x)+18, 480)
	}
	vector.DrawFilledRect(s, 55, 630, 370, 60, color.RGBA{240, 177, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, "END TURN [E / ENTER]", 160, 655)
	if g.clear {
		overlay(s, "BATTLE WON!\n\nTAP / SPACE TO RESTART")
	} else if g.over {
		overlay(s, "THE PARTY FELL!\n\nTAP / SPACE TO RETRY")
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
	ebitenutil.DebugPrintAt(s, msg, 140, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Energy Turns — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
