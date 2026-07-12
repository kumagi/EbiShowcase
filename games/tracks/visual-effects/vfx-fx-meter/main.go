// vfx-fx-meter — ADV 02: stop the meter; FX grade by accuracy.
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
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type game struct {
	shell   *vfxlive.Shell
	fx      vfxfx.System
	marker  float64
	speed   float64
	score   int
	round   int
	stopped bool
	label   string
}

func newGame() *game {
	g := &game{
		marker: vfxlive.Width / 2,
		speed:  3.4,
		shell: vfxlive.New(
			"FX meter stop",
			[]string{
				"func (g *Game) Update() {",
				"  g.updatePlay()  // marker / score  state={st}",
				"  g.fx.Update()   // FX lifetime     parts={n}",
				"}",
				"// perfect → big Burst+Flash; miss → gray puff",
			},
			&vfxlive.Param{Key: "fx", Label: "fx intensity", Value: 1.0, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "spd", Label: "meter speed", Value: 3.4, Min: 1.5, Max: 7, Format: "%.1f"},
		),
	}
	return g
}

func (g *game) updatePlay() {
	_, sy, _, sh := g.shell.Stage()
	cx := float64(vfxlive.Width) / 2
	barY := sy + sh*0.45

	if x, y, ok := vfxui.JustPressed(); ok {
		if y < sy || y > sy+sh {
			return
		}
		_ = x
		if g.stopped {
			g.stopped = false
			g.round++
			g.label = ""
			base := g.shell.Get("spd")
			g.speed = base + float64(g.round)*0.12
			if g.round%2 == 1 {
				g.speed = -g.speed
			}
			return
		}
		g.stopped = true
		d := math.Abs(g.marker - cx)
		k := g.shell.Get("fx")
		switch {
		case d <= 10:
			g.score += 100
			g.label = "PERFECT"
			g.fx.Burst(g.marker, barY, int(28*k), 3.5*k, color.RGBA{120, 255, 200, 255}, true)
			g.fx.Ring(g.marker, barY, 1.2*k, color.RGBA{180, 255, 220, 255})
			g.fx.FlashScreen(0.55*k, 60, 220, 180)
		case d <= 40:
			g.score += 40
			g.label = "OK"
			g.fx.Burst(g.marker, barY, int(12*k), 2.0*k, color.RGBA{255, 210, 90, 255}, true)
		default:
			g.label = "MISS"
			g.fx.Burst(g.marker, barY, int(8*k), 1.0, color.RGBA{130, 130, 140, 255}, false)
		}
		return
	}

	if g.stopped {
		return
	}
	left, right := 40.0, float64(vfxlive.Width-40)
	g.marker += g.speed
	if g.marker < left || g.marker > right {
		g.speed = -g.speed
		g.marker = math.Max(left, math.Min(right, g.marker))
	}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	if !ate {
		g.updatePlay()
	}
	g.fx.Update()
	st := "RUN"
	if g.stopped {
		st = "STOP"
	}
	g.shell.SetToken("st", st)
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	g.shell.Hint = "game state + fx lifetime stay separate"
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{8, 18, 36, 255})
	_, sy, _, sh := g.shell.Stage()
	cx := float64(vfxlive.Width) / 2
	barY := sy + sh*0.45
	vector.DrawFilledRect(s, 35, float32(barY-30), float32(vfxlive.Width-70), 60, color.RGBA{16, 44, 69, 255}, false)
	vector.DrawFilledRect(s, float32(cx-50), float32(barY-30), 100, 60, color.RGBA{255, 120, 80, 255}, false)
	vector.DrawFilledRect(s, float32(cx-24), float32(barY-30), 48, 60, color.RGBA{255, 205, 69, 255}, false)
	vector.DrawFilledRect(s, float32(cx-10), float32(barY-30), 20, 60, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledRect(s, float32(g.marker-3), float32(barY-42), 6, 84, color.White, false)
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %d  ROUND %d  %s", g.score, g.round+1, g.label), 12, int(sy)+8)
	tip := "TAP STAGE TO STOP"
	if g.stopped {
		tip = "TAP FOR NEXT"
	}
	ebitenutil.DebugPrintAt(s, tip, 12, int(sy+sh)-20)
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Meter — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
