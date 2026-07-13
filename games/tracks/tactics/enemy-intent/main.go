// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	W = 480
	H = 720
)

type game struct {
	target, guard, round, score, best, timer int
	message                                  string
}

func newGame() *game {
	return &game{target: 1, guard: -1, message: "Enemy intent points to the ally it will attack."}
}
func (g *game) Update() error {
	if g.timer > 0 {
		g.timer--
		if g.timer == 0 {
			if g.guard == g.target {
				g.score += 2
				g.message = "BLOCK! You read the intent."
			} else {
				g.score = max(0, g.score-1)
				g.message = "Hit! Guard the marked ally next time."
			}
			g.best = max(g.best, g.score)
			g.round++
			g.target = (g.round*2 + 1) % 3
			g.guard = -1
		}
		return nil
	}
	choice := -1
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if inpututil.IsKeyJustPressed(k) {
			choice = i
		}
	}
	if x, y, ok := press(); ok && y > 500 {
		choice = min(2, x/160)
	}
	if choice >= 0 {
		g.guard = choice
		g.timer = 40
		g.message = "Enemy advances..."
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{14, 23, 38, 255})
	ebitenutil.DebugPrintAt(s, "READ THE ENEMY INTENT", 151, 25)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ROUND %d  SCORE %d  BEST %d", g.round+1, g.score, g.best), 135, 55)
	ebitenutil.DebugPrintAt(s, g.message, 96, 82)
	enemyX := 240
	if g.timer > 0 {
		enemyX -= g.timer * 2
	}
	vector.DrawFilledCircle(s, float32(enemyX), 190, 30, color.RGBA{153, 72, 121, 255}, false)
	ebitenutil.DebugPrintAt(s, "ENEMY", enemyX-20, 185)
	for i := 0; i < 3; i++ {
		x := 80 + i*160
		c := color.RGBA{69, 139, 157, 255}
		if i == g.guard {
			c = color.RGBA{236, 177, 65, 255}
		}
		vector.DrawFilledCircle(s, float32(x), 390, 28, c, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("ALLY %d", i+1), x-22, 385)
		if i == g.target {
			vector.StrokeRect(s, float32(x-40), 340, 80, 90, 4, color.RGBA{245, 92, 82, 255}, false)
			ebitenutil.DebugPrintAt(s, "TARGET", x-22, 322)
		}
	}
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(s, float32(i*160+5), 520, 150, 90, color.RGBA{52, 86, 123, 255}, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("[%d] GUARD", i+1), i*160+45, 560)
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
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Enemy Intent")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
