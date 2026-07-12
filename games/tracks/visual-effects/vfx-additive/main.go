// vfx-additive — STEP 05: BlendLighter via live Go + mouse.
package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type game struct {
	shell  *vfxlive.Shell
	glow   *ebiten.Image
	rx, ry float64
	t      float64
}

func makeGlow() *ebiten.Image {
	const R = 64
	img := ebiten.NewImage(R*2, R*2)
	for y := 0; y < R*2; y++ {
		for x := 0; x < R*2; x++ {
			d := math.Hypot(float64(x-R)+0.5, float64(y-R)+0.5) / R
			if d >= 1 {
				continue
			}
			a := math.Pow(1-d, 2.2)
			img.Set(x, y, color.RGBA{255, 255, 255, uint8(255 * a)})
		}
	}
	return img
}

func newGame() *game {
	g := &game{
		glow: makeGlow(),
		shell: vfxlive.New(
			"BlendLighter",
			[]string{
				"op := &ebiten.DrawImageOptions{}",
				"if {add} { op.Blend = ebiten.BlendLighter }",
				"op.ColorScale.ScaleWithColor(orbColor)",
				"screen.DrawImage(glow, op) // ×3 overlapping",
			},
			&vfxlive.Param{Key: "add", Label: "additive", Value: 1, Bool: true},
			&vfxlive.Param{Key: "radius", Label: "radius", Value: 1.4, Min: 0.6, Max: 2.4, Format: "%.2f"},
			&vfxlive.Param{Key: "overlap", Label: "overlap", Value: 0.55, Min: 0.1, Max: 0.95, Format: "%.2f"},
		),
	}
	_, sy, _, sh := g.shell.Stage()
	g.rx, g.ry = 200, sy+sh/2
	return g
}

func (g *game) Update() error {
	g.t += 0.04
	ate := g.shell.Update()
	_, sy, _, sh := g.shell.Stage()
	if !ate {
		if x, y, ok := vfxui.Held(); ok && y >= sy && y <= sy+sh {
			g.rx, g.ry = x, y
		}
	}
	return nil
}

func (g *game) drawOrb(s *ebiten.Image, x, y float64, col color.Color) {
	sc := g.shell.Get("radius")
	b := g.glow.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Scale(sc, sc)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(col)
	if g.shell.Bool("add") {
		op.Blend = ebiten.BlendLighter
	}
	s.DrawImage(g.glow, op)
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 8, 16, 255})
	g.shell.FillStage(s, color.RGBA{6, 8, 16, 255})
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh/2
	spread := 90 * (1.15 - g.shell.Get("overlap"))
	bx := cx + math.Cos(g.t)*spread*0.15
	by := cy + math.Sin(g.t*0.8)*12
	gx := cx + spread
	gy := cy + 10
	g.drawOrb(s, g.rx, g.ry, color.RGBA{255, 80, 60, 255})
	g.drawOrb(s, bx, by, color.RGBA{60, 180, 255, 255})
	g.drawOrb(s, gx, gy, color.RGBA{80, 255, 120, 255})
	vector.StrokeCircle(s, float32(g.rx), float32(g.ry), 10, 2, color.RGBA{255, 200, 200, 180}, false)

	g.shell.Hint = "drag red orb  ·  toggle additive to see light pile up"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Additive — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
