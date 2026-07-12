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
	timeLimit = 60 * 60
	goal      = 8
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
type enemyMode int

const (
	patrol enemyMode = iota
	chase
	search
)

var modeNames = [...]string{"PATROL", "CHASE", "SEARCH"}
var dirs = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
var patrolRoute = [...]point{{1, 1}, {9, 1}, {9, 3}, {1, 3}, {1, 5}, {9, 5}, {9, 7}, {1, 7}, {1, 9}, {9, 9}, {9, 11}, {1, 11}}

type game struct {
	player, enemy             point
	mode                      enemyMode
	patrolIndex               int
	lastSeen                  point
	searchTimer               int
	playerCooldown, enemyTick int
	badges                    map[point]bool
	collected, frames         int
	message                   string
	won, lost                 bool
}

func newGame() *game {
	g := &game{
		player:  point{5, 7},
		enemy:   point{1, 1},
		mode:    patrol,
		badges:  map[point]bool{},
		message: "Watch the patrol route. Break line of sight behind walls.",
	}
	for _, p := range []point{{3, 1}, {7, 1}, {1, 3}, {5, 3}, {9, 3}, {3, 5}, {7, 5}, {1, 9}, {5, 9}, {9, 9}, {3, 11}, {7, 11}} {
		g.badges[p] = true
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
		g.message = "Time up while the guard still controlled the maze."
		return nil
	}
	if g.playerCooldown > 0 {
		g.playerCooldown--
	}
	if g.playerCooldown == 0 {
		if d, ok := directionInput(); ok {
			next := point{g.player.x + d.x, g.player.y + d.y}
			if passable(next) {
				g.player = next
				g.playerCooldown = 7
				g.collectBadge()
			}
		}
	}

	g.updateMode()
	g.enemyTick++
	speed := 15
	if g.mode == chase {
		speed = 10
	}
	if g.enemyTick >= speed {
		g.enemyTick = 0
		g.moveEnemy()
	}
	if g.enemy == g.player {
		g.lost = true
		g.message = "The chasing guard reached your tile."
	}
	return nil
}

func (g *game) updateMode() {
	visible := hasLineOfSight(g.enemy, g.player) && manhattan(g.enemy, g.player) <= 7
	if visible {
		if g.mode != chase {
			g.message = "Seen in a straight corridor: PATROL -> CHASE!"
		}
		g.mode = chase
		g.lastSeen = g.player
		g.searchTimer = 100
		return
	}
	if g.mode == chase {
		g.mode = search
		g.message = "Line of sight broken: searching the last seen tile."
	}
	if g.mode == search {
		g.searchTimer--
		if g.searchTimer <= 0 || g.enemy == g.lastSeen {
			g.mode = patrol
			g.message = "Search ended: returning to the patrol route."
		}
	}
}

func (g *game) moveEnemy() {
	target := patrolRoute[g.patrolIndex]
	switch g.mode {
	case chase:
		target = g.player
	case search:
		target = g.lastSeen
	case patrol:
		if g.enemy == target {
			g.patrolIndex = (g.patrolIndex + 1) % len(patrolRoute)
			target = patrolRoute[g.patrolIndex]
		}
	}
	g.enemy = simpleStep(g.enemy, target)
}

func simpleStep(from, target point) point {
	// This lesson only gives each mode a target. STEP 05 will make junction
	// choice explicit; here horizontal preference keeps the rule small.
	candidates := []point{}
	if target.x < from.x {
		candidates = append(candidates, point{-1, 0})
	} else if target.x > from.x {
		candidates = append(candidates, point{1, 0})
	}
	if target.y < from.y {
		candidates = append(candidates, point{0, -1})
	} else if target.y > from.y {
		candidates = append(candidates, point{0, 1})
	}
	candidates = append(candidates, dirs[:]...)
	for _, d := range candidates {
		n := point{from.x + d.x, from.y + d.y}
		if passable(n) && manhattan(n, target) < manhattan(from, target) {
			return n
		}
	}
	return from
}

