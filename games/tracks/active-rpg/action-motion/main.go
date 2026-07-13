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

const W, H = 480, 720

type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	timer, hp, shake, tick int
	sparks                 []spark
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
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.timer > 0 {
		g.timer--
		if g.timer == 16 {
			g.hp -= 12
			g.shake = 9
			for i := 0; i < 18; i++ {
				a := float64(i) * math.Pi / 9
				g.sparks = append(g.sparks, spark{350, 280, math.Cos(a) * 3, math.Sin(a) * 3, 28})
			}
		}
	}
	if g.timer == 0 && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0) {
		if g.hp <= 0 {
			g.hp = 60
		}
		g.timer = 38
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{22, 24, 46, 255})
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2.5) * 6
	}
	p := 0.0
	if g.timer > 0 {
		p = 1 - math.Abs(float64(g.timer-19))/19
	}
	trackatlas.DrawCentered(s, "hero", 120+p*150+ox, 280, 105)
	trackatlas.DrawCentered(s, "boss-crab", 350+ox, 280, 120)
	for _, q := range g.sparks {
		vector.DrawFilledCircle(s, float32(q.x+ox), float32(q.y), float32(2+q.life/8), color.RGBA{255, 195, 65, 255}, true)
	}
	state := "READY"
	if g.timer > 25 {
		state = "ANTICIPATION"
	} else if g.timer > 14 {
		state = "CONTACT"
	} else if g.timer > 0 {
		state = "RECOVERY"
	}
	ebitenutil.DebugPrintAt(s, "ACTION MOTION", 180, 65)
	ebitenutil.DebugPrintAt(s, state, 200, 455)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENEMY HP %d", max(0, g.hp)), 185, 500)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: ATTACK", 150, 650)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{hp: 60}); err != nil {
		panic(err)
	}
}
