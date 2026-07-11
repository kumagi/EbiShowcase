package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float64 }
type game struct {
	p               rect
	vx, vy          float64
	blocks          []rect
	grounded, clear bool
}

func newGame() *game {
	return &game{p: rect{36, 590, 28, 38}, blocks: []rect{{0, 650, 480, 70}, {90, 550, 110, 22}, {255, 470, 100, 22}, {365, 380, 90, 22}, {210, 300, 100, 22}, {35, 220, 105, 22}}}
}
func (g *game) Update() error {
	if g.clear {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	left := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				left = true
			} else {
				right = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			jump = true
		}
	}
	if left {
		g.vx -= .7
	}
	if right {
		g.vx += .7
	}
	if !left && !right {
		g.vx *= .76
	}
	g.vx = clamp(g.vx, -5.5, 5.5)
	if jump && g.grounded {
		g.vy = -12
		g.grounded = false
	}
	g.vy = math.Min(g.vy+.65, 14)
	g.p.x += g.vx
	for _, b := range g.blocks {
		if overlap(g.p, b) {
			if g.vx > 0 {
				g.p.x = b.x - g.p.w
			} else if g.vx < 0 {
				g.p.x = b.x + b.w
			}
			g.vx = 0
		}
	}
	g.p.x = clamp(g.p.x, 0, width-g.p.w)
	g.p.y += g.vy
	g.grounded = false
	for _, b := range g.blocks {
		if overlap(g.p, b) {
			if g.vy > 0 {
				g.p.y = b.y - g.p.h
				g.grounded = true
			} else {
				g.p.y = b.y + b.h
			}
			g.vy = 0
		}
	}
	if g.p.y > height {
		*g = *newGame()
	}
	if g.p.y < 210 && g.p.x < 145 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{111, 194, 231, 255})
	for _, b := range g.blocks {
		vector.DrawFilledRect(s, float32(b.x), float32(b.y), float32(b.w), float32(b.h), color.RGBA{58, 102, 69, 255}, false)
		vector.DrawFilledRect(s, float32(b.x), float32(b.y), float32(b.w), 7, color.RGBA{108, 215, 88, 255}, false)
	}
	vector.DrawFilledRect(s, 63, 175, 55, 35, color.RGBA{255, 210, 61, 255}, false)
	vector.DrawFilledRect(s, float32(g.p.x), float32(g.p.y), float32(g.p.w), float32(g.p.h), color.RGBA{239, 72, 88, 255}, false)
	vector.DrawFilledCircle(s, float32(g.p.x+20), float32(g.p.y+10), 3, color.White, false)
	ebitenutil.DebugPrintAt(s, "CLIMB TO THE GOLD FLAG", 145, 28)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH    JUMP: SPACE OR UPPER TOUCH", 50, 685)
	if g.clear {
		vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "YOU REACHED THE FLAG!\n\nTAP / SPACE TO PLAY AGAIN", 135, 330)
	}
}
func overlap(a, b rect) bool        { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Tiny Platformer — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
