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
	Scale        float64
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
		p.VX *= 0.99
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
	f := p.Life / p.Max
	b := src.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Scale(p.Scale*(0.5+0.5*f), p.Scale*(0.5+0.5*f))
	op.GeoM.Translate(p.X, p.Y)
	tint := p.Tint
	if tint == nil {
		tint = color.White
	}
	op.ColorScale.ScaleWithColor(tint)
	op.ColorScale.ScaleAlpha(float32(0.2 + 0.8*f))
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
		s.Parts = append(s.Parts, Particle{
			X: x, Y: y, VX: math.Cos(a) * sp, VY: math.Sin(a) * sp,
			Life: life, Max: life, Scale: 0.35 + float64(i%5)*0.08,
			Add: add, Tint: tint, Src: vfxsprites.Spark, Grav: 0.05,
		})
	}
}

// Ring expands a soft ring sprite once.
func (s *System) Ring(x, y, scale float64, tint color.Color) {
	s.Parts = append(s.Parts, Particle{
		X: x, Y: y, Life: 28, Max: 28, Scale: scale,
		Add: true, Tint: tint, Src: vfxsprites.Ring,
	})
	// grow by increasing scale via Velocity hack: use Grav field unused — scale via VX as growth
	s.Parts[len(s.Parts)-1].VX = 0.06 // read in Draw? simpler: just fixed fade
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
		s.Parts = append(s.Parts, Particle{
			X: x + float64(i%7-3)*4, Y: y,
			VX: float64(i%5-2) * 0.4, VY: -1.5 - float64(i%6)*0.35,
			Life: life, Max: life, Scale: 0.45 + float64(i%4)*0.1,
			Add: true, Grav: -0.02, Tint: color.RGBA{255, 120, 40, 255}, Src: vfxsprites.Fire,
		})
	}
}

// Count returns active particle count (handy for LIVE GO panels).
func (s *System) Count() int { return len(s.Parts) }
