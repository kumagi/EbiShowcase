package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct{ adventure State }

func (g *game) Update() error {
	g.adventure.Step(Input{PickUpKey: inpututil.IsKeyJustPressed(ebiten.KeyK), NextRoom: inpututil.IsKeyJustPressed(ebiten.KeySpace), OpenExit: inpututil.IsKeyJustPressed(ebiten.KeyO), Restart: inpututil.IsKeyJustPressed(ebiten.KeyR)})
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("THREE ROOMS  ROOM %d  KEY %v\nK: key  SPACE: next  O: exit", g.adventure.Room, g.adventure.HasKey), 20, 20)
	if g.adventure.Escaped {
		ebitenutil.DebugPrintAt(s, "ESCAPED! R TO RESTART", 60, 100)
	}
}
func (g *game) Layout(_, _ int) (int, int) { return 320, 180 }
func main()                                { ebiten.SetWindowSize(640, 360); _ = ebiten.RunGame(&game{adventure: NewAdventure()}) }
