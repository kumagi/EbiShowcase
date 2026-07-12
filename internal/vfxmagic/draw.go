// Package vfxmagic holds shared drawing helpers for the advanced Visual Effects
// Lab magic showcases (fire / ice / thunder / light / dark).
package vfxmagic

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle is a textured sprite with motion, spin, and lifetime.
type Particle struct {
	X, Y, VX, VY float64
	Rot, Spin    float64
	Life, Max    float64
	Scale        float64
	Add          bool
	Grav         float64
	Tint         color.Color
	Src          *ebiten.Image
	// FadeFrom/FadeTo control alpha over life (defaults 0.2→1 if zero).
	FadeFrom, FadeTo float32
	// ScaleMulFrom/ScaleMulTo multiply Scale over life (defaults 1→1).
	ScaleMulFrom, ScaleMulTo float64
}

// Update advances one particle; returns false when dead.
func (p *Particle) Update() bool {
	p.X += p.VX
	p.Y += p.VY
	p.VY += p.Grav
	p.Rot += p.Spin
	p.Life--
	return p.Life > 0
}

// Draw paints the particle centered at (X,Y).
func (p *Particle) Draw(dst *ebiten.Image) {
	if p.Src == nil || p.Life <= 0 {
		return
	}
	f := p.Life / p.Max
	sf, st := p.ScaleMulFrom, p.ScaleMulTo
	if sf == 0 && st == 0 {
		sf, st = 1, 1
	}
	scale := p.Scale * (st + (sf-st)*f)

	ff, ft := p.FadeFrom, p.FadeTo
	if ff == 0 && ft == 0 {
		ff, ft = 0.15, 1
	}
	alpha := ff + (ft-ff)*float32(f)

	b := p.Src.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(p.Rot)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(p.X, p.Y)
	tint := p.Tint
	if tint == nil {
		tint = color.White
	}
	op.ColorScale.ScaleWithColor(tint)
	op.ColorScale.ScaleAlpha(alpha)
	if p.Add {
		op.Blend = ebiten.BlendLighter
	}
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(p.Src, op)
}

// SoftFlash fills the screen with a fading additive-ish tint.
func SoftFlash(dst *ebiten.Image, amount float64, r, g, b uint8) {
	if amount <= 0 {
		return
	}
	a := amount
	if a > 1 {
		a = 1
	}
	vector.DrawFilledRect(dst, 0, 0, float32(dst.Bounds().Dx()), float32(dst.Bounds().Dy()),
		color.RGBA{uint8(float64(r) * a), uint8(float64(g) * a), uint8(float64(b) * a), uint8(200 * a)}, false)
}

// DrawSprite draws a centered textured sprite with optional additive blend.
func DrawSprite(dst, src *ebiten.Image, x, y, rot, sx, sy float64, tint color.Color, alpha float32, add bool) {
	if src == nil {
		return
	}
	b := src.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Rotate(rot)
	op.GeoM.Scale(sx, sy)
	op.GeoM.Translate(x, y)
	if tint == nil {
		tint = color.White
	}
	op.ColorScale.ScaleWithColor(tint)
	op.ColorScale.ScaleAlpha(alpha)
	if add {
		op.Blend = ebiten.BlendLighter
	}
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(src, op)
}

// JaggedBolt draws a multi-segment lightning path from (x0,y0) toward (x1,y1).
func JaggedBolt(dst *ebiten.Image, x0, y0, x1, y1 float64, segs int, wobble float64, thickness float32, col color.Color, rng func() float64) {
	if segs < 2 {
		segs = 2
	}
	px, py := x0, y0
	for i := 1; i <= segs; i++ {
		t := float64(i) / float64(segs)
		nx := x0 + (x1-x0)*t
		ny := y0 + (y1-y0)*t
		if i < segs {
			nx += (rng() - 0.5) * wobble
			ny += (rng() - 0.5) * wobble * 0.4
		}
		vector.StrokeLine(dst, float32(px), float32(py), float32(nx), float32(ny), thickness, col, false)
		px, py = nx, ny
	}
}

// SpiralOffset returns a point on a growing spiral.
func SpiralOffset(t, arms, radius, twist float64) (dx, dy float64) {
	ang := t*twist + arms*math.Pi*2*t
	r := radius * t
	return math.Cos(ang) * r, math.Sin(ang) * r
}
