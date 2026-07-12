// guard-throw — reading game with fighting-game staging: windup, range, choices.
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
)

const width, height = 480, 720

const (
	strike = iota
	guard
	throwAct
)

type game struct {
	p1, p2          float64
	enemyAct        int
	windup, resolve int // windup counts down, then resolve flash
	score, life     int
	round           int
	msg             string
	clear, over     bool
	flash           float64
	rng             *rand.Rand
	chosen          int // -1 none
}

func newGame() *game {
	g := &game{p1: 140, p2: 320, life: 3, chosen: -1, rng: rand.New(rand.NewSource(3904))}
	g.startRound()
	return g
}

func (g *game) startRound() {
	g.round++
	g.enemyAct = g.rng.Intn(3)
	g.windup = 70 + g.rng.Intn(25)
	g.resolve = 0
	g.chosen = -1
	g.msg = "Enemy winds up… pick STRIKE / GUARD / THROW!"
}

func name(v int) string { return []string{"STRIKE", "GUARD", "THROW"}[v] }

func wins(a, b int) bool {
	return (a == strike && b == throwAct) || (a == throwAct && b == guard) || (a == guard && b == strike)
}

func (g *game) pick() int {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		return strike
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		return guard
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		return throwAct
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		var x, y int
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y = ebiten.CursorPosition()
		} else {
			ids := inpututil.AppendJustPressedTouchIDs(nil)
			x, y = ebiten.TouchPosition(ids[0])
		}
		if y > 560 {
			return min(2, x/(width/3))
		}
		// Upper tap also cycles nothing — require button row
	}
	return -1
}

func (g *game) Update() error {
	if g.clear || g.over {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
			len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			*g = *newGame()
		}
		return nil
	}
	if g.flash > 0 {
		g.flash -= 0.06
	}

	// Approach during windup for fighting feel.
	ideal := 200.0
	if g.windup > 0 {
		g.p2 += (ideal + 40 - g.p2) * 0.04
		g.p1 += (ideal - 40 - g.p1) * 0.03
		if c := g.pick(); c >= 0 {
			g.chosen = c
			g.msg = "Locked " + name(c) + "! Waiting for impact…"
		}
		g.windup--
		if g.windup == 0 {
			g.resolve = 45
			you := g.chosen
			if you < 0 {
				you = g.rng.Intn(3) // panic pick if no input
				g.msg = "Too slow — random " + name(you)
			}
			if you == g.enemyAct {
				g.msg = "CLASH! Both chose " + name(you)
				g.flash = 0.4
			} else if wins(you, g.enemyAct) {
				g.score++
				g.msg = name(you) + " beats " + name(g.enemyAct) + "!"
				g.flash = 0.85
				g.p2 += 28
			} else {
				g.life--
				g.msg = name(g.enemyAct) + " beats " + name(you) + " — hit!"
				g.flash = 0.7
				g.p1 -= 28
			}
		}
		return nil
	}

	if g.resolve > 0 {
		g.resolve--
		if g.resolve == 0 {
			if g.score >= 5 {
				g.clear = true
			} else if g.life <= 0 {
				g.over = true
			} else {
				g.startRound()
			}
		}
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{16, 22, 38, 255})
	if g.flash > 0 {
		a := g.flash
		vector.DrawFilledRect(s, 0, 0, width, height, color.RGBA{uint8(255 * a), uint8(120 * a), uint8(80 * a), uint8(100 * a)}, false)
	}
	// Arena
	vector.DrawFilledRect(s, 0, 520, width, 40, color.RGBA{50, 60, 75, 255}, false)
	vector.DrawFilledRect(s, 0, 560, width, 160, color.RGBA{28, 34, 50, 255}, false)

	// Telegraph banner during windup
	if g.windup > 0 {
		shake := float32(math.Sin(float64(g.windup)*0.5) * 3)
		vector.DrawFilledRect(s, 40, 70, 400, 70, color.RGBA{90, 40, 50, 230}, false)
		ebitenutil.DebugPrintAt(s, "ENEMY WINDUP — read the stance!", 110+int(shake), 90)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("impact in %d frames — choose now!", g.windup), 100, 115)
		hint := []string{"fist high (strike?)", "arms up (guard?)", "hands low (throw?)"}[g.enemyAct]
		ebitenutil.DebugPrintAt(s, "stance hint: "+hint, 130, 160)
	}

	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %d/5   LIFE %d   ROUND %d", g.score, g.life, g.round), 120, 30)
	ebitenutil.DebugPrintAt(s, "GUARD > STRIKE > THROW > GUARD", 130, 200)

	drawFighter(s, g.p1, color.RGBA{45, 225, 194, 255}, true, g.chosen, g.windup == 0 && g.resolve > 30)
	drawFighter(s, g.p2, color.RGBA{240, 75, 91, 255}, false, g.enemyAct, g.windup == 0 && g.resolve > 30)

	vector.DrawFilledRect(s, 30, 430, 420, 55, color.RGBA{8, 14, 28, 230}, false)
	ebitenutil.DebugPrintAt(s, g.msg, 45, 450)

	labels := []string{"1 STRIKE", "2 GUARD", "3 THROW"}
	cols := []color.RGBA{{220, 70, 70, 255}, {70, 130, 230, 255}, {220, 170, 50, 255}}
	for i, l := range labels {
		x := float32(16 + i*155)
		fill := cols[i]
		if g.chosen == i {
			vector.StrokeRect(s, x-2, 575, 148, 70, 4, color.RGBA{120, 240, 220, 255}, false)
		}
		vector.DrawFilledRect(s, x, 578, 144, 64, fill, false)
		ebitenutil.DebugPrintAt(s, l, int(x)+28, 602)
	}
	ebitenutil.DebugPrintAt(s, "Tap a button during the enemy windup", 95, 690)

	if g.clear {
		overlay(s, "YOU READ EVERY MIXUP!\n\nTAP TO RESET")
	} else if g.over {
		overlay(s, "OUTREAD!\n\nTAP TO RETRY")
	}
}

func drawFighter(s *ebiten.Image, x float64, c color.RGBA, faceRight bool, act int, striking bool) {
	fx := float32(x)
	vector.DrawFilledCircle(s, fx, 455, 20, c, false)
	vector.DrawFilledRect(s, fx-16, 475, 32, 70, c, false)
	// Attack/guard pose
	arm := float32(28)
	if striking {
		arm = 70
	}
	dir := float32(1)
	if !faceRight {
		dir = -1
	}
	switch act {
	case strike:
		vector.DrawFilledRect(s, fx+dir*10, 490, dir*arm, 14, color.RGBA{255, 220, 100, 255}, false)
	case guard:
		vector.DrawFilledRect(s, fx+dir*8, 470, dir*22, 50, color.RGBA{180, 200, 255, 220}, false)
	case throwAct:
		vector.DrawFilledRect(s, fx+dir*6, 500, dir*40, 18, color.RGBA{255, 180, 80, 255}, false)
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 250, 370, 140, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 120, 295)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Guard & Throw Reading — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
