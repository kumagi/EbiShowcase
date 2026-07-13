package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	W = 480
	H = 720
)

var questions = [][3]string{{"A storm blocks the road.", "Cross it now", "Help build shelter"}, {"A locked clock hides a clue.", "Open the panel", "Ask the keeper"}, {"The lantern floats overhead.", "Leap for it", "Study its rhythm"}}

type game struct {
	q, brave, kind int
	ended          bool
}

func (g *game) pick(n int) {
	if g.ended {
		*g = game{}
		return
	}
	if n == 0 {
		g.brave++
	} else {
		g.kind++
	}
	g.q++
	if g.q == 3 {
		g.ended = true
		g.q = 2
	}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.pick(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.pick(1)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_, y := ebiten.CursorPosition()
		if y >= 500 {
			g.pick(min(1, (y-500)/70))
		}
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		_, y := ebiten.TouchPosition(ids[0])
		if y >= 500 {
			g.pick(min(1, (y-500)/70))
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{30, 28, 65, 255})
	ebitenutil.DebugPrintAt(s, "CHOICE FLAG LAB", 178, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BRAVE %d   KIND %d", g.brave, g.kind), 170, 83)
	if !g.ended {
		v := questions[g.q]
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("QUESTION %d/3", g.q+1), 190, 180)
		ebitenutil.DebugPrintAt(s, v[0], 100, 235)
		for i := 0; i < 2; i++ {
			y := 500 + i*70
			vector.DrawFilledRect(s, 40, float32(y), 400, 56, color.RGBA{55, 75 + uint8(i*10), 118, 255}, false)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("[%d] %s", i+1, v[i+1]), 65, y+22)
		}
	} else {
		ending := "BALANCED GUIDE"
		if g.brave >= 2 {
			ending = "HARBOR HERO"
		} else if g.kind >= 2 {
			ending = "TOWN LIGHT"
		}
		vector.DrawFilledRect(s, 40, 225, 400, 230, color.RGBA{5, 11, 27, 245}, false)
		ebitenutil.DebugPrintAt(s, "FLAGS CREATE AN ENDING", 150, 270)
		ebitenutil.DebugPrintAt(s, ending, 180, 325)
		ebitenutil.DebugPrintAt(s, "TAP 1 OR 2 TO TRY AGAIN", 140, 405)
	}
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
