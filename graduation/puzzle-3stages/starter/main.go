package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct{ progress Progress }

func (g *Game) Update() error {
	g.progress.Step(Input{Solve: inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft), Restart: inpututil.IsKeyJustPressed(ebiten.KeyR)})
	return nil
}
func (g *Game) Draw(s *ebiten.Image) {
	current := g.progress.Current()
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PUZZLE %s  %d/%d", current.Name, g.progress.Presses, current.Target), 20, 20)
	ebitenutil.DebugPrintAt(s, "SPACE / CLICK: solve  R: restart", 20, 42)
	if g.progress.Clear {
		ebitenutil.DebugPrintAt(s, "ALL CLEAR! Press R", 85, 100)
	}
}
func (g *Game) Layout(_, _ int) (int, int) { return 320, 180 }
func main() {
	ebiten.SetWindowSize(640, 360)
	if err := ebiten.RunGame(&Game{progress: NewProgress()}); err != nil {
		panic(err)
	}
}
