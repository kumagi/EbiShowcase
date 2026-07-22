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
	tick  int
}

// Run starts a renderer gallery for the exact game value passed by the
// capstone. Update and Layout remain the original methods; only Draw is wrapped.
func Run(inner ebiten.Game) error {
	wire, _ := ebiten.NewShader(wireframeSource)
	return ebiten.RunGame(&gallery{inner: inner, wire: wire})
}

func (g *gallery) Update() error {
	g.tick++
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

	// Keep the original snapshot at its original logical coordinates. Ebitengine
	// input is read by the unchanged inner Update, so scaling the playable view
	// would make pointer coordinates disagree with what the player sees.
	screen.DrawImage(snapshot, nil)
	drawPanelLabel(screen, 0, 0, min(w, 208), "PLAYABLE ORIGINAL DRAW")

	// The alternate renderers are presentation-only insets. They never become
	// the input surface, so the full-size game remains playable without an input
	// adapter or a second copy of game state.
	panelW := max(120, min(168, w/3))
	panelH := max(126, min(180, h/4))
	panelX := w - panelW - 8
	drawImageCell(screen, wire, panelX, 42, panelW, panelH, false)
	drawPanelLabel(screen, panelX, 42, panelW, "EDGE MAP")
	drawASCIIInset(screen, panelX, 54+panelH, panelW, panelH, g.tick)
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
	vector.DrawFilledRect(dst, float32(x+8), float32(y+8), float32(min(width-16, 184)), 24, color.RGBA{3, 10, 27, 225}, false)
	ebitenutil.DebugPrintAt(dst, label, x+16, y+15)
}

func drawASCIIInset(dst *ebiten.Image, x, y, width, height, tick int) {
	// This deliberately recognizable ASCII mini-game is a second renderer
	// concept, not a pixelated copy of the artwork. Alternating frames make the
	// character motion legible while keeping all rendering on the GPU path.
	frames := [...]string{
		"+----------------+\n| SCORE 000120    |\n|       *         |\n|   @==>     <c>  |\n|  /|\\      /\\  |\n|___|______==_____|\n| HP ***          |\n+----------------+",
		"+----------------+\n| SCORE 000120    |\n|        *        |\n|    @==>  <c>    |\n|   /|\\    /\\    |\n|___/______==_____|\n| HP ***          |\n+----------------+",
	}
	ascii := frames[(tick/24)%len(frames)]
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(width), float32(height), color.RGBA{3, 9, 18, 242}, false)
	vector.StrokeRect(dst, float32(x)+.5, float32(y)+.5, float32(width)-1, float32(height)-1, 1, color.RGBA{255, 209, 105, 210}, false)
	ebitenutil.DebugPrintAt(dst, "ASCII GAME MOCK", x+10, y+10)
	ebitenutil.DebugPrintAt(dst, ascii, x+10, y+34)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
