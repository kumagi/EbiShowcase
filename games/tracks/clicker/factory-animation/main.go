package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const width, height = 480, 720

type game struct {
	sweets         float64
	machines, tick int
	flash          float64
}

func (g *game) Update() error {
	g.tick++
	g.sweets += float64(g.machines) / 60
	if g.flash > 0 {
		g.flash--
	}
	tap := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if tap {
		_, y := ebiten.CursorPosition()
		if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
			_, y = ebiten.TouchPosition(ids[0])
		}
		if y > 500 && g.sweets >= 10 {
			g.sweets -= 10
			g.machines++
			g.flash = 12
		} else {
			g.sweets++
			g.flash = 6
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{35, 31, 55, 255})
	trackatlas.DrawCentered(s, "bakery", 240, 265, 180+math.Sin(float64(g.tick)*.16)*5)
	for i := 0; i < g.machines && i < 8; i++ {
		phase := math.Mod(float64(g.tick)*2+float64(i)*48, 400)
		trackatlas.DrawCentered(s, "coin", 40+phase, 430+math.Sin(phase*.06)*7, 18)
	}
	if g.flash > 0 {
		vector.StrokeCircle(s, 240, 265, float32(110-g.flash*2), 4, color.RGBA{255, 220, 90, 220}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SWEETS %.1f  MACHINES %d", g.sweets, g.machines), 150, 35)
	ebitenutil.DebugPrintAt(s, "EACH MACHINE ADDS A MOVING PRODUCT", 95, 65)
	vector.DrawFilledRect(s, 70, 520, 340, 90, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, "BUY MACHINE — COST 10", 145, 555)
	ebitenutil.DebugPrintAt(s, "TAP BAKERY / SPACE TO BAKE", 125, 680)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Animated Factory — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
