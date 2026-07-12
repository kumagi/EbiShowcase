// vfx-magic-thunder — Visual Effects Lab STEP 11.
// Showcase: branching lightning — jagged bolts, afterimage, sparks, white flash.
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxmagic"
	"github.com/kumagi/EbiShowcase/internal/vfxsprites"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const width, height = 480, 720

const (
	castX = 240.0
	castY = 480.0
)

type segment struct {
	x0, y0, x1, y1 float64
	life, max      float64
	thick          float32
	branch         bool
}

type game struct {
	rng   *rand.Rand
	parts []vfxmagic.Particle
	segs  []segment
	flash float64
	casts int
	t     float64
	btn   vfxui.Button
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(11)),
		btn: vfxui.Button{X: 120, Y: 640, W: 240, H: 54, Label: "CAST THUNDER", Fill: color.RGBA{40, 40, 90, 235}},
	}
}

func (g *game) addBolt(x0, y0, x1, y1 float64, depth int, life float64) {
	segs := 6 + g.rng.Intn(5)
	px, py := x0, y0
	for i := 1; i <= segs; i++ {
		tt := float64(i) / float64(segs)
		nx := x0 + (x1-x0)*tt + (g.rng.Float64()-0.5)*48
		ny := y0 + (y1-y0)*tt + (g.rng.Float64()-0.5)*18
		if i == segs {
			nx, ny = x1, y1
		}
		thick := float32(4.5 - float64(depth)*1.2)
		if thick < 1.5 {
			thick = 1.5
		}
		g.segs = append(g.segs, segment{x0: px, y0: py, x1: nx, y1: ny, life: life, max: life, thick: thick, branch: depth > 0})
		// Occasional branch.
		if depth < 2 && i > 1 && i < segs && g.rng.Float64() < 0.35 {
			bx := nx + (g.rng.Float64()-0.5)*120
			by := ny + 40 + g.rng.Float64()*80
			g.addBolt(nx, ny, bx, by, depth+1, life*0.7)
		}
		px, py = nx, ny
	}
}

func (g *game) cast() {
	g.casts++
	g.flash = 1
	// Main sky-to-hand bolt + side strikes.
	skyX := castX + (g.rng.Float64()-0.5)*100
	g.addBolt(skyX, 40, castX+(g.rng.Float64()-0.5)*30, castY-30, 0, 18+g.rng.Float64()*10)
	if g.rng.Float64() < 0.7 {
		g.addBolt(skyX+(g.rng.Float64()-0.5)*80, 60, castX+(g.rng.Float64()-0.5)*90, castY-80, 0, 14)
	}
	// Impact sparks.
	for i := 0; i < 55; i++ {
		life := 20 + g.rng.Float64()*35
		a := g.rng.Float64() * 2 * math.Pi
		sp := 2 + g.rng.Float64()*7
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX + (g.rng.Float64()-0.5)*40, Y: castY - 40 + (g.rng.Float64()-0.5)*50,
			VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: 0.2 + g.rng.Float64()*0.45,
			Add: true, Grav: 0.06,
			Tint: color.RGBA{220, 235, 255, 255}, Src: vfxsprites.Spark,
		})
	}
	// Bolt sprites along the main path for texture.
	y := 70.0
	x := skyX
	for y < castY-40 {
		life := 14 + g.rng.Float64()*10
		g.parts = append(g.parts, vfxmagic.Particle{
			X: x, Y: y, Rot: (g.rng.Float64() - 0.5) * 0.6,
			Life: life, Max: life, Scale: 0.5 + g.rng.Float64()*0.6,
			Add: true, Tint: color.White, Src: vfxsprites.Bolt, FadeFrom: 0.3, FadeTo: 1,
		})
		y += 32 + g.rng.Float64()*20
		x += (g.rng.Float64() - 0.5) * 50
	}
	g.ringsImpact()
}

func (g *game) ringsImpact() {
	g.parts = append(g.parts, vfxmagic.Particle{
		X: castX, Y: castY - 20, Life: 22, Max: 22, Scale: 1.2,
		Add: true, Tint: color.RGBA{180, 200, 255, 255}, Src: vfxsprites.Ring,
		FadeFrom: 0.05, FadeTo: 0.9, ScaleMulFrom: 2.8, ScaleMulTo: 0.6,
	})
}

func (g *game) Update() error {
	g.t += 0.1
	if g.flash > 0 {
		g.flash -= 0.07
	}
	if g.btn.Tapped() || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.cast()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame()
		return nil
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
	vector.DrawFilledRect(s, 0, 540, width, 180, color.RGBA{14, 14, 28, 255}, false)
	vfxmagic.SoftFlash(s, g.flash, 200, 210, 255)

	hero.DrawBottomCentered(s, castX, castY+55, 150)

	// Soft glow under bolts.
	for _, seg := range g.segs {
		f := seg.life / seg.max
		col := color.RGBA{160, 190, 255, uint8(180 * f)}
		if seg.branch {
			col = color.RGBA{120, 160, 255, uint8(120 * f)}
		}
		// Outer glow
		vector.StrokeLine(s, float32(seg.x0), float32(seg.y0), float32(seg.x1), float32(seg.y1), seg.thick*2.2, color.RGBA{80, 120, 255, uint8(80 * f)}, false)
		vector.StrokeLine(s, float32(seg.x0), float32(seg.y0), float32(seg.x1), float32(seg.y1), seg.thick, col, false)
		// Hot core
		vector.StrokeLine(s, float32(seg.x0), float32(seg.y0), float32(seg.x1), float32(seg.y1), seg.thick*0.35, color.RGBA{255, 255, 255, uint8(230 * f)}, false)
	}

	for i := range g.parts {
		g.parts[i].Draw(s)
	}

	ebitenutil.DebugPrintAt(s, "THUNDER MAGIC — BRANCHING BOLTS + AFTERIMAGE", 36, 18)
	ebitenutil.DebugPrintAt(s, "layers: jagged lines + bolt PNG + sparks + white flash", 16, 44)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("segments %d   particles %d   casts %d", len(g.segs), len(g.parts), g.casts), 60, 610)
	g.btn.Draw(s, len(g.segs) > 0)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE   R = reset", 150, 700)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cool Thunder Magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
