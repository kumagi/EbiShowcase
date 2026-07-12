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

const (
	screenW     = 480
	screenH     = 720
	cols        = 10
	rows        = 18
	cell        = 26
	boardX      = 110
	boardY      = 92
	targetScore = 1200
	noPiece     = -1
)

type point struct{ x, y int }

type piece struct {
	kind, rotation int
	pos            point
}

var shapeNames = [...]string{"I", "O", "T", "L", "J", "S", "Z"}
var shapeColors = [...]color.RGBA{
	{65, 190, 207, 255}, {239, 190, 62, 255}, {174, 97, 205, 255},
	{232, 142, 56, 255}, {71, 118, 209, 255}, {90, 181, 101, 255}, {224, 82, 86, 255},
}
var bases = [7][4]point{
	{{-1, 0}, {0, 0}, {1, 0}, {2, 0}},
	{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	{{-1, 0}, {0, 0}, {1, 0}, {0, 1}},
	{{-1, 0}, {0, 0}, {1, 0}, {1, 1}},
	{{-1, 0}, {0, 0}, {1, 0}, {-1, 1}},
	{{0, 0}, {1, 0}, {-1, 1}, {0, 1}},
	{{-1, 0}, {0, 0}, {0, 1}, {1, 1}},
}
var kickTests = [...]point{{0, 0}, {-1, 0}, {1, 0}, {-2, 0}, {2, 0}, {0, -1}}

type game struct {
	board                     [rows][cols]int
	active                    piece
	rng                       *rand.Rand
	bag, queue                []int
	hold                      int
	canHold                   bool
	tick, score, lines, combo int
	level                     int
	message                   string
	won, lost                 bool
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(6606)), hold: noPiece, level: 1, combo: -1}
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = noPiece
		}
	}
	g.fillQueue()
	g.spawn(g.takeNext())
	g.message = "Clear lines and score 1200!"
	return g
}

func (g *game) refillBag() {
	g.bag = g.rng.Perm(len(bases))
}

func (g *game) fillQueue() {
	for len(g.queue) < 5 {
		if len(g.bag) == 0 {
			g.refillBag()
		}
		g.queue = append(g.queue, g.bag[0])
		g.bag = g.bag[1:]
	}
}

func (g *game) takeNext() int {
	g.fillQueue()
	kind := g.queue[0]
	g.queue = g.queue[1:]
	g.fillQueue()
	return kind
}

func (g *game) spawn(kind int) {
	g.active = piece{kind: kind, pos: point{cols / 2, 0}}
	g.tick = 0
	g.canHold = true
	if g.collides(g.active) {
		g.lost = true
		g.message = "No room for the next piece."
	}
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}

	left := inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	rotate := inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	soft := inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)
	hard := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	hold := inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyShift)
	if x, y, ok := pointerPress(); ok && y >= 620 {
		switch x / 80 {
		case 0:
			left = true
		case 1:
			right = true
		case 2:
			rotate = true
		case 3:
			soft = true
		case 4:
			hard = true
		case 5:
			hold = true
		}
	}

	if left {
		g.move(-1, 0)
	}
	if right {
		g.move(1, 0)
	}
	if rotate {
		g.rotate()
	}
	if hold {
		g.swapHold()
	}
	if hard {
		distance := 0
		for g.move(0, 1) {
			distance++
		}
		g.score += distance * 2
		g.lock()
		return nil
	}
	if soft {
		if g.move(0, 1) {
			g.score++
		} else {
			g.lock()
		}
		g.tick = 0
		return nil
	}

	g.tick++
	fallEvery := max(10, 42-(g.level-1)*5)
	if g.tick >= fallEvery {
		g.tick = 0
		if !g.move(0, 1) {
			g.lock()
		}
	}
	return nil
}

func (g *game) blocks(p piece) [4]point {
	result := bases[p.kind]
	if p.kind != 1 { // the square looks identical in every rotation
		for turn := 0; turn < p.rotation; turn++ {
			for i, b := range result {
				result[i] = point{-b.y, b.x}
			}
		}
	}
	for i := range result {
		result[i].x += p.pos.x
		result[i].y += p.pos.y
	}
	return result
}

func (g *game) collides(p piece) bool {
	for _, b := range g.blocks(p) {
		if b.x < 0 || b.x >= cols || b.y < 0 || b.y >= rows || g.board[b.y][b.x] != noPiece {
			return true
		}
	}
	return false
}

func (g *game) move(dx, dy int) bool {
	next := g.active
	next.pos.x += dx
	next.pos.y += dy
	if g.collides(next) {
		return false
	}
	g.active = next
	return true
}

func (g *game) rotate() {
	next := g.active
	next.rotation = (next.rotation + 1) % 4
	for _, kick := range kickTests {
		candidate := next
		candidate.pos.x += kick.x
		candidate.pos.y += kick.y
		if !g.collides(candidate) {
			g.active = candidate
			g.message = fmt.Sprintf("Rotated with kick (%+d,%+d).", kick.x, kick.y)
			return
		}
	}
	g.message = "Rotation blocked at every kick position."
}

