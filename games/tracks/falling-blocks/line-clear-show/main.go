// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

type mote struct {
	x, y, vx, vy float64
	life         int
}
type game struct {
	board                                               [rows][cols]bool
	x, y, tick, phase, flash, shake, score, combo, best int
	motes                                               []mote
	rng                                                 *rand.Rand
}

func newGame() *game { g := &game{rng: rand.New(rand.NewSource(81))}; g.prepare(); return g }
func (g *game) prepare() {
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = false
		}
	}
	gap := g.rng.Intn(cols - 3)
	for x := 0; x < cols; x++ {
		if x < gap || x >= gap+4 {
			g.board[rows-1][x] = true
		}
	}
	g.x = gap
	g.y = 0
	g.tick = 0
	g.phase = 0
}
func (g *game) Update() error {
	for i := len(g.motes) - 1; i >= 0; i-- {
		p := &g.motes[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .09
		p.life--
		if p.life <= 0 {
			g.motes = append(g.motes[:i], g.motes[i+1:]...)
		}
	}
	if g.flash > 0 {
		g.flash--
	}
	if g.shake > 0 {
		g.shake--
	}
	if g.phase > 0 {
		g.phase--
		if g.phase == 0 {
			g.score += 100 + g.combo*40
			g.combo++
			if g.score > g.best {
				g.best = g.score
			}
			g.prepare()
		}
		return nil
	}
	left := inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	drop := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyDown)
	if px, py, ok := press(); ok && py > 610 {
		switch {
		case px < 160:
			left = true
		case px < 320:
			right = true
		default:
			drop = true
		}
	}
	if left && g.x > 0 {
		g.x--
	}
	if right && g.x < cols-4 {
		g.x++
	}
	if drop {
		for g.y < rows-1 {
			g.y++
		}
		g.land()
		return nil
	}
	g.tick++
	if g.tick >= 26 {
		g.tick = 0
		g.y++
		if g.y >= rows-1 {
			g.land()
		}
	}
	return nil
}
func (g *game) land() {
	for i := 0; i < 4; i++ {
		g.board[rows-1][g.x+i] = true
	}
	full := true
	for x := 0; x < cols; x++ {
		full = full && g.board[rows-1][x]
	}
	if full {
		g.phase = 38
		g.flash = 28
		g.shake = 12
		for i := 0; i < 36; i++ {
			a := g.rng.Float64() * math.Pi * 2
			s := 1 + g.rng.Float64()*4
			g.motes = append(g.motes, mote{bx + cols*cell/2, by + (rows-.5)*cell, math.Cos(a) * s, math.Sin(a)*s - 1, 24 + g.rng.Intn(20)})
		}
	} else {
		g.combo = 0
		g.prepare()
	}
}
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 20, 38, 255})
	ox := 0
	if g.shake > 0 {
		ox = (g.shake%3 - 1) * 4
	}
	ebitenutil.DebugPrintAt(screen, "LINE CLEAR SHOW", 174, 24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %04d  COMBO x%d  BEST %04d", g.score, g.combo+1, g.best), 112, 54)
	ebitenutil.DebugPrintAt(screen, "Fill the glowing gap with the 4-cell bar.", 94, 74)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px := float32(bx + x*cell + ox)
			py := float32(by + y*cell)
			vector.StrokeRect(screen, px, py, cell-2, cell-2, 1, color.RGBA{55, 76, 104, 255}, false)
			if g.board[y][x] {
				c := color.RGBA{70, 188, 210, 255}
				if g.phase > 0 && y == rows-1 {
					v := uint8(150 + 90*math.Abs(math.Sin(float64(g.phase)*.35)))
					c = color.RGBA{255, v, 90, 255}
				}
				vector.DrawFilledRect(screen, px+3, py+3, cell-8, cell-8, c, false)
			}
		}
	}
	if g.phase == 0 {
		for i := 0; i < 4; i++ {
			vector.DrawFilledRect(screen, float32(bx+(g.x+i)*cell+4+ox), float32(by+g.y*cell+4), cell-10, cell-10, color.RGBA{242, 180, 65, 255}, false)
		}
	}
	for _, p := range g.motes {
		vector.DrawFilledCircle(screen, float32(p.x+float64(ox)), float32(p.y), 3, color.RGBA{255, 216, 92, uint8(min(255, p.life*8))}, false)
	}
	for i, s := range []string{"LEFT", "RIGHT", "DROP"} {
		vector.DrawFilledRect(screen, float32(i*160+5), 620, 150, 64, color.RGBA{42, 83, 124, 255}, false)
		ebitenutil.DebugPrintAt(screen, s, i*160+58, 648)
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
func (g *game) Layout(int, int) (int, int) { return sw, sh }
func main() {
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowTitle("Line Clear Show — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
