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
	screenW    = 480
	screenH    = 720
	cols       = 9
	rows       = 9
	cell       = 48
	boardX     = 24
	boardY     = 100
	fuse       = 300
	blastLife  = 34
	blastRange = 3
	maxBombs   = 5
	goal       = 5
	timeLimit  = 70 * 60
)

type point struct{ x, y int }

var rays = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

type bombState int

const (
	waiting bombState = iota
	blasting
)

type bomb struct {
	at         point
	timer      int
	state      bombState
	blast      []point
	chainOrder int
}

type game struct {
	player               point
	bombs                []bomb
	targets              map[point]bool
	cleared, frames      int
	lastChain, bestChain int
	queueLog             string
	message              string
	won, lost            bool
}

func newGame() *game {
	g := &game{
		player:  point{1, 1},
		targets: map[point]bool{},
		message: "Place bombs in one blast line, then escape before the first fuse ends.",
	}
	for _, p := range []point{{7, 1}, {1, 3}, {1, 4}, {3, 3}, {5, 3}, {7, 3}, {1, 7}, {7, 7}} {
		g.targets[p] = true
	}
	return g
}

func hardWall(p point) bool {
	return p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || (p.x%2 == 0 && p.y%2 == 0)
}

func blastCells(center point) []point {
	cells := []point{center}
	for _, direction := range rays {
		for step := 1; step <= blastRange; step++ {
			next := point{center.x + direction.x*step, center.y + direction.y*step}
			if hardWall(next) {
				break
			}
			cells = append(cells, next)
		}
	}
	return cells
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
		g.message = "Time up. Put bombs within three open cells of each other."
		return nil
	}
	if direction, ok := directionInput(); ok {
		next := point{g.player.x + direction.x, g.player.y + direction.y}
		if !hardWall(next) && !g.bombAt(next) {
			g.player = next
		}
	}
	if placePressed() {
		if len(g.bombs) >= maxBombs {
			g.message = "Bomb rack empty: wait for the current bombs to finish."
		} else if g.bombAt(g.player) {
			g.message = "There is already a bomb on this cell."
		} else {
			// Later bombs have a longer natural fuse, so a quick explosion proves
			// that the event queue—not their own timers—started them.
			timer := fuse + len(g.bombs)*70
			g.bombs = append(g.bombs, bomb{at: g.player, timer: timer, state: waiting})
			g.message = fmt.Sprintf("Bomb %d placed. Link three with clear blast rays.", len(g.bombs))
		}
	}

	starts := make([]int, 0, len(g.bombs))
	for i := range g.bombs {
		if g.bombs[i].state == waiting {
			g.bombs[i].timer--
			if g.bombs[i].timer <= 0 {
				starts = append(starts, i)
			}
		} else {
			g.bombs[i].timer--
		}
	}
	if len(starts) > 0 {
		g.propagate(starts)
	}

	nextBombs := g.bombs[:0]
	for i := range g.bombs {
		b := g.bombs[i]
		if b.state == blasting {
			for _, p := range b.blast {
				if p == g.player {
					g.lost = true
					g.message = "A linked blast reached the player's cell."
				}
				if g.targets[p] {
					delete(g.targets, p)
					g.cleared++
				}
			}
		}
		if b.state == waiting || b.timer > 0 {
			nextBombs = append(nextBombs, b)
		}
	}
	g.bombs = nextBombs
	if g.cleared >= goal && g.bestChain >= 3 && !g.lost {
		g.won = true
		g.message = "Five relays cleared with a chain of three or more bombs!"
	}
	return nil
}

func (g *game) propagate(starts []int) {
	queue := append([]int(nil), starts...)
	processed := make(map[int]bool, len(g.bombs))
	order := make([]int, 0, len(g.bombs))
	for len(queue) > 0 {
		index := queue[0]
		queue = queue[1:]
		if processed[index] || g.bombs[index].state == blasting {
			continue
		}
		processed[index] = true
		order = append(order, index)
		g.bombs[index].state = blasting
		g.bombs[index].timer = blastLife
		g.bombs[index].blast = blastCells(g.bombs[index].at)
		g.bombs[index].chainOrder = len(order)
		for other := range g.bombs {
			if processed[other] || g.bombs[other].state == blasting {
				continue
			}
			if contains(g.bombs[index].blast, g.bombs[other].at) {
				queue = append(queue, other)
			}
		}
	}
	g.lastChain = len(order)
	if g.lastChain > g.bestChain {
		g.bestChain = g.lastChain
	}
	g.queueLog = "queue "
	for _, index := range order {
		g.queueLog += fmt.Sprintf("B%d -> ", index+1)
	}
	g.queueLog += "done"
	g.message = fmt.Sprintf("Queue propagated %d unique explosion event(s).", g.lastChain)
}