func (g *game) swapHold() {
	if !g.canHold {
		g.message = "HOLD can be used once per falling piece."
		return
	}
	old := g.active.kind
	if g.hold == noPiece {
		g.hold = old
		g.active = piece{kind: g.takeNext(), pos: point{cols / 2, 0}}
	} else {
		old, g.hold = g.hold, old
		g.active = piece{kind: old, pos: point{cols / 2, 0}}
	}
	g.canHold = false
	g.message = "Piece moved through HOLD."
	if g.collides(g.active) {
		g.lost = true
	}
}

func (g *game) ghost() piece {
	ghost := g.active
	for {
		next := ghost
		next.pos.y++
		if g.collides(next) {
			return ghost
		}
		ghost = next
	}
}

func (g *game) lock() {
	for _, b := range g.blocks(g.active) {
		g.board[b.y][b.x] = g.active.kind
	}
	cleared := g.clearLines()
	if cleared > 0 {
		g.combo++
		points := [...]int{0, 100, 300, 500, 800}
		g.score += points[cleared]*g.level + g.combo*50
		g.lines += cleared
		g.level = 1 + g.lines/5
		g.message = fmt.Sprintf("%d line(s)! Combo %d, level %d.", cleared, g.combo+1, g.level)
	} else {
		g.combo = -1
		g.message = "Piece locked. Build a full row."
	}
	if g.score >= targetScore {
		g.won = true
		return
	}
	g.spawn(g.takeNext())
}

func (g *game) clearLines() int {
	write, cleared := rows-1, 0
	for read := rows - 1; read >= 0; read-- {
		full := true
		for x := 0; x < cols; x++ {
			if g.board[read][x] == noPiece {
				full = false
				break
			}
		}
		if full {
			cleared++
			continue
		}
		g.board[write] = g.board[read]
		write--
	}
	for write >= 0 {
		for x := 0; x < cols; x++ {
			g.board[write][x] = noPiece
		}
		write--
	}
	return cleared
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{11, 20, 38, 255})
	ebitenutil.DebugPrintAt(screen, "EBI BLOCKS", 202, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %04d/%04d   LINES %02d   LEVEL %d", g.score, targetScore, g.lines, g.level), 112, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 92, 68)

	vector.DrawFilledRect(screen, boardX-3, boardY-3, cols*cell+6, rows*cell+6, color.RGBA{35, 48, 70, 255}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{57, 72, 97, 255}, false)
			if g.board[y][x] != noPiece {
				drawCell(screen, x, y, shapeColors[g.board[y][x]], 255)
			}
		}
	}
	if !g.won && !g.lost {
		for _, b := range g.blocks(g.ghost()) {
			drawCell(screen, b.x, b.y, shapeColors[g.active.kind], 55)
		}
		for _, b := range g.blocks(g.active) {
			drawCell(screen, b.x, b.y, shapeColors[g.active.kind], 255)
		}
	}

	drawSide(screen, 8, 110, "HOLD", g.hold)
	next := noPiece
	if len(g.queue) > 0 {
		next = g.queue[0]
	}
	drawSide(screen, 382, 110, "NEXT", next)
	ebitenutil.DebugPrintAt(screen, "7-BAG", 399, 205)
	for i := 0; i < min(4, len(g.queue)); i++ {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d  %s", i+1, shapeNames[g.queue[i]]), 404, 229+i*22)
	}

	labels := [...]string{"LEFT", "RIGHT", "TURN", "DOWN", "DROP", "HOLD"}
	for i, label := range labels {
		c := color.RGBA{55, 88, 130, 255}
		if i == 2 || i == 5 {
			c = color.RGBA{194, 116, 67, 255}
		}
		vector.DrawFilledRect(screen, float32(i*80+3), 620, 74, 70, c, false)
		ebitenutil.DebugPrintAt(screen, label, i*80+18, 650)
	}
	if g.won {
		overlay(screen, "EBI BLOCKS CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "STACK REACHED THE TOP\n\nTAP / ENTER TO RETRY")
	}
}

func drawCell(screen *ebiten.Image, x, y int, c color.RGBA, alpha uint8) {
	c.A = alpha
	px, py := float32(boardX+x*cell), float32(boardY+y*cell)
	vector.DrawFilledRect(screen, px+2, py+2, cell-4, cell-4, c, false)
	vector.StrokeRect(screen, px+4, py+4, cell-8, cell-8, 1, color.RGBA{255, 255, 255, alpha / 2}, false)
}

func drawSide(screen *ebiten.Image, x, y int, title string, kind int) {
	vector.DrawFilledRect(screen, float32(x), float32(y), 90, 78, color.RGBA{31, 45, 68, 255}, false)
	ebitenutil.DebugPrintAt(screen, title, x+26, y+12)
	value := "—"
	if kind != noPiece {
		value = shapeNames[kind]
	}
	ebitenutil.DebugPrintAt(screen, value, x+40, y+45)
}

func pointerPress() (int, int, bool) {
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
	_, _, ok := pointerPress()
	return ok
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 160, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 160, 4, color.RGBA{241, 185, 65, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 116, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Blocks — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
