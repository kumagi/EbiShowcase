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
type weapon struct {
	name                           string
	cooldown, speed, count, damage int
}
type game struct {
	px, py                                     float64
	enemies, shots                             []body
	item                                       body
	weapon                                     weapon
	rng                                        *rand.Rand
	frame, kills, life, inv                    int
	itemReady, itemTaken, evolved, clear, over bool
}

func newGame() *game {
	return &game{px: 240, py: 360, weapon: weapon{"Ebi Needle", 32, 7, 1, 1}, rng: rand.New(rand.NewSource(2205)), life: 5}
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
	dx, dy := inputDir(g.px, g.py)
	if dx != 0 || dy != 0 {
		l := math.Hypot(dx, dy)
		g.px += dx / l * 3.8
		g.py += dy / l * 3.8
	}
	g.px = clamp(g.px, 20, 460)
	g.py = clamp(g.py, 90, 690)
	if g.frame%38 == 0 {
		a := g.rng.Float64() * math.Pi * 2
		g.enemies = append(g.enemies, body{g.px + math.Cos(a)*340, g.py + math.Sin(a)*340, 0, 0, 14})
	}
	if g.frame%g.weapon.cooldown == 0 && len(g.enemies) > 0 {
		best, d := 0, math.MaxFloat64
		for i, e := range g.enemies {
			v := math.Hypot(e.x-g.px, e.y-g.py)
			if v < d {
				best, d = i, v
			}
		}
		base := math.Atan2(g.enemies[best].y-g.py, g.enemies[best].x-g.px)
		for n := 0; n < g.weapon.count; n++ {
			a := base + (float64(n)-float64(g.weapon.count-1)/2)*.18
			g.shots = append(g.shots, body{g.px, g.py, math.Cos(a) * float64(g.weapon.speed), math.Sin(a) * float64(g.weapon.speed), 5})
		}
	}
	for i := range g.enemies {
		e := &g.enemies[i]
		d := math.Hypot(g.px-e.x, g.py-e.y)
		e.x += (g.px - e.x) / d
		e.y += (g.py - e.y) / d
		if d < 24 && g.inv == 0 {
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
		rm := s.x < 0 || s.x > width || s.y < 70 || s.y > height
		for j := len(g.enemies) - 1; j >= 0 && !rm; j-- {
			if math.Hypot(s.x-g.enemies[j].x, s.y-g.enemies[j].y) < 20 {
				g.enemies = append(g.enemies[:j], g.enemies[j+1:]...)
				g.kills++
				rm = true
			}
		}
		if rm {
			g.shots = append(g.shots[:i], g.shots[i+1:]...)
		}
	}
	if g.kills >= 10 && !g.itemReady && !g.itemTaken {
		g.item = body{240, 360, 0, 0, 12}
		g.itemReady = true
	}
	if g.itemReady && math.Hypot(g.px-g.item.x, g.py-g.item.y) < 28 {
		g.itemReady = false
		g.itemTaken = true
	}
	if g.kills >= 20 && g.itemTaken && !g.evolved {
		g.evolved = true
		g.weapon = weapon{"Ebi Storm", 16, 9, 3, 2}
	}
	if g.evolved && g.kills >= 50 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 26, 42, 255})
	for _, e := range g.enemies {
		vector.DrawFilledCircle(s, float32(e.x), float32(e.y), 14, color.RGBA{228, 72, 105, 255}, false)
	}
	for _, b := range g.shots {
		vector.DrawFilledCircle(s, float32(b.x), float32(b.y), 5, color.RGBA{255, 213, 62, 255}, false)
	}
	if g.itemReady {
		vector.DrawFilledRect(s, float32(g.item.x-12), float32(g.item.y-12), 24, 24, color.RGBA{92, 170, 255, 255}, false)
	}
	if g.inv%10 < 5 {
		vector.DrawFilledCircle(s, float32(g.px), float32(g.py), 16, color.RGBA{45, 225, 194, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s   KILLS %02d/50   LIFE %d", g.weapon.name, g.kills, g.life), 70, 24)
	status := "DEFEAT 10 ENEMIES TO FIND THE BLUE CORE"
	if g.itemReady {
		status = "COLLECT THE BLUE CORE"
	} else if g.itemTaken && !g.evolved {
		status = "CORE READY — REACH 20 KILLS TO EVOLVE"
	} else if g.evolved {
		status = "EVOLVED! EBI STORM FIRES THREE SHOTS"
	}
	ebitenutil.DebugPrintAt(s, status, 80, 52)
	ebitenutil.DebugPrintAt(s, "MOVE ONLY — WEAPON FIRES AUTOMATICALLY", 90, 685)
	if g.clear {
		overlay(s, "WEAPON MASTERED!\n\nTAP / SPACE TO PLAY AGAIN")
	} else if g.over {
		overlay(s, "LAB FAILED!\n\nTAP / SPACE TO RETRY")
	}
}
func inputDir(px, py float64) (float64, float64) {
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
		dx = float64(x) - px
		dy = float64(y) - py
	}
	return dx, dy
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
	ebiten.SetWindowTitle("Weapon Evolution Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
