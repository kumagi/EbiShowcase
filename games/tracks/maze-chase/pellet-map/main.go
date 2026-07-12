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
	cols, rows    = 13, 13
	cell          = 32
	ox, oy        = 32, 108
	timeLimit     = 60 * 60
)

var levelText = [rows]string{
	"#############",
	"#...........#",
	"#.###.#.###.#",
	"#.#...#...#.#",
	"#.#.#####.#.#",
	"#...........#",
	"###.#.#.#.###",
	"#...#...#...#",
	"#.#####.###.#",
	"#.....#.....#",
	"#.###.#.###.#",
	"#...........#",
	"#############",
}

type point struct{ x, y int }

var moves = []point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

type game struct {
	walls                [rows][cols]bool
	pellets              [rows][cols]bool
	player, enemy        point
	remaining, collected int
	frames, moveTimer    int
	enemyTimer           int
	clear, over          bool
	message              string
}

func newGame() *game {
	g := &game{player: point{1, 1}, enemy: point{11, 11}}
	for y, row := range levelText {
		for x, tile := range row {
			g.walls[y][x] = tile == '#'
			g.pellets[y][x] = tile == '.'
			if g.pellets[y][x] {
				g.remaining++
			}
		}
	}
	g.removePellet(g.player)
	g.removePellet(g.enemy)
	g.message = "Walls stay fixed; every eaten pellet flips one state cell."
	return g
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.frames >= timeLimit {
		g.over = true
		g.message = "Time up before the pellet layer became empty."
		return nil
	}
	direction, moving := readDirection()
	if moving {
		g.moveTimer++
		if g.moveTimer >= 6 {
			g.moveTimer = 0
			g.tryMovePlayer(direction)
		}
	} else {
		g.moveTimer = 5
	}

	g.enemyTimer++
	if g.enemyTimer >= 16 {
		g.enemyTimer = 0
		g.moveEnemy()
	}
	if g.player == g.enemy {
		g.over = true
		g.message = "The red map bug reached your tile. Plan another route!"
	}
	return nil
}

func readDirection() (point, bool) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return moves[0], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		return moves[1], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		return moves[2], true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return moves[3], true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 575 {
			return moves[min(3, x/120)], true
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 575 {
			return moves[min(3, x/120)], true
		}
	}
	return point{}, false
}

func (g *game) tryMovePlayer(direction point) {
	next := point{g.player.x + direction.x, g.player.y + direction.y}
	if g.walls[next.y][next.x] {
		g.message = "The wall layer blocked that move; pellet data was untouched."
		return
	}
	g.player = next
	if g.pellets[next.y][next.x] {
		g.removePellet(next)
		g.collected++
		g.message = fmt.Sprintf("pellets[%d][%d] = false; remaining = %d", next.y, next.x, g.remaining)
		if g.remaining == 0 {
			g.clear = true
			g.message = "No pellets remain: the mutable layer is empty!"
		}
	}
}

func (g *game) removePellet(p point) {
	if g.pellets[p.y][p.x] {
		g.pellets[p.y][p.x] = false
		g.remaining--
	}
}

func (g *game) moveEnemy() {
	best := g.enemy
	bestDistance := manhattan(g.enemy, g.player)
	for _, direction := range moves {
		candidate := point{g.enemy.x + direction.x, g.enemy.y + direction.y}
		if g.walls[candidate.y][candidate.x] {
			continue
		}
		distance := manhattan(candidate, g.player)
		if distance < bestDistance {
			best, bestDistance = candidate, distance
		}
	}
	g.enemy = best
}

func manhattan(a, b point) int {
	dx, dy := a.x-b.x, a.y-b.y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 18, 34, 255})
	ebitenutil.DebugPrintAt(screen, "PELLET LAYER MAZE", 184, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("REMAINING %02d   COLLECTED %02d   TIME %02d", g.remaining, g.collected, max(0, 60-g.frames/60)), 119, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 44, 72)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			if g.walls[y][x] {
				vector.DrawFilledRect(screen, px+1, py+1, cell-2, cell-2, color.RGBA{37, 75, 119, 255}, false)
				vector.StrokeRect(screen, px+4, py+4, cell-8, cell-8, 2, color.RGBA{79, 139, 189, 255}, false)
			} else {
				vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{35, 48, 64, 255}, false)
				if g.pellets[y][x] {
					vector.DrawFilledCircle(screen, px+cell/2, py+cell/2, 4, color.RGBA{244, 203, 93, 255}, false)
				}
			}
		}
	}
	drawActor(screen, g.player, color.RGBA{238, 118, 73, 255}, "P")
	drawActor(screen, g.enemy, color.RGBA{218, 77, 91, 255}, "E")
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PLAYER TILE (%d,%d)  WALL=%t  PELLET=%t", g.player.x, g.player.y, g.walls[g.player.y][g.player.x], g.pellets[g.player.y][g.player.x]), 84, 542)
	labels := []string{"LEFT", "UP", "DOWN", "RIGHT"}
	for i, label := range labels {
		vector.DrawFilledRect(screen, float32(i*120+5), 575, 110, 76, color.RGBA{51, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*120+38, 609)
	}
	ebitenutil.DebugPrintAt(screen, "Hold arrows / WASD / mouse / touch", 115, 681)
	if g.clear {
		overlay(screen, "PELLET LAYER EMPTY!\n\nTAP / ENTER TO RETRY")
	} else if g.over {
		overlay(screen, "MAZE RUN FAILED\n\nTAP / ENTER TO RETRY")
	}
}

func drawActor(screen *ebiten.Image, p point, fill color.RGBA, label string) {
	cx := float32(ox + p.x*cell + cell/2)
	cy := float32(oy + p.y*cell + cell/2)
	vector.DrawFilledCircle(screen, cx, cy, 12, fill, false)
	vector.StrokeCircle(screen, cx, cy, 12, 2, color.White, false)
	ebitenutil.DebugPrintAt(screen, label, int(cx)-3, int(cy)-5)
}

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 104, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Pellet Layer Maze — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
