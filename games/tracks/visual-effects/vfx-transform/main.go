// vfx-transform — STEP 02: Rotate / Scale / pivot via live Go + sliders.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
)

type game struct {
	shell *vfxlive.Shell
}

func newGame() *game {
	return &game{
		shell: vfxlive.New(
			"Rotate + Scale + Pivot",
			[]string{
				"op := &ebiten.DrawImageOptions{}",
				"op.GeoM.Translate({ox}, {oy}) // to pivot",
				"op.GeoM.Rotate({angle})",
				"op.GeoM.Scale({scale}, {scale})",
				"op.GeoM.Translate({x}, {y})  // to screen",
				"screen.DrawImage(tenjiroh, op)",
			},
			&vfxlive.Param{Key: "angle", Label: "angle", Value: 0.4, Min: -math.Pi, Max: math.Pi, Format: "%.2f"},
			&vfxlive.Param{Key: "scale", Label: "scale", Value: 1.2, Min: 0.3, Max: 2.8, Format: "%.2f"},
			&vfxlive.Param{Key: "center", Label: "pivotCtr", Value: 1, Bool: true},
		),
	}
}

func (g *game) Update() error {
	g.shell.Update()
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{14, 18, 36, 255})
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh/2
	vector.DrawFilledCircle(s, float32(cx), float32(cy), 4, color.RGBA{120, 240, 220, 255}, false)

	img := hero.Image()
	b := img.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())
	drawH := 140.0
	base := drawH / h
	sc := g.shell.Get("scale") * base

	ox, oy := 0.0, 0.0
	if g.shell.Bool("center") {
		ox, oy = -w/2, -h/2
		g.shell.SetToken("ox", "-w/2")
		g.shell.SetToken("oy", "-h/2")
	} else {
		g.shell.SetToken("ox", "0")
		g.shell.SetToken("oy", "0")
	}
	g.shell.SetToken("x", fmt.Sprintf("%.0f", cx))
	g.shell.SetToken("y", fmt.Sprintf("%.0f", cy))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ox, oy)
	op.GeoM.Rotate(g.shell.Get("angle"))
	op.GeoM.Scale(sc, sc)
	op.GeoM.Translate(cx, cy)
	s.DrawImage(img, op)

	g.shell.Hint = "pivotCtr ON = rotate around center  ·  OFF = around corner"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Transform — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
