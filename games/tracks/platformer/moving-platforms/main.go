package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float64 }
type platform struct {
	rect
	baseX, baseY, rangeX, rangeY, phase, oldX, oldY float64
}
type game struct {
	player          rect
	vx, vy          float64
	platforms       []platform
	frame           int
	grounded, clear bool
}

func newGame() *game {
	return &game{player: rect{28, 590, 28, 38}, platforms: []platform{
		{rect: rect{0, 650, 105, 24}, baseX: 0, baseY: 650},
		{rect: rect{125, 575, 90, 18}, baseX: 125, baseY: 575, rangeY: 70, phase: .5},
		{rect: rect{250, 500, 90, 18}, baseX: 250, baseY: 500, rangeX: 75, phase: 1.4},
		{rect: rect{365, 405, 85, 18}, baseX: 365, baseY: 405, rangeY: 85, phase: 2.2},
		{rect: rect{210, 300, 95, 18}, baseX: 210, baseY: 300, rangeX: 90, phase: 3.2},
		{rect: rect{40, 205, 115, 20}, baseX: 40, baseY: 205},
	}}
}

func (g *game) Update() error {
	if g.clear {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frame++
	standing := -1
	for i := range g.platforms {
		p := &g.platforms[i]
		p.oldX, p.oldY = p.x, p.y
		t := float64(g.frame)*.025 + p.phase
		p.x = p.baseX + math.Sin(t)*p.rangeX
		p.y = p.baseY + math.Sin(t)*p.rangeY
		if math.Abs((g.player.y+g.player.h)-p.oldY) < 3 && g.player.x+g.player.w > p.oldX && g.player.x < p.oldX+p.w {
			standing = i
		}
	}
	if standing >= 0 {
		p := g.platforms[standing]
		g.player.x += p.x - p.oldX
		g.player.y += p.y - p.oldY
	}
	left := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				left = true
			} else {
				right = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			jump = true
		}
	}
	if left {
		g.vx -= .65
	}
	if right {
		g.vx += .65
	}
	if !left && !right {
		g.vx *= .78
	}
	g.vx = clamp(g.vx, -5.3, 5.3)
	if jump && g.grounded {
		g.vy = -11.8
		g.grounded = false
	}
	g.vy = math.Min(g.vy+.62, 14)
	g.player.x = clamp(g.player.x+g.vx, 0, width-g.player.w)
	oldBottom := g.player.y + g.player.h
	g.player.y += g.vy
	g.grounded = false
	if g.vy >= 0 {
		for _, p := range g.platforms {
			newBottom := g.player.y + g.player.h
			if oldBottom <= p.y+4 && newBottom >= p.y && g.player.x+g.player.w > p.x && g.player.x < p.x+p.w {
				g.player.y = p.y - g.player.h
				g.vy = 0
				g.grounded = true
			}
		}
	}
	if g.player.y > height {
		*g = *newGame()
	}
	if g.player.y < 195 && g.player.x < 160 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{23, 45, 73, 255})
	for i := 0; i < 8; i++ {
		vector.DrawFilledCircle(s, float32(i*90), 660, 90, color.RGBA{15, 29, 49, 255}, false)
	}
	for _, p := range g.platforms {
		vector.DrawFilledRect(s, float32(p.x), float32(p.y), float32(p.w), float32(p.h), color.RGBA{43, 201, 177, 255}, false)
		vector.DrawFilledRect(s, float32(p.x+6), float32(p.y+5), float32(p.w-12), 4, color.RGBA{186, 255, 229, 210}, false)
	}
	vector.DrawFilledRect(s, 72, 165, 55, 35, color.RGBA{255, 210, 61, 255}, false)
	hero.DrawBottomCentered(s, g.player.x+g.player.w/2, g.player.y+g.player.h, g.player.h*1.55)
	ebitenutil.DebugPrintAt(s, "RIDE THE MOVING PLATFORMS TO THE FLAG", 95, 28)
	ebitenutil.DebugPrintAt(s, "BOTTOM TOUCH: MOVE    TOP TOUCH: JUMP", 105, 685)
	if g.clear {
		vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "VALLEY CLEARED!\n\nTAP / SPACE TO PLAY AGAIN", 145, 330)
	}
}
func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Moving Platform Valley — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
