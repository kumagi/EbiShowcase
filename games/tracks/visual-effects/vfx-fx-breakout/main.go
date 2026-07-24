// vfx-fx-breakout — ADV 06: bricks[] vs fx.Parts[]; FlameBurst on break.
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
	brickW = 48.0
	brickH = 18.0
	ballR  = 9.0
	padHW  = 50.0
)

type brick struct {
	x, y  float64
	alive bool
	row   int
}

type game struct {
	shell   *vfxlive.Shell
	fx      vfxfx.System
	paddleX float64
	bx, by  float64
	vx, vy  float64
	bricks  []brick // world tiles — never mixed into fx.Parts
	score   int
	lives   int
	ready   bool
	ghosts  []vfxmotion.Tombstone
	hitStop int
}

func newGame() *game {
	g := &game{
		paddleX: vfxlive.Width / 2,
		lives:   3,
		shell: vfxlive.New(
			"FX breakout",
			[]string{
				"func (g *Game) Update() {",
				"  g.updatePlay() // bricks={b}  (tiles)",
				"  g.fx.Update()  // parts={n}   (FX)",
				"}",
				"// brick break → FlameBurst / Burst+Ring",
			},
			&vfxlive.Param{Key: "fx", Label: "fx intensity", Value: 1.0, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "spd", Label: "ball speed", Value: 1.0, Min: 0.6, Max: 1.5, Format: "%.2f"},
		),
	}
	g.buildBricks()
	return g
}

func (g *game) buildBricks() {
	g.bricks = g.bricks[:0]
	_, sy, _, _ := g.shell.Stage()
	for row := 0; row < 5; row++ {
		for col := 0; col < 8; col++ {
			g.bricks = append(g.bricks, brick{
				x: 16 + float64(col)*58, y: sy + 36 + float64(row)*24,
				alive: true, row: row,
			})
		}
	}
}

func (g *game) serve() {
	_, sy, _, sh := g.shell.Stage()
	g.bx = g.paddleX
	g.by = sy + sh - 48
	g.vx, g.vy = 3.2, -4.2
	g.ready = true
}

func (g *game) aliveCount() int {
	n := 0
	for _, b := range g.bricks {
		if b.alive {
			n++
		}
	}
	return n
}

func (g *game) updatePlay() {
	if !g.ready {
		g.buildBricks()
		g.serve()
	}
	_, sy, _, sh := g.shell.Stage()
	bot := sy + sh - 22
	k := g.shell.Get("fx")

	if x, _, ok := vfxui.Held(); ok {
		g.paddleX = x
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.paddleX -= 7
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.paddleX += 7
	}
	g.paddleX = math.Max(padHW, math.Min(float64(vfxlive.Width)-padHW, g.paddleX))

	sp := g.shell.Get("spd")
	g.bx += g.vx * sp
	g.by += g.vy * sp
	if g.bx < ballR || g.bx > float64(vfxlive.Width)-ballR {
		g.vx = -g.vx
		g.bx = math.Max(ballR, math.Min(float64(vfxlive.Width)-ballR, g.bx))
	}
	if g.by < sy+8 {
		g.vy = math.Abs(g.vy)
	}
	if g.vy > 0 && g.by > bot-14 && g.by < bot+8 && math.Abs(g.bx-g.paddleX) < padHW+6 {
		g.by = bot - 14
		g.vy = -math.Abs(g.vy)
		g.vx += (g.bx - g.paddleX) * 0.06
		g.fx.Burst(g.bx, g.by, int(8*k), 1.8*k, color.RGBA{120, 240, 220, 255}, true)
	}

	for i := range g.bricks {
		b := &g.bricks[i]
		if !b.alive {
			continue
		}
		hit := g.bx+ballR > b.x && g.bx-ballR < b.x+brickW &&
			g.by+ballR > b.y && g.by-ballR < b.y+brickH
		if !hit {
			continue
		}
		g.ghosts = append(g.ghosts, vfxmotion.NewTombstone(vfxmotion.Snapshot{
			ID: i, X: b.x, Y: b.y, W: brickW, H: brickH, Variant: b.row,
		}, 24))
		b.alive = false
		g.score += 10
		g.hitStop = 2
		g.vy = -g.vy
		cx, cy := b.x+brickW/2, b.y+brickH/2
		if b.row%2 == 0 {
			g.fx.FlameBurst(cx, cy, int(10*k))
		} else {
			tint := color.RGBA{uint8(80 + b.row*30), 180, 255, 255}
			g.fx.Burst(cx, cy, int(14*k), 2.4*k, tint, true)
			g.fx.Ring(cx, cy, 0.75*k, tint)
		}
		break
	}

	if g.by > sy+sh+20 {
		g.lives--
		if g.lives <= 0 || g.aliveCount() == 0 {
			g.score, g.lives = 0, 3
			g.buildBricks()
			g.fx = vfxfx.System{}
		}
		g.serve()
	}
	if g.aliveCount() == 0 {
		g.buildBricks()
		g.fx.FlashScreen(0.5*k, 255, 160, 60)
		g.fx.Confetti(vfxlive.Width/2, sy+sh*0.55, int(48*k))
		g.serve()
	}
}

