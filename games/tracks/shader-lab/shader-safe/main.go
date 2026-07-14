// Shader Lab 01 — an optional Kage pulse with an always-working DrawImage path.
package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
)

type game struct {
	t     float64
	orb   *ebiten.Image
	pulse *shaderlab.Pulse
}

func newGame() *game {
	o := ebiten.NewImage(240, 240)
	o.Fill(color.RGBA{40, 205, 225, 255})
	vector.DrawFilledCircle(o, 120, 120, 78, color.RGBA{255, 238, 120, 255}, true)
	return &game{orb: o, pulse: shaderlab.NewPulse()}
}
func (g *game) Update() error { g.t += .045; return nil }
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "SHADER LAB 01 · SAFE PULSE", 72, 28)
	ebitenutil.DebugPrintAt(s, "Kage is optional: the same orb still draws on a fallback path.", 32, 54)
	stage := ebiten.NewImage(240, 240)
	if !g.pulse.Draw(stage, g.orb, float32(g.t)) {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.Scale(float32(.84+.16*math.Sin(g.t)), 1, 1, 1)
		stage.DrawImage(g.orb, op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(120, 180)
	s.DrawImage(stage, op)
	status := "KAGE ACTIVE"
	if !g.pulse.Available() {
		status = "FALLBACK ACTIVE"
	}
	ebitenutil.DebugPrintAt(s, status, 172, 450)
	ebitenutil.DebugPrintAt(s, "Next: pass Time as a uniform and make a palette pulse.", 45, 640)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	ebiten.SetWindowTitle("Shader Lab 01")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
