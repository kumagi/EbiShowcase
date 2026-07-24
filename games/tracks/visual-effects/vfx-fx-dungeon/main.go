// vfx-fx-dungeon — ADVANCED: HP combat play vs hit flash/burst fx.
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

type slime struct {
	x, y, hp float64
	alive    bool
	reaction vfxmotion.Reaction
}

type Game struct {
	shell          *vfxlive.Shell
	fx             vfxfx.System
	px, py         float64
	fxDir, fyDir   float64
	slash          int
	php            float64
	slimes         []slime
	inv            int
	over, cleared  bool
	sy, sh         float64
	slashHit       bool
	playerReaction vfxmotion.Reaction
}

func newGame() *Game {
	g := &Game{
		px: 240, py: 360, fyDir: -1, php: 5,
		shell: vfxlive.New(
			"FX Dungeon",
			[]string{
				"type Game struct {",
				"  php float64; slimes []slime",
				"  fx System // hit VFX only",
				"}",
				"slash hit→Flash+Burst",
				"HP={hp} fx.parts={n}",
			},
			&vfxlive.Param{Key: "fx", Label: "fx", Value: 1, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "spd", Label: "speed", Value: 3.2, Min: 1.5, Max: 5, Format: "%.1f"},
		),
	}
	g.stage()
	cy := g.sy + g.sh/2
	g.py = cy
	g.slimes = []slime{
		{x: 120, y: cy - 80, hp: 3, alive: true},
		{x: 360, y: cy + 40, hp: 3, alive: true},
		{x: 240, y: cy - 140, hp: 4, alive: true},
	}
	return g
}

func (g *Game) stage() { _, g.sy, _, g.sh = g.shell.Stage() }

func (g *Game) updatePlay() {
	g.stage()
	spd := g.shell.Get("spd")
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
	if x, y, ok := vfxui.Held(); ok && y >= g.sy && y <= g.sy+g.sh {
		dx = x - g.px
		dy = y - g.py
	}
	if dx != 0 || dy != 0 {
		l := math.Hypot(dx, dy)
		g.px += dx / l * spd
		g.py += dy / l * spd
		g.fxDir, g.fyDir = dx/l, dy/l
	}
	g.px = math.Max(30, math.Min(450, g.px))
	g.py = math.Max(g.sy+30, math.Min(g.sy+g.sh-30, g.py))

	slash := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyZ)
	if x, y, ok := vfxui.JustPressed(); ok && y >= g.sy && y <= g.sy+g.sh {
		if math.Hypot(x-g.px, y-g.py) < 40 {
			slash = true
		}
	}
	if slash && g.slash <= 0 {
		g.slash = 12
		g.slashHit = false
	}
	if g.slash > 0 {
		g.slash--
		sx := g.px + g.fxDir*28
		sy := g.py + g.fyDir*28
		for i := range g.slimes {
			if g.slashHit {
				break
			}
			s := &g.slimes[i]
			if !s.alive {
				continue
			}
			if math.Hypot(s.x-sx, s.y-sy) < 36 {
				s.hp--
				s.reaction = vfxmotion.NewReaction(2, 5, 10)
				g.slashHit = true
				fx := g.shell.Get("fx")
				g.fx.Burst(s.x, s.y, int(16*fx), 2.8*fx, color.RGBA{120, 255, 140, 255}, true)
				g.fx.Shockwave(s.x, s.y, 0.45*fx, color.White, color.RGBA{100, 255, 130, 255})
				g.fx.FlashScreen(0.28*fx, 180, 255, 160)
				if s.hp <= 0 {
					s.alive = false
					g.fx.Confetti(s.x, s.y, int(14*fx))
				}
				break
			}
		}
	}
	if g.inv > 0 {
		g.inv--
	}
	alive := 0
	for i := range g.slimes {
		s := &g.slimes[i]
		if !s.alive {
			continue
		}
		alive++
		if math.Hypot(s.x-g.px, s.y-g.py) < 26 && g.inv == 0 {
			g.php--
			g.inv = 40
			g.playerReaction = vfxmotion.NewReaction(2, 6, 12)
			fx := g.shell.Get("fx")
			g.fx.Burst(g.px, g.py, int(10*fx), 2*fx, color.RGBA{255, 80, 80, 255}, true)
			g.fx.FlashScreen(0.5, 255, 40, 40)
			if g.php <= 0 {
				g.over = true
			}
		}
	}
	if alive == 0 {
		g.cleared = true
		g.fx.Confetti(vfxlive.Width/2, g.sy+g.sh/2, int(48*g.shell.Get("fx")))
	}
}

func (g *Game) reactionsInHitStop() bool {
	if g.playerReaction.Phase() == vfxmotion.ReactionHitStop {
		return true
	}
	for i := range g.slimes {
		if g.slimes[i].reaction.Phase() == vfxmotion.ReactionHitStop {
			return true
		}
	}
	return false
}

func (g *Game) advanceReactions() {
	g.playerReaction.Advance()
	for i := range g.slimes {
		g.slimes[i].reaction.Advance()
	}
}

func (g *Game) Update() error {
	ate := g.shell.Update()
	g.shell.SetToken("hp", fmt.Sprintf("%.0f", g.php))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	if g.over || g.cleared {
		if vfxui.AnyPressStart() && !ate {
			*g = *newGame()
		}
		g.fx.Update()
		return nil
	}
	_ = ate
	if !g.reactionsInHitStop() {
		g.updatePlay()
	}
	g.advanceReactions()
	g.fx.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 14, 22, 255})
	g.stage()
	g.shell.FillStage(screen, color.RGBA{28, 36, 52, 255})
	vector.StrokeRect(screen, 24, float32(g.sy+16), 432, float32(g.sh-32), 3, color.RGBA{60, 70, 90, 255}, false)
	for _, s := range g.slimes {
		if !s.alive {
			continue
		}
		x := s.x + s.reaction.Offset(7)
		tint := color.RGBA{80, 200, 100, 255}
		if s.reaction.Phase() == vfxmotion.ReactionFlash {
			tint = color.RGBA{235, 255, 235, 255}
		}
		vector.DrawFilledCircle(screen, float32(x), float32(s.y), 18, tint, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", s.hp), int(x)-4, int(s.y)-6)
	}
	if g.slash > 0 {
		base := math.Atan2(g.fyDir, g.fxDir)
		progress := 1 - float64(g.slash)/12
		for i := 0; i < 3; i++ {
			a := base - 0.8 + progress*1.6 + float64(i-1)*0.1
			x1, y1 := g.px+math.Cos(a)*20, g.py+math.Sin(a)*20
			x2, y2 := g.px+math.Cos(a)*55, g.py+math.Sin(a)*55
			vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2),
				float32(5-i), color.RGBA{255, 240, 180, uint8(220 - i*55)}, false)
		}
	}
	hero.DrawCenteredPose(screen, g.px+g.playerReaction.Offset(6), g.py, 40, hero.Pose{})
	g.fx.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %.0f", g.php), 12, int(g.sy)+6)
	if g.cleared {
		ebitenutil.DebugPrintAt(screen, "CLEAR — tap to retry", 160, int(g.sy+g.sh/2))
	}
	if g.over {
		ebitenutil.DebugPrintAt(screen, "DEFEAT — tap to retry", 155, int(g.sy+g.sh/2))
	}
	g.shell.Hint = "move: keys/drag · Space/tap self = slash · HP ≠ fx"
	g.shell.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Dungeon — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
