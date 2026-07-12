// vfx-particles — STEP 07: particle slice via live Go + mouse.
package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type particle struct {
	x, y, vx, vy, life, max, size float64
}

type game struct {
	shell *vfxlive.Shell
	glow  *ebiten.Image
	parts []particle
	rng   *rand.Rand
}

func makeGlow() *ebiten.Image {
	const R = 16
	img := ebiten.NewImage(R*2, R*2)
	for y := 0; y < R*2; y++ {
		for x := 0; x < R*2; x++ {
			d := math.Hypot(float64(x-R)+0.5, float64(y-R)+0.5) / R
			if d >= 1 {
				continue
			}
			img.Set(x, y, color.RGBA{255, 255, 255, uint8(255 * (1 - d) * (1 - d))})
		}
	}
	return img
}

func newGame() *game {
	return &game{
		glow: makeGlow(),
		rng:  rand.New(rand.NewSource(7)),
		shell: vfxlive.New(
			"Particles",
			[]string{
				"for i := 0; i < {count}; i++ { spawn() }",
				"p.vy += {grav}           // gravity",
				"p.life--; alpha := p.life/p.max",
				"if {add} { op.Blend = BlendLighter }",
			},
			&vfxlive.Param{Key: "count", Label: "count", Value: 28, Min: 4, Max: 80, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "speed", Label: "speed", Value: 3.5, Min: 0.5, Max: 8, Format: "%.1f"},
			&vfxlive.Param{Key: "grav", Label: "gravity", Value: 0.14, Min: 0, Max: 0.4, Format: "%.2f"},
			&vfxlive.Param{Key: "add", Label: "additive", Value: 1, Bool: true},
		),
	}
}

func (g *game) burst(cx, cy float64) {
	n := int(g.shell.Get("count") + 0.5)
	sp := g.shell.Get("speed")
	for i := 0; i < n; i++ {
		life := 35 + g.rng.Float64()*40
		a := g.rng.Float64() * 2 * math.Pi
		s := sp * (0.4 + g.rng.Float64())
		g.parts = append(g.parts, particle{
			x: cx, y: cy, vx: math.Cos(a) * s, vy: math.Sin(a) * s,
			life: life, max: life, size: 0.6 + g.rng.Float64()*1.2,
		})
	}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	_, sy, _, sh := g.shell.Stage()
	if !ate {
		if x, y, ok := vfxui.JustPressed(); ok && y >= sy && y <= sy+sh {
			g.burst(x, y)
		}
	}
	grav := g.shell.Get("grav")
	alive := g.parts[:0]
	for _, p := range g.parts {
		p.x += p.vx
		p.y += p.vy
		p.vy += grav
		p.vx *= 0.99
		p.life--
		if p.life > 0 {
			alive = append(alive, p)
		}
	}
	g.parts = alive
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 10, 20, 255})
	g.shell.FillStage(s, color.RGBA{8, 10, 20, 255})
	for _, p := range g.parts {
		f := p.life / p.max
		b := g.glow.Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
		op.GeoM.Scale(p.size*f, p.size*f)
		op.GeoM.Translate(p.x, p.y)
		op.ColorScale.Scale(1, 0.7+0.3*float32(f), 0.3, float32(0.2+0.8*f))
		if g.shell.Bool("add") {
			op.Blend = ebiten.BlendLighter
		}
		s.DrawImage(g.glow, op)
	}
	_, sy, _, _ := g.shell.Stage()
	vector.StrokeCircle(s, 240, float32(sy+80), 18, 2, color.RGBA{80, 100, 140, 120}, false)
	g.shell.Hint = "tap stage to burst  ·  tweak count/speed/gravity/additive"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Particles — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
