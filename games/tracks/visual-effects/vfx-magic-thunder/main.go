// vfx-magic-thunder — STEP 11: branching lightning with live Go knobs.
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

type seg struct {
	x0, y0, x1, y1 float64
	life, max      float64
	thick          float32
}

type game struct {
	shell *vfxlive.Shell
	segs  []seg
	parts []vfxmagic.Particle
	rng   *rand.Rand
	flash float64
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(11)),
		shell: vfxlive.New(
			"Branching bolts",
			[]string{
				"nx += (rand()-0.5) * {wobble}",
				"if depth < 2 && rand() < {branch} { fork() }",
				"StrokeLine(glow); StrokeLine(body); StrokeLine(core)",
				"flash = {flash}",
			},
			&vfxlive.Param{Key: "wobble", Label: "wobble", Value: 48, Min: 10, Max: 90, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "branch", Label: "branch%", Value: 0.35, Min: 0, Max: 0.8, Format: "%.2f"},
			&vfxlive.Param{Key: "thick", Label: "thick", Value: 4.5, Min: 1.5, Max: 8, Format: "%.1f"},
			&vfxlive.Param{Key: "flash", Label: "flash", Value: 1, Min: 0.2, Max: 1, Format: "%.2f"},
			&vfxlive.Param{Key: "cast", Label: "CAST", Value: 0, Bool: true},
		),
	}
}

func (g *game) addBolt(x0, y0, x1, y1 float64, depth int, life float64) {
	segs := 6 + g.rng.Intn(4)
	wobble := g.shell.Get("wobble")
	branch := g.shell.Get("branch")
	baseThick := g.shell.Get("thick")
	px, py := x0, y0
	for i := 1; i <= segs; i++ {
		tt := float64(i) / float64(segs)
		nx := x0 + (x1-x0)*tt + (g.rng.Float64()-0.5)*wobble
		ny := y0 + (y1-y0)*tt + (g.rng.Float64()-0.5)*wobble*0.35
		if i == segs {
			nx, ny = x1, y1
		}
		th := float32(baseThick - float64(depth)*1.2)
		if th < 1.2 {
			th = 1.2
		}
		g.segs = append(g.segs, seg{x0: px, y0: py, x1: nx, y1: ny, life: life, max: life, thick: th})
		if depth < 2 && i > 1 && i < segs && g.rng.Float64() < branch {
			g.addBolt(nx, ny, nx+(g.rng.Float64()-0.5)*120, ny+40+g.rng.Float64()*70, depth+1, life*0.7)
		}
		px, py = nx, ny
	}
}

func (g *game) cast() {
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.75
	skyX := cx + (g.rng.Float64()-0.5)*90
	g.addBolt(skyX, sy+10, cx+(g.rng.Float64()-0.5)*20, cy-20, 0, 18)
	g.flash = g.shell.Get("flash")
	for i := 0; i < 40; i++ {
		life := 20 + g.rng.Float64()*30
		a := g.rng.Float64() * 2 * math.Pi
		sp := 2 + g.rng.Float64()*6
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy - 30, VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: 0.3, Add: true,
			Tint: color.RGBA{220, 235, 255, 255}, Src: vfxsprites.Spark,
		})
	}
}

func (g *game) Update() error {
	g.shell.Update()
	if g.shell.Bool("cast") {
		g.cast()
		g.shell.Param("cast").Value = 0
	}
	if g.flash > 0 {
		g.flash -= 0.07
	}
	segs := g.segs[:0]
	for _, s := range g.segs {
		s.life--
		if s.life > 0 {
			segs = append(segs, s)
		}
	}
	g.segs = segs
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
	s.Fill(color.RGBA{8, 8, 18, 255})
	g.shell.FillStage(s, color.RGBA{10, 10, 22, 255})
	vfxmagic.SoftFlash(s, g.flash, 200, 210, 255)
	_, sy, _, sh := g.shell.Stage()
	hero.DrawBottomCentered(s, 240, sy+sh*0.75+40, 120)
	for _, seg := range g.segs {
		f := seg.life / seg.max
		vector.StrokeLine(s, float32(seg.x0), float32(seg.y0), float32(seg.x1), float32(seg.y1), seg.thick*2.2, color.RGBA{80, 120, 255, uint8(70 * f)}, false)
		vector.StrokeLine(s, float32(seg.x0), float32(seg.y0), float32(seg.x1), float32(seg.y1), seg.thick, color.RGBA{160, 190, 255, uint8(200 * f)}, false)
		vector.StrokeLine(s, float32(seg.x0), float32(seg.y0), float32(seg.x1), float32(seg.y1), seg.thick*0.35, color.RGBA{255, 255, 255, uint8(230 * f)}, false)
	}
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	g.shell.Hint = "wobble/branch%/thick/flash  ·  tap CAST"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Thunder magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
