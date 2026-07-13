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

type point struct{ x, y float64 }
type game struct{ stage, hits int }

var pegs = [][]point{{{240, 240}, {150, 390}, {350, 390}}, {{160, 250}, {320, 250}, {160, 430}, {320, 430}}, {{240, 170}, {90, 290}, {390, 290}, {240, 420}}}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.hits++
		if g.hits >= []int{3, 4, 5}[g.stage] {
			g.stage = (g.stage + 1) % 3
			g.hits = 0
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	bg := []color.RGBA{{15, 34, 52, 255}, {42, 29, 61, 255}, {65, 28, 37, 255}}
	s.Fill(bg[g.stage])
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE %d/3  TARGET HITS %d/%d", g.stage+1, g.hits, []int{3, 4, 5}[g.stage]), 135, 40)
	for _, p := range pegs[g.stage] {
		trackatlas.DrawCentered(s, "peg", p.x, p.y, 55)
	}
	vector.DrawFilledRect(s, 70, 560, 340, 80, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: SCORE A HIT", 135, 595)
	ebitenutil.DebugPrintAt(s, "SAME PHYSICS, DIFFERENT STAGE DATA", 95, 675)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Slingshot Stage Gimmicks — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
