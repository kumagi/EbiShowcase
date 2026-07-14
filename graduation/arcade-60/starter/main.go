package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game keeps Ebitengine input and drawing at the edge. The testable rules are
// in Round.Step in rules.go.
type Game struct{ round Round }

func (g *Game) Update() error {
	action := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	restart := inpututil.IsKeyJustPressed(ebiten.KeyR)
	// TODO 6: Add an on-screen touch button if you want a phone-friendly action.
	g.round.Step(Input{Action: action, Restart: restart})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO 7: Replace DebugPrint with your own score and result UI.
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ARCADE 60  SCORE %d  TIME %02d", g.round.Score, g.round.SecondsLeft()), 18, 18)
	ebitenutil.DebugPrintAt(screen, "SPACE / CLICK: score   R: restart", 18, 42)
	if g.round.Over {
		ebitenutil.DebugPrintAt(screen, "TIME! Press R to play again", 70, 100)
	}
}

func (g *Game) Layout(_, _ int) (int, int) { return 320, 180 }

func main() {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("Arcade 60 starter")
	if err := ebiten.RunGame(&Game{round: NewRound()}); err != nil {
		panic(err)
	}
}
