// vfx-magic-fire — Visual Effects Lab STEP 09.
// Showcase: layered fire plume — rising flame tongues, embers, heat ring, flash.
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

type ring struct {
	life, max, scale float64
}

type game struct {
	rng     *rand.Rand
	parts   []vfxmagic.Particle
	rings   []ring
	flash   float64
	heat    float64
	charge  float64
	casting bool
	t       float64
	casts   int
	btn     vfxui.Button
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(9)),
		btn: vfxui.Button{X: 140, Y: 640, W: 200, H: 54, Label: "CAST FIRE", Fill: color.RGBA{90, 28, 12, 235}},
	}
}

func (g *game) burst(power float64) {
	g.casts++
	nFlame := int(22 + power*40)
	nEmber := int(30 + power*55)
	for i := 0; i < nFlame; i++ {
		life := 50 + g.rng.Float64()*70
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX + (g.rng.Float64()-0.5)*(18+power*40), Y: castY + g.rng.Float64()*10,
			VX: (g.rng.Float64() - 0.5) * (1.2 + power), VY: -1.4 - g.rng.Float64()*(2.2+power*2.5),
			Rot: (g.rng.Float64() - 0.5) * 0.5, Spin: (g.rng.Float64() - 0.5) * 0.05,
			Life: life, Max: life, Scale: 0.45 + g.rng.Float64()*0.85 + power*0.25,
			Add: true, Grav: -0.012 - power*0.01,
			Tint: color.RGBA{255, uint8(70 + g.rng.Intn(100)), uint8(10 + g.rng.Intn(40)), 255},
			Src:  vfxsprites.Fire, FadeFrom: 0.1, FadeTo: 1, ScaleMulFrom: 1.15, ScaleMulTo: 0.55,
		})
	}
	for i := 0; i < nEmber; i++ {
		life := 30 + g.rng.Float64()*50
		a := -math.Pi/2 + (g.rng.Float64()-0.5)*(1.0+power)
		sp := 1.2 + g.rng.Float64()*(3+power*4)
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX + (g.rng.Float64()-0.5)*20, Y: castY - 8,
			VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: 0.2 + g.rng.Float64()*0.5,
			Add: true, Grav: -0.035,
			Tint: color.RGBA{255, uint8(100 + g.rng.Intn(120)), 40, 255},
			Src:  vfxsprites.Spark, FadeFrom: 0.05, FadeTo: 1,
		})
	}
	g.rings = append(g.rings, ring{life: 28 + power*20, max: 28 + power*20, scale: 0.6 + power})
	g.flash = 0.35 + power*0.55
	g.heat = math.Min(1, g.heat+0.35+power*0.4)
}

func (g *game) Update() error {
	g.t += 0.08
	if g.flash > 0 {
		g.flash -= 0.045
	}
	if g.heat > 0 {
		g.heat -= 0.008
	}

	holding := false
	if x, y, ok := vfxui.Held(); ok && g.btn.Contains(x, y) {
		holding = true
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.Key1) {
		holding = true
	}
	if holding {
		g.charge = math.Min(1, g.charge+0.025)
		g.casting = true
		// Continuous ember trickle while charging.
		if g.rng.Float64() < 0.45+g.charge*0.4 {
			life := 25 + g.rng.Float64()*30
			g.parts = append(g.parts, vfxmagic.Particle{
				X: castX + (g.rng.Float64()-0.5)*28, Y: castY,
				VX: (g.rng.Float64() - 0.5) * 0.8, VY: -1 - g.rng.Float64()*2.5*g.charge,
				Life: life, Max: life, Scale: 0.25 + g.rng.Float64()*0.4*g.charge,
				Add: true, Grav: -0.03,
				Tint: color.RGBA{255, 140, 40, 255}, Src: vfxsprites.Spark,
			})
		}
	} else if g.casting {
		g.burst(0.35 + g.charge*0.65)
		g.charge = 0
		g.casting = false
	}
	if g.btn.Tapped() && !holding {
		// tap without hold still casts a medium burst on release handled above;
		// ensure instant cast if charge was near zero: handled by casting edge.
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame()
		return nil
	}

	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		// Flame sway.
		p.X += math.Sin(g.t*2.2+p.Y*0.05) * 0.4
		if p.Update() {
			alive = append(alive, p)
		}
	}
	g.parts = alive

	rings := g.rings[:0]
	for _, r := range g.rings {
		r.life--
		r.scale += 0.08
		if r.life > 0 {
			rings = append(rings, r)
		}
	}
	g.rings = rings
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{14, 8, 10, 255})
	// Warm floor gradient.
	vector.DrawFilledRect(s, 0, 520, width, 200, color.RGBA{28, 12, 8, 255}, false)
	vfxmagic.SoftFlash(s, g.flash, 255, 120, 40)
	if g.heat > 0 {
		vfxmagic.SoftFlash(s, g.heat*0.25, 255, 80, 20)
	}

	hero.DrawBottomCentered(s, castX, castY+55, 150)

	// Charge aura.
	if g.charge > 0.05 {
		sc := 0.8 + g.charge*1.6
		vfxmagic.DrawSprite(s, vfxsprites.Fire, castX, castY-10, g.t*0.3, sc*0.7, sc, color.RGBA{255, 160, 60, 255}, float32(0.35+g.charge*0.5), true)
		vfxmagic.DrawSprite(s, vfxsprites.Spark, castX, castY, 0, sc*1.4, sc*0.5, color.RGBA{255, 100, 30, 255}, float32(g.charge*0.6), true)
	}

	for _, r := range g.rings {
		f := r.life / r.max
		vfxmagic.DrawSprite(s, vfxsprites.Ring, castX, castY-20, 0, r.scale*1.8, r.scale*0.7,
			color.RGBA{255, uint8(80 + 100*f), 20, 255}, float32(0.55*f), true)
	}
	for i := range g.parts {
		g.parts[i].Draw(s)
	}

	ebitenutil.DebugPrintAt(s, "FIRE MAGIC — HOLD TO CHARGE, RELEASE TO BURST", 28, 18)
	ebitenutil.DebugPrintAt(s, "layers: flame PNG + embers + heat ring + screen flash", 16, 44)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("charge %.0f%%   particles %d   casts %d", g.charge*100, len(g.parts), g.casts), 70, 610)
	g.btn.Draw(s, g.casting || g.charge > 0.1)
	ebitenutil.DebugPrintAt(s, "HOLD BUTTON / SPACE   R = reset", 120, 700)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cool Fire Magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
