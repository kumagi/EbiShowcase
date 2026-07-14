package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	room  int
	key   bool
	clear bool
}

func (g *game) Update() error {
	if g.clear {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*g = game{room: 1}
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.room == 1 {
			g.key = true
			g.room = 2
		} else if g.room == 2 {
			g.room = 3
		} else if g.key {
			g.clear = true
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("THREE ROOM EXPLORE  ROOM %d  KEY %v\nSPACE: interact", g.room, g.key), 20, 20)
	if g.clear {
		ebitenutil.DebugPrintAt(s, "ESCAPED! SPACE TO RESTART", 60, 100)
	}
}
func (g *game) Layout(_, _ int) (int, int) { return 320, 180 }
func main()                                { _ = ebiten.RunGame(&game{room: 1}) }
