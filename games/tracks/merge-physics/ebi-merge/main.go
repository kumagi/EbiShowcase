package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

//go:embed assets/pearl-nursery-v2.png
var nurseryPNG []byte

//go:embed assets/merge-creatures-v2.png
var creaturesPNG []byte

var nurseryArt *ebiten.Image
var creatureRects []image.Rectangle
var creatureSprites [7]*ebiten.Image

func loadGeneratedArt() {
	decode := func(data []byte) image.Image {
		im, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		return im
	}
	nursery := decode(nurseryPNG)
	creatures := decode(creaturesPNG)
	nurseryArt = ebiten.NewImageFromImage(nursery)
	creatureRects = alphaColumnCells(creatures, 7)
	subImager, ok := creatures.(interface {
		SubImage(image.Rectangle) image.Image
	})
	if !ok {
		panic("merge-creatures image does not support cropping")
	}
	for i, rect := range creatureRects {
		creatureSprites[i] = ebiten.NewImageFromImage(subImager.SubImage(rect))
	}
}

// alphaColumnCells finds each creature from its solid core, then expands the
// crop using every non-zero alpha pixel inside that creature's gutter-bounded
// cell. The old core-only crop discarded faint antialiasing and glow pixels,
// which made a hard rectangular cut visible after downscaling. A creature may
// contain disconnected details (a pearl, antenna, or crown), so connected-
// component detection is not a valid sprite boundary.
func alphaColumnCells(atlas image.Image, cells int) []image.Rectangle {
	bounds := atlas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	alpha := func(x, y int) uint32 {
		_, _, _, a := atlas.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
		return a
	}
	columnHasCoreInk := func(x int) bool {
		for y := 0; y < h; y++ {
			if alpha(x, y) >= 0x1000 {
				return true
			}
		}
		return false
	}
	type run struct{ minX, maxX int }
	runs := make([]run, 0, cells)
	start := -1
	for x := 0; x <= w; x++ {
		hasInk := x < w && columnHasCoreInk(x)
		if hasInk && start < 0 {
			start = x
		}
		if !hasInk && start >= 0 {
			runs = append(runs, run{start, x})
			start = -1
		}
	}
	if len(runs) != cells {
		panic(fmt.Sprintf("merge-creatures atlas has %d transparent-gutter cells, want %d", len(runs), cells))
	}

	result := make([]image.Rectangle, cells)
	for i, r := range runs {
		cellMinX, cellMaxX := 0, w
		if i > 0 {
			cellMinX = (runs[i-1].maxX + r.minX) / 2
		}
		if i+1 < len(runs) {
			cellMaxX = (r.maxX + runs[i+1].minX) / 2
		}
		minX, minY, maxX, maxY := cellMaxX, h, cellMinX, 0
		for x := cellMinX; x < cellMaxX; x++ {
			for y := 0; y < h; y++ {
				if alpha(x, y) > 0 {
					minX = min(minX, x)
					minY = min(minY, y)
					maxX = max(maxX, x+1)
					maxY = max(maxY, y+1)
				}
			}
		}
		if minX >= maxX || minY >= maxY {
			panic(fmt.Sprintf("merge-creatures atlas cell %d has no visible pixels", i))
		}
		// Keep a generous transparent apron so linear filtering never clamps a
		// glow, whisker, fin, or shadow against the crop edge.
		const apron = 12
		result[i] = image.Rect(max(cellMinX, minX-apron), minY-apron, min(cellMaxX, maxX+apron), maxY+apron).
			Intersect(image.Rect(0, 0, w, h)).Add(bounds.Min)
	}
	return result
}

const width, height = 480, 720

const (
	maxTier      = 6
	dangerLineY  = 120
	dangerLimit  = 180
	worldGravity = .48
)

// Each merge roughly preserves area: the next radius is about sqrt(2) times
// the previous one. The old late tiers grew too slowly, so two large pieces
// collapsed into much less occupied area and made the board unusually easy.
var radii = [...]float64{15, 21, 30, 42, 60, 85, 120}

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
	fruits                               []fruit
	rng                                  *rand.Rand
	next, after, score, danger, cooldown int
	cursor                               float64
	over                                 bool
	combo, comboTimer, shake, best       int
	sparks                               []spark
	audio                                *audio.Context
	gate                                 audiolab.Gate
	pulse                                *shaderlab.Pulse
	cam                                  cameralab.State
	badge                                *ebiten.Image
}

