package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
)

const W, H = 480, 720

type game struct {
	room, fromRoom int
	seen           [8]bool
	move           vfxmotion.Tween
	fx             vfxfx.System
}

func (g *game) Update() error {
	if g.move.Done() {
		next := g.room
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			next = (g.room + 1) % 8
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			next = (g.room + 7) % 8
		}
		if next != g.room {
			g.fromRoom = g.room
			g.room = next
			g.move = vfxmotion.NewTween(14)
			if !g.seen[g.room] {
				g.seen[g.room] = true
				x, y := roomCenter(g.room)
				g.fx.Burst(x, y, 18, 2.5, color.RGBA{90, 235, 205, 255}, true)
				g.fx.Ring(x, y, .7, color.RGBA{255, 215, 90, 255})
			}
		}
	}
	g.move.Advance()
	g.fx.Update()
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{13, 24, 42, 255})
	for i := 0; i < 8; i++ {
		x := 35 + float32(i%4)*110
		y := 180 + float32(i/4)*150
		c := color.RGBA{42, 48, 65, 255}
		if g.seen[i] {
			c = color.RGBA{75, 158, 147, 255}
		}
		vector.DrawFilledRect(s, x, y, 90, 110, c, false)
		if i == g.room {
			vector.StrokeRect(s, x-4, y-4, 98, 118, 2, color.RGBA{255, 205, 70, 120}, true)
		}
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("R%d", i+1), int(x)+35, int(y)+45)
	}
	fromX, fromY := roomCenter(g.fromRoom)
	toX, toY := roomCenter(g.room)
	progress := vfxmotion.EaseInOutCubic(g.move.Progress())
	markerX := vfxmotion.Lerp(fromX, toX, progress)
	markerY := vfxmotion.Lerp(fromY, toY, progress)
	vector.DrawFilledCircle(s, float32(markerX), float32(markerY), 11, color.RGBA{255, 205, 70, 255}, true)
	vector.StrokeCircle(s, float32(markerX), float32(markerY), 16, 2, color.White, true)
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DISCOVERED %d/8", count(g.seen)), 175, 60)
	ebitenutil.DebugPrintAt(s, "LEFT/RIGHT OR TAP: EXPLORE", 135, 670)
}
func roomCenter(room int) (float64, float64) {
	return 80 + float64(room%4)*110, 235 + float64(room/4)*150
}
func count(a [8]bool) int {
	n := 0
	for _, v := range a {
		if v {
			n++
		}
	}
	return n
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	g := &game{}
	g.seen[0] = true
	g.move = vfxmotion.NewTween(1)
	g.move.Advance()
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Exploration Map")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
