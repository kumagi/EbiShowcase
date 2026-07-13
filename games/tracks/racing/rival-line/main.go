package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const W, H = 480, 720

type game struct {
	x, y, a   float64
	next, lap int
}

var path = [][2]float64{{240, 110}, {390, 220}, {390, 500}, {240, 610}, {90, 500}, {90, 220}}

func (g *game) Update() error {
	q := path[g.next]
	target := math.Atan2(q[1]-g.y, q[0]-g.x) + math.Pi/2
	d := math.Mod(target-g.a+math.Pi, 2*math.Pi) - math.Pi
	g.a += math.Max(-.025, math.Min(.025, d))
	g.x += math.Sin(g.a) * 3
	g.y -= math.Cos(g.a) * 3
	if math.Hypot(g.x-q[0], g.y-q[1]) < 35 {
		g.next = (g.next + 1) % len(path)
		if g.next == 0 {
			g.lap++
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{31, 75, 87, 255})
	for i := 0; i < len(path); i++ {
		a := path[i]
		b := path[(i+1)%len(path)]
		vector.StrokeLine(s, float32(a[0]), float32(a[1]), float32(b[0]), float32(b[1]), 4, color.RGBA{255, 255, 255, 70}, true)
	}
	vector.DrawFilledRect(s, float32(g.x)-12, float32(g.y)-19, 24, 38, color.RGBA{76, 166, 232, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("AI LAP %d  TARGET GATE %d", g.lap, g.next+1), 145, 50)
	ebitenutil.DebugPrintAt(s, "ANGLE DIFF -> LIMITED STEER -> MOVE", 100, 680)
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Rival Racing Line")
	if err := ebiten.RunGame(&game{x: 240, y: 145}); err != nil {
		panic(err)
	}
}
