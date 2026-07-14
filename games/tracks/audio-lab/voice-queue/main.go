package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"image/color"
)

type game struct {
	q      audiolab.Queue
	played int
}

func newGame() *game { return &game{q: audiolab.NewQueue(4)} }
func (g *game) Update() error {
	hit := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if hit {
		for i := 0; i < 3; i++ {
			g.q.Push("HIT")
		}
	}
	if _, ok := g.q.Pop(); ok {
		g.played++
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "AUDIO LAB 04 · FIXED VOICE QUEUE", 66, 32)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: enqueue 3 hits", 125, 58)
	for i := 0; i < 4; i++ {
		c := color.RGBA{43, 78, 120, 255}
		if i < g.q.Len() {
			c = color.RGBA{255, 190, 75, 255}
		}
		vector.DrawFilledRect(s, float32(60+i*92), 280, 74, 100, c, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("QUEUE %d / 4   PLAYED %d", g.q.Len(), g.played), 150, 440)
	ebitenutil.DebugPrintAt(s, "Fixed slots: repeated hits replace old requests instead of growing allocations.", 35, 590)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
