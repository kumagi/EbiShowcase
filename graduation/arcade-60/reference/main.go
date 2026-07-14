package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct{ round Round }

func (g *game) Update() error {
	g.round.Step(Input{
		Action:  inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
		Restart: inpututil.IsKeyJustPressed(ebiten.KeyR),
	})
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("60 SECOND ARCADE  SCORE %d  TIME %02d", g.round.Score, g.round.SecondsLeft()), 20, 20)
	ebitenutil.DebugPrintAt(screen, "SPACE / CLICK: SCORE   R: RETRY", 20, 44)
	if g.round.Over {
		ebitenutil.DebugPrintAt(screen, "TIME! PRESS R TO RETRY", 80, 100)
	}
}

func (g *game) Layout(_, _ int) (int, int) { return 320, 180 }

func main() {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("60 Second Arcade")
	if err := ebiten.RunGame(&game{round: NewRound()}); err != nil {
		panic(err)
	}
}
