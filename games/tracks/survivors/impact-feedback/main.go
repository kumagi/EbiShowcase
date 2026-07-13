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
	x, y, hp float64
	flash    int
}
type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	px, py             float64
	mobs               []mob
	sparks             []spark
	tick, kills, shake int
	rng                *rand.Rand
}

func newGame() *game { return &game{px: 240, py: 360, rng: rand.New(rand.NewSource(44))} }
func (g *game) burst(x, y float64) {
	for i := 0; i < 10; i++ {
		a := float64(i) * math.Pi / 5
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * (1 + float64(i%3)), math.Sin(a) * (1 + float64(i%3)), 28})
	}
}
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
		g.px += dx / d * 3.8
		g.py += dy / d * 3.8
	}
	g.px = math.Max(20, math.Min(460, g.px))
	g.py = math.Max(90, math.Min(690, g.py))
	if g.tick%42 == 0 {
		a := g.rng.Float64() * math.Pi * 2
		g.mobs = append(g.mobs, mob{g.px + math.Cos(a)*300, g.py + math.Sin(a)*300, 2, 0})
	}
	for i := len(g.mobs) - 1; i >= 0; i-- {
		m := &g.mobs[i]
		d := math.Hypot(g.px-m.x, g.py-m.y)
		if d > 0 {
			m.x += (g.px - m.x) / d
			m.y += (g.py - m.y) / d
		}
		if m.flash > 0 {
			m.flash--
		}
		if d < 68 && g.tick%15 == 0 {
			m.hp--
			m.flash = 6
			g.burst(m.x, m.y)
			if m.hp <= 0 {
				g.kills++
				g.shake = 5
				g.mobs = append(g.mobs[:i], g.mobs[i+1:]...)
			}
		}
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.shake > 0 {
		g.shake--
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{14, 30, 50, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 5
	}
	for _, m := range g.mobs {
		if m.flash > 0 {
			trackatlas.DrawTinted(s, "swarm", m.x+ox, m.y, 30, 1, 1, .3, 1)
		} else {
			trackatlas.DrawCentered(s, "swarm", m.x+ox, m.y, 30)
		}
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 214, 66, 255}, true)
	}
	pulse := 64 + math.Sin(float64(g.tick)*.3)*5
	vector.StrokeCircle(s, float32(g.px+ox), float32(g.py), float32(pulse), 3, color.RGBA{255, 211, 62, 190}, true)
	hero.DrawCentered(s, g.px+ox, g.py, 36)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("KILLS %03d", g.kills), 205, 25)
	ebitenutil.DebugPrintAt(s, "HIT FLASH -> PARTICLES -> SHAKE", 120, 52)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / ARROWS / DRAG", 130, 685)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Impact Feedback Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
