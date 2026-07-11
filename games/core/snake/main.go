package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height, cell = 480, 720, 24

type point struct{ x, y int }
type game struct {
	body         []point
	dir, next    point
	food         point
	frame, score int
	rng          *rand.Rand
	over         bool
}

func newGame() *game {
	g := &game{body: []point{{10, 14}, {9, 14}, {8, 14}}, dir: point{1, 0}, next: point{1, 0}, rng: rand.New(rand.NewSource(31))}
	g.placeFood()
	return g
}
func (g *game) placeFood() {
	for {
		p := point{g.rng.Intn(width / cell), 2 + g.rng.Intn(height/cell-4)}
		ok := true
		for _, b := range g.body {
			if b == p {
				ok = false
			}
		}
		if ok {
			g.food = p
			return
		}
	}
}
func (g *game) setDir(d point) {
	if d.x+g.dir.x == 0 && d.y+g.dir.y == 0 {
		return
	}
	g.next = d
}
func (g *game) Update() error {
	if g.over {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			*g = *newGame()
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.setDir(point{0, -1})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.setDir(point{0, 1})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.setDir(point{-1, 0})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.setDir(point{1, 0})
	}
	px, py, pressed := 0, 0, false
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		px, py = ebiten.CursorPosition()
		pressed = true
	}
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		px, py = ebiten.TouchPosition(ids[0])
		pressed = true
	}
	if pressed {
		h := g.body[0]
		dx, dy := px-(h.x*cell+cell/2), py-(h.y*cell+cell/2)
		if abs(dx) > abs(dy) {
			if dx < 0 {
				g.setDir(point{-1, 0})
			} else {
				g.setDir(point{1, 0})
			}
		} else {
			if dy < 0 {
				g.setDir(point{0, -1})
			} else {
				g.setDir(point{0, 1})
			}
		}
	}
	g.frame++
	if g.frame%max(4, 10-g.score/3) != 0 {
		return nil
	}
	g.dir = g.next
	head := point{g.body[0].x + g.dir.x, g.body[0].y + g.dir.y}
	if head.x < 0 || head.x >= width/cell || head.y < 2 || head.y >= height/cell-2 {
		g.over = true
		return nil
	}
	for _, b := range g.body {
		if head == b {
			g.over = true
			return nil
		}
	}
	g.body = append([]point{head}, g.body...)
	if head == g.food {
		g.score++
		g.placeFood()
	} else {
		g.body = g.body[:len(g.body)-1]
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{7, 20, 38, 255})
	for x := 0; x < width; x += cell {
		for y := cell * 2; y < height-cell*2; y += cell {
			vector.StrokeRect(s, float32(x), float32(y), cell, cell, 1, color.RGBA{45, 226, 194, 20}, false)
		}
	}
	vector.DrawFilledCircle(s, float32(g.food.x*cell+cell/2), float32(g.food.y*cell+cell/2), 9, color.RGBA{255, 105, 79, 255}, false)
	for i, b := range g.body {
		size := float32(cell - 4)
		c := color.RGBA{45, 226, 194, 255}
		if i == 0 {
			c = color.RGBA{255, 210, 72, 255}
		}
		vector.DrawFilledRect(s, float32(b.x*cell+2), float32(b.y*cell+2), size, size, c, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %02d", g.score), 205, 18)
	if g.over {
		ebitenutil.DebugPrintAt(s, "GAME OVER\n\nTAP / SPACE TO RETRY", 160, 340)
	} else {
		ebitenutil.DebugPrintAt(s, "ARROWS / TAP A DIRECTION", 150, 685)
	}
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Snake — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
