package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float64 }
type enemy struct {
	x, y, vx, vy float64
	state, timer int
	alive        bool
}
type game struct {
	player           rect
	walls            []rect
	enemies          []enemy
	facingX, facingY float64
	attack, life     int
	clear, over      bool
}

func newGame() *game {
	g := &game{player: rect{48, 610, 28, 28}, facingY: -1, life: 5}
	g.walls = []rect{{0, 72, 480, 24}, {0, 696, 480, 24}, {0, 72, 24, 648}, {456, 72, 24, 648}, {110, 170, 260, 24}, {110, 170, 24, 150}, {250, 265, 206, 24}, {345, 265, 24, 165}, {24, 405, 230, 24}, {110, 520, 260, 24}, {235, 405, 24, 115}}
	g.enemies = []enemy{{180, 115, 0, 0, 0, 90, true}, {405, 220, 0, 0, 0, 30, true}, {70, 350, 0, 0, 0, 150, true}, {400, 470, 0, 0, 0, 60, true}, {175, 600, 0, 0, 0, 120, true}}
	return g
}
func (g *game) Update() error {
	if g.clear || g.over {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		dy++
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y > height/2 {
			dx = float64(x) - g.player.x - g.player.w/2
			dy = float64(y) - g.player.y - g.player.h/2
		}
	}
	if dx != 0 || dy != 0 {
		l := math.Hypot(dx, dy)
		dx, dy = dx/l*3.3, dy/l*3.3
		g.facingX, g.facingY = dx/3.3, dy/3.3
	}
	g.movePlayer(dx, 0)
	g.movePlayer(0, dy)
	if g.attack > 0 {
		g.attack--
	}
	attackPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y <= height/2 {
			attackPressed = true
		}
	}
	if attackPressed && g.attack == 0 {
		g.attack = 18
	}
	for i := range g.enemies {
		e := &g.enemies[i]
		if !e.alive {
			continue
		}
		dist := math.Hypot((g.player.x+14)-e.x, (g.player.y+14)-e.y)
		if dist < 165 {
			e.state = 1
		} else if dist > 230 {
			e.state = 0
		}
		if e.state == 0 {
			e.timer--
			if e.timer <= 0 {
				e.vx, e.vy = -e.vy, e.vx
				if e.vx == 0 && e.vy == 0 {
					e.vx = 1.1
				}
				e.timer = 100
			}
		} else {
			e.vx = ((g.player.x + 14) - e.x) / dist * 1.35
			e.vy = ((g.player.y + 14) - e.y) / dist * 1.35
		}
		g.moveEnemy(e, e.vx, 0)
		g.moveEnemy(e, 0, e.vy)
		if g.attack > 8 && math.Hypot((g.player.x+14+g.facingX*30)-e.x, (g.player.y+14+g.facingY*30)-e.y) < 35 {
			e.alive = false
		}
		if dist < 22 && g.attack == 0 {
			g.life--
			g.attack = 55
			if g.life <= 0 {
				g.over = true
			}
		}
	}
	if g.aliveCount() == 0 && g.player.y < 115 && g.player.x > 390 {
		g.clear = true
	}
	return nil
}
func (g *game) movePlayer(dx, dy float64) {
	g.player.x += dx
	g.player.y += dy
	for _, w := range g.walls {
		if overlap(g.player, w) {
			g.player.x -= dx
			g.player.y -= dy
			return
		}
	}
}
func (g *game) moveEnemy(e *enemy, dx, dy float64) {
	r := rect{e.x - 13 + dx, e.y - 13 + dy, 26, 26}
	for _, w := range g.walls {
		if overlap(r, w) {
			e.vx, e.vy = -e.vy, e.vx
			return
		}
	}
	e.x += dx
	e.y += dy
}
func (g *game) aliveCount() int {
	n := 0
	for _, e := range g.enemies {
		if e.alive {
			n++
		}
	}
	return n
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 18, 30, 255})
	for y := 96; y < 696; y += 32 {
		for x := 24; x < 456; x += 32 {
			c := color.RGBA{26, 39, 52, 255}
			if (x/32+y/32)%2 == 0 {
				c = color.RGBA{30, 44, 58, 255}
			}
			vector.DrawFilledRect(s, float32(x), float32(y), 32, 32, c, false)
		}
	}
	for _, w := range g.walls {
		vector.DrawFilledRect(s, float32(w.x), float32(w.y), float32(w.w), float32(w.h), color.RGBA{74, 78, 97, 255}, false)
		vector.StrokeRect(s, float32(w.x+3), float32(w.y+3), float32(w.w-6), float32(w.h-6), 2, color.RGBA{121, 123, 145, 255}, false)
	}
	exitColor := color.RGBA{115, 45, 61, 255}
	if g.aliveCount() == 0 {
		exitColor = color.RGBA{53, 210, 153, 255}
	}
	vector.DrawFilledRect(s, 397, 76, 55, 36, exitColor, false)
	for _, e := range g.enemies {
		if !e.alive {
			continue
		}
		c := color.RGBA{178, 86, 221, 255}
		if e.state == 1 {
			c = color.RGBA{245, 76, 98, 255}
		}
		vector.DrawFilledCircle(s, float32(e.x), float32(e.y), 13, c, false)
		vector.DrawFilledCircle(s, float32(e.x-5), float32(e.y-3), 2, color.White, false)
		vector.DrawFilledCircle(s, float32(e.x+5), float32(e.y-3), 2, color.White, false)
	}
	px, py := float32(g.player.x+14), float32(g.player.y+14)
	vector.DrawFilledCircle(s, px, py, 15, color.RGBA{45, 225, 195, 255}, false)
	vector.DrawFilledCircle(s, px+float32(g.facingX*7), py+float32(g.facingY*7), 3, color.RGBA{6, 31, 42, 255}, false)
	if g.attack > 8 {
		vector.StrokeLine(s, px+float32(g.facingX*18), py+float32(g.facingY*18), px+float32(g.facingX*48), py+float32(g.facingY*48), 7, color.RGBA{255, 222, 87, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LIFE %d     MONSTERS %d", g.life, g.aliveCount()), 18, 24)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / LOWER TOUCH    SWORD: SPACE / UPPER TOUCH", 52, 675)
	if g.aliveCount() == 0 {
		ebitenutil.DebugPrintAt(s, "THE EXIT IS OPEN!", 190, 52)
	}
	if g.clear {
		overlay(s, "DUNGEON CLEAR!\n\nTAP / SPACE TO RETURN")
	} else if g.over {
		overlay(s, "YOU FELL...\n\nTAP / SPACE TO RETRY")
	}
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, msg, 145, 330)
}
func overlap(a, b rect) bool { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Top-down Dungeon — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
