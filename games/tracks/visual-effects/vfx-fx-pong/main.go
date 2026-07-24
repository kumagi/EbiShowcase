// vfx-fx-pong — ADV 05: paddle/ball; sparks on hit + additive trail.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const (
	padHW = 52.0
	ballR = 10.0
)

type game struct {
	shell          *vfxlive.Shell
	fx             vfxfx.System
	px, cx         float64
	bx, by         float64
	vx, vy         float64
	pscore, cscore int
	frame          int
	placed         bool
	trail          vfxmotion.Trail
	playerKick     int
	cpuKick        int
}

func newGame() *game {
	return &game{
		px:    vfxlive.Width / 2,
		cx:    vfxlive.Width / 2,
		trail: vfxmotion.NewTrail(16),
		shell: vfxlive.New(
			"FX pong",
			[]string{
				"func (g *Game) Update() {",
				"  g.updatePlay() // paddle / bounce  score={sc}",
				"  g.fx.Update()  // sparks + trail   parts={n}",
				"}",
				"// hit → Burst; every few frames → trail spark",
			},
			&vfxlive.Param{Key: "fx", Label: "fx intensity", Value: 1.0, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "spd", Label: "ball speed", Value: 1.0, Min: 0.6, Max: 1.6, Format: "%.2f"},
		),
	}
}

func (g *game) ensureBall(dir float64) {
	if g.placed {
		return
	}
	_, sy, _, sh := g.shell.Stage()
	g.bx = vfxlive.Width / 2
	g.by = sy + sh/2
	g.vx, g.vy = 3.0, 3.8*dir
	g.trail.Clear()
	g.trail.Push(vfxmotion.Point{X: g.bx, Y: g.by})
	g.placed = true
}

func (g *game) spark(x, y float64) {
	k := g.shell.Get("fx")
	g.fx.Burst(x, y, int(10*k), 2.2*k, color.RGBA{120, 220, 255, 255}, true)
	g.fx.Ring(x, y, 0.35*k, color.RGBA{180, 235, 255, 255})
}

func (g *game) updatePlay() {
	g.ensureBall(1)
	_, sy, _, sh := g.shell.Stage()
	topY := sy + 28
	botY := sy + sh - 28

	if x, _, ok := vfxui.Held(); ok {
		g.px = x
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.px -= 7
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.px += 7
	}
	g.px = math.Max(padHW, math.Min(float64(vfxlive.Width)-padHW, g.px))
	if g.cx < g.bx-6 {
		g.cx += 3.2
	} else if g.cx > g.bx+6 {
		g.cx -= 3.2
	}

	sp := g.shell.Get("spd")
	g.bx += g.vx * sp
	g.by += g.vy * sp
	g.trail.Push(vfxmotion.Point{X: g.bx, Y: g.by})

	if g.bx < ballR || g.bx > float64(vfxlive.Width)-ballR {
		g.vx = -g.vx
		g.bx = math.Max(ballR, math.Min(float64(vfxlive.Width)-ballR, g.bx))
		g.spark(g.bx, g.by)
	}
	if g.vy > 0 && g.by >= botY-8 && g.by <= botY+10 && math.Abs(g.bx-g.px) < padHW+8 {
		g.by = botY - 8
		g.vy = -math.Abs(g.vy) * 1.02
		g.vx += (g.bx - g.px) * 0.05
		g.playerKick = 8
		g.spark(g.bx, g.by)
	}
	if g.vy < 0 && g.by <= topY+8 && g.by >= topY-10 && math.Abs(g.bx-g.cx) < padHW+8 {
		g.by = topY + 8
		g.vy = math.Abs(g.vy) * 1.01
		g.vx += (g.bx - g.cx) * 0.04
		g.cpuKick = 8
		g.spark(g.bx, g.by)
	}
	g.vx = math.Max(-7, math.Min(7, g.vx))

	g.frame++
	if g.playerKick > 0 {
		g.playerKick--
	}
	if g.cpuKick > 0 {
		g.cpuKick--
	}
	k := g.shell.Get("fx")
	if g.frame%4 == 0 {
		g.fx.Burst(g.bx, g.by, int(2*k), 0.6, color.RGBA{80, 180, 255, 180}, true)
	}

	if g.by < sy-20 {
		g.pscore++
		g.placed = false
		g.ensureBall(1)
	}
	if g.by > sy+sh+20 {
		g.cscore++
		g.placed = false
		g.ensureBall(-1)
	}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	if !ate {
		g.updatePlay()
	} else {
		g.ensureBall(1)
	}
	g.fx.Update()
	g.shell.SetToken("sc", fmt.Sprintf("%d-%d", g.pscore, g.cscore))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	g.shell.Hint = "bounce sparks + trail particles live in fx"
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{8, 20, 38, 255})
	_, sy, _, sh := g.shell.Stage()
	topY := sy + 28
	botY := sy + sh - 28
	cpuSquash := 1 + 0.12*math.Sin(float64(g.cpuKick)/8*math.Pi)
	playerSquash := 1 + 0.12*math.Sin(float64(g.playerKick)/8*math.Pi)
	vector.DrawFilledRect(s, float32(g.cx-padHW*cpuSquash), float32(topY), float32(padHW*2*cpuSquash), 12, color.RGBA{255, 120, 90, 255}, false)
	vector.DrawFilledRect(s, float32(g.px-padHW*playerSquash), float32(botY), float32(padHW*2*playerSquash), 12, color.RGBA{45, 226, 194, 255}, false)
	points := g.trail.Points()
	for i, point := range points {
		t := float64(i+1) / float64(max(1, len(points)))
		alpha := uint8(25 + 130*t)
		vector.DrawFilledCircle(s, float32(point.X), float32(point.Y), float32(2+ballR*t*0.65),
			color.RGBA{90, 190, 255, alpha}, false)
	}
	vector.DrawFilledCircle(s, float32(g.bx), float32(g.by), ballR, color.White, false)
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("YOU %d   CPU %d", g.pscore, g.cscore), 12, int(sy)+8)
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Pong — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
