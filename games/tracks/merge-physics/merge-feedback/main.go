package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const width, height = 480, 720

type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	tier, combo, tick, shake int
	sparks                   []spark
}

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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.tier = (g.tier + 1) % 7
		g.combo++
		g.shake = 5
		for i := 0; i < 16; i++ {
			a := float64(i) * math.Pi / 8
			g.sparks = append(g.sparks, spark{240, 330, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 28})
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{24, 33, 52, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 5
	}
	pulse := 100 + math.Sin(float64(g.tick)*.2)*6
	trackatlas.DrawCentered(s, trackatlas.Merge(g.tier+1), 240+ox, 330, pulse)
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TIER %d   CHAIN x%d", g.tier+1, g.combo), 170, 55)
	ebitenutil.DebugPrintAt(s, "MERGE -> PULSE -> PARTICLES -> SHAKE", 90, 470)
	vector.DrawFilledRect(s, 70, 560, 340, 80, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: MERGE", 155, 595)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Merge Feedback — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
