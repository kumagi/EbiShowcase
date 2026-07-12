// vfx-fx-sokoban — ADVANCED: tile map play vs fx dust/sparkle.
package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

type pt struct{ x, y int }

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
	np := pt{g.player.x + dx, g.player.y + dy}
	if g.walls[np] {
		return
	}
	if i := g.boxAt(np); i >= 0 {
		bp := pt{np.x + dx, np.y + dy}
		if g.walls[bp] || g.boxAt(bp) >= 0 {
			return
		}
		g.boxes[i] = bp
		cx, cy := g.center(np)
		fx := g.shell.Get("fx")
		n := int(g.shell.Get("dust") * fx)
		g.fx.Burst(cx, cy, n, 1.8*fx, color.RGBA{180, 160, 120, 255}, false)
		if g.onGoal(bp) {
			gx, gy := g.center(bp)
			g.fx.Ring(gx, gy, 1.1*fx, color.RGBA{120, 255, 200, 255})
			g.fx.Burst(gx, gy, int(12*fx), 2.2*fx, color.RGBA{180, 255, 220, 255}, true)
		}
	}
	g.player = np
	g.moves++
	g.cleared = true
	for _, b := range g.boxes {
		if !g.onGoal(b) {
			g.cleared = false
			break
		}
	}
	if g.cleared {
		g.fx.FlashScreen(0.5*g.shell.Get("fx"), 120, 255, 200)
	}
}

func (g *Game) readInput(ate bool) {
	if ate || g.cleared {
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
	g.shell.SetToken("b", fmt.Sprintf("%d", len(g.boxes)))
	g.shell.SetToken("n", fmt.Sprintf("%d", g.fx.Count()))
	if g.cleared {
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
		vector.StrokeCircle(screen, float32(cx), float32(cy), float32(g.ts*0.28), 3, color.RGBA{80, 220, 160, 255}, false)
	}
	for _, b := range g.boxes {
		cx, cy := g.center(b)
		col := color.RGBA{200, 140, 70, 255}
		if g.onGoal(b) {
			col = color.RGBA{120, 230, 160, 255}
		}
		vector.DrawFilledRect(screen, float32(cx-g.ts*0.32), float32(cy-g.ts*0.32),
			float32(g.ts*0.64), float32(g.ts*0.64), col, false)
	}
	px, py := g.center(g.player)
	hero.DrawCentered(screen, px, py, g.ts*0.85)
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
