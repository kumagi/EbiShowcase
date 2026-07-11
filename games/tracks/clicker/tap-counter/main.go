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

type game struct {
	count, pulse int
	clear        bool
}

func (g *game) Update() error {
	if g.clear {
		if pressed() {
			*g = game{}
		}
		return nil
	}
	if pressed() {
		g.count++
		g.pulse = 10
		if g.count >= 30 {
			g.clear = true
		}
	}
	if g.pulse > 0 {
		g.pulse--
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{34, 23, 45, 255})
	ebitenutil.DebugPrintAt(s, "EBI SWEET COUNTER", 175, 55)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SWEETS %03d / 030", g.count), 165, 130)
	scale := 1 + float64(g.pulse)*.012
	r := float32(115 * scale)
	vector.DrawFilledCircle(s, 240, 380, r, color.RGBA{221, 142, 76, 255}, false)
	vector.StrokeCircle(s, 240, 380, r, 8, color.RGBA{255, 202, 113, 255}, false)
	for i := 0; i < 9; i++ {
		a := float64(i) * math.Pi * 2 / 9
		vector.DrawFilledCircle(s, 240+float32(math.Cos(a)*70), 380+float32(math.Sin(a)*70), 10, color.RGBA{91, 50, 51, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, "TAP THE BIG SWEET / PRESS SPACE", 125, 570)
	if g.clear {
		vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "30 SWEETS BAKED!\n\nTAP / SPACE TO COUNT AGAIN", 135, 330)
	}
}
func pressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Sweet Counter — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
