// Package vfxui holds tiny shared helpers for the Visual Effects Lab toys:
// unified pointer input (mouse + touch) and on-screen tap buttons so every
// lesson works the same way on desktop and on phones.
package vfxui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Held reports the current pointer position while a mouse button or a finger
// is held down. It does not care whether the press just started.
func Held() (x, y float64, ok bool) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		return float64(cx), float64(cy), true
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		tx, ty := ebiten.TouchPosition(ids[0])
		return float64(tx), float64(ty), true
	}
	return 0, 0, false
}

// JustPressed reports the pointer position on the first frame of a press.
func JustPressed() (x, y float64, ok bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		return float64(cx), float64(cy), true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		tx, ty := ebiten.TouchPosition(ids[0])
		return float64(tx), float64(ty), true
	}
	return 0, 0, false
}

// AnyPressStart reports whether any pointer or the space key started this frame.
func AnyPressStart() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

// Button is a simple rounded tap target with a centered label.
type Button struct {
	X, Y, W, H float64
	Label      string
	Fill       color.RGBA
}

// Contains reports whether the point (px, py) is inside the button.
func (b Button) Contains(px, py float64) bool {
	return px >= b.X && px <= b.X+b.W && py >= b.Y && py <= b.Y+b.H
}

// Tapped reports whether a press started inside the button this frame.
func (b Button) Tapped() bool {
	if x, y, ok := JustPressed(); ok {
		return b.Contains(x, y)
	}
	return false
}

// Draw paints the button. When active is true it uses a brighter border.
func (b Button) Draw(dst *ebiten.Image, active bool) {
	fill := b.Fill
	if fill.A == 0 {
		fill = color.RGBA{40, 54, 82, 235}
	}
	vector.DrawFilledRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), fill, false)
	edge := color.RGBA{96, 116, 158, 255}
	if active {
		edge = color.RGBA{120, 240, 220, 255}
	}
	vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 3, edge, false)
	tx := int(b.X + b.W/2 - float64(len(b.Label))*3)
	ty := int(b.Y + b.H/2 - 8)
	ebitenutil.DebugPrintAt(dst, b.Label, tx, ty)
}