func (g *game) updateGhosts() {
	alive := g.ghosts[:0]
	for i := range g.ghosts {
		g.ghosts[i].Advance()
		if !g.ghosts[i].Done() {
			alive = append(alive, g.ghosts[i])
		}
	}
	g.ghosts = alive
}

func (g *game) Update() error {
	ate := g.shell.Update()
	if !ate {
		if g.hitStop > 0 {
			g.hitStop--
		} else {
			g.updatePlay()
		}
	} else if !g.ready {
		g.buildBricks()
		g.serve()
	}
	g.updateGhosts()
	g.fx.Update()
	g.shell.SetToken("b", fmt.Sprintf("%d", g.aliveCount()))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	g.shell.Hint = "bricks[] world · fx.Parts[] sparks — keep separate"
	return nil
}

var brickTint = []color.RGBA{
	{255, 105, 79, 255}, {255, 180, 70, 255}, {255, 220, 90, 255},
	{45, 226, 194, 255}, {100, 180, 255, 255},
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{6, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{8, 18, 36, 255})
	_, sy, _, sh := g.shell.Stage()
	bot := sy + sh - 22
	for _, b := range g.bricks {
		if !b.alive {
			continue
		}
		vector.DrawFilledRect(s, float32(b.x), float32(b.y), brickW, brickH, brickTint[b.row%len(brickTint)], false)
	}
	for _, ghost := range g.ghosts {
		snapshot := ghost.Snapshot
		t := ghost.Progress()
		alpha := uint8(220 * (1 - t))
		base := brickTint[snapshot.Variant%len(brickTint)]
		c := color.RGBA{base.R, base.G, base.B, alpha}
		for shard := 0; shard < 4; shard++ {
			sideX := -1.0
			if shard%2 == 1 {
				sideX = 1
			}
			sideY := -1.0
			if shard >= 2 {
				sideY = 1
			}
			x := snapshot.X + snapshot.W/2 + sideX*t*(12+float64(shard)*3)
			y := snapshot.Y + snapshot.H/2 + sideY*t*(8+float64(shard)*2) + t*t*18
			vector.DrawFilledRect(s, float32(x-snapshot.W/5), float32(y-snapshot.H/5),
				float32(snapshot.W/2.5), float32(snapshot.H/2.5), c, false)
		}
	}
	vector.DrawFilledRect(s, float32(g.paddleX-padHW), float32(bot), float32(padHW*2), 12, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledCircle(s, float32(g.bx), float32(g.by), ballR, color.White, false)
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %d  LIVES %d  bricks=%d", g.score, g.lives, g.aliveCount()), 12, int(sy)+8)
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Breakout — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
