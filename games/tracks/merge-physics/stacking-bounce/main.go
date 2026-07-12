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

type ball struct {
	x, y, vx, vy, r, mass float64
	c                     color.RGBA
}
type game struct {
	balls        []ball
	rng          *rand.Rand
	drops, quiet int
	clear        bool
}

func newGame() *game { return &game{rng: rand.New(rand.NewSource(4503))} }
func (g *game) Update() error {
	if g.clear {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	if x, ok := drop(); ok && g.drops < 12 {
		r := float64(18 + g.rng.Intn(25))
		g.balls = append(g.balls, ball{float64(x), 90, 0, 0, r, r * r, color.RGBA{uint8(80 + g.rng.Intn(160)), uint8(90 + g.rng.Intn(150)), uint8(100 + g.rng.Intn(140)), 255}})
		g.drops++
	}
	for i := range g.balls {
		b := &g.balls[i]
		b.vy += .48
		b.vx *= .998
		b.x += b.vx
		b.y += b.vy
		if b.x-b.r < 25 {
			b.x = 25 + b.r
			b.vx = math.Abs(b.vx) * .45
		}
		if b.x+b.r > 455 {
			b.x = 455 - b.r
			b.vx = -math.Abs(b.vx) * .45
		}
		if b.y+b.r > 660 {
			b.y = 660 - b.r
			if b.vy > 1 {
				b.vy = -b.vy * .25
			} else {
				b.vy = 0
			}
			b.vx *= .88
		}
	}
	for iteration := 0; iteration < 5; iteration++ {
		for i := 0; i < len(g.balls); i++ {
			for j := i + 1; j < len(g.balls); j++ {
				resolve(&g.balls[i], &g.balls[j])
			}
		}
	}
	energy := 0.0
	for _, b := range g.balls {
		energy += math.Abs(b.vx) + math.Abs(b.vy)
	}
	if g.drops == 12 && energy < 4 {
		g.quiet++
	} else {
		g.quiet = 0
	}
	if g.quiet > 90 {
		g.clear = true
	}
	return nil
}
func resolve(a, b *ball) {
	dx, dy := b.x-a.x, b.y-a.y
	d := math.Hypot(dx, dy)
	minimum := a.r + b.r
	if d <= 0 || d >= minimum {
		return
	}
	nx, ny := dx/d, dy/d
	overlap := minimum - d
	invA, invB := 1/a.mass, 1/b.mass
	sum := invA + invB
	a.x -= nx * overlap * invA / sum
	a.y -= ny * overlap * invA / sum
	b.x += nx * overlap * invB / sum
	b.y += ny * overlap * invB / sum
	relative := (b.vx-a.vx)*nx + (b.vy-a.vy)*ny
	if relative >= 0 {
		return
	}
	impulse := -(1 + .25) * relative / sum
	a.vx -= impulse * nx * invA
	a.vy -= impulse * ny * invA
	b.vx += impulse * nx * invB
	b.vy += impulse * ny * invB
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 28, 44, 255})
	vector.DrawFilledRect(s, 20, 70, 440, 610, color.RGBA{35, 44, 61, 255}, false)
	vector.StrokeRect(s, 20, 70, 440, 610, 5, color.RGBA{110, 128, 151, 255}, false)
	for _, b := range g.balls {
		vector.DrawFilledCircle(s, float32(b.x), float32(b.y), float32(b.r), b.c, false)
		vector.StrokeCircle(s, float32(b.x), float32(b.y), float32(b.r), 2, color.White, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BALLS %02d/12   SOLVER PASSES 5   QUIET %02d/90", g.drops, g.quiet), 75, 28)
	ebitenutil.DebugPrintAt(s, "TAP TO DROP — LARGE CIRCLES HAVE MORE MASS", 85, 690)
	if g.clear {
		overlay(s, "THE STACK IS STABLE!\n\nTAP / SPACE TO RESET")
	}
}
func drop() (int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		return max(55, min(425, x)), true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		return max(55, min(425, x)), true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return 240, true
	}
	return 0, false
}
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Stacking and Bounce — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
