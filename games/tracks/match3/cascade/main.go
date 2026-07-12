package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const width, height = 480, 720
const cols, rows, cell, ox, oy = 6, 7, 64, 48, 130
const empty = -1

var colors = []color.RGBA{{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255}, {105, 194, 119, 255}, {177, 94, 218, 255}}

type point struct{ x, y int }
type game struct {
	board               [rows][cols]int
	marked              map[point]bool
	combo, score, round int
	settled, clear      bool
	message             string
}

func newGame() *game { g := &game{}; g.seed(); return g }
func (g *game) seed() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = (x*2 + y) % 5
		}
	}
	for x := 1; x < 4; x++ {
		g.board[5][x] = 0
	}
	g.marked = scan(g.board)
	g.combo = 0
	g.settled = false
	g.message = "Press RESOLVE to clear the first match."
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
		if g.settled {
			g.round++
			if g.round >= 3 {
				g.clear = true
			} else {
				g.seed()
				g.message = "New move: resolve every cascade again."
			}
		} else {
			g.resolveOne()
		}
	}
	return nil
}
func (g *game) resolveOne() {
	count := len(g.marked)
	g.combo++
	g.score += count * 10 * g.combo
	for p := range g.marked {
		g.board[p.y][p.x] = empty
	}
	g.fall()
	g.refill()
	g.marked = scan(g.board)
	if len(g.marked) == 0 {
		g.settled = true
		g.message = fmt.Sprintf("Settled after %d-chain! Press NEXT MOVE.", g.combo)
	} else {
		g.message = fmt.Sprintf("New match found: %d-chain, multiplier x%d!", g.combo+1, g.combo+1)
	}
}
func (g *game) fall() {
	for x := 0; x < cols; x++ {
		w := rows - 1
		for y := rows - 1; y >= 0; y-- {
			if g.board[y][x] != empty {
				g.board[w][x] = g.board[y][x]
				if w != y {
					g.board[y][x] = empty
				}
				w--
			}
		}
	}
}
func (g *game) refill() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if g.board[y][x] == empty {
				g.board[y][x] = (x*2 + y + g.combo) % 5
			}
		}
	}
	if g.combo < 3 {
		kind := g.combo % 5
		for x := 0; x < 3; x++ {
			g.board[0][x] = kind
		}
	}
}
func scan(b [rows][cols]int) map[point]bool {
	m := map[point]bool{}
	for y := 0; y < rows; y++ {
		st := 0
		for x := 1; x <= cols; x++ {
			if x == cols || b[y][x] != b[y][st] {
				if b[y][st] != empty && x-st >= 3 {
					for i := st; i < x; i++ {
						m[point{i, y}] = true
					}
				}
				st = x
			}
		}
	}
	for x := 0; x < cols; x++ {
		st := 0
		for y := 1; y <= rows; y++ {
			if y == rows || b[y][x] != b[st][x] {
				if b[st][x] != empty && y-st >= 3 {
					for i := st; i < y; i++ {
						m[point{x, i}] = true
					}
				}
				st = y
			}
		}
	}
	return m
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(s, "CASCADE REACTOR", 180, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MOVE %d/3   CHAIN x%d   SCORE %04d", g.round+1, max(1, g.combo), g.score), 120, 70)
	ebitenutil.DebugPrintAt(s, g.message, 55, 100)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.DrawFilledRect(s, px+3, py+3, cell-6, cell-6, colors[g.board[y][x]], false)
			if g.marked[point{x, y}] {
				vector.StrokeRect(s, px+2, py+2, cell-4, cell-4, 6, color.White, false)
			}
		}
	}
	vector.DrawFilledRect(s, 55, 615, 370, 65, color.RGBA{240, 177, 65, 255}, false)
	label := "RESOLVE NEXT MATCH [SPACE]"
	if g.settled {
		label = "NEXT MOVE [SPACE]"
	}
	ebitenutil.DebugPrintAt(s, label, 135, 642)
	if g.clear {
		overlay(s, "CASCADE MASTERED!\n\nTAP / SPACE TO RESTART")
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
	ebitenutil.DebugPrintAt(s, m, 115, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cascade Reactor — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
