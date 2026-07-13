package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const W, H = 480, 720

type runner struct {
	name               string
	speed, gauge, wins int
	c                  color.RGBA
}
type game struct {
	rs   [3]runner
	tick int
}

func (g *game) Update() error {
	g.tick++
	for i := range g.rs {
		r := &g.rs[i]
		r.gauge += r.speed
		if r.gauge >= 1000 {
			r.wins++
			r.gauge = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.rs[0].speed = 5 + (g.rs[0].speed-4)%11
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{13, 24, 43, 255})
	ebitenutil.DebugPrintAt(s, "ACTION GAUGE RACE", 170, 65)
	for i, r := range g.rs {
		y := 190 + i*130
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  SPD %02d  READY %d", r.name, r.speed, r.wins), 70, y)
		vector.DrawFilledRect(s, 70, float32(y+38), 340, 30, color.RGBA{35, 49, 68, 255}, false)
		vector.DrawFilledRect(s, 70, float32(y+38), float32(340*r.gauge/1000), 30, r.c, false)
	}
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: change TENJIROH speed", 110, 630)
	ebitenutil.DebugPrintAt(s, "Gauge += speed every Update", 135, 665)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	g := &game{rs: [3]runner{{"TENJIROH", 8, 0, 0, color.RGBA{70, 165, 230, 255}}, {"MAGE", 13, 0, 0, color.RGBA{180, 95, 220, 255}}, {"SHELL", 5, 0, 0, color.RGBA{100, 195, 120, 255}}}}
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
