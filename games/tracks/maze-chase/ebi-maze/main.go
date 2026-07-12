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
	timeLimit = 90 * 60
	goal      = 20
)

var maze = [...]string{
	"###########", "#.........#", "#.###.###.#", "#.........#",
	"###.#.#.###", "#.........#", "#.###.###.#", "#.........#",
	"###.#.#.###", "#.........#", "#.###.###.#", "#.........#", "###########",
}

type point struct{ x, y int }
type mode int

const (
	patrol mode = iota
	chase
	search
)

var dirs = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
var dirNames = [...]string{"LEFT", "UP", "DOWN", "RIGHT"}
var modeNames = [...]string{"PATROL", "CHASE", "SEARCH"}
var patrolGoals = [2][4]point{
	{{1, 1}, {9, 1}, {9, 11}, {1, 11}},
	{{9, 11}, {1, 11}, {1, 1}, {9, 1}},
}

type runner struct {
	tile, target point
	dir, wanted  int
	moving       bool
	progress     float64
}

type guard struct {
	tile, dir, lastSeen point
	mode                mode
	routeIndex          int
	searchTimer, tick   int
}

type game struct {
	player               runner
	guards               [2]guard
	pellets              map[point]bool
	collected, lives     int
	frames, invulnerable int
	message              string
	won, lost            bool
}

func newGame() *game {
	g := &game{pellets: map[point]bool{}, lives: 3}
	g.resetActors()
	n := 0
	for y := 1; y < rows-1 && n < goal; y++ {
		for x := 1; x < cols-1 && n < goal; x++ {
			p := point{x, y}
			if passable(p) && p != g.player.tile && p != g.guards[0].tile && p != g.guards[1].tile && (x+y)%2 == 1 {
				g.pellets[p] = true
				n++
			}
		}
	}
	g.message = "Buffer turns, read guard modes, and collect every pearl."
	return g
}

func (g *game) resetActors() {
	g.player = runner{tile: point{5, 7}, target: point{5, 7}, dir: 3, wanted: 3}
	g.guards[0] = guard{tile: point{1, 1}, dir: point{1, 0}, mode: patrol}
	g.guards[1] = guard{tile: point{9, 11}, dir: point{-1, 0}, mode: patrol}
	g.invulnerable = 90
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
		g.message = "Time up with pearls still in the maze."
		return nil
	}
	if g.invulnerable > 0 {
		g.invulnerable--
	}
	if d, ok := directionInput(); ok {
		g.player.wanted = d
	}
	g.updatePlayer()
	for i := range g.guards {
		g.updateGuard(i)
	}
	g.checkContacts()
	return nil
}

func (g *game) updatePlayer() {
	p := &g.player
	if p.moving {
		p.progress += 0.15
		if p.progress < 1 {
			return
		}
		p.tile, p.target = p.target, p.target
		p.progress = 0
		p.moving = false
		g.collect()
	}
	chosen := p.wanted
	if !passable(add(p.tile, dirs[chosen])) {
		chosen = p.dir
	}
	next := add(p.tile, dirs[chosen])
	if passable(next) {
		p.dir = chosen
		p.target = next
		p.moving = true
	}
}

func (g *game) updateGuard(index int) {
	e := &g.guards[index]
	visible := hasLineOfSight(e.tile, g.player.tile) && manhattan(e.tile, g.player.tile) <= 7
	if visible {
		e.mode = chase
		e.lastSeen = g.player.tile
		e.searchTimer = 110
	} else if e.mode == chase {
		e.mode = search
	} else if e.mode == search {
		e.searchTimer--
		if e.searchTimer <= 0 || e.tile == e.lastSeen {
			e.mode = patrol
		}
	}
	e.tick++
	stepEvery := 15
	if e.mode == chase {
		stepEvery = 11
	}
	if e.tick < stepEvery {
		return
	}
	e.tick = 0
	target := patrolGoals[index][e.routeIndex]
	if e.mode == chase {
		target = g.player.tile
	} else if e.mode == search {
		target = e.lastSeen
	} else if e.tile == target {
		e.routeIndex = (e.routeIndex + 1) % len(patrolGoals[index])
		target = patrolGoals[index][e.routeIndex]
	}
	e.dir = chooseAtJunction(e.tile, e.dir, target)
	e.tile = add(e.tile, e.dir)
}

func chooseAtJunction(from, current, target point) point {
	best, bestDistance := point{}, 1<<30
	reverse := point{-current.x, -current.y}
	legal := []point{}
	for _, d := range dirs {
		if passable(add(from, d)) {
			legal = append(legal, d)
		}
	}
	for _, d := range legal {
		if len(legal) > 1 && d == reverse {
			continue
		}
		distance := manhattan(add(from, d), target)
		if distance < bestDistance {
			best, bestDistance = d, distance
		}
	}
	if best == (point{}) && len(legal) > 0 {
		best = legal[0]
	}
	return best
}

