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

const width, height = 480, 720
const cols, rows, cell, ox, oy = 6, 7, 64, 48, 130

var pieceColors = []color.RGBA{{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255}, {105, 194, 119, 255}, {177, 94, 218, 255}}

type point struct{ x, y int }
type game struct {
	board          [rows][cols]int
	marked         map[point]bool
	rng            *rand.Rand
	round, total   int
	scanned, clear bool
	message        string
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(5702)), marked: map[point]bool{}, message: "Predict the lines, then press SCAN."}
	g.makeBoard()
	return g
}
func (g *game) makeBoard() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = g.rng.Intn(len(pieceColors))
		}
	}
	kind := g.round % len(pieceColors)
	row := 1 + g.round%5
	for x := 1; x < 4+(g.round%2); x++ {
		g.board[row][x] = kind
	}
	col := 4 - g.round%3
	for y := 3; y < 6; y++ {
		g.board[y][col] = (kind + 2) % len(pieceColors)
	}
	g.marked = map[point]bool{}
	g.scanned = false
}
func (g *game) Update() error {
	if g.clear {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	activate := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if _, y, ok := press(); ok && y > 610 {
		activate = true
	}
	if activate {
		if !g.scanned {
			count := g.scan()
			g.total += count
			g.scanned = true
			g.message = fmt.Sprintf("Found %d unique matched cells. Press NEXT.", count)
		} else {
			g.round++
			if g.round >= 4 {
				g.clear = true
			} else {
				g.makeBoard()
				g.message = "New board: predict first, then scan."
			}
		}
	}
	return nil
}
func (g *game) scan() int {
	g.marked = map[point]bool{}
	for y := 0; y < rows; y++ {
		start := 0
		for x := 1; x <= cols; x++ {
			if x == cols || g.board[y][x] != g.board[y][start] {
				if x-start >= 3 {
					for i := start; i < x; i++ {
						g.marked[point{i, y}] = true
					}
				}
				start = x
			}
		}
	}
	for x := 0; x < cols; x++ {
		start := 0
		for y := 1; y <= rows; y++ {
			if y == rows || g.board[y][x] != g.board[start][x] {
				if y-start >= 3 {
					for i := start; i < y; i++ {
						g.marked[point{x, i}] = true
					}
				}
				start = y
			}
		}
	}
	return len(g.marked)
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(s, "MATCH SCANNER", 185, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BOARD %d/4   MATCHED CELLS TOTAL %02d", g.round+1, g.total), 115, 70)
	ebitenutil.DebugPrintAt(s, g.message, 55, 100)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.DrawFilledRect(s, px+3, py+3, cell-6, cell-6, pieceColors[g.board[y][x]], false)
			vector.DrawFilledCircle(s, px+32, py+32, 17, color.RGBA{255, 255, 255, 45}, false)
			if g.marked[point{x, y}] {
				vector.StrokeRect(s, px+2, py+2, cell-4, cell-4, 6, color.White, false)
			}
		}
	}
	vector.DrawFilledRect(s, 55, 615, 370, 65, color.RGBA{240, 177, 65, 255}, false)
	label := "SCAN ROWS + COLUMNS [SPACE]"
	if g.scanned {
		label = "NEXT BOARD [SPACE]"
	}
	ebitenutil.DebugPrintAt(s, label, 130, 642)
	if g.clear {
		overlay(s, "FOUR BOARDS SCANNED!\n\nTAP / SPACE TO RESTART")
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
	ebitenutil.DebugPrintAt(s, m, 110, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Match Scanner — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
