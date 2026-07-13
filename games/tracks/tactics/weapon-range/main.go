// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
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

const (
	W  = 480
	H  = 720
	N  = 7
	T  = 58
	OX = 37
	OY = 100
)

type pt struct{ x, y int }
type game struct {
	hero                              pt
	enemies                           []pt
	weapon, score, best, flash, frame int
	rng                               *rand.Rand
	message                           string
}

func newGame() *game {
	g := &game{hero: pt{3, 5}, rng: rand.New(rand.NewSource(34)), message: "Choose BLADE (range 1) or BOW (range 2)."}
	g.spawn()
	return g
}
func (g *game) spawn() {
	// Two targets fit the blade ring and two fit the bow ring, so both tools
	// always have a readable answer instead of depending on a lucky spawn.
	g.enemies = []pt{{2, 5}, {4, 5}, {3, 3}, {5, 5}}
}
func (g *game) Update() error {
	g.frame++
	if g.flash > 0 {
		g.flash--
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.weapon = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.weapon = 1
	}
	if x, y, ok := press(); ok {
		if y > 590 {
			g.weapon = min(1, x/240)
			return nil
		}
		if x >= OX && x < OX+N*T && y >= OY && y < OY+N*T {
			p := pt{(x - OX) / T, (y - OY) / T}
			for i, e := range g.enemies {
				if e == p {
					d := abs(e.x-g.hero.x) + abs(e.y-g.hero.y)
					r := g.weapon + 1
					if d == r {
						g.score += 10 + r*5
						g.best = max(g.best, g.score)
						g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
						g.flash = 16
						g.message = fmt.Sprintf("Perfect range %d hit!", r)
						if len(g.enemies) == 0 {
							g.spawn()
							g.message = "Wave cleared! New targets arrived."
						}
					} else {
						g.score = max(0, g.score-5)
						g.message = fmt.Sprintf("Distance %d is outside exact range %d.", d, r)
					}
					break
				}
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{13, 22, 38, 255})
	ebitenutil.DebugPrintAt(s, "WEAPON RANGE DRILL", 166, 24)
	names := []string{"BLADE", "BOW"}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s RANGE %d  SCORE %d  BEST %d", names[g.weapon], g.weapon+1, g.score, g.best), 119, 53)
	ebitenutil.DebugPrintAt(s, g.message, 86, 78)
	for y := 0; y < N; y++ {
		for x := 0; x < N; x++ {
			d := abs(x-g.hero.x) + abs(y-g.hero.y)
			c := color.RGBA{43, 67, 83, 255}
			if d == g.weapon+1 {
				c = color.RGBA{155, 94, 64, 255}
			}
			vector.DrawFilledRect(s, float32(OX+x*T+1), float32(OY+y*T+1), T-2, T-2, c, false)
		}
	}
	hx := float32(OX + g.hero.x*T + T/2)
	hy := float32(OY + g.hero.y*T + T/2)
	vector.DrawFilledCircle(s, hx, hy, 19, color.RGBA{241, 175, 65, 255}, false)
	for i, e := range g.enemies {
		x := float32(OX + e.x*T + T/2)
		y := float32(OY + e.y*T + T/2 + int(math.Sin(float64(g.frame+i*8)*.15)*3))
		c := color.RGBA{172, 78, 111, 255}
		if g.flash > 0 && i == 0 {
			c = color.RGBA{255, 255, 255, 255}
		}
		vector.DrawFilledCircle(s, x, y, 17, c, false)
	}
	for i, n := range names {
		c := color.RGBA{54, 85, 122, 255}
		if i == g.weapon {
			c = color.RGBA{188, 126, 57, 255}
		}
		vector.DrawFilledRect(s, float32(i*240+6), 590, 228, 78, c, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("[%d] %s", i+1, n), i*240+80, 623)
	}
}
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
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
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Weapon Range Drill")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
