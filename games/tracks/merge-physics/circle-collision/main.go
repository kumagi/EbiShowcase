package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type ball struct {
	x, y, vx, vy, r float64
	c               color.RGBA
}
type game struct {
	a, b            ball
	contacts, flash int
	clear           bool
}

func newGame() *game {
	return &game{a: ball{135, 340, 2.5, 1.3, 52, color.RGBA{45, 225, 194, 255}}, b: ball{345, 390, -2.1, -1.1, 72, color.RGBA{240, 91, 115, 255}}}
}
func (g *game) Update() error {
	if g.clear {
		if pressed() {
			*g = *newGame()
		}
		return nil
	}
	if x, y, ok := press(); ok {
		dx, dy := float64(x)-g.a.x, float64(y)-g.a.y
		l := math.Hypot(dx, dy)
		if l > 0 {
			g.a.vx = dx / l * 5
			g.a.vy = dy / l * 5
		}
	}
	move := func(p *ball) {
		p.x += p.vx
		p.y += p.vy
		if p.x-p.r < 24 {
			p.x = 24 + p.r
			p.vx = math.Abs(p.vx)
		}
		if p.x+p.r > 456 {
			p.x = 456 - p.r
			p.vx = -math.Abs(p.vx)
		}
		if p.y-p.r < 90 {
			p.y = 90 + p.r
			p.vy = math.Abs(p.vy)
		}
		if p.y+p.r > 660 {
			p.y = 660 - p.r
			p.vy = -math.Abs(p.vy)
		}
	}
	move(&g.a)
	move(&g.b)
	dx, dy := g.b.x-g.a.x, g.b.y-g.a.y
	distance := math.Hypot(dx, dy)
	minimum := g.a.r + g.b.r
	if distance < minimum && distance > 0 {
		nx, ny := dx/distance, dy/distance
		overlap := minimum - distance
		g.a.x -= nx * overlap / 2
		g.a.y -= ny * overlap / 2
		g.b.x += nx * overlap / 2
		g.b.y += ny * overlap / 2
		relative := (g.b.vx-g.a.vx)*nx + (g.b.vy-g.a.vy)*ny
		if relative < 0 {
			impulse := -relative * .9
			g.a.vx -= impulse * nx
			g.a.vy -= impulse * ny
			g.b.vx += impulse * nx
			g.b.vy += impulse * ny
			g.contacts++
			g.flash = 10
			if g.contacts >= 10 {
				g.clear = true
			}
		}
	}
	if g.flash > 0 {
		g.flash--
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 28, 44, 255})
	vector.StrokeRect(s, 20, 86, 440, 578, 5, color.RGBA{110, 128, 151, 255}, false)
	vector.DrawFilledCircle(s, float32(g.a.x), float32(g.a.y), float32(g.a.r), g.a.c, false)
	vector.DrawFilledCircle(s, float32(g.b.x), float32(g.b.y), float32(g.b.r), g.b.c, false)
	dx, dy := g.b.x-g.a.x, g.b.y-g.a.y
	d := math.Hypot(dx, dy)
	if d > 0 {
		nx, ny := dx/d, dy/d
		line := color.RGBA{100, 165, 255, 255}
		if g.flash > 0 {
			line = color.RGBA{255, 220, 70, 255}
		}
		vector.StrokeLine(s, float32(g.a.x), float32(g.a.y), float32(g.a.x+nx*g.a.r), float32(g.a.y+ny*g.a.r), 4, line, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("CONTACTS %02d/10   DISTANCE %.1f   RADII %.1f", g.contacts, d, g.a.r+g.b.r), 75, 30)
	ebitenutil.DebugPrintAt(s, "TAP TO LAUNCH THE GREEN CIRCLE", 125, 690)
	if g.clear {
		overlay(s, "TEN COLLISIONS SOLVED!\n\nTAP / SPACE TO RESET")
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return 345, 390, true
	}
	return 0, 0, false
}
func pressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Circle Collision — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
