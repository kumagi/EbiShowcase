package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	x, y, score, frames int
	over                bool
}

func (g *game) Update() error {
	if g.over {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*g = game{x: 120, y: 120}
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.x -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.x += 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.y -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.y += 3
	}
	g.frames++
	if g.x > 250 {
		g.score += 10
		g.x = 30
	}
	if g.frames >= 60*60 {
		g.over = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("60 SECOND ARCADE  SCORE %d  TIME %d", g.score, max(0, 60-g.frames/60)), 20, 20)
	if g.over {
		ebitenutil.DebugPrintAt(s, "TIME! SPACE TO RETRY", 80, 100)
	}
}
func (g *game) Layout(_, _ int) (int, int) { return 320, 180 }
func main()                                { ebiten.SetWindowSize(640, 360); _ = ebiten.RunGame(&game{x: 120, y: 120}) }
