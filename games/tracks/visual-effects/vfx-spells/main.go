// vfx-spells — Visual Effects Lab STEP 08 (capstone).
// Combine transforms, color, alpha, additive blending, and textured sprites into
// three toy spells cast by Ebi Tenjiroh: fire, water, and lightning.
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
	"github.com/kumagi/EbiShowcase/internal/vfxsprites"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const width, height = 480, 720

const (
	castX = 240.0
	castY = 470.0
)

type pKind int

const (
	pFlame pKind = iota
	pEmber
	pDrop
	pSplash
	pSpark
)

type particle struct {
	x, y, vx, vy float64
	rot, spin    float64
	life, max    float64
	scale        float64
	kind         pKind
	add          bool
	grav         float64
}

type bolt struct {
	x, y, rot, scale float64
	life, max        int
}

type game struct {
	rng      *rand.Rand
	parts    []particle
	bolts    []bolt
	flash    float64
	castMask int
	buttons  []vfxui.Button
	t        float64
	clear    bool
	banner   string
	bannerT  float64
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(3))}
	w := 130.0
	gap := 12.0
	x := (width - (w*3 + gap*2)) / 2
	for _, l := range []string{"FIRE", "WATER", "THUNDER"} {
		g.buttons = append(g.buttons, vfxui.Button{X: x, Y: 636, W: w, H: 54, Label: l})
		x += w + gap
	}
	return g
}

func (g *game) announce(msg string) {
	g.banner = msg
	g.bannerT = 2.2
}

func (g *game) cast(kind int) {
	g.castMask |= 1 << kind
	switch kind {
	case 0: // Fire plume: flame tongues + rising embers (additive).
		for i := 0; i < 28; i++ {
			life := 55 + g.rng.Float64()*50
			g.parts = append(g.parts, particle{
				x: castX + (g.rng.Float64()-0.5)*36, y: castY + g.rng.Float64()*8,
				vx:  (g.rng.Float64() - 0.5) * 1.4,
				vy:  -1.6 - g.rng.Float64()*2.8,
				rot: (g.rng.Float64() - 0.5) * 0.4, spin: (g.rng.Float64() - 0.5) * 0.04,
				life: life, max: life, scale: 0.55 + g.rng.Float64()*0.7,
				kind: pFlame, add: true, grav: -0.015,
			})
		}
		for i := 0; i < 40; i++ {
			life := 35 + g.rng.Float64()*40
			a := -math.Pi/2 + (g.rng.Float64()-0.5)*1.2
			sp := 1.5 + g.rng.Float64()*3.2
			g.parts = append(g.parts, particle{
				x: castX + (g.rng.Float64()-0.5)*24, y: castY,
				vx: math.Cos(a) * sp, vy: math.Sin(a) * sp,
				life: life, max: life, scale: 0.25 + g.rng.Float64()*0.45,
				kind: pEmber, add: true, grav: -0.04,
			})
		}
		g.flash = 0.55
		g.announce("FIRE = flame sprite + additive + rise")
	case 1: // Water: textured droplets under gravity + splash sparks.
		for i := 0; i < 36; i++ {
			life := 70 + g.rng.Float64()*45
			a := -math.Pi/2 + (g.rng.Float64()-0.5)*1.5
			sp := 3.5 + g.rng.Float64()*4.5
			g.parts = append(g.parts, particle{
				x: castX + (g.rng.Float64()-0.5)*10, y: castY - 10,
				vx: math.Cos(a) * sp, vy: math.Sin(a) * sp,
				rot: g.rng.Float64() * math.Pi, spin: (g.rng.Float64() - 0.5) * 0.12,
				life: life, max: life, scale: 0.45 + g.rng.Float64()*0.55,
				kind: pDrop, add: false, grav: 0.24,
			})
		}
		for i := 0; i < 20; i++ {
			life := 28 + g.rng.Float64()*24
			g.parts = append(g.parts, particle{
				x: castX + (g.rng.Float64()-0.5)*40, y: castY - 20,
				vx: (g.rng.Float64() - 0.5) * 5, vy: -1 - g.rng.Float64()*2,
				life: life, max: life, scale: 0.2 + g.rng.Float64()*0.3,
				kind: pSplash, add: true, grav: 0.18,
			})
		}
		g.flash = 0.28
		g.announce("WATER = droplet sprite + alpha + gravity")
	case 2: // Lightning: bolt sprites + white flash + sparks.
		x := castX + (g.rng.Float64()-0.5)*90
		y := 70.0
		for y < castY-20 {
			life := 16 + g.rng.Intn(8)
			ang := (g.rng.Float64() - 0.5) * 0.55
			sc := 0.55 + g.rng.Float64()*0.55
			g.bolts = append(g.bolts, bolt{x: x, y: y, rot: ang, scale: sc, life: life, max: life})
			y += 28 + g.rng.Float64()*18
			x += (g.rng.Float64() - 0.5) * 55
		}
		// Final bolt aimed at caster hand.
		g.bolts = append(g.bolts, bolt{x: castX, y: castY - 50, rot: 0, scale: 0.7, life: 18, max: 18})
		g.flash = 1
		for i := 0; i < 48; i++ {
			life := 24 + g.rng.Float64()*28
			a := g.rng.Float64() * 2 * math.Pi
			sp := 2 + g.rng.Float64()*5
			g.parts = append(g.parts, particle{
				x: castX + (g.rng.Float64()-0.5)*30, y: castY - 40 + (g.rng.Float64()-0.5)*40,
				vx: math.Cos(a) * sp, vy: math.Sin(a) * sp,
				life: life, max: life, scale: 0.2 + g.rng.Float64()*0.4,
				kind: pSpark, add: true, grav: 0.04,
			})
		}
		g.announce("THUNDER = bolt sprite + flash + sparks")
	}
	if g.castMask == 0b111 {
		g.clear = true
	}
}

