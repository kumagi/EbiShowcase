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

type game struct {
	step        int
	dash, wings bool
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.step++
		if g.step == 2 {
			g.dash = true
		}
		if g.step == 4 {
			g.wings = true
		}
		if g.step > 6 {
			g.step = 0
			g.dash = false
			g.wings = false
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{20, 27, 48, 255})
	nodes := []string{"START", "MAP", "DASH", "OLD GATE", "WINGS", "HIGH LEDGE", "RELIC"}
	for i, n := range nodes {
		x := 40 + float32(i%4)*110
		y := 180 + float32(i/4)*180
		c := color.RGBA{55, 63, 84, 255}
		if i <= g.step {
			c = color.RGBA{45, 205, 181, 255}
		}
		vector.DrawFilledRect(s, x, y, 90, 100, c, false)
		ebitenutil.DebugPrintAt(s, n, int(x)+10, int(y)+45)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DASH %v   WINGS %v", g.dash, g.wings), 150, 60)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: ADVANCE EXPLORATION", 110, 670)
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ability Route")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
