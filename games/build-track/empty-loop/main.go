// Build Track STEP 01: an empty Ebitengine game loop.
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 480
	screenHeight = 720
)

// game will hold the game's numbers in the next steps.
type game struct{}

// Update is where input and game-state changes will belong.
func (g *game) Update() error {
	return nil
}

// Draw only paints the current game state. There is no state to paint yet.
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 19, 48, 255})
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Build Track 01 — Empty Loop")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
