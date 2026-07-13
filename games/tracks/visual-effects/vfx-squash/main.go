// vfx-squash — squash, stretch, anticipation, and overshoot with one sprite.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type game struct {
	shell *vfxlive.Shell
	timer float64
}

func newGame() *game {
	return &game{timer: math.Pi / 2, shell: vfxlive.New("Squash & Stretch", []string{
		"s := sin({phase}) * {amount}",
		"scaleX, scaleY := 1+s, 1-s",
		"op.GeoM.Scale(scaleX, scaleY)",
		"// anchor the feet, not the center",
	},
		&vfxlive.Param{Key: "amount", Label: "squash", Value: .24, Min: 0, Max: .42, Format: "%.2f"},
		&vfxlive.Param{Key: "speed", Label: "speed", Value: .13, Min: .04, Max: .3, Format: "%.2f"},
		&vfxlive.Param{Key: "auto", Label: "auto", Value: 1, Bool: true},
	)}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	_, sy, _, sh := g.shell.Stage()
	if !ate {
		if _, y, ok := vfxui.JustPressed(); ok && y >= sy && y <= sy+sh {
			g.timer = 0
		}
	}
	if g.shell.Bool("auto") || g.timer < math.Pi*2 {
		g.timer += g.shell.Get("speed")
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 14, 30, 255})
	g.shell.FillStage(screen, color.RGBA{19, 29, 55, 255})
	_, sy, _, sh := g.shell.Stage()
	phase := math.Mod(g.timer, math.Pi*2)
	// Use two harmonics: a readable anticipation squash and a quick stretch.
	s := math.Sin(phase) * g.shell.Get("amount")
	sx, syScale := 1+s, 1-s

	img := hero.Image()
	b := img.Bounds()
	base := 140 / float64(b.Dy())
	feetX, feetY := 240.0, sy+sh*.72
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy()))
	op.GeoM.Scale(base*sx, base*syScale)
	op.GeoM.Translate(feetX, feetY)
	screen.DrawImage(img, op)

	g.shell.SetToken("phase", format(phase))
	g.shell.Hint = "tap stage to replay  ·  watch the feet stay planted"
	g.shell.Draw(screen)
}

func format(v float64) string { return fmt.Sprintf("%.2f", v) }

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Squash & Stretch — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
