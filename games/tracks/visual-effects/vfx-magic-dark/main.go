// vfx-magic-dark — Visual Effects Lab STEP 13.
// Showcase: void vortex — spiral wisps, purple fringe, absorb, expanding ring.
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
	castY = 450.0
)

type game struct {
	rng     *rand.Rand
	parts   []vfxmagic.Particle
	flash   float64
	vortex  float64
	casting bool
	casts   int
	t       float64
	btn     vfxui.Button
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(13)),
		btn: vfxui.Button{X: 140, Y: 640, W: 200, H: 54, Label: "CAST DARK", Fill: color.RGBA{36, 12, 48, 235}},
	}
}

func (g *game) cast() {
	g.casts++
	g.flash = 0.55
	g.vortex = 1
	g.casting = true
	// Outer ring shockwave (dark-tinted).
	g.parts = append(g.parts, vfxmagic.Particle{
		X: castX, Y: castY - 30, Life: 50, Max: 50, Scale: 1,
		Add: true, Tint: color.RGBA{120, 40, 180, 255}, Src: vfxsprites.Ring,
		FadeFrom: 0.05, FadeTo: 0.8, ScaleMulFrom: 3.5, ScaleMulTo: 0.5,
	})
	// Spiral wisps forming a vortex.
	for i := 0; i < 80; i++ {
		life := 70 + g.rng.Float64()*60
		t0 := g.rng.Float64()
		arms := float64(g.rng.Intn(3) + 2)
		dx, dy := vfxmagic.SpiralOffset(t0, arms, 160, 4.5)
		// Velocity pulls inward + tangential swirl.
		ang := math.Atan2(dy, dx)
		spIn := 1.2 + g.rng.Float64()*2.5
		spTan := 1.5 + g.rng.Float64()*2
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX + dx, Y: castY - 30 + dy,
			VX:  -math.Cos(ang)*spIn - math.Sin(ang)*spTan,
			VY:  -math.Sin(ang)*spIn + math.Cos(ang)*spTan,
			Rot: ang, Spin: 0.08 + g.rng.Float64()*0.1,
			Life: life, Max: life, Scale: 0.35 + g.rng.Float64()*0.7,
			Add: i%3 != 0, Grav: 0,
			Tint: color.RGBA{uint8(40 + g.rng.Intn(80)), uint8(g.rng.Intn(40)), uint8(80 + g.rng.Intn(120)), 255},
			Src:  vfxsprites.Dark, FadeFrom: 0.15, FadeTo: 1, ScaleMulFrom: 1.2, ScaleMulTo: 0.3,
		})
	}
	// Violet fringe sparks on the rim.
	for i := 0; i < 40; i++ {
		life := 35 + g.rng.Float64()*40
		a := g.rng.Float64() * 2 * math.Pi
		r := 90 + g.rng.Float64()*40
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX + math.Cos(a)*r, Y: castY - 30 + math.Sin(a)*r*0.7,
			VX: -math.Cos(a) * (1 + g.rng.Float64()*2), VY: -math.Sin(a) * (1 + g.rng.Float64()*2),
			Life: life, Max: life, Scale: 0.2 + g.rng.Float64()*0.4,
			Add: true, Tint: color.RGBA{180, 80, 255, 255}, Src: vfxsprites.Spark,
		})
	}
}

func (g *game) Update() error {
	g.t += 0.07
	if g.flash > 0 {
		g.flash -= 0.035
	}
	if g.vortex > 0 {
		g.vortex -= 0.01
	} else {
		g.casting = false
	}

	if g.btn.Tapped() || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.cast()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame()
		return nil
	}

	// Continuous swirl toward center while vortex is active.
	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		if g.vortex > 0.15 && p.Src == vfxsprites.Dark {
			dx := castX - p.X
			dy := castY - 30 - p.Y
			p.VX += dx * 0.004
			p.VY += dy * 0.004
			// Extra tangential force.
			p.VX += -dy * 0.0025
			p.VY += dx * 0.0025
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
	vector.DrawFilledRect(s, 0, 540, width, 180, color.RGBA{12, 6, 20, 255}, false)
	vfxmagic.SoftFlash(s, g.flash, 60, 20, 90)

	hero.DrawBottomCentered(s, castX, castY+70, 150)

	cx, cy := castX, castY-30
	// Dark core that grows then collapses.
	if g.vortex > 0 {
		sc := 0.6 + g.vortex*1.8
		vfxmagic.DrawSprite(s, vfxsprites.Dark, cx, cy, g.t, sc, sc, color.RGBA{20, 0, 40, 255}, float32(0.7*g.vortex), false)
		vfxmagic.DrawSprite(s, vfxsprites.Dark, cx, cy, -g.t*1.3, sc*1.3, sc*1.3, color.RGBA{100, 30, 160, 255}, float32(0.45*g.vortex), true)
		vfxmagic.DrawSprite(s, vfxsprites.Ring, cx, cy, 0, 0.8+g.vortex*2.2, 0.8+g.vortex*2.2,
			color.RGBA{140, 60, 220, 255}, float32(0.35*g.vortex), true)
	}

	for i := range g.parts {
		g.parts[i].Draw(s)
	}

	ebitenutil.DebugPrintAt(s, "DARK MAGIC — VOID VORTEX + PURPLE FRINGE", 42, 18)
	ebitenutil.DebugPrintAt(s, "layers: spiral pull + dark PNG + additive rim + ring", 16, 44)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("vortex %.0f%%   particles %d   casts %d", g.vortex*100, len(g.parts), g.casts), 60, 610)
	g.btn.Draw(s, g.casting || g.vortex > 0.1)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE   R = reset", 150, 700)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cool Dark Magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
