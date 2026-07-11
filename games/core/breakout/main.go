package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type brick struct {
	x, y  float64
	alive bool
	hue   int
}
type game struct {
	paddle, bx, by, vx, vy float64
	bricks                 []brick
	score, lives           int
}

func newGame() *game {
	g := &game{paddle: width / 2, lives: 3}
	for row := 0; row < 6; row++ {
		for col := 0; col < 8; col++ {
			g.bricks = append(g.bricks, brick{x: 22 + float64(col)*56, y: 90 + float64(row)*32, alive: true, hue: row})
		}
	}
	g.serve()
	return g
}
func (g *game) serve() { g.bx, g.by, g.vx, g.vy = width/2, 590, 3.4, -4.5 }
func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.paddle -= 6
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.paddle += 6
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		g.paddle = float64(x)
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		g.paddle = float64(x)
	}
	g.paddle = math.Max(55, math.Min(width-55, g.paddle))
	g.bx += g.vx
	g.by += g.vy
	if g.bx < 12 || g.bx > width-12 {
		g.vx = -g.vx
		g.bx = math.Max(12, math.Min(width-12, g.bx))
	}
	if g.by < 52 {
		g.vy = math.Abs(g.vy)
	}
	if g.vy > 0 && g.by > 630 && g.by < 650 && math.Abs(g.bx-g.paddle) < 64 {
		g.by = 630
		g.vy = -math.Abs(g.vy)
		g.vx += (g.bx - g.paddle) * .06
	}
	for i := range g.bricks {
		b := &g.bricks[i]
		if !b.alive {
			continue
		}
		if g.bx+10 > b.x && g.bx-10 < b.x+50 && g.by+10 > b.y && g.by-10 < b.y+24 {
			b.alive = false
			g.score += 10
			g.vy = -g.vy
			break
		}
	}
	if g.by > height+15 {
		g.lives--
		if g.lives <= 0 {
			*g = *newGame()
		} else {
			g.serve()
		}
	}
	alive := 0
	for _, b := range g.bricks {
		if b.alive {
			alive++
		}
	}
	if alive == 0 {
		*g = *newGame()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{7, 20, 38, 255})
	cols := []color.RGBA{{255, 105, 79, 255}, {255, 145, 72, 255}, {255, 210, 72, 255}, {45, 226, 194, 255}, {71, 161, 220, 255}, {151, 113, 230, 255}}
	for _, b := range g.bricks {
		if b.alive {
			vector.DrawFilledRect(s, float32(b.x), float32(b.y), 50, 24, cols[b.hue], false)
		}
	}
	vector.DrawFilledRect(s, float32(g.paddle-58), 640, 116, 14, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledCircle(s, float32(g.bx), float32(g.by), 10, color.White, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %04d   LIVES %d", g.score, g.lives), 160, 25)
	ebitenutil.DebugPrintAt(s, "ARROWS / DRAG", 184, 690)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Breakout — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
