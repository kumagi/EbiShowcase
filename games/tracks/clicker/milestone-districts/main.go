package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
)

const width, height = 480, 720

type game struct {
	total         float64
	workers, tick int
}

func (g *game) Update() error {
	g.tick++
	g.total += float64(g.workers) / 30
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.total += 5
		if int(g.total)/50 > g.workers && g.workers < 6 {
			g.workers++
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	district := 0
	if g.total >= 50 {
		district = 1
	}
	if g.total >= 150 {
		district = 2
	}
	if g.total >= 400 {
		district = 3
	}
	bg := []color.RGBA{{35, 28, 50, 255}, {39, 48, 78, 255}, {41, 74, 68, 255}, {91, 62, 39, 255}}
	s.Fill(bg[district])
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TOTAL %.0f   DISTRICTS %d/3", g.total, district), 145, 35)
	for i := 0; i < 4; i++ {
		x := 35 + float32(i)*110
		c := color.RGBA{50, 55, 75, 255}
		if i <= district {
			c = color.RGBA{45, 205, 181, 255}
		}
		vector.DrawFilledRect(s, x, 130, 80, 320, c, false)
		if i <= district {
			trackatlas.DrawCentered(s, "bakery", float64(x)+40, 235, 70)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("D%d", i), int(x)+32, 330)
		}
	}
	next := []float64{50, 50, 150, 400}[district]
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("NEXT MILESTONE %.0f", next), 160, 500)
	vector.DrawFilledRect(s, 60, 550, 360, 90, color.RGBA{255, 192, 66, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: SHIP +5", 145, 585)
	ebitenutil.DebugPrintAt(s, "MILESTONES CHANGE THE WHOLE SCENE", 100, 680)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Milestone Districts — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
