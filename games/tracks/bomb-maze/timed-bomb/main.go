package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	screenW   = 480
	screenH   = 720
	cols      = 9
	rows      = 9
	cell      = 48
	boardX    = 24
	boardY    = 100
	fuse      = 150
	blastLife = 28
)

type point struct{ x, y int }
type bomb struct {
	at       point
	timer    int
	blasting bool
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
	g := &game{player: point{1, 1}, targets: map[point]bool{}, message: "Move, place one bomb, then leave its cell before the timer ends."}
	for _, p := range []point{{3, 1}, {5, 2}, {7, 3}, {2, 5}, {4, 6}, {7, 7}} {
		g.targets[p] = true
	}
	return g
}
func wall(p point) bool {
	return p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || (p.x%2 == 0 && p.y%2 == 0)
}
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	d, ok := inputDir()
	if ok {
		n := point{g.player.x + d.x, g.player.y + d.y}
		if !wall(n) && !g.bombAt(n) {
			g.player = n
		}
	}
	if placePressed() {
		if !g.bombAt(g.player) {
			g.bombs = append(g.bombs, bomb{at: g.player, timer: fuse})
			g.message = "Bomb placed: 2.5 seconds. Step away!"
		} else {
			g.message = "Only one bomb can occupy a cell."
		}
	}
	next := g.bombs[:0]
	for i := range g.bombs {
		b := g.bombs[i]
		b.timer--
		if !b.blasting && b.timer <= 0 {
			b.blasting = true
			b.timer = blastLife
			g.message = "Timer reached zero: state changed to BLAST"
		}
		if b.blasting {
			if b.at == g.player {
				g.lost = true
				g.message = "The blast caught the player."
			}
			if g.targets[b.at] {
				delete(g.targets, b.at)
				g.cleared++
				if g.cleared >= 5 {
					g.won = true
					g.message = "Five marked cells cleared with timed bombs!"
				}
			}
		}
		if b.timer > 0 {
			next = append(next, b)
		}
	}
	g.bombs = next
	if g.frames >= 60*60 {
		g.lost = true
		g.message = "Time up."
	}
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
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{9, 18, 34, 255})
	ebitenutil.DebugPrintAt(s, "TIMED BOMB LAB", 184, 18)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TARGETS %d/5  TIME %02d", g.cleared, max(0, 60-g.frames/60)), 156, 45)
	ebitenutil.DebugPrintAt(s, g.message, 35, 72)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			c := color.RGBA{29, 48, 63, 255}
			if wall(point{x, y}) {
				c = color.RGBA{52, 78, 116, 255}
			}
			vector.DrawFilledRect(s, px+1, py+1, cell-2, cell-2, c, false)
			if g.targets[point{x, y}] {
				vector.StrokeRect(s, px+8, py+8, cell-16, cell-16, 3, color.RGBA{244, 189, 65, 255}, false)
			}
		}
	}
	for _, b := range g.bombs {
		x, y := float32(boardX+b.at.x*cell+cell/2), float32(boardY+b.at.y*cell+cell/2)
		if b.blasting {
			vector.DrawFilledCircle(s, x, y, 22, color.RGBA{255, 151, 53, 220}, false)
		} else {
			vector.DrawFilledCircle(s, x, y, 14, color.RGBA{28, 31, 40, 255}, false)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("%.1f", float64(b.timer)/60), int(x)-10, int(y)-5)
		}
	}
	px, py := float32(boardX+g.player.x*cell+cell/2), float32(boardY+g.player.y*cell+cell/2)
	vector.DrawFilledCircle(s, px, py, 13, color.RGBA{231, 89, 78, 255}, false)
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT", "BOMB"}
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*96+3), 600, 90, 62, color.RGBA{52, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*96+27, 627)
	}
	ebitenutil.DebugPrintAt(s, "Arrows/WASD + Space, or tap controls", 105, 688)
	if g.won {
		overlay(s, "TIMER MASTER!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(s, "BLASTED / TIME UP\n\nTAP / ENTER TO RETRY")
	}
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
func retry() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	ebitenutil.DebugPrintAt(s, t, 112, 328)
}
func (g *game) Layout(int, int) (int, int) { return screenW, screenH }
func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Timed Bomb Lab")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
