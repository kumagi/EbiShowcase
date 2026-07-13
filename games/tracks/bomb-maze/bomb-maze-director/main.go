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

type pt struct{ x, y int }
type game struct{ stage, broken int }

var layouts = [][]pt{{{2, 1}, {4, 1}, {2, 3}, {4, 3}}, {{1, 2}, {3, 2}, {5, 2}, {7, 2}, {3, 5}}, {{2, 1}, {4, 1}, {6, 1}, {2, 3}, {6, 3}, {2, 5}, {4, 5}, {6, 5}}}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.broken++
		if g.broken >= len(layouts[g.stage]) {
			g.stage = (g.stage + 1) % 3
			g.broken = 0
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	bg := []color.RGBA{{15, 28, 44, 255}, {39, 25, 55, 255}, {60, 25, 32, 255}}
	s.Fill(bg[g.stage])
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MAZE %d/3  WALLS %d/%d", g.stage+1, g.broken, len(layouts[g.stage])), 155, 40)
	for i, p := range layouts[g.stage] {
		x := 48 + float64(p.x)*48
		y := 110 + float64(p.y)*58
		if i < g.broken {
			trackatlas.DrawCentered(s, "flame", x, y, 40)
		} else {
			trackatlas.DrawCentered(s, "tile-crate", x, y, 46)
		}
	}
	vector.DrawFilledRect(s, 70, 580, 340, 75, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: BREAK WALL", 145, 610)
	ebitenutil.DebugPrintAt(s, "LAYOUT + GOAL + ENEMY SPEED = STAGE DATA", 65, 685)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Bomb Maze Director — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
