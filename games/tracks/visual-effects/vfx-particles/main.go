// vfx-particles — Visual Effects Lab STEP 07.
// A particle system: a slice of short-lived dots that spawn, move, and fade.
package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const width, height = 480, 720

const glowR = 32

type particle struct {
	x, y, vx, vy float64
	life, max    float64
	col          color.RGBA
}

type game struct {
	glow     *ebiten.Image
	rng      *rand.Rand
	parts    []particle
	gravity  bool
	additive bool
	bursts   int
	buttons  []vfxui.Button
	clear    bool
}

func makeGlow() *ebiten.Image {
	size := glowR * 2
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			d := math.Hypot(float64(x-glowR)+0.5, float64(y-glowR)+0.5) / glowR
			a := 0.0
			if d < 1 {
				a = (1 - d) * (1 - d)
			}
			img.SetRGBA(x, y, color.RGBA{255, 255, 255, uint8(a * 255)})
		}
	}
	return ebiten.NewImageFromImage(img)
}

func newGame() *game {
	g := &game{glow: makeGlow(), rng: rand.New(rand.NewSource(11)), gravity: true, additive: true}
	w := 130.0
	gap := 12.0
	x := (width - (w*2 + gap)) / 2
	for _, l := range []string{"GRAVITY: ON", "ADD: ON"} {
		g.buttons = append(g.buttons, vfxui.Button{X: x, Y: 636, W: w, H: 54, Label: l})
		x += w + gap
	}
	return g
}

func (g *game) burst(cx, cy float64) {
	hue := g.rng.Float64()
	for i := 0; i < 28; i++ {
		a := g.rng.Float64() * 2 * math.Pi
		sp := 1.5 + g.rng.Float64()*5
		life := 40 + g.rng.Float64()*40
		g.parts = append(g.parts, particle{
			x: cx, y: cy,
			vx: math.Cos(a) * sp, vy: math.Sin(a)*sp - 2,
			life: life, max: life,
			col: sparkColor(hue + g.rng.Float64()*0.15),
		})
	}
	g.bursts++
	if g.bursts >= 6 {
		g.clear = true
	}
}

func sparkColor(h float64) color.RGBA {
	h = math.Mod(h, 1)
	switch {
	case h < 0.33:
		return color.RGBA{255, 150, 70, 255}
	case h < 0.66:
		return color.RGBA{90, 220, 255, 255}
	default:
		return color.RGBA{200, 130, 255, 255}
	}
}

func (g *game) Update() error {
	if g.clear {
		if vfxui.AnyPressStart() {
			*g = *newGame()
		}
		return nil
	}
	if x, y, ok := vfxui.JustPressed(); ok {
		if g.buttons[0].Contains(x, y) {
			g.gravity = !g.gravity
			g.buttons[0].Label = label("GRAVITY", g.gravity)
		} else if g.buttons[1].Contains(x, y) {
			g.additive = !g.additive
			g.buttons[1].Label = label("ADD", g.additive)
		} else if y < 620 {
			g.burst(x, y)
		}
	}
	alive := g.parts[:0]
	for _, p := range g.parts {
		p.x += p.vx
		p.y += p.vy
		if g.gravity {
			p.vy += 0.14
		}
		p.vx *= 0.99
		p.life--
		if p.life > 0 {
			alive = append(alive, p)
		}
	}
	g.parts = alive
	return nil
}

func label(name string, on bool) string {
	if on {
		return name + ": ON"
	}
	return name + ": OFF"
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 10, 20, 255})
	for _, p := range g.parts {
		f := p.life / p.max
		scale := 0.4 + f*1.1
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(p.x-glowR*scale, p.y-glowR*scale)
		op.ColorScale.ScaleWithColor(p.col)
		op.ColorScale.ScaleAlpha(float32(f))
		if g.additive {
			op.Blend = ebiten.BlendLighter
		}
		s.DrawImage(g.glow, op)
	}
	ebitenutil.DebugPrintAt(s, "TAP TO BURST SPARKS", 156, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("particles alive: %d   bursts %d/6", len(g.parts), g.bursts), 16, 52)
	ebitenutil.DebugPrintAt(s, "EACH SPARK: pos += vel, life--, fade by life:", 44, 604)
	for i := range g.buttons {
		on := (i == 0 && g.gravity) || (i == 1 && g.additive)
		g.buttons[i].Draw(s, on)
	}
	if g.clear {
		overlay(s, "SIX BURSTS!\n\nONE SLICE OF DOTS =\nSMOKE, FIRE, SPARKLE.\nTAP / SPACE TO RESET")
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 250, 370, 170, color.RGBA{8, 12, 24, 245}, false)
	vector.StrokeRect(s, 55, 250, 370, 170, 3, color.RGBA{120, 240, 220, 255}, false)
	ebitenutil.DebugPrintAt(s, msg, 95, 285)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Particle Burst — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
