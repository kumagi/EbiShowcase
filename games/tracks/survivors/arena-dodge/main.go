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

type mob struct{ x, y float64 }
type game struct {
	px, py           float64
	mobs             []mob
	rng              *rand.Rand
	frame, life, inv int
	clear, over      bool
}

func newGame() *game { return &game{px: 240, py: 360, rng: rand.New(rand.NewSource(1801)), life: 5} }
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
		g.px += dx / l * 4
		g.py += dy / l * 4
	}
	g.px = clamp(g.px, 20, 460)
	g.py = clamp(g.py, 90, 690)
	if g.frame%45 == 0 {
		side := g.rng.Intn(4)
		x, y := 0.0, 0.0
		switch side {
		case 0:
			x = float64(g.rng.Intn(width))
			y = 75
		case 1:
			x = 470
			y = float64(80 + g.rng.Intn(620))
		case 2:
			x = float64(g.rng.Intn(width))
			y = 715
		case 3:
			x = 10
			y = float64(80 + g.rng.Intn(620))
		}
		g.mobs = append(g.mobs, mob{x, y})
	}
	for i := range g.mobs {
		m := &g.mobs[i]
		d := math.Hypot(g.px-m.x, g.py-m.y)
		m.x += (g.px - m.x) / d * 1.15
		m.y += (g.py - m.y) / d * 1.15
		if d < 24 && g.inv == 0 {
			g.life--
			g.inv = 90
			if g.life <= 0 {
				g.over = true
			}
		}
	}
	if g.frame >= 60*30 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{13, 28, 44, 255})
	for x := 0; x < width; x += 40 {
		for y := 80; y < height; y += 40 {
			vector.StrokeRect(s, float32(x), float32(y), 40, 40, 1, color.RGBA{61, 92, 111, 90}, false)
		}
	}
	for _, m := range g.mobs {
		vector.DrawFilledCircle(s, float32(m.x), float32(m.y), 14, color.RGBA{226, 72, 105, 255}, false)
	}
	if g.inv%10 < 5 {
		vector.DrawFilledCircle(s, float32(g.px), float32(g.py), 16, color.RGBA{45, 225, 194, 255}, false)
		vector.DrawFilledCircle(s, float32(g.px+6), float32(g.py-4), 3, color.RGBA{5, 27, 38, 255}, false)
	}
	left := max(0, 30-g.frame/60)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SURVIVE %02d   LIFE %d   ENEMIES %02d", left, g.life, len(g.mobs)), 90, 25)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / ARROWS / DRAG", 125, 685)
	if g.clear {
		overlay(s, "30 SECONDS SURVIVED!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "CAUGHT!\n\nTAP / SPACE TO RETRY")
	}
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, msg, 135, 330)
}
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Arena Dodge — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
