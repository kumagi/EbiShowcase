// vfx-outline — a readable outline made from tinted offset copies.
package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
)

type game struct{ shell *vfxlive.Shell }

func newGame() *game {
	return &game{shell: vfxlive.New("Outline by Copies", []string{
		"for i := 0; i < {copies}; i++ {",
		"  angle := float64(i) * 2*Pi/{copies}",
		"  drawTinted(cos(angle)*{width}, sin(angle)*{width})",
		"}",
		"screen.DrawImage(sprite, mainOp)",
	},
		&vfxlive.Param{Key: "width", Label: "outline px", Value: 4, Min: 1, Max: 10, Step: 1, Format: "%.0f"},
		&vfxlive.Param{Key: "copies", Label: "copies", Value: 8, Min: 4, Max: 16, Step: 4, Format: "%.0f"},
		&vfxlive.Param{Key: "light", Label: "light edge", Value: 0, Bool: true},
	)}
}

func (g *game) Update() error { g.shell.Update(); return nil }

func (g *game) draw(screen *ebiten.Image, dx, dy float64, tint color.Color) {
	img := hero.Image()
	b := img.Bounds()
	scale := 145 / float64(b.Dy())
	_, sy, _, sh := g.shell.Stage()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(240+dx, sy+sh/2+dy)
	if tint != nil {
		op.ColorScale.ScaleWithColor(tint)
	}
	screen.DrawImage(img, op)
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{235, 210, 130, 255})
	g.shell.FillStage(screen, color.RGBA{235, 210, 130, 255})
	n := int(g.shell.Get("copies"))
	r := g.shell.Get("width")
	outline := color.RGBA{12, 20, 42, 255}
	if g.shell.Bool("light") {
		outline = color.RGBA{255, 255, 245, 255}
	}
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.draw(screen, math.Cos(a)*r, math.Sin(a)*r, outline)
	}
	g.draw(screen, 0, 0, nil)
	g.shell.Hint = "the hitbox never changes  ·  only extra DrawImage calls"
	g.shell.Draw(screen)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }
func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Outline — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
