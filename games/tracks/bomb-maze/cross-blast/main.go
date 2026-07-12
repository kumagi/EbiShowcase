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
	fuse       = 150
	blastLife  = 32
	blastRange = 3
	goal       = 6
	timeLimit  = 60 * 60
)

type point struct{ x, y int }

var rays = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

type bomb struct {
	at    point
	timer int
	blast []point
}

type game struct {
	player          point
	bombs           []bomb
	targets         map[point]bool
	cleared, frames int
	message         string
	won, lost       bool
}

func newGame() *game {
	g := &game{
		player:  point{1, 1},
		targets: map[point]bool{},
		message: "A blast walks one cell at a time and stops at hard walls.",
	}
	for _, p := range []point{{4, 1}, {7, 1}, {1, 3}, {5, 3}, {3, 5}, {7, 7}, {1, 7}, {5, 7}} {
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
		g.message = "Time up. Use walls to control each blast ray."
		return nil
	}
	if direction, ok := directionInput(); ok {
		next := point{g.player.x + direction.x, g.player.y + direction.y}
		if !hardWall(next) && !g.bombAt(next) {
			g.player = next
		}
	}
	if placePressed() {
		if len(g.bombs) == 0 {
			g.bombs = append(g.bombs, bomb{at: g.player, timer: fuse})
			g.message = "Bomb placed. The preview shows its four range-limited rays."
		} else {
			g.message = "One bomb at a time in this lesson. Wait for its blast."
		}
	}

	nextBombs := g.bombs[:0]
	for i := range g.bombs {
		b := g.bombs[i]
		b.timer--
		if b.timer == 0 {
			b.blast = blastCells(b.at)
			g.message = fmt.Sprintf("Ray traversal marked %d blast cells.", len(b.blast))
		}
		if b.timer <= 0 {
			for _, p := range b.blast {
				if p == g.player {
					g.lost = true
					g.message = "The player stood on a marked blast cell."
				}
				if g.targets[p] {
					delete(g.targets, p)
					g.cleared++
				}
			}
			if g.cleared >= goal && !g.lost {
				g.won = true
				g.message = "Six beacons cleared by aiming cross-shaped blasts!"
			}
		}
		if b.timer > -blastLife {
			nextBombs = append(nextBombs, b)
		}
	}
	g.bombs = nextBombs
	return nil
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
	ebitenutil.DebugPrintAt(screen, "CROSS BLAST", 196, 17)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BEACONS %d/%d   TIME %02d   RANGE %d", g.cleared, goal, seconds, blastRange), 118, 43)
	ebitenutil.DebugPrintAt(screen, g.message, 33, 72)
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
				vector.StrokeCircle(screen, px+cell/2, py+cell/2, 11, 3, color.RGBA{250, 202, 68, 255}, false)
			}
		}
	}
	for _, b := range g.bombs {
		cells := b.blast
		preview := b.timer > 0
		if preview {
			cells = blastCells(b.at)
		}
		for _, p := range cells {
			px, py := float32(boardX+p.x*cell), float32(boardY+p.y*cell)
			if preview {
				vector.StrokeRect(screen, px+5, py+5, cell-10, cell-10, 2, color.RGBA{255, 172, 69, 140}, false)
			} else {
				vector.DrawFilledRect(screen, px+4, py+4, cell-8, cell-8, color.RGBA{255, 145, 45, 220}, false)
				vector.DrawFilledCircle(screen, px+cell/2, py+cell/2, 9, color.RGBA{255, 236, 127, 255}, false)
			}
		}
		x, y := center(b.at)
		if preview {
			vector.DrawFilledCircle(screen, x, y, 14, color.RGBA{25, 29, 39, 255}, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f", float64(b.timer)/60), int(x)-10, int(y)-5)
		}
	}
	px, py := center(g.player)
	vector.DrawFilledCircle(screen, px, py, 13, color.RGBA{229, 88, 76, 255}, false)
	drawControls(screen)
	if g.won {
		drawOverlay(screen, "RAY MASTER!", g.message)
	}
	if g.lost {
		drawOverlay(screen, "BLASTED!", g.message)
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
	ebitenutil.DebugPrintAt(screen, title, 190, 304)
	ebitenutil.DebugPrintAt(screen, message, 68, 336)
	ebitenutil.DebugPrintAt(screen, "TAP / ENTER / R TO RETRY", 132, 380)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Cross Blast - Ebi Showcase")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
