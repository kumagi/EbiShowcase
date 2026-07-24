// vfx-fx-snake — ADVANCED: grid snake; body[] is play, fx is separate.
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
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
	visualFrom  []cell
	stepTween   vfxmotion.Tween
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
	g.visualFrom = append([]cell(nil), g.body...)
	g.stepTween = vfxmotion.NewTween(int(g.shell.Get("spd")))
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
	oldBody := append([]cell(nil), g.body...)
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
	g.visualFrom = make([]cell, len(g.body))
	g.visualFrom[0] = oldBody[0]
	for i := 1; i < len(g.visualFrom); i++ {
		source := i - 1
		if source >= len(oldBody) {
			source = len(oldBody) - 1
		}
		g.visualFrom[i] = oldBody[source]
	}
	g.stepTween = vfxmotion.NewTween(int(g.shell.Get("spd")))
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
	if !g.stepTween.Done() {
		g.stepTween.Advance()
	}
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
	foodPulse := 0.32 + 0.06*mathSin(float64(g.tick)*0.5)
	vector.DrawFilledCircle(screen, float32(fx), float32(fy), float32(g.cs*foodPulse), color.RGBA{255, 90, 90, 255}, false)
	progress := vfxmotion.EaseInOutCubic(g.stepTween.Progress())
	positions := make([]vfxmotion.Point, len(g.body))
	for i, b := range g.body {
		from := b
		if i < len(g.visualFrom) {
			from = g.visualFrom[i]
		}
		x0, y0 := g.cellCenter(from)
		x1, y1 := g.cellCenter(b)
		positions[i] = vfxmotion.Point{X: vfxmotion.Lerp(x0, x1, progress), Y: vfxmotion.Lerp(y0, y1, progress)}
	}
	for i := len(positions) - 1; i > 0; i-- {
		vector.StrokeLine(screen, float32(positions[i].X), float32(positions[i].Y),
			float32(positions[i-1].X), float32(positions[i-1].Y),
			float32(g.cs*0.62), color.RGBA{70, 200, 125, 255}, false)
	}
	for i, point := range positions {
		c := color.RGBA{80, 220, 140, 255}
		if i == 0 {
			c = color.RGBA{120, 255, 180, 255}
		}
		radius := g.cs * 0.34
		if i > 0 && i+1 < len(g.body) {
			pose := vfxmotion.PoseForSegment(
				vfxmotion.Cell{X: g.body[i-1].x, Y: g.body[i-1].y},
				vfxmotion.Cell{X: g.body[i].x, Y: g.body[i].y},
				vfxmotion.Cell{X: g.body[i+1].x, Y: g.body[i+1].y},
			)
			if pose.Kind == vfxmotion.SegmentCorner {
				radius *= 1.1
			}
		}
		vector.DrawFilledCircle(screen, float32(point.X), float32(point.Y), float32(radius), c, false)
	}
	g.fx.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("score %d", g.score), 12, int(g.oy)-18)
	if g.over {
		ebitenutil.DebugPrintAt(screen, "GAME OVER — tap to retry", 140, int(g.oy)+g.rows*int(g.cs)/2)
	}
	g.shell.Hint = "arrows/WASD or tap stage to turn · eat → burst (body≠fx)"
	g.shell.Draw(screen)
}

func mathSin(v float64) float64 {
	// Small wrapper keeps the drawing formula readable beside grid logic.
	return math.Sin(v)
}

func (g *Game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Snake — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
