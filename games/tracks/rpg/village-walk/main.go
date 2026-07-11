package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const width, height, tile = 480, 720, 48

type point struct{ x, y int }
type game struct {
	p, facing point
	walls     map[point]bool
	stars     map[point]bool
	collected int
	clear     bool
}

func newGame() *game {
	g := &game{p: point{1, 11}, facing: point{0, -1}, walls: map[point]bool{}, stars: map[point]bool{{2, 2}: true, {7, 4}: true, {4, 9}: true}}
	rows := []string{"##########", "#   ##   #", "#        #", "# ##  ## #", "#        #", "#  ####  #", "#        #", "# ##  ## #", "#        #", "#   ##   #", "#        #", "#        #", "##########"}
	for y, row := range rows {
		for x, c := range row {
			if c == '#' {
				g.walls[point{x, y}] = true
			}
		}
	}
	return g
}
func (g *game) Update() error {
	if g.clear {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	d := point{}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		d.x = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		d.x = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		d.y = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		d.y = 1
	}
	if x, y, ok := press(); ok {
		dx, dy := x-(g.p.x*tile+tile/2), y-(74+g.p.y*tile+tile/2)
		if abs(dx) > abs(dy) {
			if dx < 0 {
				d.x = -1
			} else {
				d.x = 1
			}
		} else {
			if dy < 0 {
				d.y = -1
			} else {
				d.y = 1
			}
		}
	}
	if d.x != 0 || d.y != 0 {
		g.facing = d
		n := point{g.p.x + d.x, g.p.y + d.y}
		if !g.walls[n] {
			g.p = n
			if g.stars[n] {
				delete(g.stars, n)
				g.collected++
			}
		}
	}
	if g.p == (point{8, 1}) && g.collected == 3 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{11, 28, 39, 255})
	oy := 74
	for y := 0; y < 13; y++ {
		for x := 0; x < 10; x++ {
			p := point{x, y}
			c := color.RGBA{86, 164, 91, 255}
			if (x+y)%2 == 0 {
				c = color.RGBA{94, 174, 96, 255}
			}
			vector.DrawFilledRect(s, float32(x*tile), float32(oy+y*tile), tile, tile, c, false)
			if g.walls[p] {
				vector.DrawFilledRect(s, float32(x*tile+3), float32(oy+y*tile+3), tile-6, tile-6, color.RGBA{65, 91, 112, 255}, false)
			}
		}
	}
	for p := range g.stars {
		vector.DrawFilledCircle(s, float32(p.x*tile+24), float32(oy+p.y*tile+24), 10, color.RGBA{255, 215, 62, 255}, false)
	}
	gate := color.RGBA{123, 70, 78, 255}
	if g.collected == 3 {
		gate = color.RGBA{45, 225, 194, 255}
	}
	vector.DrawFilledRect(s, 8*tile+8, float32(oy+1*tile+8), 32, 32, gate, false)
	px, py := float32(g.p.x*tile+24), float32(oy+g.p.y*tile+24)
	vector.DrawFilledCircle(s, px, py, 15, color.RGBA{240, 74, 90, 255}, false)
	vector.DrawFilledCircle(s, px+float32(g.facing.x*8), py+float32(g.facing.y*8), 3, color.White, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("VILLAGE STARS %d/3", g.collected), 170, 25)
	ebitenutil.DebugPrintAt(s, "COLLECT 3 STARS, THEN REACH THE GREEN SHRINE", 85, 50)
	ebitenutil.DebugPrintAt(s, "ARROWS / WASD / TAP A DIRECTION", 125, 690)
	if g.clear {
		overlay(s, "VILLAGE EXPLORED!\n\nTAP / SPACE TO WALK AGAIN")
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
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 240}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Village Walk — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
