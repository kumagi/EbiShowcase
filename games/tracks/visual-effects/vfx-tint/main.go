// vfx-tint — STEP 03: ColorScale modes via live Go + mouse.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
)

type game struct {
	shell *vfxlive.Shell
	t     float64
}

func newGame() *game {
	return &game{
		shell: vfxlive.New(
			"ColorScale",
			[]string{
				"op := &ebiten.DrawImageOptions{}",
				"// mode {mode}: 0 normal · 1 tint · 2 flash · 3 shadow",
				"op.ColorScale.Scale({r}, {g}, {b}, {a})",
				"screen.DrawImage(tenjiroh, op)",
			},
			&vfxlive.Param{Key: "mode", Label: "mode", Value: 1, Min: 0, Max: 3, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "hue", Label: "hue", Value: 0.2, Min: 0, Max: 1, Format: "%.2f"},
			&vfxlive.Param{Key: "amount", Label: "amount", Value: 0.85, Min: 0.1, Max: 1, Format: "%.2f"},
		),
	}
}

func (g *game) Update() error {
	g.t += 0.03
	g.shell.Update()
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{14, 18, 36, 255})
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh/2

	mode := int(g.shell.Get("mode") + 0.5)
	amt := g.shell.Get("amount")
	hue := g.shell.Get("hue")
	var r, gg, b, a float64 = 1, 1, 1, 1
	switch mode {
	case 1:
		r = 0.4 + 0.6*math.Sin(hue*math.Pi*2)
		gg = 0.4 + 0.6*math.Sin(hue*math.Pi*2+2.1)
		b = 0.4 + 0.6*math.Sin(hue*math.Pi*2+4.2)
		a = amt
	case 2:
		f := 1 + 2*amt*(0.5+0.5*math.Sin(g.t*4))
		r, gg, b, a = f, f, f, 1
	case 3:
		r, gg, b, a = amt*0.25, amt*0.25, amt*0.35, amt
	}
	g.shell.SetToken("r", fmt.Sprintf("%.2f", r))
	g.shell.SetToken("g", fmt.Sprintf("%.2f", gg))
	g.shell.SetToken("b", fmt.Sprintf("%.2f", b))
	g.shell.SetToken("a", fmt.Sprintf("%.2f", a))

	img := hero.Image()
	bb := img.Bounds()
	h := float64(bb.Dy())
	sc := 160 / h
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(bb.Dx())/2, -h/2)
	op.GeoM.Scale(sc, sc)
	op.GeoM.Translate(cx, cy)
	op.ColorScale.Scale(float32(r), float32(gg), float32(b), float32(a))
	s.DrawImage(img, op)

	g.shell.Hint = "drag mode / hue / amount — ColorScale numbers update live"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: ColorScale — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
