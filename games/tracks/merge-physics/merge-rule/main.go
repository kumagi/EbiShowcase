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
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
)

const width, height = 480, 720

var radii = []float64{16, 22, 30, 40, 52, 67}
var colors = []color.RGBA{{79, 210, 181, 255}, {250, 187, 65, 255}, {241, 104, 85, 255}, {151, 92, 215, 255}, {82, 166, 235, 255}, {255, 222, 95, 255}}

type fruit struct {
	x, y, vx, vy float64
	tier         int
	dead         bool
	born         int
}

type mergeVisual struct {
	fromAX, fromAY float64
	fromBX, fromBY float64
	toX, toY       float64
	tier           int
	tween          vfxmotion.Tween
}

type game struct {
	fruits              []fruit
	rng                 *rand.Rand
	next, score, merges int
	clear               bool
	mergeVisuals        []mergeVisual
	fx                  vfxfx.System
}

func newGame() *game { return &game{rng: rand.New(rand.NewSource(4705))} }
func (g *game) Update() error {
	g.updatePresentation()
	if g.clear {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	if x, ok := drop(); ok {
		g.fruits = append(g.fruits, fruit{x: float64(x), y: 85, tier: g.next})
		g.next = g.rng.Intn(2)
	}
	for i := range g.fruits {
		f := &g.fruits[i]
		if f.born > 0 {
			f.born--
		}
		f.vy += .46
		f.x += f.vx
		f.y += f.vy
		r := radii[f.tier]
		if f.x-r < 25 {
			f.x = 25 + r
			f.vx = math.Abs(f.vx) * .25
		}
		if f.x+r > 455 {
			f.x = 455 - r
			f.vx = -math.Abs(f.vx) * .25
		}
		if f.y+r > 660 {
			f.y = 660 - r
			if f.vy > 1 {
				f.vy = -f.vy * .15
			} else {
				f.vy = 0
			}
			f.vx *= .85
		}
	}
	spawned := []fruit{}
	for pass := 0; pass < 5; pass++ {
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
				minimum := radii[a.tier] + radii[b.tier]
				if d <= 0 || d >= minimum {
					continue
				}
				if pass == 0 && a.tier == b.tier && a.tier < 5 {
					tier := a.tier + 1
					x, y := (a.x+b.x)/2, (a.y+b.y)/2
					g.mergeVisuals = append(g.mergeVisuals, mergeVisual{
						fromAX: a.x, fromAY: a.y, fromBX: b.x, fromBY: b.y,
						toX: x, toY: y, tier: a.tier, tween: vfxmotion.NewTween(14),
					})
					spawned = append(spawned, fruit{x: x, y: y, vx: (a.vx + b.vx) / 2, vy: -2, tier: tier, born: 14})
					a.dead, b.dead = true, true
					g.merges++
					g.score += 10 * (tier + 1)
					g.fx.Shockwave(x, y, 0.5+float64(tier)*0.12, color.White, colors[tier])
					g.fx.Burst(x, y, 10+tier*3, 1.8+float64(tier)*0.25, colors[tier], true)
					if tier == 5 {
						g.clear = true
					}
					break
				}
				nx, ny := dx/d, dy/d
				over := minimum - d
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
	g.fruits = append(alive, spawned...)
	return nil
}

func (g *game) updatePresentation() {
	alive := g.mergeVisuals[:0]
	for i := range g.mergeVisuals {
		g.mergeVisuals[i].tween.Advance()
		if !g.mergeVisuals[i].tween.Done() {
			alive = append(alive, g.mergeVisuals[i])
		}
	}
	g.mergeVisuals = alive
	g.fx.Update()
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 28, 44, 255})
	vector.DrawFilledRect(s, 20, 70, 440, 610, color.RGBA{35, 44, 61, 255}, false)
	vector.StrokeRect(s, 20, 70, 440, 610, 5, color.RGBA{110, 128, 151, 255}, false)
	for _, f := range g.fruits {
		r := radii[f.tier]
		scale := 1.0
		if f.born > 0 {
			progress := 1 - float64(f.born)/14
			scale = 0.55 + vfxmotion.EaseOutCubic(progress)*0.45 + math.Sin(progress*math.Pi)*0.12
		}
		vector.DrawFilledCircle(s, float32(f.x), float32(f.y), float32(r*scale), colors[f.tier], false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d", f.tier+1), int(f.x)-3, int(f.y)-5)
	}
	for _, visual := range g.mergeVisuals {
		t := vfxmotion.EaseInOutCubic(visual.tween.Progress())
		ax := vfxmotion.Lerp(visual.fromAX, visual.toX, t)
		ay := vfxmotion.Lerp(visual.fromAY, visual.toY, t)
		bx := vfxmotion.Lerp(visual.fromBX, visual.toX, t)
		by := vfxmotion.Lerp(visual.fromBY, visual.toY, t)
		r := radii[visual.tier] * (1 - 0.65*t)
		vector.DrawFilledCircle(s, float32(ax), float32(ay), float32(r), colors[visual.tier], false)
		vector.DrawFilledCircle(s, float32(bx), float32(by), float32(r), colors[visual.tier], false)
	}
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("NEXT TIER %d   MERGES %02d   SCORE %04d", g.next+1, g.merges, g.score), 85, 28)
	ebitenutil.DebugPrintAt(s, "TAP TO DROP — MATCH EQUAL NUMBERS", 105, 690)
	if g.clear {
		overlay(s, "TIER 6 CREATED!\n\nTAP / SPACE TO RESET")
	}
}
func drop() (int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		return max(50, min(430, x)), true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		return max(50, min(430, x)), true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return 240, true
	}
	return 0, false
}
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 140, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Merge Rule — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
