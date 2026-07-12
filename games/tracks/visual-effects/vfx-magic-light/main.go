// vfx-magic-light — STEP 12: radial flare with live Go knobs.
package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxmagic"
	"github.com/kumagi/EbiShowcase/internal/vfxsprites"
)

type ray struct {
	ang, len, life, max float64
	width               float32
}

type game struct {
	shell *vfxlive.Shell
	rays  []ray
	parts []vfxmagic.Particle
	rng   *rand.Rand
	flash float64
	t     float64
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(12)),
		shell: vfxlive.New(
			"Radial flare",
			[]string{
				"for i := 0; i < {rays}; i++ {",
				"  ang := i * 2π / rays",
				"  StrokeLine(cx, cy, cos*len, sin*len)",
				"}",
				"DrawImage(LightPNG) // bloom ×{bloom}",
			},
			&vfxlive.Param{Key: "rays", Label: "rays", Value: 12, Min: 4, Max: 24, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "len", Label: "length", Value: 160, Min: 60, Max: 280, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "bloom", Label: "bloom", Value: 1.4, Min: 0.5, Max: 2.8, Format: "%.1f"},
			&vfxlive.Param{Key: "sparks", Label: "sparks", Value: 50, Min: 10, Max: 100, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "cast", Label: "CAST", Value: 0, Bool: true},
		),
	}
}

func (g *game) cast() {
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.45
	n := int(g.shell.Get("rays") + 0.5)
	length := g.shell.Get("len")
	base := g.rng.Float64() * math.Pi
	g.rays = nil
	for i := 0; i < n; i++ {
		life := 40.0
		g.rays = append(g.rays, ray{
			ang: base + float64(i)*math.Pi*2/float64(n),
			len: length, life: life, max: life, width: 2.5,
		})
	}
	bloom := g.shell.Get("bloom")
	for i := 0; i < 4; i++ {
		life := 30 + float64(i)*6
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, Life: life, Max: life, Scale: bloom * (0.6 + float64(i)*0.25),
			Add: true, Tint: color.RGBA{255, 236, 180, 255}, Src: vfxsprites.Light,
			FadeFrom: 0.05, FadeTo: 0.9, ScaleMulFrom: 1.6, ScaleMulTo: 0.5,
		})
	}
	ns := int(g.shell.Get("sparks") + 0.5)
	for i := 0; i < ns; i++ {
		life := 35 + g.rng.Float64()*40
		a := g.rng.Float64() * 2 * math.Pi
		sp := 0.8 + g.rng.Float64()*4
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: 0.25, Add: true,
			Tint: color.RGBA{255, 230, 160, 255}, Src: vfxsprites.Spark,
		})
	}
	g.flash = 0.85
}

func (g *game) Update() error {
	g.t += 0.05
	g.shell.Update()
	if g.shell.Bool("cast") {
		g.cast()
		g.shell.Param("cast").Value = 0
	}
	if g.flash > 0 {
		g.flash -= 0.04
	}
	rays := g.rays[:0]
	for _, r := range g.rays {
		r.life--
		r.len *= 1.01
		r.ang += 0.006
		if r.life > 0 {
			rays = append(rays, r)
		}
	}
	g.rays = rays
	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		if p.Update() {
			alive = append(alive, p)
		}
	}
	g.parts = alive
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{14, 16, 32, 255})
	vfxmagic.SoftFlash(s, g.flash, 255, 230, 160)
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.45
	hero.DrawBottomCentered(s, cx, sy+sh*0.78+30, 120)
	for _, r := range g.rays {
		f := r.life / r.max
		x1 := cx + math.Cos(r.ang)*r.len
		y1 := cy + math.Sin(r.ang)*r.len
		vector.StrokeLine(s, float32(cx), float32(cy), float32(x1), float32(y1), r.width*1.6, color.RGBA{255, 200, 80, uint8(60 * f)}, false)
		vector.StrokeLine(s, float32(cx), float32(cy), float32(x1), float32(y1), r.width, color.RGBA{255, 230, 160, uint8(200 * f)}, false)
	}
	vfxmagic.DrawSprite(s, vfxsprites.Light, cx, cy, g.t, 1.0, 1.0, color.RGBA{255, 236, 180, 255}, 0.3, true)
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	g.shell.Hint = "rays/length/bloom/sparks  ·  tap CAST"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Light magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
