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
	cols      = 9
	rows      = 9
	cell      = 48
	boardX    = 24
	boardY    = 106
	fuse      = 150
	blastLife = 28
	moveEvery = 11
)

type point struct{ x, y int }
type bomb struct {
	at    point
	timer int
}
type flame struct {
	at    point
	timer int
}
type node struct {
	at   point
	step int
	path []point
}

type game struct {
	player, enemy               point
	bombs                       []bomb
	flames                      []flame
	danger                      map[point]int
	path                        []point
	frames, enemyTick, survived int
	message                     string
	won, lost                   bool
}

func newGame() *game {
	g := &game{player: point{1, 1}, enemy: point{7, 7}, danger: map[point]int{}, message: "Place bombs. The blue scout predicts their blast times."}
	g.predictAndPlan()
	return g
}

func wall(p point) bool {
	return p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || p.x == 0 || p.y == 0 || p.x == cols-1 || p.y == rows-1 || (p.x%2 == 0 && p.y%2 == 0)
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if d, ok := inputDir(); ok {
		n := point{g.player.x + d.x, g.player.y + d.y}
		if !wall(n) && !g.bombAt(n) {
			g.player = n
		}
	}
	if placePressed() && len(g.bombs) < 3 && !g.bombAt(g.player) {
		g.bombs = append(g.bombs, bomb{g.player, fuse})
		g.message = "Prediction updated: numbers show frames until danger."
		g.predictAndPlan()
	}
	g.updateBombs()
	g.updateFlames()
	if g.frames%6 == 0 {
		g.predictAndPlan()
	}
	g.enemyTick++
	if g.enemyTick >= moveEvery && len(g.path) > 1 {
		g.enemy = g.path[1]
		g.enemyTick = 0
		g.predictAndPlan()
	}
	if g.player == g.enemy {
		g.lost = true
		g.message = "You bumped into the scout. Give it room to escape."
	}
	if g.frames >= 75*60 {
		g.lost = true
		g.message = "Time up: detonate six bombs to complete the test."
	}
	return nil
}

func (g *game) updateBombs() {
	next := g.bombs[:0]
	for _, b := range g.bombs {
		b.timer--
		if b.timer <= 0 {
			for _, p := range blastCells(b.at) {
				g.flames = append(g.flames, flame{p, blastLife})
			}
			g.survived++
			if g.survived >= 6 {
				g.won = true
				g.message = "The scout escaped six predicted blasts!"
			}
		} else {
			next = append(next, b)
		}
	}
	g.bombs = next
}

func (g *game) updateFlames() {
	next := g.flames[:0]
	for _, f := range g.flames {
		f.timer--
		if f.at == g.player || f.at == g.enemy {
			g.lost = true
			g.message = "A blast caught someone. Retry and leave a safe route."
		}
		if f.timer > 0 {
			next = append(next, f)
		}
	}
	g.flames = next
}

func blastCells(origin point) []point {
	cells := []point{origin}
	for _, d := range []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		for n := 1; n <= 2; n++ {
			p := point{origin.x + d.x*n, origin.y + d.y*n}
			if wall(p) {
				break
			}
			cells = append(cells, p)
		}
	}
	return cells
}

func (g *game) predictAndPlan() {
	g.danger = map[point]int{}
	for _, b := range g.bombs {
		for _, p := range blastCells(b.at) {
			if old, ok := g.danger[p]; !ok || b.timer < old {
				g.danger[p] = b.timer
			}
		}
	}
	g.path = g.safePath()
	if len(g.path) == 0 {
		g.path = []point{g.enemy}
	}
}

func (g *game) safePath() []point {
	queue := []node{{at: g.enemy, path: []point{g.enemy}}}
	best := map[point]int{g.enemy: 0}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		arrival := cur.step * moveEvery
		if cur.step > 0 && g.safeForWindow(cur.at, arrival, arrival+50) {
			return cur.path
		}
		if cur.step >= 14 {
			continue
		}
		for _, d := range []point{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			next := point{cur.at.x + d.x, cur.at.y + d.y}
			nextStep := cur.step + 1
			if wall(next) || g.bombAt(next) || !g.safeAt(next, nextStep*moveEvery) {
				continue
			}
			if seen, ok := best[next]; ok && seen <= nextStep {
				continue
			}
			best[next] = nextStep
			path := append(append([]point{}, cur.path...), next)
			queue = append(queue, node{next, nextStep, path})
		}
	}
	return nil
}

