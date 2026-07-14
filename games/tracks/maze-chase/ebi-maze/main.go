package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const (
	screenW  = 480
	screenH  = 720
	cols     = 11
	rows     = 13
	tileSize = 36
	mazeX    = 42
	mazeY    = 92
)

type point struct{ x, y int }
type stageData struct {
	name                string
	maze                [rows]string
	guards, sight, goal int
	step                int
	wall                color.RGBA
}

var stages = []stageData{
	{"PEARL GARDEN", [rows]string{"###########", "#.........#", "#.###.###.#", "#.........#", "###.#.#.###", "#.........#", "#.###.###.#", "#.........#", "###.#.#.###", "#.........#", "#.###.###.#", "#.........#", "###########"}, 1, 5, 16, 16, color.RGBA{50, 92, 142, 255}},
	{"CROSSROAD VAULT", [rows]string{"###########", "#...#.....#", "#.#.#.###.#", "#.#.......#", "#.###.#.#.#", "#.....#.#.#", "###.#...#.#", "#...###...#", "#.#.....#.#", "#.#####.#.#", "#.......#.#", "#...#.....#", "###########"}, 2, 7, 20, 13, color.RGBA{89, 69, 151, 255}},
	{"MOON LABYRINTH", [rows]string{"###########", "#.........#", "#.###.#.#.#", "#...#.#.#.#", "###.#...#.#", "#...###...#", "#.#.....#.#", "#.#.###.#.#", "#...#...#.#", "#.###.###.#", "#.........#", "#...#.....#", "###########"}, 3, 9, 24, 10, color.RGBA{138, 61, 101, 255}},
}
var dirs = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
var dirNames = [...]string{"LEFT", "UP", "DOWN", "RIGHT"}

type mode int

const (
	patrol mode = iota
	chase
	search
	ambush
)

var modeNames = [...]string{"PATROL", "CHASE", "SEARCH", "AMBUSH"}

type runner struct {
	tile, target point
	dir, wanted  int
	moving       bool
	progress     float64
}
type guard struct {
	tile, dir, lastSeen point
	mode                mode
	route, tick         int
	phase               float64
}
type particle struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}
type game struct {
	stage                                                                                  int
	player                                                                                 runner
	guards                                                                                 []guard
	pellets                                                                                map[point]bool
	collected, total, lives, frames, stageFrames, invuln, score, best, combo, shake, flash int
	particles                                                                              []particle
	message                                                                                string
	won, lost                                                                              bool
	rng                                                                                    *rand.Rand
	audio                                                                                  *audio.Context
	gate                                                                                   audiolab.Gate
	pulse                                                                                  *shaderlab.Pulse
	cam                                                                                    cameralab.State
	badge                                                                                  *ebiten.Image
}

var sessionBest int

