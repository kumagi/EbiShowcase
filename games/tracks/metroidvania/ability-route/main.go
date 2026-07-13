package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	W         = 480
	H         = 720
	worldW    = 1320.0
	floorY    = 590.0
	playerW   = 30.0
	playerH   = 42.0
	moveSpeed = 4.2
)

type platform struct{ x, y, w, h float64 }

var platforms = []platform{
	{0, floorY, worldW, 130},
	{75, 455, 170, 22},
	{0, 315, 145, 22},
	{1060, 500, 160, 22},
}

type game struct {
	x, y, vx, vy float64
	camX         float64
	onGround     bool
	hasWings     bool
	jumpsLeft    int
	won          bool
	tick         int
	message      string
	particles    []particle
}

type particle struct{ x, y, vx, vy, life float64 }

func newGame() *game {
	return &game{x: 355, y: floorY - playerH, message: "Go RIGHT and collect the DOUBLE JUMP, then return LEFT."}
}

func (g *game) Update() error {
	g.tick++
	if g.won && retryPressed() {
		*g = *newGame()
		return nil
	}

	left, right, jump := controls()
	g.vx = 0
	if left {
		g.vx = -moveSpeed
	}
	if right {
		g.vx = moveSpeed
	}
	if jump && (g.onGround || g.jumpsLeft > 0) {
		g.vy = -11.2
		if !g.onGround {
			g.jumpsLeft--
			g.spawnBurst(g.x+playerW/2, g.y+playerH, 10)
		}
		g.onGround = false
	}

	g.x += g.vx
	g.x = math.Max(0, math.Min(worldW-playerW, g.x))
	g.vy += .58
	if g.vy > 13 {
		g.vy = 13
	}
	oldBottom := g.y + playerH
	g.y += g.vy
	g.onGround = false
	for _, p := range platforms {
		if g.x+playerW <= p.x || g.x >= p.x+p.w {
			continue
		}
		newBottom := g.y + playerH
		if g.vy >= 0 && oldBottom <= p.y && newBottom >= p.y {
			g.y, g.vy, g.onGround = p.y-playerH, 0, true
			if g.hasWings {
				g.jumpsLeft = 1
			}
		}
	}

	// The ability waits at the easy, low route on the right.
	if !g.hasWings && distance(g.x+playerW/2, g.y+playerH/2, 1135, 460) < 48 {
		g.hasWings, g.jumpsLeft = true, 1
		g.message = "DOUBLE JUMP GET! Return LEFT and jump again in mid-air."
		g.spawnBurst(1135, 460, 28)
	}
	// The relic is visible from the start but normal jump cannot reach it.
	if g.hasWings && distance(g.x+playerW/2, g.y+playerH/2, 70, 270) < 54 {
		g.won = true
		g.message = "RELIC FOUND! The old dead end became a route."
		g.spawnBurst(70, 270, 40)
	}

	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .05
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}
	target := g.x - W/2 + playerW/2
	g.camX += (target - g.camX) * .09
	g.camX = math.Max(0, math.Min(worldW-W, g.camX))
	return nil
}

func controls() (left, right, jump bool) {
	left = ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right = ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump = inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		left = left || (my > 620 && mx < 120)
		right = right || (my > 620 && mx >= 120 && mx < 250)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && my > 620 && mx >= 300 {
		jump = true
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > 620 && x < 120 {
			left = true
		} else if y > 620 && x < 250 {
			right = true
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > 620 && x >= 300 {
			jump = true
		}
	}
	return
}

func retryPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) spawnBurst(x, y float64, count int) {
	for i := 0; i < count; i++ {
		a := float64(i) * math.Pi * 2 / float64(count)
		g.particles = append(g.particles, particle{x, y, math.Cos(a) * (1.5 + float64(i%4)), math.Sin(a) * (1.5 + float64(i%4)), 34})
	}
}

