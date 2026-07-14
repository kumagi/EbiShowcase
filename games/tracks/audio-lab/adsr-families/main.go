package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"image/color"
)

type game struct{ i int }

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.i = (g.i + 1) % 3
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	names := []string{"JUMP", "MAGIC", "HIT"}
	_, _, a := audiolab.Family([]string{"jump", "magic", "hit"}[g.i])
	ebitenutil.DebugPrintAt(s, "AUDIO LAB 03 · ADSR FAMILIES", 76, 32)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: "+names[g.i], 150, 58)
	pts := []float64{a.Attack, a.Attack + a.Decay, .65, .65 + a.Release}
	x := 70.
	for n, p := range pts {
		y := 440.
		if n == 0 {
			y = 170
		}
		if n == 1 {
			y = 260
		}
		if n == 2 {
			y = 330
		}
		vector.StrokeLine(s, float32(x), 440, float32(x+p*350), float32(y), 5, color.RGBA{80, 220, 190, 255}, true)
		x += p * 350
	}
	ebitenutil.DebugPrintAt(s, "Attack → Decay → Sustain → Release", 110, 540)
	ebitenutil.DebugPrintAt(s, "Same oscillator, different envelope parameters = different game meaning.", 38, 600)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
