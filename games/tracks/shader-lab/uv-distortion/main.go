// Shader Lab 03 — three coordinate offsets sample one scene in different ways.
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"image/color"
	"math"
)

type game struct {
	t, impact float64
	mode      int
	src       *ebiten.Image
	fx        *shaderlab.Distort
}

func newGame() *game {
	i := ebiten.NewImage(240, 240)
	i.Fill(color.RGBA{20, 65, 105, 255})
	for y := 20; y < 240; y += 34 {
		vector.StrokeLine(i, 0, float32(y), 240, float32(y), 2, color.RGBA{105, 225, 240, 255}, true)
	}
	vector.DrawFilledCircle(i, 120, 115, 42, color.RGBA{255, 180, 74, 255}, true)
	return &game{src: i, fx: shaderlab.NewDistort()}
}
func (g *game) Update() error {
	g.t += .035
	g.impact += .018
	if g.impact > 1 {
		g.impact = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.mode = (g.mode + 1) % 3
		g.impact = 0
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	names := []string{"WATER", "HEAT HAZE", "IMPACT WAVE"}
	ebitenutil.DebugPrintAt(s, "SHADER LAB 03 · UV DISTORTION", 72, 28)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: "+names[g.mode], 135, 55)
	d := ebiten.NewImage(240, 240)
	if !g.fx.Draw(d, g.src, float32(g.t), float32(g.mode), float32(g.impact)) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(math.Sin(g.t*2)*float64(g.mode+1), 0)
		d.DrawImage(g.src, op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(120, 180)
	s.DrawImage(d, op)
	ebitenutil.DebugPrintAt(s, "Change uv before imageSrc0At(uv): one idea, three effects.", 43, 600)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
