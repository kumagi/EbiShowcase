// vfx-transform — Visual Effects Lab STEP 02.
// Rotate and scale with GeoM, and why the pivot (order of operations) matters.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const width, height = 480, 720

const (
	targetAngle = 0.7
	targetScale = 1.7
	drawH       = 150.0
)

type game struct {
	angle   float64
	scale   float64
	center  bool
	buttons []vfxui.Button
	clear   bool
	glow    float64
}

func newGame() *game {
	g := &game{angle: 0, scale: 1, center: true}
	w := 86.0
	gap := 8.0
	x := (width - (w*5 + gap*4)) / 2
	labels := []string{"ROT-", "ROT+", "SCL-", "SCL+", "PIVOT"}
	for _, l := range labels {
		g.buttons = append(g.buttons, vfxui.Button{X: x, Y: 636, W: w, H: 54, Label: l})
		x += w + gap
	}
	return g
}

func (g *game) Update() error {
	g.glow += 0.1
	if g.clear {
		if vfxui.AnyPressStart() {
			*g = *newGame()
		}
		return nil
	}
	if g.buttons[0].Tapped() || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.angle -= 0.05
	}
	if g.buttons[1].Tapped() || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.angle += 0.05
	}
	if g.buttons[2].Tapped() || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.scale = math.Max(0.4, g.scale-0.05)
	}
	if g.buttons[3].Tapped() || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.scale = math.Min(2.6, g.scale+0.05)
	}
	if g.buttons[4].Tapped() || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.center = !g.center
	}
	if math.Abs(g.angle-targetAngle) < 0.06 && math.Abs(g.scale-targetScale) < 0.06 {
		g.clear = true
	}
	return nil
}

func drawHero(dst *ebiten.Image, cx, cy, angle, scale float64, center bool) {
	sprite := hero.Image()
	b := sprite.Bounds()
	sw, sh := float64(b.Dx()), float64(b.Dy())
	base := drawH / sh
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	if center {
		op.GeoM.Translate(-sw/2, -sh/2)
		op.GeoM.Scale(base*scale, base*scale)
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(cx, cy)
	} else {
		op.GeoM.Scale(base*scale, base*scale)
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(cx, cy)
	}
	dst.DrawImage(sprite, op)
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{16, 24, 40, 255})
	cx, cy := float64(width/2), 330.0

	// Pivot marker.
	vector.StrokeLine(s, float32(cx-14), float32(cy), float32(cx+14), float32(cy), 2, color.RGBA{90, 108, 150, 255}, false)
	vector.StrokeLine(s, float32(cx), float32(cy-14), float32(cx), float32(cy+14), 2, color.RGBA{90, 108, 150, 255}, false)

	// Ghost target: faint hero at the goal angle & scale.
	{
		sprite := hero.Image()
		b := sprite.Bounds()
		sw, sh := float64(b.Dx()), float64(b.Dy())
		base := drawH / sh
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterLinear
		op.GeoM.Translate(-sw/2, -sh/2)
		op.GeoM.Scale(base*targetScale, base*targetScale)
		op.GeoM.Rotate(targetAngle)
		op.GeoM.Translate(cx, cy)
		op.ColorScale.Scale(0.35, 0.55, 0.65, 0.35)
		s.DrawImage(sprite, op)
	}

	drawHero(s, cx, cy, g.angle, g.scale, g.center)

	pivot := "CENTER"
	if !g.center {
		pivot = "CORNER"
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("angle %.2f rad   scale x%.2f", g.angle, g.scale), 16, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PIVOT %s   (match the faint ghost)", pivot), 16, 46)
	ebitenutil.DebugPrintAt(s, "ROT / SCL / PIVOT BUTTONS  OR  ARROW KEYS + P", 66, 604)
	for i := range g.buttons {
		g.buttons[i].Draw(s, i == 4 && g.center)
	}
	if g.clear {
		overlay(s, "MATCHED THE GHOST!\n\nROTATE + SCALE = GeoM.\nTAP / SPACE TO RESET")
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 250, 370, 150, color.RGBA{8, 16, 32, 240}, false)
	vector.StrokeRect(s, 55, 250, 370, 150, 3, color.RGBA{120, 240, 220, 255}, false)
	ebitenutil.DebugPrintAt(s, msg, 95, 285)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Rotate & Scale — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
