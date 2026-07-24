// vfx-fx-bullethell — ADVANCED climax: bullets[] vs fx clear bomb.
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

type bullet struct{ x, y, vx, vy float64 }

type Game struct {
	shell          *vfxlive.Shell
	fx             vfxfx.System
	px, py         float64
	bullets        []bullet
	angle          float64
	frame, lives   int
	bombs, inv     int
	score          int
	over           bool
	rng            *rand.Rand
	sy, sh         float64
	bombBtn        vfxui.Button
	bombScript     vfxmotion.EffectScript
	scriptActive   bool
	clearedBullets []bullet
	freezeFrames   int
}

func newGame() *Game {
	g := &Game{
		rng:   rand.New(rand.NewSource(99)),
		lives: 3,
		bombs: 3,
		shell: vfxlive.New(
			"FX BulletHell",
			[]string{
				"type Game struct {",
				"  bullets []bullet // play",
				"  fx System        // clear VFX",
				"}",
				"bomb: Ring+Flash+Burst",
				"then delete nearby bullets",
				"bullets={b} fx.parts={n}",
			},
			&vfxlive.Param{Key: "fx", Label: "fx", Value: 1, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "dens", Label: "emit", Value: 1, Min: 0.4, Max: 2, Format: "%.1f"},
		),
	}
	g.stage()
	g.px, g.py = 240, g.sy+g.sh-70
	g.bombBtn = vfxui.Button{X: 360, Y: g.sy + 8, W: 100, H: 28, Label: "BOMB", Fill: color.RGBA{80, 40, 60, 235}}
	return g
}

func (g *Game) stage() {
	_, g.sy, _, g.sh = g.shell.Stage()
	g.bombBtn.Y = g.sy + 8
}

func (g *Game) emit() {
	ex, ey := 240.0, g.sy+50
	n := int(6 * g.shell.Get("dens"))
	for i := 0; i < n; i++ {
		a := g.angle + float64(i)*math.Pi*2/float64(n)
		sp := 1.6 + g.rng.Float64()*0.6
		g.bullets = append(g.bullets, bullet{
			x: ex, y: ey,
			vx: math.Cos(a) * sp, vy: math.Sin(a) * sp,
		})
	}
	g.angle += 0.28
}

func (g *Game) doBomb() {
	if g.bombs <= 0 {
		return
	}
	g.bombs--
	fx := g.shell.Get("fx")
	r := 110 * fx
	keep := g.bullets[:0]
	g.clearedBullets = g.clearedBullets[:0]
	for _, b := range g.bullets {
		if math.Hypot(b.x-g.px, b.y-g.py) > r {
			keep = append(keep, b)
		} else {
			g.clearedBullets = append(g.clearedBullets, b)
			g.score++
		}
	}
	g.bullets = keep
	g.inv = 30
	g.bombScript = vfxmotion.BombScript(fx)
	g.scriptActive = true
}

func (g *Game) updateBombScript() {
	if !g.scriptActive {
		return
	}
	fx := g.shell.Get("fx")
	for _, cue := range g.bombScript.Current() {
		switch cue.Kind {
		case vfxmotion.CueFreeze:
			g.freezeFrames = int(cue.Strength)
		case vfxmotion.CueFlash:
			g.fx.FlashScreen(cue.Strength, 255, 220, 160)
		case vfxmotion.CueShockwave:
			g.fx.Shockwave(g.px, g.py, cue.Strength, color.White, color.RGBA{255, 170, 70, 255})
			g.fx.Burst(g.px, g.py, int(32*fx), 4.5*fx, color.RGBA{255, 180, 80, 255}, true)
		case vfxmotion.CueDissolve:
			limit := min(len(g.clearedBullets), 48)
			for i := 0; i < limit; i++ {
				b := g.clearedBullets[i]
				g.fx.Burst(b.x, b.y, 2, 1.2*fx, color.RGBA{255, 180, 220, 255}, true)
			}
		case vfxmotion.CueConfetti:
			g.fx.Confetti(g.px, g.py, int(28*fx))
		}
	}
	g.bombScript.Advance()
	if g.bombScript.Done() {
		g.scriptActive = false
		g.clearedBullets = g.clearedBullets[:0]
	}
}

