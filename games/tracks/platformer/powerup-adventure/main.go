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
type game struct {
	p                                rect
	vx, vy, camera                   float64
	grounds                          []rect
	coins                            []rect
	power                            rect
	foes                             []foe
	score, stage, life               int
	grounded, big, powerTaken, clear bool
}

func newGame() *game { g := &game{stage: 1, life: 3}; g.load(); return g }
func (g *game) load() {
	g.p = rect{35, 580, 28, 38}
	g.vx, g.vy, g.camera = 0, 0, 0
	g.grounds = []rect{{0, 640, 350, 80}, {410, 570, 260, 150}, {730, 640, 300, 80}, {1090, 550, 230, 170}, {1380, 640, 320, 80}, {200, 500, 100, 20}, {520, 430, 110, 20}, {830, 500, 100, 20}, {1150, 400, 100, 20}, {1450, 490, 110, 20}}
	g.coins = nil
	for _, x := range []float64{230, 550, 850, 1180, 1480} {
		g.coins = append(g.coins, rect{x, 365 + math.Mod(x, 130), 14, 14})
	}
	g.foes = []foe{{rect{470, 542, 28, 28}, 1.2, true}, {rect{780, 612, 28, 28}, -1.3, true}, {rect{1420, 612, 28, 28}, 1.4, true}}
	g.power = rect{860, 470, 22, 22}
	g.powerTaken = false
	g.big = false
}
func (g *game) Update() error {
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
			g.coins = append(g.coins[:i], g.coins[i+1:]...)
			g.score += 100
		}
	}
	if !g.powerTaken && overlap(g.p, g.power) {
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
				g.load()
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
		if g.stage == 1 {
			g.stage = 2
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
func (g *game) Draw(s *ebiten.Image) {
	sky := color.RGBA{102, 189, 231, 255}
	if g.stage == 2 {
		sky = color.RGBA{99, 91, 173, 255}
	}
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
		trackatlas.DrawCentered(s, "coin", c.x-g.camera+7, c.y+7, 20)
	}
	if !g.powerTaken {
		trackatlas.DrawCentered(s, "power-star", g.power.x-g.camera+11, g.power.y+11, 26)
	}
	for _, e := range g.foes {
		if e.alive {
			trackatlas.DrawCentered(s, "slug", e.x-g.camera+14, e.y+14, 30)
		}
	}
	trackatlas.DrawCentered(s, "hero", g.p.x-g.camera+g.p.w/2, g.p.y+g.p.h/2, g.p.h)
	flag := 1650 - g.camera
	trackatlas.Draw(s, "flag", flag, 480, 140)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE %d/2   LIFE %d   SCORE %05d   COINS %d", g.stage, g.life, g.score, len(g.coins)), 45, 22)
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
