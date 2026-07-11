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
	rect
	vx    float64
	alive bool
}
type game struct {
	player                rect
	vx, vy                float64
	grounds               []rect
	enemies               []enemy
	grounded, clear, over bool
	life                  int
}

func newGame() *game {
	return &game{player: rect{30, 590, 28, 38}, life: 3, grounds: []rect{{0, 650, 480, 70}, {70, 535, 130, 20}, {245, 445, 120, 20}, {355, 345, 110, 20}, {180, 255, 115, 20}, {25, 175, 110, 20}}, enemies: []enemy{{rect: rect{115, 507, 28, 28}, vx: 1.1, alive: true}, {rect: rect{285, 417, 28, 28}, vx: -1.25, alive: true}, {rect: rect{390, 317, 28, 28}, vx: 1.35, alive: true}, {rect: rect{220, 227, 28, 28}, vx: -1.45, alive: true}}}
}
func (g *game) Update() error {
	if g.clear || g.over {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	left := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				left = true
			} else {
				right = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			jump = true
		}
	}
	if left {
		g.vx -= .7
	}
	if right {
		g.vx += .7
	}
	if !left && !right {
		g.vx *= .78
	}
	g.vx = clamp(g.vx, -5.5, 5.5)
	if jump && g.grounded {
		g.vy = -12
		g.grounded = false
	}
	g.vy = math.Min(g.vy+.65, 14)
	g.player.x = clamp(g.player.x+g.vx, 0, width-g.player.w)
	oldBottom := g.player.y + g.player.h
	g.player.y += g.vy
	g.grounded = false
	for _, p := range g.grounds {
		if g.vy >= 0 && oldBottom <= p.y+3 && g.player.y+g.player.h >= p.y && g.player.x+g.player.w > p.x && g.player.x < p.x+p.w {
			g.player.y = p.y - g.player.h
			g.vy = 0
			g.grounded = true
		}
	}
	for i := range g.enemies {
		e := &g.enemies[i]
		if !e.alive {
			continue
		}
		e.x += e.vx
		support := false
		for _, p := range g.grounds {
			if e.x+e.w/2 >= p.x && e.x+e.w/2 <= p.x+p.w && math.Abs(e.y+e.h-p.y) < 3 {
				support = true
				if e.x <= p.x || e.x+e.w >= p.x+p.w {
					e.vx = -e.vx
				}
			}
		}
		if !support {
			e.vx = -e.vx
			e.x += e.vx * 2
		}
		if overlap(g.player, e.rect) {
			if g.vy > 1 && oldBottom <= e.y+8 {
				e.alive = false
				g.player.y = e.y - g.player.h
				g.vy = -8
			} else {
				g.life--
				g.player.x = 30
				g.player.y = 590
				g.vx, g.vy = 0, 0
				if g.life <= 0 {
					g.over = true
				}
			}
		}
	}
	if g.player.y > height {
		g.life--
		g.player.x, g.player.y = 30, 590
		if g.life <= 0 {
			g.over = true
		}
	}
	if g.player.y < 165 && g.player.x < 140 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{116, 195, 231, 255})
	for _, p := range g.grounds {
		vector.DrawFilledRect(s, float32(p.x), float32(p.y), float32(p.w), float32(p.h), color.RGBA{62, 106, 68, 255}, false)
		vector.DrawFilledRect(s, float32(p.x), float32(p.y), float32(p.w), 7, color.RGBA{113, 220, 89, 255}, false)
	}
	for _, e := range g.enemies {
		if !e.alive {
			continue
		}
		vector.DrawFilledCircle(s, float32(e.x+14), float32(e.y+14), 14, color.RGBA{143, 75, 194, 255}, false)
		vector.DrawFilledCircle(s, float32(e.x+9), float32(e.y+10), 2, color.White, false)
		vector.DrawFilledCircle(s, float32(e.x+19), float32(e.y+10), 2, color.White, false)
	}
	vector.DrawFilledRect(s, 62, 135, 55, 35, color.RGBA{255, 211, 62, 255}, false)
	vector.DrawFilledRect(s, float32(g.player.x), float32(g.player.y), float32(g.player.w), float32(g.player.h), color.RGBA{241, 72, 88, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LIFE %d   ENEMIES %d", g.life, g.alive()), 18, 22)
	ebitenutil.DebugPrintAt(s, "JUMP ON ENEMIES FROM ABOVE", 145, 48)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH    JUMP: SPACE OR UPPER TOUCH", 50, 685)
	if g.clear {
		overlay(s, "GARDEN CLEAR!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "GAME OVER\n\nTAP / SPACE TO RETRY")
	}
}
func (g *game) alive() int {
	n := 0
	for _, e := range g.enemies {
		if e.alive {
			n++
		}
	}
	return n
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, msg, 145, 330)
}
func overlap(a, b rect) bool        { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Patrol Enemies — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
