package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float64 }
type stage struct {
	name   string
	sky    color.RGBA
	blocks []rect
	goal   rect
}
type game struct {
	p        rect
	vx, vy   float64
	grounded bool
	stage    int
	stages   []stage
	clear    bool
}

func newGame() *game {
	g := &game{stages: []stage{
		{"STAIRS", color.RGBA{98, 185, 229, 255}, []rect{{0, 640, 480, 80}, {100, 540, 100, 20}, {230, 440, 100, 20}, {360, 340, 100, 20}}, rect{405, 285, 35, 55}},
		{"ISLANDS", color.RGBA{239, 151, 103, 255}, []rect{{0, 640, 140, 80}, {205, 555, 100, 20}, {355, 470, 100, 20}, {210, 360, 90, 20}, {55, 270, 100, 20}}, rect{85, 215, 35, 55}},
		{"TOWER", color.RGBA{48, 54, 108, 255}, []rect{{0, 640, 480, 80}, {30, 520, 100, 20}, {190, 430, 100, 20}, {350, 340, 100, 20}, {190, 250, 100, 20}, {30, 160, 100, 20}}, rect{65, 105, 35, 55}},
	}}
	g.load()
	return g
}
func (g *game) load() { g.p = rect{25, 590, 30, 48}; g.vx = 0; g.vy = 0; g.clear = false }
func (g *game) Update() error {
	if g.clear {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			g.stage++
			if g.stage >= len(g.stages) {
				g.stage = 0
			}
			g.load()
		}
		return nil
	}
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
		g.vx -= .7
	}
	if r {
		g.vx += .7
	}
	if !l && !r {
		g.vx *= .76
	}
	g.vx = math.Max(-5.5, math.Min(5.5, g.vx))
	if j && g.grounded {
		g.vy = -12
	}
	g.vy += .65
	g.p.x += g.vx
	g.p.x = math.Max(0, math.Min(width-g.p.w, g.p.x))
	old := g.p.y + g.p.h
	g.p.y += g.vy
	g.grounded = false
	for _, b := range g.stages[g.stage].blocks {
		if g.vy >= 0 && old <= b.y+3 && g.p.y+g.p.h >= b.y && g.p.x+g.p.w > b.x && g.p.x < b.x+b.w {
			g.p.y = b.y - g.p.h
			g.vy = 0
			g.grounded = true
		}
	}
	if g.p.y > height {
		g.load()
	}
	if overlap(g.p, g.stages[g.stage].goal) {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	st := g.stages[g.stage]
	s.Fill(st.sky)
	for _, b := range st.blocks {
		vector.DrawFilledRect(s, float32(b.x), float32(b.y), float32(b.w), float32(b.h), color.RGBA{48, 91, 65, 255}, false)
		vector.DrawFilledRect(s, float32(b.x), float32(b.y), float32(b.w), 6, color.RGBA{111, 218, 91, 255}, false)
	}
	trackatlas.Draw(s, "flag", st.goal.x-25, st.goal.y-45, 100)
	trackatlas.DrawCentered(s, "hero", g.p.x+15, g.p.y+24, 46)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE DATA %d/3: %s", g.stage+1, st.name), 145, 28)
	ebitenutil.DebugPrintAt(s, "SAME RULES + DIFFERENT DATA", 140, 52)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH   JUMP: SPACE OR UPPER TOUCH", 45, 685)
	if g.clear {
		vector.DrawFilledRect(s, 60, 285, 360, 125, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "STAGE CLEAR!\n\nSPACE / TAP: LOAD NEXT DATA", 130, 320)
	}
}
func overlap(a, b rect) bool               { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Stage Data Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
