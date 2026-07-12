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
	cols, rows    = 8, 12
	cell          = 46
	ox, oy        = 56, 76
	dropEvery     = 42
	targetLines   = 4
)

type game struct {
	board         [rows][cols]bool
	x, y, timer   int
	lines, pieces int
	cleared, over bool
	message       string
}

func newGame() *game {
	g := &game{}
	// Four nearly complete rows make the lesson goal easy to see. Each gap is
	// exactly the width of the falling domino, but appears in a new column.
	gaps := [4]int{1, 5, 3, 0}
	for i, gap := range gaps {
		y := rows - 1 - i
		for x := 0; x < cols; x++ {
			g.board[y][x] = x != gap && x != gap+1
		}
	}
	g.spawn()
	g.message = "Fill each two-cell gap and clear 4 lines!"
	return g
}

func (g *game) spawn() {
	g.x, g.y, g.timer = cols/2-1, 0, 0
	if g.blocked(g.x, g.y) {
		g.over = true
		g.message = "No room to spawn — the board topped out."
	}
}

// blocked checks both squares of the horizontal domino.
func (g *game) blocked(x, y int) bool {
	if x < 0 || x+1 >= cols || y >= rows {
		return true
	}
	return y >= 0 && (g.board[y][x] || g.board[y][x+1])
}

func (g *game) move(dx int) {
	if !g.blocked(g.x+dx, g.y) {
		g.x += dx
	}
}

func (g *game) fallOne() {
	if !g.blocked(g.x, g.y+1) {
		g.y++
		return
	}
	g.board[g.y][g.x], g.board[g.y][g.x+1] = true, true
	g.pieces++
	n := g.clearFullRows()
	if n > 0 {
		g.lines += n
		g.message = fmt.Sprintf("Full row found! %d line(s) compacted.", n)
	} else {
		g.message = "Locked into the board. Find the next gap."
	}
	if g.lines >= targetLines {
		g.cleared = true
		g.message = "Four full rows cleared — compaction mastered!"
		return
	}
	g.spawn()
}

// clearFullRows copies every non-full row downward. write is the next row
// that should receive data, so empty rows naturally collect at the top.
func (g *game) clearFullRows() int {
	write, cleared := rows-1, 0
	for read := rows - 1; read >= 0; read-- {
		full := true
		for x := 0; x < cols; x++ {
			if !g.board[read][x] {
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
		g.board[write] = [cols]bool{}
		write--
	}
	return cleared
}

func (g *game) Update() error {
	if g.cleared || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	left := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	down := inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if x, y, ok := pressPosition(); ok && y >= 648 {
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
		g.fallOne()
		g.timer = 0
		return nil
	}
	g.timer++
	if g.timer >= dropEvery {
		g.timer = 0
		g.fallOne()
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 18, 36, 255})
	ebitenutil.DebugPrintAt(screen, "LOCK & LINE WORKSHOP", 160, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("LINES %d/%d   PIECES %d", g.lines, targetLines, g.pieces), 156, 47)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.StrokeRect(screen, px, py, cell-2, cell-2, 1, color.RGBA{47, 68, 101, 255}, false)
			if g.board[y][x] {
				drawBlock(screen, px, py, color.RGBA{65, 157, 207, 255})
			}
		}
	}
	if !g.cleared && !g.over {
		for dx := 0; dx < 2; dx++ {
			drawBlock(screen, float32(ox+(g.x+dx)*cell), float32(oy+g.y*cell), color.RGBA{239, 174, 54, 255})
		}
	}
	ebitenutil.DebugPrintAt(screen, g.message, 62, 612)
	button(screen, 10, "LEFT")
	button(screen, 170, "DOWN")
	button(screen, 330, "RIGHT")
	if g.cleared || g.over {
		title := "LINE WORKSHOP CLEAR!"
		if g.over {
			title = "BOARD TOPPED OUT!"
		}
		vector.DrawFilledRect(screen, 40, 276, 400, 160, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 154, 322)
		ebitenutil.DebugPrintAt(screen, g.message, 76, 354)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 397)
	}
}

func drawBlock(screen *ebiten.Image, x, y float32, c color.Color) {
	vector.DrawFilledRect(screen, x+3, y+3, cell-8, cell-8, c, false)
	vector.StrokeRect(screen, x+6, y+6, cell-14, cell-14, 2, color.RGBA{255, 255, 255, 125}, false)
}

func button(screen *ebiten.Image, x int, label string) {
	vector.DrawFilledRect(screen, float32(x), 648, 140, 54, color.RGBA{52, 79, 117, 255}, false)
	ebitenutil.DebugPrintAt(screen, label, x+48, 670)
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

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Lock & Line Workshop — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
