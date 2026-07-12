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
	cols, rows    = 8, 10
	cell, ox, oy  = 44, 64, 112
	gravityFrames = 42
)

type point struct{ x, y int }

type shape struct {
	name   string
	blocks [4]point
	color  color.RGBA
}

var shapes = []shape{
	{"I", [4]point{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, color.RGBA{65, 190, 220, 255}},
	{"O", [4]point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, color.RGBA{246, 190, 55, 255}},
	{"T", [4]point{{0, 0}, {1, 0}, {2, 0}, {1, 1}}, color.RGBA{180, 100, 220, 255}},
	{"S", [4]point{{1, 0}, {2, 0}, {0, 1}, {1, 1}}, color.RGBA{91, 193, 112, 255}},
}

var targetX = []int{0, 5, 2, 4}

type game struct {
	shapeIndex int
	x, y       int
	tick       int
	successes  int
	misses     int
	waiting    bool
	lastFit    bool
	win        bool
	gameOver   bool
	message    string
}

func newGame() *game {
	g := &game{}
	g.spawn()
	return g
}

func (g *game) spawn() {
	g.x, g.y = 2, 0
	g.tick = 0
	g.waiting = false
	g.message = "Move all four cells over the dotted landing shape."
}

func (g *game) Update() error {
	if g.win || g.gameOver {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}

	dx, drop, pressed := g.readControls()
	if g.waiting {
		if pressed || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if !g.lastFit {
				g.spawn()
			} else if g.shapeIndex == len(shapes)-1 {
				g.win = true
			} else {
				g.shapeIndex++
				g.spawn()
			}
		}
		return nil
	}

	if dx != 0 && g.canPlace(g.x+dx, g.y) {
		g.x += dx
	}
	if drop {
		for g.canPlace(g.x, g.y+1) {
			g.y++
		}
		g.land()
		return nil
	}

	g.tick++
	if g.tick >= gravityFrames || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.tick = 0
		if g.canPlace(g.x, g.y+1) {
			g.y++
		} else {
			g.land()
		}
	}
	return nil
}

func (g *game) readControls() (int, bool, bool) {
	dx, drop, pressed := 0, false, false
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dx, pressed = -1, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dx, pressed = 1, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		drop, pressed = true, true
	}
	if x, y, ok := justPressed(); ok && y >= 588 {
		pressed = true
		switch {
		case x < 160:
			dx = -1
		case x < 320:
			drop = true
		default:
			dx = 1
		}
	}
	return dx, drop, pressed
}

func (g *game) canPlace(x, y int) bool {
	for _, offset := range shapes[g.shapeIndex].blocks {
		px, py := x+offset.x, y+offset.y
		if px < 0 || px >= cols || py < 0 || py >= rows {
			return false
		}
	}
	return true
}

func (g *game) landingY() int {
	maxY := 0
	for _, p := range shapes[g.shapeIndex].blocks {
		if p.y > maxY {
			maxY = p.y
		}
	}
	return rows - 1 - maxY
}

func (g *game) land() {
	g.y = g.landingY()
	if g.x == targetX[g.shapeIndex] {
		g.successes++
		g.lastFit = true
		g.message = "Perfect fit! The four offsets stayed together."
	} else {
		g.misses++
		g.lastFit = false
		g.message = fmt.Sprintf("Missed the outline. Mistakes %d/3.", g.misses)
		if g.misses >= 3 {
			g.gameOver = true
			return
		}
	}
	g.waiting = true
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{16, 25, 43, 255})
	ebitenutil.DebugPrintAt(screen, "FOUR-CELL SHAPE DOCK", 165, 28)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SHAPE %d/4: %s   FITS %d   MISSES %d/3", g.shapeIndex+1, shapes[g.shapeIndex].name, g.successes, g.misses), 96, 55)
	ebitenutil.DebugPrintAt(screen, g.message, 62, 80)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{50, 70, 98, 255}, false)
		}
	}
	shape := shapes[g.shapeIndex]
	ly := g.landingY()
	for _, offset := range shape.blocks {
		drawCell(screen, targetX[g.shapeIndex]+offset.x, ly+offset.y, color.RGBA{52, 70, 93, 255}, true)
	}
	for _, offset := range shape.blocks {
		drawCell(screen, g.x+offset.x, g.y+offset.y, shape.color, false)
	}

	button(screen, 28, 590, 132, "LEFT", color.RGBA{55, 111, 145, 255})
	button(screen, 174, 590, 132, "DROP", color.RGBA{240, 177, 65, 255})
	button(screen, 320, 590, 132, "RIGHT", color.RGBA{55, 111, 145, 255})
	ebitenutil.DebugPrintAt(screen, "Arrow keys / A D   Space = hard drop", 101, 675)
	if g.waiting {
		ebitenutil.DebugPrintAt(screen, "TAP A BUTTON / SPACE FOR NEXT SHAPE", 102, 560)
	}
	if g.win {
		overlay(screen, "ALL FOUR SHAPES DELIVERED!\n\nTAP / SPACE TO RESTART")
	} else if g.gameOver {
		overlay(screen, "DOCK CLOSED: THREE MISSES\n\nTAP / SPACE TO RETRY")
	}
}

func drawCell(screen *ebiten.Image, x, y int, fill color.RGBA, outline bool) {
	px, py := float32(ox+x*cell), float32(oy+y*cell)
	if !outline {
		vector.DrawFilledRect(screen, px+3, py+3, cell-6, cell-6, fill, false)
	}
	line := color.RGBA{190, 218, 235, 255}
	if !outline {
		line = color.RGBA{245, 250, 255, 255}
	}
	vector.StrokeRect(screen, px+3, py+3, cell-6, cell-6, 3, line, false)
}

func button(screen *ebiten.Image, x, y, w int, label string, fill color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), 64, fill, false)
	ebitenutil.DebugPrintAt(screen, label, x+w/2-len(label)*3, y+28)
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
	vector.DrawFilledRect(screen, 48, 280, 384, 160, color.RGBA{5, 15, 31, 245}, false)
	ebitenutil.DebugPrintAt(screen, message, 104, 333)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Four-cell Shape Dock — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
