// vfx-stamp — STEP 01: GeoM.Translate via live Go + mouse sliders.
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type game struct {
	shell  *vfxlive.Shell
	stamps [][2]float64
}

func newGame() *game {
	return &game{
		shell: vfxlive.New(
			"Translate",
			[]string{
				"op := &ebiten.DrawImageOptions{}",
				"op.GeoM.Translate({x}, {y})",
				"screen.DrawImage(tenjiroh, op)",
			},
			&vfxlive.Param{Key: "x", Label: "x", Value: 240, Min: 40, Max: 440, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "y", Label: "y", Value: 360, Min: 220, Max: 520, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "size", Label: "size", Value: 120, Min: 48, Max: 200, Step: 1, Format: "%.0f"},
		),
	}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	if ate {
		return nil
	}
	// Tap the stage to jump Translate to that point and leave a stamp ghost.
	if x, y, ok := vfxui.JustPressed(); ok {
		_, sy, _, sh := g.shell.Stage()
		if y >= sy && y <= sy+sh {
			g.shell.Param("x").Set(x)
			g.shell.Param("y").Set(y)
			g.stamps = append(g.stamps, [2]float64{x, y})
			if len(g.stamps) > 12 {
				g.stamps = g.stamps[len(g.stamps)-12:]
			}
		}
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{14, 18, 36, 255})
	_, sy, _, sh := g.shell.Stage()
	// Guide crosshair.
	x, y := g.shell.Get("x"), g.shell.Get("y")
	vector.StrokeLine(s, float32(x), float32(sy), float32(x), float32(sy+sh), 1, color.RGBA{40, 60, 100, 120}, false)
	vector.StrokeLine(s, 0, float32(y), 480, float32(y), 1, color.RGBA{40, 60, 100, 120}, false)
	for _, p := range g.stamps {
		hero.DrawCentered(s, p[0], p[1], 56)
	}
	hero.DrawCentered(s, x, y, g.shell.Get("size"))
	g.shell.Hint = "tap stage to stamp  ·  drag x/y to move Translate"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Translate — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
