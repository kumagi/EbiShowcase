// vfx-stamp — Visual Effects Lab STEP 01.
// One image, drawn in many places with DrawImageOptions + GeoM.Translate.
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const width, height = 480, 720

type stamp struct{ x, y float64 }

type game struct {
	rng     *rand.Rand
	stamps  []stamp
	targetX float64
	targetY float64
	pointer struct{ x, y float64 }
	score   int
	clear   bool
	pulse   float64
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(7))}
	g.pointer.x, g.pointer.y = width/2, height/2
	g.newTarget()
	return g
}

func (g *game) newTarget() {
	g.targetX = 90 + g.rng.Float64()*300
	g.targetY = 150 + g.rng.Float64()*380
}

func (g *game) Update() error {
	g.pulse += 0.08
	if g.clear {
		if vfxui.AnyPressStart() {
			*g = *newGame()
		}
		return nil
	}
	if x, y, ok := vfxui.Held(); ok {
		g.pointer.x, g.pointer.y = x, y
	}
	if x, y, ok := vfxui.JustPressed(); ok {
		g.pointer.x, g.pointer.y = x, y
		g.stamps = append(g.stamps, stamp{x, y})
		if math.Hypot(x-g.targetX, y-g.targetY) < 46 {
			g.score++
			g.newTarget()
			if g.score >= 6 {
				g.clear = true
			}
		}
		if len(g.stamps) > 40 {
			g.stamps = g.stamps[len(g.stamps)-40:]
		}
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{16, 24, 40, 255})
	for x := 0; x <= width; x += 48 {
		vector.StrokeLine(s, float32(x), 90, float32(x), 660, 1, color.RGBA{30, 42, 66, 255}, false)
	}
	for y := 90; y <= 660; y += 48 {
		vector.StrokeLine(s, 0, float32(y), width, float32(y), 1, color.RGBA{30, 42, 66, 255}, false)
	}

	// Target ring.
	r := float32(38 + math.Sin(g.pulse)*6)
	vector.StrokeCircle(s, float32(g.targetX), float32(g.targetY), r, 4, color.RGBA{120, 240, 220, 255}, false)

	// Faint history: the SAME image drawn in many places.
	for _, st := range g.stamps {
		hero.DrawCentered(s, st.x, st.y, 46)
	}
	// Live pointer preview.
	hero.DrawCentered(s, g.pointer.x, g.pointer.y, 64)

	ebitenutil.DebugPrintAt(s, fmt.Sprintf("op.GeoM.Translate(%.0f, %.0f)", g.pointer.x, g.pointer.y), 16, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAMPS %d   ON TARGET %d/6", len(g.stamps), g.score), 16, 46)
	ebitenutil.DebugPrintAt(s, "TAP / CLICK TO STAMP TENJIROH ON THE RING", 90, 684)

	if g.clear {
		overlay(s, "SIX STAMPS ON TARGET!\n\nONE IMAGE, MANY PLACES.\nTAP / SPACE TO RESET")
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 285, 370, 150, color.RGBA{8, 16, 32, 240}, false)
	vector.StrokeRect(s, 55, 285, 370, 150, 3, color.RGBA{120, 240, 220, 255}, false)
	ebitenutil.DebugPrintAt(s, msg, 95, 320)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Stamp & Move — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
