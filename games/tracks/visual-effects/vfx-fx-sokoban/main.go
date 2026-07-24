// vfx-fx-sokoban — ADVANCED: tile map play vs fx dust/sparkle.
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

type pt struct{ x, y int }

type moveAnimation struct {
	plan        vfxmotion.MovePlan
	tween       vfxmotion.Tween
	reachedGoal bool
}

type Game struct {
	shell      *vfxlive.Shell
	fx         vfxfx.System
	player     pt
	boxes      []pt
	goals      []pt
	walls      map[pt]bool
	moves      int
	cleared    bool
	ox, oy, ts float64
	w, h       int
	move       *moveAnimation
	frame      int
}

var level = []string{
	"#######",
	"#  .  #",
	"# $@$ #",
	"#  .  #",
	"#######",
}

func newGame() *Game {
	g := &Game{
		walls: map[pt]bool{},
		shell: vfxlive.New(
			"FX Sokoban",
			[]string{
				"type Game struct {",
				"  map tiles; boxes []pt",
				"  fx System // dust ≠ tiles",
				"}",
				"push→Burst; onGoal→Ring",
				"boxes={b} fx.parts={n}",
			},
			&vfxlive.Param{Key: "fx", Label: "fx", Value: 1, Min: 0.3, Max: 1.5, Format: "%.2f"},
			&vfxlive.Param{Key: "dust", Label: "dustN", Value: 10, Min: 4, Max: 24, Step: 1, Format: "%.0f"},
		),
	}
	for y, row := range level {
		g.h = y + 1
		g.w = len(row)
		for x, c := range row {
			p := pt{x, y}
			switch c {
			case '#':
				g.walls[p] = true
			case '@':
				g.player = p
			case '$':
				g.boxes = append(g.boxes, p)
			case '.':
				g.goals = append(g.goals, p)
			}
		}
	}
	g.layout()
	return g
}

func (g *Game) layout() {
	_, sy, sw, sh := g.shell.Stage()
	g.ts = mathMin(sw/float64(g.w), sh/float64(g.h)) * 0.92
	g.ox = (sw - float64(g.w)*g.ts) / 2
	g.oy = sy + (sh-float64(g.h)*g.ts)/2
}

func mathMin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (g *Game) center(p pt) (float64, float64) {
	return g.ox + (float64(p.x)+0.5)*g.ts, g.oy + (float64(p.y)+0.5)*g.ts
}

func (g *Game) boxAt(p pt) int {
	for i, b := range g.boxes {
		if b == p {
			return i
		}
	}
	return -1
}

func (g *Game) onGoal(p pt) bool {
	for _, g0 := range g.goals {
		if g0 == p {
			return true
		}
	}
	return false
}

func (g *Game) tryMove(dx, dy int) {
	if g.move != nil {
		return
	}
	walls := make(map[vfxmotion.Cell]bool, len(g.walls))
	for wall := range g.walls {
		walls[vfxmotion.Cell{X: wall.x, Y: wall.y}] = true
	}
	boxes := make([]vfxmotion.Cell, len(g.boxes))
	for i, box := range g.boxes {
		boxes[i] = vfxmotion.Cell{X: box.x, Y: box.y}
	}
	plan := vfxmotion.PlanSokobanMove(
		vfxmotion.Cell{X: g.player.x, Y: g.player.y},
		vfxmotion.Cell{X: dx, Y: dy},
		walls,
		boxes,
	)
	if !plan.Allowed {
		return
	}

	g.player = pt{plan.PlayerTo.X, plan.PlayerTo.Y}
	reachedGoal := false
	if plan.BoxIndex >= 0 {
		g.boxes[plan.BoxIndex] = pt{plan.BoxTo.X, plan.BoxTo.Y}
		reachedGoal = g.onGoal(g.boxes[plan.BoxIndex])
	}
	g.moves++
	g.move = &moveAnimation{
		plan: plan, tween: vfxmotion.NewTween(9), reachedGoal: reachedGoal,
	}
	g.cleared = true
	for _, b := range g.boxes {
		if !g.onGoal(b) {
			g.cleared = false
			break
		}
	}
}

func (g *Game) updateMoveAnimation() {
	if g.move == nil {
		return
	}
	g.move.tween.Advance()
	if !g.move.tween.Done() {
		return
	}
	fx := g.shell.Get("fx")
	plan := g.move.plan
	if plan.BoxIndex >= 0 {
		x, y := g.center(pt{plan.BoxTo.X, plan.BoxTo.Y})
		g.fx.Dust(x, y+g.ts*0.28, float64(plan.BoxTo.X-plan.BoxFrom.X), int(g.shell.Get("dust")*fx),
			color.RGBA{180, 160, 120, 255})
		if g.move.reachedGoal {
			g.fx.Shockwave(x, y, 0.65*fx, color.RGBA{210, 255, 230, 255}, color.RGBA{80, 230, 170, 255})
			g.fx.Burst(x, y, int(12*fx), 2.2*fx, color.RGBA{180, 255, 220, 255}, true)
		}
	}
	g.move = nil
	if g.cleared {
		g.fx.FlashScreen(0.5*fx, 120, 255, 200)
		g.fx.Confetti(vfxlive.Width/2, g.oy+float64(g.h)*g.ts/2, int(40*fx))
	}
}

