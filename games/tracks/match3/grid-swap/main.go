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
const cols, rows = 6, 7
const cell = 64
const ox, oy = 48, 130

var colors = []color.RGBA{{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255}, {105, 194, 119, 255}, {177, 94, 218, 255}}

type point struct{ x, y int }
type game struct {
	board         [rows][cols]int
	selected      point
	hasSelection  bool
	swaps, misses int
	message       string
	clear         bool
}

func newGame() *game {
	g := &game{selected: point{-1, -1}, message: "Select a piece, then one beside it."}
	r := rand.New(rand.NewSource(5601))
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = r.Intn(len(colors))
		}
	}
	return g
}
func (g *game) Update() error {
	if g.clear {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.hasSelection = false
	}
	if x, y, ok := press(); ok {
		cx, cy := (x-ox)/cell, (y-oy)/cell
		if x >= ox && y >= oy && cx >= 0 && cx < cols && cy >= 0 && cy < rows {
			g.choose(point{cx, cy})
		}
	}
	return nil
}
func (g *game) choose(p point) {
	if !g.hasSelection {
		g.selected = p
		g.hasSelection = true
		g.message = fmt.Sprintf("Selected row %d, column %d.", p.y+1, p.x+1)
		return
	}
	d := abs(p.x-g.selected.x) + abs(p.y-g.selected.y)
	if d == 1 {
		a := g.selected
		g.board[a.y][a.x], g.board[p.y][p.x] = g.board[p.y][p.x], g.board[a.y][a.x]
		g.swaps++
		g.message = "Adjacent pieces swapped!"
		if g.swaps >= 10 {
			g.clear = true
		}
	} else {
		g.misses++
		g.message = "Too far apart. Choose a direct neighbor."
	}
	g.selected = p
	g.hasSelection = true
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(s, "GRID SWAP LAB", 185, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("VALID SWAPS %02d/10   NOT ADJACENT %02d", g.swaps, g.misses), 105, 70)
	ebitenutil.DebugPrintAt(s, g.message, 70, 100)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.DrawFilledRect(s, px+3, py+3, cell-6, cell-6, colors[g.board[y][x]], false)
			vector.DrawFilledCircle(s, px+32, py+32, 17, color.RGBA{255, 255, 255, 45}, false)
			if g.hasSelection && g.selected == (point{x, y}) {
				vector.StrokeRect(s, px+1, py+1, cell-2, cell-2, 5, color.White, false)
			}
		}
	}
	ebitenutil.DebugPrintAt(s, "Tap one cell, then a cell sharing an edge.", 90, 610)
	ebitenutil.DebugPrintAt(s, "Diagonal and distant cells are rejected.", 100, 640)
	if g.clear {
		overlay(s, "TEN CLEAN SWAPS!\n\nTAP / SPACE TO RESTART")
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
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func overlay(s *ebiten.Image, m string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, m, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Grid Swap — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
