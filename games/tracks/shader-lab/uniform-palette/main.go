// Shader Lab 02 — uniforms make a time palette and a separate hit flash.
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"image/color"
	"math"
)

type game struct {
	t, flash float64
	src      *ebiten.Image
	fx       *shaderlab.Palette
}

func newGame() *game {
	i := ebiten.NewImage(280, 280)
	i.Fill(color.RGBA{20, 62, 110, 255})
	vector.DrawFilledCircle(i, 140, 140, 90, color.RGBA{255, 200, 70, 255}, true)
	return &game{src: i, fx: shaderlab.NewPalette()}
}
func (g *game) Update() error {
	g.t += .04
	g.flash = math.Max(0, g.flash-.04)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.flash = 1
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "SHADER LAB 02 · UNIFORMS", 92, 28)
	ebitenutil.DebugPrintAt(s, "TIME changes palette · tap / SPACE sends FLASH = 1", 55, 55)
	d := ebiten.NewImage(280, 280)
	if !g.fx.Draw(d, g.src, float32(g.t), float32(g.flash)) {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.Scale(float32(.55+.45*math.Sin(g.t)), .7, 1, 1)
		d.DrawImage(g.src, op)
		if g.flash > 0 {
			vector.DrawFilledRect(d, 0, 0, 280, 280, color.RGBA{255, 255, 255, uint8(g.flash * 110)}, true)
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(100, 180)
	s.DrawImage(d, op)
	ebitenutil.DebugPrintAt(s, "Uniform map: { Time: seconds, Flash: hitStrength }", 58, 570)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
