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
	room int
	seen [8]bool
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.room = (g.room + 1) % 8
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.room = (g.room + 7) % 8
	}
	g.seen[g.room] = true
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{13, 24, 42, 255})
	for i := 0; i < 8; i++ {
		x := 35 + float32(i%4)*110
		y := 180 + float32(i/4)*150
		c := color.RGBA{42, 48, 65, 255}
		if g.seen[i] {
			c = color.RGBA{75, 158, 147, 255}
		}
		vector.DrawFilledRect(s, x, y, 90, 110, c, false)
		if i == g.room {
			vector.StrokeRect(s, x-4, y-4, 98, 118, 4, color.RGBA{255, 205, 70, 255}, true)
		}
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("R%d", i+1), int(x)+35, int(y)+45)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DISCOVERED %d/8", count(g.seen)), 175, 60)
	ebitenutil.DebugPrintAt(s, "LEFT/RIGHT OR TAP: EXPLORE", 135, 670)
}
func count(a [8]bool) int {
	n := 0
	for _, v := range a {
		if v {
			n++
		}
	}
	return n
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	g := &game{}
	g.seen[0] = true
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Exploration Map")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
