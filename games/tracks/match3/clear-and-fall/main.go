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
const empty = -1

var colors = []color.RGBA{
	{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255},
	{105, 194, 119, 255}, {177, 94, 218, 255},
}

type point struct{ x, y int }

// fallPiece remembers where a gem came from so Draw can lerp while Update
// only advances progress. Logical board cells are already compacted.
type fallPiece struct {
	kind       int
	fromY, toY int
	x          int
}

type game struct {
	board                [rows][cols]int
	marked               map[point]bool
	rng                  *rand.Rand
	phase, cycles, score int
	message              string
	clear                bool

	// phase 1 animation: board is compacted; Draw uses fallers + progress.
	falling  bool
	progress float64
	fallers  []fallPiece
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
	g.falling = false
	g.progress = 0
	g.fallers = nil
	g.message = "STEP 1: clear the outlined match."
}

func (g *game) Update() error {
	if g.clear {
		if restart() {
			*g = *newGame()
		}
		return nil
	}

	// Mid-fall: ignore Space / taps; only advance the tween.
	if g.falling {
		g.progress += 0.07
		if g.progress >= 1 {
			g.falling = false
			g.progress = 0
			g.fallers = nil
			g.phase = 2
			g.message = "STEP 3: refill empty cells from the top."
		}
		return nil
	}

	goNext := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if _, y, ok := press(); ok && y > 610 {
		goNext = true
	}
	if !goNext {
		return nil
	}

	switch g.phase {
	case 0:
		for p := range g.marked {
			g.board[p.y][p.x] = empty
			g.score += 10
		}
		g.phase = 1
		g.message = "STEP 2: watch Update drive progress, Draw lerp."
		g.beginFall()
	case 2:
		g.refill()
		g.cycles++
		if g.cycles >= 3 {
			g.clear = true
		} else {
			g.seed()
		}
	}
	return nil
}

func (g *game) beginFall() {
	// Snapshot old column heights before compacting so Draw can animate.
	old := g.board
	g.fall()
	g.fallers = nil
	for x := 0; x < cols; x++ {
		// Map each surviving kind at new row to its previous row.
		srcRows := []int{}
		for y := 0; y < rows; y++ {
			if old[y][x] != empty {
				srcRows = append(srcRows, y)
			}
		}
		si := 0
		for y := 0; y < rows; y++ {
			if g.board[y][x] == empty {
				continue
			}
			from := srcRows[si]
			si++
			if from != y {
				g.fallers = append(g.fallers, fallPiece{
					kind: g.board[y][x], fromY: from, toY: y, x: x,
				})
			}
		}
	}
	if len(g.fallers) == 0 {
		g.phase = 2
		g.message = "STEP 3: refill empty cells from the top."
		return
	}
	g.falling = true
	g.progress = 0
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
		for write >= 0 {
			g.board[write][x] = empty
			write--
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
				if b[y][start] != empty && x-start >= 3 {
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
				if b[start][x] != empty && y-start >= 3 {
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
	ebitenutil.DebugPrintAt(s, g.message, 40, 100)

	animating := map[point]bool{}
	if g.falling {
		for _, f := range g.fallers {
			animating[point{f.x, f.toY}] = true
			py := float32(oy) + float32(float64(f.fromY)+float64(f.toY-f.fromY)*g.progress)*cell
			px := float32(ox + f.x*cell)
			vector.DrawFilledRect(s, px+3, py+3, cell-6, cell-6, colors[f.kind], false)
		}
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			if g.board[y][x] == empty {
				vector.StrokeRect(s, px+4, py+4, cell-8, cell-8, 2, color.RGBA{70, 83, 105, 255}, false)
				continue
			}
			if animating[point{x, y}] {
				continue // already drawn at lerp height
			}
			vector.DrawFilledRect(s, px+3, py+3, cell-6, cell-6, colors[g.board[y][x]], false)
			if g.phase == 0 && g.marked[point{x, y}] {
				vector.StrokeRect(s, px+2, py+2, cell-4, cell-4, 6, color.White, false)
			}
		}
	}

	label := "CLEAR MATCH [SPACE]"
	switch {
	case g.falling:
		label = fmt.Sprintf("FALLING… %.0f%%  (input ignored)", g.progress*100)
	case g.phase == 2:
		label = "REFILL TOP [SPACE]"
	case g.phase == 1:
		label = "FALL DOWN [SPACE]"
	}
	vector.DrawFilledRect(s, 55, 615, 370, 65, color.RGBA{240, 177, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, label, 110, 642)
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
