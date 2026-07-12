package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width, height = 480, 720
	cols, rows    = 6, 9
	cell          = 48
	boardX        = 96
	boardY        = 112
	empty         = -1
	maxMisses     = 3
)

var pieceColors = []color.RGBA{{232, 86, 87, 255}, {72, 158, 219, 255}, {239, 183, 59, 255}, {94, 188, 108, 255}}
var colorNames = []string{"RED", "BLUE", "YELLOW", "GREEN"}

type point struct{ x, y int }

var directions = []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var childOffsets = []point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type game struct {
	board                   [rows][cols]int
	x, y, rotation, timer   int
	pivotKind, childKind    int
	mission, groups, misses int
	visitedOrder            []point
	marked                  map[point]bool
	scanIndex               int
	scanning, waiting       bool
	clear, over             bool
	message                 string
}

func newGame() *game {
	g := &game{}
	g.loadMission(0)
	return g
}

func (g *game) loadMission(index int) {
	g.mission = index
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = empty
		}
	}
	switch index {
	case 0:
		g.board[8][0], g.board[8][1], g.board[7][0] = 0, 0, 0
		g.pivotKind, g.childKind = 0, 1
	case 1:
		g.board[8][4], g.board[8][5], g.board[7][5] = 1, 1, 1
		g.pivotKind, g.childKind = 1, 2
	case 2:
		g.board[8][2], g.board[7][2], g.board[6][2] = 3, 3, 3
		g.pivotKind, g.childKind = 3, 0
	}
	g.visitedOrder = nil
	g.marked = map[point]bool{}
	g.scanning, g.waiting = false, false
	g.spawn()
	g.message = fmt.Sprintf("Mission %d: place %s beside the outlined group of three.", index+1, colorNames[g.pivotKind])
}

func (g *game) spawn() {
	g.x, g.y, g.rotation, g.timer = cols/2, 1, 0, 0
	if g.blocked(g.x, g.y, g.rotation) {
		g.over = true
		g.message = "No room for a new pair. The board topped out."
	}
}

func (g *game) cellsAt(x, y, rotation int) [2]point {
	offset := childOffsets[rotation%4]
	return [2]point{{x, y}, {x + offset.x, y + offset.y}}
}

func (g *game) blocked(x, y, rotation int) bool {
	for _, p := range g.cellsAt(x, y, rotation) {
		if p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || g.board[p.y][p.x] != empty {
			return true
		}
	}
	return false
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	if g.scanning {
		if g.scanIndex < len(g.visitedOrder) {
			g.timer++
			if g.timer >= 5 {
				g.scanIndex++
				g.timer = 0
			}
			g.message = fmt.Sprintf("BFS queue processed %d/%d colored cells.", g.scanIndex, len(g.visitedOrder))
			return nil
		}
		g.scanning = false
		if len(g.marked) >= 4 {
			g.groups++
			g.waiting = true
			g.message = fmt.Sprintf("Connected group size %d >= 4! Press NEXT.", len(g.marked))
		} else {
			g.misses++
			if g.misses >= maxMisses {
				g.over = true
				g.message = "Three pairs missed every group of four."
			} else {
				g.message = fmt.Sprintf("No group reached four. Misses %d/%d.", g.misses, maxMisses)
				g.spawn()
			}
		}
		return nil
	}
	if g.waiting {
		if anyActionPressed() {
			if g.mission == 2 {
				g.clear = true
				g.message = "Three groups found by four-way connected search!"
			} else {
				g.loadMission(g.mission + 1)
			}
		}
		return nil
	}

	left := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	turn := inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	drop := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if x, y, ok := pressPosition(); ok && y >= 590 {
		switch {
		case x < 120:
			left = true
		case x < 240:
			turn = true
		case x < 360:
			right = true
		default:
			drop = true
		}
	}
	if left && !g.blocked(g.x-1, g.y, g.rotation) {
		g.x--
	}
	if right && !g.blocked(g.x+1, g.y, g.rotation) {
		g.x++
	}
	if turn {
		g.rotate()
	}
	if drop {
		for !g.blocked(g.x, g.y+1, g.rotation) {
			g.y++
		}
		g.lockPair()
		return nil
	}
	g.timer++
	if g.timer >= 90 {
		g.timer = 0
		if !g.blocked(g.x, g.y+1, g.rotation) {
			g.y++
		} else {
			g.lockPair()
		}
	}
	return nil
}

