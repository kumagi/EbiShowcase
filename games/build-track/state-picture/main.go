// Build Track STEP 02: Draw projects a position stored in game.
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 720
)

// game is the single source of truth for the square's position.
type game struct {
	x int
	y int
}

func newGame() *game {
	return &game{x: 184, y: 316}
}

// Update owns future state changes. This step deliberately changes nothing.
func (g *game) Update() error {
	return nil
}

// Draw reads game.x and game.y, then projects them as pixels.
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 19, 48, 255})
	vector.DrawFilledRect(screen, float32(g.x), float32(g.y), 112, 112, color.RGBA{46, 230, 200, 255}, false)
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Build Track 02 — State Picture")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
