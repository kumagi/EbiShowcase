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
	width, height    = 480, 720
	cols, rows, cell = 6, 6, 64
	ox, oy           = 48, 160
	stageInspect     = 0
	stageSpecial     = 1
	stageBlast       = 2
	specialRow       = 0
	specialColor     = 1
	specialArea      = 2
)

var pieceColors = []color.RGBA{
	{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255},
	{105, 194, 119, 255}, {177, 94, 218, 255},
}

var recipeNames = []string{"4 IN A ROW", "5 IN A ROW", "L SHAPE"}
var specialNames = []string{"ROW ROCKET", "COLOR WAVE", "AREA BOMB"}

type point struct{ x, y int }

type game struct {
	board      [rows][cols]int
	recipe     int
	stage      int
	special    int
	specialPos point
	cleared    map[point]bool
	done       [3]bool
	message    string
	win        bool
}

func newGame() *game {
	g := &game{}
	g.loadRecipe(0)
	return g
}

func (g *game) loadRecipe(recipe int) {
	g.recipe = recipe
	g.stage = stageInspect
	g.cleared = map[point]bool{}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = (x + y*2) % len(pieceColors)
		}
	}
	switch recipe {
	case specialRow:
		for x := 1; x <= 4; x++ {
			g.board[3][x] = 0
		}
	case specialColor:
		for x := 0; x <= 4; x++ {
			g.board[3][x] = 1
		}
	case specialArea:
		for y := 2; y <= 4; y++ {
			g.board[y][1] = 2
		}
		for x := 1; x <= 3; x++ {
			g.board[4][x] = 2
		}
	}
	g.message = "Inspect the bright outline, then press FORGE."
}

func (g *game) Update() error {
	if g.win {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}

	for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if inpututil.IsKeyJustPressed(key) {
			g.loadRecipe(i)
		}
	}
	action := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if x, y, ok := justPressed(); ok {
		if y >= 72 && y <= 128 {
			for i := 0; i < 3; i++ {
				left := 24 + i*148
				if x >= left && x < left+136 {
					g.loadRecipe(i)
				}
			}
		}
		if y >= 605 {
			action = true
		}
	}
	if action {
		g.advance()
	}
	return nil
}

func (g *game) advance() {
	switch g.stage {
	case stageInspect:
		kind, pos, ok := detectSpecial(g.board)
		if !ok || kind != g.recipe {
			g.message = "No matching shape yet. Inspect the outline again."
			return
		}
		g.special, g.specialPos = kind, pos
		g.stage = stageSpecial
		g.message = specialNames[kind] + " forged! Press ACTIVATE."
	case stageSpecial:
		g.activate()
		g.done[g.recipe] = true
		g.stage = stageBlast
		g.message = fmt.Sprintf("Effect reached %d cells. Press NEXT RECIPE.", len(g.cleared))
	case stageBlast:
		if g.done[0] && g.done[1] && g.done[2] {
			g.win = true
			return
		}
		for step := 1; step <= 3; step++ {
			next := (g.recipe + step) % 3
			if !g.done[next] {
				g.loadRecipe(next)
				return
			}
		}
	}
}

func detectSpecial(board [rows][cols]int) (int, point, bool) {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			h := runLength(board, x, y, 1, 0)
			v := runLength(board, x, y, 0, 1)
			if h >= 3 && v >= 3 {
				return specialArea, point{x, y}, true
			}
			if h >= 5 || v >= 5 {
				return specialColor, point{x, y}, true
			}
			if h >= 4 || v >= 4 {
				return specialRow, point{x, y}, true
			}
		}
	}
	return 0, point{}, false
}

func runLength(board [rows][cols]int, x, y, dx, dy int) int {
	kind := board[y][x]
	count := 1
	for px, py := x-dx, y-dy; px >= 0 && px < cols && py >= 0 && py < rows && board[py][px] == kind; px, py = px-dx, py-dy {
		count++
	}
	for px, py := x+dx, y+dy; px >= 0 && px < cols && py >= 0 && py < rows && board[py][px] == kind; px, py = px+dx, py+dy {
		count++
	}
	return count
}

