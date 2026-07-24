// Package vfxfx is a tiny effect layer for Visual Effects Lab "advanced" demos.
// Gameplay state stays in each lesson's Game; particles / flashes / rings live here
// so Update/Draw can call fx.Update() and fx.Draw() as a separate concern.
package vfxfx

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxsprites"
)

// Particle is one short-lived sprite or soft glow.
type Particle struct {
	X, Y, VX, VY float64
	Life, Max    float64
	Scale        float64 // size at birth
	EndScale     float64 // size at death; zero means a small fade-out
	Rotation     float64
	Spin         float64
	Drag         float64
	Grav         float64
	Add          bool
	Tint         color.Color
	Src          *ebiten.Image
}

// System owns all transient visuals for a demo.
type System struct {
	Parts      []Particle
	Flash      float64
	FR, FG, FB float64
	MaxParts   int
}

const defaultMaxParts = 900

func (s *System) particleLimit() int {
	if s.MaxParts > 0 {
		return s.MaxParts
	}
	return defaultMaxParts
}

func (s *System) spawn(p Particle) {
	if len(s.Parts) >= s.particleLimit() {
		return
	}
	if p.Max <= 0 {
		p.Max = p.Life
	}
	if p.Life <= 0 {
		return
	}
	if p.Drag == 0 {
		p.Drag = 0.99
	}
	s.Parts = append(s.Parts, p)
}

// Update advances every effect. Call once per tick from Game.Update, after play logic.
func (s *System) Update() {
	if s.Flash > 0 {
		s.Flash -= 0.05
		if s.Flash < 0 {
			s.Flash = 0
		}
	}
	alive := s.Parts[:0]
	for _, p := range s.Parts {
		p.X += p.VX
		p.Y += p.VY
		p.VY += p.Grav
		p.VX *= p.Drag
		p.Rotation += p.Spin
		p.Life--
		if p.Life > 0 {
			alive = append(alive, p)
		}
	}
	s.Parts = alive
}

// Draw paints the screen flash (behind) then all particles. Call after the world.
func (s *System) Draw(dst *ebiten.Image) {
	if s.Flash > 0 {
		a := s.Flash
		if a > 1 {
			a = 1
		}
		r := s.FR
		g := s.FG
		b := s.FB
		if r == 0 && g == 0 && b == 0 {
			r, g, b = 200, 220, 255
		}
		vector.DrawFilledRect(dst, 0, 0, float32(dst.Bounds().Dx()), float32(dst.Bounds().Dy()),
			color.RGBA{uint8(r * a), uint8(g * a), uint8(b * a), uint8(180 * a)}, false)
	}
	for i := range s.Parts {
		s.drawPart(dst, &s.Parts[i])
	}
}

func (s *System) drawPart(dst *ebiten.Image, p *Particle) {
	src := p.Src
	if src == nil {
		src = vfxsprites.Spark
	}
	f := p.Life / math.Max(1, p.Max)
	progress := 1 - f
	endScale := p.EndScale
	if endScale == 0 {
		endScale = p.Scale * 0.18
	}
	scale := p.Scale + (endScale-p.Scale)*(1-math.Pow(1-progress, 3))
	b := src.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(p.Rotation)
	op.GeoM.Translate(p.X, p.Y)
	tint := p.Tint
	if tint == nil {
		tint = color.White
	}
	op.ColorScale.ScaleWithColor(tint)
	fadeIn := math.Min(1, progress*5)
	op.ColorScale.ScaleAlpha(float32(f * fadeIn))
	if p.Add {
		op.Blend = ebiten.BlendLighter
	}
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(src, op)
}

// Burst spawns n sparks from (x,y).
func (s *System) Burst(x, y float64, n int, speed float64, tint color.Color, add bool) {
	for i := 0; i < n; i++ {
		a := float64(i)/float64(n)*2*math.Pi + float64(i%7)*0.17
		sp := speed * (0.4 + 0.6*float64((i*37)%10)/10)
		life := 22 + float64((i*13)%30)
		scale := 0.35 + float64(i%5)*0.08
		s.spawn(Particle{
			X: x, Y: y, VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: scale, EndScale: scale * 0.12,
			Rotation: a, Spin: float64(i%5-2) * 0.035,
			Add: add, Tint: tint, Src: vfxsprites.Spark, Grav: 0.05, Drag: 0.985,
		})
	}
}

