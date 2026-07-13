// vfx-faux-bloom — cheap bloom from larger translucent additive copies.
package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
)

type game struct {
	shell *vfxlive.Shell
	glow  *ebiten.Image
	t     float64
}

func makeGlow() *ebiten.Image {
	const size = 128
	img := ebiten.NewImage(size, size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			d := math.Hypot(float64(x-64), float64(y-64)) / 64
			if d < 1 {
				img.Set(x, y, color.RGBA{255, 255, 255, uint8(100 * math.Pow(1-d, 2))})
			}
		}
	}
	return img
}

func newGame() *game {
	return &game{glow: makeGlow(), shell: vfxlive.New("Cheap Bloom", []string{
		"for i := 0; i < {copies}; i++ {",
		"  op.GeoM.Scale(1 + float64(i)*{spread}, ...)",
		"  op.ColorScale.ScaleAlpha({alpha})",
		"  op.Blend = ebiten.BlendLighter",
		"  screen.DrawImage(glow, op)",
		"}",
	},
		&vfxlive.Param{Key: "copies", Label: "glow copies", Value: 5, Min: 1, Max: 9, Step: 1, Format: "%.0f"},
		&vfxlive.Param{Key: "spread", Label: "spread", Value: .22, Min: .05, Max: .5, Format: "%.2f"},
		&vfxlive.Param{Key: "alpha", Label: "alpha", Value: .32, Min: .08, Max: .65, Format: "%.2f"},
	)}
}

func (g *game) Update() error { g.t += .04; g.shell.Update(); return nil }

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{4, 7, 18, 255})
	g.shell.FillStage(screen, color.RGBA{5, 9, 22, 255})
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh/2
	pulse := 1 + .05*math.Sin(g.t)
	n := int(g.shell.Get("copies"))
	for i := n - 1; i >= 0; i-- {
		scale := pulse * (1 + float64(i)*g.shell.Get("spread"))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-64, -64)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(cx, cy)
		op.ColorScale.Scale(0.25, 0.75, 1, float32(g.shell.Get("alpha")))
		op.Blend = ebiten.BlendLighter
		screen.DrawImage(g.glow, op)
	}
	vector.DrawFilledCircle(screen, float32(cx), float32(cy), 30, color.RGBA{220, 250, 255, 255}, true)
	vector.DrawFilledCircle(screen, float32(cx-9), float32(cy-8), 7, color.White, true)
	g.shell.Hint = "larger faint copies imitate blur  ·  no shader needed"
	g.shell.Draw(screen)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }
func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Faux Bloom — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
