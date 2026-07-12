package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"image/color"
	"math"
	"math/rand"
)

const width, height = 480, 720

type mob struct {
	x, y   float64
	active bool
}
type game struct {
	px, py                  float64
	mobs                    [100]mob
	rng                     *rand.Rand
	frame, kills, life, inv int
	clear, over             bool
}

func newGame() *game {
	g := &game{px: 240, py: 360, rng: rand.New(rand.NewSource(2003)), life: 5}
	for i := range g.mobs {
		g.spawn(i)
	}
	return g
}
func (g *game) spawn(i int) {
	a := g.rng.Float64() * math.Pi * 2
	r := 280 + g.rng.Float64()*120
	g.mobs[i] = mob{g.px + math.Cos(a)*r, g.py + math.Sin(a)*r, true}
}
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
	if dx != 0 || dy != 0 {
		l := math.Hypot(dx, dy)
		g.px += dx / l * 4
		g.py += dy / l * 4
	}
	g.px = clamp(g.px, 20, 460)
	g.py = clamp(g.py, 90, 690)
	for i := range g.mobs {
		m := &g.mobs[i]
		d := math.Hypot(g.px-m.x, g.py-m.y)
		m.x += (g.px - m.x) / d * .72
		m.y += (g.py - m.y) / d * .72
		if d < 62 && g.frame%18 == 0 {
			g.kills++
			g.spawn(i)
			continue
		}
		if d < 22 && g.inv == 0 {
			g.life--
			g.inv = 90
			g.spawn(i)
			if g.life <= 0 {
				g.over = true
			}
		}
	}
	if g.kills >= 100 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{11, 25, 40, 255})
	for _, m := range g.mobs {
		vector.DrawFilledCircle(s, float32(m.x), float32(m.y), 9, color.RGBA{220, 70, 104, 230}, false)
	}
	vector.DrawFilledCircle(s, float32(g.px), float32(g.py), 62, color.RGBA{255, 211, 61, 35}, false)
	vector.StrokeCircle(s, float32(g.px), float32(g.py), 62, 2, color.RGBA{255, 211, 61, 180}, false)
	if g.inv%10 < 5 {
		hero.DrawCentered(s, g.px, g.py, 34)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("POOL 100/100   DEFEATED %03d/100   LIFE %d", g.kills, g.life), 52, 24)
	ebitenutil.DebugPrintAt(s, "THE GOLD AURA ATTACKS EVERY 18 FRAMES", 90, 50)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / ARROWS / TOUCH", 120, 685)
	if g.clear {
		overlay(s, "100 ENEMIES DEFEATED!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "SWARMED!\n\nTAP / SPACE TO RETRY")
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
	ebiten.SetWindowTitle("100 Enemy Swarm — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
