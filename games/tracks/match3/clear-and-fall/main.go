package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

const width, height = 480, 720
const cols, rows, cell, ox, oy = 6, 7, 64, 48, 130
const empty = -1

var colors = []color.RGBA{{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255}, {105, 194, 119, 255}, {177, 94, 218, 255}}

type point struct{ x, y int }
type game struct {
	board                [rows][cols]int
	marked               map[point]bool
	rng                  *rand.Rand
	phase, cycles, score int
	message              string
	clear                bool
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(5803)), marked: map[point]bool{}}
	g.seed()
	return g
}
func (g *game) seed() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = g.rng.Intn(len(colors))
		}
	}
	kind := g.cycles % 5
	for x := 1; x < 4; x++ {
		g.board[5][x] = kind
	}
	g.marked = scan(g.board)
	g.phase = 0
	g.message = "STEP 1: clear the outlined match."
}
func (g *game) Update() error {
	if g.clear {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	goNext := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if _, y, ok := press(); ok && y > 610 {
		goNext = true
	}
	if goNext {
		switch g.phase {
		case 0:
			for p := range g.marked {
				g.board[p.y][p.x] = empty
				g.score += 10
			}
			g.phase = 1
			g.message = "STEP 2: compact each column downward."
		case 1:
			g.fall()
			g.phase = 2
			g.message = "STEP 3: refill empty cells from the top."
		case 2:
			g.refill()
			g.cycles++
			if g.cycles >= 3 {
				g.clear = true
			} else {
				g.seed()
			}
		}
	}
	return nil
}
func (g *game) fall() {
	for x := 0; x < cols; x++ {
		write := rows - 1
		for y := rows - 1; y >= 0; y-- {
			if g.board[y][x] != empty {
				g.board[write][x] = g.board[y][x]
				if write != y {
					g.board[y][x] = empty
				}
				write--
			}
		}
	}
}
func (g *game) refill() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if g.board[y][x] == empty {
				g.board[y][x] = g.rng.Intn(len(colors))
			}
		}
	}
}
func scan(b [rows][cols]int) map[point]bool {
	m := map[point]bool{}
	for y := 0; y < rows; y++ {
		start := 0
		for x := 1; x <= cols; x++ {
			if x == cols || b[y][x] != b[y][start] {
				if x-start >= 3 {
					for i := start; i < x; i++ {
						m[point{i, y}] = true
					}
				}
				start = x
			}
		}
	}
	for x := 0; x < cols; x++ {
		start := 0
		for y := 1; y <= rows; y++ {
			if y == rows || b[y][x] != b[start][x] {
				if y-start >= 3 {
					for i := start; i < y; i++ {
						m[point{x, i}] = true
					}
				}
				start = y
			}
		}
	}
	return m
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(s, "CLEAR / FALL / REFILL", 165, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("CYCLE %d/3   SCORE %03d   PHASE %d/3", g.cycles+1, g.score, g.phase+1), 120, 70)
	ebitenutil.DebugPrintAt(s, g.message, 65, 100)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			if g.board[y][x] == empty {
				vector.StrokeRect(s, px+4, py+4, cell-8, cell-8, 2, color.RGBA{70, 83, 105, 255}, false)
				continue
			}
			vector.DrawFilledRect(s, px+3, py+3, cell-6, cell-6, colors[g.board[y][x]], false)
			if g.phase == 0 && g.marked[point{x, y}] {
				vector.StrokeRect(s, px+2, py+2, cell-4, cell-4, 6, color.White, false)
			}
		}
	}
	vector.DrawFilledRect(s, 55, 615, 370, 65, color.RGBA{240, 177, 65, 255}, false)
	labels := []string{"CLEAR MATCH [SPACE]", "FALL DOWN [SPACE]", "REFILL TOP [SPACE]"}
	ebitenutil.DebugPrintAt(s, labels[g.phase], 150, 642)
	if g.clear {
		overlay(s, "THREE CYCLES COMPLETE!\n\nTAP / SPACE TO RESTART")
	}
}
func press() (int, int, bool) {
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
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, m string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, m, 105, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Clear and Fall — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
