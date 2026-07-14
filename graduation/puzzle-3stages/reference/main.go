package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	stage, presses int
	clear          bool
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.clear {
			*g = game{stage: 1}
			return nil
		}
		g.presses++
		if g.presses >= g.stage+1 {
			g.presses = 0
			g.stage++
			if g.stage > 3 {
				g.clear = true
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("THREE STAGE PUZZLE  STAGE %d/3\nSPACE: solve this tiny puzzle", min(3, g.stage)), 20, 20)
	if g.clear {
		ebitenutil.DebugPrintAt(s, "ALL PUZZLES CLEAR! SPACE TO REPLAY", 40, 100)
	}
}
func (g *game) Layout(_, _ int) (int, int) { return 320, 180 }
func main()                                { _ = ebiten.RunGame(&game{stage: 1}) }
