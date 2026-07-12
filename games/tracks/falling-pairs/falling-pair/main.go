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
	boardY  = 100
	empty   = -1
	goal    = 10
)

var pieceColors = [...]color.RGBA{
	{233, 88, 89, 255}, {75, 158, 218, 255}, {239, 184, 59, 255}, {96, 190, 111, 255},
}

type game struct {
	board                [rows][cols]int
	x, y, timer          int // x,y is the pivot; the child is always at y-1
	pivotKind, childKind int
	locked               int
	message              string
	won, lost            bool
}

func newGame() *game {
	g := &game{}
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = empty
		}
	}
	g.spawn()
	g.message = "Move the vertical pair and keep the spawn area open."
	return g
}

func (g *game) spawn() {
	g.x, g.y, g.timer = cols/2, 1, 0
	g.pivotKind = (g.locked * 2) % len(pieceColors)
	g.childKind = (g.locked*2 + 1) % len(pieceColors)
	if g.blocked(g.x, g.y) {
		g.lost = true
		g.message = "The new pair has no room to enter."
	}
}

func (g *game) blocked(x, pivotY int) bool {
	childY := pivotY - 1
	if x < 0 || x >= cols || pivotY >= rows || childY < 0 {
		return true
	}
	return g.board[pivotY][x] != empty || g.board[childY][x] != empty
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
	down := inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)
	drop := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if x, y, ok := pressPosition(); ok && y >= 610 {
		switch {
		case x < 120:
			left = true
		case x < 240:
			down = true
		case x < 360:
			right = true
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
	if drop {
		for g.move(0, 1) {
		}
		g.lockPair()
		return nil
	}
	if down {
		if !g.move(0, 1) {
			g.lockPair()
		}
		g.timer = 0
		return nil
	}
	g.timer++
	if g.timer >= 40 {
		g.timer = 0
		if !g.move(0, 1) {
			g.lockPair()
		}
	}
	return nil
}

func (g *game) move(dx, dy int) bool {
	if g.blocked(g.x+dx, g.y+dy) {
		return false
	}
	g.x += dx
	g.y += dy
	return true
}

func (g *game) lockPair() {
	g.board[g.y][g.x] = g.pivotKind
	g.board[g.y-1][g.x] = g.childKind
	g.locked++
	g.message = fmt.Sprintf("Pair fixed into two board cells. %d/%d", g.locked, goal)
	if g.locked >= goal {
		g.won = true
		return
	}
	g.spawn()
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 20, 38, 255})
	ebitenutil.DebugPrintAt(screen, "FALLING PAIR WORKSHOP", 164, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PAIRS LOCKED %02d/%02d", g.locked, goal), 174, 48)
	ebitenutil.DebugPrintAt(screen, g.message, 68, 74)
	vector.DrawFilledRect(screen, boardX-4, boardY-4, cols*cell+8, rows*cell+8, color.RGBA{33, 47, 69, 255}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{62, 78, 102, 255}, false)
			if g.board[y][x] != empty {
				drawPiece(screen, x, y, pieceColors[g.board[y][x]], "")
			}
		}
	}
	if !g.won && !g.lost {
		drawPiece(screen, g.x, g.y, pieceColors[g.pivotKind], "P")
		drawPiece(screen, g.x, g.y-1, pieceColors[g.childKind], "C")
		vector.StrokeLine(screen, float32(boardX+g.x*cell+cell/2), float32(boardY+(g.y-1)*cell+cell/2), float32(boardX+g.x*cell+cell/2), float32(boardY+g.y*cell+cell/2), 4, color.White, false)
	}
	ebitenutil.DebugPrintAt(screen, "C = CHILD (pivot + 0,-1)   P = PIVOT", 98, 590)
	buttons := [...]string{"LEFT", "DOWN", "RIGHT", "DROP"}
	for i, label := range buttons {
		vector.DrawFilledRect(screen, float32(i*120+5), 610, 110, 62, color.RGBA{52, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*120+39, 636)
	}
	ebitenutil.DebugPrintAt(screen, "Arrows / A,D,S / Space / click / tap", 104, 690)
	if g.won {
		overlay(screen, "TEN PAIRS LOCKED!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "SPAWN AREA BLOCKED\n\nTAP / ENTER TO RETRY")
	}
}

func drawPiece(screen *ebiten.Image, x, y int, c color.RGBA, label string) {
	cx := float32(boardX + x*cell + cell/2)
	cy := float32(boardY + y*cell + cell/2)
	vector.DrawFilledCircle(screen, cx, cy, cell/2-4, c, false)
	vector.StrokeCircle(screen, cx, cy, cell/2-4, 2, color.White, false)
	if label != "" {
		ebitenutil.DebugPrintAt(screen, label, int(cx)-3, int(cy)-5)
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
	ebitenutil.DebugPrintAt(screen, text, 108, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Falling Pair Workshop — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
