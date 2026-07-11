package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"time"
)

const width, height = 480, 720

type game struct {
	points float64
	labs   int
	last   time.Time
	clear  bool
}

func newGame() *game          { return &game{points: 10, last: time.Now()} }
func (g *game) cost() float64 { return 10 * math.Pow(1.18, float64(g.labs)) }
func (g *game) rate() float64 { return float64(g.labs*g.labs) * 8 }
func (g *game) Update() error {
	now := time.Now()
	dt := math.Min(.25, now.Sub(g.last).Seconds())
	g.last = now
	if !g.clear {
		g.points += g.rate() * dt
	}
	if g.clear {
		if any() {
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
			g.points += 1 + float64(g.labs)
		} else if g.points >= g.cost() {
			g.points -= g.cost()
			g.labs++
		}
	}
	if g.points >= 1_000_000 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{25, 27, 47, 255})
	ebitenutil.DebugPrintAt(s, "GROWTH CURVE LAB", 170, 45)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENERGY %s / 1.00M", short(g.points)), 155, 90)
	vector.DrawFilledCircle(s, 240, 285, 105, color.RGBA{104, 88, 220, 255}, false)
	vector.StrokeCircle(s, 240, 285, 105, 7, color.RGBA{178, 163, 255, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TAP POWER +%d", 1+g.labs), 190, 280)
	vector.DrawFilledRect(s, 55, 480, 370, 150, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BUY LAB [B]\n\nCOST %s     OWNED %d\nPRODUCTION %s / SECOND", short(g.cost()), g.labs, short(g.rate())), 130, 505)
	if g.clear {
		overlay(s, "ONE MILLION ENERGY!\n\nTAP / SPACE TO RESTART")
	}
}
func short(v float64) string {
	switch {
	case v >= 1e9:
		return fmt.Sprintf("%.2fB", v/1e9)
	case v >= 1e6:
		return fmt.Sprintf("%.2fM", v/1e6)
	case v >= 1e3:
		return fmt.Sprintf("%.2fK", v/1e3)
	default:
		return fmt.Sprintf("%.1f", v)
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
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 240}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Growth Curve Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
