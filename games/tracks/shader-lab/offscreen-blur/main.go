// Shader Lab 05 — render once offscreen, then choose a cheap presentation pass.
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

type game struct {
	t     float64
	mode  int
	small *ebiten.Image
}

func newGame() *game { return &game{small: ebiten.NewImage(120, 120)} }
func (g *game) Update() error {
	g.t += .04
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.mode = (g.mode + 1) % 3
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	g.small.Fill(color.RGBA{14, 42, 78, 255})
	x, y := float32(60+math.Sin(g.t)*18), float32(60)
	vector.DrawFilledCircle(g.small, x, y, 24, color.RGBA{255, 210, 85, 255}, true)
	names := []string{"CRISP UPSCALE", "DOWNSAMPLE", "FAUX BLUR"}
	ebitenutil.DebugPrintAt(s, "SHADER LAB 05 · OFFSCREEN PASSES", 65, 28)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: "+names[g.mode], 145, 55)
	for i := 0; i < 4; i++ {
		if g.mode != 2 {
			break
		}
		op := &ebiten.DrawImageOptions{}
		scale := 2.0 + float64(i)*.10
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(240-float64(120)*scale/2, 290-float64(120)*scale/2)
		op.ColorScale.Scale(0.15, 0.75, 1, float32(.10))
		op.Blend = ebiten.BlendLighter
		s.DrawImage(g.small, op)
	}
	op := &ebiten.DrawImageOptions{}
	scale := 2.0
	if g.mode == 1 {
		scale = 2.8
	}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(240-float64(120)*scale/2, 290-float64(120)*scale/2)
	s.DrawImage(g.small, op)
	ebitenutil.DebugPrintAt(s, "world → small offscreen image → chosen final pass", 82, 580)
	ebitenutil.DebugPrintAt(s, "Use downsampling or a few copies before reaching for expensive blur.", 33, 610)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
