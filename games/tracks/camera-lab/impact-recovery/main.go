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
	i     cameralab.Impact
	frame int
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.i.Trigger(12, 5, 9)
	}
	g.i.Tick()
	g.frame++
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	o := g.i.Offset(g.frame)
	ebitenutil.DebugPrintAt(s, "CAMERA LAB 04 · HIT STOP + SHAKE", 70, 30)
	state := "RECOVERY"
	if g.i.Frozen() {
		state = "HIT STOP"
	}
	vector.DrawFilledCircle(s, 240+float32(o.X), 330+float32(o.Y), 48, color.RGBA{255, 105, 80, 255}, true)
	ebitenutil.DebugPrintAt(s, state, 195, 460)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: stop briefly, then deterministic shake fades out.", 55, 590)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
