package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

const width, height = 480, 720

type body struct{ x, y, vx, vy, r float64 }
type game struct {
	px, py                  float64
	enemies, shots          []body
	rng                     *rand.Rand
	frame, life, score, inv int
	clear, over             bool
}

func newGame() *game { return &game{px: 240, py: 360, rng: rand.New(rand.NewSource(1902)), life: 5} }
func (g *game) Update() error {
	if g.clear || g.over {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	g.frame++
	if g.inv > 0 {
		g.inv--
	}
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy++
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		dx = float64(x) - g.px
		dy = float64(y) - g.py
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		dx = float64(x) - g.px
		dy = float64(y) - g.py
	}
	if dx != 0 || dy != 0 {
		l := math.Hypot(dx, dy)
		g.px += dx / l * 3.8
		g.py += dy / l * 3.8
	}
	g.px = clamp(g.px, 20, 460)
	g.py = clamp(g.py, 90, 690)
	if g.frame%42 == 0 {
		a := g.rng.Float64() * math.Pi * 2
		g.enemies = append(g.enemies, body{g.px + math.Cos(a)*360, g.py + math.Sin(a)*360, 0, 0, 14})
	}
	if g.frame%28 == 0 && len(g.enemies) > 0 {
		best, dist := 0, math.MaxFloat64
		for i, e := range g.enemies {
			d := math.Hypot(e.x-g.px, e.y-g.py)
			if d < dist {
				best, dist = i, d
			}
		}
		e := g.enemies[best]
		g.shots = append(g.shots, body{g.px, g.py, (e.x - g.px) / dist * 7, (e.y - g.py) / dist * 7, 5})
	}
	for i := range g.enemies {
		e := &g.enemies[i]
		d := math.Hypot(g.px-e.x, g.py-e.y)
		e.x += (g.px - e.x) / d * 1.05
		e.y += (g.py - e.y) / d * 1.05
		if d < 25 && g.inv == 0 {
			g.life--
			g.inv = 90
			if g.life <= 0 {
				g.over = true
			}
		}
	}
	for i := len(g.shots) - 1; i >= 0; i-- {
		s := &g.shots[i]
		s.x += s.vx
		s.y += s.vy
		remove := s.x < 0 || s.x > width || s.y < 70 || s.y > height
		for j := len(g.enemies) - 1; j >= 0 && !remove; j-- {
			if math.Hypot(s.x-g.enemies[j].x, s.y-g.enemies[j].y) < s.r+g.enemies[j].r {
				g.enemies = append(g.enemies[:j], g.enemies[j+1:]...)
				g.score++
				remove = true
			}
		}
		if remove {
			g.shots = append(g.shots[:i], g.shots[i+1:]...)
		}
	}
	if g.score >= 30 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 27, 43, 255})
	for _, e := range g.enemies {
		vector.DrawFilledCircle(s, float32(e.x), float32(e.y), 14, color.RGBA{229, 73, 106, 255}, false)
	}
	for _, b := range g.shots {
		vector.DrawFilledCircle(s, float32(b.x), float32(b.y), 5, color.RGBA{255, 216, 62, 255}, false)
	}
	if g.inv%10 < 5 {
		vector.DrawFilledCircle(s, float32(g.px), float32(g.py), 17, color.RGBA{45, 225, 194, 255}, false)
		vector.StrokeCircle(s, float32(g.px), float32(g.py), 25, 3, color.RGBA{255, 215, 63, 220}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TARGETS %02d/30   LIFE %d   COOLDOWN %02d", g.score, g.life, 28-g.frame%28), 65, 25)
	ebitenutil.DebugPrintAt(s, "MOVE ONLY — THE TURRET FIRES BY ITSELF", 90, 685)
	if g.clear {
		overlay(s, "30 TARGETS DEFEATED!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "OVERRUN!\n\nTAP / SPACE TO RETRY")
	}
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Auto Turret Survival — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
