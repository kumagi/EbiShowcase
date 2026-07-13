package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const W, H = 480, 720

type game struct{ next, lap int }

var gates = [][2]float32{{240, 120}, {390, 250}, {390, 500}, {240, 610}, {90, 500}, {90, 250}}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.next = (g.next + 1) % len(gates)
		if g.next == 0 {
			g.lap++
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{37, 100, 65, 255})
	for i, p := range gates {
		c := color.RGBA{80, 90, 105, 255}
		if i == g.next {
			c = color.RGBA{255, 205, 70, 255}
		}
		vector.StrokeCircle(s, p[0], p[1], 30, 6, c, true)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d", i+1), int(p[0])-4, int(p[1])-5)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LAP %d  NEXT GATE %d", g.lap, g.next+1), 160, 50)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: PASS THE GLOWING GATE", 90, 680)
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Lap Gates")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
