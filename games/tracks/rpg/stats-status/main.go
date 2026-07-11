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

type game struct {
	hp, enemyHP, attack, defense, buff, poison, turn int
	message                                          string
	clear, over                                      bool
	rng                                              *rand.Rand
}

func newGame() *game {
	return &game{hp: 55, enemyHP: 80, attack: 15, defense: 7, message: "Choose: Attack, Power Up, or Antidote", rng: rand.New(rand.NewSource(3204))}
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
		c = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		c = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		c = 2
	}
	if x, y, ok := press(); ok && y > 520 {
		c = min(2, x/(width/3))
	}
	if c < 0 {
		return nil
	}
	g.turn++
	switch c {
	case 0:
		atk := g.attack
		if g.buff > 0 {
			atk += 8
		}
		d := max(1, atk-6+g.rng.Intn(5))
		g.enemyHP -= d
		g.message = fmt.Sprintf("Attack %d - enemy DEF 6 = %d damage", atk, d)
	case 1:
		g.buff = 3
		g.message = "Power Up! ATK +8 for three turns"
	case 2:
		if g.poison > 0 {
			g.poison = 0
			g.message = "Antidote removed poison"
		} else {
			g.hp = min(55, g.hp+8)
			g.message = "No poison. Recovered 8 HP instead"
		}
	}
	if g.enemyHP <= 0 {
		g.clear = true
		return nil
	}
	enemyAttack := 13 + g.rng.Intn(5)
	d := max(1, enemyAttack-g.defense)
	g.hp -= d
	if g.turn%3 == 0 {
		g.poison = 4
	}
	if g.poison > 0 {
		g.hp -= 3
		g.poison--
	}
	if g.buff > 0 {
		g.buff--
	}
	if g.hp <= 0 {
		g.over = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{19, 29, 48, 255})
	vector.DrawFilledCircle(s, 240, 190, 70, color.RGBA{169, 91, 205, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("VENOM SLIME HP %02d/80  DEF 6", max(0, g.enemyHP)), 140, 95)
	vector.DrawFilledCircle(s, 240, 400, 36, color.RGBA{240, 74, 90, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI HP %02d/55  ATK %d  DEF %d", max(0, g.hp), g.attack, g.defense), 145, 340)
	effects := "NONE"
	if g.buff > 0 {
		effects = fmt.Sprintf("POWER UP %d", g.buff)
	}
	if g.poison > 0 {
		effects += fmt.Sprintf("  POISON %d", g.poison)
	}
	ebitenutil.DebugPrintAt(s, "STATUS: "+effects, 155, 455)
	vector.DrawFilledRect(s, 20, 490, 440, 190, color.RGBA{5, 17, 35, 245}, false)
	ebitenutil.DebugPrintAt(s, g.message, 35, 515)
	labels := []string{"1 ATTACK", "2 POWER UP", "3 ANTIDOTE"}
	for i, l := range labels {
		x := float32(25 + i*145)
		vector.DrawFilledRect(s, x, 575, 135, 65, color.RGBA{45, 225, 194, 255}, false)
		ebitenutil.DebugPrintAt(s, l, int(x)+18, 602)
	}
	if g.clear {
		overlay(s, "STATUS BATTLE WON!\n\nTAP / SPACE TO RESTART")
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
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Stats and Status — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
