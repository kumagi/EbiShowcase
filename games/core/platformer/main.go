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
type game struct {
	player                rect
	vx, vy, camera        float64
	platforms             []rect
	coins                 []rect
	collected             int
	onGround, clear, over bool
}

func newGame() *game {
	g := &game{player: rect{60, 560, 28, 38}}
	g.platforms = []rect{
		{0, 640, 420, 80}, {470, 590, 230, 130}, {760, 640, 260, 80}, {1080, 565, 210, 155}, {1350, 640, 300, 80}, {1710, 570, 220, 150}, {1990, 640, 430, 80},
		{250, 510, 120, 24}, {590, 450, 110, 24}, {850, 500, 100, 24}, {1160, 415, 100, 24}, {1450, 490, 120, 24}, {1800, 420, 110, 24}, {2110, 500, 110, 24},
	}
	for _, p := range []rect{{290, 475, 14, 14}, {625, 415, 14, 14}, {890, 465, 14, 14}, {1198, 380, 14, 14}, {1500, 455, 14, 14}, {1845, 385, 14, 14}, {2155, 465, 14, 14}} {
		g.coins = append(g.coins, p)
	}
	return g
}

func (g *game) Update() error {
	if g.clear || g.over {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	left := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_, y := ebiten.CursorPosition()
		if y < height/2 {
			jump = true
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			jump = true
		}
	}
	if left {
		g.vx -= .65
	}
	if right {
		g.vx += .65
	}
	if !left && !right {
		g.vx *= .78
	}
	g.vx = clamp(g.vx, -6, 6)
	if jump && g.onGround {
		g.vy = -12.5
		g.onGround = false
	}
	g.vy = math.Min(g.vy+.62, 14)
	g.player.x += g.vx
	for _, p := range g.platforms {
		if overlap(g.player, p) {
			if g.vx > 0 {
				g.player.x = p.x - g.player.w
			} else if g.vx < 0 {
				g.player.x = p.x + p.w
			}
			g.vx = 0
		}
	}
	g.player.y += g.vy
	g.onGround = false
	for _, p := range g.platforms {
		if overlap(g.player, p) {
			if g.vy > 0 {
				g.player.y = p.y - g.player.h
				g.onGround = true
			} else if g.vy < 0 {
				g.player.y = p.y + p.h
			}
			g.vy = 0
		}
	}
	for i := len(g.coins) - 1; i >= 0; i-- {
		if overlap(g.player, g.coins[i]) {
			g.coins = append(g.coins[:i], g.coins[i+1:]...)
			g.collected++
		}
	}
	if g.player.y > height+100 {
		g.over = true
	}
	if g.player.x > 2320 {
		g.clear = true
	}
	target := g.player.x - width*.38
	g.camera += (target - g.camera) * .09
	g.camera = clamp(g.camera, 0, 1940)
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{91, 184, 230, 255})
	for i := 0; i < 9; i++ {
		x := float32(i*320) - float32(math.Mod(g.camera*.25, 320))
		vector.DrawFilledCircle(s, x, 560, 150, color.RGBA{71, 145, 155, 255}, false)
	}
	for _, p := range g.platforms {
		x := float32(p.x - g.camera)
		vector.DrawFilledRect(s, x, float32(p.y), float32(p.w), float32(p.h), color.RGBA{51, 99, 68, 255}, false)
		vector.DrawFilledRect(s, x, float32(p.y), float32(p.w), 10, color.RGBA{98, 205, 91, 255}, false)
	}
	for _, c := range g.coins {
		x, y := float32(c.x-g.camera+c.w/2), float32(c.y+c.h/2)
		vector.DrawFilledCircle(s, x, y, 9, color.RGBA{255, 215, 60, 255}, false)
		vector.StrokeCircle(s, x, y, 12, 2, color.RGBA{255, 244, 170, 255}, false)
	}
	px, py := float32(g.player.x-g.camera), float32(g.player.y)
	vector.DrawFilledRect(s, px, py, float32(g.player.w), float32(g.player.h), color.RGBA{240, 73, 89, 255}, false)
	vector.DrawFilledRect(s, px-3, py+30, 34, 8, color.RGBA{30, 57, 86, 255}, false)
	vector.DrawFilledCircle(s, px+20, py+11, 3, color.White, false)
	flagX := float32(2340 - g.camera)
	vector.DrawFilledRect(s, flagX, 440, 6, 200, color.RGBA{235, 240, 245, 255}, false)
	vector.DrawFilledRect(s, flagX+6, 450, 70, 42, color.RGBA{255, 211, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("COINS %d/7", g.collected), 18, 18)
	ebitenutil.DebugPrintAt(s, "RUN: A/D OR BOTTOM TOUCH    JUMP: SPACE OR TOP TOUCH", 56, 680)
	if g.clear {
		overlay(s, "COURSE CLEAR!\n\nTAP / SPACE TO RUN AGAIN")
	} else if g.over {
		overlay(s, "MISS!\n\nTAP / SPACE TO RETRY")
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, msg, 145, 330)
}
func overlap(a, b rect) bool          { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Platformer — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
