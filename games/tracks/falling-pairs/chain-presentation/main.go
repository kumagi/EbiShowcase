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

type spark struct {
	x, y, vx, vy, life float64
	kind               int
}
type game struct {
	phase, tick, chain, score, shake int
	sparks                           []spark
}

func (g *game) start() { g.phase = 90; g.chain = 1 }
func (g *game) burst(chain int) {
	g.shake = 4 + chain*3
	g.score += 40 * chain
	for i := 0; i < 22; i++ {
		a := float64(i) * math.Pi / 11
		g.sparks = append(g.sparks, spark{240, 310, math.Cos(a) * (2 + float64(chain)*.4), math.Sin(a) * (2 + float64(chain)*.4), 30, i % 4})
	}
}
func (g *game) Update() error {
	g.tick++
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .05
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.phase > 0 {
		g.phase--
		if g.phase == 65 {
			g.burst(1)
		}
		if g.phase == 30 {
			g.chain = 2
			g.burst(2)
		}
	}
	if g.phase == 0 && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0) {
		g.start()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 24, 48, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2.4) * float64(g.shake)
	}
	vector.DrawFilledRect(s, 90, 160, 300, 330, color.RGBA{28, 47, 69, 255}, false)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			kind := (x + y) % 4
			scale := 54.0
			if g.phase > 55 && g.phase < 70 || g.phase > 20 && g.phase < 35 {
				scale *= 1 + math.Sin(float64(g.tick)*.7)*.12
			}
			trackatlas.DrawCentered(s, trackatlas.Gem(kind), 150+float64(x)*60+ox, 240+float64(y)*60, scale)
		}
	}
	for _, p := range g.sparks {
		c := []color.RGBA{{240, 90, 85, 255}, {70, 165, 235, 255}, {245, 185, 60, 255}, {100, 195, 120, 255}}[p.kind]
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/10), c, true)
	}
	state := "READY"
	if g.phase > 65 {
		state = "ANTICIPATION"
	} else if g.phase > 30 {
		state = "1 CHAIN / IMPACT"
	} else if g.phase > 0 {
		state = "2 CHAIN / BIG IMPACT"
	}
	ebitenutil.DebugPrintAt(s, "CHAIN PRESENTATION LAB", 145, 70)
	ebitenutil.DebugPrintAt(s, state, 155, 530)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("CHAIN %d   SCORE %04d", g.chain, g.score), 155, 575)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: PLAY A 2-CHAIN", 120, 665)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Chain Presentation — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
