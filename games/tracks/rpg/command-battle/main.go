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
	command = iota
	playerTurn
	enemyTurn
	result
)

type game struct {
	hp, maxHP, enemyHP, state, timer, choice int
	guard, clear, over                       bool
	message                                  string
	rng                                      *rand.Rand
}

func newGame() *game {
	return &game{hp: 40, maxHP: 40, enemyHP: 55, state: command, message: "Choose a command.", rng: rand.New(rand.NewSource(3103))}
}
func (g *game) Update() error {
	if g.clear || g.over {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	switch g.state {
	case command:
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
		if x, y, ok := press(); ok && y > 500 {
			c = min(2, x/(width/3))
		}
		if c >= 0 {
			g.choice = c
			g.state = playerTurn
			g.timer = 35
			g.guard = false
			// Apply the chosen effect as soon as the player-turn motion starts
			// (not when the timer later hits zero).
			switch c {
			case 0:
				d := 8 + g.rng.Intn(7)
				g.enemyHP -= d
				g.message = fmt.Sprintf("Ebi attacks! %d damage!", d)
			case 1:
				h := 8 + g.rng.Intn(6)
				g.hp = min(g.maxHP, g.hp+h)
				g.message = fmt.Sprintf("Ebi heals %d HP!", h)
			case 2:
				g.guard = true
				g.message = "Ebi braces for the attack!"
			}
		}
	case playerTurn:
		g.timer--
		if g.timer <= 0 {
			if g.enemyHP <= 0 {
				g.state = result
				g.clear = true
				g.message = "Victory!"
			} else {
				// Start enemy motion AND apply damage on this frame (attack start).
				g.state = enemyTurn
				g.timer = 40
				d := 7 + g.rng.Intn(8)
				if g.guard {
					d = (d + 1) / 2
				}
				g.hp -= d
				g.message = fmt.Sprintf("Slime attacks! %d damage!", d)
			}
		}
	case enemyTurn:
		// Timer is follow-through only; HP already changed when the motion began.
		g.timer--
		if g.timer <= 0 {
			if g.hp <= 0 {
				g.state = result
				g.over = true
				g.message = "Ebi was defeated..."
			} else {
				g.state = command
				g.message = "Choose the next command."
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 28, 48, 255})
	vector.DrawFilledCircle(s, 240, 210, 75, color.RGBA{112, 190, 235, 255}, false)
	vector.DrawFilledCircle(s, 210, 195, 8, color.White, false)
	vector.DrawFilledCircle(s, 270, 195, 8, color.White, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SLIME HP %02d/55", max(0, g.enemyHP)), 185, 105)
	vector.DrawFilledCircle(s, 240, 430, 35, color.RGBA{240, 74, 90, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI HP %02d/%02d", max(0, g.hp), g.maxHP), 190, 380)
	vector.DrawFilledRect(s, 20, 480, 440, 190, color.RGBA{5, 17, 35, 245}, false)
	ebitenutil.DebugPrintAt(s, g.message, 45, 505)
	labels := []string{"1 ATTACK", "2 HEAL", "3 GUARD"}
	for i, l := range labels {
		x := float32(25 + i*145)
		c := color.RGBA{45, 225, 194, 255}
		if g.state != command {
			c = color.RGBA{75, 84, 101, 255}
		}
		vector.DrawFilledRect(s, x, 560, 135, 70, c, false)
		ebitenutil.DebugPrintAt(s, l, int(x)+25, 588)
	}
	if g.clear {
		overlay(s, "BATTLE WON!\n\nTAP / SPACE TO FIGHT AGAIN")
	} else if g.over {
		overlay(s, "TRY AGAIN!\n\nTAP / SPACE TO RETRY")
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
	ebitenutil.DebugPrintAt(s, msg, 145, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Command Battle — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
