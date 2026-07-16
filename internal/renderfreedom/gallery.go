// Copyright 2026 Ebi Showcase contributors
// SPDX-License-Identifier: Apache-2.0

// Package renderfreedom wraps an existing Ebitengine game without replacing
// its Update or Layout methods. Its Draw projects one snapshot three ways so a
// learner can see that the renderer is not the game state.
package renderfreedom

import (
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed wireframe.kage
var wireframeSource []byte

type gallery struct {
	inner ebiten.Game
	wire  *ebiten.Shader
}

// Run starts a renderer gallery for the exact game value passed by the
// capstone. Update and Layout remain the original methods; only Draw is wrapped.
func Run(inner ebiten.Game) error {
	wire, _ := ebiten.NewShader(wireframeSource)
	return ebiten.RunGame(&gallery{inner: inner, wire: wire})
}

func (g *gallery) Update() error {
	return g.inner.Update()
}

func (g *gallery) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.inner.Layout(outsideWidth, outsideHeight)
}

func (g *gallery) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if w < 2 || h < 2 {
		return
	}

	// The original Draw runs once. Every panel below observes these exact same
	// pixels from the exact same game snapshot.
	snapshot := ebiten.NewImage(w, h)
	g.inner.Draw(snapshot)

	wire := ebiten.NewImage(w, h)
	if g.wire != nil {
		op := &ebiten.DrawRectShaderOptions{}
		op.Images[0] = snapshot
		wire.DrawRectShader(w, h, g.wire, op)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.Scale(.2, .9, 1, 1)
		wire.DrawImage(snapshot, op)
	}

	// Scaling down once and back with nearest filtering turns the same snapshot
	// into visible rectangular cells. No rule or state is reconstructed here.
	tinyW, tinyH := max(12, w/18), max(12, h/18)
	blocks := ebiten.NewImage(tinyW, tinyH)
	down := &ebiten.DrawImageOptions{}
	down.GeoM.Scale(float64(tinyW)/float64(w), float64(tinyH)/float64(h))
	down.Filter = ebiten.FilterLinear
	blocks.DrawImage(snapshot, down)

	screen.Fill(color.RGBA{4, 9, 22, 255})
	cellW, cellH := w/2, h/2
	drawImageCell(screen, snapshot, 0, 0, cellW, cellH, false)
	drawImageCell(screen, wire, cellW, 0, w-cellW, cellH, false)
	drawImageCell(screen, blocks, 0, cellH, cellW, h-cellH, true)
	drawExplanation(screen, cellW, cellH, w-cellW, h-cellH)
	drawPanelLabel(screen, 0, 0, cellW, "POLISHED DRAW")
	drawPanelLabel(screen, cellW, 0, w-cellW, "WIREFRAME DRAW")
	drawPanelLabel(screen, 0, cellH, cellW, "RECTANGLE DRAW")
}

func drawImageCell(dst, src *ebiten.Image, x, y, width, height int, nearest bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(width)/float64(src.Bounds().Dx()), float64(height)/float64(src.Bounds().Dy()))
	op.GeoM.Translate(float64(x), float64(y))
	if nearest {
		op.Filter = ebiten.FilterNearest
	} else {
		op.Filter = ebiten.FilterLinear
	}
	dst.DrawImage(src, op)
	vector.StrokeRect(dst, float32(x)+.5, float32(y)+.5, float32(width)-1, float32(height)-1, 1, color.RGBA{105, 221, 232, 180}, false)
}

func drawPanelLabel(dst *ebiten.Image, x, y, width int, label string) {
	vector.DrawFilledRect(dst, float32(x+8), float32(y+8), float32(min(width-16, 148)), 24, color.RGBA{3, 10, 27, 225}, false)
	ebitenutil.DebugPrintAt(dst, label, x+16, y+15)
}

func drawExplanation(dst *ebiten.Image, x, y, width, height int) {
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(width), float32(height), color.RGBA{8, 17, 39, 255}, false)
	vector.StrokeRect(dst, float32(x)+.5, float32(y)+.5, float32(width)-1, float32(height)-1, 1, color.RGBA{255, 209, 105, 190}, false)
	left := x + max(12, width/12)
	top := y + max(18, height/7)
	ebitenutil.DebugPrintAt(dst, "ONE GAME SNAPSHOT", left, top)
	ebitenutil.DebugPrintAt(dst, "", left, top+20)
	ebitenutil.DebugPrintAt(dst, "UPDATE: ORIGINAL", left, top+40)
	ebitenutil.DebugPrintAt(dst, "LAYOUT: ORIGINAL", left, top+58)
	ebitenutil.DebugPrintAt(dst, "DRAW: REPLACEABLE", left, top+76)
	vector.StrokeLine(dst, float32(left), float32(top+102), float32(min(x+width-14, left+180)), float32(top+102), 2, color.RGBA{94, 234, 212, 220}, false)
	ebitenutil.DebugPrintAt(dst, "SAME RULES / SAME STATE", left, top+116)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
