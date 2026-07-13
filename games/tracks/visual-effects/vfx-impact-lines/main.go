// vfx-impact-lines — expanding rings and converging speed lines.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type game struct {
	shell       *vfxlive.Shell
	cx, cy, age float64
}

func newGame() *game {
	g := &game{shell: vfxlive.New("Shockwave + Focus Lines", []string{
		"radius := age * {speed}",
		"StrokeCircle(screen, x, y, radius, {thick}, white)",
		"for i := 0; i < {rays}; i++ { drawFocusLine(i) }",
		"alpha := 1 - age/{life}",
	},
		&vfxlive.Param{Key: "rays", Label: "focus lines", Value: 18, Min: 6, Max: 36, Step: 2, Format: "%.0f"},
		&vfxlive.Param{Key: "speed", Label: "ring speed", Value: 7, Min: 2, Max: 13, Format: "%.1f"},
		&vfxlive.Param{Key: "thick", Label: "ring width", Value: 5, Min: 1, Max: 12, Format: "%.1f"},
	)}
	_, sy, _, sh := g.shell.Stage()
	g.cx, g.cy, g.age = 240, sy+sh/2, 0
	return g
}

func (g *game) Update() error {
	ate := g.shell.Update()
	_, sy, _, sh := g.shell.Stage()
	if !ate {
		if x, y, ok := vfxui.JustPressed(); ok && y >= sy && y <= sy+sh {
			g.cx, g.cy, g.age = x, y, 0
		}
	}
	g.age++
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 13, 28, 255})
	g.shell.FillStage(screen, color.RGBA{17, 22, 48, 255})
	_, _, _, sh := g.shell.Stage()
	life := 42.0
	t := math.Min(1, g.age/life)
	a := uint8(255 * (1 - t))
	radius := g.age * g.shell.Get("speed")
	if g.age < life {
		vector.StrokeCircle(screen, float32(g.cx), float32(g.cy), float32(radius), float32(g.shell.Get("thick")), color.RGBA{110, 240, 220, a}, true)
		vector.StrokeCircle(screen, float32(g.cx), float32(g.cy), float32(radius*.62), 2, color.RGBA{255, 220, 105, a}, true)
		n := int(g.shell.Get("rays"))
		maxR := math.Hypot(480, sh)
		for i := 0; i < n; i++ {
			ang := float64(i)*math.Pi*2/float64(n) + .13*math.Sin(float64(i)*2.7)
			inner := radius + 25 + float64((i*17)%35)
			outer := maxR
			vector.StrokeLine(screen, float32(g.cx+math.Cos(ang)*inner), float32(g.cy+math.Sin(ang)*inner), float32(g.cx+math.Cos(ang)*outer), float32(g.cy+math.Sin(ang)*outer), 2, color.RGBA{255, 245, 220, a / 2}, true)
		}
	}
	vector.DrawFilledCircle(screen, float32(g.cx), float32(g.cy), float32(10+8*(1-t)), color.RGBA{255, 245, 190, max(a, 40)}, true)
	g.shell.SetToken("life", fmt.Sprintf("%.0f", life))
	g.shell.Hint = "tap anywhere on stage  ·  circles + lines make the impact"
	g.shell.Draw(screen)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }
func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Impact Lines — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
