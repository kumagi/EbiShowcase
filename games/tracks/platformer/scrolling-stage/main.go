package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float64 }
type game struct {
	p                          rect
	vx, vy, camera, checkpoint float64
	grounds                    []rect
	flags                      []float64
	grounded, clear            bool
	visible                    int
}

func newGame() *game {
	return &game{p: rect{35, 580, 28, 38}, checkpoint: 35, grounds: []rect{{0, 640, 350, 80}, {410, 580, 260, 140}, {735, 640, 310, 80}, {1110, 545, 260, 175}, {1430, 640, 320, 80}, {1810, 565, 250, 155}, {2120, 640, 380, 80}, {180, 505, 100, 20}, {520, 430, 105, 20}, {840, 500, 105, 20}, {1190, 395, 110, 20}, {1530, 480, 110, 20}, {1900, 410, 105, 20}, {2240, 500, 105, 20}}, flags: []float64{780, 1500, 2200}}
}
func (g *game) Update() error {
	if g.clear {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	l := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	r := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	j := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				l = true
			} else {
				r = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			j = true
		}
	}
	if l {
		g.vx -= .7
	}
	if r {
		g.vx += .7
	}
	if !l && !r {
		g.vx *= .8
	}
	g.vx = clamp(g.vx, -6, 6)
	if j && g.grounded {
		g.vy = -12.5
	}
	g.vy = math.Min(g.vy+.65, 14)
	g.p.x = clamp(g.p.x+g.vx, 0, 2470)
	old := g.p.y + g.p.h
	g.p.y += g.vy
	g.grounded = false
	for _, b := range g.grounds {
		if g.vy >= 0 && old <= b.y+3 && g.p.y+g.p.h >= b.y && g.p.x+g.p.w > b.x && g.p.x < b.x+b.w {
			g.p.y = b.y - g.p.h
			g.vy = 0
			g.grounded = true
		}
	}
	for _, f := range g.flags {
		if g.p.x > f {
			g.checkpoint = f
		}
	}
	if g.p.y > height {
		g.p.x = g.checkpoint
		g.p.y = 500
		g.vx, g.vy = 0, 0
	}
	if g.p.x > 2440 {
		g.clear = true
	}
	target := g.p.x - width*.4
	g.camera = clamp(g.camera+(target-g.camera)*.08, 0, 2020)
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{100, 188, 230, 255})
	g.visible = 0
	for _, b := range g.grounds {
		x := b.x - g.camera
		if x+b.w < 0 || x > width {
			continue
		}
		g.visible++
		vector.DrawFilledRect(s, float32(x), float32(b.y), float32(b.w), float32(b.h), color.RGBA{55, 101, 66, 255}, false)
		vector.DrawFilledRect(s, float32(x), float32(b.y), float32(b.w), 8, color.RGBA{108, 214, 87, 255}, false)
	}
	for _, f := range g.flags {
		x := f - g.camera
		if x < 0 || x > width {
			continue
		}
		c := color.RGBA{255, 210, 60, 255}
		if f <= g.checkpoint {
			c = color.RGBA{45, 224, 192, 255}
		}
		vector.DrawFilledRect(s, float32(x), 500, 5, 140, color.RGBA{240, 245, 250, 255}, false)
		vector.DrawFilledRect(s, float32(x+5), 505, 48, 30, c, false)
	}
	vector.DrawFilledRect(s, float32(g.p.x-g.camera), float32(g.p.y), float32(g.p.w), float32(g.p.h), color.RGBA{240, 73, 89, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("WORLD X %04d   CAMERA %04d   DRAWN %02d", int(g.p.x), int(g.camera), g.visible), 55, 22)
	ebitenutil.DebugPrintAt(s, "GREEN FLAGS ARE CHECKPOINTS", 145, 48)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH    JUMP: SPACE OR UPPER TOUCH", 50, 685)
	if g.clear {
		vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "HILLS CLEARED!\n\nTAP / SPACE TO RUN AGAIN", 145, 330)
	}
}
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Scrolling Hills — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
