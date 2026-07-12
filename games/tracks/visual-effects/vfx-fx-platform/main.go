// vfx-fx-platform — ADVANCED: gravity player vs jump/land fx.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type rect struct{ x, y, w, h float64 }

type Game struct {
	shell         *vfxlive.Shell
	fx            vfxfx.System
	player        rect
	vx, vy        float64
	plats         []rect
	onGround      bool
	wasGround     bool
	coins         []rect
	got, goal     int
	over, cleared bool
	sy, sh        float64
}

func newGame() *Game {
	g := &Game{
		player: rect{60, 200, 26, 34},
		goal:   3,
		shell: vfxlive.New(
			"FX Platform",
			[]string{
				"type Game struct {",
				"  px,py,vx,vy float64 // physics",
				"  fx System           // dust",
				"}",
				"jump→Burst down; land→Ring",
				"coins={c} fx.parts={n}",
			},
			&vfxlive.Param{Key: "fx", Label: "fx", Value: 1, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "jump", Label: "jump", Value: 11, Min: 8, Max: 15, Format: "%.1f"},
		),
	}
	g.stage()
	floor := g.sy + g.sh - 28
	g.plats = []rect{
		{0, floor, 480, 28},
		{80, floor - 90, 110, 16},
		{220, floor - 160, 100, 16},
		{340, floor - 110, 100, 16},
		{160, floor - 240, 90, 16},
	}
	g.coins = []rect{
		{120, floor - 120, 12, 12},
		{255, floor - 190, 12, 12},
		{380, floor - 140, 12, 12},
	}
	g.player.y = floor - g.player.h - 2
	g.onGround = true
	return g
}

func (g *Game) stage() { _, g.sy, _, g.sh = g.shell.Stage() }

func overlap(a, b rect) bool {
	return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y
}

func (g *Game) updatePlay() {
	g.stage()
	left := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	if x, y, ok := vfxui.JustPressed(); ok && y >= g.sy && y <= g.sy+g.sh {
		if y > g.sy+g.sh*0.65 {
			jump = true
		} else if x < g.player.x {
			left = true
		} else if x > g.player.x+g.player.w {
			right = true
		}
	}
	if x, y, ok := vfxui.Held(); ok && y >= g.sy && y <= g.sy+g.sh*0.65 {
		if x < g.player.x {
			left = true
		} else if x > g.player.x+g.player.w {
			right = true
		}
	}
	ax := 0.0
	if left {
		ax = -0.55
	}
	if right {
		ax = 0.55
	}
	g.vx = math.Max(-4.2, math.Min(4.2, g.vx*0.86+ax*3))
	if jump && g.onGround {
		g.vy = -g.shell.Get("jump")
		g.onGround = false
		fx := g.shell.Get("fx")
		g.fx.Burst(g.player.x+g.player.w/2, g.player.y+g.player.h, int(14*fx), 2.2*fx, color.RGBA{200, 180, 140, 255}, false)
	}
	g.vy += 0.55
	g.player.x += g.vx
	g.player.y += g.vy
	g.player.x = math.Max(0, math.Min(480-g.player.w, g.player.x))
	g.wasGround = g.onGround
	g.onGround = false
	for _, p := range g.plats {
		if overlap(g.player, p) && g.vy >= 0 && g.player.y+g.player.h-g.vy <= p.y+4 {
			g.player.y = p.y - g.player.h
			g.vy = 0
			g.onGround = true
		}
	}
	if g.onGround && !g.wasGround {
		fx := g.shell.Get("fx")
		cx := g.player.x + g.player.w/2
		cy := g.player.y + g.player.h
		g.fx.Ring(cx, cy, 0.55*fx, color.RGBA{220, 200, 160, 255})
		g.fx.Burst(cx, cy, int(8*fx), 1.4*fx, color.RGBA{180, 170, 140, 255}, false)
	}
	keep := g.coins[:0]
	for _, c := range g.coins {
		if overlap(g.player, c) {
			g.got++
			fx := g.shell.Get("fx")
			g.fx.Burst(c.x+6, c.y+6, int(12*fx), 2*fx, color.RGBA{255, 220, 80, 255}, true)
			continue
		}
		keep = append(keep, c)
	}
	g.coins = keep
	if g.got >= g.goal {
		g.cleared = true
		g.fx.FlashScreen(0.45, 255, 220, 80)
	}
	if g.player.y > g.sy+g.sh+40 {
		g.over = true
	}
}

func (g *Game) Update() error {
	ate := g.shell.Update()
	g.shell.SetToken("c", fmt.Sprintf("%d", g.got))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	if g.over || g.cleared {
		if vfxui.AnyPressStart() && !ate {
			*g = *newGame()
		}
		g.fx.Update()
		return nil
	}
	_ = ate
	g.updatePlay()
	g.fx.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 28, 48, 255})
	g.stage()
	g.shell.FillStage(screen, color.RGBA{40, 120, 200, 255})
	for _, p := range g.plats {
		vector.DrawFilledRect(screen, float32(p.x), float32(p.y), float32(p.w), float32(p.h), color.RGBA{60, 90, 70, 255}, false)
	}
	for _, c := range g.coins {
		vector.DrawFilledCircle(screen, float32(c.x+6), float32(c.y+6), 7, color.RGBA{255, 210, 60, 255}, false)
	}
	hero.DrawBottomCentered(screen, g.player.x+g.player.w/2, g.player.y+g.player.h, g.player.h+4)
	g.fx.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("coins %d/%d", g.got, g.goal), 12, int(g.sy)+6)
	if g.cleared {
		ebitenutil.DebugPrintAt(screen, "CLEAR — tap to retry", 160, int(g.sy+g.sh/2))
	}
	if g.over {
		ebitenutil.DebugPrintAt(screen, "FALL — tap to retry", 165, int(g.sy+g.sh/2))
	}
	g.shell.Hint = "A/D move · Space/tap low = jump · physics ≠ fx"
	g.shell.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Platform — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
