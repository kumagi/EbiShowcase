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
	rows    = 10
	cell    = 44
	boardX  = 108
	boardY  = 118
	empty   = -1
)

type point struct{ x, y int }
type phase int

const (
	planning phase = iota
	searching
	clearing
	falling
)

var pieceColors = [...]color.RGBA{{235, 82, 91, 255}, {69, 155, 222, 255}, {241, 183, 54, 255}, {91, 190, 116, 255}}

type game struct {
	board                                [rows][cols]int
	marked                               map[point]bool
	phase                                phase
	timer, round, chain, score, mistakes int
	lastGain                             int
	message                              string
	won, lost                            bool
}

func newGame() *game { g := &game{}; g.seedRound(); return g }

func (g *game) seedRound() {
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = empty
		}
	}
	shift := g.round % 3
	a := g.round % len(pieceColors)
	b := (g.round + 1) % len(pieceColors)
	// Three pieces wait on the floor. The correct plan supplies the fourth.
	for x := 0; x < 3; x++ {
		g.board[rows-1][(x+shift)%cols] = a
	}
	// Three pieces of the next color hover one row up. A fourth rests higher;
	// after the first clear, column compaction joins all four.
	for x := 0; x < 3; x++ {
		g.board[rows-2][(x+shift)%cols] = b
	}
	g.board[rows-3][(3+shift)%cols] = b
	g.phase = planning
	g.chain = 0
	g.lastGain = 0
	g.message = "Predict the plan that completes the floor group."
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	if g.phase == planning {
		choice := -1
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			choice = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			choice = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			choice = 2
		}
		if x, y, ok := pressPosition(); ok && y >= 612 {
			choice = x / 160
		}
		if choice >= 0 {
			g.choose(choice)
		}
		return nil
	}
	g.timer++
	if g.timer < 32 {
		return nil
	}
	g.timer = 0
	switch g.phase {
	case searching:
		g.marked = findGroups(g.board)
		if len(g.marked) == 0 {
			g.finishChain()
			return nil
		}
		g.chain++
		g.lastGain = len(g.marked) * 10 * g.chain
		g.score += g.lastGain
		g.message = fmt.Sprintf("CHAIN %d: %d cells x10 x%d = +%d", g.chain, len(g.marked), g.chain, g.lastGain)
		g.phase = clearing
	case clearing:
		for p := range g.marked {
			g.board[p.y][p.x] = empty
		}
		g.marked = nil
		g.phase = falling
		g.message = "Clear finished. Gravity compacts every column."
	case falling:
		compact(&g.board)
		g.phase = searching
		g.message = "Pieces settled. Search the whole board again."
	}
	return nil
}

func (g *game) choose(choice int) {
	shift := g.round % 3
	target := (3 + shift) % cols
	base := g.round % len(pieceColors)
	columns := [3]int{target, (target + 1) % cols, (target + 2) % cols}
	kinds := [3]int{base, (base + 2) % len(pieceColors), (base + 3) % len(pieceColors)}
	x := columns[choice]
	y := rows - 1
	for y >= 0 && g.board[y][x] != empty {
		y--
	}
	if y < 0 {
		g.mistakes++
		g.message = "That column is full."
		return
	}
	g.board[y][x] = kinds[choice]
	g.phase = searching
	g.timer = 0
	g.chain = 0
	g.message = "Plan placed. SEARCH begins."
}

func (g *game) finishChain() {
	if g.chain >= 2 {
		g.round++
		if g.round >= 3 {
			g.won = true
			g.message = "Three chain plans solved!"
			return
		}
		g.seedRound()
		return
	}
	g.mistakes++
	if g.mistakes >= 3 {
		g.lost = true
		g.message = "Three plans ended before a 2-chain."
		return
	}
	g.seedRound()
	g.message = fmt.Sprintf("Only %d chain. Try the next board. Mistakes %d/3", g.chain, g.mistakes)
}

func findGroups(board [rows][cols]int) map[point]bool {
	visited := map[point]bool{}
	marked := map[point]bool{}
	dirs := [...]point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			s := point{x, y}
			if board[y][x] == empty || visited[s] {
				continue
			}
			kind := board[y][x]
			q := []point{s}
			visited[s] = true
			group := []point{}
			for len(q) > 0 {
				p := q[0]
				q = q[1:]
				group = append(group, p)
				for _, d := range dirs {
					n := point{p.x + d.x, p.y + d.y}
					if n.x < 0 || n.x >= cols || n.y < 0 || n.y >= rows || visited[n] || board[n.y][n.x] != kind {
						continue
					}
					visited[n] = true
					q = append(q, n)
				}
			}
			if len(group) >= 4 {
				for _, p := range group {
					marked[p] = true
				}
			}
		}
	}
	return marked
}

func compact(board *[rows][cols]int) {
	for x := 0; x < cols; x++ {
		w := rows - 1
		for r := rows - 1; r >= 0; r-- {
			if board[r][x] == empty {
				continue
			}
			board[w][x] = board[r][x]
			if w != r {
				board[r][x] = empty
			}
			w--
		}
		for w >= 0 {
			board[w][x] = empty
			w--
		}
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 18, 35, 255})
	names := [...]string{"PLAN", "SEARCH", "CLEAR", "GRAVITY"}
	ebitenutil.DebugPrintAt(screen, "CHAIN SCORE LAB", 186, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BOARD %d/3  PHASE %s  CHAIN %d  SCORE %d", g.round+1, names[g.phase], g.chain, g.score), 72, 46)
	ebitenutil.DebugPrintAt(screen, g.message, 42, 76)
	vector.DrawFilledRect(screen, boardX-4, boardY-4, cols*cell+8, rows*cell+8, color.RGBA{30, 45, 67, 255}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px := float32(boardX + x*cell)
			py := float32(boardY + y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{59, 77, 101, 255}, false)
			k := g.board[y][x]
			if k == empty {
				continue
			}
			vector.DrawFilledCircle(screen, px+cell/2, py+cell/2, cell/2-4, pieceColors[k], false)
			if g.marked[point{x, y}] {
				vector.StrokeCircle(screen, px+cell/2, py+cell/2, cell/2-2, 4, color.White, false)
			}
		}
	}
	labels := [...]string{"[1] COMPLETE", "[2] OFFSET", "[3] MIX"}
	for i, l := range labels {
		c := color.RGBA{48, 78, 116, 255}
		if i == 0 {
			c = color.RGBA{174, 104, 55, 255}
		}
		vector.DrawFilledRect(screen, float32(i*160+5), 612, 150, 58, c, false)
		ebitenutil.DebugPrintAt(screen, l, i*160+23, 635)
	}
	ebitenutil.DebugPrintAt(screen, "After every fall: SEARCH again. Multiplier = chain number.", 55, 690)
	if g.won {
		overlay(screen, "THREE 2-CHAINS BUILT!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "CHAIN PLAN FAILED\n\nTAP / ENTER TO RETRY")
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
	vector.DrawFilledRect(screen, 35, 278, 410, 142, color.RGBA{5, 12, 25, 235}, false)
	ebitenutil.DebugPrintAt(screen, text, 132, 320)
}
func (g *game) Layout(int, int) (int, int) { return screenW, screenH }
func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Chain Score Lab")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
