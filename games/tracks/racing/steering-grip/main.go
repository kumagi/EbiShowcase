package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const W, H = 480, 720

type game struct{ x, y, a, speed float64 }

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.speed += .05
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.a -= .02 * (.3 + g.speed/5)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.a += .02 * (.3 + g.speed/5)
	}
	g.speed *= .99
	g.speed = math.Min(5, g.speed)
	g.x += math.Sin(g.a) * g.speed
	g.y -= math.Cos(g.a) * g.speed
	if g.x < 20 {
		g.x = 460
	}
	if g.x > 460 {
		g.x = 20
	}
	if g.y < 100 {
		g.y = 620
	}
	if g.y > 620 {
		g.y = 100
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{27, 70, 72, 255})
	for i := 0; i < 8; i++ {
		vector.StrokeCircle(s, 240, 360, float32(60+i*28), 1, color.RGBA{255, 255, 255, 35}, true)
	}
	op := &ebiten.DrawImageOptions{}
	img := ebiten.NewImage(24, 38)
	img.Fill(color.RGBA{235, 91, 76, 255})
	op.GeoM.Translate(-12, -19)
	op.GeoM.Rotate(g.a)
	op.GeoM.Translate(g.x, g.y)
	s.DrawImage(img, op)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SPEED %.1f  TURN RATE %.3f", g.speed, .02*(.3+g.speed/5)), 145, 50)
	ebitenutil.DebugPrintAt(s, "UP: GAS   LEFT/RIGHT: STEER", 125, 670)
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Steering and Grip")
	if err := ebiten.RunGame(&game{x: 240, y: 560}); err != nil {
		panic(err)
	}
}
