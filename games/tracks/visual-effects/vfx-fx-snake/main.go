// vfx-fx-snake — ADVANCED: grid snake; body[] is play, fx is separate.
package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type cell struct{ x, y int }

type Game struct {
	shell       *vfxlive.Shell
	fx          vfxfx.System
	body        []cell
	dir, next   cell
	food        cell
	tick, score int
	over        bool
	rng         *rand.Rand
	ox, oy      float64
	cols, rows  int
	cs          float64
}

func newGame() *Game {
	g := &Game{
		rng:  rand.New(rand.NewSource(71)),
		dir:  cell{1, 0},
		next: cell{1, 0},
		body: []cell{{4, 5}, {3, 5}, {2, 5}},
		shell: vfxlive.New(
			"FX Snake",
			[]string{
				"type Game struct { body []cell; fx System }",
				"func (g *Game) Update() {",
				"  g.updatePlay() // body grows on eat",
				"  g.fx.Update()  // particles fade",
				"}",
				"body={len}  fx.parts={n}",
			},
			&vfxlive.Param{Key: "fx", Label: "fx", Value: 1, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "spd", Label: "tick", Value: 8, Min: 4, Max: 16, Step: 1, Format: "%.0f"},
		),
	}
	g.layoutGrid()
	g.placeFood()
	return g
}

func (g *Game) layoutGrid() {
	_, sy, sw, sh := g.shell.Stage()
	g.cs = 20
	g.cols = int(sw / g.cs)
	g.rows = int(sh / g.cs)
	g.ox = (sw - float64(g.cols)*g.cs) / 2
	g.oy = sy + (sh-float64(g.rows)*g.cs)/2
}

func (g *Game) placeFood() {
	for {
		c := cell{g.rng.Intn(g.cols), g.rng.Intn(g.rows)}
		ok := true
		for _, b := range g.body {
			if b == c {
				ok = false
				break
			}
		}
		if ok {
			g.food = c
			return
		}
	}
}

func (g *Game) setDir(d cell) {
	if d.x+g.dir.x == 0 && d.y+g.dir.y == 0 {
		return
	}
	g.next = d
}

func (g *Game) readInput(ate bool) {
	if ate {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.setDir(cell{0, -1})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.setDir(cell{0, 1})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.setDir(cell{-1, 0})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.setDir(cell{1, 0})
	}
	if x, y, ok := vfxui.JustPressed(); ok {
		_, sy, _, sh := g.shell.Stage()
		if y < sy || y > sy+sh {
			return
		}
		mx, my := x-g.ox, y-g.oy
		cx, cy := float64(g.cols)*g.cs/2, float64(g.rows)*g.cs/2
		dx, dy := mx-cx, my-cy
		if absF(dx) > absF(dy) {
			if dx < 0 {
				g.setDir(cell{-1, 0})
			} else {
				g.setDir(cell{1, 0})
			}
		} else {
			if dy < 0 {
				g.setDir(cell{0, -1})
			} else {
				g.setDir(cell{0, 1})
			}
		}
	}
}

func absF(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func (g *Game) cellCenter(c cell) (float64, float64) {
	return g.ox + (float64(c.x)+0.5)*g.cs, g.oy + (float64(c.y)+0.5)*g.cs
}

func (g *Game) updatePlay() {
	g.tick++
	if g.tick < int(g.shell.Get("spd")) {
		return
	}
	g.tick = 0
	g.dir = g.next
	h := g.body[0]
	nh := cell{h.x + g.dir.x, h.y + g.dir.y}
	if nh.x < 0 || nh.y < 0 || nh.x >= g.cols || nh.y >= g.rows {
		g.over = true
		return
	}
	for _, b := range g.body {
		if b == nh {
			g.over = true
			return
		}
	}
	g.body = append([]cell{nh}, g.body...)
	if nh == g.food {
		g.score++
		fx := g.shell.Get("fx")
		cx, cy := g.cellCenter(nh)
		g.fx.Burst(cx, cy, int(18*fx), 2.8*fx, color.RGBA{255, 200, 80, 255}, true)
		g.fx.Ring(cx, cy, 0.9*fx, color.RGBA{255, 220, 120, 255})
		g.placeFood()
	} else {
		g.body = g.body[:len(g.body)-1]
	}
}

func (g *Game) Update() error {
	ate := g.shell.Update()
	g.layoutGrid()
	g.shell.SetToken("len", fmt.Sprintf("%d", len(g.body)))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	if g.over {
		if vfxui.AnyPressStart() && !ate {
			*g = *newGame()
		}
		g.fx.Update()
		return nil
	}
	g.readInput(ate)
	g.updatePlay()
	g.fx.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 14, 28, 255})
	g.shell.FillStage(screen, color.RGBA{14, 22, 40, 255})
	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.cols; x++ {
			if (x+y)%2 == 0 {
				vector.DrawFilledRect(screen, float32(g.ox+float64(x)*g.cs), float32(g.oy+float64(y)*g.cs),
					float32(g.cs), float32(g.cs), color.RGBA{18, 28, 48, 255}, false)
			}
		}
	}
	fx, fy := g.cellCenter(g.food)
	vector.DrawFilledCircle(screen, float32(fx), float32(fy), float32(g.cs*0.35), color.RGBA{255, 90, 90, 255}, false)
	for i, b := range g.body {
		c := color.RGBA{80, 220, 140, 255}
		if i == 0 {
			c = color.RGBA{120, 255, 180, 255}
		}
		vector.DrawFilledRect(screen, float32(g.ox+float64(b.x)*g.cs+2), float32(g.oy+float64(b.y)*g.cs+2),
			float32(g.cs-4), float32(g.cs-4), c, false)
	}
	g.fx.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("score %d", g.score), 12, int(g.oy)-18)
	if g.over {
		ebitenutil.DebugPrintAt(screen, "GAME OVER — tap to retry", 140, int(g.oy)+g.rows*int(g.cs)/2)
	}
	g.shell.Hint = "arrows/WASD or tap stage to turn · eat → burst (body≠fx)"
	g.shell.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Snake — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
