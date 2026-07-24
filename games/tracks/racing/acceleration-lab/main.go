package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
)

const W, H = 480, 720

type game struct {
	x, speed float64
	tick     int
	trail    vfxmotion.Trail
	fx       vfxfx.System
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
		g.trail.Clear()
		g.fx.Shockwave(g.x, 374, .55, color.RGBA{255, 220, 90, 255}, color.RGBA{80, 235, 210, 255})
	}
	g.trail.Push(vfxmotion.Point{X: g.x, Y: 374})
	g.fx.Update()
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{34, 99, 66, 255})
	vector.DrawFilledRect(s, 0, 300, 480, 150, color.RGBA{65, 70, 82, 255}, false)
	points := g.trail.Points()
	for i, point := range points {
		alpha := uint8(25 + 100*i/max(1, len(points)))
		radius := float32(3 + 5*i/max(1, len(points)))
		vector.DrawFilledCircle(s, float32(point.X), float32(point.Y), radius, color.RGBA{255, 210, 90, alpha}, true)
	}
	bob := math.Sin(float64(g.tick)*(.12+g.speed*.08)) * math.Min(3, g.speed*.55)
	stretch := 1 + g.speed*.035
	carW, carH := 30*stretch, 55/stretch
	vector.DrawFilledRect(s, float32(g.x-carW/2), float32(372-carH/2+bob), float32(carW), float32(carH), color.RGBA{235, 91, 76, 255}, false)
	if g.speed > 3.5 {
		for i := 0; i < 3; i++ {
			y := float32(342 + i*18)
			vector.StrokeLine(s, float32(g.x-28-float64(i*7)), y, float32(g.x-55-float64(i*10)), y, 2, color.RGBA{155, 235, 220, 150}, true)
		}
	}
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SPEED %.2f   ACCEL +0.06   DRAG x0.985", g.speed), 95, 60)
	ebitenutil.DebugPrintAt(s, "HOLD SPACE / SCREEN: ACCELERATE", 105, 650)
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Acceleration Lab")
	if err := ebiten.RunGame(&game{x: 40, trail: vfxmotion.NewTrail(18)}); err != nil {
		panic(err)
	}
}
