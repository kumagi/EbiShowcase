package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kumagi/EbiShowcase/internal/uilab"
	"image/color"
	"math"
)

type game struct{ v float64 }

func (g *game) Update() error {
	g.v += .015
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.v = 0
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "UI LAB 03 · REUSABLE WIDGETS", 78, 30)
	uilab.Panel(s, 48, 130, 384, 310, color.RGBA{27, 54, 90, 255}, color.RGBA{80, 220, 190, 255})
	ebitenutil.DebugPrintAt(s, "★  RELIC RUN", 78, 175)
	ebitenutil.DebugPrintAt(s, "HP", 78, 235)
	uilab.Gauge(s, 125, 215, 255, 30, float32(.5+.45*math.Sin(g.v)), color.RGBA{245, 90, 100, 255})
	ebitenutil.DebugPrintAt(s, "MP", 78, 305)
	uilab.Gauge(s, 125, 285, 255, 30, float32(.5+.45*math.Cos(g.v)), color.RGBA{70, 165, 245, 255})
	ebitenutil.DebugPrintAt(s, "TAP / SPACE resets gauges", 130, 500)
	ebitenutil.DebugPrintAt(s, "Panel + corners + gauge + icon label are one reusable UI vocabulary.", 35, 590)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
