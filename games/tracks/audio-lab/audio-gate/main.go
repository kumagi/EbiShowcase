// Audio Lab 01 — start original sound only after an explicit input gesture.
package main

import (
	"encoding/binary"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"image/color"
	"math"
)

type game struct {
	gate audiolab.Gate
	ctx  *audio.Context
	p    *audio.Player
	t    float64
}

func beep() []byte {
	b := make([]byte, 48000/5*4)
	for i := 0; i < len(b)/4; i++ {
		v := float32(math.Sin(float64(i)*2*math.Pi*440/48000) * .18)
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(v))
	}
	return b
}
func newGame() *game { return &game{ctx: audiolab.Context()} }
func (g *game) Update() error {
	g.t += .04
	hit := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if g.gate.Arm(hit) && g.p == nil {
		g.p = g.ctx.NewPlayerF32FromBytes(beep())
		g.p.Play()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "AUDIO LAB 01 · USER GESTURE GATE", 60, 32)
	on := g.gate.Armed()
	c := color.RGBA{80, 105, 145, 255}
	label := "TAP / SPACE TO ENABLE SOUND"
	if on {
		c = color.RGBA{50, 210, 165, 255}
		label = "SOUND ARMED · ORIGINAL BEEP PLAYED"
	}
	vector.DrawFilledCircle(s, 240, 330, float32(75+math.Sin(g.t)*8), c, true)
	ebitenutil.DebugPrintAt(s, label, 95, 530)
	ebitenutil.DebugPrintAt(s, "The game shows a silent visual before input; only a user gesture creates a player.", 23, 590)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
