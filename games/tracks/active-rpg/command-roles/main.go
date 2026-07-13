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

const W, H = 480, 720

type game struct {
	hero, enemy, turn, flash int
	message                  string
}

func (g *game) Update() error {
	if g.hero <= 0 || g.enemy <= 0 {
		_, _, tapped := press()
		if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || tapped {
			*g = game{60, 90, 0, 0, "Choose a command."}
		}
		return nil
	}
	a := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		a = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		a = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		a = 2
	}
	if x, y, ok := press(); ok && y > 570 {
		a = x / 160
	}
	if a >= 0 {
		switch a {
		case 0:
			g.enemy -= 12
			g.message = "ATTACK: safe damage"
		case 1:
			g.enemy -= 23
			g.hero -= 5
			g.message = "BURST: big damage, small cost"
		case 2:
			g.hero = min(60, g.hero+18)
			g.message = "HEAL: restore the lowest ally"
		}
		g.hero -= 8
		g.flash = 12
		g.turn++
	}
	if g.flash > 0 {
		g.flash--
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{16, 25, 44, 255})
	trackatlas.DrawCentered(s, "hero", 130, 270, 110)
	trackatlas.DrawCentered(s, "boss-crab", 350, 270, 120)
	if g.flash > 0 {
		vector.StrokeCircle(s, 350, 270, float32(55+(12-g.flash)*5), 5, color.RGBA{255, 205, 70, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HERO HP %d     ENEMY HP %d", max(0, g.hero), max(0, g.enemy)), 105, 100)
	ebitenutil.DebugPrintAt(s, g.message, 120, 430)
	labels := []string{"[1] ATTACK", "[2] BURST", "[3] HEAL"}
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*160+5), 570, 150, 75, color.RGBA{170, 92, 52, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*160+25, 603)
	}
	if g.hero <= 0 || g.enemy <= 0 {
		ebitenutil.DebugPrintAt(s, "BATTLE ENDED / R TO RETRY", 135, 500)
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
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{hero: 60, enemy: 90, message: "Choose a command."}); err != nil {
		panic(err)
	}
}
