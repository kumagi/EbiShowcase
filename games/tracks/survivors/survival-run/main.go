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
	x, y, hp, r float64
	boss        bool
}
type game struct {
	px, py                   float64
	mobs                     []mob
	rng                      *rand.Rand
	frame, life, inv, kills  int
	bossSpawned, clear, over bool
}

func newGame() *game { return &game{px: 240, py: 360, rng: rand.New(rand.NewSource(2306)), life: 6} }
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
	sec := g.frame / 60
	interval := max(14, 42-sec/2)
	if g.frame%interval == 0 && sec < 35 {
		a := g.rng.Float64() * math.Pi * 2
		g.mobs = append(g.mobs, mob{g.px + math.Cos(a)*350, g.py + math.Sin(a)*350, 1, 13, false})
	}
	if sec >= 35 && !g.bossSpawned {
		g.bossSpawned = true
		g.mobs = append(g.mobs, mob{240, 95, 60, 38, true})
	}
	speed := .85 + float64(sec)*.018
	for i := len(g.mobs) - 1; i >= 0; i-- {
		m := &g.mobs[i]
		d := math.Hypot(g.px-m.x, g.py-m.y)
		ms := speed
		if m.boss {
			ms = .65
		}
		m.x += (g.px - m.x) / d * ms
		m.y += (g.py - m.y) / d * ms
		if d < 68 && g.frame%14 == 0 {
			m.hp--
			if m.hp <= 0 {
				if m.boss {
					g.clear = true
				} else {
					g.kills++
				}
				g.mobs = append(g.mobs[:i], g.mobs[i+1:]...)
				continue
			}
		}
		if d < m.r+14 && g.inv == 0 {
			g.life--
			g.inv = 90
			if g.life <= 0 {
				g.over = true
			}
		}
	}
	if sec >= 45 && g.bossSpawned {
		bossAlive := false
		for _, m := range g.mobs {
			if m.boss {
				bossAlive = true
			}
		}
		if !bossAlive {
			g.clear = true
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 24, 39, 255})
	for _, m := range g.mobs {
		c := color.RGBA{226, 70, 104, 255}
		if m.boss {
			c = color.RGBA{153, 70, 210, 255}
		}
		vector.DrawFilledCircle(s, float32(m.x), float32(m.y), float32(m.r), c, false)
		if m.boss {
			vector.DrawFilledRect(s, float32(m.x-45), float32(m.y-52), 90, 7, color.RGBA{50, 48, 70, 255}, false)
			vector.DrawFilledRect(s, float32(m.x-45), float32(m.y-52), float32(90*m.hp/60), 7, color.RGBA{255, 211, 62, 255}, false)
		}
	}
	vector.DrawFilledCircle(s, float32(g.px), float32(g.py), 68, color.RGBA{255, 211, 62, 28}, false)
	vector.StrokeCircle(s, float32(g.px), float32(g.py), 68, 2, color.RGBA{255, 211, 62, 170}, false)
	if g.inv%10 < 5 {
		hero.DrawCentered(s, g.px, g.py, 34)
	}
	sec := g.frame / 60
	wave := min(3, sec/15+1)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TIME %02d/45   WAVE %d/3   LIFE %d   KILLS %03d", sec, wave, g.life, g.kills), 55, 24)
	msg := "WAVE 1 — LEARN THE ARENA"
	if wave == 2 {
		msg = "WAVE 2 — ENEMIES SPEED UP"
	}
	if wave == 3 {
		msg = "FINAL WAVE — DEFEAT THE BOSS"
	}
	ebitenutil.DebugPrintAt(s, msg, 120, 52)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / ARROWS / TOUCH    AURA AUTO ATTACKS", 70, 685)
	if g.clear {
		overlay(s, "EBI SURVIVORS CLEAR!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "RUN ENDED!\n\nTAP / SPACE TO RETRY")
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
	ebiten.SetWindowTitle("Ebi Survivors — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
