// frame-attack — startup / active / recovery punch with clear feedback.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float32 }

type game struct {
	frame, hits, clock int
	hitThis, clear     bool
	flash              float64
	buffered           bool
}

func (g *game) wantPunch() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Update() error {
	if g.clear {
		if g.wantPunch() {
			*g = game{}
		}
		return nil
	}
	if g.flash > 0 {
		g.flash -= 0.08
	}
	g.clock++

	pressed := g.wantPunch()
	if pressed {
		if g.frame == 0 {
			g.frame = 1
			g.hitThis = false
			g.buffered = false
		} else {
			g.buffered = true
		}
	}
	if g.frame == 0 && g.buffered {
		g.frame = 1
		g.hitThis = false
		g.buffered = false
	}
	if g.frame > 0 {
		g.frame++
		if g.frame >= 31 {
			g.frame = 0
		}
	}
	if g.frame >= 9 && g.frame <= 12 && !g.hitThis {
		if overlap(g.attackBox(), g.target()) {
			g.hits++
			g.hitThis = true
			g.flash = 1
			if g.hits >= 5 {
				g.clear = true
			}
		}
	}
	return nil
}

func (g *game) attackBox() rect {
	reach := float32(0)
	if g.frame >= 9 && g.frame <= 12 {
		reach = 110
	} else if g.frame >= 5 && g.frame <= 8 {
		reach = float32(g.frame-4) * 18
	}
	return rect{200, 495, reach, 35}
}

func (g *game) target() rect {
	x := float32(340 + math.Sin(float64(g.clock)*.04)*65)
	return rect{x - 24, 470, 48, 115}
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{19, 28, 44, 255})
	if g.flash > 0 {
		a := g.flash
		vector.DrawFilledRect(s, 0, 0, width, height, color.RGBA{uint8(255 * a), uint8(90 * a), uint8(70 * a), uint8(140 * a)}, false)
	}
	vector.DrawFilledRect(s, 0, 590, 480, 130, color.RGBA{45, 54, 68, 255}, false)

	vector.DrawFilledCircle(s, 190, 485, 22, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 172, 507, 36, 78, color.RGBA{45, 225, 194, 255}, false)
	ab := g.attackBox()
	if ab.w > 0 {
		col := color.RGBA{255, 200, 80, 255}
		if g.frame >= 9 && g.frame <= 12 {
			col = color.RGBA{255, 70, 75, 255}
			vector.StrokeRect(s, ab.x, ab.y, ab.w, ab.h, 3, color.White, false)
		}
		vector.DrawFilledRect(s, ab.x, ab.y+6, ab.w, ab.h-12, col, false)
	}
	t := g.target()
	vector.DrawFilledRect(s, t.x, t.y, t.w, t.h, color.RGBA{240, 75, 91, 255}, false)

	phase := "READY — TAP / CLICK / SPACE TO PUNCH"
	c := color.RGBA{70, 90, 120, 255}
	switch {
	case g.frame >= 1 && g.frame <= 8:
		phase, c = "STARTUP — arm extending…", color.RGBA{255, 205, 70, 255}
	case g.frame >= 9 && g.frame <= 12:
		phase, c = "ACTIVE — HITBOX LIVE!", color.RGBA{255, 70, 75, 255}
	case g.frame >= 13:
		phase, c = "RECOVERY — wait for READY", color.RGBA{100, 155, 255, 255}
	}
	vector.DrawFilledRect(s, 24, 64, 432, 100, c, false)
	ebitenutil.DebugPrintAt(s, phase, 55, 90)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("FRAME %02d/30   HITS %d/5", g.frame, g.hits), 145, 120)
	ebitenutil.DebugPrintAt(s, "Hit the moving red foe while the red arm is ACTIVE", 55, 185)

	vector.DrawFilledRect(s, 140, 625, 200, 54, color.RGBA{40, 54, 82, 245}, false)
	vector.StrokeRect(s, 140, 625, 200, 54, 3, color.RGBA{120, 240, 220, 255}, false)
	ebitenutil.DebugPrintAt(s, "PUNCH", 210, 645)

	if g.clear {
		vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
		ebitenutil.DebugPrintAt(s, "FIVE PUNCHES LANDED!\n\nTAP / SPACE TO RESET", 110, 320)
	}
}

func overlap(a, b rect) bool {
	return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Three-phase Punch — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
