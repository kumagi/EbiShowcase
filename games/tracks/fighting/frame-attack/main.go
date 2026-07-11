package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float32 }
type game struct {
	frame, hits, clock int
	hitThis, clear     bool
}

func (g *game) Update() error {
	if g.clear {
		if press() {
			*g = game{}
		}
		return nil
	}
	g.clock++
	if g.frame == 0 && press() {
		g.frame = 1
		g.hitThis = false
	} else if g.frame > 0 {
		g.frame++
		if g.frame >= 31 {
			g.frame = 0
		}
	}
	if g.frame >= 9 && g.frame <= 12 && !g.hitThis {
		target := g.target()
		attack := rect{235, 495, 105, 35}
		if overlap(attack, target) {
			g.hits++
			g.hitThis = true
			if g.hits >= 5 {
				g.clear = true
			}
		}
	}
	return nil
}
func (g *game) target() rect {
	x := float32(345 + math.Sin(float64(g.clock)*.035)*75)
	return rect{x - 22, 475, 44, 110}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{19, 28, 44, 255})
	vector.DrawFilledRect(s, 0, 590, 480, 130, color.RGBA{55, 66, 78, 255}, false)
	vector.DrawFilledCircle(s, 190, 485, 22, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 172, 507, 36, 78, color.RGBA{45, 225, 194, 255}, false)
	t := g.target()
	vector.DrawFilledRect(s, t.x, t.y, t.w, t.h, color.RGBA{240, 75, 91, 255}, false)
	phase := "READY"
	c := color.RGBA{90, 100, 115, 255}
	if g.frame >= 1 && g.frame <= 8 {
		phase = "STARTUP"
		c = color.RGBA{255, 205, 70, 255}
	} else if g.frame >= 9 && g.frame <= 12 {
		phase = "ACTIVE"
		c = color.RGBA{255, 70, 75, 255}
		vector.StrokeRect(s, 235, 495, 105, 35, 4, c, false)
	} else if g.frame >= 13 {
		phase = "RECOVERY"
		c = color.RGBA{100, 155, 255, 255}
	}
	vector.DrawFilledRect(s, 50, 90, 380, 80, c, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  FRAME %02d / 30", phase, g.frame), 150, 120)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HITS %d/5", g.hits), 205, 200)
	ebitenutil.DebugPrintAt(s, "STARTUP 1-8   ACTIVE 9-12   RECOVERY 13-30", 85, 235)
	ebitenutil.DebugPrintAt(s, "PRESS SPACE / X / TAP — TIME THE MOVING TARGET", 75, 680)
	if g.clear {
		overlay(s, "FIVE PUNCHES LANDED!\n\nTAP / SPACE TO RESET")
	}
}
func press() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlap(a, b rect) bool { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Three-phase Punch — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
