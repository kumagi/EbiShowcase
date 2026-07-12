// vfx-additive — Visual Effects Lab STEP 05.
// Additive blending (op.Blend = BlendLighter): overlapping light adds to white.
package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const width, height = 480, 720

const glowR = 150

type orb struct {
	x, y   float64
	col    color.RGBA
	phase  float64
	anchor float64
}

type game struct {
	glow     *ebiten.Image
	orbs     []orb
	additive bool
	button   vfxui.Button
	t        float64
	clear    bool
}

func makeGlow() *ebiten.Image {
	size := glowR * 2
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x-glowR) + 0.5
			dy := float64(y-glowR) + 0.5
			d := math.Hypot(dx, dy) / glowR
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
	g := &game{glow: makeGlow(), additive: true}
	g.orbs = []orb{
		{x: 170, y: 300, col: color.RGBA{255, 70, 90, 255}, anchor: 170},
		{x: 310, y: 300, col: color.RGBA{70, 255, 120, 255}, phase: 2, anchor: 310},
		{x: 240, y: 300, col: color.RGBA{90, 140, 255, 255}, phase: 4, anchor: 240},
	}
	g.button = vfxui.Button{X: 130, Y: 636, W: 220, H: 54, Label: "BLEND: ADD"}
	return g
}

func (g *game) Update() error {
	g.t += 0.03
	if g.clear {
		if vfxui.AnyPressStart() {
			*g = *newGame()
		}
		return nil
	}
	if g.button.Tapped() || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.additive = !g.additive
		if g.additive {
			g.button.Label = "BLEND: ADD"
		} else {
			g.button.Label = "BLEND: NORMAL"
		}
	}
	// First orb follows the pointer; others drift so overlaps happen naturally.
	if x, y, ok := vfxui.Held(); ok && y < 620 {
		g.orbs[0].x, g.orbs[0].y = x, y
	}
	for i := 1; i < len(g.orbs); i++ {
		g.orbs[i].x = g.orbs[i].anchor + math.Sin(g.t+g.orbs[i].phase)*70
		g.orbs[i].y = 320 + math.Cos(g.t*1.3+g.orbs[i].phase)*70
	}
	if g.additive && g.close(0, 1) && g.close(1, 2) && g.close(0, 2) {
		g.clear = true
	}
	return nil
}

func (g *game) close(i, j int) bool {
	return math.Hypot(g.orbs[i].x-g.orbs[j].x, g.orbs[i].y-g.orbs[j].y) < 90
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 8, 16, 255})
	for _, o := range g.orbs {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(o.x-glowR, o.y-glowR)
		op.ColorScale.ScaleWithColor(o.col)
		if g.additive {
			op.Blend = ebiten.BlendLighter
		}
		s.DrawImage(g.glow, op)
	}

	mode := "ADDITIVE (light adds up)"
	if !g.additive {
		mode = "NORMAL (front hides back)"
	}
	ebitenutil.DebugPrintAt(s, "THREE COLORED LIGHTS", 148, 24)
	ebitenutil.DebugPrintAt(s, mode, 16, 52)
	ebitenutil.DebugPrintAt(s, "DRAG THE RED LIGHT OVER THE OTHERS (SPACE toggles)", 58, 604)
	g.button.Draw(s, g.additive)
	if g.clear {
		overlay(s, "WHITE-HOT OVERLAP!\n\nBlendLighter ADDS COLOR,\nSO LIGHT PILES UP TO WHITE.\nTAP / SPACE TO RESET")
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
	ebiten.SetWindowTitle("Additive Glow — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
