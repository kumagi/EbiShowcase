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

type level struct {
	name        string
	moves, goal int
	bg          color.RGBA
}

var levels = []level{{"CORAL COVE", 8, 300, color.RGBA{15, 35, 55, 255}}, {"TIDE TEMPLE", 7, 500, color.RGBA{10, 55, 60, 255}}, {"STARLIGHT REEF", 6, 700, color.RGBA{35, 20, 65, 255}}}

type game struct {
	level, moves, score, best, tick, burst int
	cleared                                bool
}

func (g *game) reset(i int) {
	g.level = i
	g.moves = levels[i].moves
	g.score = 0
	g.cleared = false
	g.burst = 0
}
func (g *game) Update() error {
	g.tick++
	if g.burst > 0 {
		g.burst--
	}
	pressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if !pressed {
		return nil
	}
	if g.cleared {
		if g.level < 2 {
			g.reset(g.level + 1)
		} else {
			g.reset(0)
		}
		return nil
	}
	if g.moves > 0 {
		g.moves--
		g.score += 80 + (g.tick%4)*20
		g.burst = 16
		if g.score >= levels[g.level].goal {
			g.cleared = true
			if g.score+g.moves*40 > g.best {
				g.best = g.score + g.moves*40
			}
		}
	}
	if g.moves == 0 && !g.cleared {
		g.reset(g.level)
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	l := levels[g.level]
	s.Fill(l.bg)
	vector.DrawFilledRect(s, 55, 125, 370, 440, color.RGBA{8, 18, 34, 180}, false)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			bob := math.Sin(float64(g.tick)*.06+float64(x+y)) * .8
			trackatlas.DrawCentered(s, trackatlas.Gem((x+y+g.level)%5), 132+float64(x)*72, 220+float64(y)*72+bob, 55)
		}
	}
	p := float32(g.score) / float32(l.goal)
	if p > 1 {
		p = 1
	}
	vector.DrawFilledRect(s, 80, 505, 320, 18, color.RGBA{45, 60, 80, 255}, false)
	vector.DrawFilledRect(s, 80, 505, 320*p, 18, color.RGBA{245, 190, 70, 255}, false)
	if g.burst > 0 {
		vector.StrokeCircle(s, 240, 350, float32(25+(16-g.burst)*8), 5, color.RGBA{255, 220, 100, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("REEF %d/3 / %s", g.level+1, l.name), 145, 65)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MOVES %d   SCORE %d/%d", g.moves, g.score, l.goal), 125, 95)
	msg := "TAP / SPACE: MAKE A MATCH"
	if g.cleared {
		msg = fmt.Sprintf("CLEAR! BEST %d / NEXT REEF", g.best)
	}
	ebitenutil.DebugPrintAt(s, msg, 120, 625)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	x := &game{}
	x.reset(0)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Reef Director — Ebitengine")
	if err := ebiten.RunGame(x); err != nil {
		panic(err)
	}
}
