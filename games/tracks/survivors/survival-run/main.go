package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
)

const width, height = 480, 720

type mob struct {
	x, y, hp, r float64
	boss        bool
}

type gem struct{ x, y float64 }

type game struct {
	px, py                             float64
	mobs                               []mob
	gems                               []gem
	rng                                *rand.Rand
	frame, life, inv, kills, xp, need  int
	level                              int
	speed, aura                        float64
	auraTick                           int
	bossSpawned, drafting, clear, over bool
	pickA, pickB, pickC                string
}

func newGame() *game {
	return &game{
		px: 240, py: 360,
		rng:  rand.New(rand.NewSource(2306)),
		life: 4, speed: 3.6, aura: 46, auraTick: 18,
		level: 1, need: 4,
	}
}

func (g *game) Update() error {
	if g.clear || g.over {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	if g.drafting {
		return g.updateDraft()
	}
	g.frame++
	if g.inv > 0 {
		g.inv--
	}
	g.movePlayer()
	sec := g.frame / 60
	interval := max(10, 34-sec/2)
	if g.frame%interval == 0 && sec < 35 {
		a := g.rng.Float64() * math.Pi * 2
		dist := 280.0 + g.rng.Float64()*80
		g.mobs = append(g.mobs, mob{g.px + math.Cos(a)*dist, g.py + math.Sin(a)*dist, 1, 14, false})
	}
	if sec >= 35 && !g.bossSpawned {
		g.bossSpawned = true
		g.mobs = append(g.mobs, mob{240, 95, 80, 40, true})
	}
	chase := 1.15 + float64(sec)*0.028
	for i := len(g.mobs) - 1; i >= 0; i-- {
		m := &g.mobs[i]
		d := math.Hypot(g.px-m.x, g.py-m.y)
		if d < 1 {
			d = 1
		}
		ms := chase
		if m.boss {
			ms = 0.85
		}
		m.x += (g.px - m.x) / d * ms
		m.y += (g.py - m.y) / d * ms

		// Contact hurts first — standing still is unsafe.
		if d < m.r+14 && g.inv == 0 {
			g.life--
			g.inv = 50
			if g.life <= 0 {
				g.over = true
				return nil
			}
		}

		if d < g.aura && g.frame%g.auraTick == 0 {
			m.hp--
			if m.hp <= 0 {
				if m.boss {
					g.clear = true
				} else {
					g.kills++
					g.gems = append(g.gems, gem{m.x, m.y})
				}
				g.mobs = append(g.mobs[:i], g.mobs[i+1:]...)
			}
		}
	}
	for i := len(g.gems) - 1; i >= 0; i-- {
		gem := &g.gems[i]
		d := math.Hypot(g.px-gem.x, g.py-gem.y)
		if d < 120 {
			gem.x += (g.px - gem.x) * 0.12
			gem.y += (g.py - gem.y) * 0.12
		}
		if d < 22 {
			g.xp++
			g.gems = append(g.gems[:i], g.gems[i+1:]...)
			if g.xp >= g.need {
				g.openDraft()
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

func (g *game) movePlayer() {
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
		g.px += dx / l * g.speed
		g.py += dy / l * g.speed
	}
	g.px = clamp(g.px, 20, 460)
	g.py = clamp(g.py, 90, 690)
}

func (g *game) openDraft() {
	g.drafting = true
	opts := []string{"AURA+", "SPEED+", "LIFE+"}
	g.rng.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	g.pickA, g.pickB, g.pickC = opts[0], opts[1], opts[2]
}

func (g *game) updateDraft() error {
	choice := ""
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		choice = g.pickA
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		choice = g.pickB
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		choice = g.pickC
	}
	touch := inpututil.AppendJustPressedTouchIDs(nil)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(touch) > 0 {
		x, y := ebiten.CursorPosition()
		if len(touch) > 0 {
			x, y = ebiten.TouchPosition(touch[0])
		}
		switch {
		case y > 300 && y < 370 && x < 160:
			choice = g.pickA
		case y > 300 && y < 370 && x < 320:
			choice = g.pickB
		case y > 300 && y < 370:
			choice = g.pickC
		}
	}
	if choice == "" {
		return nil
	}
	switch choice {
	case "AURA+":
		g.aura = math.Min(90, g.aura+10)
		g.auraTick = max(10, g.auraTick-2)
	case "SPEED+":
		g.speed = math.Min(6.2, g.speed+0.45)
	case "LIFE+":
		g.life++
	}
	g.level++
	g.xp = 0
	g.need = min(12, 4+g.level)
	g.drafting = false
	g.inv = 40
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 24, 39, 255})
	for _, gem := range g.gems {
		vector.DrawFilledCircle(s, float32(gem.x), float32(gem.y), 7, color.RGBA{90, 220, 255, 255}, false)
	}
	for _, m := range g.mobs {
		c := color.RGBA{226, 70, 104, 255}
		if m.boss {
			c = color.RGBA{153, 70, 210, 255}
		}
		vector.DrawFilledCircle(s, float32(m.x), float32(m.y), float32(m.r), c, false)
		if m.boss {
			vector.DrawFilledRect(s, float32(m.x-45), float32(m.y-52), 90, 7, color.RGBA{50, 48, 70, 255}, false)
			vector.DrawFilledRect(s, float32(m.x-45), float32(m.y-52), float32(90*m.hp/80), 7, color.RGBA{255, 211, 62, 255}, false)
		}
	}
	vector.DrawFilledCircle(s, float32(g.px), float32(g.py), float32(g.aura), color.RGBA{255, 211, 62, 28}, false)
	vector.StrokeCircle(s, float32(g.px), float32(g.py), float32(g.aura), 2, color.RGBA{255, 211, 62, 170}, false)
	if g.inv%10 < 5 {
		hero.DrawCentered(s, g.px, g.py, 34)
	}
	sec := g.frame / 60
	wave := min(3, sec/15+1)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TIME %02d/45  LV%d  XP %d/%d  LIFE %d  KILLS %03d", sec, g.level, g.xp, g.need, g.life, g.kills), 28, 24)
	msg := "WAVE 1 — KEEP MOVING"
	if wave == 2 {
		msg = "WAVE 2 — ENEMIES SPEED UP"
	}
	if wave == 3 {
		msg = "FINAL WAVE — DEFEAT THE BOSS"
	}
	ebitenutil.DebugPrintAt(s, msg, 120, 52)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / DRAG    GEMS → LEVEL UP PICK", 55, 685)
	if g.drafting {
		vector.DrawFilledRect(s, 40, 250, 400, 180, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "LEVEL UP — PICK ONE", 150, 270)
		drawPick(s, 55, 310, g.pickA)
		drawPick(s, 175, 310, g.pickB)
		drawPick(s, 295, 310, g.pickC)
		ebitenutil.DebugPrintAt(s, "1/2/3 or TAP a card", 150, 400)
	}
	if g.clear {
		overlay(s, "EBI SURVIVORS CLEAR!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "RUN ENDED!\n\nTAP / SPACE TO RETRY")
	}
}

func drawPick(s *ebiten.Image, x, y float64, label string) {
	vector.DrawFilledRect(s, float32(x), float32(y), 110, 55, color.RGBA{40, 70, 120, 255}, false)
	ebitenutil.DebugPrintAt(s, label, int(x)+28, int(y)+20)
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
