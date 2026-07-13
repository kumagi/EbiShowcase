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
	phase, tick, enemy, flash, shake int
	sparks                           []spark
}

func (g *game) Update() error {
	g.tick++
	if g.phase > 0 {
		g.phase--
		if g.phase == 10 {
			g.enemy -= 8
			g.flash = 7
			g.shake = 5
			for i := 0; i < 14; i++ {
				a := float64(i) * math.Pi / 7
				g.sparks = append(g.sparks, spark{240, 170, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 28})
			}
		}
	}
	if g.flash > 0 {
		g.flash--
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
	if g.phase == 0 && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0) {
		g.phase = 30
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{22, 29, 49, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 5
	}
	if g.flash > 0 {
		trackatlas.DrawTinted(s, "boss-crab", 240+ox, 170, 125, 1, 1, .3, 1)
	} else {
		trackatlas.DrawCentered(s, "boss-crab", 240+ox, 170, 125)
	}
	t := 0.0
	if g.phase > 0 {
		t = 1 - math.Abs(float64(g.phase-15))/15
	}
	cy := 540 - t*270
	scale := 70 + t*22
	vector.DrawFilledRect(s, float32(205+ox), float32(cy), float32(scale), 110, color.RGBA{220, 85, 73, 255}, false)
	trackatlas.DrawCentered(s, "card-attack", 240+ox, cy+35, 44)
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENEMY HP %d", g.enemy), 195, 50)
	ebitenutil.DebugPrintAt(s, "SELECT -> FLY -> CONTACT -> RETURN", 100, 400)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: PLAY CARD", 145, 675)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Card Motion Lab — Ebitengine")
	if err := ebiten.RunGame(&game{enemy: 40}); err != nil {
		panic(err)
	}
}
