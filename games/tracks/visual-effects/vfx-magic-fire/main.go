// vfx-magic-fire — STEP 09: charged fire with live Go knobs.
package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxmagic"
	"github.com/kumagi/EbiShowcase/internal/vfxsprites"
)

type game struct {
	shell *vfxlive.Shell
	parts []vfxmagic.Particle
	rng   *rand.Rand
	flash float64
	t     float64
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(9)),
		shell: vfxlive.New(
			"Fire layers",
			[]string{
				"power := {power}",
				"spawnFlames(n={count}, scale={scale})",
				"p.Add = true; p.Grav = {grav}",
				"flash = 0.35 + power*0.55",
				"DrawImage(FirePNG) // BlendLighter",
			},
			&vfxlive.Param{Key: "power", Label: "power", Value: 0.7, Min: 0.15, Max: 1, Format: "%.2f"},
			&vfxlive.Param{Key: "count", Label: "count", Value: 40, Min: 10, Max: 90, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "scale", Label: "scale", Value: 0.8, Min: 0.3, Max: 1.8, Format: "%.2f"},
			&vfxlive.Param{Key: "grav", Label: "rise", Value: -0.02, Min: -0.06, Max: 0, Format: "%.3f"},
			&vfxlive.Param{Key: "burst", Label: "BURST", Value: 0, Bool: true},
		),
	}
}

func (g *game) burst() {
	power := g.shell.Get("power")
	n := int(g.shell.Get("count") + 0.5)
	sc := g.shell.Get("scale")
	grav := g.shell.Get("grav")
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.75
	for i := 0; i < n; i++ {
		life := 45 + g.rng.Float64()*55
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx + (g.rng.Float64()-0.5)*(20+power*40), Y: cy,
			VX: (g.rng.Float64() - 0.5) * (1 + power), VY: -1.2 - g.rng.Float64()*(2+power*3),
			Rot: (g.rng.Float64() - 0.5) * 0.4, Spin: (g.rng.Float64() - 0.5) * 0.04,
			Life: life, Max: life, Scale: sc * (0.5 + g.rng.Float64()*0.8),
			Add: true, Grav: grav, Tint: color.RGBA{255, uint8(80 + g.rng.Intn(100)), 30, 255},
			Src: vfxsprites.Fire, FadeFrom: 0.1, FadeTo: 1, ScaleMulFrom: 1.1, ScaleMulTo: 0.5,
		})
	}
	for i := 0; i < n/2; i++ {
		life := 30 + g.rng.Float64()*40
		a := -math.Pi/2 + (g.rng.Float64()-0.5)*(0.8+power)
		sp := 1.2 + g.rng.Float64()*(3+power*3)
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: sc * 0.35, Add: true, Grav: grav * 1.5,
			Tint: color.RGBA{255, 140, 40, 255}, Src: vfxsprites.Spark,
		})
	}
	g.parts = append(g.parts, vfxmagic.Particle{
		X: cx, Y: cy - 20, Life: 28, Max: 28, Scale: 0.8 + power,
		Add: true, Tint: color.RGBA{255, 100, 30, 255}, Src: vfxsprites.Ring,
		FadeFrom: 0.05, FadeTo: 0.8, ScaleMulFrom: 2.5, ScaleMulTo: 0.5,
	})
	g.flash = 0.35 + power*0.55
}

func (g *game) Update() error {
	g.t += 0.08
	g.shell.Update()
	if g.shell.Bool("burst") {
		g.burst()
		g.shell.Param("burst").Value = 0
	}
	if g.flash > 0 {
		g.flash -= 0.04
	}
	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		p.X += math.Sin(g.t*2+p.Y*0.04) * 0.35
		if p.Update() {
			alive = append(alive, p)
		}
	}
	g.parts = alive
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{14, 8, 10, 255})
	g.shell.FillStage(s, color.RGBA{18, 10, 10, 255})
	vfxmagic.SoftFlash(s, g.flash, 255, 120, 40)
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.75
	hero.DrawBottomCentered(s, cx, cy+45, 130)
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	g.shell.Hint = "set power/count/scale/rise  ·  tap BURST to cast"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Fire magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
