package main

import (
	"fmt"
	"image/color"
	"math"

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
	goal      = 12
)

var maze = [...]string{
	"###########",
	"#.........#",
	"#.###.###.#",
	"#.........#",
	"###.#.#.###",
	"#.........#",
	"#.###.###.#",
	"#.........#",
	"###.#.#.###",
	"#.........#",
	"#.###.###.#",
	"#.........#",
	"###########",
}

type point struct{ x, y int }

var directions = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

type game struct {
	tile, target      point
	moving            bool
	progress          float64
	pearls            map[point]bool
	collected, frames int
	message           string
	won, lost         bool
}

func newGame() *game {
	g := &game{
		tile:    point{1, 1},
		target:  point{1, 1},
		pearls:  map[point]bool{},
		message: "At a tile center, hold a legal direction.",
	}
	n := 0
	for y := 1; y < rows-1; y++ {
		for x := 1; x < cols-1; x++ {
			p := point{x, y}
			if maze[y][x] == '.' && p != g.tile && (x+y)%2 == 1 {
				g.pearls[p] = true
				n++
				if n >= 24 {
					return g
				}
			}
		}
	}
	return g
}

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
		g.message = "Time up before all twelve pearls were collected."
		return nil
	}
	if g.moving {
		g.progress += 0.105
		if g.progress >= 1 {
			g.tile = g.target
			g.target = g.tile
			g.progress = 0
			g.moving = false
			g.collect()
		}
		return nil
	}

	dir, ok := directionInput()
	if !ok {
		return nil
	}
	next := point{g.tile.x + dir.x, g.tile.y + dir.y}
	if !passable(next) {
		g.message = fmt.Sprintf("Tile (%d,%d) is a wall. Stay at the center.", next.x, next.y)
		return nil
	}
	g.target = next
	g.moving = true
	g.progress = 0
	g.message = fmt.Sprintf("Legal tile: moving center (%d,%d) -> (%d,%d).", g.tile.x, g.tile.y, next.x, next.y)
	return nil
}

func passable(p point) bool {
	return p.x >= 0 && p.x < cols && p.y >= 0 && p.y < rows && maze[p.y][p.x] != '#'
}

func (g *game) collect() {
	if g.pearls[g.tile] {
		delete(g.pearls, g.tile)
		g.collected++
		g.message = fmt.Sprintf("Pearl collected at tile (%d,%d): %d/%d", g.tile.x, g.tile.y, g.collected, goal)
		if g.collected >= goal {
			g.won = true
			g.message = "Twelve pearls collected along legal tile centers!"
		}
	}
}

func (g *game) screenPosition() (float64, float64) {
	fromX, fromY := center(g.tile)
	toX, toY := center(g.target)
	return fromX + (toX-fromX)*g.progress, fromY + (toY-fromY)*g.progress
}

func center(p point) (float64, float64) {
	return float64(mazeX + p.x*tileSize + tileSize/2), float64(mazeY + p.y*tileSize + tileSize/2)
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 18, 35, 255})
	seconds := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, "TILE CENTER RUNNER", 181, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PEARLS %02d/%02d   TIME %02d   TILE (%d,%d)", g.collected, goal, seconds, g.tile.x, g.tile.y), 103, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 45, 70)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(mazeX+x*tileSize), float32(mazeY+y*tileSize)
			if maze[y][x] == '#' {
				vector.DrawFilledRect(screen, px+1, py+1, tileSize-2, tileSize-2, color.RGBA{48, 76, 121, 255}, false)
			} else {
				vector.DrawFilledRect(screen, px+1, py+1, tileSize-2, tileSize-2, color.RGBA{24, 42, 57, 255}, false)
				vector.DrawFilledCircle(screen, px+tileSize/2, py+tileSize/2, 2, color.RGBA{111, 139, 153, 255}, false)
			}
			if g.pearls[point{x, y}] {
				vector.DrawFilledCircle(screen, px+tileSize/2, py+tileSize/2, 6, color.RGBA{247, 201, 81, 255}, false)
			}
		}
	}
	if g.moving {
		tx, ty := center(g.target)
		vector.StrokeCircle(screen, float32(tx), float32(ty), 13, 3, color.RGBA{247, 201, 81, 180}, false)
	}
	px, py := g.screenPosition()
	vector.DrawFilledCircle(screen, float32(px), float32(py), 12, color.RGBA{232, 92, 79, 255}, false)
	vector.StrokeCircle(screen, float32(px), float32(py), 12, 3, color.White, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CENTER ERROR %.1f px", centerError(px, py)), 170, 574)

	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT"}
	for i, label := range labels {
		vector.DrawFilledRect(screen, float32(i*120+5), 610, 110, 62, color.RGBA{52, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*120+40, 636)
	}
	ebitenutil.DebugPrintAt(screen, "Arrows / WASD / hold mouse or touch button", 86, 691)
	if g.won {
		overlay(screen, "MAZE ROUTE CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "TIME UP\n\nTAP / ENTER TO RETRY")
	}
}

func centerError(x, y float64) float64 {
	nearestX := float64(mazeX+tileSize/2) + math.Round((x-float64(mazeX+tileSize/2))/tileSize)*tileSize
	nearestY := float64(mazeY+tileSize/2) + math.Round((y-float64(mazeY+tileSize/2))/tileSize)*tileSize
	return math.Hypot(x-nearestX, y-nearestY)
}

func directionInput() (point, bool) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return directions[0], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		return directions[1], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		return directions[2], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return directions[3], true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			return directions[min(3, x/120)], true
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 610 {
			return directions[min(3, x/120)], true
		}
	}
	return point{}, false
}

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 112, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Tile Center Runner — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
