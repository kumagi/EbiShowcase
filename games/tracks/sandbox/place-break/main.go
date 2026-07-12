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
	cols, rows    = 10, 10
	cell          = 40
	ox, oy        = 40, 96
	empty         = 0
	dirt          = 1
	stone         = 2
	actionLimit   = 18
)

type point struct{ x, y int }

var blueprint = map[point]int{
	{3, 3}: dirt, {4, 3}: dirt, {5, 3}: stone,
	{3, 4}: dirt, {4, 4}: dirt, {5, 4}: stone,
}

type game struct {
	board             [rows][cols]int
	stock             [3]int
	selected, actions int
	breakMode         bool
	clear, over       bool
	message           string
}

func newGame() *game {
	g := &game{selected: dirt, message: "Break supplies, then copy the six-cell blueprint."}
	g.stock[dirt], g.stock[stone] = 2, 1
	// Extra material starts as blocks, so breaking changes the map and inventory.
	g.board[6][3], g.board[6][4] = dirt, dirt
	g.board[6][5] = stone
	return g
}

func inReach(x, y int) bool {
	// The builder stands at (4,5). Manhattan distance makes the range easy to see.
	dx, dy := x-4, y-5
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx+dy <= 3
}

func (g *game) useCell(x, y int) {
	if x < 0 || x >= cols || y < 0 || y >= rows || !inReach(x, y) {
		g.message = "That cell is outside the builder's three-step reach."
		return
	}
	if g.breakMode {
		kind := g.board[y][x]
		if kind == empty {
			g.message = "Nothing to break in that cell."
			return
		}
		g.board[y][x] = empty
		g.stock[kind]++
		g.actions++
		g.message = "Block returned to inventory."
	} else {
		if g.board[y][x] != empty {
			g.message = "Break the old block before placing here."
			return
		}
		if g.stock[g.selected] == 0 {
			g.message = "No selected blocks left — break one to recover it."
			return
		}
		g.board[y][x] = g.selected
		g.stock[g.selected]--
		g.actions++
		g.message = "Placed one block and spent one from inventory."
	}
	if g.matchesBlueprint() {
		g.clear = true
		g.message = "Blueprint complete — every tile value matches!"
	} else if g.actions >= actionLimit {
		g.over = true
		g.message = "Action limit reached. Plan the tile changes first!"
	}
}

func (g *game) matchesBlueprint() bool {
	for p, kind := range blueprint {
		if g.board[p.y][p.x] != kind {
			return false
		}
	}
	return true
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.breakMode = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.breakMode = false
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.selected = dirt
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.selected = stone
	}
	if x, y, ok := pressPosition(); ok {
		if y >= 600 {
			switch {
			case x < 160:
				g.breakMode = !g.breakMode
			case x < 320:
				g.selected = dirt
			default:
				g.selected = stone
			}
			return nil
		}
		if x < ox || x >= ox+cols*cell || y < oy || y >= oy+rows*cell {
			g.message = "Tap a cell inside the board."
			return nil
		}
		g.useCell((x-ox)/cell, (y-oy)/cell)
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 24, 38, 255})
	ebitenutil.DebugPrintAt(screen, "TIDEPOOL BUILDER", 184, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ACTIONS %02d/%02d  DIRT %d  STONE %d", g.actions, actionLimit, g.stock[dirt], g.stock[stone]), 106, 49)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			base := color.RGBA{24, 48, 61, 255}
			if inReach(x, y) {
				base = color.RGBA{31, 66, 73, 255}
			}
			vector.DrawFilledRect(screen, px+1, py+1, cell-2, cell-2, base, false)
			if want, ok := blueprint[point{x, y}]; ok && g.board[y][x] != want {
				vector.StrokeRect(screen, px+5, py+5, cell-10, cell-10, 3, tileColor(want), false)
			}
			if g.board[y][x] != empty {
				vector.DrawFilledRect(screen, px+4, py+4, cell-8, cell-8, tileColor(g.board[y][x]), false)
			}
		}
	}
	// Builder marker: its tile is still editable; only the range uses this center.
	vector.DrawFilledCircle(screen, ox+4*cell+cell/2, oy+5*cell+cell/2, 8, color.RGBA{101, 220, 226, 255}, true)
	ebitenutil.DebugPrintAt(screen, g.message, 54, 520)
	mode := "MODE: PLACE [P/B]"
	if g.breakMode {
		mode = "MODE: BREAK [P/B]"
	}
	button(screen, 10, mode, g.breakMode)
	button(screen, 170, "DIRT [1]", g.selected == dirt)
	button(screen, 330, "STONE [2]", g.selected == stone)
	ebitenutil.DebugPrintAt(screen, "Tap a highlighted cell. Outlines show the blueprint.", 62, 681)
	if g.clear || g.over {
		title := "TIDEPOOL HOME BUILT!"
		if g.over {
			title = "OUT OF ACTIONS!"
		}
		vector.DrawFilledRect(screen, 38, 260, 404, 166, color.RGBA{5, 14, 27, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 148, 307)
		ebitenutil.DebugPrintAt(screen, g.message, 74, 343)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 389)
	}
}

func tileColor(kind int) color.RGBA {
	if kind == stone {
		return color.RGBA{133, 150, 166, 255}
	}
	return color.RGBA{205, 137, 72, 255}
}

func button(screen *ebiten.Image, x int, label string, active bool) {
	c := color.RGBA{52, 78, 102, 255}
	if active {
		c = color.RGBA{219, 153, 61, 255}
	}
	vector.DrawFilledRect(screen, float32(x), 600, 140, 56, c, false)
	ebitenutil.DebugPrintAt(screen, label, x+18, 622)
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
	ebiten.SetWindowTitle("Tidepool Builder — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
