package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
)

const W, H = 480, 720

type game struct {
	next, lap  int
	carX, carY float64
	move       *vfxmotion.Proxy
	fx         vfxfx.System
	tick       int
}

var gates = [][2]float32{{240, 120}, {390, 250}, {390, 500}, {240, 610}, {90, 500}, {90, 250}}

func (g *game) Update() error {
	g.tick++
	if g.move != nil {
		g.move.Advance()
		g.carX, g.carY = g.move.Position()
		if g.move.Done() {
			g.fx.Shockwave(g.carX, g.carY, .55, color.RGBA{255, 220, 90, 255}, color.RGBA{80, 235, 210, 255})
			g.move = nil
		}
		g.fx.Update()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		target := gates[g.next]
		proxy := vfxmotion.NewProxy(g.next, g.carX, g.carY, float64(target[0]), float64(target[1]), 22)
		g.move = &proxy
		g.next = (g.next + 1) % len(gates)
		if g.next == 0 {
			g.lap++
		}
	}
	g.fx.Update()
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{37, 100, 65, 255})
	for i, p := range gates {
		c := color.RGBA{80, 90, 105, 255}
		if i == g.next {
			c = color.RGBA{255, 205, 70, 255}
		}
		radius := float32(30)
		if i == g.next {
			radius += float32(4 + math.Sin(float64(g.tick)*.16)*3)
		}
		vector.StrokeCircle(s, p[0], p[1], radius, 6, c, true)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d", i+1), int(p[0])-4, int(p[1])-5)
	}
	vector.DrawFilledCircle(s, float32(g.carX), float32(g.carY), 12, color.RGBA{235, 91, 76, 255}, true)
	vector.StrokeCircle(s, float32(g.carX), float32(g.carY), 16, 2, color.RGBA{255, 255, 255, 180}, true)
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LAP %d  NEXT GATE %d", g.lap, g.next+1), 160, 50)
	if g.move == nil {
		ebitenutil.DebugPrintAt(s, "TAP / SPACE: PASS THE GLOWING GATE", 90, 680)
	} else {
		ebitenutil.DebugPrintAt(s, "COMMIT ONCE — PRESENT THE DRIVE", 105, 680)
	}
}
func (g *game) Layout(_, _ int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Lap Gates")
	start := gates[len(gates)-1]
	if err := ebiten.RunGame(&game{carX: float64(start[0]), carY: float64(start[1])}); err != nil {
		panic(err)
	}
}
