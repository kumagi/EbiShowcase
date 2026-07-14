package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"image/color"
)

type game struct {
	cam    cameralab.State
	target cameralab.Vec
}

func newGame() *game {
	return &game{cam: cameralab.State{Pos: cameralab.Vec{240, 360}, ViewW: 480, ViewH: 720}, target: cameralab.Vec{240, 360}}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.target.X += 180
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.target.X -= 180
	}
	g.cam.Follow(g.target, .08, 1200, 900)
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "CAMERA LAB 02 · FOLLOW + CLAMP", 72, 30)
	for x := 0; x < 1200; x += 80 {
		p := g.cam.WorldToScreen(cameralab.Vec{float64(x), 400})
		vector.StrokeLine(s, float32(p.X), 120, float32(p.X), 620, 1, color.RGBA{80, 100, 130, 255}, true)
	}
	p := g.cam.WorldToScreen(g.target)
	vector.DrawFilledCircle(s, float32(p.X), float32(p.Y), 22, color.RGBA{255, 210, 85, 255}, true)
	ebitenutil.DebugPrintAt(s, "Left/Right or tap: move target · camera follows but stops at world edge.", 28, 650)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
