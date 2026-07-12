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
	W    = 480
	H    = 720
	cols = 6
	rows = 10
	cell = 54
	ox   = 78
	oy   = 78
)

var dirs = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type game struct {
	b                       [rows][cols]bool
	x, y, dir, tick, placed int
	clear, over             bool
	msg                     string
}

func newGame() *game {
	g := &game{}
	g.spawn()
	g.msg = "Rotate the partner around its pivot; lock 8 pairs."
	return g
}
func (g *game) cells(x, y, d int) (int, int, int, int) { v := dirs[d]; return x, y, x + v[0], y + v[1] }
func (g *game) blocked(x, y, d int) bool {
	a, b, c, e := g.cells(x, y, d)
	return a < 0 || a >= cols || c < 0 || c >= cols || b < 0 || b >= rows || e < 0 || e >= rows || g.b[b][a] || g.b[e][c]
}
func (g *game) spawn() {
	g.x, g.y, g.dir, g.tick = 2, 1, 2, 0
	if g.blocked(g.x, g.y, g.dir) {
		g.over = true
	}
}
func (g *game) rotate() {
	nd := (g.dir + 1) % 4
	for _, k := range [][2]int{{0, 0}, {-1, 0}, {1, 0}, {0, -1}} {
		if !g.blocked(g.x+k[0], g.y+k[1], nd) {
			g.x += k[0]
			g.y += k[1]
			g.dir = nd
			g.msg = "Rotation accepted; wall correction tried in order."
			return
		}
	}
	g.msg = "Rotation blocked: every correction candidate collided."
}
func (g *game) down() {
	if !g.blocked(g.x, g.y+1, g.dir) {
		g.y++
		return
	}
	a, b, c, d := g.cells(g.x, g.y, g.dir)
	g.b[b][a], g.b[d][c] = true, true
	g.placed++
	if g.placed >= 8 {
		g.clear = true
		return
	}
	g.spawn()
}
func (g *game) Update() error {
	if g.clear || g.over {
		if retry() {
			*g = *newGame()
		}
		return nil
	}
	l := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	r := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	rot := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	drop := inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)
	if x, y, ok := press(); ok && y > 625 {
		if x < 120 {
			l = true
		} else if x < 240 {
			rot = true
		} else if x < 360 {
			drop = true
		} else {
			r = true
		}
	}
	if l && !g.blocked(g.x-1, g.y, g.dir) {
		g.x--
	}
	if r && !g.blocked(g.x+1, g.y, g.dir) {
		g.x++
	}
	if rot {
		g.rotate()
	}
	if drop {
		g.down()
		g.tick = 0
	}
	g.tick++
	if g.tick >= 42 {
		g.tick = 0
		g.down()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 22, 42, 255})
	ebitenutil.DebugPrintAt(s, "ORBIT PAIR WORKSHOP", 170, 20)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LOCKED %d/8  DIR %d", g.placed, g.dir), 180, 47)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.StrokeRect(s, px, py, cell-2, cell-2, 1, color.RGBA{55, 77, 105, 255}, false)
			if g.b[y][x] {
				vector.DrawFilledCircle(s, px+cell/2, py+cell/2, 20, color.RGBA{86, 155, 210, 255}, true)
			}
		}
	}
	if !g.clear && !g.over {
		a, b, c, d := g.cells(g.x, g.y, g.dir)
		vector.DrawFilledCircle(s, float32(ox+a*cell+cell/2), float32(oy+b*cell+cell/2), 20, color.RGBA{243, 178, 59, 255}, true)
		vector.DrawFilledCircle(s, float32(ox+c*cell+cell/2), float32(oy+d*cell+cell/2), 20, color.RGBA{220, 91, 113, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, g.msg, 48, 592)
	for i, t := range []string{"LEFT", "ROTATE", "DOWN", "RIGHT"} {
		vector.DrawFilledRect(s, float32(i*120), 625, 115, 62, color.RGBA{55, 85, 121, 255}, false)
		ebitenutil.DebugPrintAt(s, t, i*120+35, 650)
	}
	if g.clear || g.over {
		vector.DrawFilledRect(s, 45, 275, 390, 150, color.RGBA{5, 14, 29, 245}, false)
		m := "8 PAIRS LOCKED!"
		if g.over {
			m = "BOARD TOPPED OUT!"
		}
		ebitenutil.DebugPrintAt(s, m, 160, 320)
		ebitenutil.DebugPrintAt(s, "TAP / ENTER TO RETRY", 145, 380)
	}
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
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
