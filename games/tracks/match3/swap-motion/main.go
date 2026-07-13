package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const width, height = 480, 720

type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	phase, tick, score, flash, shake int
	valid, resolved                  bool
	sparks                           []spark
}

func (g *game) start() { g.phase, g.valid, g.resolved = 42, (g.score/100)%2 == 0, false }
func (g *game) Update() error {
	g.tick++
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.phase > 0 {
		g.phase--
		if g.phase == 18 && !g.resolved {
			g.resolved = true
			if g.valid {
				g.score += 100
				g.flash = 10
				g.shake = 7
				for i := 0; i < 16; i++ {
					a := float64(i) * math.Pi / 8
					g.sparks = append(g.sparks, spark{240, 340, math.Cos(a) * 3, math.Sin(a) * 3, 28})
				}
			}
		}
	}
	if g.flash > 0 {
		g.flash--
	}
	if g.phase == 0 && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0) {
		g.start()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{15, 25, 45, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 5
	}
	t := 0.0
	if g.phase > 18 {
		t = float64(42-g.phase) / 24
	} else if g.phase > 0 && !g.valid {
		t = float64(g.phase) / 18
	} else if g.phase > 0 {
		t = 1
	}
	x1, x2 := 150+96*t, 330-96*t
	vector.DrawFilledRect(s, 105, 270, 270, 140, color.RGBA{27, 48, 70, 255}, false)
	trackatlas.DrawCentered(s, "gem-red", x1+ox, 340, 72)
	trackatlas.DrawCentered(s, "gem-blue", x2+ox, 340, 72)
	if g.flash > 0 {
		vector.StrokeCircle(s, 240, 340, float32(48+(10-g.flash)*5), 5, color.RGBA{255, 220, 95, 255}, true)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), 4, color.RGBA{255, 200, 70, 255}, true)
	}
	state := "READY"
	if g.phase > 18 {
		state = "ANTICIPATION / TRAVEL"
	} else if g.phase > 0 && g.valid {
		state = "CONTACT / CLEAR"
	} else if g.phase > 0 {
		state = "INVALID / RETURN"
	}
	ebitenutil.DebugPrintAt(s, "SWAP MOTION LAB", 170, 70)
	ebitenutil.DebugPrintAt(s, state, 150, 465)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %04d", g.score), 195, 510)
	ebitenutil.DebugPrintAt(s, "Every second swap returns", 140, 570)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: SWAP", 155, 665)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Swap Motion Lab — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
