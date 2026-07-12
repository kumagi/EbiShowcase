// vfx-magic-dark — STEP 13: void vortex with live Go knobs.
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
	shell  *vfxlive.Shell
	parts  []vfxmagic.Particle
	rng    *rand.Rand
	vortex float64
	flash  float64
	t      float64
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(13)),
		shell: vfxlive.New(
			"Void vortex",
			[]string{
				"dx, dy := SpiralOffset(t, arms={arms}, r={radius})",
				"p.VX += (cx-p.X) * {pull}   // absorb",
				"p.VX += -(cy-p.Y) * {spin}  // swirl",
				"DrawImage(DarkPNG) // core + purple fringe",
			},
			&vfxlive.Param{Key: "arms", Label: "arms", Value: 3, Min: 2, Max: 6, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "radius", Label: "radius", Value: 140, Min: 60, Max: 220, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "pull", Label: "pull", Value: 0.004, Min: 0.001, Max: 0.012, Format: "%.3f"},
			&vfxlive.Param{Key: "spin", Label: "spin", Value: 0.0025, Min: 0, Max: 0.01, Format: "%.4f"},
			&vfxlive.Param{Key: "cast", Label: "CAST", Value: 0, Bool: true},
		),
	}
}

func (g *game) cast() {
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.45
	arms := g.shell.Get("arms")
	radius := g.shell.Get("radius")
	g.vortex = 1
	g.flash = 0.5
	g.parts = append(g.parts, vfxmagic.Particle{
		X: cx, Y: cy, Life: 50, Max: 50, Scale: 1,
		Add: true, Tint: color.RGBA{120, 40, 180, 255}, Src: vfxsprites.Ring,
		FadeFrom: 0.05, FadeTo: 0.75, ScaleMulFrom: 3.2, ScaleMulTo: 0.5,
	})
	for i := 0; i < 70; i++ {
		life := 70 + g.rng.Float64()*50
		t0 := g.rng.Float64()
		dx, dy := vfxmagic.SpiralOffset(t0, arms, radius, 4.5)
		ang := math.Atan2(dy, dx)
		spIn := 1 + g.rng.Float64()*2
		spTan := 1.2 + g.rng.Float64()*2
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx + dx, Y: cy + dy,
			VX:  -math.Cos(ang)*spIn - math.Sin(ang)*spTan,
			VY:  -math.Sin(ang)*spIn + math.Cos(ang)*spTan,
			Rot: ang, Spin: 0.08, Life: life, Max: life,
			Scale: 0.4 + g.rng.Float64()*0.6, Add: i%3 != 0,
			Tint: color.RGBA{uint8(40 + g.rng.Intn(70)), 10, uint8(80 + g.rng.Intn(100)), 255},
			Src:  vfxsprites.Dark, FadeFrom: 0.15, FadeTo: 1, ScaleMulFrom: 1.2, ScaleMulTo: 0.3,
		})
	}
}

func (g *game) Update() error {
	g.t += 0.06
	g.shell.Update()
	if g.shell.Bool("cast") {
		g.cast()
		g.shell.Param("cast").Value = 0
	}
	if g.flash > 0 {
		g.flash -= 0.03
	}
	if g.vortex > 0 {
		g.vortex -= 0.008
	}
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.45
	pull := g.shell.Get("pull")
	spin := g.shell.Get("spin")
	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		if g.vortex > 0.1 && p.Src == vfxsprites.Dark {
			dx := cx - p.X
			dy := cy - p.Y
			p.VX += dx * pull
			p.VY += dy * pull
			p.VX += -dy * spin
			p.VY += dx * spin
		}
		if p.Update() {
			alive = append(alive, p)
		}
	}
	g.parts = alive
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 4, 12, 255})
	g.shell.FillStage(s, color.RGBA{8, 6, 16, 255})
	vfxmagic.SoftFlash(s, g.flash, 60, 20, 90)
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.45
	hero.DrawBottomCentered(s, cx, sy+sh*0.78+30, 120)
	if g.vortex > 0 {
		sc := 0.5 + g.vortex*1.6
		vfxmagic.DrawSprite(s, vfxsprites.Dark, cx, cy, g.t, sc, sc, color.RGBA{20, 0, 40, 255}, float32(0.65*g.vortex), false)
		vfxmagic.DrawSprite(s, vfxsprites.Dark, cx, cy, -g.t, sc*1.25, sc*1.25, color.RGBA{100, 30, 160, 255}, float32(0.4*g.vortex), true)
	}
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	g.shell.Hint = "arms/radius/pull/spin  ·  tap CAST"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Dark magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
