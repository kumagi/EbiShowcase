package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
)

const width, height = 480, 720

type game struct {
	hp, enemy, turn int
	message         string
	over            bool
}

func newGame() *game { return &game{hp: 50, enemy: 70, message: "Read NEXT, then choose."} }
func (g *game) Update() error {
	if g.over {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			*g = *newGame()
		}
		return nil
	}
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y > 520 {
			choice = x / (width / 2)
		}
	}
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y > 520 {
			choice = x / (width / 2)
		}
	}
	if choice < 0 {
		return nil
	}
	guard := choice == 1
	if !guard {
		g.enemy -= 12
	}
	intent := g.turn % 3
	damage := []int{8, 20, 5}[intent]
	if guard {
		damage = (damage + 1) / 2
	}
	g.hp -= damage
	g.message = fmt.Sprintf("Enemy dealt %d. Guard=%v", damage, guard)
	g.turn++
	if g.enemy <= 0 || g.hp <= 0 {
		g.over = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{28, 31, 58, 255})
	trackatlas.DrawCentered(s, "boss-crab", 240, 210, 145)
	intent := []string{"NORMAL 8", "HEAVY 20 — GUARD!", "QUICK 5"}[g.turn%3]
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PARTY %d/50   ENEMY %d/70", max(0, g.hp), max(0, g.enemy)), 135, 40)
	ebitenutil.DebugPrintAt(s, "NEXT: "+intent, 160, 370)
	ebitenutil.DebugPrintAt(s, g.message, 130, 440)
	vector.DrawFilledRect(s, 30, 530, 200, 90, color.RGBA{45, 205, 181, 255}, false)
	vector.DrawFilledRect(s, 250, 530, 200, 90, color.RGBA{80, 145, 225, 255}, false)
	ebitenutil.DebugPrintAt(s, "1 ATTACK", 95, 570)
	ebitenutil.DebugPrintAt(s, "2 GUARD", 315, 570)
	if g.over {
		ebitenutil.DebugPrintAt(s, "BATTLE END — SPACE / TAP", 135, 670)
	}
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Enemy Intent Battle — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
