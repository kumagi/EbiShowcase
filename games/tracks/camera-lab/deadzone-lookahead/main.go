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
	cam         cameralab.State
	target, vel cameralab.Vec
}

func newGame() *game {
	return &game{cam: cameralab.State{Pos: cameralab.Vec{240, 360}, ViewW: 480, ViewH: 720}, target: cameralab.Vec{240, 360}}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.vel = cameralab.Vec{X: 6}
		g.target.X += 80
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.vel = cameralab.Vec{X: -6}
		g.target.X -= 80
	}
	g.cam.FollowDeadZone(g.target, g.vel, 70, 80, 7)
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "CAMERA LAB 03 · DEAD ZONE + LOOK AHEAD", 45, 30)
	vector.StrokeRect(s, 170, 260, 140, 160, 3, color.RGBA{80, 220, 190, 255}, true)
	p := g.cam.WorldToScreen(g.target)
	vector.DrawFilledCircle(s, float32(p.X), float32(p.Y), 20, color.RGBA{255, 210, 85, 255}, true)
	ebitenutil.DebugPrintAt(s, "box = dead zone · velocity moves camera aim ahead", 70, 550)
	ebitenutil.DebugPrintAt(s, "Left/Right or tap: switch movement direction", 95, 610)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
