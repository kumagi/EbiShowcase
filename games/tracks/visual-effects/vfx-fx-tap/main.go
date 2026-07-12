// vfx-fx-tap — ADV 01: tap targets; FX spawn on hit/miss, then fx.Update().
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const startSec = 25

type game struct {
	shell      *vfxlive.Shell
	fx         vfxfx.System
	rng        *rand.Rand
	cx, cy, r  float64
	score      int
	framesLeft int
	over       bool
}

func newGame() *game {
	g := &game{
		rng: rand.New(rand.NewSource(11)),
		r:   36,
		shell: vfxlive.New(
			"FX on tap",
			[]string{
				"func (g *Game) Update() {",
				"  g.updatePlay()   // score / timer / targets",
				"  g.fx.Update()    // separate FX layer",
				"}",
				"// on hit: Burst+Ring+Flash  fx.parts={n}",
			},
			&vfxlive.Param{Key: "fx", Label: "fx intensity", Value: 1.0, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "time", Label: "seconds", Value: float64(startSec), Min: 10, Max: 40, Step: 1, Format: "%.0f"},
		),
	}
	g.reset()
	return g
}

func (g *game) reset() {
	g.score = 0
	g.framesLeft = int(g.shell.Get("time")) * 60
	g.over = false
	g.r = 36
	g.fx = vfxfx.System{}
	g.moveTarget()
}

func (g *game) moveTarget() {
	_, sy, _, sh := g.shell.Stage()
	g.cx = 50 + g.rng.Float64()*(vfxlive.Width-100)
	g.cy = sy + 40 + g.rng.Float64()*(sh-80)
}

func (g *game) intensity() float64 { return g.shell.Get("fx") }

func (g *game) updatePlay() {
	if g.over {
		if x, y, ok := vfxui.JustPressed(); ok {
			_, sy, _, sh := g.shell.Stage()
			if y >= sy && y <= sy+sh {
				_ = x
				g.reset()
			}
		}
		return
	}
	g.framesLeft--
	if g.framesLeft <= 0 {
		g.over = true
		return
	}
	x, y, ok := vfxui.JustPressed()
	if !ok {
		return
	}
	_, sy, _, sh := g.shell.Stage()
	if y < sy || y > sy+sh {
		return
	}
	hit := math.Hypot(x-g.cx, y-g.cy) <= g.r
	k := g.intensity()
	if hit {
		g.score++
		g.r = math.Max(20, 36-float64(g.score)*0.4)
		n := int(18 * k)
		g.fx.Burst(g.cx, g.cy, n, 2.8*k, color.RGBA{120, 240, 220, 255}, true)
		g.fx.Ring(g.cx, g.cy, 0.9*k, color.RGBA{180, 255, 240, 255})
		g.fx.FlashScreen(0.35*k, 80, 200, 180)
		g.moveTarget()
	} else {
		g.fx.Burst(x, y, int(6*k), 1.2, color.RGBA{140, 140, 150, 255}, false)
	}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	if !ate {
		g.updatePlay()
	}
	g.fx.Update() // separate
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	g.shell.Hint = "play updates score; FX spawn on events, then fx.Update()"
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{8, 18, 36, 255})
	pulse := 0.75 + 0.25*math.Sin(float64(g.framesLeft)*0.12)
	vector.DrawFilledCircle(s, float32(g.cx), float32(g.cy), float32(g.r*pulse), color.RGBA{45, 226, 194, 220}, false)
	vector.StrokeCircle(s, float32(g.cx), float32(g.cy), float32(g.r+6), 2, color.RGBA{120, 240, 220, 160}, false)
	g.fx.Draw(s)
	msg := fmt.Sprintf("SCORE %d   TIME %.1fs", g.score, float64(g.framesLeft)/60)
	if g.over {
		msg = fmt.Sprintf("DONE  SCORE %d  — tap stage to retry", g.score)
	}
	_, sy, _, _ := g.shell.Stage()
	ebitenutil.DebugPrintAt(s, msg, 12, int(sy)+8)
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Tap — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
