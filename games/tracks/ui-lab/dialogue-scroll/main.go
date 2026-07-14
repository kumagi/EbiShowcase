package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/uilab"
	"image/color"
)

type game struct {
	scroll uilab.Scroll
	status uilab.Status
}

func newGame() *game { return &game{scroll: uilab.Scroll{Max: 2}} }
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.scroll.Move(1)
		g.status.Set("Dialogue advanced", 80)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.scroll.Move(-1)
		g.status.Set("Scrolled back", 80)
	}
	g.status.Tick()
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	ebitenutil.DebugPrintAt(s, "UI LAB 05 · DIALOGUE + STATUS", 75, 30)
	uilab.Panel(s, 40, 105, 400, 380, color.RGBA{25, 48, 80, 255}, color.RGBA{80, 220, 190, 255})
	lines := []string{"Ebi: Every action needs words, not only color.", "Player: The story list can scroll.", "Ebi: A status says what just happened.", "Player: Touch and arrows use the same action.", "Ebi: Clear feedback makes menus kinder."}
	for i := 0; i < 3; i++ {
		ebitenutil.DebugPrintAt(s, lines[g.scroll.Offset+i], 65, 165+i*85)
	}
	vector.DrawFilledRect(s, 410, 140, 10, 280, color.RGBA{50, 75, 110, 255}, true)
	vector.DrawFilledRect(s, 410, 140+float32(g.scroll.Offset*70), 10, 70, color.RGBA{255, 210, 85, 255}, true)
	if g.status.Text != "" {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("STATUS: %s", g.status.Text), 120, 540)
	}
	ebitenutil.DebugPrintAt(s, "Down/tap: next · Up: back · text status remains after each action.", 40, 600)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