func (g *game) safeAt(p point, arrival int) bool {
	t, threatened := g.danger[p]
	return !threatened || arrival+8 < t || arrival > t+blastLife
}

func (g *game) safeForWindow(p point, from, until int) bool {
	t, threatened := g.danger[p]
	return !threatened || until < t || from > t+blastLife
}

func (g *game) bombAt(p point) bool {
	for _, b := range g.bombs {
		if b.at == p {
			return true
		}
	}
	return false
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{7, 15, 30, 255})
	ebitenutil.DebugPrintAt(s, "FUTURE-SAFE SCOUT", 175, 17)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SURVIVED %d/6   BOMBS %d/3   TIME %02d", g.survived, len(g.bombs), max(0, 75-g.frames/60)), 105, 43)
	ebitenutil.DebugPrintAt(s, g.message, 28, 72)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := point{x, y}
			px := float32(boardX + x*cell)
			py := float32(boardY + y*cell)
			c := color.RGBA{25, 48, 52, 255}
			if wall(p) {
				c = color.RGBA{63, 83, 117, 255}
			}
			vector.DrawFilledRect(s, px+1, py+1, cell-2, cell-2, c, false)
			if t, ok := g.danger[p]; ok {
				vector.DrawFilledRect(s, px+5, py+5, cell-10, cell-10, color.RGBA{168, 51, 58, 95}, false)
				ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d", max(0, t)), int(px)+15, int(py)+18)
			}
		}
	}
	for i, p := range g.path {
		if i == 0 {
			continue
		}
		x := float32(boardX + p.x*cell + cell/2)
		y := float32(boardY + p.y*cell + cell/2)
		vector.DrawFilledCircle(s, x, y, 5, color.RGBA{83, 224, 238, 210}, false)
	}
	for _, b := range g.bombs {
		x := float32(boardX + b.at.x*cell + cell/2)
		y := float32(boardY + b.at.y*cell + cell/2)
		vector.DrawFilledCircle(s, x, y, 14, color.RGBA{18, 21, 30, 255}, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%.1f", float64(b.timer)/60), int(x)-10, int(y)-5)
	}
	for _, f := range g.flames {
		x := float32(boardX + f.at.x*cell + cell/2)
		y := float32(boardY + f.at.y*cell + cell/2)
		vector.DrawFilledCircle(s, x, y, 20, color.RGBA{255, 157, 48, 225}, false)
	}
	drawActor(s, g.player, color.RGBA{235, 91, 76, 255}, "P")
	drawActor(s, g.enemy, color.RGBA{72, 190, 235, 255}, "AI")
	drawControls(s)
	if g.won {
		overlay(s, "PREDICTION COMPLETE!\n\nTAP / ENTER TO RETRY")
	} else if g.lost {
		overlay(s, "ESCAPE TEST FAILED\n\nTAP / ENTER TO RETRY")
	}
}

func drawActor(s *ebiten.Image, p point, c color.RGBA, label string) {
	x := float32(boardX + p.x*cell + cell/2)
	y := float32(boardY + p.y*cell + cell/2)
	vector.DrawFilledCircle(s, x, y, 14, c, false)
	ebitenutil.DebugPrintAt(s, label, int(x)-7, int(y)-5)
}

func drawControls(s *ebiten.Image) {
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT", "BOMB"}
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*96+3), 600, 90, 62, color.RGBA{45, 78, 113, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*96+27, 627)
	}
	ebitenutil.DebugPrintAt(s, "Red numbers = danger time | Cyan dots = AI route", 67, 681)
}

func inputDir() (point, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return point{-1, 0}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return point{0, -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return point{0, 1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		return point{1, 0}, true
	}
	if x, y, ok := press(); ok && y >= 600 && x < 384 {
		return [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}[min(3, x/96)], true
	}
	return point{}, false
}

func placePressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		return true
	}
	x, y, ok := press()
	return ok && y >= 600 && x >= 384
}
func press() (int, int, bool) {
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
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 40, 270, 400, 160, color.RGBA{4, 12, 27, 245}, false)
	ebitenutil.DebugPrintAt(s, t, 112, 330)
}
func (g *game) Layout(int, int) (int, int) { return screenW, screenH }
func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Future-safe Scout")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
