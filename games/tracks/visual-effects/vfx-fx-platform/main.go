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
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
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
	coins         []rect
	got, goal     int
	over, cleared bool
	sy, sh        float64
	landPulse     vfxmotion.Tween
	pose          vfxmotion.Locomotion
	runFrame      int
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
		g.fx.Dust(g.player.x+g.player.w/2, g.player.y+g.player.h, -g.vx, int(14*fx), color.RGBA{200, 180, 140, 255})
	}
	g.vy += 0.55
	g.player.x += g.vx
	g.player.y += g.vy
	g.player.x = math.Max(0, math.Min(480-g.player.w, g.player.x))
	wasGround := g.onGround
	g.onGround = false
	for _, p := range g.plats {
		if overlap(g.player, p) && g.vy >= 0 && g.player.y+g.player.h-g.vy <= p.y+4 {
			g.player.y = p.y - g.player.h
			g.vy = 0
			g.onGround = true
		}
	}
	edges := vfxmotion.DetectGroundEdges(wasGround, g.onGround)
	if edges.Landed {
		fx := g.shell.Get("fx")
		cx := g.player.x + g.player.w/2
		cy := g.player.y + g.player.h
		g.landPulse = vfxmotion.NewTween(10)
		g.fx.Shockwave(cx, cy, 0.38*fx, color.RGBA{235, 220, 190, 255}, color.RGBA{160, 150, 125, 255})
		g.fx.Dust(cx, cy, g.vx, int(14*fx), color.RGBA{180, 170, 140, 255})
	}
	if !g.landPulse.Done() {
		g.landPulse.Advance()
	}
	landedFrames := 0
	if !g.landPulse.Done() {
		landedFrames = g.landPulse.Frames - g.landPulse.Frame
	}
	g.pose = vfxmotion.PoseForPlatform(g.vx, g.vy, g.onGround, landedFrames)
	g.runFrame++
	keep := g.coins[:0]
	for _, c := range g.coins {
		if overlap(g.player, c) {
			g.got++
			fx := g.shell.Get("fx")
			g.fx.Burst(c.x+6, c.y+6, int(12*fx), 2*fx, color.RGBA{255, 220, 80, 255}, true)
			g.fx.Ring(c.x+6, c.y+6, 0.45*fx, color.RGBA{255, 235, 130, 255})
			continue
		}
		keep = append(keep, c)
	}
	g.coins = keep
	if g.got >= g.goal {
		g.cleared = true
		g.fx.FlashScreen(0.45, 255, 220, 80)
		g.fx.Confetti(g.player.x+g.player.w/2, g.player.y, int(36*g.shell.Get("fx")))
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
	pose := hero.Pose{}
	switch g.pose {
	case vfxmotion.LocomotionRun:
		bob := math.Sin(float64(g.runFrame)*0.38) * 0.05
		pose.ScaleX, pose.ScaleY = 1-bob, 1+bob
	case vfxmotion.LocomotionRise:
		pose.ScaleX, pose.ScaleY = 0.9, 1.13
	case vfxmotion.LocomotionFall:
		pose.ScaleX, pose.ScaleY = 1.08, 0.94
	case vfxmotion.LocomotionLand:
		squash := math.Sin(g.landPulse.Progress()*math.Pi) * 0.18
		pose.ScaleX, pose.ScaleY = 1+squash, 1-squash
	}
	hero.DrawBottomCenteredPose(screen, g.player.x+g.player.w/2, g.player.y+g.player.h, g.player.h+4, pose)
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
