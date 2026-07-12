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
	screenW = 480
	screenH = 720
	cols    = 6
	rows    = 12
	cell    = 40
	boardX  = 120
	boardY  = 94
	empty   = -1
)

type point struct{ x, y int }

type phase int

const (
	clearPhase phase = iota
	gravityPhase
	settledPhase
)

var colors = [...]color.RGBA{
	{233, 88, 89, 255}, {72, 157, 219, 255}, {239, 183, 58, 255}, {92, 188, 111, 255},
}

type game struct {
	board                  [rows][cols]int
	marked                 map[point]bool
	phase                  phase
	round, cleared, errors int
	groups                 int
	message                string
	won, lost              bool
}

func newGame() *game {
	g := &game{}
	g.seedRound()
	return g
}

func (g *game) seedRound() {
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = empty
		}
	}
	shift := g.round % 2
	a := g.round % len(colors)
	b := (g.round + 1) % len(colors)
	c := (g.round + 2) % len(colors)
	for _, p := range []point{{shift, 10}, {shift, 11}, {shift + 1, 10}, {shift + 1, 11}} {
		g.board[p.y][p.x] = a
	}
	for _, p := range []point{{4 - shift, 9}, {5 - shift, 9}, {4 - shift, 10}, {5 - shift, 10}} {
		g.board[p.y][p.x] = b
	}
	// These pieces float above the groups and make column compaction visible.
	for _, p := range []point{{shift, 5}, {shift, 7}, {shift + 1, 6}, {4 - shift, 4}, {5 - shift, 6}, {3, 8}} {
		g.board[p.y][p.x] = c
	}
	g.marked, g.groups = findClearGroups(g.board)
	g.phase = clearPhase
	g.message = fmt.Sprintf("%d groups are outlined. CLEAR them together.", g.groups)
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	action := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		action = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) || inpututil.IsKeyJustPressed(ebiten.Key2) {
		action = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) || inpututil.IsKeyJustPressed(ebiten.Key3) {
		action = 2
	}
	if x, y, ok := pressPosition(); ok && y >= 610 {
		action = x / 160
	}
	if action >= 0 {
		g.act(action)
	}
	return nil
}

func (g *game) act(action int) {
	expected := int(g.phase)
	if action != expected {
		g.errors++
		g.message = fmt.Sprintf("Wrong order! CLEAR -> GRAVITY -> NEXT. Errors %d/3", g.errors)
		if g.errors >= 3 {
			g.lost = true
		}
		return
	}
	switch g.phase {
	case clearPhase:
		count := len(g.marked)
		for p := range g.marked {
			g.board[p.y][p.x] = empty
		}
		g.cleared += count
		g.marked = nil
		g.phase = gravityPhase
		g.message = fmt.Sprintf("All %d groups cleared at once: %d holes. Now GRAVITY.", g.groups, count)
	case gravityPhase:
		moved := g.compactColumns()
		g.phase = settledPhase
		g.message = fmt.Sprintf("Columns compacted bottom-up: %d pieces fell. Press NEXT.", moved)
	case settledPhase:
		g.round++
		if g.round >= 4 {
			g.won = true
			g.message = "Four boards resolved in the correct two-stage order!"
			return
		}
		g.seedRound()
	}
}

func findClearGroups(board [rows][cols]int) (map[point]bool, int) {
	visited := map[point]bool{}
	marked := map[point]bool{}
	groups := 0
	directions := [...]point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			start := point{x, y}
			if board[y][x] == empty || visited[start] {
				continue
			}
			colorID := board[y][x]
			queue := []point{start}
			visited[start] = true
			group := []point{}
			for len(queue) > 0 {
				p := queue[0]
				queue = queue[1:]
				group = append(group, p)
				for _, d := range directions {
					n := point{p.x + d.x, p.y + d.y}
					if n.x < 0 || n.x >= cols || n.y < 0 || n.y >= rows || visited[n] || board[n.y][n.x] != colorID {
						continue
					}
					visited[n] = true
					queue = append(queue, n)
				}
			}
			if len(group) >= 4 {
				groups++
				for _, p := range group {
					marked[p] = true
				}
			}
		}
	}
	return marked, groups
}

func (g *game) compactColumns() int {
	moved := 0
	for x := 0; x < cols; x++ {
		write := rows - 1
		for read := rows - 1; read >= 0; read-- {
			if g.board[read][x] == empty {
				continue
			}
			if write != read {
				g.board[write][x] = g.board[read][x]
				g.board[read][x] = empty
				moved++
			}
			write--
		}
		for write >= 0 {
			g.board[write][x] = empty
			write--
		}
	}
	return moved
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 20, 38, 255})
	phaseNames := [...]string{"CLEAR", "GRAVITY", "NEXT"}
	ebitenutil.DebugPrintAt(screen, "CLEAR + GRAVITY LAB", 169, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BOARD %d/4   PHASE %s   CLEARED %d   ERRORS %d/3", g.round+1, phaseNames[g.phase], g.cleared, g.errors), 68, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 44, 71)
	vector.DrawFilledRect(screen, boardX-4, boardY-4, cols*cell+8, rows*cell+8, color.RGBA{33, 47, 69, 255}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{62, 78, 102, 255}, false)
			kind := g.board[y][x]
			if kind == empty {
				continue
			}
			vector.DrawFilledCircle(screen, px+cell/2, py+cell/2, cell/2-4, colors[kind], false)
			if g.marked[point{x, y}] {
				vector.StrokeCircle(screen, px+cell/2, py+cell/2, cell/2-2, 4, color.White, false)
			}
		}
	}
	labels := [...]string{"[1/C] CLEAR ALL", "[2/G] GRAVITY", "[3/N] NEXT BOARD"}
	for i, label := range labels {
		c := color.RGBA{52, 84, 122, 255}
		if i == int(g.phase) {
			c = color.RGBA{192, 119, 61, 255}
		}
		vector.DrawFilledRect(screen, float32(i*160+5), 610, 150, 64, c, false)
		ebitenutil.DebugPrintAt(screen, label, i*160+15, 636)
	}
	ebitenutil.DebugPrintAt(screen, "Click / tap or keys C,G,N in that order", 101, 692)
	if g.won {
		overlay(screen, "TWO-STAGE PROCESS CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "THREE ORDER ERRORS\n\nTAP / ENTER TO RETRY")
	}
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	_, _, ok := pressPosition()
	return ok
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 91, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Clear and Gravity Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