func newGame() *game {
	g := &game{lives: 3, best: sessionBest, rng: rand.New(rand.NewSource(70))}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: screenW / 2, Y: screenH / 2}, ViewW: screenW, ViewH: screenH}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{255, 220, 104, 255})
	g.loadStage(0)
	return g
}
func (g *game) loadStage(n int) {
	g.stage = n
	g.stageFrames = 0
	g.collected = 0
	g.pellets = map[point]bool{}
	g.player = runner{tile: point{1, 11}, target: point{1, 11}, dir: 3, wanted: 3}
	spawns := []point{{9, 1}, {9, 11}, {1, 1}}
	g.guards = nil
	for i := 0; i < stages[n].guards; i++ {
		g.guards = append(g.guards, guard{tile: spawns[i], dir: point{-1, 0}, route: i})
	}
	count := 0
	for y := 1; y < rows-1 && count < stages[n].goal; y++ {
		for x := 1; x < cols-1 && count < stages[n].goal; x++ {
			p := point{x, y}
			if g.passable(p) && p != g.player.tile && !g.guardAt(p) && (x+y+n)%2 == 1 {
				g.pellets[p] = true
				count++
			}
		}
	}
	g.total = count
	g.invuln = 100
	g.message = "Read the guard label, buffer a turn, and collect every pearl."
}
func (g *game) guardAt(p point) bool {
	for _, e := range g.guards {
		if e.tile == p {
			return true
		}
	}
	return false
}
func (g *game) passable(p point) bool {
	return p.x >= 0 && p.x < cols && p.y >= 0 && p.y < rows && stages[g.stage].maze[p.y][p.x] != '#'
}
func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	g.stageFrames++
	if g.invuln > 0 {
		g.invuln--
	}
	if d, ok := directionInput(); ok {
		g.player.wanted = d
	}
	g.updatePlayer()
	for i := range g.guards {
		g.updateGuard(i)
	}
	g.contacts()
	g.updateParticles()
	if g.shake > 0 {
		g.shake--
	}
	if g.flash > 0 {
		g.flash--
	}
	return nil
}
func (g *game) updatePlayer() {
	p := &g.player
	if p.moving {
		p.progress += .17
		if p.progress < 1 {
			return
		}
		p.tile = p.target
		p.progress = 0
		p.moving = false
		g.collect()
	}
	chosen := p.wanted
	if !g.passable(add(p.tile, dirs[chosen])) {
		chosen = p.dir
	}
	next := add(p.tile, dirs[chosen])
	if g.passable(next) {
		p.dir = chosen
		p.target = next
		p.moving = true
	}
}
func (g *game) collect() {
	if !g.pellets[g.player.tile] {
		return
	}
	delete(g.pellets, g.player.tile)
	g.collected++
	g.combo++
	g.score += 50 + g.combo*3
	g.play(500)
	x, y := tileCenter(g.player.tile)
	for i := 0; i < 8; i++ {
		a := float64(i) * math.Pi / 4
		g.particles = append(g.particles, particle{float64(x), float64(y), math.Cos(a) * 1.8, math.Sin(a) * 1.8, 22, color.RGBA{255, 220, 104, 255}})
	}
	if g.collected >= g.total {
		g.score += 800 + max(0, 1200-g.stageFrames/3) + g.lives*100
		g.flash = 14
		for i := 0; i < 35; i++ {
			g.particles = append(g.particles, particle{240, 330, g.rng.Float64()*6 - 3, g.rng.Float64()*-4 - 1, 50, color.RGBA{80 + uint8(i*5), 200, 230, 255}})
		}
		if g.stage == len(stages)-1 {
			g.won = true
			if g.score > sessionBest {
				sessionBest = g.score
			}
			g.best = sessionBest
			g.message = "All three mazes cleared! Find a faster route for a new best."
		} else {
			g.loadStage(g.stage + 1)
			g.message = "Maze clear! The next guards remember more and move faster."
		}
	} else {
		g.message = fmt.Sprintf("Pearl chain x%d — %d/%d", g.combo, g.collected, g.total)
	}
}
func (g *game) updateGuard(i int) {
	e := &g.guards[i]
	d := stages[g.stage]
	visible := g.lineOfSight(e.tile, g.player.tile) && manhattan(e.tile, g.player.tile) <= d.sight
	if visible {
		e.mode = chase
		e.lastSeen = g.player.tile
	} else if e.mode == chase {
		e.mode = search
	} else if e.mode == search && e.tile == e.lastSeen {
		e.mode = patrol
	}
	if g.stage == 2 && i == 2 {
		e.mode = ambush
	}
	e.tick++
	wait := d.step
	if e.mode == chase || e.mode == ambush {
		wait -= 2
	}
	if e.tick < wait {
		return
	}
	e.tick = 0
	target := []point{{1, 1}, {9, 1}, {5, 6}, {9, 11}}[(e.route+g.stage)%4]
	switch e.mode {
	case chase:
		target = g.player.tile
	case search:
		target = e.lastSeen
	case ambush:
		target = g.player.tile
		for k := 0; k < 3; k++ {
			next := add(target, dirs[g.player.dir])
			if g.passable(next) {
				target = next
			}
		}
	}
	e.dir = g.choose(e.tile, e.dir, target)
	e.tile = add(e.tile, e.dir)
	e.phase += .8
}
func (g *game) choose(from, current, target point) point {
	best := point{}
	distance := 1 << 30
	reverse := point{-current.x, -current.y}
	legal := []point{}
	for _, d := range dirs {
		if g.passable(add(from, d)) {
			legal = append(legal, d)
		}
	}
	for _, d := range legal {
		if len(legal) > 1 && d == reverse {
			continue
		}
		v := manhattan(add(from, d), target)
		if v < distance {
			best, distance = d, v
		}
	}
	if best == (point{}) && len(legal) > 0 {
		best = legal[0]
	}
	return best
}
func (g *game) contacts() {
	if g.invuln > 0 {
		return
	}
	for _, e := range g.guards {
		if e.tile == g.player.tile || (g.player.moving && e.tile == g.player.target) {
			g.lives--
			g.play(160)
			g.combo = 0
			g.shake = 14
			g.flash = 10
			for i := 0; i < 18; i++ {
				a := float64(i) * .6
				g.particles = append(g.particles, particle{240, 350, math.Cos(a) * 3, math.Sin(a) * 3, 32, color.RGBA{255, 83, 105, 255}})
			}
			if g.lives <= 0 {
				g.lost = true
				g.message = "Caught three times. Read each AI mode and try another route."
				return
			}
			g.resetActors()
			g.message = fmt.Sprintf("Caught! Pearls stay collected. %d lives remain.", g.lives)
			return
		}
	}
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .05)).Play()
}
func (g *game) resetActors() {
	g.player = runner{tile: point{1, 11}, target: point{1, 11}, dir: 3, wanted: 3}
	spawns := []point{{9, 1}, {9, 11}, {1, 1}}
	for i := range g.guards {
		g.guards[i].tile = spawns[i]
		g.guards[i].mode = patrol
	}
	g.invuln = 100
}
func (g *game) updateParticles() {
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .03
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}
}
func (g *game) lineOfSight(a, b point) bool {
	if a.x == b.x {
		for y := min(a.y, b.y) + 1; y < max(a.y, b.y); y++ {
			if !g.passable(point{a.x, y}) {
				return false
			}
		}
		return true
	}
	if a.y == b.y {
		for x := min(a.x, b.x) + 1; x < max(a.x, b.x); x++ {
			if !g.passable(point{x, a.y}) {
				return false
			}
		}
		return true
	}
	return false
}
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 15, 31, 255})
	dx, dy := 0, 0
	if g.shake > 0 {
		dx = g.rng.Intn(7) - 3
		dy = g.rng.Intn(7) - 3
	}
	world := ebiten.NewImage(screenW, screenH)
	d := stages[g.stage]
	g.drawTitle(world, d.name)
	g.drawEffectBadge(world)
	ebitenutil.DebugPrintAt(world, fmt.Sprintf("PEARLS %02d/%02d  LIFE %d  SCORE %05d  %s>%s", g.collected, g.total, g.lives, g.score, dirNames[g.player.dir], dirNames[g.player.wanted]), 35, 44)
	ebitenutil.DebugPrintAt(world, g.message, 35, 69)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(mazeX+x*tileSize), float32(mazeY+y*tileSize)
			if d.maze[y][x] == '#' {
				vector.DrawFilledRect(world, px+1, py+1, tileSize-2, tileSize-2, d.wall, false)
				vector.StrokeRect(world, px+4, py+4, tileSize-8, tileSize-8, 2, color.RGBA{190, 220, 235, 80}, false)
			} else {
				vector.DrawFilledRect(world, px+1, py+1, tileSize-2, tileSize-2, color.RGBA{16, 29, 49, 255}, false)
			}
			if g.pellets[point{x, y}] {
				pulse := float64(13) + math.Sin(float64(g.frames)*.14+float64(x+y))*2
				trackatlas.DrawCentered(world, "pearl", float64(px+18), float64(py+18), pulse)
			}
		}
	}
	px, py := runnerPosition(g.player)
	bob := math.Sin(float64(g.frames)*.25) * 2
	alpha := float32(1)
	if g.invuln > 0 && g.invuln%10 < 5 {
		alpha = .3
	}
	trackatlas.DrawTinted(world, "hero", float64(px), float64(py)+bob, 30+math.Sin(float64(g.frames)*.22)*2, 1, 1, 1, alpha)
	for i, e := range g.guards {
		ex, ey := tileCenter(e.tile)
		sprite := "ghost-patrol"
		if e.mode == chase || e.mode == ambush {
			sprite = "ghost-chase"
		} else if e.mode == search {
			sprite = "ghost-search"
		}
		trackatlas.DrawCentered(world, sprite, float64(ex), float64(ey)+math.Sin(float64(g.frames)*.18+float64(i))*3, 30)
		ebitenutil.DebugPrintAt(world, modeNames[e.mode], int(ex)-22, int(ey)-28)
	}
	for _, p := range g.particles {
		vector.DrawFilledCircle(world, float32(p.x), float32(p.y), 3, p.c, false)
	}
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT"}
	for i, label := range labels {
		vector.DrawFilledRect(world, float32(i*120+5), 610, 110, 62, color.RGBA{52, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(world, label, i*120+40, 636)
	}
	ebitenutil.DebugPrintAt(world, "Buffer the next turn / arrows, WASD, or touch", 75, 691)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	screen.DrawImage(world, op)
	if g.flash > 0 {
		vector.DrawFilledRect(screen, 0, 0, screenW, screenH, color.RGBA{255, 255, 255, 60}, false)
	}
	if g.won {
		overlay(screen, fmt.Sprintf("THREE MAZES CLEAR!\n\nSCORE %05d  BEST %05d\n\nTAP / ENTER TO RACE AGAIN", g.score, g.best))
	}
	if g.lost {
		overlay(screen, fmt.Sprintf("MAZE RUN ENDED\n\nSCORE %05d  BEST %05d\n\nTAP / ENTER TO RETRY", g.score, sessionBest))
	}
}
func (g *game) drawTitle(screen *ebiten.Image, name string) {
	label := fmt.Sprintf("EBI MAZE  %d/3  %s", g.stage+1, name)
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(125, 5)
		text.Draw(screen, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(screen, label, 125, 17)
}
func (g *game) drawEffectBadge(screen *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frames)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenW-34, 10)
	screen.DrawImage(fx, op)
}
func runnerPosition(r runner) (float32, float32) {
	x, y := tileCenter(r.tile)
	tx, ty := tileCenter(r.target)
	return x + (tx-x)*float32(r.progress), y + (ty-y)*float32(r.progress)
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
func add(a, b point) point     { return point{a.x + b.x, a.y + b.y} }
func manhattan(a, b point) int { return abs(a.x-b.x) + abs(a.y-b.y) }
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func overlay(s *ebiten.Image, text string) {
	vector.DrawFilledRect(s, 32, 235, 416, 245, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(s, 32, 235, 416, 245, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(s, text, 104, 285)
}
func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }
func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Maze Marathon")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
