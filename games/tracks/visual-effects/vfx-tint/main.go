// vfx-tint — Visual Effects Lab STEP 03.
// Change or drain color with op.ColorScale: tint, hit-flash, and silhouette.
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

const (
	modeNormal = iota
	modeTint
	modeFlash
	modeShadow
	modeCount
)

type game struct {
	mode    int
	hue     float64
	visited [modeCount]bool
	buttons []vfxui.Button
	clear   bool
	t       float64
}

func newGame() *game {
	g := &game{}
	g.visited[modeNormal] = true
	w := 104.0
	gap := 10.0
	x := (width - (w*4 + gap*3)) / 2
	for _, l := range []string{"NORMAL", "TINT", "FLASH", "SHADOW"} {
		g.buttons = append(g.buttons, vfxui.Button{X: x, Y: 636, W: w, H: 54, Label: l})
		x += w + gap
	}
	return g
}

func (g *game) setMode(m int) {
	g.mode = m
	g.visited[m] = true
	allSeen := true
	for _, v := range g.visited {
		if !v {
			allSeen = false
		}
	}
	if allSeen {
		g.clear = true
	}
}

func (g *game) Update() error {
	g.t += 0.12
	if g.clear {
		if vfxui.AnyPressStart() {
			*g = *newGame()
		}
		return nil
	}
	for i := range g.buttons {
		if g.buttons[i].Tapped() {
			g.setMode(i)
		}
	}
	if g.mode == modeTint {
		g.hue += 0.02
	}
	for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if ebiten.IsKeyPressed(key) {
			g.setMode(i)
		}
	}
	return nil
}

func hueRGB(h float64) (r, gc, b float64) {
	h = math.Mod(h, 1)
	switch {
	case h < 1.0/3:
		return 1, 0.35, 0.4
	case h < 2.0/3:
		return 0.4, 1, 0.6
	default:
		return 0.5, 0.6, 1
	}
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 26, 44, 255})
	cx, cy := float64(width/2), 340.0

	sprite := hero.Image()
	b := sprite.Bounds()
	sw, sh := float64(b.Dx()), float64(b.Dy())
	scale := 260.0 / sh
	tx := cx - sw*scale/2
	ty := cy - sh*scale/2

	desc := "op.ColorScale (none)"
	switch g.mode {
	case modeShadow:
		// Draw a color-drained copy first: RGB scaled to 0 = a shadow.
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterLinear
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(tx+18, ty+22)
		op.ColorScale.Scale(0, 0, 0, 0.55)
		s.DrawImage(sprite, op)
		desc = "op.ColorScale.Scale(0, 0, 0, 0.55)"
	}

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(tx, ty)
	switch g.mode {
	case modeTint:
		r, gg, bb := hueRGB(g.hue)
		op.ColorScale.Scale(float32(r), float32(gg), float32(bb), 1)
		desc = fmt.Sprintf("op.ColorScale.Scale(%.1f, %.1f, %.1f, 1)", r, gg, bb)
	case modeFlash:
		f := 3.0 + math.Abs(math.Sin(g.t))*4
		op.ColorScale.Scale(float32(f), float32(f), float32(f), 1)
		desc = "op.ColorScale.Scale(6, 6, 6, 1) // clamps to white"
	}
	s.DrawImage(sprite, op)

	ebitenutil.DebugPrintAt(s, "SAME PIXELS, RECOLORED BY ColorScale", 96, 24)
	ebitenutil.DebugPrintAt(s, desc, 16, 52)
	ebitenutil.DebugPrintAt(s, "TAP A MODE (or keys 1-4). SEE ALL 4 TO CLEAR.", 74, 604)
	for i := range g.buttons {
		g.buttons[i].Draw(s, i == g.mode)
	}
	if g.clear {
		overlay(s, "YOU TRIED EVERY MODE!\n\nTINT, FLASH, AND SHADOW\nARE ALL ColorScale.\nTAP / SPACE TO RESET")
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
	ebiten.SetWindowTitle("Tint & Drain — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