func (g *game) rotate() {
	next := (g.rotation + 1) % 4
	for _, kick := range []int{0, -1, 1} {
		if !g.blocked(g.x+kick, g.y, next) {
			g.x += kick
			g.rotation = next
			return
		}
	}
}

func (g *game) lockPair() {
	cells := g.cellsAt(g.x, g.y, g.rotation)
	g.board[cells[0].y][cells[0].x] = g.pivotKind
	g.board[cells[1].y][cells[1].x] = g.childKind
	g.findGroups()
	g.scanning = true
	g.scanIndex = 0
	g.timer = 0
	g.message = "Pair locked. BFS begins from every unvisited colored cell."
}

func (g *game) findGroups() {
	visited := [rows][cols]bool{}
	g.visitedOrder = nil
	g.marked = map[point]bool{}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if g.board[y][x] == empty || visited[y][x] {
				continue
			}
			kind := g.board[y][x]
			queue := []point{{x, y}}
			visited[y][x] = true
			group := []point{}
			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				group = append(group, current)
				g.visitedOrder = append(g.visitedOrder, current)
				for _, direction := range directions {
					next := point{current.x + direction.x, current.y + direction.y}
					if next.x < 0 || next.x >= cols || next.y < 0 || next.y >= rows || visited[next.y][next.x] || g.board[next.y][next.x] != kind {
						continue
					}
					visited[next.y][next.x] = true
					queue = append(queue, next)
				}
			}
			if len(group) >= 4 {
				for _, p := range group {
					g.marked[p] = true
				}
			}
		}
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 20, 38, 255})
	ebitenutil.DebugPrintAt(screen, "FOUR-WAY COLOR SEARCH", 165, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MISSION %d/3   GROUPS %d/3   MISSES %d/%d", g.mission+1, g.groups, g.misses, maxMisses), 102, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 55, 72)
	vector.DrawFilledRect(screen, boardX-4, boardY-4, cols*cell+8, rows*cell+8, color.RGBA{30, 45, 67, 255}, false)
	visitedNow := map[point]int{}
	for i := 0; i < min(g.scanIndex, len(g.visitedOrder)); i++ {
		visitedNow[g.visitedOrder[i]] = i + 1
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{62, 78, 102, 255}, false)
			if g.board[y][x] != empty {
				drawPiece(screen, x, y, pieceColors[g.board[y][x]], "")
			}
			p := point{x, y}
			if order, ok := visitedNow[p]; ok {
				vector.StrokeRect(screen, px+4, py+4, cell-8, cell-8, 3, color.RGBA{247, 211, 84, 255}, false)
				ebitenutil.DebugPrintAt(screen, fmt.Sprint(order), int(px)+19, int(py)+19)
			}
			if g.marked[p] && !g.scanning {
				vector.StrokeRect(screen, px+2, py+2, cell-4, cell-4, 6, color.White, false)
			}
		}
	}
	if !g.scanning && !g.waiting && !g.clear && !g.over {
		cells := g.cellsAt(g.x, g.y, g.rotation)
		drawPiece(screen, cells[0].x, cells[0].y, pieceColors[g.pivotKind], "P")
		drawPiece(screen, cells[1].x, cells[1].y, pieceColors[g.childKind], "C")
	}
	labels := []string{"LEFT", "TURN", "RIGHT", "DROP / NEXT"}
	for i, label := range labels {
		vector.DrawFilledRect(screen, float32(i*120+5), 590, 110, 72, color.RGBA{51, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*120+24, 622)
	}
	ebitenutil.DebugPrintAt(screen, "Arrows / A,D,W / Space / mouse / touch", 102, 684)
	if g.clear {
		overlay(screen, "THREE GROUPS FOUND!\n\nTAP / ENTER TO RETRY")
	} else if g.over {
		overlay(screen, "COLOR SEARCH FAILED\n\nTAP / ENTER TO RETRY")
	}
}

func drawPiece(screen *ebiten.Image, x, y int, c color.RGBA, label string) {
	cx := float32(boardX + x*cell + cell/2)
	cy := float32(boardY + y*cell + cell/2)
	vector.DrawFilledCircle(screen, cx, cy, cell/2-5, c, false)
	vector.StrokeCircle(screen, cx, cy, cell/2-5, 2, color.White, false)
	if label != "" {
		ebitenutil.DebugPrintAt(screen, label, int(cx)-3, int(cy)-5)
	}
}

func anyActionPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return true
	}
	_, _, ok := pressPosition()
	return ok
}

func pressPosition() (int, int, bool) {
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

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 108, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Four-way Color Search — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
