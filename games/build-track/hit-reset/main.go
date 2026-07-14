// Build Track STEP 04: hit testing and relocation belong in Update.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 720
	radius       = 58.0
)

var targets = [][2]float64{{240, 330}, {120, 510}, {360, 480}, {240, 210}}

type game struct {
	targetX, targetY float64
	next             int
	score            int
}

func newGame() *game {
	return &game{targetX: targets[0][0], targetY: targets[0][1], next: 1}
}

func justPressedPosition() (x, y int, pressed bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y = ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y = ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}

// hit is pure, so a test can prove the collision rule without opening a game.
func hit(px, py, cx, cy, r float64) bool {
	return math.Hypot(px-cx, py-cy) <= r
}

func (g *game) moveTarget() {
	p := targets[g.next%len(targets)]
	g.targetX, g.targetY = p[0], p[1]
	g.next++
}

// Update reads input, decides the collision, then mutates score and position.
func (g *game) Update() error {
	px, py, pressed := justPressedPosition()
	if pressed && hit(float64(px), float64(py), g.targetX, g.targetY, radius) {
		g.score++
		g.moveTarget()
	}
	return nil
}

// Draw only projects the state Update has already decided.
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 19, 48, 255})
	vector.DrawFilledCircle(screen, float32(g.targetX), float32(g.targetY), radius+10, color.RGBA{46, 230, 200, 45}, false)
	vector.DrawFilledCircle(screen, float32(g.targetX), float32(g.targetY), radius, color.RGBA{46, 230, 200, 255}, false)
	ebitenutil.DebugPrintAt(screen, "TAP THE CIRCLE", 165, 64)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %02d", g.score), 180, 650)
}

func (g *game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Build Track 04 — Hit and Reset")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
