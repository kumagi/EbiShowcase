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
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const width, height = 480, 720

var radii = []float64{15, 20, 27, 36, 47, 60, 74}

type fruit struct {
	x, y, vx, vy float64
	tier         int
	dead         bool
}
type spark struct {
	x, y, vx, vy, life float64
	c                  color.RGBA
}
type game struct {
	fruits                                []fruit
	rng                                   *rand.Rand
	next, after, score, danger, cooldown  int
	cursor                                float64
	clear, over                           bool
	level, combo, comboTimer, shake, best int
	sparks                                []spark
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(4806)), cursor: 240, level: 1}
	g.next = g.rng.Intn(3)
	g.after = g.rng.Intn(3)
	return g
}
func (g *game) Update() error {
	if g.clear || g.over {
		if any() {
			if g.clear && g.level < 3 {
				g.level++
				g.fruits = nil
				g.danger = 0
				g.clear = false
				g.next = g.rng.Intn(3)
				g.after = g.rng.Intn(3)
			} else {
				best := g.best
				*g = *newGame()
				g.best = best
			}
		}
		return nil
	}
	if g.cooldown > 0 {
		g.cooldown--
	}
	if g.comboTimer > 0 {
		g.comboTimer--
	} else {
		g.combo = 0
	}
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .04
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if x, ok := pointerX(); ok {
		g.cursor = math.Max(45, math.Min(435, float64(x)))
	}
	if x, ok := drop(); ok && g.cooldown == 0 {
		g.cursor = math.Max(45, math.Min(435, float64(x)))
		g.fruits = append(g.fruits, fruit{x: g.cursor, y: 96, tier: g.next})
		g.next, g.after = g.after, g.rng.Intn(3)
		g.cooldown = 18
	}
	for i := range g.fruits {
		f := &g.fruits[i]
		f.vy += []float64{.40, .48, .56}[g.level-1]
		f.x += f.vx
		f.y += f.vy
		r := radii[f.tier]
		if f.x-r < 25 {
			f.x = 25 + r
			f.vx = math.Abs(f.vx) * .2
		}
		if f.x+r > 455 {
			f.x = 455 - r
			f.vx = -math.Abs(f.vx) * .2
		}
		if f.y+r > 665 {
			f.y = 665 - r
			if f.vy > 1 {
				f.vy = -f.vy * .12
			} else {
				f.vy = 0
			}
			f.vx *= .84
		}
	}
	spawn := []fruit{}
	for pass := 0; pass < 6; pass++ {
		for i := 0; i < len(g.fruits); i++ {
			if g.fruits[i].dead {
				continue
			}
			for j := i + 1; j < len(g.fruits); j++ {
				if g.fruits[j].dead {
					continue
				}
				a, b := &g.fruits[i], &g.fruits[j]
				dx, dy := b.x-a.x, b.y-a.y
				d := math.Hypot(dx, dy)
				minD := radii[a.tier] + radii[b.tier]
				if d <= 0 || d >= minD {
					continue
				}
				if pass == 0 && a.tier == b.tier && a.tier < 6 {
					tier := a.tier + 1
					a.dead, b.dead = true, true
					spawn = append(spawn, fruit{x: (a.x + b.x) / 2, y: (a.y + b.y) / 2, vx: (a.vx + b.vx) / 2, vy: -2.5, tier: tier})
					g.score += (tier + 1) * (tier + 1) * 10
					g.combo++
					g.comboTimer = 100
					g.score += g.combo * 15
					g.shake = min(8, 3+tier)
					g.burst((a.x+b.x)/2, (a.y+b.y)/2, 8+tier*2)
					target := []int{4, 5, 6}[g.level-1]
					if tier == target {
						g.clear = true
						if g.score > g.best {
							g.best = g.score
						}
					}
					break
				}
				nx, ny := dx/d, dy/d
				over := minD - d
				a.x -= nx * over / 2
				a.y -= ny * over / 2
				b.x += nx * over / 2
				b.y += ny * over / 2
				rel := (b.vx-a.vx)*nx + (b.vy-a.vy)*ny
				if rel < 0 {
					imp := -rel * .55
					a.vx -= imp * nx
					a.vy -= imp * ny
					b.vx += imp * nx
					b.vy += imp * ny
				}
			}
		}
	}
	alive := g.fruits[:0]
	for _, f := range g.fruits {
		if !f.dead {
			alive = append(alive, f)
		}
	}
	g.fruits = append(alive, spawn...)
	above := false
	for _, f := range g.fruits {
		if f.y-radii[f.tier] < 175 && math.Hypot(f.vx, f.vy) < 1.2 {
			above = true
		}
	}
	if above {
		g.danger++
	} else {
		g.danger = max(0, g.danger-3)
	}
	if g.danger >= 180 {
		g.over = true
	}
	return nil
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 28 + float64(i%9), color.RGBA{255, 211, 62, 255}})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	bgs := []color.RGBA{{18, 28, 44, 255}, {39, 28, 55, 255}, {52, 29, 31, 255}}
	s.Fill(bgs[g.level-1])
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.cooldown+g.comboTimer)*2) * 5
	}
	vector.DrawFilledRect(s, 20+float32(ox), 70, 440, 615, color.RGBA{35, 44, 61, 255}, false)
	line := color.RGBA{240, 90, 95, 90}
	if g.danger > 0 {
		line = color.RGBA{255, 70, 75, 220}
	}
	vector.StrokeLine(s, 25, 175, 455, 175, 3, line, false)
	for _, f := range g.fruits {
		r := radii[f.tier]
		pulse := 1.0
		if g.comboTimer > 90 {
			pulse = 1.08
		}
		trackatlas.DrawCentered(s, trackatlas.Merge(f.tier+1), f.x+ox, f.y, r*2*pulse)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d", f.tier+1), int(f.x)-3, int(f.y)-5)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/15), p.c, true)
	}
	vector.StrokeLine(s, float32(g.cursor), 75, float32(g.cursor), 135, 2, color.RGBA{255, 255, 255, 130}, false)
	trackatlas.DrawCentered(s, trackatlas.Merge(g.next+1), g.cursor, 92, radii[g.next]*2)
	target := []int{5, 6, 7}[g.level-1]
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE %d/3 SCORE %05d BEST %05d COMBO x%d", g.level, g.score, g.best, g.combo), 55, 25)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("NEXT %d AFTER %d DANGER %03d/180 TARGET TIER %d", g.next+1, g.after+1, g.danger, target), 65, 48)
	ebitenutil.DebugPrintAt(s, "MOVE POINTER, TAP TO DROP — BUILD THE TARGET", 65, 700)
	if g.clear {
		msg := "STAGE CLEAR!\n\nTAP / SPACE FOR NEXT STAGE"
		if g.level == 3 {
			msg = "EBI MERGE COMPLETE!\n\nTAP / SPACE TO PLAY AGAIN"
		}
		overlay(s, msg)
	} else if g.over {
		overlay(s, "STACK CROSSED THE LINE!\n\nTAP / SPACE TO RETRY")
	}
}
func pointerX() (int, bool) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		return x, true
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		return x, true
	}
	return 0, false
}
func drop() (int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		return x, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		return x, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return int(240), true
	}
	return 0, false
}
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 115, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Merge — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