func newGame() *game {
	if nurseryArt == nil {
		loadGeneratedArt()
	}
	b := ebiten.NewImage(20, 20)
	b.Fill(color.RGBA{255, 130, 80, 255})
	g := &game{rng: rand.New(rand.NewSource(4806)), cursor: 240, audio: audiolab.Context(), pulse: shaderlab.NewPulse(), cam: cameralab.State{Pos: cameralab.Vec{240, 360}, ViewW: width, ViewH: height}, badge: b}
	g.next = g.rng.Intn(3)
	g.after = g.rng.Intn(3)
	// A prepared opening board communicates the merge goal immediately.
	g.fruits = []fruit{{x: 130, y: 642, tier: 1}, {x: 175, y: 638, tier: 1}, {x: 310, y: 632, tier: 2}, {x: 368, y: 640, tier: 0}}
	return g
}
func (g *game) Update() error {
	if g.over {
		if any() {
			best := g.best
			*g = *newGame()
			g.best = best
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
		f.vy += worldGravity
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
				if pass == 0 && a.tier == b.tier && a.tier < maxTier {
					tier := a.tier + 1
					a.dead, b.dead = true, true
					spawn = append(spawn, fruit{x: (a.x + b.x) / 2, y: (a.y + b.y) / 2, vx: (a.vx + b.vx) / 2, vy: -2.5, tier: tier})
					g.gate.Arm(true)
					g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, 420+float64(tier)*80, .09)).Play()
					g.combo++
					g.comboTimer = 100
					g.score += mergePoints(tier, g.combo)
					g.best = max(g.best, g.score)
					g.shake = min(8, 3+tier)
					g.burst((a.x+b.x)/2, (a.y+b.y)/2, 8+tier*2)
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
		if f.y-radii[f.tier] < dangerLineY && math.Hypot(f.vx, f.vy) < 1.2 {
			above = true
		}
	}
	if above {
		g.danger++
	} else {
		g.danger = max(0, g.danger-3)
	}
	if g.danger >= dangerLimit {
		g.over = true
	}
	return nil
}

func mergePoints(tier, combo int) int {
	return (tier+1)*(tier+1)*10 + combo*15
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 28 + float64(i%9), color.RGBA{255, 211, 62, 255}})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	if g.pulse.Available() {
		fx := ebiten.NewImage(20, 20)
		if g.pulse.Draw(fx, g.badge, float32(g.score)*.01) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(440, 12)
			s.DrawImage(fx, op)
		}
	}
	drawCover(s, nurseryArt)
	vector.DrawFilledRect(s, 0, 0, width, height, color.RGBA{2, 13, 35, 50}, false)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.cooldown+g.comboTimer)*2) * 5
	}
	// The tank is glass, not an opaque debug rectangle: keep the generated
	// nursery visible while a restrained tint preserves sprite contrast.
	vector.DrawFilledRect(s, 20+float32(ox), 70, 440, 615, color.RGBA{10, 28, 56, 128}, false)
	vector.StrokeRect(s, 20+float32(ox), 70, 440, 615, 5, color.RGBA{245, 184, 84, 180}, false)
	// Glass shine and bin lip sell the physical toy-box presentation.
	vector.DrawFilledRect(s, 30+float32(ox), 82, 18, 570, color.RGBA{255, 255, 255, 18}, false)
	vector.DrawFilledRect(s, 30+float32(ox), 650, 420, 25, color.RGBA{91, 65, 73, 230}, false)
	line := color.RGBA{240, 90, 95, 90}
	if g.danger > 0 {
		line = color.RGBA{255, 70, 75, 220}
	}
	vector.StrokeLine(s, 25, dangerLineY, 455, dangerLineY, 3, line, false)
	for _, f := range g.fruits {
		r := radii[f.tier]
		pulse := 1.0
		if g.comboTimer > 90 {
			pulse = 1.08
		}
		drawCreature(s, f.tier, f.x+ox, f.y, r*2*pulse)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/15), p.c, true)
	}
	vector.StrokeLine(s, float32(g.cursor), 75, float32(g.cursor), 135, 2, color.RGBA{255, 255, 255, 130}, false)
	drawCreature(s, g.next, g.cursor, 92, radii[g.next]*2)
	vector.DrawFilledRect(s, 340, 82, 96, 70, color.RGBA{7, 13, 28, 210}, false)
	ebitenutil.DebugPrintAt(s, "UP NEXT", 355, 91)
	drawCreature(s, g.after, 390, 126, radii[g.after]*1.2)
	vector.DrawFilledRect(s, 12, 10, 456, 50, color.RGBA{5, 12, 27, 225}, false)
	drawCreature(s, maxTier, 34, 35, 40)
	if f, e := uilab.Face("en", 16); e == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(55, 25)
		text.Draw(s, fmt.Sprintf("SCORE %05d  BEST %05d  COMBO x%d", g.score, g.best, g.combo), f, op)
	} else {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %05d  BEST %05d  COMBO x%d", g.score, g.best, g.combo), 55, 25)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EVERY MERGE SCORES          DANGER %03d/%03d", g.danger, dangerLimit), 86, 48)
	ebitenutil.DebugPrintAt(s, "MOVE POINTER, TAP TO DROP — KEEP MERGING", 75, 700)
	if g.over {
		overlay(s, "STACK CROSSED THE LINE!\n\nTAP / SPACE TO RETRY")
	}
}

func drawCreature(dst *ebiten.Image, tier int, cx, cy, size float64) {
	if tier < 0 || tier >= len(radii) || creatureSprites[tier] == nil {
		return
	}
	src := creatureSprites[tier]
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	scale := size / float64(max(w, h))
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(cx-float64(w)*scale/2, cy-float64(h)*scale/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(src, op)
}

func drawCover(dst, src *ebiten.Image) {
	w, h := float64(src.Bounds().Dx()), float64(src.Bounds().Dy())
	scale := math.Max(width/w, height/h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((width-w*scale)/2, (height-h*scale)/2)
	dst.DrawImage(src, op)
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
