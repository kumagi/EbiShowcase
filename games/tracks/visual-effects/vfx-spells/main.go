// vfx-spells — STEP 08: composed spells via live Go + mouse knobs.
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

type game struct {
	shell  *vfxlive.Shell
	parts  []vfxmagic.Particle
	bolts  []vfxmagic.Particle
	rng    *rand.Rand
	flash  float64
	last   int
	castCD int
}

func newGame() *game {
	g := &game{
		rng:  rand.New(rand.NewSource(8)),
		last: -1,
		shell: vfxlive.New(
			"Compose spells",
			[]string{
				"// spell {spell}: 0 fire · 1 water · 2 thunder",
				"p.Add = {add}  // BlendLighter when true",
				"p.Grav = {grav}",
				"p.Scale = {scale}",
				"spawn({count}) // textured sprites",
			},
			&vfxlive.Param{Key: "spell", Label: "spell", Value: 0, Min: 0, Max: 2, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "count", Label: "count", Value: 30, Min: 8, Max: 70, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "scale", Label: "scale", Value: 0.7, Min: 0.25, Max: 1.6, Format: "%.2f"},
			&vfxlive.Param{Key: "auto", Label: "autoCast", Value: 1, Bool: true},
		),
	}
	g.shell.SetToken("add", "true")
	g.shell.SetToken("grav", "-0.02")
	return g
}

func (g *game) cast() {
	kind := int(g.shell.Get("spell") + 0.5)
	n := int(g.shell.Get("count") + 0.5)
	sc := g.shell.Get("scale")
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.72
	switch kind {
	case 0:
		g.shell.SetToken("add", "true")
		g.shell.SetToken("grav", "-0.02")
		for i := 0; i < n; i++ {
			life := 45 + g.rng.Float64()*50
			g.parts = append(g.parts, vfxmagic.Particle{
				X: cx + (g.rng.Float64()-0.5)*30, Y: cy,
				VX: (g.rng.Float64() - 0.5) * 1.4, VY: -1.5 - g.rng.Float64()*2.5,
				Life: life, Max: life, Scale: sc * (0.5 + g.rng.Float64()*0.7),
				Add: true, Grav: -0.015, Tint: color.RGBA{255, 120, 40, 255}, Src: vfxsprites.Fire,
			})
		}
		g.flash = 0.45
	case 1:
		g.shell.SetToken("add", "false")
		g.shell.SetToken("grav", "0.24")
		for i := 0; i < n; i++ {
			life := 60 + g.rng.Float64()*40
			a := -math.Pi/2 + (g.rng.Float64()-0.5)*1.4
			sp := 3 + g.rng.Float64()*4
			g.parts = append(g.parts, vfxmagic.Particle{
				X: cx, Y: cy - 10, VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
				Rot: g.rng.Float64() * math.Pi, Spin: (g.rng.Float64() - 0.5) * 0.1,
				Life: life, Max: life, Scale: sc * (0.45 + g.rng.Float64()*0.5),
				Add: false, Grav: 0.24, Tint: color.RGBA{180, 220, 255, 255}, Src: vfxsprites.Water,
			})
		}
		g.flash = 0.25
	default:
		g.shell.SetToken("add", "true")
		g.shell.SetToken("grav", "0.04")
		y := sy + 20
		x := cx + (g.rng.Float64()-0.5)*80
		for y < cy-20 {
			life := 16.0
			g.bolts = append(g.bolts, vfxmagic.Particle{
				X: x, Y: y, Rot: (g.rng.Float64() - 0.5) * 0.5,
				Life: life, Max: life, Scale: sc * (0.5 + g.rng.Float64()*0.5),
				Add: true, Src: vfxsprites.Bolt, Tint: color.White,
			})
			y += 30
			x += (g.rng.Float64() - 0.5) * 50
		}
		for i := 0; i < n/2; i++ {
			life := 24 + g.rng.Float64()*20
			a := g.rng.Float64() * 2 * math.Pi
			sp := 2 + g.rng.Float64()*4
			g.parts = append(g.parts, vfxmagic.Particle{
				X: cx, Y: cy - 40, VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
				Life: life, Max: life, Scale: sc * 0.4, Add: true, Grav: 0.04,
				Tint: color.RGBA{210, 230, 255, 255}, Src: vfxsprites.Spark,
			})
		}
		g.flash = 0.9
	}
}

func (g *game) Update() error {
	g.shell.Update()
	kind := int(g.shell.Get("spell") + 0.5)
	if kind != g.last {
		g.last = kind
		g.cast()
		g.castCD = 50
	}
	if g.shell.Bool("auto") {
		g.castCD--
		if g.castCD <= 0 {
			g.cast()
			g.castCD = 55
		}
	}
	if g.flash > 0 {
		g.flash -= 0.05
	}
	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		if p.Update() {
			alive = append(alive, p)
		}
	}
	g.parts = alive
	bolts := g.bolts[:0]
	for i := range g.bolts {
		p := g.bolts[i]
		if p.Update() {
			bolts = append(bolts, p)
		}
	}
	g.bolts = bolts
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 12, 24, 255})
	g.shell.FillStage(s, color.RGBA{10, 12, 24, 255})
	vfxmagic.SoftFlash(s, g.flash, 180, 200, 255)
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.72
	hero.DrawBottomCentered(s, cx, cy+50, 130)
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	for i := range g.bolts {
		g.bolts[i].Draw(s)
	}
	vector.StrokeCircle(s, float32(cx), float32(cy), 8, 2, color.RGBA{120, 240, 220, 100}, false)
	g.shell.Hint = "change spell 0/1/2 — recipe tokens + look update together"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Spells — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
