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
	round, hp, tick, freeze, shake int
	sparks                         []spark
}

func newGame() *game { return &game{round: 1, hp: 100} }
func (g *game) Update() error {
	g.tick++
	if g.freeze > 0 {
		g.freeze--
		return nil
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
		g.hp -= 20
		g.freeze = 7
		g.shake = 6
		for i := 0; i < 18; i++ {
			a := float64(i) * math.Pi / 9
			g.sparks = append(g.sparks, spark{330, 340, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 28})
		}
		if g.hp <= 0 {
			g.round++
			if g.round > 3 {
				g.round = 1
			}
			g.hp = 100
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	bg := []color.RGBA{{24, 33, 52, 255}, {52, 28, 55, 255}, {26, 61, 58, 255}}
	s.Fill(bg[g.round-1])
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 6
	}
	trackatlas.DrawCentered(s, "fighter-p1", 130+ox, 520, 130)
	trackatlas.DrawCentered(s, "fighter-p2", 340+ox, 520, 130)
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/15), color.RGBA{255, 211, 62, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ROUND %d/3   TARGET HP %d", g.round, g.hp), 145, 40)
	ebitenutil.DebugPrintAt(s, "HITSTOP + PARTICLES + SHAKE", 120, 75)
	vector.DrawFilledRect(s, 70, 600, 340, 70, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: LAND A HIT", 140, 630)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Round Presentation — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
