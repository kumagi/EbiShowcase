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

type dot struct{ x, y float64 }
type game struct {
	px, py                            float64
	enemies, gems                     []dot
	rng                               *rand.Rand
	frame, level, xp, need, life, inv int
	speed, radius                     float64
	draft, clear, over                bool
}

func newGame() *game {
	return &game{px: 240, py: 360, rng: rand.New(rand.NewSource(2104)), level: 1, need: 5, life: 5, speed: 3.8, radius: 55}
}
func (g *game) Update() error {
	if g.clear || g.over {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	if g.draft {
		return g.updateDraft()
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
		g.px += dx / l * g.speed
		g.py += dy / l * g.speed
	}
	g.px = clamp(g.px, 20, 460)
	g.py = clamp(g.py, 90, 690)
	if g.frame%35 == 0 {
		a := g.rng.Float64() * math.Pi * 2
		g.enemies = append(g.enemies, dot{g.px + math.Cos(a)*330, g.py + math.Sin(a)*330})
	}
	for i := len(g.enemies) - 1; i >= 0; i-- {
		e := &g.enemies[i]
		d := math.Hypot(g.px-e.x, g.py-e.y)
		e.x += (g.px - e.x) / d
		e.y += (g.py - e.y) / d
		if d < g.radius && g.frame%15 == 0 {
			g.gems = append(g.gems, *e)
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
		} else if d < 22 && g.inv == 0 {
			g.life--
			g.inv = 90
			if g.life <= 0 {
				g.over = true
			}
		}
	}
	for i := len(g.gems) - 1; i >= 0; i-- {
		d := math.Hypot(g.px-g.gems[i].x, g.py-g.gems[i].y)
		if d < 90 {
			g.gems[i].x += (g.px - g.gems[i].x) / d * 5
			g.gems[i].y += (g.py - g.gems[i].y) / d * 5
		}
		if d < 18 {
			g.xp++
			g.gems = append(g.gems[:i], g.gems[i+1:]...)
			if g.xp >= g.need {
				g.draft = true
				break
			}
		}
	}
	return nil
}
func (g *game) updateDraft() error {
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		choice = 2
	}
	if x, _, ok := press(); ok {
		choice = min(2, x/(width/3))
	}
	if choice < 0 {
		return nil
	}
	switch choice {
	case 0:
		g.speed += .65
	case 1:
		g.radius += 14
	case 2:
		g.life++
	}
	g.level++
	g.xp = 0
	g.need += 3
	g.draft = false
	if g.level >= 5 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 26, 42, 255})
	for _, e := range g.enemies {
		vector.DrawFilledCircle(s, float32(e.x), float32(e.y), 12, color.RGBA{228, 72, 105, 255}, false)
	}
	for _, d := range g.gems {
		vector.DrawFilledCircle(s, float32(d.x), float32(d.y), 6, color.RGBA{95, 171, 255, 255}, false)
	}
	vector.DrawFilledCircle(s, float32(g.px), float32(g.py), float32(g.radius), color.RGBA{255, 211, 62, 28}, false)
	if g.inv%10 < 5 {
		vector.DrawFilledCircle(s, float32(g.px), float32(g.py), 15, color.RGBA{45, 225, 194, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LEVEL %d/5   XP %02d/%02d   LIFE %d", g.level, g.xp, g.need, g.life), 90, 24)
	ebitenutil.DebugPrintAt(s, "COLLECT BLUE EXPERIENCE GEMS", 125, 685)
	if g.draft {
		vector.DrawFilledRect(s, 24, 245, 432, 220, color.RGBA{6, 18, 37, 245}, false)
		ebitenutil.DebugPrintAt(s, "LEVEL UP! CHOOSE 1 / 2 / 3", 125, 270)
		cards := []string{"1  SPEED\n\nMOVE FASTER", "2  AURA\n\nWIDER ATTACK", "3  HEART\n\n+1 LIFE"}
		for i, t := range cards {
			x := float32(35 + i*140)
			vector.StrokeRect(s, x, 310, 130, 125, 2, color.RGBA{45, 225, 194, 255}, false)
			ebitenutil.DebugPrintAt(s, t, int(x)+18, 330)
		}
	}
	if g.clear {
		overlay(s, "LEVEL 5 REACHED!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "DEFEATED!\n\nTAP / SPACE TO RETRY")
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
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, msg, 140, 330)
}
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Experience Draft — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
