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

type bullet struct {
	x, y, vx, vy, r float64
	enemy           bool
}
type game struct {
	px, py, vx, vy                  float64
	bullets                         []bullet
	frame, bossHP, life, invincible int
	clear, over                     bool
}

func newGame() *game { return &game{px: 240, py: 620, bossHP: 320, life: 5} }
func (g *game) Update() error {
	if g.clear || g.over {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frame++
	if g.invincible > 0 {
		g.invincible--
	}
	ax, ay := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		ax--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		ax++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		ay--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		ay++
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		ax = float64(x) - g.px
		ay = float64(y) - g.py
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		ax = float64(x) - g.px
		ay = float64(y) - g.py
	}
	if ax != 0 || ay != 0 {
		l := math.Hypot(ax, ay)
		g.vx += ax / l * .8
		g.vy += ay / l * .8
	}
	g.vx *= .8
	g.vy *= .8
	speed := math.Hypot(g.vx, g.vy)
	if speed > 5.2 {
		g.vx *= 5.2 / speed
		g.vy *= 5.2 / speed
	}
	g.px = clamp(g.px+g.vx, 18, 462)
	g.py = clamp(g.py+g.vy, 90, 690)
	if g.frame%7 == 0 {
		g.bullets = append(g.bullets, bullet{g.px - 7, g.py - 15, -.25, -10, 4, false}, bullet{g.px + 7, g.py - 15, .25, -10, 4, false})
	}
	if g.frame%38 == 0 {
		count := 18
		spin := float64(g.frame) * .017
		for i := 0; i < count; i++ {
			a := spin + float64(i)*math.Pi*2/float64(count)
			speed := 2.0 + float64((g.frame/38)%3)*.35
			g.bullets = append(g.bullets, bullet{240, 128, math.Cos(a) * speed, math.Sin(a) * speed, 6, true})
		}
	}
	if g.frame%91 == 0 {
		base := math.Atan2(g.py-128, g.px-240)
		for i := -3; i <= 3; i++ {
			a := base + float64(i)*.11
			g.bullets = append(g.bullets, bullet{240, 128, math.Cos(a) * 3.0, math.Sin(a) * 3.0, 7, true})
		}
	}
	for i := len(g.bullets) - 1; i >= 0; i-- {
		b := &g.bullets[i]
		b.x += b.vx
		b.y += b.vy
		remove := b.x < -30 || b.x > 510 || b.y < -30 || b.y > 750
		if !b.enemy && math.Hypot(b.x-240, b.y-128) < 48 {
			g.bossHP--
			remove = true
			if g.bossHP <= 0 {
				g.clear = true
			}
		}
		if b.enemy && g.invincible == 0 && math.Hypot(b.x-g.px, b.y-g.py) < b.r+4 {
			g.life--
			g.invincible = 90
			remove = true
			if g.life <= 0 {
				g.over = true
			}
		}
		if remove {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{5, 8, 24, 255})
	for i := 0; i < 36; i++ {
		x := float32((i*83 + g.frame/3) % width)
		y := float32((i*137 + g.frame) % height)
		vector.DrawFilledCircle(s, x, y, 1, color.RGBA{150, 170, 230, 170}, false)
	}
	vector.DrawFilledCircle(s, 240, 128, 42, color.RGBA{143, 55, 208, 255}, false)
	vector.StrokeCircle(s, 240, 128, 52, 4, color.RGBA{246, 80, 133, 220}, false)
	vector.DrawFilledCircle(s, 226, 122, 5, color.RGBA{255, 224, 90, 255}, false)
	vector.DrawFilledCircle(s, 254, 122, 5, color.RGBA{255, 224, 90, 255}, false)
	vector.DrawFilledRect(s, 55, 55, 370, 14, color.RGBA{44, 48, 72, 255}, false)
	vector.DrawFilledRect(s, 55, 55, float32(370*max(g.bossHP, 0)/320), 14, color.RGBA{247, 72, 113, 255}, false)
	for _, b := range g.bullets {
		c := color.RGBA{255, 218, 72, 255}
		if b.enemy {
			c = color.RGBA{255, 76, 145, 255}
			vector.StrokeCircle(s, float32(b.x), float32(b.y), float32(b.r+2), 1, color.RGBA{139, 105, 255, 220}, false)
		}
		vector.DrawFilledCircle(s, float32(b.x), float32(b.y), float32(b.r), c, false)
	}
	if g.invincible%10 < 5 {
		vector.DrawFilledRect(s, float32(g.px-5), float32(g.py-17), 10, 32, color.RGBA{48, 228, 201, 255}, false)
		vector.DrawFilledRect(s, float32(g.px-18), float32(g.py+4), 36, 8, color.RGBA{48, 228, 201, 255}, false)
		vector.DrawFilledCircle(s, float32(g.px), float32(g.py), 4, color.White, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BOSS %03d     LIFE %d     BULLETS %03d", max(g.bossHP, 0), g.life, len(g.bullets)), 70, 20)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / ARROWS / DRAG    AUTO FIRE", 96, 688)
	if g.clear {
		overlay(s, "BOSS DEFEATED!\n\nTAP / SPACE TO FIGHT AGAIN")
	} else if g.over {
		overlay(s, "CONTINUE?\n\nTAP / SPACE TO RETRY")
	}
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 12, 35, 240}, false)
	ebitenutil.DebugPrintAt(s, msg, 145, 330)
}
func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Bullet Hell Boss — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
