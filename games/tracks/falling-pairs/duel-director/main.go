package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
)

const width, height = 480, 720

type rule struct {
	name                      string
	goal, missLimit, incoming int
	bg                        color.RGBA
}

var rules = []rule{{"SUNNY LAGOON", 6, 3, 1, color.RGBA{12, 40, 60, 255}}, {"CURRENT CAVE", 9, 3, 2, color.RGBA{8, 55, 58, 255}}, {"STORM PALACE", 12, 2, 2, color.RGBA{35, 20, 65, 255}}}

type game struct {
	stage, sent, misses, best, tick, burst int
	clear                                  bool
}

func (g *game) reset(n int) { g.stage = n; g.sent = 0; g.misses = 0; g.clear = false }
func (g *game) Update() error {
	g.tick++
	if g.burst > 0 {
		g.burst--
	}
	pressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if !pressed {
		return nil
	}
	r := rules[g.stage]
	if g.clear {
		g.reset((g.stage + 1) % 3)
		return nil
	}
	g.sent += 2 + (g.tick % 3)
	g.burst = 15
	if g.sent >= r.goal {
		g.clear = true
		grade := g.sent*100 - g.misses*50
		if grade > g.best {
			g.best = grade
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	r := rules[g.stage]
	s.Fill(r.bg)
	vector.DrawFilledRect(s, 65, 130, 350, 390, color.RGBA{6, 17, 33, 180}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DUEL %d/3 / %s", g.stage+1, r.name), 145, 65)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("GOAL %d  MISS LIMIT %d  INCOMING %d", r.goal, r.missLimit, r.incoming), 105, 95)
	for i := 0; i < r.goal; i++ {
		x := 145 + (i%4)*65
		y := 210 + (i/4)*65
		if i < g.sent {
			trackatlas.DrawCentered(s, "gem-trash", float64(x), float64(y), 44)
		} else {
			vector.DrawFilledCircle(s, float32(x), float32(y), 20, color.RGBA{55, 70, 88, 255}, false)
		}
	}
	if g.burst > 0 {
		vector.StrokeCircle(s, 240, 360, float32(30+(15-g.burst)*7), 5, color.RGBA{255, 205, 70, 255}, true)
	}
	msg := "TAP / SPACE: SEND A CHAIN"
	if g.clear {
		msg = fmt.Sprintf("DUEL CLEAR! BEST %d / NEXT", g.best)
	}
	ebitenutil.DebugPrintAt(s, msg, 125, 590)
	ebitenutil.DebugPrintAt(s, "Rules change; Update stays the same", 110, 635)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Duel Director — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
