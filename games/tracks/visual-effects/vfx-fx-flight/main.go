// vfx-fx-flight — ADV 04: flappy-like; gravity play vs FX on flap/score.
package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const (
	birdX   = 110.0
	birdR   = 14.0
	pipeW   = 54.0
	gravity = 0.38
	flapVY  = -6.8
)

type pipe struct {
	x, gapY float64
	scored  bool
}

type game struct {
	shell    *vfxlive.Shell
	fx       vfxfx.System
	rng      *rand.Rand
	birdY    float64
	vy       float64
	pipes    []pipe
	score    int
	started  bool
	gameOver bool
	ready    bool
}

func newGame() *game {
	g := &game{
		rng: rand.New(rand.NewSource(7)),
		shell: vfxlive.New(
			"FX flight",
			[]string{
				"func (g *Game) Update() {",
				"  g.updatePlay() // gravity / pipes  score={sc}",
				"  g.fx.Update()  // whoosh / rings   parts={n}",
				"}",
				"// gameplay gravity ≠ particle sparks",
			},
			&vfxlive.Param{Key: "fx", Label: "fx intensity", Value: 1.0, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "gap", Label: "pipe gap", Value: 150, Min: 110, Max: 200, Step: 5, Format: "%.0f"},
		),
	}
	return g
}

func (g *game) reset() {
	_, sy, _, sh := g.shell.Stage()
	g.birdY = sy + sh*0.4
	g.vy = 0
	g.score = 0
	g.started = false
	g.gameOver = false
	g.fx = vfxfx.System{}
	base := float64(vfxlive.Width) + 40
	g.pipes = []pipe{
		{x: base, gapY: sy + sh*0.45},
		{x: base + 200, gapY: sy + sh*0.55},
		{x: base + 400, gapY: sy + sh*0.4},
	}
	g.ready = true
}

func (g *game) stagePress() bool {
	_, sy, _, sh := g.shell.Stage()
	if _, y, ok := vfxui.JustPressed(); ok && y >= sy && y <= sy+sh {
		return true
	}
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp)
}

func (g *game) hitPipe(gap float64) bool {
	for _, p := range g.pipes {
		if birdX+birdR <= p.x || birdX-birdR >= p.x+pipeW {
			continue
		}
		top, bot := p.gapY-gap/2, p.gapY+gap/2
		if g.birdY-birdR < top || g.birdY+birdR > bot {
			return true
		}
	}
	return false
}

func (g *game) updatePlay() {
	if !g.ready {
		g.reset()
	}
	_, sy, _, sh := g.shell.Stage()
	ground := sy + sh - 6
	gap := g.shell.Get("gap")
	k := g.shell.Get("fx")
	pressed := g.stagePress()

	if g.gameOver {
		if pressed {
			g.reset()
		}
		return
	}
	if !g.started {
		if pressed {
			g.started = true
			g.vy = flapVY
			g.fx.Burst(birdX, g.birdY, int(10*k), 2.2*k, color.RGBA{180, 220, 255, 255}, true)
		}
		return
	}
	if pressed {
		g.vy = flapVY
		g.fx.Burst(birdX-6, g.birdY+8, int(12*k), 2.4*k, color.RGBA{160, 210, 255, 255}, true)
	}
	g.vy += gravity
	g.birdY += g.vy

	for i := range g.pipes {
		p := &g.pipes[i]
		p.x -= 2.6
		if !p.scored && p.x+pipeW < birdX {
			p.scored = true
			g.score++
			g.fx.Ring(birdX+20, g.birdY, 0.85*k, color.RGBA{120, 240, 220, 255})
			g.fx.FlashScreen(0.25*k, 80, 200, 220)
		}
		if p.x < -pipeW {
			p.x += 600
			p.gapY = sy + 50 + g.rng.Float64()*(sh-100)
			p.scored = false
		}
	}
	if g.birdY > ground-birdR || g.birdY < sy+birdR || g.hitPipe(gap) {
		g.gameOver = true
		g.fx.Burst(birdX, g.birdY, int(14*k), 2.0, color.RGBA{255, 120, 100, 255}, true)
	}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	if !ate {
		g.updatePlay()
	} else if !g.ready {
		g.reset()
	}
	g.fx.Update()
	g.shell.SetToken("sc", fmt.Sprintf("%d", g.score))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	g.shell.Hint = "gravity is play; sparks are fx — Update both"
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{12, 28, 48, 255})
	_, sy, _, sh := g.shell.Stage()
	gap := g.shell.Get("gap")
	for _, p := range g.pipes {
		top, bot := p.gapY-gap/2, p.gapY+gap/2
		vector.DrawFilledRect(s, float32(p.x), float32(sy), pipeW, float32(top-sy), color.RGBA{45, 180, 120, 255}, false)
		vector.DrawFilledRect(s, float32(p.x), float32(bot), pipeW, float32(sy+sh-bot), color.RGBA{45, 180, 120, 255}, false)
	}
	hero.DrawCentered(s, birdX, g.birdY, 42)
	g.fx.Draw(s)
	msg := fmt.Sprintf("SCORE %d", g.score)
	if !g.started {
		msg = "TAP TO FLAP"
	}
	if g.gameOver {
		msg = fmt.Sprintf("CRASH  SCORE %d — tap retry", g.score)
	}
	ebitenutil.DebugPrintAt(s, msg, 12, int(sy)+8)
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Flight — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
