// Build Track STEP 03: input and score changes belong in Update.
package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 720
)

type game struct {
	score int
}

func justPressed() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

// addScore is pure: the same old score always makes the same new score.
func addScore(score int) int {
	return score + 1
}

// Update owns input and the score rule.
func (g *game) Update() error {
	if justPressed() {
		g.score = addScore(g.score)
	}
	return nil
}

// Draw only projects the score already decided by Update.
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 19, 48, 255})
	vector.DrawFilledCircle(screen, screenWidth/2, screenHeight/2, 76, color.RGBA{46, 230, 200, 255}, false)
	ebitenutil.DebugPrintAt(screen, "TAP / CLICK ANYWHERE", 132, 250)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %02d", g.score), 180, 470)
}

func (g *game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Build Track 03 — Tap Score")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
