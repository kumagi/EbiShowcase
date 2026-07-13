package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
	"math/rand"
)

const width, height = 480, 720

type mob struct {
	x, y, speed, size float64
	boss              bool
}
type game struct {
	px, py     float64
	mobs       []mob
	tick, life int
	rng        *rand.Rand
	boss       bool
}

func newGame() *game { return &game{px: 240, py: 360, life: 5, rng: rand.New(rand.NewSource(88))} }
func (g *game) Update() error {
	g.tick++
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
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		dx = float64(x) - g.px
		dy = float64(y) - g.py
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		dx = float64(x) - g.px
		dy = float64(y) - g.py
	}
	if dx != 0 || dy != 0 {
		d := math.Hypot(dx, dy)
		g.px += dx / d * 4
		g.py += dy / d * 4
	}
	g.px = math.Max(20, math.Min(460, g.px))
	g.py = math.Max(90, math.Min(690, g.py))
	sec := g.tick / 60
	wave := min(3, sec/10+1)
	interval := []int{55, 36, 22}[wave-1]
	if g.tick%interval == 0 && sec < 30 {
		a := g.rng.Float64() * math.Pi * 2
		g.mobs = append(g.mobs, mob{g.px + math.Cos(a)*300, g.py + math.Sin(a)*300, []float64{.8, 1.15, 1.55}[wave-1], []float64{26, 34, 22}[wave-1], false})
	}
	if sec >= 30 && !g.boss {
		g.boss = true
		g.mobs = append(g.mobs, mob{240, 100, .7, 75, true})
	}
	for i := len(g.mobs) - 1; i >= 0; i-- {
		m := &g.mobs[i]
		d := math.Hypot(g.px-m.x, g.py-m.y)
		speed := m.speed
		if m.boss && g.tick%210 > 165 {
			speed = 3
		}
		if d > 0 {
			m.x += (g.px - m.x) / d * speed
			m.y += (g.py - m.y) / d * speed
		}
		if d < 55 && g.tick%50 == 0 {
			g.life--
			if !m.boss {
				g.mobs = append(g.mobs[:i], g.mobs[i+1:]...)
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	sec := g.tick / 60
	wave := min(3, sec/10+1)
	colors := []color.RGBA{{13, 38, 55, 255}, {45, 25, 65, 255}, {68, 29, 28, 255}}
	s.Fill(colors[wave-1])
	for _, m := range g.mobs {
		sprite := "swarm"
		if m.boss {
			sprite = "boss-crab"
			if g.tick%210 > 130 && g.tick%210 <= 165 {
				vector.StrokeCircle(s, float32(m.x), float32(m.y), float32(55+(g.tick%10)*3), 3, color.RGBA{255, 95, 75, 220}, true)
			}
		}
		trackatlas.DrawCentered(s, sprite, m.x, m.y, m.size)
	}
	hero.DrawCentered(s, g.px, g.py, 36)
	label := []string{"SCOUT WAVE", "ARMORED WAVE", "RUSH WAVE"}[wave-1]
	if sec >= 30 {
		label = "BOSS: RING MEANS DASH"
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DIRECTOR %02d SEC  LIFE %d", sec, g.life), 145, 25)
	ebitenutil.DebugPrintAt(s, label, 145, 52)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / ARROWS / DRAG", 130, 685)
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Wave Director Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
