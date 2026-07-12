package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const (
	screenWidth  = 480
	screenHeight = 720
	cols         = 6
	rows         = 7
	cell         = 64
	boardX       = 48
	boardY       = 150
	empty        = -1
)

var pieceColors = []color.RGBA{
	{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255},
	{105, 194, 119, 255}, {177, 94, 218, 255},
}

type point struct{ x, y int }

// stage is data, not a rule. A designer can make another level by changing
// these values without rewriting Update or the match scanner.
type stage struct {
	name        string
	moves       int
	targetScore int
	seed        int64
	board       [rows][cols]int
}

var firstStage = stage{
	name:        "CORAL COVE",
	moves:       12,
	targetScore: 650,
	seed:        6106,
	board: [rows][cols]int{
		{0, 1, 2, 3, 4, 0},
		{1, 2, 0, 4, 0, 1},
		{2, 0, 3, 0, 1, 2},
		{3, 2, 0, 1, 2, 3},
		{4, 0, 1, 2, 3, 4},
		{0, 1, 2, 3, 4, 0},
		{1, 2, 3, 4, 0, 1},
	},
}

type faller struct {
	kind          int
	x, fromY, toY int
}

type game struct {
	level            stage
	board            [rows][cols]int
	rng              *rand.Rand
	cursor, selected point
	hasSelection     bool
	moves, score     int
	combo            int
	message          string
	won, lost        bool

	// Juice: flash matched cells, then lerp falls. Input is ignored while busy.
	busy      bool
	flash     map[point]bool
	flashLeft int
	falling   bool
	progress  float64
	fallers   []faller
	pending   bool // continue cascade after anim
}

func newGame(level stage) *game {
	g := &game{level: level, rng: rand.New(rand.NewSource(level.seed))}
	g.board = level.board
	g.moves = level.moves
	g.cursor = point{2, 3}
	g.message = "Swap neighbors. Match 3 or more!"
	return g
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame(g.level)
		}
		return nil
	}

	if g.busy {
		if g.flashLeft > 0 {
			g.flashLeft--
			if g.flashLeft == 0 {
				g.clear(g.flash)
				g.flash = nil
				g.beginFall()
			}
			return nil
		}
		if g.falling {
			g.progress += 0.1
			if g.progress >= 1 {
				g.falling = false
				g.progress = 0
				g.fallers = nil
				g.refillEmpties()
				g.busy = false
				if g.pending {
					g.resolveMatches()
				}
			}
			return nil
		}
		g.busy = false
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.cursor.x > 0 {
		g.cursor.x--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.cursor.x < cols-1 {
		g.cursor.x++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && g.cursor.y > 0 {
		g.cursor.y--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && g.cursor.y < rows-1 {
		g.cursor.y++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.choose(g.cursor)
	}
	if p, ok := pressedCell(); ok {
		g.cursor = p
		g.choose(p)
	}
	return nil
}

func (g *game) choose(p point) {
	if !g.hasSelection {
		g.selected = p
		g.hasSelection = true
		g.message = "Now choose a neighbor."
		return
	}
	if p == g.selected {
		g.hasSelection = false
		g.message = "Selection canceled."
		return
	}
	if abs(p.x-g.selected.x)+abs(p.y-g.selected.y) != 1 {
		g.selected = p
		g.message = "Choose a touching neighbor."
		return
	}

	a := g.selected
	g.hasSelection = false
	g.board[a.y][a.x], g.board[p.y][p.x] = g.board[p.y][p.x], g.board[a.y][a.x]
	if len(scan(g.board)) == 0 {
		g.board[a.y][a.x], g.board[p.y][p.x] = g.board[p.y][p.x], g.board[a.y][a.x]
		g.message = "No line—swap returned. Try again!"
		return
	}

	g.moves--
	g.combo = 0
	g.resolveMatches()
}

func (g *game) resolveMatches() {
	matches := scan(g.board)
	if len(matches) == 0 {
		g.pending = false
		if !hasValidSwap(g.board) {
			g.makePlayableBoard()
		}
		if g.score >= g.level.targetScore {
			g.won = true
		} else if g.moves == 0 {
			g.lost = true
		}
		return
	}
	g.combo++
	g.score += len(matches) * 10 * g.combo
	g.message = fmt.Sprintf("%d pieces! Chain x%d", len(matches), g.combo)
	g.flash = matches
	g.flashLeft = 10
	g.busy = true
	g.pending = true
}

func (g *game) beginFall() {
	old := g.board
	g.fallOnly()
	g.fallers = nil
	for x := 0; x < cols; x++ {
		src := []int{}
		for y := 0; y < rows; y++ {
			if old[y][x] != empty {
				src = append(src, y)
			}
		}
		si := 0
		for y := 0; y < rows; y++ {
			if g.board[y][x] == empty {
				continue
			}
			from := src[si]
			si++
			if from != y {
				g.fallers = append(g.fallers, faller{kind: g.board[y][x], x: x, fromY: from, toY: y})
			}
		}
	}
	if len(g.fallers) == 0 {
		g.refillEmpties()
		g.busy = false
		if g.pending {
			g.resolveMatches()
		}
		return
	}
	g.falling = true
	g.progress = 0
}

func (g *game) fallOnly() {
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

func (g *game) refillEmpties() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if g.board[y][x] == empty {
				g.board[y][x] = g.rng.Intn(len(pieceColors))
			}
		}
	}
}

func scan(board [rows][cols]int) map[point]bool {
	found := map[point]bool{}
	for y := 0; y < rows; y++ {
		start := 0
		for x := 1; x <= cols; x++ {
			if x == cols || board[y][x] != board[y][start] {
				if board[y][start] != empty && x-start >= 3 {
					for i := start; i < x; i++ {
						found[point{i, y}] = true
					}
				}
				start = x
			}
		}
	}
	for x := 0; x < cols; x++ {
		start := 0
		for y := 1; y <= rows; y++ {
			if y == rows || board[y][x] != board[start][x] {
				if board[start][x] != empty && y-start >= 3 {
					for i := start; i < y; i++ {
						found[point{x, i}] = true
					}
				}
				start = y
			}
		}
	}
	return found
}

func (g *game) clear(matches map[point]bool) {
	for p := range matches {
		g.board[p.y][p.x] = empty
	}
}

func (g *game) fallAndRefill() {
	g.fallOnly()
	g.refillEmpties()
}

// makePlayableBoard is the safety net used by a complete level. Random refills
// can occasionally leave no useful swap, so generate a settled, playable board.
func (g *game) makePlayableBoard() {
	for {
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				for {
					kind := g.rng.Intn(len(pieceColors))
					if x >= 2 && g.board[y][x-1] == kind && g.board[y][x-2] == kind {
						continue
					}
					if y >= 2 && g.board[y-1][x] == kind && g.board[y-2][x] == kind {
						continue
					}
					g.board[y][x] = kind
					break
				}
			}
		}
		if hasValidSwap(g.board) {
			return
		}
	}
}

