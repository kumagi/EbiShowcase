// Shader Lab 04 — status feedback is display data, not combat logic.
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
	status int
	damage float64
	src    *ebiten.Image
	fx     *shaderlab.Status
}

func newGame() *game {
	i := ebiten.NewImage(240, 240)
	i.Fill(color.RGBA{35, 58, 110, 255})
	vector.DrawFilledCircle(i, 120, 122, 78, color.RGBA{255, 220, 100, 255}, true)
	vector.DrawFilledCircle(i, 92, 104, 9, color.RGBA{25, 32, 65, 255}, true)
	vector.DrawFilledCircle(i, 148, 104, 9, color.RGBA{25, 32, 65, 255}, true)
	return &game{src: i, fx: shaderlab.NewStatus()}
}
func (g *game) Update() error {
	g.damage = math.Max(0, g.damage-.04)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.status = (g.status + 1) % 3
		g.damage = 1
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	names := []string{"POISON", "NORMAL", "FREEZE"}
	ebitenutil.DebugPrintAt(s, "SHADER LAB 04 · STATUS COLORS", 73, 28)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: "+names[g.status], 145, 55)
	d := ebiten.NewImage(240, 240)
	if !g.fx.Draw(d, g.src, float32(g.status), float32(g.damage)) {
		op := &ebiten.DrawImageOptions{}
		if g.status == 0 {
			op.ColorScale.Scale(1, .35, .3, 1)
		}
		if g.status == 2 {
			op.ColorScale.Scale(.4, .9, 1, 1)
		}
		d.DrawImage(g.src, op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(120, 180)
	s.DrawImage(d, op)
	ebitenutil.DebugPrintAt(s, "Damage separates red/blue samples; status tints after rules resolve.", 37, 600)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
