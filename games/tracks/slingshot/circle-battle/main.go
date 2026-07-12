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
	width, height = 480, 720
	playerRadius  = 20.0
	enemyRadius   = 30.0
	timeLimit     = 45 * 60
)

type enemy struct {
	x, y     float64
	hp       int
	cooldown int
}

type game struct {
	x, y, vx, vy float64
	enemies      []enemy
	frames       int
	hits         int
	clear, over  bool
	message      string
}

func newGame() *game {
	return &game{
		x: 240, y: 560,
		enemies: []enemy{{110, 190, 3, 0}, {365, 270, 3, 0}, {235, 390, 3, 0}},
		message: "Build speed, then bump every coral circle!",
	}
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.frames >= timeLimit {
		g.over = true
		g.message = "Time up — try a faster route!"
		return nil
	}

	ax, ay := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		ax--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		ax++
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		ay--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		ay++
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		px, py := ebiten.CursorPosition()
		ax, ay = float64(px)-g.x, float64(py)-g.y
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		px, py := ebiten.TouchPosition(ids[0])
		ax, ay = float64(px)-g.x, float64(py)-g.y
	}
	if length := math.Hypot(ax, ay); length > 0 {
		g.vx += ax / length * 0.34
		g.vy += ay / length * 0.34
	}
	speed := math.Hypot(g.vx, g.vy)
	if speed > 9 {
		g.vx, g.vy = g.vx*9/speed, g.vy*9/speed
	}
	g.vx *= 0.985
	g.vy *= 0.985
	g.x += g.vx
	g.y += g.vy
	g.wallBounce()

	alive := 0
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp <= 0 {
			continue
		}
		alive++
		if e.cooldown > 0 {
			e.cooldown--
		}
		g.collide(e)
	}
	if alive == 0 {
		g.clear = true
		g.message = "All coral guards cleared!"
	}
	return nil
}

func (g *game) wallBounce() {
	if g.x < playerRadius {
		g.x, g.vx = playerRadius, math.Abs(g.vx)*0.8
	}
	if g.x > width-playerRadius {
		g.x, g.vx = width-playerRadius, -math.Abs(g.vx)*0.8
	}
	if g.y < 82+playerRadius {
		g.y, g.vy = 82+playerRadius, math.Abs(g.vy)*0.8
	}
	if g.y > 635-playerRadius {
		g.y, g.vy = 635-playerRadius, -math.Abs(g.vy)*0.8
	}
}

func (g *game) collide(e *enemy) {
	dx, dy := g.x-e.x, g.y-e.y
	distance := math.Hypot(dx, dy)
	minDistance := playerRadius + enemyRadius
	if distance >= minDistance {
		return
	}
	if distance == 0 {
		dx, dy, distance = 1, 0, 1
	}
	nx, ny := dx/distance, dy/distance
	// Push the player exactly to the edge so circles cannot remain overlapped.
	g.x, g.y = e.x+nx*minDistance, e.y+ny*minDistance
	dot := g.vx*nx + g.vy*ny
	impactSpeed := math.Abs(dot)
	if dot < 0 {
		g.vx -= 1.7 * dot * nx
		g.vy -= 1.7 * dot * ny
	}
	if e.cooldown == 0 && impactSpeed >= 1.5 {
		e.hp--
		e.cooldown = 24
		g.hits++
		g.message = fmt.Sprintf("Impact %.1f: damage! Cooldown started.", impactSpeed)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 22, 41, 255})
	ebitenutil.DebugPrintAt(screen, "CIRCLE CONTACT ARENA", 160, 22)
	seconds := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME %02d   HITS %d", seconds, g.hits), 188, 50)
	vector.StrokeRect(screen, 8, 82, 464, 553, 4, color.RGBA{48, 76, 111, 255}, false)
	for _, e := range g.enemies {
		if e.hp <= 0 {
			continue
		}
		c := color.RGBA{225, 92, 91, 255}
		if e.cooldown > 0 {
			c = color.RGBA{245, 177, 66, 255}
		}
		vector.DrawFilledCircle(screen, float32(e.x), float32(e.y), enemyRadius, c, true)
		vector.StrokeCircle(screen, float32(e.x), float32(e.y), enemyRadius, 3, color.White, true)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %d", e.hp), int(e.x)-14, int(e.y)-4)
	}
	vector.DrawFilledCircle(screen, float32(g.x), float32(g.y), playerRadius, color.RGBA{89, 203, 224, 255}, true)
	vector.StrokeCircle(screen, float32(g.x), float32(g.y), playerRadius, 3, color.White, true)
	vector.StrokeLine(screen, float32(g.x), float32(g.y), float32(g.x+g.vx*7), float32(g.y+g.vy*7), 4, color.RGBA{250, 210, 81, 255}, true)
	ebitenutil.DebugPrintAt(screen, g.message, 68, 652)
	ebitenutil.DebugPrintAt(screen, "HOLD/TAP A DIRECTION OR USE ARROWS / WASD", 72, 683)
	if g.clear || g.over {
		title := "CONTACT BATTLE CLEAR!"
		if g.over {
			title = "TIME UP!"
		}
		vector.DrawFilledRect(screen, 40, 276, 400, 158, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 148, 322)
		ebitenutil.DebugPrintAt(screen, g.message, 102, 355)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 397)
	}
}

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Circle Contact Arena — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