func hasValidSwap(board [rows][cols]int) bool {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			for _, d := range [...]point{{1, 0}, {0, 1}} {
				nx, ny := x+d.x, y+d.y
				if nx >= cols || ny >= rows {
					continue
				}
				board[y][x], board[ny][nx] = board[ny][nx], board[y][x]
				if len(scan(board)) > 0 {
					return true
				}
				board[y][x], board[ny][nx] = board[ny][nx], board[y][x]
			}
		}
	}
	return false
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 25, 45, 255})
	ebitenutil.DebugPrintAt(screen, "EBI MATCH / "+g.level.name, 145, 28)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MOVES %02d     SCORE %04d / %04d", g.moves, g.score, g.level.targetScore), 105, 67)
	barWidth := float32(360) * float32(g.score) / float32(g.level.targetScore)
	if barWidth > 360 {
		barWidth = 360
	}
	vector.DrawFilledRect(screen, 60, 94, 360, 14, color.RGBA{45, 61, 86, 255}, false)
	vector.DrawFilledRect(screen, 60, 94, barWidth, 14, color.RGBA{245, 190, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, g.message, 112, 122)

	animating := map[point]bool{}
	if g.falling {
		for _, f := range g.fallers {
			animating[point{f.x, f.toY}] = true
			py := float32(boardY) + float32(float64(f.fromY)+float64(f.toY-f.fromY)*g.progress)*cell
			px := float32(boardX + f.x*cell)
			trackatlas.Draw(screen, trackatlas.Gem(f.kind), float64(px+3), float64(py+3), float64(cell-6))
		}
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px := float32(boardX + x*cell)
			py := float32(boardY + y*cell)
			if g.board[y][x] == empty {
				vector.StrokeRect(screen, px+8, py+8, cell-16, cell-16, 2, color.RGBA{60, 75, 100, 255}, false)
				continue
			}
			if animating[point{x, y}] {
				continue
			}
			if g.flashLeft > 0 && g.flash[point{x, y}] {
				trackatlas.DrawTinted(screen, trackatlas.Gem(g.board[y][x]), float64(px+cell/2), float64(py+cell/2), float64(cell-6), 2.4, 2.4, 2.4, 1)
			} else {
				trackatlas.Draw(screen, trackatlas.Gem(g.board[y][x]), float64(px+3), float64(py+3), float64(cell-6))
			}
			if g.hasSelection && g.selected == (point{x, y}) {
				vector.StrokeRect(screen, px+2, py+2, cell-4, cell-4, 6, color.White, false)
			}
			if !g.busy && g.cursor == (point{x, y}) {
				vector.StrokeRect(screen, px+7, py+7, cell-14, cell-14, 3, color.RGBA{25, 32, 48, 255}, false)
			}
		}
	}
	ebitenutil.DebugPrintAt(screen, "Tap two neighbors  |  Arrows + Space", 90, 628)
	ebitenutil.DebugPrintAt(screen, "Reach the gold score before moves run out", 74, 654)
	if g.won {
		overlay(screen, "STAGE CLEAR!\n\nTAP / SPACE TO PLAY AGAIN")
	}
	if g.lost {
		overlay(screen, "OUT OF MOVES\n\nTAP / SPACE TO RETRY")
	}
}

func pressedCell() (point, bool) {
	x, y, ok := pointerPress()
	if !ok || x < boardX || y < boardY || x >= boardX+cols*cell || y >= boardY+rows*cell {
		return point{}, false
	}
	return point{(x - boardX) / cell, (y - boardY) / cell}, true
}

func pointerPress() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	touches := inpututil.AppendJustPressedTouchIDs(nil)
	if len(touches) > 0 {
		x, y := ebiten.TouchPosition(touches[0])
		return x, y, true
	}
	return 0, 0, false
}

func retryPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) {
		return true
	}
	_, _, ok := pointerPress()
	return ok
}

func overlay(screen *ebiten.Image, message string) {
	vector.DrawFilledRect(screen, 50, 278, 380, 160, color.RGBA{5, 16, 34, 245}, false)
	vector.StrokeRect(screen, 50, 278, 380, 160, 4, color.RGBA{245, 190, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, message, 115, 330)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (g *game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebi Match — Ebitengine")
	if err := ebiten.RunGame(newGame(firstStage)); err != nil {
		panic(err)
	}
}
