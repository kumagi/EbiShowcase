// vfx-fx-catch — ADV 03: catch stars; stars[] and fx.Parts[] stay separate.
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type star struct {
	id          int
	x, y, speed float64
}

type game struct {
	shell       *vfxlive.Shell
	fx          vfxfx.System
	rng         *rand.Rand
	basketX     float64
	stars       []star // gameplay entities — NOT mixed with fx.Parts
	frame       int
	score       int
	lives       int
	over        bool
	nextID      int
	proxies     []vfxmotion.Proxy
	basketPulse vfxmotion.Tween
}

func newGame() *game {
	g := &game{
		rng:     rand.New(rand.NewSource(19)),
		basketX: vfxlive.Width / 2,
		lives:   3,
		shell: vfxlive.New(
			"FX catch",
			[]string{
				"func (g *Game) Update() {",
				"  g.updatePlay() // stars={s}  (game list)",
				"  g.fx.Update()  // parts={n}  (FX list)",
				"}",
				"// mixing stars+sparks in one slice = messy",
			},
			&vfxlive.Param{Key: "fx", Label: "fx intensity", Value: 1.0, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "rate", Label: "spawn rate", Value: 40, Min: 20, Max: 80, Step: 5, Format: "%.0f"},
		),
	}
	return g
}

func (g *game) reset() {
	g.stars = nil
	g.fx = vfxfx.System{}
	g.score, g.lives, g.frame = 0, 3, 0
	g.over = false
	g.basketX = vfxlive.Width / 2
	g.nextID = 0
	g.proxies = nil
	g.basketPulse = vfxmotion.Tween{}
}

func (g *game) moveBasket() {
	if x, _, ok := vfxui.Held(); ok {
		g.basketX = x
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.basketX -= 6
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.basketX += 6
	}
	g.basketX = math.Max(40, math.Min(float64(vfxlive.Width-40), g.basketX))
}

func (g *game) updatePlay() {
	_, sy, _, sh := g.shell.Stage()
	floor := sy + sh - 8
	basketY := floor - 28

	if g.over {
		if _, y, ok := vfxui.JustPressed(); ok && y >= sy && y <= sy+sh {
			g.reset()
		}
		return
	}
	g.moveBasket()
	g.frame++
	rate := int(g.shell.Get("rate"))
	if g.frame%rate == 0 {
		g.nextID++
		g.stars = append(g.stars, star{
			id: g.nextID,
			x:  30 + g.rng.Float64()*float64(vfxlive.Width-60),
			y:  sy + 10, speed: 2.4 + g.rng.Float64()*2,
		})
	}
	k := g.shell.Get("fx")
	next := g.stars[:0]
	for _, s := range g.stars {
		s.y += s.speed
		caught := s.y > basketY-10 && s.y < basketY+24 && math.Abs(s.x-g.basketX) < 52
		if caught {
			g.score++
			g.proxies = append(g.proxies, vfxmotion.NewProxy(s.id, s.x, s.y, 28, sy+18, 24))
			g.basketPulse = vfxmotion.NewTween(10)
			g.fx.Burst(s.x, s.y, int(16*k), 2.6*k, color.RGBA{255, 230, 120, 255}, true)
			g.fx.Shockwave(s.x, s.y, 0.55*k, color.White, color.RGBA{255, 220, 100, 255})
			continue
		}
		if s.y > floor {
			g.lives--
			g.fx.Burst(s.x, floor-4, int(8*k), 1.1, color.RGBA{150, 140, 130, 255}, false)
			if g.lives <= 0 {
				g.over = true
			}
			continue
		}
		next = append(next, s)
	}
	g.stars = next
	visible := g.proxies[:0]
	for i := range g.proxies {
		g.proxies[i].Advance()
		if !g.proxies[i].Done() {
			visible = append(visible, g.proxies[i])
		}
	}
	g.proxies = visible
	if !g.basketPulse.Done() {
		g.basketPulse.Advance()
	}
}

func (g *game) Update() error {
	ate := g.shell.Update()
	if !ate {
		g.updatePlay()
	}
	g.fx.Update()
	g.shell.SetToken("s", fmt.Sprintf("%d", len(g.stars)))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	g.shell.Hint = "stars[] gameplay · fx.Parts[] visuals — keep apart"
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{10, 22, 42, 255})
	_, sy, _, sh := g.shell.Stage()
	floor := sy + sh - 8
	basketY := floor - 28
	for _, st := range g.stars {
		vector.DrawFilledCircle(s, float32(st.x), float32(st.y), 8, color.RGBA{255, 220, 90, 255}, false)
	}
	for _, proxy := range g.proxies {
		x, y := proxy.Position()
		radius := 4 + 5*(1-proxy.Tween.Progress())
		vector.DrawFilledCircle(s, float32(x), float32(y), float32(radius), color.RGBA{255, 240, 150, 220}, false)
	}
	squash := 0.0
	if !g.basketPulse.Done() {
		squash = math.Sin(g.basketPulse.Progress()*math.Pi) * 0.16
	}
	hero.DrawBottomCenteredPose(s, g.basketX, basketY+8, 56, hero.Pose{ScaleX: 1 + squash, ScaleY: 1 - squash})
	vector.DrawFilledRect(s, float32(g.basketX-48*(1+squash)), float32(basketY), float32(96*(1+squash)), 14, color.RGBA{255, 105, 79, 255}, false)
	g.fx.Draw(s)
	msg := fmt.Sprintf("SCORE %d  LIVES %d", g.score, g.lives)
	if g.over {
		msg = fmt.Sprintf("OVER  SCORE %d — tap to retry", g.score)
	}
	ebitenutil.DebugPrintAt(s, msg, 12, int(sy)+8)
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Catch — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
