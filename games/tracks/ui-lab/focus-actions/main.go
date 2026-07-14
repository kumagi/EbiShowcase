package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/uilab"
	"image/color"
)

type game struct {
	f       uilab.Focus
	message string
}

func newGame() *game {
	return &game{f: uilab.Focus{Count: 3, Disabled: map[int]bool{1: true}}, message: "SELECT START"}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.f.Move(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.f.Move(-1)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.f.Move(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.f.Activate() {
			g.message = "ACTION: " + []string{"START", "LOCKED", "QUIT"}[g.f.Index]
		} else {
			g.message = "DISABLED: LOCKED"
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "UI LAB 04 · FOCUS + ACTION MAP", 70, 30)
	labels := []string{"START", "LOCKED", "QUIT"}
	for i, l := range labels {
		c := color.RGBA{48, 85, 130, 255}
		if g.f.Disabled[i] {
			c = color.RGBA{55, 58, 70, 255}
		}
		if g.f.Index == i {
			c = color.RGBA{55, 215, 180, 255}
		}
		vector.DrawFilledRect(s, 100, float32(160+i*100), 280, 72, c, true)
		ebitenutil.DebugPrintAt(s, l, 210, 185+i*100)
	}
	ebitenutil.DebugPrintAt(s, g.message, 145, 510)
	ebitenutil.DebugPrintAt(s, "Arrows/tap: focus · Enter/Space: action · disabled stays visibly unavailable.", 25, 590)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