func hasLineOfSight(a, b point) bool {
	if a.x == b.x {
		lo, hi := min(a.y, b.y), max(a.y, b.y)
		for y := lo + 1; y < hi; y++ {
			if maze[y][a.x] == '#' {
				return false
			}
		}
		return true
	}
	if a.y == b.y {
		lo, hi := min(a.x, b.x), max(a.x, b.x)
		for x := lo + 1; x < hi; x++ {
			if maze[a.y][x] == '#' {
				return false
			}
		}
		return true
	}
	return false
}

func passable(p point) bool {
	return p.x >= 0 && p.x < cols && p.y >= 0 && p.y < rows && maze[p.y][p.x] != '#'
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

func (g *game) collectBadge() {
	if g.badges[g.player] {
		delete(g.badges, g.player)
		g.collected++
		g.message = fmt.Sprintf("Research badge collected: %d/%d", g.collected, goal)
		if g.collected >= goal {
			g.won = true
			g.message = "Eight badges recovered by reading the guard's modes!"
		}
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 18, 35, 255})
	seconds := max(0, (timeLimit-g.frames+59)/60)
	visible := hasLineOfSight(g.enemy, g.player) && manhattan(g.enemy, g.player) <= 7
	ebitenutil.DebugPrintAt(screen, "PATROL / CHASE MAZE", 173, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BADGES %d/%d   TIME %02d   ENEMY %s   VISIBLE %t", g.collected, goal, seconds, modeNames[g.mode], visible), 67, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 46, 70)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(mazeX+x*tileSize), float32(mazeY+y*tileSize)
			c := color.RGBA{24, 42, 57, 255}
			if maze[y][x] == '#' {
				c = color.RGBA{48, 76, 121, 255}
			}
			vector.DrawFilledRect(screen, px+1, py+1, tileSize-2, tileSize-2, c, false)
			if g.badges[point{x, y}] {
				vector.DrawFilledCircle(screen, px+tileSize/2, py+tileSize/2, 6, color.RGBA{244, 197, 72, 255}, false)
			}
		}
	}
	if visible {
		ex, ey := tileCenter(g.enemy)
		px, py := tileCenter(g.player)
		vector.StrokeLine(screen, ex, ey, px, py, 4, color.RGBA{242, 88, 86, 180}, false)
	}
	px, py := tileCenter(g.player)
	ex, ey := tileCenter(g.enemy)
	vector.DrawFilledCircle(screen, px, py, 12, color.RGBA{76, 194, 215, 255}, false)
	vector.StrokeCircle(screen, px, py, 12, 3, color.White, false)
	enemyColor := color.RGBA{239, 170, 62, 255}
	if g.mode == chase {
		enemyColor = color.RGBA{229, 76, 82, 255}
	}
	if g.mode == search {
		enemyColor = color.RGBA{171, 93, 201, 255}
	}
	vector.DrawFilledCircle(screen, ex, ey, 13, enemyColor, false)
	vector.StrokeCircle(screen, ex, ey, 13, 3, color.White, false)
	ebitenutil.DebugPrintAt(screen, modeNames[g.mode], int(ex)-20, int(ey)-28)

	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT"}
	for i, label := range labels {
		vector.DrawFilledRect(screen, float32(i*120+5), 610, 110, 62, color.RGBA{52, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*120+40, 636)
	}
	ebitenutil.DebugPrintAt(screen, "Arrows / WASD / hold mouse or touch button", 86, 691)
	if g.won {
		overlay(screen, "STEALTH SURVEY CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "THE GUARD CAUGHT YOU\n\nTAP / ENTER TO RETRY")
	}
}

func tileCenter(p point) (float32, float32) {
	return float32(mazeX + p.x*tileSize + tileSize/2), float32(mazeY + p.y*tileSize + tileSize/2)
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

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 102, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Patrol and Chase Maze — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
