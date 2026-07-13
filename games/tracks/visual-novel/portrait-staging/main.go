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

const (
	W = 480
	H = 720
)

var names = []string{"FADE", "SLIDE", "POP", "BOUNCE"}

type game struct{ kind, frames, seen int }

func (g *game) choose(k int) { g.kind = k; g.frames = 0; g.seen |= 1 << k }
func (g *game) Update() error {
	g.frames++
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if inpututil.IsKeyJustPressed(k) {
			g.choose(i)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y > 600 {
			g.choose(min(3, x/120))
		}
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y > 600 {
			g.choose(min(3, x/120))
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{27, 34, 70, 255})
	t := min(1, float64(g.frames)/45)
	x, y, scale := 240.0, 285.0, 1.0
	alpha := uint8(255)
	switch g.kind {
	case 0:
		alpha = uint8(t * 255)
	case 1:
		x = 560 - 320*t
	case 2:
		scale = .6 + .4*(1-math.Pow(1-t, 3))
	case 3:
		y -= math.Sin(t*math.Pi) * 70
	}
	vector.DrawFilledRect(s, float32(x-75*scale), float32(y), float32(150*scale), 220, color.RGBA{226, 105, 144, alpha}, false)
	vector.DrawFilledCircle(s, float32(x), float32(y-35), float32(70*scale), color.RGBA{248, 205, 172, alpha}, false)
	vector.DrawFilledCircle(s, float32(x-22), float32(y-45), 4, color.RGBA{40, 32, 52, alpha}, false)
	vector.DrawFilledCircle(s, float32(x+22), float32(y-45), 4, color.RGBA{40, 32, 52, alpha}, false)
	ebitenutil.DebugPrintAt(s, "PORTRAIT STAGING", 169, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  seen %d/4", names[g.kind], bits(g.seen)), 174, 85)
	if g.seen == 15 {
		ebitenutil.DebugPrintAt(s, "ALL ENTRANCES DISCOVERED!", 135, 555)
	}
	for i, n := range names {
		vector.DrawFilledRect(s, float32(i*120+4), 610, 112, 62, color.RGBA{53, 77, 120, 255}, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s", i+1, n), i*120+25, 637)
	}
}
func bits(v int) int {
	n := 0
	for v > 0 {
		n += v & 1
		v >>= 1
	}
	return n
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