func (g *Game) updatePlay(ate bool) {
	g.stage()
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		dy++
	}
	if x, y, ok := vfxui.Held(); ok && y >= g.sy && y <= g.sy+g.sh && !g.bombBtn.Contains(x, y) {
		dx = x - g.px
		dy = y - g.py
	}
	if dx != 0 || dy != 0 {
		l := math.Hypot(dx, dy)
		g.px += dx / l * 3.6
		g.py += dy / l * 3.6
	}
	g.px = math.Max(24, math.Min(456, g.px))
	g.py = math.Max(g.sy+40, math.Min(g.sy+g.sh-24, g.py))

	wantBomb := inpututil.IsKeyJustPressed(ebiten.KeySpace) || (!ate && g.bombBtn.Tapped())
	if wantBomb {
		g.doBomb()
	}

	g.frame++
	period := int(18 / math.Max(0.4, g.shell.Get("dens")))
	if g.frame%period == 0 {
		g.emit()
	}

	alive := g.bullets[:0]
	for _, b := range g.bullets {
		b.x += b.vx
		b.y += b.vy
		if b.x < -20 || b.x > 500 || b.y < g.sy-20 || b.y > g.sy+g.sh+20 {
			continue
		}
		if g.inv == 0 && math.Hypot(b.x-g.px, b.y-g.py) < 10 {
			g.lives--
			g.inv = 50
			fx := g.shell.Get("fx")
			g.fx.Burst(g.px, g.py, int(14*fx), 2.5*fx, color.RGBA{255, 80, 100, 255}, true)
			g.fx.FlashScreen(0.6, 255, 50, 70)
			if g.lives <= 0 {
				g.over = true
			}
			continue
		}
		alive = append(alive, b)
	}
	g.bullets = alive
	if g.inv > 0 {
		g.inv--
	}
	g.score++
}

func (g *Game) Update() error {
	ate := g.shell.Update()
	g.shell.SetToken("b", fmt.Sprintf("%d", len(g.bullets)))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	if g.over {
		if vfxui.AnyPressStart() && !ate {
			*g = *newGame()
		}
		g.fx.Update()
		return nil
	}
	if g.freezeFrames > 0 {
		g.freezeFrames--
	} else {
		g.updatePlay(ate)
	}
	g.updateBombScript()
	g.fx.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 6, 16, 255})
	g.stage()
	g.shell.FillStage(screen, color.RGBA{14, 10, 28, 255})
	vector.DrawFilledCircle(screen, 240, float32(g.sy+50), 22, color.RGBA{180, 60, 90, 255}, false)
	for _, b := range g.bullets {
		vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), 4, color.RGBA{255, 160, 200, 255}, false)
	}
	if g.scriptActive {
		t := float64(g.bombScript.Frame) / 13
		alpha := uint8(180 * math.Max(0, 1-t))
		for i, b := range g.clearedBullets {
			if i >= 80 {
				break
			}
			x := b.x + b.vx*t*18
			y := b.y + b.vy*t*18
			vector.StrokeCircle(screen, float32(x), float32(y), float32(4+t*7), 2,
				color.RGBA{255, 180, 230, alpha}, false)
		}
	}
	if g.inv%4 < 2 {
		hero.DrawCentered(screen, g.px, g.py, 32)
	}
	g.fx.Draw(screen)
	g.bombBtn.Draw(screen, g.bombs > 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("life %d  bombs %d  score %d", g.lives, g.bombs, g.score), 12, int(g.sy)+10)
	if g.over {
		ebitenutil.DebugPrintAt(screen, "GAME OVER — tap to retry", 150, int(g.sy+g.sh/2))
	}
	g.shell.Hint = "dodge · Space/BOMB clears bullets[] then paints fx"
	g.shell.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX BulletHell — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
