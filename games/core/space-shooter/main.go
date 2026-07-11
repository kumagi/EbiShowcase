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
)

const screenWidth, screenHeight = 480, 720

type body struct{ x, y, vx, vy, radius float64 }
type enemy struct {
	body
	phase, reload float64
}
type star struct{ x, y, speed float64 }
type game struct {
	player                                               body
	bullets, enemyBullets                                []body
	enemies                                              []enemy
	stars                                                []star
	rng                                                  *rand.Rand
	frame, score, lives, wave, shootCooldown, invincible int
	over                                                 bool
}

func newGame() *game {
	g := &game{player: body{x: screenWidth / 2, y: screenHeight - 92, radius: 15}, rng: rand.New(rand.NewSource(808)), lives: 3}
	for i := 0; i < 70; i++ {
		g.stars = append(g.stars, star{g.rng.Float64() * screenWidth, g.rng.Float64() * screenHeight, .5 + g.rng.Float64()*2})
	}
	g.spawnWave()
	return g
}

func (g *game) spawnWave() {
	g.wave++
	count := min(4+g.wave, 10)
	for i := 0; i < count; i++ {
		g.enemies = append(g.enemies, enemy{body: body{x: float64(54 + i%5*93), y: float64(-40 - i/5*70), vy: .65 + float64(g.wave)*.05, radius: 17}, phase: float64(i) * .8, reload: float64(80 + g.rng.Intn(100))})
	}
}

func (g *game) Update() error {
	if g.over {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frame++
	if g.invincible > 0 {
		g.invincible--
	}
	for i := range g.stars {
		g.stars[i].y += g.stars[i].speed
		if g.stars[i].y > screenHeight {
			g.stars[i].y = 0
			g.stars[i].x = g.rng.Float64() * screenWidth
		}
	}
	targetX, targetY, pointer := pointerPosition()
	const acceleration, drag, maxSpeed = .7, .82, 7.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.vx -= acceleration
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.vx += acceleration
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.vy -= acceleration
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.vy += acceleration
	}
	if pointer {
		g.player.vx += clamp((targetX-g.player.x)*.08, -acceleration*1.8, acceleration*1.8)
		g.player.vy += clamp((targetY-g.player.y)*.08, -acceleration*1.8, acceleration*1.8)
	}
	g.player.vx = clamp(g.player.vx*drag, -maxSpeed, maxSpeed)
	g.player.vy = clamp(g.player.vy*drag, -maxSpeed, maxSpeed)
	g.player.x = clamp(g.player.x+g.player.vx, 22, screenWidth-22)
	g.player.y = clamp(g.player.y+g.player.vy, 90, screenHeight-52)
	if g.shootCooldown > 0 {
		g.shootCooldown--
	}
	shooting := ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || len(ebiten.AppendTouchIDs(nil)) > 0
	if shooting && g.shootCooldown == 0 {
		g.bullets = append(g.bullets, body{x: g.player.x, y: g.player.y - 20, vy: -10, radius: 5})
		g.shootCooldown = 9
	}
	for i := range g.bullets {
		g.bullets[i].x += g.bullets[i].vx
		g.bullets[i].y += g.bullets[i].vy
	}
	for i := range g.enemyBullets {
		g.enemyBullets[i].x += g.enemyBullets[i].vx
		g.enemyBullets[i].y += g.enemyBullets[i].vy
	}
	for i := range g.enemies {
		e := &g.enemies[i]
		e.y += e.vy
		e.x += math.Sin(float64(g.frame)*.025+e.phase) * 1.2
		e.reload--
		if e.reload <= 0 && e.y > 60 {
			dx, dy := g.player.x-e.x, g.player.y-e.y
			length := math.Hypot(dx, dy)
			g.enemyBullets = append(g.enemyBullets, body{x: e.x, y: e.y + 15, vx: dx / length * 2.3, vy: dy / length * 2.3, radius: 7})
			e.reload = float64(max(42, 125-g.wave*5) + g.rng.Intn(70))
		}
	}
	for bi := len(g.bullets) - 1; bi >= 0; bi-- {
		hit := false
		for ei := len(g.enemies) - 1; ei >= 0; ei-- {
			if overlaps(g.bullets[bi], g.enemies[ei].body) {
				g.enemies = append(g.enemies[:ei], g.enemies[ei+1:]...)
				g.score += 100
				hit = true
				break
			}
		}
		if hit || g.bullets[bi].y < -20 {
			g.bullets = append(g.bullets[:bi], g.bullets[bi+1:]...)
		}
	}
	for i := len(g.enemyBullets) - 1; i >= 0; i-- {
		if g.enemyBullets[i].y > screenHeight+20 || g.enemyBullets[i].x < -20 || g.enemyBullets[i].x > screenWidth+20 {
			g.enemyBullets = append(g.enemyBullets[:i], g.enemyBullets[i+1:]...)
			continue
		}
		if g.invincible == 0 && overlaps(g.player, g.enemyBullets[i]) {
			g.enemyBullets = append(g.enemyBullets[:i], g.enemyBullets[i+1:]...)
			g.hurt()
		}
	}
	for i := len(g.enemies) - 1; i >= 0; i-- {
		if g.enemies[i].y > screenHeight+30 {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			g.hurt()
		}
	}
	if len(g.enemies) == 0 {
		g.spawnWave()
	}
	return nil
}

