package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"image/color"
)

type game struct {
	gate  audiolab.Gate
	ctx   *audio.Context
	wave  int
	flash int
}

func newGame() *game { return &game{ctx: audio.NewContext(audiolab.SampleRate)} }
func (g *game) Update() error {
	hit := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if hit {
		g.gate.Arm(true)
		g.ctx.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Wave(g.wave), []float64{440, 660, 180}[g.wave], .18)).Play()
		g.wave = (g.wave + 1) % 3
		g.flash = 18
	}
	if g.flash > 0 {
		g.flash--
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	names := []string{"SINE · round", "SQUARE · bright", "NOISE · hit"}
	ebitenutil.DebugPrintAt(s, "AUDIO LAB 02 · PURE-GO ONE-SHOT", 60, 32)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: "+names[g.wave], 145, 58)
	for x := 20; x < 460; x += 8 {
		vector.StrokeLine(s, float32(x), 330, float32(x+8), float32(330+(x%31)-15), 2, color.RGBA{55, 215, 190, 255}, true)
	}
	if g.flash > 0 {
		vector.DrawFilledCircle(s, 240, 330, float32(g.flash*7), color.RGBA{255, 210, 90, 90}, true)
	}
	ebitenutil.DebugPrintAt(s, "waveform + frequency + decay → float32 PCM → Player", 72, 580)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
