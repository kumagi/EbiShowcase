package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

const width, height = 480, 720
const (
	strike = iota
	guard
	throw
)

type game struct {
	enemy, score, life, round int
	message                   string
	clear, over               bool
	rng                       *rand.Rand
}

func newGame() *game { g := &game{life: 3, rng: rand.New(rand.NewSource(3904))}; g.next(); return g }
func (g *game) next() {
	g.enemy = g.rng.Intn(3)
	g.round++
	g.message = "Enemy prepares " + name(g.enemy) + ". Choose a counter!"
}
func (g *game) Update() error {
	if g.clear || g.over {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	c := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		c = strike
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		c = guard
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		c = throw
	}
	if x, y, ok := press(); ok && y > 520 {
		c = min(2, x/(width/3))
	}
	if c < 0 {
		return nil
	}
	if c == g.enemy {
		g.message = "Same choice — clash!"
	} else if wins(c, g.enemy) {
		g.score++
		g.message = name(c) + " beats " + name(g.enemy) + "! Point!"
	} else {
		g.life--
		g.message = name(g.enemy) + " beats " + name(c) + "! Damage!"
	}
	if g.score >= 5 {
		g.clear = true
	} else if g.life <= 0 {
		g.over = true
	} else {
		g.next()
	}
	return nil
}
func wins(a, b int) bool {
	return (a == strike && b == throw) || (a == throw && b == guard) || (a == guard && b == strike)
}
func name(v int) string { return []string{"STRIKE", "GUARD", "THROW"}[v] }
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{19, 28, 44, 255})
	vector.DrawFilledCircle(s, 130, 350, 45, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledCircle(s, 350, 350, 45, color.RGBA{240, 75, 91, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %d/5   LIFE %d   ROUND %d", g.score, g.life, g.round), 135, 40)
	ebitenutil.DebugPrintAt(s, "ENEMY TELEGRAPH: "+name(g.enemy), 145, 120)
	ebitenutil.DebugPrintAt(s, "GUARD > STRIKE > THROW > GUARD", 125, 165)
	vector.DrawFilledRect(s, 25, 450, 430, 70, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, g.message, 45, 480)
	labels := []string{"1 STRIKE", "2 GUARD", "3 THROW"}
	cols := []color.RGBA{{240, 85, 80, 255}, {75, 145, 240, 255}, {225, 170, 65, 255}}
	for i, l := range labels {
		x := float32(20 + i*150)
		vector.DrawFilledRect(s, x, 555, 140, 80, cols[i], false)
		ebitenutil.DebugPrintAt(s, l, int(x)+30, 590)
	}
	if g.clear {
		overlay(s, "READING GAME WON!\n\nTAP / SPACE TO RESET")
	} else if g.over {
		overlay(s, "OUTREAD!\n\nTAP / SPACE TO RETRY")
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
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Guard and Throw — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
