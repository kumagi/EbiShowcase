package uilab

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// Panel draws a nine-slice-style frame: stable corners/border with a stretchable centre.
func Panel(dst *ebiten.Image, x, y, w, h float32, fill, border color.Color) {
	vector.DrawFilledRect(dst, x, y, w, h, fill, true)
	vector.StrokeRect(dst, x, y, w, h, 3, border, true)
	for _, p := range [][2]float32{{x + 7, y + 7}, {x + w - 7, y + 7}, {x + 7, y + h - 7}, {x + w - 7, y + h - 7}} {
		vector.DrawFilledCircle(dst, p[0], p[1], 3, border, true)
	}
}
func Gauge(dst *ebiten.Image, x, y, w, h, value float32, fill color.Color) {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	vector.DrawFilledRect(dst, x, y, w, h, color.RGBA{15, 24, 43, 255}, true)
	vector.DrawFilledRect(dst, x+3, y+3, (w-6)*value, h-6, fill, true)
	vector.StrokeRect(dst, x, y, w, h, 2, color.White, true)
}
