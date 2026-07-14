// Package shaderlab provides small, optional Kage effects. Every caller can
// keep drawing when shader creation is unavailable (for example old GPUs).
package shaderlab

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed pulse.kage
var pulseSource []byte

//go:embed palette.kage
var paletteSource []byte

//go:embed distort.kage
var distortSource []byte

//go:embed status.kage
var statusSource []byte

type Pulse struct{ shader *ebiten.Shader }

func NewPulse() *Pulse {
	s, err := ebiten.NewShader(pulseSource)
	if err != nil {
		return &Pulse{}
	}
	return &Pulse{shader: s}
}
func (p *Pulse) Available() bool { return p != nil && p.shader != nil }

// Draw returns true when Kage rendered the image. Callers draw their normal
// image on false, so a visual enhancement never makes a game unplayable.
func (p *Pulse) Draw(dst, src *ebiten.Image, time float32) bool {
	if !p.Available() {
		return false
	}
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	op.Uniforms = map[string]any{"Time": time}
	dst.DrawRectShader(src.Bounds().Dx(), src.Bounds().Dy(), p.shader, op)
	return true
}

// Palette demonstrates named uniforms: Time drives color, Flash adds a brief
// hit response. It has the same opt-in contract as Pulse.
type Palette struct{ shader *ebiten.Shader }

func NewPalette() *Palette {
	s, err := ebiten.NewShader(paletteSource)
	if err != nil {
		return &Palette{}
	}
	return &Palette{shader: s}
}
func (p *Palette) Available() bool { return p != nil && p.shader != nil }
func (p *Palette) Draw(dst, src *ebiten.Image, time, flash float32) bool {
	if !p.Available() {
		return false
	}
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	op.Uniforms = map[string]any{"Time": time, "Flash": flash}
	dst.DrawRectShader(src.Bounds().Dx(), src.Bounds().Dy(), p.shader, op)
	return true
}

// Distort keeps effect selection in uniforms: water, heat, and an impact wave
// are all coordinate changes before sampling the same source image.
type Distort struct{ shader *ebiten.Shader }

func NewDistort() *Distort {
	s, err := ebiten.NewShader(distortSource)
	if err != nil {
		return &Distort{}
	}
	return &Distort{shader: s}
}
func (d *Distort) Available() bool { return d != nil && d.shader != nil }
func (d *Distort) Draw(dst, src *ebiten.Image, time, mode, impact float32) bool {
	if !d.Available() {
		return false
	}
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	op.Uniforms = map[string]any{"Time": time, "Mode": mode, "Impact": impact}
	dst.DrawRectShader(src.Bounds().Dx(), src.Bounds().Dy(), d.shader, op)
	return true
}

// Status separates feedback state (poison/freeze/damage) from the game rule
// that caused it. The renderer only receives display strengths.
type Status struct{ shader *ebiten.Shader }

func NewStatus() *Status {
	s, err := ebiten.NewShader(statusSource)
	if err != nil {
		return &Status{}
	}
	return &Status{shader: s}
}
func (s *Status) Available() bool { return s != nil && s.shader != nil }
func (s *Status) Draw(dst, src *ebiten.Image, status, damage float32) bool {
	if !s.Available() {
		return false
	}
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	op.Uniforms = map[string]any{"Status": status, "Damage": damage}
	dst.DrawRectShader(src.Bounds().Dx(), src.Bounds().Dy(), s.shader, op)
	return true
}
