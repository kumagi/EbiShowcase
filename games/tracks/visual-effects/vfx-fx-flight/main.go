// vfx-fx-flight — ADV 04: flappy-like; gravity play vs FX on flap/score.
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
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
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
	shell     *vfxlive.Shell
	fx        vfxfx.System
	rng       *rand.Rand
	birdY     float64
	vy        float64
	pipes     []pipe
	score     int
	started   bool
	gameOver  bool
	ready     bool
	sinceFlap int
	frame     int
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
	g.sinceFlap = 99
	g.frame = 0
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
	g.frame++
	g.sinceFlap++

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
			g.sinceFlap = 0
			g.fx.Burst(birdX, g.birdY, int(10*k), 2.2*k, color.RGBA{180, 220, 255, 255}, true)
		}
		return
	}
	if pressed {
		g.vy = flapVY
		g.sinceFlap = 0
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
			cx, cy := birdX+18, g.birdY
			g.fx.Shockwave(cx, cy, 1.05*k, color.RGBA{120, 255, 230, 255}, color.RGBA{255, 230, 120, 255})
			g.fx.Burst(cx, cy, int(36*k), 3.6*k, color.RGBA{100, 255, 220, 255}, true)
			g.fx.Burst(cx, cy, int(18*k), 2.4*k, color.RGBA{255, 220, 100, 255}, true)
			g.fx.Confetti(cx, cy, int(12*k))
			g.fx.FlashScreen(0.55*k, 90, 230, 210)
		}
		if p.x < -pipeW {
			p.x += 600
			p.gapY = sy + 50 + g.rng.Float64()*(sh-100)
			p.scored = false
		}
	}
	if g.birdY > ground-birdR || g.birdY < sy+birdR || g.hitPipe(gap) {
		g.gameOver = true
		g.fx.Shockwave(birdX, g.birdY, 0.7*k, color.White, color.RGBA{255, 100, 90, 255})
		g.fx.Burst(birdX, g.birdY, int(18*k), 2.4, color.RGBA{255, 120, 100, 255}, true)
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
	for i := 0; i < 6; i++ {
		x := math.Mod(float64(g.frame*3+i*91), float64(vfxlive.Width))
		y := sy + 24 + float64((i*53)%int(math.Max(1, sh-48)))
		vector.StrokeLine(s, float32(x), float32(y), float32(x-18), float32(y), 1, color.RGBA{150, 210, 235, 80}, false)
	}
	_, pose := vfxmotion.PoseForFlight(g.vy, g.sinceFlap, g.gameOver)
	hero.DrawCenteredPose(s, birdX, g.birdY, 42, hero.Pose{
		ScaleX: pose.ScaleX, ScaleY: pose.ScaleY, Rotation: pose.Rotation,
	})
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
