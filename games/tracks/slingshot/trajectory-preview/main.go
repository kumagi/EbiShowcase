package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const width, height = 480, 720

type game struct {
	x, y, vx, vy float64
	drag         bool
	dx, dy       float64
	shots        int
}

func (g *game) Update() error {
	if !g.drag && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if math.Hypot(float64(x)-g.x, float64(y)-g.y) < 55 {
			g.drag = true
		}
	}
	if !g.drag {
		if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
			x, y := ebiten.TouchPosition(ids[0])
			if math.Hypot(float64(x)-g.x, float64(y)-g.y) < 55 {
				g.drag = true
			}
		}
	}
	if g.drag {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			g.dx = g.x - float64(x)
			g.dy = g.y - float64(y)
		}
		if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
			x, y := ebiten.TouchPosition(ids[0])
			g.dx = g.x - float64(x)
			g.dy = g.y - float64(y)
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) || len(inpututil.AppendJustReleasedTouchIDs(nil)) > 0 {
			g.vx = g.dx * .09
			g.vy = g.dy * .09
			g.drag = false
			g.shots++
		}
	} else {
		g.x += g.vx
		g.y += g.vy
		g.vx *= .985
		g.vy *= .985
		if g.x < 25 || g.x > 455 {
			g.vx = -g.vx
		}
		if g.y < 100 || g.y > 600 {
			g.vy = -g.vy
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{15, 31, 50, 255})
	vector.StrokeRect(s, 20, 90, 440, 530, 4, color.RGBA{80, 145, 170, 255}, false)
	if g.drag {
		for i := 1; i <= 12; i++ {
			t := float64(i) * .5
			vector.DrawFilledCircle(s, float32(g.x+g.dx*.09*t), float32(g.y+g.dy*.09*t), 3, color.RGBA{255, 255, 255, uint8(230 - i*12)}, true)
		}
	}
	trackatlas.DrawCentered(s, "ally", g.x, g.y, 44)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SHOTS %d  POWER %.1f", g.shots, math.Hypot(g.dx, g.dy)), 165, 35)
	ebitenutil.DebugPrintAt(s, "DRAG ALLY — DOTS PREVIEW THE FIRST FRAMES", 85, 675)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Trajectory Preview — Ebitengine")
	if err := ebiten.RunGame(&game{x: 240, y: 520}); err != nil {
		panic(err)
	}
}
