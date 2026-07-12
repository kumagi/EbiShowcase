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
	screenW   = 480
	screenH   = 720
	cols      = 11
	rows      = 13
	tileSize  = 36
	mazeX     = 42
	mazeY     = 92
	timeLimit = 45 * 60
)

var maze = [...]string{"###########", "#.........#", "#.###.###.#", "#.........#", "###.#.#.###", "#.........#", "#.###.###.#", "#.........#", "###.#.#.###", "#.........#", "#.###.###.#", "#.........#", "###########"}

type point struct{ x, y int }

var dirs = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
var dirNames = [...]string{"LEFT", "UP", "DOWN", "RIGHT"}

type game struct {
	tile, target, current, queued point
	moving                        bool
	progress                      float64
	gates                         map[point]bool
	passed, turns, frames         int
	message                       string
	won, lost                     bool
}

func newGame() *game {
	g := &game{tile: point{1, 1}, target: point{1, 1}, current: point{1, 0}, queued: point{1, 0}, gates: map[point]bool{}, message: "Press the next direction early. QUEUED remembers it."}
	for _, p := range []point{{5, 1}, {5, 3}, {1, 3}, {1, 5}, {5, 5}, {9, 5}, {9, 7}, {5, 7}, {5, 9}, {1, 9}, {5, 11}, {9, 11}} {
		g.gates[p] = true
	}
	return g
}
func passable(p point) bool {
	return p.x >= 0 && p.x < cols && p.y >= 0 && p.y < rows && maze[p.y][p.x] != '#'
}
func add(a, b point) point { return point{a.x + b.x, a.y + b.y} }
func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.frames >= timeLimit {
		g.lost = true
		g.message = "Time up: queued turns were not ready."
		return nil
	}
	if d, ok := directionInput(); ok {
		g.queued = d
		g.message = "QUEUED " + nameOf(d) + " — it will be tested at the next center."
	}
	if g.moving {
		g.progress += .12
		if g.progress < 1 {
			return nil
		}
		g.tile = g.target
		g.progress = 0
		g.moving = false
		if g.gates[g.tile] {
			delete(g.gates, g.tile)
			g.passed++
			if g.passed >= 10 {
				g.won = true
				g.message = "Ten gates crossed with buffered turns!"
				return nil
			}
		}
	}
	// Only at a tile center: try the queued direction first, then continue straight.
	if passable(add(g.tile, g.queued)) {
		if g.queued != g.current {
			g.turns++
		}
		g.current = g.queued
		g.message = "Center reached: queued turn accepted."
	} else if !passable(add(g.tile, g.current)) {
		g.message = "Both queued and current directions are blocked. Choose a route."
		return nil
	}
	g.target = add(g.tile, g.current)
	g.moving = true
	return nil
}
func (g *game) screenPos() (float64, float64) {
	fx, fy := center(g.tile)
	tx, ty := center(g.target)
	return fx + (tx-fx)*g.progress, fy + (ty-fy)*g.progress
}
func center(p point) (float64, float64) {
	return float64(mazeX + p.x*tileSize + tileSize/2), float64(mazeY + p.y*tileSize + tileSize/2)
}
func nameOf(d point) string {
	for i, v := range dirs {
		if v == d {
			return dirNames[i]
		}
	}
	return "NONE"
}
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 18, 35, 255})
	sec := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, "BUFFERED TURN RUNNER", 170, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("GATES %02d/10 TIME %02d CURRENT %s QUEUED %s", g.passed, sec, nameOf(g.current), nameOf(g.queued)), 50, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 38, 70)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(mazeX+x*tileSize), float32(mazeY+y*tileSize)
			c := color.RGBA{23, 41, 57, 255}
			if maze[y][x] == '#' {
				c = color.RGBA{47, 75, 121, 255}
			}
			vector.DrawFilledRect(screen, px+1, py+1, tileSize-2, tileSize-2, c, false)
			if g.gates[point{x, y}] {
				vector.StrokeCircle(screen, px+tileSize/2, py+tileSize/2, 10, 3, color.RGBA{244, 193, 70, 255}, false)
			}
		}
	}
	px, py := g.screenPos()
	vector.DrawFilledCircle(screen, float32(px), float32(py), 12, color.RGBA{232, 91, 80, 255}, false)
	vector.StrokeCircle(screen, float32(px), float32(py), 12, 3, color.White, false)
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT"}
	for i, l := range labels {
		c := color.RGBA{51, 83, 122, 255}
		if dirs[i] == g.queued {
			c = color.RGBA{178, 106, 57, 255}
		}
		vector.DrawFilledRect(screen, float32(i*120+5), 610, 110, 62, c, false)
		ebitenutil.DebugPrintAt(screen, l, i*120+40, 636)
	}
	ebitenutil.DebugPrintAt(screen, "Press before the corner: arrows / WASD / touch", 89, 691)
	if g.won {
		overlay(screen, "BUFFERED ROUTE CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "TIME UP\n\nTAP / ENTER TO RETRY")
	}
}
func directionInput() (point, bool) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return dirs[0], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		return dirs[1], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		return dirs[2], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return dirs[3], true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			return dirs[min(3, x/120)], true
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 610 {
			return dirs[min(3, x/120)], true
		}
	}
	return point{}, false
}
func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	ebitenutil.DebugPrintAt(s, t, 112, 328)
}
func (g *game) Layout(int, int) (int, int) { return screenW, screenH }
func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Buffered Turn Runner")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
