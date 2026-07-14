package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct{ adventure State }

func (g *Game) Update() error {
	// TODO 6: Make a touch-friendly on-screen interaction button.
	g.adventure.Step(Input{
		PickUpKey: inpututil.IsKeyJustPressed(ebiten.KeyK),
		NextRoom:  inpututil.IsKeyJustPressed(ebiten.KeySpace),
		OpenExit:  inpututil.IsKeyJustPressed(ebiten.KeyO),
		Restart:   inpututil.IsKeyJustPressed(ebiten.KeyR),
	})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO 7: Draw each room and the key instead of using DebugPrint.
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("THREE ROOMS  ROOM %d  KEY %v", g.adventure.Room, g.adventure.HasKey), 18, 18)
	ebitenutil.DebugPrintAt(screen, "K: key  SPACE: next  O: exit  R: restart", 18, 42)
	if g.adventure.Escaped {
		ebitenutil.DebugPrintAt(screen, "ESCAPED! Press R for a new adventure", 55, 100)
	}
}

func (g *Game) Layout(_, _ int) (int, int) { return 320, 180 }

func main() {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("Three rooms starter")
	if err := ebiten.RunGame(&Game{adventure: NewAdventure()}); err != nil {
		panic(err)
	}
}
