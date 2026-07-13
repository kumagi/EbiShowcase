package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float64 }
type foe struct {
	rect
	vx    float64
	alive bool
}
type spark struct {
	x, y, vx, vy, life float64
	c                  color.RGBA
}
type game struct {
	p                                rect
	vx, vy, camera                   float64
	grounds                          []rect
	coins                            []rect
	power                            rect
	foes                             []foe
	score, stage, life               int
	grounded, big, powerTaken, clear bool
	tick, finishTimer                int
	sparks                           []spark
}

func newGame() *game { g := &game{stage: 1, life: 3}; g.load(); return g }
func (g *game) load() {
	g.p = rect{35, 580, 28, 38}
	g.vx, g.vy, g.camera = 0, 0, 0
	stages := [][]rect{
		{{0, 640, 350, 80}, {410, 570, 260, 150}, {730, 640, 300, 80}, {1090, 550, 230, 170}, {1380, 640, 320, 80}, {200, 500, 100, 20}, {520, 430, 110, 20}, {830, 500, 100, 20}, {1150, 400, 100, 20}, {1450, 490, 110, 20}},
		{{0, 640, 220, 80}, {300, 560, 180, 160}, {570, 470, 150, 250}, {810, 590, 230, 130}, {1130, 480, 190, 240}, {1400, 640, 300, 80}, {130, 450, 90, 18}, {390, 360, 100, 18}, {680, 300, 100, 18}, {990, 400, 100, 18}},
		{{0, 640, 280, 80}, {360, 640, 170, 80}, {610, 640, 180, 80}, {870, 640, 170, 80}, {1120, 640, 180, 80}, {1380, 640, 320, 80}, {180, 500, 90, 18}, {430, 430, 90, 18}, {690, 350, 90, 18}, {950, 430, 90, 18}, {1210, 350, 90, 18}},
		{{0, 640, 180, 80}, {260, 540, 170, 180}, {520, 430, 170, 290}, {780, 320, 170, 400}, {1040, 430, 170, 290}, {1300, 540, 170, 180}, {1550, 640, 150, 80}, {110, 450, 80, 18}, {1430, 390, 80, 18}},
	}
	g.grounds = stages[g.stage-1]
	g.coins = nil
	coinSets := [][]float64{{230, 550, 850, 1180, 1480}, {150, 410, 660, 930, 1180, 1510}, {200, 450, 700, 960, 1210, 1480}, {130, 350, 610, 870, 1130, 1390, 1600}}
	for i, x := range coinSets[g.stage-1] {
		g.coins = append(g.coins, rect{x, 380 - float64((i+g.stage)%3)*55, 14, 14})
	}
	foeSets := [][]foe{
		{{rect{470, 542, 28, 28}, 1.2, true}, {rect{780, 612, 28, 28}, -1.3, true}, {rect{1420, 612, 28, 28}, 1.4, true}},
		{{rect{320, 532, 28, 28}, 1.4, true}, {rect{840, 562, 28, 28}, -1.5, true}, {rect{1160, 452, 28, 28}, 1.3, true}},
		{{rect{380, 612, 28, 28}, 1.7, true}, {rect{640, 612, 28, 28}, -1.7, true}, {rect{900, 612, 28, 28}, 1.8, true}, {rect{1410, 612, 28, 28}, -1.6, true}},
		{{rect{290, 512, 28, 28}, 1.8, true}, {rect{550, 402, 28, 28}, -1.8, true}, {rect{810, 292, 28, 28}, 2, true}, {rect{1070, 402, 28, 28}, -1.9, true}, {rect{1330, 512, 28, 28}, 2, true}},
	}
	g.foes = foeSets[g.stage-1]
	g.power = rect{[]float64{860, 650, 950, 1120}[g.stage-1], []float64{470, 270, 400, 380}[g.stage-1], 22, 22}
	g.powerTaken = false
	g.big = false
}
func (g *game) Update() error {
	g.tick++
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .08
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.clear {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	l := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	r := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	j := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				l = true
			} else {
				r = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			j = true
		}
	}
	if l {
		g.vx -= .7
	}
	if r {
		g.vx += .7
	}
	if !l && !r {
		g.vx *= .8
	}
	g.vx = clamp(g.vx, -6, 6)
	if j && g.grounded {
		g.vy = -12.5
		g.burst(g.p.x+g.p.w/2, g.p.y+g.p.h, color.RGBA{235, 224, 170, 255}, 5)
	}
	g.vy = math.Min(g.vy+.65, 14)
	g.p.x = clamp(g.p.x+g.vx, 0, 1670)
	old := g.p.y + g.p.h
	g.p.y += g.vy
	g.grounded = false
	for _, b := range g.grounds {
		if g.vy >= 0 && old <= b.y+3 && g.p.y+g.p.h >= b.y && g.p.x+g.p.w > b.x && g.p.x < b.x+b.w {
			g.p.y = b.y - g.p.h
			g.vy = 0
			g.grounded = true
		}
	}
	for i := len(g.coins) - 1; i >= 0; i-- {
		if overlap(g.p, g.coins[i]) {
			g.burst(g.coins[i].x+7, g.coins[i].y+7, color.RGBA{255, 220, 62, 255}, 8)
			g.coins = append(g.coins[:i], g.coins[i+1:]...)
			g.score += 100
		}
	}
	if !g.powerTaken && overlap(g.p, g.power) {
		g.burst(g.power.x+11, g.power.y+11, color.RGBA{104, 255, 168, 255}, 16)
		g.powerTaken = true
		g.big = true
		g.p.y -= 16
		g.p.h = 54
	}
	for i := range g.foes {
		e := &g.foes[i]
		if !e.alive {
			continue
		}
		e.x += e.vx
		if e.x < 0 || e.x > 1670 {
			e.vx = -e.vx
		}
		if overlap(g.p, e.rect) {
			if g.vy > 0 && old <= e.y+8 {
				g.burst(e.x+14, e.y+14, color.RGBA{255, 154, 88, 255}, 12)
				e.alive = false
				g.vy = -8
				g.score += 200
			} else if g.big {
				g.big = false
				g.p.h = 38
				e.vx = -e.vx
				g.p.x -= 20
			} else {
				g.life--
				if g.life <= 0 {
					*g = *newGame()
				} else {
					g.load()
				}
				return nil
			}
		}
	}
	if g.p.y > height {
		g.life--
		if g.life <= 0 {
			*g = *newGame()
		} else {
			g.load()
		}
	}
	if g.p.x > 1640 {
		if g.stage < 4 {
			g.stage++
			g.score += 500
			g.load()
		} else {
			g.clear = true
		}
	}
	target := g.p.x - width*.4
	g.camera = clamp(g.camera+(target-g.camera)*.08, 0, 1220)
	return nil
}
func (g *game) burst(x, y float64, c color.RGBA, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * (1 + float64(i%3)), math.Sin(a) * (1 + float64(i%3)), 25 + float64(i%8), c})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	skies := []color.RGBA{{102, 189, 231, 255}, {99, 91, 173, 255}, {241, 151, 98, 255}, {28, 40, 88, 255}}
	sky := skies[g.stage-1]
	s.Fill(sky)
	for _, b := range g.grounds {
		x := b.x - g.camera
		if x+b.w < 0 || x > width {
			continue
		}
		vector.DrawFilledRect(s, float32(x), float32(b.y), float32(b.w), float32(b.h), color.RGBA{55, 101, 66, 255}, false)
		for tx := 0.0; tx < b.w; tx += 40 {
			w := math.Min(40, b.w-tx)
			trackatlas.Draw(s, "tile-platform", x+tx, b.y, w)
		}
	}
	for _, c := range g.coins {
		pulse := 20 + math.Sin(float64(g.tick)*.12+c.x)*3
		trackatlas.DrawCentered(s, "coin", c.x-g.camera+7, c.y+7, pulse)
	}
	if !g.powerTaken {
		trackatlas.DrawCentered(s, "power-star", g.power.x-g.camera+11, g.power.y+11, 26)
	}
	for _, e := range g.foes {
		if e.alive {
			bob := math.Sin(float64(g.tick)*.16+e.x) * 2
			trackatlas.DrawCentered(s, "slug", e.x-g.camera+14, e.y+14+bob, 30)
		}
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x-g.camera), float32(p.y), float32(2+p.life/14), p.c, true)
	}
	runBob := 0.0
	if g.grounded && math.Abs(g.vx) > .3 {
		runBob = math.Abs(math.Sin(float64(g.tick)*.32)) * 4
	}
	heroSize := g.p.h - runBob
	trackatlas.DrawCentered(s, "hero", g.p.x-g.camera+g.p.w/2, g.p.y+g.p.h-heroSize/2, heroSize)
	flag := 1650 - g.camera
	trackatlas.Draw(s, "flag", flag, 480, 140)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE %d/4   LIFE %d   SCORE %05d   COINS %d", g.stage, g.life, g.score, len(g.coins)), 45, 22)
	ebitenutil.DebugPrintAt(s, "GREEN ORB MAKES EBI BIG", 150, 48)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH    JUMP: SPACE OR UPPER TOUCH", 50, 685)
	if g.clear {
		vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "EBI ADVENTURE COMPLETE!\n\nTAP / SPACE TO PLAY AGAIN", 125, 330)
	}
}
func overlap(a, b rect) bool        { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Adventure — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
