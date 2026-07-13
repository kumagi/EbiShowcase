package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const width, height = 480, 720

type game struct {
	x, y, vx, vy float64
	grounded     bool
	tick, frame  int
	stars        [5]bool
}

func (g *game) Update() error {
	left := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	right := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				left = true
			} else {
				right = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			jump = true
		}
	}
	if left {
		g.vx -= .65
	}
	if right {
		g.vx += .65
	}
	if !left && !right {
		g.vx *= .78
	}
	g.vx = math.Max(-6, math.Min(6, g.vx))
	if jump && g.grounded {
		g.vy = -12
		g.grounded = false
	}
	g.vy += .65
	g.x += g.vx
	g.y += g.vy
	if g.y >= 570 {
		g.y, g.vy, g.grounded = 570, 0, true
	}
	if g.x < 0 {
		g.x = 430
	}
	if g.x > 430 {
		g.x = 0
	}
	g.tick++
	if !g.grounded {
		g.frame = 3
	} else if math.Abs(g.vx) > .4 {
		g.frame = (g.tick / 7) % 3
	} else {
		g.frame = 0
	}
	for i := range g.stars {
		sx := 55 + float64(i)*92
		if !g.stars[i] && math.Abs((g.x+24)-sx) < 25 && math.Abs((g.y+24)-555) < 35 {
			g.stars[i] = true
		}
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{72, 154, 210, 255})
	vector.DrawFilledRect(s, 0, 620, width, 100, color.RGBA{42, 91, 67, 255}, false)
	for i, taken := range g.stars {
		if !taken {
			trackatlas.DrawCentered(s, "coin", 55+float64(i)*92, 555, 24)
		}
	}
	sizes := []float64{46, 42, 47, 43}
	bob := []float64{0, 4, 0, -5}
	trackatlas.DrawCentered(s, "hero", g.x+24, g.y+24+bob[g.frame], sizes[g.frame])
	for i := 0; i < 4; i++ {
		x := 70 + float64(i)*92
		c := color.RGBA{21, 36, 63, 170}
		if i == g.frame {
			c = color.RGBA{255, 204, 61, 230}
		}
		vector.DrawFilledRect(s, float32(x-28), 70, 56, 70, c, false)
		trackatlas.DrawCentered(s, "hero", x, 104+bob[i], sizes[i]*.75)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("FRAME %d   COLLECTED %d/5", g.frame, count(g.stars)), 145, 28)
	ebitenutil.DebugPrintAt(s, "IDLE -> STEP -> SQUASH -> AIR", 120, 155)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH   JUMP: SPACE OR UPPER TOUCH", 45, 685)
}
func count(a [5]bool) int {
	n := 0
	for _, v := range a {
		if v {
			n++
		}
	}
	return n
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Run Animation Lab — Ebitengine")
	if err := ebiten.RunGame(&game{x: 30, y: 570}); err != nil {
		panic(err)
	}
}
