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

type game struct{ stage, merges int }

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.merges++
		if g.merges >= []int{3, 5, 7}[g.stage] {
			g.stage = (g.stage + 1) % 3
			g.merges = 0
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	bg := []color.RGBA{{20, 36, 54, 255}, {48, 31, 61, 255}, {67, 32, 32, 255}}
	s.Fill(bg[g.stage])
	target := []int{3, 5, 7}[g.stage]
	gravity := []string{"LIGHT GRAVITY", "NORMAL GRAVITY", "HEAVY GRAVITY"}[g.stage]
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("CHALLENGE %d/3 — %s", g.stage+1, gravity), 120, 45)
	for i := 0; i < target; i++ {
		x := 70 + float64(i%4)*110
		y := 220 + float64(i/4)*130
		alpha := uint8(90)
		if i < g.merges {
			alpha = 255
		}
		trackatlas.DrawTinted(s, trackatlas.Merge(min(6, g.stage+3)), x, y, 70, 1, 1, 1, float32(alpha)/255)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MERGES %d/%d", g.merges, target), 185, 500)
	vector.DrawFilledRect(s, 70, 560, 340, 80, color.RGBA{255, 190, 75, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: COMPLETE MERGE", 120, 595)
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Merge Challenge Stages — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
