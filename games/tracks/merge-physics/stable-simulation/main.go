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

const width, height, cell = 480, 720, 80

type ball struct {
	x, y, vx, vy, r, mass float64
	quiet                 int
	sleep                 bool
	c                     color.RGBA
}
type game struct {
	balls         []ball
	rng           *rand.Rand
	drops, checks int
	clear         bool
}

func newGame() *game { return &game{rng: rand.New(rand.NewSource(4604))} }
func (g *game) Update() error {
	if g.clear {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	if x, ok := drop(); ok && g.drops < 30 {
		r := float64(12 + g.rng.Intn(16))
		g.balls = append(g.balls, ball{x: float64(x), y: 85, r: r, mass: r * r, c: color.RGBA{uint8(80 + g.rng.Intn(160)), uint8(90 + g.rng.Intn(150)), uint8(110 + g.rng.Intn(130)), 255}})
		g.drops++
	}
	for i := range g.balls {
		b := &g.balls[i]
		if b.sleep {
			continue
		}
		b.vy += .45
		b.x += b.vx
		b.y += b.vy
		b.vx *= .995
		if b.x-b.r < 25 {
			b.x = 25 + b.r
			b.vx = math.Abs(b.vx) * .2
		}
		if b.x+b.r > 455 {
			b.x = 455 - b.r
			b.vx = -math.Abs(b.vx) * .2
		}
		if b.y+b.r > 660 {
			b.y = 660 - b.r
			if b.vy > 1 {
				b.vy = -b.vy * .12
			} else {
				b.vy = 0
			}
			b.vx *= .8
		}
	}
	g.checks = 0
	for pass := 0; pass < 7; pass++ {
		buckets := map[[2]int][]int{}
		for i, b := range g.balls {
			k := [2]int{int(b.x) / cell, int(b.y) / cell}
			buckets[k] = append(buckets[k], i)
		}
		for i := range g.balls {
			a := &g.balls[i]
			cx, cy := int(a.x)/cell, int(a.y)/cell
			for oy := -1; oy <= 1; oy++ {
				for ox := -1; ox <= 1; ox++ {
					for _, j := range buckets[[2]int{cx + ox, cy + oy}] {
						if j <= i {
							continue
						}
						g.checks++
						resolve(a, &g.balls[j])
					}
				}
			}
		}
	}
	sleeping := 0
	for i := range g.balls {
		b := &g.balls[i]
		speed := math.Hypot(b.vx, b.vy)
		if speed < .12 && b.y+b.r > 150 {
			b.quiet++
		} else {
			b.quiet = 0
			b.sleep = false
		}
		if b.quiet > 90 {
			b.sleep = true
			b.vx, b.vy = 0, 0
		}
		if b.sleep {
			sleeping++
		}
	}
	if g.drops == 30 && sleeping == 30 {
		g.clear = true
	}
	return nil
}
func resolve(a, b *ball) {
	dx, dy := b.x-a.x, b.y-a.y
	d := math.Hypot(dx, dy)
	minD := a.r + b.r
	if d <= 0 || d >= minD {
		return
	}
	a.sleep, b.sleep = false, false
	a.quiet, b.quiet = 0, 0
	nx, ny := dx/d, dy/d
	invA, invB := 1/a.mass, 1/b.mass
	sum := invA + invB
	over := minD - d
	a.x -= nx * over * invA / sum
	a.y -= ny * over * invA / sum
	b.x += nx * over * invB / sum
	b.y += ny * over * invB / sum
	rel := (b.vx-a.vx)*nx + (b.vy-a.vy)*ny
	if rel < 0 {
		j := -(1 + .08) * rel / sum
		a.vx -= j * nx * invA
		a.vy -= j * ny * invA
		b.vx += j * nx * invB
		b.vy += j * ny * invB
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 28, 44, 255})
	vector.DrawFilledRect(s, 20, 70, 440, 610, color.RGBA{35, 44, 61, 255}, false)
	for x := 20; x < 460; x += cell {
		vector.StrokeLine(s, float32(x), 70, float32(x), 680, 1, color.RGBA{100, 140, 170, 35}, false)
	}
	for y := 70; y < 680; y += cell {
		vector.StrokeLine(s, 20, float32(y), 460, float32(y), 1, color.RGBA{100, 140, 170, 35}, false)
	}
	sleeping := 0
	for _, b := range g.balls {
		vector.DrawFilledCircle(s, float32(b.x), float32(b.y), float32(b.r), b.c, false)
		if b.sleep {
			sleeping++
			vector.StrokeCircle(s, float32(b.x), float32(b.y), float32(b.r+2), 2, color.RGBA{100, 180, 255, 255}, false)
		}
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BALLS %02d/30   SLEEP %02d   NEAR CHECKS %04d", g.drops, sleeping, g.checks), 75, 28)
	ebitenutil.DebugPrintAt(s, "TAP TO DROP — BLUE RING MEANS SLEEPING", 95, 690)
	if g.clear {
		overlay(s, "THIRTY CIRCLES SLEEPING!\n\nTAP / SPACE TO RESET")
	}
}
func drop() (int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		return max(45, min(435, x)), true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		return max(45, min(435, x)), true
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
	ebitenutil.DebugPrintAt(s, msg, 115, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Stable Simulation — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