func (g *game) hurt() {
	if g.invincible > 0 || g.over {
		return
	}
	g.lives--
	g.invincible = 120
	if g.lives <= 0 {
		g.over = true
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{4, 11, 29, 255})
	for _, s := range g.stars {
		vector.DrawFilledCircle(screen, float32(s.x), float32(s.y), float32(.8+s.speed*.35), color.RGBA{160, 210, 255, 210}, false)
	}
	for _, b := range g.bullets {
		vector.DrawFilledRect(screen, float32(b.x-2), float32(b.y-10), 4, 18, color.RGBA{255, 221, 86, 255}, false)
	}
	for _, b := range g.enemyBullets {
		vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), float32(b.radius), color.RGBA{255, 94, 122, 255}, false)
		vector.StrokeCircle(screen, float32(b.x), float32(b.y), float32(b.radius+3), 1, color.RGBA{255, 170, 190, 150}, false)
	}
	for _, e := range g.enemies {
		vector.DrawFilledCircle(screen, float32(e.x), float32(e.y), 17, color.RGBA{255, 91, 109, 255}, false)
		vector.DrawFilledRect(screen, float32(e.x-24), float32(e.y-4), 48, 9, color.RGBA{130, 54, 210, 255}, false)
		vector.DrawFilledCircle(screen, float32(e.x), float32(e.y+3), 5, color.RGBA{255, 222, 91, 255}, false)
	}
	if g.invincible%12 < 6 {
		x, y := float32(g.player.x), float32(g.player.y)
		vector.DrawFilledRect(screen, x-6, y-20, 12, 38, color.RGBA{44, 226, 194, 255}, false)
		vector.DrawFilledRect(screen, x-22, y+4, 44, 10, color.RGBA{44, 226, 194, 255}, false)
		vector.DrawFilledCircle(screen, x, y-8, 6, color.RGBA{255, 223, 83, 255}, false)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %06d", g.score), 18, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("WAVE %02d", g.wave), 205, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("LIFE %d", g.lives), 390, 18)
	if g.over {
		vector.DrawFilledRect(screen, 55, 280, 370, 150, color.RGBA{5, 14, 34, 235}, false)
		ebitenutil.DebugPrintAt(screen, "MISSION FAILED\n\nTAP / SPACE TO RETRY", 155, 330)
	} else {
		ebitenutil.DebugPrintAt(screen, "MOVE: ARROWS / DRAG    FIRE: SPACE / TOUCH", 72, 690)
	}
}

func pointerPosition() (float64, float64, bool) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return float64(x), float64(y), true
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return float64(x), float64(y), true
	}
	return 0, 0, false
}
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlaps(a, b body) bool              { return math.Hypot(a.x-b.x, a.y-b.y) < a.radius+b.radius }
func clamp(v, low, high float64) float64   { return math.Max(low, math.Min(high, v)) }
func (g *game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }
func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Space Shooter — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