func (g *Game) readInput(ate bool) {
	if ate || g.cleared || g.move != nil {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.tryMove(0, -1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.tryMove(0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.tryMove(-1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.tryMove(1, 0)
	}
	if x, y, ok := vfxui.JustPressed(); ok {
		px, py := g.center(g.player)
		dx, dy := x-px, y-py
		if abs(dx) < g.ts*0.3 && abs(dy) < g.ts*0.3 {
			return
		}
		if abs(dx) > abs(dy) {
			if dx < 0 {
				g.tryMove(-1, 0)
			} else {
				g.tryMove(1, 0)
			}
		} else {
			if dy < 0 {
				g.tryMove(0, -1)
			} else {
				g.tryMove(0, 1)
			}
		}
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func (g *Game) Update() error {
	ate := g.shell.Update()
	g.layout()
	g.frame++
	g.shell.SetToken("b", fmt.Sprintf("%d", len(g.boxes)))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	g.updateMoveAnimation()
	if g.cleared && g.move == nil {
		if vfxui.AnyPressStart() && !ate {
			*g = *newGame()
		}
		g.fx.Update()
		return nil
	}
	g.readInput(ate)
	g.fx.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 16, 28, 255})
	g.shell.FillStage(screen, color.RGBA{20, 28, 44, 255})
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			p := pt{x, y}
			c := color.RGBA{36, 48, 70, 255}
			if g.walls[p] {
				c = color.RGBA{70, 82, 110, 255}
			}
			vector.DrawFilledRect(screen, float32(g.ox+float64(x)*g.ts), float32(g.oy+float64(y)*g.ts),
				float32(g.ts-1), float32(g.ts-1), c, false)
		}
	}
	for _, gl := range g.goals {
		cx, cy := g.center(gl)
		pulse := 1 + 0.08*math.Sin(float64(g.frame)*0.08)
		vector.StrokeCircle(screen, float32(cx), float32(cy), float32(g.ts*0.28*pulse), 3, color.RGBA{80, 220, 160, 255}, false)
	}
	for i, b := range g.boxes {
		cx, cy := g.center(b)
		if g.move != nil && g.move.plan.BoxIndex == i {
			fromX, fromY := g.center(pt{g.move.plan.BoxFrom.X, g.move.plan.BoxFrom.Y})
			toX, toY := g.center(pt{g.move.plan.BoxTo.X, g.move.plan.BoxTo.Y})
			t := vfxmotion.EaseInOutCubic(g.move.tween.Progress())
			cx, cy = vfxmotion.Lerp(fromX, toX, t), vfxmotion.Lerp(fromY, toY, t)
		}
		col := color.RGBA{200, 140, 70, 255}
		if g.onGoal(b) {
			col = color.RGBA{120, 230, 160, 255}
		}
		vector.DrawFilledRect(screen, float32(cx-g.ts*0.32), float32(cy-g.ts*0.32),
			float32(g.ts*0.64), float32(g.ts*0.64), col, false)
	}
	px, py := g.center(g.player)
	pose := hero.Pose{}
	if g.move != nil {
		fromX, fromY := g.center(pt{g.move.plan.PlayerFrom.X, g.move.plan.PlayerFrom.Y})
		toX, toY := g.center(pt{g.move.plan.PlayerTo.X, g.move.plan.PlayerTo.Y})
		t := vfxmotion.EaseInOutCubic(g.move.tween.Progress())
		px, py = vfxmotion.Lerp(fromX, toX, t), vfxmotion.Lerp(fromY, toY, t)
		bob := math.Sin(t * math.Pi)
		py -= bob * g.ts * 0.08
		pose.ScaleX, pose.ScaleY = 1+0.08*bob, 1-0.06*bob
	}
	hero.DrawCenteredPose(screen, px, py, g.ts*0.85, pose)
	g.fx.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("moves %d", g.moves), 12, int(g.oy)-18)
	if g.cleared {
		ebitenutil.DebugPrintAt(screen, "CLEAR — tap to retry", 160, int(g.oy+float64(g.h)*g.ts/2))
	}
	g.shell.Hint = "arrows/WASD or tap toward move · map ≠ fx"
	g.shell.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("ADV: FX Sokoban — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
