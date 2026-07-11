package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const width, height = 480, 720

type game struct {
	sweets, ovens, cost, flash int
	clear                      bool
}

func newGame() *game { return &game{cost: 10} }
func (g *game) Update() error {
	if g.clear {
		if anyPress() {
			*g = *newGame()
		}
		return nil
	}
	_, y, ok := press()
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ok = true
		y = 300
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		ok = true
		y = 570
	}
	if ok {
		if y < 470 {
			g.sweets++
			g.flash = 8
		} else if g.sweets >= g.cost {
			g.sweets -= g.cost
			g.ovens++
			g.cost += 5
			g.flash = 12
		} else {
			g.flash = -15
		}
	}
	if g.flash > 0 {
		g.flash--
	} else if g.flash < 0 {
		g.flash++
	}
	if g.ovens >= 5 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{35, 24, 46, 255})
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SWEETS %03d", g.sweets), 190, 55)
	r := float32(105 + max(g.flash, 0))
	vector.DrawFilledCircle(s, 240, 290, r, color.RGBA{221, 143, 77, 255}, false)
	vector.StrokeCircle(s, 240, 290, r, 7, color.RGBA{255, 204, 115, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: BAKE +1", 155, 420)
	c := color.RGBA{45, 225, 194, 255}
	if g.sweets < g.cost {
		c = color.RGBA{91, 93, 110, 255}
	}
	vector.DrawFilledRect(s, 65, 500, 350, 120, c, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BUY OVEN  [B]\n\nCOST %d SWEETS     OWNED %d/5", g.cost, g.ovens), 145, 525)
	if g.flash < 0 {
		ebitenutil.DebugPrintAt(s, "NOT ENOUGH SWEETS", 160, 650)
	}
	if g.clear {
		overlay(s, "SHOP COMPLETE!\n\nTAP / SPACE TO PLAY AGAIN")
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
func anyPress() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 240}, false)
	ebitenutil.DebugPrintAt(s, msg, 140, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("First Shop — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
