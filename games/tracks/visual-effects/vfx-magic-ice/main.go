// vfx-magic-ice — STEP 10: form→shatter with live Go knobs.
package main

import (
	"fmt"
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
	shell   *vfxlive.Shell
	parts   []vfxmagic.Particle
	rng     *rand.Rand
	phase   int
	timer   float64
	crystal float64
	flash   float64
	bloom   float64
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(10)),
		shell: vfxlive.New(
			"Ice form→shatter",
			[]string{
				"phase := {phase} // 0 idle 1 form 2 shatter",
				"f.x += (cx-f.x) * {pull}  // gather",
				"spawnShards(n={count}, grav={grav})",
				"bloomRing(scale={bloom})",
			},
			&vfxlive.Param{Key: "pull", Label: "pull", Value: 0.06, Min: 0.02, Max: 0.15, Format: "%.2f"},
			&vfxlive.Param{Key: "count", Label: "shards", Value: 48, Min: 12, Max: 90, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "grav", Label: "grav", Value: 0.12, Min: 0.02, Max: 0.3, Format: "%.2f"},
			&vfxlive.Param{Key: "bloom", Label: "bloom", Value: 1.2, Min: 0.4, Max: 2.5, Format: "%.1f"},
			&vfxlive.Param{Key: "cast", Label: "CAST", Value: 0, Bool: true},
		),
	}
}

func (g *game) start() {
	if g.phase != 0 {
		return
	}
	g.phase = 1
	g.timer = 0
	g.crystal = 0
	g.flash = 0.35
	g.parts = nil
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.55
	for i := 0; i < 20; i++ {
		a := g.rng.Float64() * 2 * math.Pi
		r := 40 + g.rng.Float64()*80
		life := 80.0
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx + math.Cos(a)*r, Y: cy + math.Sin(a)*r*0.55,
			Rot: a, Life: life, Max: life, Scale: 0.4, Add: true,
			Tint: color.RGBA{200, 240, 255, 255}, Src: vfxsprites.Ice,
		})
	}
}

func (g *game) shatter() {
	g.phase = 2
	g.timer = 0
	g.flash = 0.85
	g.bloom = g.shell.Get("bloom")
	n := int(g.shell.Get("count") + 0.5)
	grav := g.shell.Get("grav")
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.55
	g.parts = nil
	for i := 0; i < n; i++ {
		life := 50 + g.rng.Float64()*50
		a := g.rng.Float64() * 2 * math.Pi
		sp := 2.5 + g.rng.Float64()*6
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, VX: math.Cos(a) * sp, VY: math.Sin(a)*sp*0.85 - 1,
			Rot: a, Spin: (g.rng.Float64() - 0.5) * 0.15,
			Life: life, Max: life, Scale: 0.4 + g.rng.Float64()*0.7,
			Add: false, Grav: grav, Tint: color.RGBA{200, 230, 255, 255}, Src: vfxsprites.Ice,
			FadeFrom: 0.2, FadeTo: 1, ScaleMulFrom: 1.1, ScaleMulTo: 0.4,
		})
	}
	for i := 0; i < 40; i++ {
		life := 25 + g.rng.Float64()*35
		a := g.rng.Float64() * 2 * math.Pi
		sp := 1 + g.rng.Float64()*4
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, VX: math.Cos(a) * sp, VY: math.Sin(a)*sp - 1,
			Life: life, Max: life, Scale: 0.25, Add: true,
			Tint: color.RGBA{200, 240, 255, 255}, Src: vfxsprites.Spark,
		})
	}
}

func (g *game) Update() error {
	g.shell.Update()
	if g.shell.Bool("cast") {
		g.start()
		g.shell.Param("cast").Value = 0
	}
	g.shell.SetToken("phase", fmt.Sprintf("%d", g.phase))
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.55
	pull := g.shell.Get("pull")

	switch g.phase {
	case 1:
		g.timer++
		g.crystal = math.Min(1, g.timer/36)
		for i := range g.parts {
			g.parts[i].X += (cx - g.parts[i].X) * pull
			g.parts[i].Y += (cy - g.parts[i].Y) * pull
			g.parts[i].Rot += 0.08
		}
		if g.timer >= 40 {
			g.shatter()
		}
	case 2:
		g.timer++
		if g.bloom > 0 {
			g.bloom -= 0.015
		}
		if g.timer > 100 && len(g.parts) < 5 {
			g.phase = 0
		}
	}
	if g.flash > 0 {
		g.flash -= 0.05
	}
	if g.phase == 2 {
		alive := g.parts[:0]
		for i := range g.parts {
			p := g.parts[i]
			if p.Update() {
				alive = append(alive, p)
			}
		}
		g.parts = alive
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 12, 28, 255})
	g.shell.FillStage(s, color.RGBA{8, 16, 36, 255})
	vfxmagic.SoftFlash(s, g.flash, 160, 210, 255)
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.55
	hero.DrawBottomCentered(s, cx, sy+sh*0.75+40, 120)
	if g.bloom > 0 {
		vfxmagic.DrawSprite(s, vfxsprites.Ring, cx, cy, 0, g.bloom*2, g.bloom*0.8,
			color.RGBA{120, 200, 255, 255}, float32(0.5*g.bloom), true)
	}
	if g.phase == 1 || (g.phase == 2 && g.timer < 15) {
		sc := 0.4 + g.crystal*1.3
		a := float32(0.55 + g.crystal*0.4)
		vfxmagic.DrawSprite(s, vfxsprites.Ice, cx, cy, g.timer*0.05, sc, sc*1.1, color.White, a, false)
	}
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	g.shell.Hint = "tweak pull/shards/grav/bloom  ·  tap CAST"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Ice magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
