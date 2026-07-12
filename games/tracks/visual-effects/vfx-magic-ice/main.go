// vfx-magic-ice — Visual Effects Lab STEP 10.
// Showcase: crystal shatter — ice shards, frost ring, ground bloom, sparkle.
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
	castY = 470.0
)

type frost struct {
	x, y, rot, scale float64
	life, max        float64
}

type game struct {
	rng     *rand.Rand
	parts   []vfxmagic.Particle
	frosts  []frost
	flash   float64
	bloom   float64
	phase   int // 0 idle, 1 form, 2 shatter
	timer   float64
	crystal float64
	casts   int
	t       float64
	btn     vfxui.Button
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(10)),
		btn: vfxui.Button{X: 140, Y: 640, W: 200, H: 54, Label: "CAST ICE", Fill: color.RGBA{12, 48, 90, 235}},
	}
}

func (g *game) cast() {
	if g.phase != 0 {
		return
	}
	g.casts++
	g.phase = 1
	g.timer = 0
	g.crystal = 0
	g.flash = 0.4
	// Forming frost flakes around caster.
	for i := 0; i < 24; i++ {
		a := g.rng.Float64() * 2 * math.Pi
		r := 30 + g.rng.Float64()*70
		life := 40 + g.rng.Float64()*30
		g.frosts = append(g.frosts, frost{
			x: castX + math.Cos(a)*r, y: castY - 40 + math.Sin(a)*r*0.55,
			rot: a, scale: 0.35 + g.rng.Float64()*0.5, life: life, max: life,
		})
	}
}

func (g *game) shatter() {
	g.phase = 2
	g.timer = 0
	g.flash = 0.85
	g.bloom = 1
	// Outward crystal shards.
	for i := 0; i < 48; i++ {
		life := 55 + g.rng.Float64()*55
		a := g.rng.Float64() * 2 * math.Pi
		sp := 2.5 + g.rng.Float64()*6.5
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX + (g.rng.Float64()-0.5)*20, Y: castY - 50 + (g.rng.Float64()-0.5)*20,
			VX: math.Cos(a) * sp, VY: math.Sin(a)*sp*0.85 - 1.5,
			Rot: a, Spin: (g.rng.Float64() - 0.5) * 0.15,
			Life: life, Max: life, Scale: 0.35 + g.rng.Float64()*0.75,
			Add: false, Grav: 0.12,
			Tint: color.RGBA{uint8(180 + g.rng.Intn(60)), uint8(220 + g.rng.Intn(35)), 255, 255},
			Src:  vfxsprites.Ice, FadeFrom: 0.2, FadeTo: 1, ScaleMulFrom: 1.1, ScaleMulTo: 0.4,
		})
	}
	// Additive sparkle.
	for i := 0; i < 60; i++ {
		life := 25 + g.rng.Float64()*40
		a := g.rng.Float64() * 2 * math.Pi
		sp := 1 + g.rng.Float64()*5
		g.parts = append(g.parts, vfxmagic.Particle{
			X: castX, Y: castY - 50,
			VX: math.Cos(a) * sp, VY: math.Sin(a)*sp - 1,
			Life: life, Max: life, Scale: 0.15 + g.rng.Float64()*0.4,
			Add: true, Grav: 0.05,
			Tint: color.RGBA{200, 240, 255, 255}, Src: vfxsprites.Spark,
		})
	}
	g.frosts = append(g.frosts, frost{x: castX, y: castY - 40, rot: 0, scale: 0.5, life: 45, max: 45})
}

func (g *game) Update() error {
	g.t += 0.07
	if g.flash > 0 {
		g.flash -= 0.05
	}
	if g.bloom > 0 {
		g.bloom -= 0.012
	}

	if g.btn.Tapped() || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		if g.phase == 0 {
			g.cast()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame()
		return nil
	}

	switch g.phase {
	case 1:
		g.timer++
		g.crystal = math.Min(1, g.timer/36)
		// Pull frost inward.
		for i := range g.frosts {
			f := &g.frosts[i]
			f.x += (castX - f.x) * 0.06
			f.y += (castY - 50 - f.y) * 0.06
			f.rot += 0.08
			f.scale += 0.01
		}
		if g.timer >= 40 {
			g.shatter()
		}
	case 2:
		g.timer++
		if g.timer > 90 && len(g.parts) < 8 {
			g.phase = 0
			g.frosts = nil
		}
	}

	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		if p.Update() {
			alive = append(alive, p)
		}
	}
	g.parts = alive

	frosts := g.frosts[:0]
	for _, f := range g.frosts {
		f.life--
		if f.life > 0 {
			frosts = append(frosts, f)
		}
	}
	g.frosts = frosts
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 12, 28, 255})
	vector.DrawFilledRect(s, 0, 540, width, 180, color.RGBA{10, 24, 48, 255}, false)
	vfxmagic.SoftFlash(s, g.flash, 160, 210, 255)
	if g.bloom > 0 {
		vfxmagic.DrawSprite(s, vfxsprites.Ring, castX, castY-30, 0, 1.2+g.bloom*3.5, 0.55+g.bloom*1.2,
			color.RGBA{120, 200, 255, 255}, float32(0.55*g.bloom), true)
		vfxmagic.DrawSprite(s, vfxsprites.Light, castX, castY-50, g.t, 0.8+g.bloom*1.5, 0.8+g.bloom*1.5,
			color.RGBA{180, 230, 255, 255}, float32(0.4*g.bloom), true)
	}

	hero.DrawBottomCentered(s, castX, castY+60, 150)

	// Growing crystal core while forming.
	if g.phase == 1 || (g.phase == 2 && g.timer < 20) {
		sc := 0.4 + g.crystal*1.4
		a := float32(0.5 + g.crystal*0.5)
		if g.phase == 2 {
			a *= float32(1 - g.timer/25)
		}
		vfxmagic.DrawSprite(s, vfxsprites.Ice, castX, castY-55, g.t*0.4, sc, sc*1.15, color.White, a, false)
		vfxmagic.DrawSprite(s, vfxsprites.Ice, castX, castY-55, -g.t*0.55, sc*0.7, sc*0.7, color.RGBA{160, 220, 255, 255}, a*0.7, true)
	}

	for _, f := range g.frosts {
		ff := f.life / f.max
		vfxmagic.DrawSprite(s, vfxsprites.Ice, f.x, f.y, f.rot, f.scale*(0.6+0.4*ff), f.scale*(0.6+0.4*ff),
			color.RGBA{200, 240, 255, 255}, float32(0.35+0.55*ff), true)
	}
	for i := range g.parts {
		g.parts[i].Draw(s)
	}

	ebitenutil.DebugPrintAt(s, "ICE MAGIC — FORM CRYSTAL, THEN SHATTER", 48, 18)
	ebitenutil.DebugPrintAt(s, "layers: ice PNG + frost gather + shatter + bloom ring", 16, 44)
	phase := "idle"
	switch g.phase {
	case 1:
		phase = "forming"
	case 2:
		phase = "shatter"
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("phase %s   particles %d   casts %d", phase, len(g.parts), g.casts), 70, 610)
	g.btn.Draw(s, g.phase != 0)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE   R = reset", 150, 700)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cool Ice Magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
