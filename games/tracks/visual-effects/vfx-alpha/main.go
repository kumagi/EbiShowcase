// vfx-alpha — Visual Effects Lab STEP 04.
// Translucency and draw order: a fading afterimage trail from stacked copies.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const width, height = 480, 720

type vec struct{ x, y float64 }

type game struct {
	pos      vec
	target   vec
	trail    []vec
	trailLen int
	active   int
	rings    [2]vec
	combo    int
	buttons  []vfxui.Button
	clear    bool
}

func newGame() *game {
	g := &game{
		pos:      vec{width / 2, height / 2},
		target:   vec{width / 2, height / 2},
		trailLen: 14,
		rings:    [2]vec{{100, 250}, {380, 470}},
	}
	w := 130.0
	gap := 12.0
	x := (width - (w*2 + gap)) / 2
	for _, l := range []string{"TRAIL +", "TRAIL -"} {
		g.buttons = append(g.buttons, vfxui.Button{X: x, Y: 636, W: w, H: 54, Label: l})
		x += w + gap
	}
	return g
}

func (g *game) Update() error {
	if g.clear {
		if vfxui.AnyPressStart() {
			*g = *newGame()
		}
		return nil
	}
	if x, y, ok := vfxui.Held(); ok {
		if y < 620 {
			g.target = vec{x, y}
		}
	}
	if g.buttons[0].Tapped() {
		g.trailLen = int(math.Min(28, float64(g.trailLen)+2))
	}
	if g.buttons[1].Tapped() {
		g.trailLen = int(math.Max(2, float64(g.trailLen)-2))
	}
	// Ease toward the target so motion (and the trail) stays smooth.
	g.pos.x += (g.target.x - g.pos.x) * 0.2
	g.pos.y += (g.target.y - g.pos.y) * 0.2
	g.trail = append(g.trail, g.pos)
	if len(g.trail) > g.trailLen {
		g.trail = g.trail[len(g.trail)-g.trailLen:]
	}
	r := g.rings[g.active]
	if math.Hypot(g.pos.x-r.x, g.pos.y-r.y) < 44 {
		g.active = 1 - g.active
		g.combo++
		if g.combo >= 6 {
			g.clear = true
		}
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{14, 22, 38, 255})

	// Glowing rings to visit alternately.
	for i, r := range g.rings {
		c := color.RGBA{70, 90, 130, 255}
		if i == g.active {
			c = color.RGBA{120, 240, 220, 255}
		}
		vector.StrokeCircle(s, float32(r.x), float32(r.y), 40, 4, c, false)
	}

	// Oldest copies are the most transparent: alpha grows toward the head.
	sprite := hero.Image()
	b := sprite.Bounds()
	sw, sh := float64(b.Dx()), float64(b.Dy())
	scale := 120.0 / sh
	for i, p := range g.trail {
		alpha := float64(i+1) / float64(len(g.trail))
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterLinear
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(p.x-sw*scale/2, p.y-sh*scale/2)
		op.ColorScale.ScaleAlpha(float32(alpha * 0.9))
		s.DrawImage(sprite, op)
	}

	ebitenutil.DebugPrintAt(s, "DRAG TENJIROH. OLD COPIES FADE OUT.", 108, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("op.ColorScale.ScaleAlpha(i/len)   TRAIL %d", g.trailLen), 16, 52)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("VISIT RINGS: %d/6", g.combo), 16, 74)
	ebitenutil.DebugPrintAt(s, "MORE / FEWER AFTERIMAGES:", 60, 604)
	for i := range g.buttons {
		g.buttons[i].Draw(s, false)
	}
	if g.clear {
		overlay(s, "SIX RINGS, ONE SWOOSH!\n\nSTACKED TRANSPARENT COPIES\n= MOTION BLUR.\nTAP / SPACE TO RESET")
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 250, 370, 170, color.RGBA{8, 16, 32, 240}, false)
	vector.StrokeRect(s, 55, 250, 370, 170, 3, color.RGBA{120, 240, 220, 255}, false)
	ebitenutil.DebugPrintAt(s, msg, 95, 285)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Alpha & Afterimage — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
