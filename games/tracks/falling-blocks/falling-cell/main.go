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
	screenWidth  = 480
	screenHeight = 720
	cols         = 8
	rows         = 12
	cellSize     = 46
	boardX       = 56
	boardY       = 90
	dropFrames   = 38
	targetCells  = 16
)

type game struct {
	board               [rows][cols]bool
	x, y, timer, locked int
	clear, over         bool
	message             string
}

func newGame() *game {
	g := &game{}
	g.spawn()
	g.message = "Guide the cell. Lock 16 cells to clear!"
	return g
}

func (g *game) spawn() {
	g.x, g.y, g.timer = cols/2, 0, 0
	if g.board[g.y][g.x] {
		g.over = true
		g.message = "The spawn cell is blocked. Try again!"
	}
}

func (g *game) blocked(x, y int) bool {
	return x < 0 || x >= cols || y >= rows || (y >= 0 && g.board[y][x])
}

func (g *game) move(dx int) {
	if !g.blocked(g.x+dx, g.y) {
		g.x += dx
	}
}

func (g *game) stepDown() {
	if !g.blocked(g.x, g.y+1) {
		g.y++
		return
	}
	g.board[g.y][g.x] = true
	g.locked++
	if g.locked >= targetCells {
		g.clear = true
		g.message = "16 cells locked — timer mastered!"
		return
	}
	g.spawn()
}

func (g *game) Update() error {
	if g.clear || g.over {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}

	left := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	down := inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)
	if x, y, ok := justPressedPosition(); ok && y >= 648 {
		switch {
		case x < 160:
			left = true
		case x < 320:
			down = true
		default:
			right = true
		}
	}
	if left {
		g.move(-1)
	}
	if right {
		g.move(1)
	}
	if down {
		g.stepDown()
		g.timer = 0
		return nil
	}

	g.timer++
	if g.timer >= dropFrames {
		g.timer = 0
		g.stepDown()
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 19, 38, 255})
	ebitenutil.DebugPrintAt(screen, "FALLING CELL LAB", 174, 24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("LOCKED %02d / %02d", g.locked, targetCells), 190, 52)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px := float32(boardX + x*cellSize)
			py := float32(boardY + y*cellSize)
			vector.StrokeRect(screen, px, py, cellSize-2, cellSize-2, 1, color.RGBA{45, 69, 102, 255}, false)
			if g.board[y][x] {
				drawCell(screen, px, py, color.RGBA{83, 172, 226, 255})
			}
		}
	}
	if !g.clear && !g.over {
		drawCell(screen, float32(boardX+g.x*cellSize), float32(boardY+g.y*cellSize), color.RGBA{245, 181, 62, 255})
	}
	progress := float32(g.timer) / dropFrames
	vector.DrawFilledRect(screen, 56, 616, 368, 8, color.RGBA{35, 53, 82, 255}, false)
	vector.DrawFilledRect(screen, 56, 616, 368*progress, 8, color.RGBA{245, 181, 62, 255}, false)
	vector.DrawFilledRect(screen, 10, 648, 140, 54, color.RGBA{49, 84, 126, 255}, false)
	vector.DrawFilledRect(screen, 170, 648, 140, 54, color.RGBA{186, 115, 48, 255}, false)
	vector.DrawFilledRect(screen, 330, 648, 140, 54, color.RGBA{49, 84, 126, 255}, false)
	ebitenutil.DebugPrintAt(screen, "LEFT", 64, 670)
	ebitenutil.DebugPrintAt(screen, "DOWN", 221, 670)
	ebitenutil.DebugPrintAt(screen, "RIGHT", 380, 670)
	if g.clear || g.over {
		title := "TIME KEEPER CLEAR!"
		if g.over {
			title = "STACKED TOO HIGH!"
		}
		vector.DrawFilledRect(screen, 45, 272, 390, 160, color.RGBA{6, 15, 31, 246}, false)
		ebitenutil.DebugPrintAt(screen, title, 160, 320)
		ebitenutil.DebugPrintAt(screen, g.message, 85, 350)
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE TO RETRY", 145, 390)
	}
}

func drawCell(screen *ebiten.Image, x, y float32, c color.Color) {
	vector.DrawFilledRect(screen, x+3, y+3, cellSize-8, cellSize-8, c, false)
	vector.StrokeRect(screen, x+5, y+5, cellSize-12, cellSize-12, 2, color.RGBA{255, 255, 255, 130}, false)
}

func justPressedPosition() (int, int, bool) {
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

func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Falling Cell Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
