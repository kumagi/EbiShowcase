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
	cols    = 10
	rows    = 14
	cell    = 32
	boardX  = 80
	boardY  = 120
)

type point struct{ x, y int }

var shapes = [][4]point{
	{{-1, 0}, {0, 0}, {1, 0}, {0, 1}}, // T
	{{-1, 0}, {0, 0}, {1, 0}, {1, 1}}, // L
}

// Each offset is tried in order after a rotation collides. This short,
// teachable kick table is intentionally simpler than a commercial ruleset.
var kickTests = [...]point{{0, 0}, {-1, 0}, {1, 0}, {-2, 0}, {2, 0}, {0, -1}}

type game struct {
	board                         [rows][cols]bool
	pos                           point
	rotation, shape, tick, locked int
	rotations, kicks              int
	lastTests, message            string
	won, lost                     bool
}

func newGame() *game {
	g := &game{}
	g.resetBoard()
	g.spawn()
	g.message = "Move beside a wall, then ROTATE."
	g.lastTests = "KICK TESTS: waiting"
	return g
}

func (g *game) resetBoard() {
	// A small staircase makes upward and sideways kicks visible.
	for x := 0; x < cols; x++ {
		g.board[rows-1][x] = true
	}
	for _, p := range []point{{0, 12}, {1, 12}, {8, 12}, {9, 12}, {0, 11}, {9, 11}, {4, 12}, {5, 12}} {
		g.board[p.y][p.x] = true
	}
}

func (g *game) spawn() {
	g.shape = g.locked % len(shapes)
	g.rotation = 0
	g.pos = point{cols / 2, 0}
	g.tick = 0
	if g.collides(g.pos, g.rotation) {
		g.lost = true
	}
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}

	left := inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	right := inpututil.IsKeyJustPressed(ebiten.KeyRight)
	rotate := inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
	drop := inpututil.IsKeyJustPressed(ebiten.KeyDown)
	if x, y, ok := pointerPress(); ok && y >= 610 {
		switch {
		case x < 120:
			left = true
		case x < 240:
			right = true
		case x < 360:
			rotate = true
		default:
			drop = true
		}
	}
	if left {
		g.move(-1, 0)
	}
	if right {
		g.move(1, 0)
	}
	if rotate {
		g.tryRotate()
	}
	if drop {
		for g.move(0, 1) {
		}
		g.lock()
		return nil
	}

	g.tick++
	if g.tick >= 38 {
		g.tick = 0
		if !g.move(0, 1) {
			g.lock()
		}
	}
	return nil
}

func (g *game) move(dx, dy int) bool {
	next := point{g.pos.x + dx, g.pos.y + dy}
	if g.collides(next, g.rotation) {
		return false
	}
	g.pos = next
	return true
}

func (g *game) tryRotate() {
	nextRotation := (g.rotation + 1) % 4
	tried := "KICK TESTS:"
	for i, kick := range kickTests {
		candidate := point{g.pos.x + kick.x, g.pos.y + kick.y}
		if i > 0 {
			tried += " ->"
		}
		tried += fmt.Sprintf(" (%+d,%+d)", kick.x, kick.y)
		if !g.collides(candidate, nextRotation) {
			g.pos = candidate
			g.rotation = nextRotation
			g.rotations++
			if kick != (point{}) {
				g.kicks++
				g.message = fmt.Sprintf("WALL KICK! used (%+d,%+d)", kick.x, kick.y)
			} else {
				g.message = "Rotated at the same position."
			}
			g.lastTests = tried + "  OK"
			g.checkGoal()
			return
		}
	}
	g.lastTests = tried + "  BLOCKED"
	g.message = "Every kick candidate collided."
}

func (g *game) blocks(pos point, rotation int) [4]point {
	result := shapes[g.shape]
	for turn := 0; turn < rotation; turn++ {
		for i, p := range result {
			result[i] = point{-p.y, p.x}
		}
	}
	for i := range result {
		result[i].x += pos.x
		result[i].y += pos.y
	}
	return result
}

func (g *game) collides(pos point, rotation int) bool {
	for _, p := range g.blocks(pos, rotation) {
		if p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || g.board[p.y][p.x] {
			return true
		}
	}
	return false
}

func (g *game) lock() {
	for _, p := range g.blocks(g.pos, g.rotation) {
		if p.y >= 0 && p.y < rows && p.x >= 0 && p.x < cols {
			g.board[p.y][p.x] = true
		}
	}
	g.locked++
	if g.locked >= 8 && !g.won {
		g.lost = true
		return
	}
	g.spawn()
}

func (g *game) checkGoal() {
	if g.rotations >= 8 && g.kicks >= 3 {
		g.won = true
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{14, 24, 43, 255})
	ebitenutil.DebugPrintAt(screen, "ROTATION + WALL KICK LAB", 135, 24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ROTATIONS %d/8   KICKS %d/3   PIECES %d/8", g.rotations, g.kicks, g.locked), 75, 55)
	ebitenutil.DebugPrintAt(screen, g.message, 95, 82)
	ebitenutil.DebugPrintAt(screen, g.lastTests, 22, 103)

	vector.DrawFilledRect(screen, boardX-4, boardY-4, cols*cell+8, rows*cell+8, color.RGBA{35, 48, 70, 255}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{62, 78, 103, 255}, false)
			if g.board[y][x] {
				drawBlock(screen, x, y, color.RGBA{82, 101, 128, 255})
			}
		}
	}
	for _, p := range g.blocks(g.pos, g.rotation) {
		drawBlock(screen, p.x, p.y, color.RGBA{239, 110, 92, 255})
	}

	drawButton(screen, 8, "LEFT", color.RGBA{65, 119, 180, 255})
	drawButton(screen, 126, "RIGHT", color.RGBA{65, 119, 180, 255})
	drawButton(screen, 244, "ROTATE", color.RGBA{214, 139, 57, 255})
	drawButton(screen, 362, "DROP", color.RGBA{91, 163, 112, 255})
	if g.won {
		overlay(screen, "KICK TRAINING CLEAR!\n\nTAP / SPACE TO RETRY")
	}
	if g.lost {
		overlay(screen, "BOARD FILLED UP\n\nTAP / SPACE TO RETRY")
	}
}

func drawBlock(screen *ebiten.Image, x, y int, c color.RGBA) {
	px, py := float32(boardX+x*cell), float32(boardY+y*cell)
	vector.DrawFilledRect(screen, px+2, py+2, cell-4, cell-4, c, false)
	vector.StrokeRect(screen, px+3, py+3, cell-6, cell-6, 2, color.RGBA{255, 255, 255, 90}, false)
}

func drawButton(screen *ebiten.Image, x int, label string, c color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), 612, 110, 66, c, false)
	ebitenutil.DebugPrintAt(screen, label, x+33, 639)
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}
	_, _, ok := pointerPress()
	return ok
}

func overlay(screen *ebiten.Image, message string) {
	vector.DrawFilledRect(screen, 45, 280, 390, 150, color.RGBA{5, 15, 33, 245}, false)
	vector.StrokeRect(screen, 45, 280, 390, 150, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, message, 115, 330)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Rotation and Wall Kicks — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
