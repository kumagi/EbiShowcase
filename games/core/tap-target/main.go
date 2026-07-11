package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type game struct {
	x, y       float64
	radius     float64
	score      int
	framesLeft int
	started    bool
	rng        *rand.Rand
}

func newGame() *game {
	g := &game{radius: 38, framesLeft: 30 * 60, rng: rand.New(rand.NewSource(11))}
	g.moveTarget()
	return g
}

func (g *game) moveTarget() {
	g.x = 60 + g.rng.Float64()*(width-120)
	g.y = 130 + g.rng.Float64()*(height-240)
}

func pressedPosition() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}

func (g *game) Update() error {
	px, py, pressed := pressedPosition()
	if !g.started {
		if pressed {
			g.started = true
		}
		return nil
	}
	if g.framesLeft <= 0 {
		if pressed {
			g.score, g.framesLeft, g.started = 0, 30*60, false
			g.moveTarget()
		}
		return nil
	}
	g.framesLeft--
	if pressed && math.Hypot(float64(px)-g.x, float64(py)-g.y) <= g.radius {
		g.score++
		g.radius = math.Max(22, 38-float64(g.score)/2)
		g.moveTarget()
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})
	for i := 0; i < 12; i++ {
		vector.StrokeCircle(screen, float32(width/2), float32(height/2), float32(55+i*38), 1, color.RGBA{45, 226, 194, 25}, false)
	}
	vector.DrawFilledCircle(screen, float32(g.x), float32(g.y), float32(g.radius+8), color.RGBA{45, 226, 194, 45}, false)
	vector.DrawFilledCircle(screen, float32(g.x), float32(g.y), float32(g.radius), color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledCircle(screen, float32(g.x-10), float32(g.y-10), float32(g.radius/4), color.RGBA{230, 255, 249, 220}, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %02d", g.score), 24, 24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME  %02d", max(0, g.framesLeft/60)), 380, 24)
	if !g.started {
		ebitenutil.DebugPrintAt(screen, "TAP THE TARGET\n\nCLICK / TOUCH TO START", 145, 320)
	} else if g.framesLeft <= 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME UP!  SCORE %d\n\nTAP TO RETRY", g.score), 160, 320)
	}
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Tap the Target — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
