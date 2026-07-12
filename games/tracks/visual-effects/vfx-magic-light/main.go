// vfx-magic-light — Visual Effects Lab STEP 12.
// Showcase: holy flare — radial rays, soft bloom, expanding rings, twinkles.
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
	castY = 460.0
)

type ray struct {
	ang, len, life, max float64
	width               float32
}

type game struct {
	rng   *rand.Rand
	parts []vfxmagic.Particle
	rays  []ray
	flash float64
	pulse float64
	casts int
	t     float64
	btn   vfxui.Button
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(12)),
		btn: vfxui.Button{X: 140, Y: 640, W: 200, H: 54, Label: "CAST LIGHT", Fill: color.RGBA{90, 70, 20, 235}},
	}
}

func (g *game) cast() {
	g.casts++
	g.flash = 0.9
	g.pulse = 1
	// Long radial rays.
	n := 10 + g.rng.Intn(6)
	base := g.rng.Float64() * math.Pi
	for i := 0; i < n; i++ {
		life := 35 + g.rng.Float64()*25
		g.rays = append(g.rays, ray{
			ang: base + float64(i)*math.Pi*2/float64(n) + (g.rng.Float64()-0.5)*0.15,
			len: 120 + g.rng.Float64()*160, life: life, max: life,
			width: float32(2 + g.rng.Float64()*3),
		})
	}
	// Soft flare core + rings.
	for i := 0; i < 5; i++ {
		life := 28 + float64(i)*8
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX, Y: castY - 40, Rot: g.rng.Float64() * math.Pi,
			Life: life, Max: life, Scale: 0.8 + float64(i)*0.35,
			Add: true, Tint: color.RGBA{255, 236, 180, 255}, Src: vfxsprites.Light,
			FadeFrom: 0.05, FadeTo: 0.95, ScaleMulFrom: 1.8 + float64(i)*0.4, ScaleMulTo: 0.5,
		})
	}
	g.parts = append(g.parts, vfxmagic.Particle{
		X: castX, Y: castY - 40, Life: 40, Max: 40, Scale: 1,
		Add: true, Tint: color.RGBA{255, 220, 140, 255}, Src: vfxsprites.Ring,
		FadeFrom: 0.05, FadeTo: 0.85, ScaleMulFrom: 3.2, ScaleMulTo: 0.4,
	})
	// Twinkling sparks outward.
	for i := 0; i < 70; i++ {
		life := 40 + g.rng.Float64()*50
		a := g.rng.Float64() * 2 * math.Pi
		sp := 0.8 + g.rng.Float64()*4.5
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX, Y: castY - 40,
			VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: 0.15 + g.rng.Float64()*0.4,
			Add: true, Grav: -0.01,
			Tint: color.RGBA{255, uint8(220 + g.rng.Intn(35)), uint8(150 + g.rng.Intn(80)), 255},
			Src:  vfxsprites.Spark, FadeFrom: 0.1, FadeTo: 1,
		})
	}
}

func (g *game) Update() error {
	g.t += 0.06
	if g.flash > 0 {
		g.flash -= 0.04
	}
	if g.pulse > 0 {
		g.pulse -= 0.015
	}
	if g.btn.Tapped() || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.cast()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame()
		return nil
	}

	rays := g.rays[:0]
	for _, r := range g.rays {
		r.life--
		r.len *= 1.012
		r.ang += 0.008
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
	vector.DrawFilledRect(s, 0, 520, width, 200, color.RGBA{22, 20, 36, 255}, false)
	vfxmagic.SoftFlash(s, g.flash, 255, 230, 160)
	if g.pulse > 0 {
		vfxmagic.SoftFlash(s, g.pulse*0.2, 255, 240, 200)
	}

	hero.DrawBottomCentered(s, castX, castY+70, 150)

	cx, cy := castX, castY-40
	for _, r := range g.rays {
		f := r.life / r.max
		x1 := cx + math.Cos(r.ang)*r.len
		y1 := cy + math.Sin(r.ang)*r.len
		col := color.RGBA{255, 230, 160, uint8(200 * f)}
		vector.StrokeLine(s, float32(cx), float32(cy), float32(x1), float32(y1), r.width*1.8, color.RGBA{255, 200, 80, uint8(70 * f)}, false)
		vector.StrokeLine(s, float32(cx), float32(cy), float32(x1), float32(y1), r.width, col, false)
	}

	// Persistent soft halo.
	vfxmagic.DrawSprite(s, vfxsprites.Light, cx, cy, g.t*0.2, 1.1+0.15*math.Sin(g.t), 1.1+0.15*math.Sin(g.t),
		color.RGBA{255, 236, 180, 255}, 0.25, true)

	for i := range g.parts {
		g.parts[i].Draw(s)
	}

	ebitenutil.DebugPrintAt(s, "LIGHT MAGIC — FLARE RAYS + SOFT BLOOM", 55, 18)
	ebitenutil.DebugPrintAt(s, "layers: radial lines + light PNG + rings + twinkles", 16, 44)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("rays %d   particles %d   casts %d", len(g.rays), len(g.parts), g.casts), 80, 610)
	g.btn.Draw(s, len(g.rays) > 0 || g.pulse > 0.2)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE   R = reset", 150, 700)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cool Light Magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
