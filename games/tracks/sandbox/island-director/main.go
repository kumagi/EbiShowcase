package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	width  = 480
	height = 720
)

type stage struct {
	name, rule                   string
	sky, ground                  color.RGBA
	wood, stone, crystal, lights int
}

var stages = []stage{{"MOSS CAMP", "Gather one of each resource.", color.RGBA{31, 76, 91, 255}, color.RGBA{55, 102, 76, 255}, 2, 2, 1, 1}, {"CRYSTAL CAVE", "Crystal is plentiful; two crawlers patrol.", color.RGBA{27, 35, 76, 255}, color.RGBA{55, 61, 95, 255}, 2, 2, 3, 1}, {"EMBER ISLE", "Build two lights before the long night.", color.RGBA{92, 43, 46, 255}, color.RGBA{104, 68, 55, 255}, 3, 3, 2, 2}}

type game struct {
	stage, found, score, best int
	done                      bool
}

var best int

func (g *game) advance() {
	if g.done {
		*g = game{best: best}
		return
	}
	g.found++
	g.score += 100 + g.stage*50
	if g.found >= 4+g.stage {
		g.stage++
		g.found = 0
		if g.stage >= len(stages) {
			g.done = true
			g.stage = len(stages) - 1
			if g.score > best {
				best = g.score
			}
			g.best = best
		}
	}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.advance()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	d := stages[g.stage]
	s.Fill(d.sky)
	vector.DrawFilledRect(s, 0, 300, width, 420, d.ground, false)
	ebitenutil.DebugPrintAt(s, "ISLAND DATA DIRECTOR", 157, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ISLAND %d/3  %s", g.stage+1, d.name), 145, 78)
	ebitenutil.DebugPrintAt(s, d.rule, 80, 110)
	for i := 0; i < 10; i++ {
		x := float32(24 + i*46)
		h := float32(70 + (i*g.stage*17+i*9)%100)
		vector.DrawFilledRect(s, x, 480-h, 38, h, color.RGBA{uint8(88 + i*3), uint8(110 + g.stage*18), 92, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DATA: wood %d  stone %d  crystal %d  lights %d", d.wood, d.stone, d.crystal, d.lights), 73, 520)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DISCOVERY %d/%d   SCORE %04d", g.found, 4+g.stage, g.score), 118, 558)
	vector.DrawFilledRect(s, 70, 610, 340, 64, color.RGBA{230, 164, 58, 255}, false)
	label := "TAP / SPACE: DISCOVER"
	if g.done {
		label = "ALL ISLANDS CLEAR — REPLAY"
	}
	ebitenutil.DebugPrintAt(s, label, 135, 637)
	if g.done {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("BEST EXPEDITION %04d", best), 166, 590)
	}
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Island Data Director")
	if err := ebiten.RunGame(&game{best: best}); err != nil {
		panic(err)
	}
}
