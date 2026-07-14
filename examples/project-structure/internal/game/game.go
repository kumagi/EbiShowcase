// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/examples/project-structure/internal/rules"
)

type Game struct{ score, gems int }

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.score, g.gems = rules.Collect(g.score, g.gems)
	}
	return nil
}
func (g *Game) Draw(*ebiten.Image)         {}
func (g *Game) Layout(_, _ int) (int, int) { return 320, 180 }
func Run()                                 { _ = ebiten.RunGame(&Game{}) }
