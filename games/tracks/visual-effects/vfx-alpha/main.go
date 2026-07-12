// vfx-alpha — STEP 04: ScaleAlpha + afterimage trail via live Go + mouse.
package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type game struct {
	shell *vfxlive.Shell
	trail [][2]float64
	x, y  float64
}

func newGame() *game {
	g := &game{
		shell: vfxlive.New(
			"ScaleAlpha + trail",
			[]string{
				"for i, p := range trail {",
				"  a := float32(i) / float32({trail}) * {alpha}",
				"  op.ColorScale.ScaleAlpha(a)",
				"  op.GeoM.Translate(p.x, p.y)",
				"  screen.DrawImage(tenjiroh, op)",
				"}",
			},
			&vfxlive.Param{Key: "trail", Label: "trail", Value: 10, Min: 1, Max: 24, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "alpha", Label: "alpha", Value: 0.85, Min: 0.15, Max: 1, Format: "%.2f"},
			&vfxlive.Param{Key: "ease", Label: "ease", Value: 0.2, Min: 0.05, Max: 0.5, Format: "%.2f"},
		),
	}
	_, sy, _, sh := g.shell.Stage()
	g.x, g.y = 240, sy+sh/2
	return g
}

func (g *game) Update() error {
	ate := g.shell.Update()
	_, sy, _, sh := g.shell.Stage()
	if !ate {
		if x, y, ok := vfxui.Held(); ok && y >= sy && y <= sy+sh {
			ease := g.shell.Get("ease")
			g.x += (x - g.x) * ease
			g.y += (y - g.y) * ease
		}
	}
	g.trail = append(g.trail, [2]float64{g.x, g.y})
	n := int(g.shell.Get("trail") + 0.5)
	if len(g.trail) > n {
		g.trail = g.trail[len(g.trail)-n:]
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{14, 18, 36, 255})
	amax := g.shell.Get("alpha")
	n := len(g.trail)
	for i, p := range g.trail {
		a := float32(1)
		if n > 1 {
			a = float32(i) / float32(n-1) * float32(amax)
		}
		img := hero.Image()
		bb := img.Bounds()
		h := float64(bb.Dy())
		sc := 110 / h
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(bb.Dx())/2, -h/2)
		op.GeoM.Scale(sc, sc)
		op.GeoM.Translate(p[0], p[1])
		op.ColorScale.ScaleAlpha(a)
		s.DrawImage(img, op)
	}
	g.shell.SetToken("trail", fmt.Sprintf("%d", int(g.shell.Get("trail")+0.5)))
	g.shell.Hint = "drag on stage to move  ·  trail/alpha change the afterimage"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Alpha — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
