package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const width, height = 480, 720

type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	fuse, blast, tick, shake int
	sparks                   []spark
}

func (g *game) Update() error {
	g.tick++
	if g.fuse > 0 {
		g.fuse--
		if g.fuse == 0 {
			g.blast = 30
			g.shake = 7
			for i := 0; i < 24; i++ {
				a := float64(i) * math.Pi / 12
				g.sparks = append(g.sparks, spark{240, 330, math.Cos(a) * float64(1+i%4), math.Sin(a) * float64(1+i%4), 30})
			}
		}
	}
	if g.blast > 0 {
		g.blast--
	}
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.fuse == 0 && g.blast == 0 && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0) {
		g.fuse = 90
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 24, 42, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 6
	}
	if g.fuse > 0 {
		size := 38 + math.Sin(float64(g.fuse)*.35)*5
		trackatlas.DrawCentered(s, "bomb", 240+ox, 330, size)
	}
	if g.blast > 0 {
		for _, d := range [][2]float64{{0, 0}, {48, 0}, {-48, 0}, {0, 48}, {0, -48}} {
			trackatlas.DrawCentered(s, "flame", 240+d[0]+ox, 330+d[1], 42+math.Sin(float64(g.blast)*.8)*6)
		}
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 175, 65, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("FUSE %d  BLAST %d", g.fuse, g.blast), 170, 55)
	ebitenutil.DebugPrintAt(s, "PULSE -> FLASH -> PARTICLES -> SHAKE", 90, 470)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: ARM BOMB", 150, 675)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Blast Feedback — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
