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

type game struct {
	x, speed float64
	tick     int
}

func (g *game) Update() error {
	g.tick++
	gas := ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || len(ebiten.AppendTouchIDs(nil)) > 0 {
		gas = true
	}
	if gas {
		g.speed += .06
	}
	g.speed *= .985
	g.speed = math.Min(6, g.speed)
	g.x += g.speed
	if g.x > 440 {
		g.x = 40
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{34, 99, 66, 255})
	vector.DrawFilledRect(s, 0, 300, 480, 150, color.RGBA{65, 70, 82, 255}, false)
	vector.DrawFilledRect(s, float32(g.x)-15, 345, 30, 55, color.RGBA{235, 91, 76, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SPEED %.2f   ACCEL +0.06   DRAG x0.985", g.speed), 95, 60)
	ebitenutil.DebugPrintAt(s, "HOLD SPACE / SCREEN: ACCELERATE", 105, 650)
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Acceleration Lab")
	if err := ebiten.RunGame(&game{x: 40}); err != nil {
		panic(err)
	}
}
