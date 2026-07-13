package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const width, height = 480, 720

type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	enemy, tick, phase, flash, shake int
	sparks                           []spark
}

func (g *game) burst() {
	for i := 0; i < 16; i++ {
		a := float64(i) * math.Pi / 8
		g.sparks = append(g.sparks, spark{330, 250, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 30})
	}
}
func (g *game) Update() error {
	g.tick++
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.phase = 18
		g.flash = 8
		g.shake = 5
		g.enemy = max(0, g.enemy-12)
		g.burst()
	}
	if g.phase > 0 {
		g.phase--
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{26, 27, 52, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 5
	}
	heroX := 130.0
	if g.phase > 8 {
		heroX += float64(18-g.phase) * 12
	} else if g.phase > 0 {
		heroX += float64(g.phase) * 12
	}
	hero.DrawCentered(s, heroX+ox, 360, 65)
	if g.flash > 0 {
		trackatlas.DrawTinted(s, "boss-crab", 330+ox, 250, 130, 1, 1, .25, 1)
	} else {
		trackatlas.DrawCentered(s, "boss-crab", 330+ox, 250, 130)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/15), color.RGBA{255, 211, 62, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENEMY HP %d/60", g.enemy), 180, 60)
	ebitenutil.DebugPrintAt(s, "ANTICIPATION -> CONTACT -> RECOVERY", 95, 460)
	vector.DrawFilledRect(s, 70, 540, 340, 90, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: ATTACK", 155, 575)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Battle Animation Lab — Ebitengine")
	if err := ebiten.RunGame(&game{enemy: 60}); err != nil {
		panic(err)
	}
}
