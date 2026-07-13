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

type game struct{ x, cam float64 }

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.x += 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.x -= 4
	}
	g.x = math.Max(0, math.Min(1600, g.x))
	g.cam += (g.x - 240 - g.cam) * .1
	g.cam = math.Max(0, math.Min(1120, g.cam))
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 22, 39, 255})
	for r := 0; r < 4; r++ {
		x := float32(float64(r*400) - g.cam)
		vector.DrawFilledRect(s, x, 160, 390, 400, []color.RGBA{{38, 65, 70, 255}, {43, 49, 79, 255}, {66, 45, 69, 255}, {70, 51, 39, 255}}[r], false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("ROOM %d", r+1), int(x)+160, 200)
	}
	vector.DrawFilledRect(s, float32(g.x-g.cam)-12, 500, 24, 36, color.RGBA{235, 91, 76, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("WORLD X %.0f  CAMERA %.0f  SCREEN X %.0f", g.x, g.cam, g.x-g.cam), 85, 50)
	ebitenutil.DebugPrintAt(s, "A/D: WALK THROUGH A WORLD WIDER THAN SCREEN", 75, 675)
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Camera Rooms")
	if err := ebiten.RunGame(&game{x: 80}); err != nil {
		panic(err)
	}
}
