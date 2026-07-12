// vfx-fx-shooter — ADVANCED: ship/enemies/bullets lists vs fx.Parts.
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
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type ball struct{ x, y, vx, vy, r float64 }

type Game struct {
	shell        *vfxlive.Shell
	fx           vfxfx.System
	shipX, shipY float64
	bullets      []ball
	enemies      []ball
	cool, score  int
	lives        int
	over         bool
	rng          *rand.Rand
	sy, sh       float64
}

func newGame() *Game {
	g := &Game{
		rng:   rand.New(rand.NewSource(88)),
		lives: 3,
		shell: vfxlive.New(
			"FX Shooter",
			[]string{
				"type Game struct {",
				"  bullets, enemies []ball",
				"  fx vfxfx.System // separate",
				"}",
				"shot→muzzle; kill→Burst+Flash",
				"bullets={b} enemies={e} fx={n}",
			},
			&vfxlive.Param{Key: "fx", Label: "fx", Value: 1, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "rate", Label: "fire", Value: 10, Min: 4, Max: 20, Step: 1, Format: "%.0f"},
		),
	}
	_, sy, _, sh := g.shell.Stage()
	g.sy, g.sh = sy, sh
	g.shipX, g.shipY = 240, sy+sh-48
	return g
}

func (g *Game) stage() {
	_, g.sy, _, g.sh = g.shell.Stage()
}

func (g *Game) spawnEnemy() {
	g.enemies = append(g.enemies, ball{
		x: 40 + g.rng.Float64()*400, y: g.sy + 20,
		vy: 1.2 + g.rng.Float64()*1.4, r: 14,
	})
}

func (g *Game) updatePlay() {
	g.stage()
	// steer
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.shipX -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.shipX += 4
	}
	if x, y, ok := vfxui.Held(); ok && y >= g.sy && y <= g.sy+g.sh {
		g.shipX += (x - g.shipX) * 0.2
	}
	g.shipX = math.Max(24, math.Min(456, g.shipX))
	g.shipY = g.sy + g.sh - 48

	g.cool++
	fire := int(g.shell.Get("rate"))
	wantShot := g.cool >= fire && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || len(ebiten.AppendTouchIDs(nil)) > 0)
	if wantShot {
		g.cool = 0
		g.bullets = append(g.bullets, ball{x: g.shipX, y: g.shipY - 18, vy: -8, r: 4})
		fx := g.shell.Get("fx")
		g.fx.Burst(g.shipX, g.shipY-18, int(8*fx), 1.6*fx, color.RGBA{255, 220, 120, 255}, true)
	}

	if g.rng.Float64() < 0.025 {
		g.spawnEnemy()
	}

	aliveB := g.bullets[:0]
	for _, b := range g.bullets {
		b.y += b.vy
		if b.y > g.sy-10 && b.y < g.sy+g.sh+10 {
			aliveB = append(aliveB, b)
		}
	}
	g.bullets = aliveB

	aliveE := g.enemies[:0]
	for _, e := range g.enemies {
		e.y += e.vy
		hit := false
		keepB := g.bullets[:0]
		for _, b := range g.bullets {
			if math.Hypot(b.x-e.x, b.y-e.y) < e.r+b.r {
				hit = true
				continue
			}
			keepB = append(keepB, b)
		}
		g.bullets = keepB
		if hit {
			g.score++
			fx := g.shell.Get("fx")
			g.fx.Burst(e.x, e.y, int(22*fx), 3.2*fx, color.RGBA{255, 140, 60, 255}, true)
			g.fx.FlashScreen(0.55*fx, 255, 160, 80)
			continue
		}
		if math.Hypot(e.x-g.shipX, e.y-g.shipY) < e.r+16 {
			g.lives--
			fx := g.shell.Get("fx")
			g.fx.Burst(g.shipX, g.shipY, int(16*fx), 2.5*fx, color.RGBA{255, 80, 80, 255}, true)
			g.fx.FlashScreen(0.7, 255, 40, 40)
			if g.lives <= 0 {
				g.over = true
			}
			continue
		}
		if e.y < g.sy+g.sh+20 {
			aliveE = append(aliveE, e)
		}
	}
	g.enemies = aliveE
}

func (g *Game) Update() error {
	ate := g.shell.Update()
	g.shell.SetToken("b", fmt.Sprintf("%d", len(g.bullets)))
	g.shell.SetToken("e", fmt.Sprintf("%d", len(g.enemies)))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	if g.over {
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
	screen.Fill(color.RGBA{6, 8, 18, 255})
	g.stage()
	g.shell.FillStage(screen, color.RGBA{8, 12, 28, 255})
	for _, e := range g.enemies {
		vector.DrawFilledCircle(screen, float32(e.x), float32(e.y), float32(e.r), color.RGBA{220, 70, 90, 255}, false)
	}
	for _, b := range g.bullets {
		vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), float32(b.r), color.RGBA{255, 240, 140, 255}, false)
	}
	hero.DrawCentered(screen, g.shipX, g.shipY, 36)
	g.fx.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("score %d  lives %d", g.score, g.lives), 12, int(g.sy)+6)
	if g.over {
		ebitenutil.DebugPrintAt(screen, "GAME OVER — tap to retry", 150, int(g.sy+g.sh/2))
	}
	g.shell.Hint = "drag/arrows to move · hold to fire · lists ≠ fx.Parts"
	g.shell.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Shooter — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