// Ring expands a soft ring sprite once.
func (s *System) Ring(x, y, scale float64, tint color.Color) {
	s.spawn(Particle{
		X: x, Y: y, Life: 28, Max: 28,
		Scale: scale * 0.18, EndScale: scale * 2.2,
		Add: true, Tint: tint, Src: vfxsprites.Ring,
	})
}

// Shockwave layers two expanding rings with different lifetimes. It is still
// presentation-only: callers must resolve collisions before requesting it.
func (s *System) Shockwave(x, y, scale float64, inner, outer color.Color) {
	s.spawn(Particle{
		X: x, Y: y, Life: 22, Max: 22,
		Scale: scale * 0.12, EndScale: scale * 2.0,
		Add: true, Tint: inner, Src: vfxsprites.Ring,
	})
	s.spawn(Particle{
		X: x, Y: y, Life: 38, Max: 38,
		Scale: scale * 0.08, EndScale: scale * 3.0,
		Add: true, Tint: outer, Src: vfxsprites.Ring,
	})
}

// Dust emits a low, sideways fan useful for landings and pushes.
func (s *System) Dust(x, y, direction float64, n int, tint color.Color) {
	for i := 0; i < n; i++ {
		side := -1.0
		if i%2 == 1 {
			side = 1
		}
		spread := 0.4 + float64((i*17)%10)/10
		life := 20 + float64(i%16)
		scale := 0.28 + float64(i%4)*0.07
		s.spawn(Particle{
			X: x + side*float64(i%4)*2, Y: y,
			VX:   side*(0.7+spread*1.8) + direction*0.25,
			VY:   -0.5 - spread*1.3,
			Life: life, Max: life, Scale: scale, EndScale: scale * 0.45,
			Grav: 0.07, Drag: 0.94, Tint: tint, Src: vfxsprites.Spark,
		})
	}
}

// Confetti creates a capped celebratory shower.
func (s *System) Confetti(x, y float64, n int) {
	palette := []color.Color{
		color.RGBA{255, 110, 110, 255},
		color.RGBA{255, 220, 90, 255},
		color.RGBA{80, 235, 200, 255},
		color.RGBA{130, 170, 255, 255},
	}
	for i := 0; i < n; i++ {
		a := -math.Pi*0.9 + math.Pi*0.8*float64((i*29)%100)/100
		speed := 2.2 + float64((i*31)%25)/5
		life := 55 + float64(i%35)
		scale := 0.26 + float64(i%3)*0.08
		s.spawn(Particle{
			X: x, Y: y,
			VX: math.Cos(a) * speed, VY: math.Sin(a) * speed,
			Life: life, Max: life, Scale: scale, EndScale: scale * 0.6,
			Rotation: a, Spin: float64(i%7-3) * 0.08,
			Grav: 0.11, Drag: 0.995, Tint: palette[i%len(palette)], Src: vfxsprites.Spark,
		})
	}
}

// FlashScreen sets a brief full-screen tint.
func (s *System) FlashScreen(amount, r, g, b float64) {
	s.Flash = amount
	s.FR, s.FG, s.FB = r, g, b
}

// FlameBurst rising fire tongues (for dramatic clears).
func (s *System) FlameBurst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		life := 35 + float64(i%20)
		scale := 0.45 + float64(i%4)*0.1
		s.spawn(Particle{
			X: x + float64(i%7-3)*4, Y: y,
			VX: float64(i%5-2) * 0.4, VY: -1.5 - float64(i%6)*0.35,
			Life: life, Max: life, Scale: scale, EndScale: scale * 1.35,
			Add: true, Grav: -0.02, Drag: 0.98,
			Spin: float64(i%5-2) * 0.02,
			Tint: color.RGBA{255, 120, 40, 255}, Src: vfxsprites.Fire,
		})
	}
}

// Count returns active particle count (handy for LIVE GO panels).
func (s *System) Count() int { return len(s.Parts) }