func (g *game) activate() {
	g.cleared = map[point]bool{}
	switch g.special {
	case specialRow:
		for x := 0; x < cols; x++ {
			g.cleared[point{x, g.specialPos.y}] = true
		}
	case specialColor:
		kind := g.board[g.specialPos.y][g.specialPos.x]
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				if g.board[y][x] == kind {
					g.cleared[point{x, y}] = true
				}
			}
		}
	case specialArea:
		for y := g.specialPos.y - 1; y <= g.specialPos.y+1; y++ {
			for x := g.specialPos.x - 1; x <= g.specialPos.x+1; x++ {
				if x >= 0 && x < cols && y >= 0 && y < rows {
					g.cleared[point{x, y}] = true
				}
			}
		}
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(screen, "SPECIAL PIECE WORKSHOP", 154, 28)
	ebitenutil.DebugPrintAt(screen, g.message, 62, 53)
	for i, name := range recipeNames {
		x := float32(24 + i*148)
		fill := color.RGBA{40, 57, 82, 255}
		if i == g.recipe {
			fill = color.RGBA{46, 112, 140, 255}
		}
		vector.DrawFilledRect(screen, x, 72, 136, 56, fill, false)
		mark := "[ ]"
		if g.done[i] {
			mark = "[OK]"
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %s", i+1, name), int(x)+8, 85)
		ebitenutil.DebugPrintAt(screen, mark, int(x)+52, 105)
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := point{x, y}
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			fill := pieceColors[g.board[y][x]]
			if g.cleared[p] {
				fill = color.RGBA{245, 249, 255, 255}
			}
			vector.DrawFilledRect(screen, px+4, py+4, cell-8, cell-8, fill, false)
			if g.stage == stageInspect && isRecipeCell(g.recipe, x, y) {
				vector.StrokeRect(screen, px+2, py+2, cell-4, cell-4, 5, color.White, false)
			}
		}
	}
	if g.stage >= stageSpecial {
		px := float32(ox + g.specialPos.x*cell + cell/2)
		py := float32(oy + g.specialPos.y*cell + cell/2)
		vector.DrawFilledCircle(screen, px, py, 18, color.RGBA{20, 28, 48, 255}, false)
		vector.StrokeCircle(screen, px, py, 18, 4, color.White, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("S%d", g.special+1), int(px)-8, int(py)-7)
	}

	vector.DrawFilledRect(screen, 55, 605, 370, 66, color.RGBA{240, 177, 65, 255}, false)
	label := "FORGE SPECIAL [SPACE]"
	if g.stage == stageSpecial {
		label = "ACTIVATE EFFECT [SPACE]"
	} else if g.stage == stageBlast {
		label = "NEXT RECIPE [SPACE]"
	}
	ebitenutil.DebugPrintAt(screen, label, 145, 632)
	ebitenutil.DebugPrintAt(screen, "Tap recipe cards or press 1 / 2 / 3", 103, 685)
	if g.win {
		overlay(screen, "ALL THREE EFFECTS MASTERED!\n\nTAP / SPACE TO RESTART")
	}
}

func isRecipeCell(recipe, x, y int) bool {
	switch recipe {
	case specialRow:
		return y == 3 && x >= 1 && x <= 4
	case specialColor:
		return y == 3 && x <= 4
	case specialArea:
		return (x == 1 && y >= 2 && y <= 4) || (y == 4 && x >= 1 && x <= 3)
	}
	return false
}

func justPressed() (int, int, bool) {
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
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func overlay(screen *ebiten.Image, message string) {
	vector.DrawFilledRect(screen, 48, 285, 384, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(screen, message, 103, 334)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Special Piece Workshop — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