func (g *game) Update() error {
	g.t += 0.1
	if g.flash > 0 {
		g.flash -= 0.06
	}
	if g.bannerT > 0 {
		g.bannerT -= 1.0 / 60
	}
	if g.clear {
		if vfxui.AnyPressStart() {
			*g = *newGame()
		}
		return nil
	}
	if x, y, ok := vfxui.JustPressed(); ok {
		for i := range g.buttons {
			if g.buttons[i].Contains(x, y) {
				g.cast(i)
			}
		}
	}
	for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if inpututil.IsKeyJustPressed(key) {
			g.cast(i)
		}
	}

	alive := g.parts[:0]
	for _, p := range g.parts {
		p.x += p.vx
		p.y += p.vy
		p.vy += p.grav
		p.rot += p.spin
		// Flame sway.
		if p.kind == pFlame {
			p.x += math.Sin(g.t*2.5+p.y*0.04) * 0.35
		}
		p.life--
		if p.life > 0 {
			alive = append(alive, p)
		}
	}
	g.parts = alive

	bolts := g.bolts[:0]
	for _, b := range g.bolts {
		b.life--
		if b.life > 0 {
			bolts = append(bolts, b)
		}
	}
	g.bolts = bolts
	return nil
}

func (g *game) drawPart(s *ebiten.Image, p particle) {
	f := p.life / p.max
	var src *ebiten.Image
	var tint color.Color = color.White
	scale := p.scale
	switch p.kind {
	case pFlame:
		src = vfxsprites.Fire
		scale *= 0.55 + 0.55*f
		// Cool toward red as it dies.
		tint = color.RGBA{255, uint8(80 + 140*f), uint8(20 + 40*f), 255}
	case pEmber:
		src = vfxsprites.Spark
		scale *= 0.4 + 0.8*f
		tint = color.RGBA{255, uint8(90 + 120*f), 30, 255}
	case pDrop:
		src = vfxsprites.Water
		scale *= 0.5 + 0.5*f
		tint = color.RGBA{180, 220, 255, 255}
	case pSplash:
		src = vfxsprites.Spark
		scale *= 0.35 + 0.5*f
		tint = color.RGBA{120, 200, 255, 255}
	case pSpark:
		src = vfxsprites.Spark
		scale *= 0.35 + 0.7*f
		tint = color.RGBA{210, 230, 255, 255}
	}
	if src == nil {
		return
	}
	b := src.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(p.rot)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(p.x, p.y)
	op.ColorScale.ScaleWithColor(tint)
	alpha := float32(0.25 + 0.75*f)
	if !p.add {
		alpha *= 0.85
	}
	op.ColorScale.ScaleAlpha(alpha)
	if p.add {
		op.Blend = ebiten.BlendLighter
	}
	op.Filter = ebiten.FilterLinear
	s.DrawImage(src, op)
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 12, 24, 255})
	if g.flash > 0 {
		a := g.flash
		vector.DrawFilledRect(s, 0, 0, width, height, color.RGBA{uint8(140 * a), uint8(160 * a), uint8(220 * a), uint8(220 * a)}, false)
	}

	hero.DrawBottomCentered(s, castX, castY+60, 150)

	// Soft ground glow under caster hand.
	if g.castMask != 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2.2, 0.7)
		op.GeoM.Translate(castX-28*2.2, castY-10)
		op.ColorScale.Scale(1, 0.55, 0.2, 0.35)
		op.Blend = ebiten.BlendLighter
		s.DrawImage(vfxsprites.Spark, op)
	}

	for _, p := range g.parts {
		g.drawPart(s, p)
	}
	for _, b := range g.bolts {
		f := float64(b.life) / float64(b.max)
		src := vfxsprites.Bolt
		bb := src.Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(bb.Dx())/2, -float64(bb.Dy())/2)
		op.GeoM.Rotate(b.rot)
		op.GeoM.Scale(b.scale*(0.8+0.4*f), b.scale*(0.8+0.4*f))
		op.GeoM.Translate(b.x, b.y)
		op.ColorScale.Scale(1, 1, 1, float32(0.4+0.6*f))
		op.Blend = ebiten.BlendLighter
		op.Filter = ebiten.FilterLinear
		s.DrawImage(src, op)
	}

	ebitenutil.DebugPrintAt(s, "CAST A SPELL — TEXTURED FIRE / WATER / LIGHTNING", 40, 20)
	tips := []string{"FIRE  = flame PNG + additive rise", "WATER = drop PNG + alpha + gravity", "THUNDER = bolt PNG + flash + sparks"}
	for i, t := range tips {
		mark := "[ ]"
		if g.castMask&(1<<i) != 0 {
			mark = "[x]"
		}
		ebitenutil.DebugPrintAt(s, mark+" "+t, 16, 48+i*18)
	}
	if g.bannerT > 0 && g.banner != "" {
		ebitenutil.DebugPrintAt(s, g.banner, 50, 118)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("particles %d   TAP A SPELL (or keys 1-3)", len(g.parts)), 70, 604)
	for i := range g.buttons {
		g.buttons[i].Draw(s, g.castMask&(1<<i) != 0)
	}
	if g.clear {
		overlay(s, "SPELLBOOK COMPLETE!\n\nTEXTURE + BLEND + MOTION\n= FIRE, WATER, LIGHTNING.\nTAP / SPACE TO RESET")
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 45, 250, 390, 170, color.RGBA{8, 12, 24, 245}, false)
	vector.StrokeRect(s, 45, 250, 390, 170, 3, color.RGBA{120, 240, 220, 255}, false)
	ebitenutil.DebugPrintAt(s, msg, 70, 285)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cast: Fire/Water/Lightning — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
