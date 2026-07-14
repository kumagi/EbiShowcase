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

type game struct{ m audiolab.Mixer }

func newGame() *game { return &game{m: audiolab.NewMixer()} }
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.m.TriggerImportantSE(90)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.m.Paused = !g.m.Paused
	}
	g.m.Tick()
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "AUDIO LAB 05 · MIXER + DUCKING", 70, 32)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: important SE · P: pause", 90, 58)
	w := float32(g.m.BGMGain() * 320)
	vector.DrawFilledRect(s, 80, 250, w, 54, color.RGBA{65, 135, 225, 255}, true)
	vector.DrawFilledRect(s, 80, 350, float32(g.m.SEGain()*320), 54, color.RGBA{55, 215, 170, 255}, true)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BGM gain %.2f", g.m.BGMGain()), 90, 280)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SE gain %.2f", g.m.SEGain()), 90, 380)
	state := "RUNNING"
	if g.m.Paused {
		state = "PAUSED"
	}
	ebitenutil.DebugPrintAt(s, state, 200, 480)
	ebitenutil.DebugPrintAt(s, "Looping BGM follows BGM gain; important SE temporarily ducks it.", 55, 590)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
