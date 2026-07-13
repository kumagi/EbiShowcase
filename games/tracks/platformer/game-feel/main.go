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
	x, y, vx, vy              float64
	grounded                  bool
	coins                     [8]bool
	tick, frames, best, shake int
	sparks                    []spark
	clear                     bool
}

func (g *game) reset() {
	g.x = 20
	g.y = 580
	g.vx = 0
	g.vy = 0
	g.coins = [8]bool{}
	g.frames = 0
	g.clear = false
	g.sparks = nil
}
func (g *game) burst(x, y float64) {
	for i := 0; i < 14; i++ {
		a := float64(i) * math.Pi * 2 / 14
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 30})
	}
}
func (g *game) Update() error {
	g.tick++
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .1
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.shake > 0 {
		g.shake--
	}
	if g.clear {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			g.reset()
		}
		return nil
	}
	g.frames++
	l := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	r := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	j := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				l = true
			} else {
				r = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			j = true
		}
	}
	if l {
		g.vx -= .75
	}
	if r {
		g.vx += .75
	}
	if !l && !r {
		g.vx *= .8
	}
	g.vx = math.Max(-7, math.Min(7, g.vx))
	if j && g.grounded {
		g.vy = -12.5
		g.burst(g.x+22, 625)
	}
	g.vy += .7
	g.x += g.vx
	g.y += g.vy
	if g.y >= 580 {
		g.y = 580
		g.vy = 0
		g.grounded = true
	}
	g.x = math.Max(0, math.Min(440, g.x))
	for i := range g.coins {
		cx := 45 + float64(i)*52
		if !g.coins[i] && math.Abs(g.x+20-cx) < 23 && math.Abs(g.y+20-(555-float64(i%2)*75)) < 35 {
			g.coins[i] = true
			g.burst(cx, 555-float64(i%2)*75)
			g.shake = 5
		}
	}
	if count(g.coins) == 8 {
		g.clear = true
		if g.best == 0 || g.frames < g.best {
			g.best = g.frames
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2) * 4
	}
	s.Fill(color.RGBA{28, 45, 86, 255})
	vector.DrawFilledRect(s, float32(ox), 630, width, 90, color.RGBA{39, 94, 69, 255}, false)
	for i, v := range g.coins {
		if !v {
			trackatlas.DrawCentered(s, "coin", ox+45+float64(i)*52, 555-float64(i%2)*75, 24)
		}
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(ox+p.x), float32(p.y), float32(2+p.life/12), color.RGBA{255, 213, 70, 255}, true)
	}
	trackatlas.DrawCentered(s, "hero", ox+g.x+20, g.y+24, 46)
	best := "--"
	if g.best > 0 {
		best = fmt.Sprintf("%.2f", float64(g.best)/60)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("COINS %d/8   TIME %.2f   BEST %s", count(g.coins), float64(g.frames)/60, best), 90, 28)
	ebitenutil.DebugPrintAt(s, "PARTICLES + SHAKE + BEST TIME", 125, 52)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH   JUMP: SPACE OR UPPER TOUCH", 45, 685)
	if g.clear {
		vector.DrawFilledRect(s, 65, 285, 350, 135, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("FINISH %.2f SEC!\n\nSPACE / TAP: BEAT YOUR BEST", float64(g.frames)/60), 125, 320)
	}
}
func count(a [8]bool) int {
	n := 0
	for _, v := range a {
		if v {
			n++
		}
	}
	return n
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	g := &game{}
	g.reset()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Game Feel Sprint — Ebitengine")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
