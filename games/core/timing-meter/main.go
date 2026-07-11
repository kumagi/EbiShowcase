package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type game struct {
	marker, speed float64
	score, round  int
	stopped       bool
}

func pressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Update() error {
	if pressed() {
		if g.stopped {
			g.stopped = false
			g.round++
			g.speed = 3.2 + float64(g.round)*0.18
			if g.round%2 == 1 {
				g.speed = -g.speed
			}
		} else {
			g.stopped = true
			distance := math.Abs(g.marker - width/2)
			switch {
			case distance <= 8:
				g.score += 100
			case distance <= 28:
				g.score += 50
			case distance <= 55:
				g.score += 10
			}
		}
	}
	if g.stopped {
		return nil
	}
	g.marker += g.speed
	if g.marker < 45 || g.marker > width-45 {
		g.speed = -g.speed
		g.marker = math.Max(45, math.Min(width-45, g.marker))
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})
	vector.DrawFilledRect(screen, 35, 310, width-70, 80, color.RGBA{16, 44, 69, 255}, false)
	vector.DrawFilledRect(screen, width/2-55, 310, 110, 80, color.RGBA{255, 105, 79, 255}, false)
	vector.DrawFilledRect(screen, width/2-28, 310, 56, 80, color.RGBA{255, 205, 69, 255}, false)
	vector.DrawFilledRect(screen, width/2-8, 310, 16, 80, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledRect(screen, float32(g.marker-4), 285, 8, 130, color.White, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %04d   ROUND %02d", g.score, g.round+1), 145, 35)
	if g.stopped {
		d := math.Abs(g.marker - width/2)
		label := "MISS"
		if d <= 8 {
			label = "PERFECT +100"
		} else if d <= 28 {
			label = "GREAT +50"
		} else if d <= 55 {
			label = "GOOD +10"
		}
		ebitenutil.DebugPrintAt(screen, label+"\n\nTAP FOR NEXT ROUND", 165, 470)
	} else {
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE TO STOP", 165, 470)
	}
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Stop the Meter — Ebitengine")
	if err := ebiten.RunGame(&game{marker: 45, speed: 3.2}); err != nil {
		panic(err)
	}
}