func distance(ax, ay, bx, by float64) float64 { return math.Hypot(ax-bx, ay-by) }

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 17, 35, 255})
	// Parallax cave shapes make camera movement readable.
	for i := 0; i < 9; i++ {
		x := float32(float64(i*180)-g.camX*.25) - 60
		vector.DrawFilledCircle(s, x, 205+float32((i%3)*42), 115, color.RGBA{18, 38, 66, 255}, false)
	}
	for _, p := range platforms {
		x := p.x - g.camX
		if x+p.w < 0 || x > W {
			continue
		}
		vector.DrawFilledRect(s, float32(x), float32(p.y), float32(p.w), float32(p.h), color.RGBA{33, 68, 82, 255}, false)
		vector.DrawFilledRect(s, float32(x), float32(p.y), float32(p.w), 5, color.RGBA{57, 205, 181, 255}, false)
	}

	// Old route and its goal are always visible, so the player remembers it.
	relicX := float32(70 - g.camX)
	vector.DrawFilledCircle(s, relicX, 270, 16+float32(math.Sin(float64(g.tick)*.1)*3), color.RGBA{255, 201, 75, 255}, true)
	if !g.hasWings {
		vector.StrokeCircle(s, relicX, 270, 27, 3, color.RGBA{255, 201, 75, 150}, true)
	}
	orbX := float32(1135 - g.camX)
	if !g.hasWings {
		vector.DrawFilledCircle(s, orbX, 460, 22, color.RGBA{153, 119, 255, 255}, true)
		vector.StrokeCircle(s, orbX, 460, 34, 4, color.RGBA{88, 224, 255, 190}, true)
	}
	for _, p := range g.particles {
		vector.DrawFilledCircle(s, float32(p.x-g.camX), float32(p.y), float32(2+p.life/12), color.RGBA{255, 211, 90, uint8(math.Min(255, p.life*7))}, true)
	}

	px := float32(g.x - g.camX)
	vector.DrawFilledRect(s, px, float32(g.y), playerW, playerH, color.RGBA{255, 113, 125, 255}, true)
	vector.DrawFilledCircle(s, px+8, float32(g.y)+8, 3, color.White, true)
	vector.DrawFilledCircle(s, px+22, float32(g.y)+8, 3, color.White, true)
	if g.hasWings {
		vector.DrawFilledCircle(s, px-4, float32(g.y)+22, 10, color.RGBA{146, 222, 255, 210}, true)
		vector.DrawFilledCircle(s, px+34, float32(g.y)+22, 10, color.RGBA{146, 222, 255, 210}, true)
	}

	vector.DrawFilledRect(s, 0, 0, W, 116, color.RGBA{7, 14, 31, 235}, true)
	ebitenutil.DebugPrintAt(s, "ABILITY ROUTE / RETURN TO THE OLD LEDGE", 72, 18)
	ebitenutil.DebugPrintAt(s, g.message, 24, 48)
	ability := "LOCKED"
	if g.hasWings {
		ability = "READY (1 AIR JUMP)"
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DOUBLE JUMP: %s", ability), 24, 78)
	progress := g.x / worldW
	vector.DrawFilledRect(s, 24, 98, 432, 6, color.RGBA{42, 55, 86, 255}, true)
	vector.DrawFilledRect(s, 24, 98, float32(432*progress), 6, color.RGBA{57, 205, 181, 255}, true)

	drawButton(s, 16, 630, 104, 64, "LEFT")
	drawButton(s, 128, 630, 104, 64, "RIGHT")
	drawButton(s, 312, 630, 152, 64, "JUMP")
	if g.won {
		vector.DrawFilledRect(s, 34, 210, 412, 190, color.RGBA{12, 22, 48, 245}, true)
		vector.StrokeRect(s, 34, 210, 412, 190, 4, color.RGBA{255, 201, 75, 255}, true)
		ebitenutil.DebugPrintAt(s, "RELIC FOUND!", 184, 254)
		ebitenutil.DebugPrintAt(s, "The new ability changed an old dead end.", 82, 300)
		ebitenutil.DebugPrintAt(s, "TAP / R: RETRY", 176, 350)
	}
}

func drawButton(s *ebiten.Image, x, y, w, h float32, label string) {
	vector.DrawFilledRect(s, x, y, w, h, color.RGBA{28, 48, 78, 245}, true)
	vector.StrokeRect(s, x, y, w, h, 2, color.RGBA{73, 151, 174, 255}, true)
	ebitenutil.DebugPrintAt(s, label, int(x+w/2)-len(label)*3, int(y+28))
}

func (g *game) Layout(_, _ int) (int, int) { return W, H }

func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ability Route")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
