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
	sw   = 480
	sh   = 720
	cols = 8
	rows = 12
	cell = 42
	bx   = 72
	by   = 92
)

type rule struct {
	name              string
	fall, goal, holes int
	bg                color.RGBA
}

var rules = []rule{{"CALM BAY", 34, 2, 0, color.RGBA{9, 28, 48, 255}}, {"REEF RUSH", 24, 3, 2, color.RGBA{30, 18, 54, 255}}, {"STORM DECK", 15, 4, 4, color.RGBA{7, 45, 50, 255}}}

type game struct {
	board                           [rows][cols]bool
	x, y, tick, stage, clears, best int
	won, lost                       bool
}

func newGame() *game { g := &game{}; g.setup(); return g }
func (g *game) setup() {
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = false
		}
	}
	for n := 0; n < rules[g.stage].holes; n++ {
		y := rows - 1 - n
		gap := (n*3 + 2) % (cols - 1)
		for x := 0; x < cols; x++ {
			g.board[y][x] = x != gap && x != gap+1
		}
	}
	g.x = cols/2 - 1
	g.y = 0
	g.tick = 0
}
func (g *game) Update() error {
	if g.won || g.lost {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			b := g.best
			*g = *newGame()
			g.best = b
		}
		return nil
	}
	left := inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	drop := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyDown)
	if px, py, ok := press(); ok && py > 610 {
		if px < 160 {
			left = true
		} else if px < 320 {
			right = true
		} else {
			drop = true
		}
	}
	if left && g.can(g.x-1, g.y) {
		g.x--
	}
	if right && g.can(g.x+1, g.y) {
		g.x++
	}
	if drop {
		for g.can(g.x, g.y+1) {
			g.y++
		}
		g.lock()
		return nil
	}
	g.tick++
	if g.tick >= rules[g.stage].fall {
		g.tick = 0
		if g.can(g.x, g.y+1) {
			g.y++
		} else {
			g.lock()
		}
	}
	return nil
}
func (g *game) can(x, y int) bool {
	for i := 0; i < 2; i++ {
		xx := x + i
		if xx < 0 || xx >= cols || y < 0 || y >= rows || g.board[y][xx] {
			return false
		}
	}
	return true
}
func (g *game) lock() {
	for i := 0; i < 2; i++ {
		if g.y < 0 {
			g.lost = true
			return
		}
		g.board[g.y][g.x+i] = true
	}
	n := 0
	w := rows - 1
	for r := rows - 1; r >= 0; r-- {
		full := true
		for x := 0; x < cols; x++ {
			full = full && g.board[r][x]
		}
		if full {
			n++
			continue
		}
		g.board[w] = g.board[r]
		w--
	}
	for w >= 0 {
		g.board[w] = [cols]bool{}
		w--
	}
	g.clears += n
	if g.clears >= rules[g.stage].goal {
		g.best = max(g.best, g.stage+1)
		if g.stage == len(rules)-1 {
			g.won = true
			return
		}
		g.stage++
		g.clears = 0
		g.setup()
		return
	}
	g.x = cols/2 - 1
	g.y = 0
	if !g.can(g.x, g.y) {
		g.lost = true
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(rules[g.stage].bg)
	ebitenutil.DebugPrintAt(s, "RULE STAGE VOYAGE", 166, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE %d/3  %s  ROWS %d/%d", g.stage+1, rules[g.stage].name, g.clears, rules[g.stage].goal), 105, 52)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DROP TIMER %d   STARTING GARBAGE %d", rules[g.stage].fall, rules[g.stage].holes), 108, 72)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px := float32(bx + x*cell)
			py := float32(by + y*cell)
			vector.StrokeRect(s, px, py, cell-2, cell-2, 1, color.RGBA{60, 82, 108, 255}, false)
			if g.board[y][x] {
				vector.DrawFilledRect(s, px+3, py+3, cell-8, cell-8, color.RGBA{83, 174, 190, 255}, false)
			}
		}
	}
	if !g.won && !g.lost {
		for i := 0; i < 2; i++ {
			vector.DrawFilledRect(s, float32(bx+(g.x+i)*cell+3), float32(by+g.y*cell+3), cell-8, cell-8, color.RGBA{244, 177, 62, 255}, false)
		}
	}
	for i, t := range []string{"LEFT", "RIGHT", "DROP"} {
		vector.DrawFilledRect(s, float32(i*160+5), 620, 150, 64, color.RGBA{52, 88, 128, 255}, false)
		ebitenutil.DebugPrintAt(s, t, i*160+58, 648)
	}
	if g.won {
		over(s, "3 STAGES CLEAR!\nTAP / ENTER TO REPLAY")
	} else if g.lost {
		over(s, "BLOCKED AT THE TOP\nTAP / ENTER TO RETRY")
	}
}
func over(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 48, 285, 384, 120, color.RGBA{5, 13, 30, 245}, false)
	ebitenutil.DebugPrintAt(s, t, 135, 330)
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
func (g *game) Layout(int, int) (int, int) { return sw, sh }
func main() {
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowTitle("Rule Stage Voyage — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
