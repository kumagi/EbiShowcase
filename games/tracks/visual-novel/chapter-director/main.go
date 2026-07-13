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
	W = 480
	H = 720
)

type chapter struct {
	name, place, rule string
	scenes            int
	c                 color.RGBA
}

var chapters = []chapter{{"THE MISSING LIGHT", "Harbor", "Meet the cast and find one clue.", 3, color.RGBA{35, 74, 95, 255}}, {"THE STOPPED CLOCK", "Market", "Choices change flags and available scenes.", 3, color.RGBA{81, 59, 113, 255}}, {"THE STORM RING", "Observatory", "Combine flags into one of four endings.", 3, color.RGBA{51, 55, 103, 255}}}

type game struct {
	chapter, scene, score, best int
	ended                       bool
}

var best int

func (g *game) step() {
	if g.ended {
		*g = game{best: best}
		return
	}
	g.scene++
	g.score += 100 + g.chapter*50
	if g.scene >= chapters[g.chapter].scenes {
		g.chapter++
		g.scene = 0
		if g.chapter == 3 {
			g.chapter = 2
			g.ended = true
			if g.score > best {
				best = g.score
			}
			g.best = best
		}
	}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.step()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	d := chapters[g.chapter]
	s.Fill(d.c)
	ebitenutil.DebugPrintAt(s, "CHAPTER DATA DIRECTOR", 149, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("CHAPTER %d/3  %s", g.chapter+1, d.name), 110, 82)
	ebitenutil.DebugPrintAt(s, d.place, 205, 120)
	vector.DrawFilledRect(s, 45, 190, 390, 250, color.RGBA{8, 16, 35, 210}, false)
	ebitenutil.DebugPrintAt(s, d.rule, 82, 235)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCENE %d/%d", g.scene, d.scenes), 190, 300)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("RUN SCORE %04d", g.score), 180, 345)
	if g.ended {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("ALL CHAPTERS COMPLETE / BEST %04d", best), 107, 400)
	}
	vector.DrawFilledRect(s, 70, 610, 340, 64, color.RGBA{230, 164, 58, 255}, false)
	label := "TAP / SPACE: NEXT SCENE"
	if g.ended {
		label = "START ANOTHER STORY RUN"
	}
	ebitenutil.DebugPrintAt(s, label, 130, 637)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