func contains(cells []point, target point) bool {
	for _, p := range cells {
		if p == target {
			return true
		}
	}
	return false
}

func (g *game) bombAt(p point) bool {
	for _, b := range g.bombs {
		if b.at == p {
			return true
		}
	}
	return false
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 18, 34, 255})
	seconds := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, "CHAIN EXPLOSION", 176, 16)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("RELAYS %d/%d  BEST CHAIN %d/3  TIME %02d", g.cleared, goal, g.bestChain, seconds), 91, 41)
	ebitenutil.DebugPrintAt(screen, g.message, 31, 68)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := point{x, y}
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			c := color.RGBA{27, 47, 61, 255}
			if hardWall(p) {
				c = color.RGBA{55, 79, 116, 255}
			}
			vector.DrawFilledRect(screen, px+1, py+1, cell-2, cell-2, c, false)
			if g.targets[p] {
				vector.StrokeRect(screen, px+11, py+11, cell-22, cell-22, 3, color.RGBA{250, 202, 68, 255}, false)
			}
		}
	}
	for i := range g.bombs {
		b := &g.bombs[i]
		if b.state == blasting {
			for _, p := range b.blast {
				px, py := float32(boardX+p.x*cell), float32(boardY+p.y*cell)
				vector.DrawFilledRect(screen, px+4, py+4, cell-8, cell-8, color.RGBA{255, 140, 43, 220}, false)
			}
		}
		x, y := center(b.at)
		if b.state == waiting {
			vector.DrawFilledCircle(screen, x, y, 14, color.RGBA{25, 29, 39, 255}, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("B%d %.1f", i+1, float64(b.timer)/60), int(x)-20, int(y)-5)
		} else {
			vector.DrawFilledCircle(screen, x, y, 12, color.RGBA{255, 235, 120, 255}, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("#%d", b.chainOrder), int(x)-8, int(y)-5)
		}
	}
	px, py := center(g.player)
	vector.DrawFilledCircle(screen, px, py, 13, color.RGBA{229, 88, 76, 255}, false)
	ebitenutil.DebugPrintAt(screen, g.queueLog, 35, 550)
	ebitenutil.DebugPrintAt(screen, "Goal: clear 5 relays AND trigger a 3-bomb chain", 57, 572)
	drawControls(screen)
	if g.won {
		drawOverlay(screen, "CHAIN MASTER!", g.message)
	} else if g.lost {
		drawOverlay(screen, "CHAIN CAUGHT YOU", g.message)
	}
}

func center(p point) (float32, float32) {
	return float32(boardX + p.x*cell + cell/2), float32(boardY + p.y*cell + cell/2)
}

func directionInput() (point, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return point{-1, 0}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return point{0, -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return point{0, 1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		return point{1, 0}, true
	}
	if x, y, ok := justPressed(); ok && y >= 600 && x < 384 {
		return rays[min(3, x/96)], true
	}
	return point{}, false
}

func placePressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		return true
	}
	x, y, ok := justPressed()
	return ok && y >= 600 && x >= 384
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
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func drawControls(screen *ebiten.Image) {
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT", "BOMB"}
	for i, label := range labels {
		vector.DrawFilledRect(screen, float32(i*96+3), 600, 90, 62, color.RGBA{50, 82, 120, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*96+27, 627)
	}
	ebitenutil.DebugPrintAt(screen, "Arrows/WASD + Space, or tap controls", 105, 688)
}

func drawOverlay(screen *ebiten.Image, title, message string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	ebitenutil.DebugPrintAt(screen, title, 175, 304)
	ebitenutil.DebugPrintAt(screen, message, 59, 336)
	ebitenutil.DebugPrintAt(screen, "TAP / ENTER / R TO RETRY", 132, 380)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Chain Explosion - Ebi Showcase")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
