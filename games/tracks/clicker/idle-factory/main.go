package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"time"
)

const width, height = 480, 720

type game struct {
	sweets   float64
	machines int
	cost     float64
	last     time.Time
	clear    bool
}

func newGame() *game { return &game{cost: 15, last: time.Now()} }
func (g *game) Update() error {
	now := time.Now()
	dt := now.Sub(g.last).Seconds()
	g.last = now
	if dt > .25 {
		dt = .25
	}
	if !g.clear {
		g.sweets += float64(g.machines) * dt
	}
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
		} else if g.sweets >= g.cost {
			g.sweets -= g.cost
			g.machines++
			g.cost += 10
		}
	}
	if g.sweets >= 200 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{26, 30, 48, 255})
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SWEETS %06.1f / 200", g.sweets), 150, 55)
	vector.DrawFilledCircle(s, 240, 275, 100, color.RGBA{222, 143, 76, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: BAKE +1", 155, 405)
	vector.DrawFilledRect(s, 60, 480, 360, 145, color.RGBA{45, 196, 174, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BUY AUTO OVEN [B]\n\nCOST %.0f     OWNED %d\nPRODUCTION %.1f SWEETS / SECOND", g.cost, g.machines, float64(g.machines)), 135, 505)
	if g.clear {
		overlay(s, "200 SWEETS PRODUCED!\n\nTAP / SPACE TO RESTART")
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
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Idle Factory — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
