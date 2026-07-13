package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const W, H = 480, 720

type game struct {
	x       float64
	dash    bool
	message string
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.x += 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.x -= 3
	}
	if !g.dash && g.x > 230 {
		g.x = 230
		g.message = "SEALED: collect the crest on the left"
	}
	if g.x < 90 {
		g.dash = true
		g.message = "DASH CREST: press X at the seal"
	}
	if g.dash && inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.x += 90
	}
	g.x = maxf(30, minf(450, g.x))
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{20, 31, 53, 255})
	vector.DrawFilledRect(s, 0, 540, 480, 180, color.RGBA{44, 72, 71, 255}, false)
	vector.DrawFilledRect(s, 250, 220, 25, 320, color.RGBA{90, 120, 180, 180}, false)
	vector.DrawFilledCircle(s, 70, 500, 15, color.RGBA{255, 205, 70, 255}, true)
	vector.DrawFilledRect(s, float32(g.x)-12, 500, 24, 40, color.RGBA{235, 91, 76, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DASH %v", g.dash), 205, 50)
	ebitenutil.DebugPrintAt(s, g.message, 90, 90)
	ebitenutil.DebugPrintAt(s, "A/D MOVE   X DASH", 165, 670)
}
func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Dash Ability Gate")
	if err := ebiten.RunGame(&game{x: 180, message: "The seal will not open yet."}); err != nil {
		panic(err)
	}
}