func (g *game) collect() {
	if g.pellets[g.player.tile] {
		delete(g.pellets, g.player.tile)
		g.collected++
		g.message = fmt.Sprintf("Pearl %d/%d collected from the data layer.", g.collected, goal)
		if g.collected >= goal {
			g.won = true
			g.message = "Every pearl collected while both guards were active!"
		}
	}
}

func (g *game) checkContacts() {
	if g.invulnerable > 0 || g.won {
		return
	}
	for _, e := range g.guards {
		if e.tile == g.player.tile || (g.player.moving && e.tile == g.player.target) {
			g.lives--
			if g.lives <= 0 {
				g.lost = true
				g.message = "All three rescue attempts were used."
				return
			}
			g.message = fmt.Sprintf("Caught! Actors reset, but pearls stay collected. Lives %d.", g.lives)
			g.resetActors()
			return
		}
	}
}

func passable(p point) bool {
	return p.x >= 0 && p.x < cols && p.y >= 0 && p.y < rows && maze[p.y][p.x] != '#'
}
func add(a, b point) point     { return point{a.x + b.x, a.y + b.y} }
func manhattan(a, b point) int { return abs(a.x-b.x) + abs(a.y-b.y) }
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func hasLineOfSight(a, b point) bool {
	if a.x == b.x {
		for y := min(a.y, b.y) + 1; y < max(a.y, b.y); y++ {
			if maze[y][a.x] == '#' {
				return false
			}
		}
		return true
	}
	if a.y == b.y {
		for x := min(a.x, b.x) + 1; x < max(a.x, b.x); x++ {
			if maze[a.y][x] == '#' {
				return false
			}
		}
		return true
	}
	return false
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 18, 35, 255})
	seconds := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, "EBI MAZE / INTEGRATED RUN", 149, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PEARLS %02d/%02d   LIVES %d   TIME %02d   DIR %s -> %s", g.collected, goal, g.lives, seconds, dirNames[g.player.dir], dirNames[g.player.wanted]), 44, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 43, 70)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(mazeX+x*tileSize), float32(mazeY+y*tileSize)
			c := color.RGBA{23, 42, 58, 255}
			if maze[y][x] == '#' {
				c = color.RGBA{47, 74, 119, 255}
			}
			vector.DrawFilledRect(screen, px+1, py+1, tileSize-2, tileSize-2, c, false)
			if g.pellets[point{x, y}] {
				vector.DrawFilledCircle(screen, px+tileSize/2, py+tileSize/2, 5, color.RGBA{245, 199, 75, 255}, false)
			}
		}
	}
	px, py := runnerPosition(g.player)
	playerColor := color.RGBA{74, 195, 216, 255}
	if g.invulnerable > 0 && g.invulnerable%10 < 5 {
		playerColor.A = 80
	}
	vector.DrawFilledCircle(screen, px, py, 12, playerColor, false)
	vector.StrokeCircle(screen, px, py, 12, 3, color.White, false)
	for i, e := range g.guards {
		ex, ey := tileCenter(e.tile)
		c := color.RGBA{237, 169, 62, 255}
		if e.mode == chase {
			c = color.RGBA{226, 75, 82, 255}
		}
		if e.mode == search {
			c = color.RGBA{170, 91, 202, 255}
		}
		vector.DrawFilledCircle(screen, ex, ey, 13, c, false)
		vector.StrokeCircle(screen, ex, ey, 13, 3, color.White, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("G%d %s", i+1, modeNames[e.mode]), int(ex)-25, int(ey)-28)
	}
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT"}
	for i, label := range labels {
		vector.DrawFilledRect(screen, float32(i*120+5), 610, 110, 62, color.RGBA{52, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*120+40, 636)
	}
	ebitenutil.DebugPrintAt(screen, "Input during motion is buffered for the next center", 74, 691)
	if g.won {
		overlay(screen, "EBI MAZE CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "MAZE RESCUE FAILED\n\nTAP / ENTER TO RETRY")
	}
}

func runnerPosition(r runner) (float32, float32) {
	fx, fy := tileCenter(r.tile)
	tx, ty := tileCenter(r.target)
	return fx + (tx-fx)*float32(r.progress), fy + (ty-fy)*float32(r.progress)
}
func tileCenter(p point) (float32, float32) {
	return float32(mazeX + p.x*tileSize + tileSize/2), float32(mazeY + p.y*tileSize + tileSize/2)
}

func directionInput() (int, bool) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return 0, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		return 1, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		return 2, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return 3, true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			return min(3, x/120), true
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 610 {
			return min(3, x/120), true
		}
	}
	return 0, false
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
	ebiten.SetWindowTitle("Ebi Maze — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
