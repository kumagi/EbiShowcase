// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
	"math/rand"
)

const (
	w = 480
	h = 720
)

type mote struct {
	x, y, vx, vy float64
	life         int
}
type game struct {
	hp, phase, timer, streak, best, attempts, frames int
	success                                          bool
	motes                                            []mote
	rng                                              *rand.Rand
}

func newGame() *game { return &game{hp: 40, rng: rand.New(rand.NewSource(212))} }
func (g *game) Update() error {
	g.frames++
	for i := len(g.motes) - 1; i >= 0; i-- {
		p := &g.motes[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .07
		p.life--
		if p.life <= 0 {
			g.motes = append(g.motes[:i], g.motes[i+1:]...)
		}
	}
	if g.phase == 8 {
		g.timer--
		if g.timer <= 0 {
			g.hp = 40
			g.phase = 0
		}
		return nil
	}
	if g.phase > 0 {
		g.timer--
		if g.timer == 36 || g.timer == 20 || g.timer == 8 {
			g.phase++
		}
		if g.timer <= 0 {
			if g.success {
				g.streak++
				g.best = max(g.best, g.streak)
				g.burst()
				g.phase = 8
				g.timer = 70
			} else {
				g.streak = 0
				g.phase = 0
			}
		}
		return nil
	}
	hit := inpututil.IsKeyJustPressed(ebiten.Key1)
	orb := inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if x, y, ok := press(); ok && y > 590 {
		if x < 240 {
			hit = true
		} else {
			orb = true
		}
	}
	if hit {
		g.hp = max(1, g.hp-9)
	}
	if orb {
		g.attempts++
		chance := 25 + (40-g.hp)*2
		g.success = g.rng.Intn(100) < chance
		g.phase = 1
		g.timer = 52
	}
	return nil
}
func (g *game) burst() {
	for i := 0; i < 42; i++ {
		a := g.rng.Float64() * math.Pi * 2
		s := 1 + g.rng.Float64()*4
		g.motes = append(g.motes, mote{350, 255, math.Cos(a) * s, math.Sin(a)*s - 1, 25 + g.rng.Intn(22)})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 22, 40, 255})
	ebitenutil.DebugPrintAt(s, "CAPTURE MOMENT", 181, 24)
	chance := 25 + (40-g.hp)*2
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HP %d/40  CHANCE %d%%  STREAK %d  BEST %d", g.hp, chance, g.streak, g.best), 95, 52)
	vector.DrawFilledRect(s, 30, 100, 420, 390, color.RGBA{31, 68, 78, 255}, false)
	bob := math.Sin(float64(g.frames)*.1) * 5
	if g.phase == 0 {
		trackatlas.DrawCentered(s, trackatlas.Species(3), 350, 255+bob, 135)
	} else {
		x := 150.0
		y := 305.0
		if g.phase < 8 {
			p := 1 - float64(g.timer)/52
			x = 150 + 200*p
			y = 305 - math.Sin(p*math.Pi)*130
		}
		if g.phase >= 2 && g.phase < 8 {
			x = 350 + math.Sin(float64(g.timer)*.8)*float64(8-g.phase)
		}
		vector.DrawFilledCircle(s, float32(x), float32(y), 20, color.RGBA{245, 184, 62, 255}, false)
		vector.StrokeCircle(s, float32(x), float32(y), 20, 4, color.White, false)
	}
	for _, p := range g.motes {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), 3, color.RGBA{255, 220, 90, uint8(min(255, p.life*8))}, false)
	}
	msg := "Weaken first, then throw the orb."
	if g.phase > 0 && g.phase < 8 {
		msg = fmt.Sprintf("ORB ROCK %d / 3", min(3, max(1, g.phase-1)))
	}
	if g.phase == 8 {
		msg = "CAPTURED! Can you build a longer streak?"
	}
	ebitenutil.DebugPrintAt(s, msg, 102, 455)
	for i, t := range []string{"[1] QUICK FIN", "[2] THROW ORB"} {
		c := color.RGBA{50, 91, 128, 255}
		if i == 1 {
			c = color.RGBA{185, 124, 53, 255}
		}
		vector.DrawFilledRect(s, float32(i*240+6), 590, 228, 82, c, false)
		ebitenutil.DebugPrintAt(s, t, i*240+58, 625)
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
func (g *game) Layout(int, int) (int, int) { return w, h }
func main() {
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Capture Sequence — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
