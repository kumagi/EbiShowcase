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
	"github.com/kumagi/EbiShowcase/internal/hero"
)

const (
	screenWidth  = 480
	screenHeight = 720
)

// 丸い当たり判定つきの物体
type body struct {
	x, y   float64
	vx, vy float64
	radius float64
}

type enemy struct {
	body
	phase  float64
	reload float64 // 次に撃つまでの待ち
}

type star struct {
	x, y  float64
	speed float64
}

type game struct {
	player       body
	bullets      []body
	enemyBullets []body
	enemies      []enemy
	stars        []star
	rng          *rand.Rand
	frame        int
	score        int
	lives        int
	wave         int
	shootWait    int
	invincible   int // 無敵フレーム残り
	gameOver     bool
}

func newGame() *game {
	g := &game{
		player: body{x: screenWidth / 2, y: screenHeight - 92, radius: 15},
		rng:    rand.New(rand.NewSource(808)),
		lives:  3,
	}
	for i := 0; i < 70; i++ {
		g.stars = append(g.stars, star{
			x:     g.rng.Float64() * screenWidth,
			y:     g.rng.Float64() * screenHeight,
			speed: 0.5 + g.rng.Float64()*2,
		})
	}
	g.spawnWave()
	return g
}

func clamp(v, low, high float64) float64 {
	return math.Max(low, math.Min(high, v))
}

func overlaps(a, b body) bool {
	return math.Hypot(a.x-b.x, a.y-b.y) < a.radius+b.radius
}

// aimVelocity returns a vector with the requested length.  When the target is
// exactly on the shooter, there is no direction to normalize, so return zero
// instead of dividing by zero and producing NaN.
func aimVelocity(dx, dy, speed float64) (float64, float64) {
	length := math.Hypot(dx, dy)
	if length == 0 {
		return 0, 0
	}
	return dx / length * speed, dy / length * speed
}

func restartPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		return true
	}
	return false
}

func pointerPosition() (x, y float64, ok bool) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		px, py := ebiten.CursorPosition()
		return float64(px), float64(py), true
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		px, py := ebiten.TouchPosition(ids[0])
		return float64(px), float64(py), true
	}
	return 0, 0, false
}

func (g *game) spawnWave() {
	g.wave++
	count := min(4+g.wave, 10)
	for i := 0; i < count; i++ {
		g.enemies = append(g.enemies, enemy{
			body: body{
				x:      float64(54 + i%5*93),
				y:      float64(-40 - i/5*70),
				vy:     0.65 + float64(g.wave)*0.05,
				radius: 17,
			},
			phase:  float64(i) * 0.8,
			reload: float64(80 + g.rng.Intn(100)),
		})
	}
}

func (g *game) hurt() {
	if g.invincible > 0 || g.gameOver {
		return
	}
	g.lives--
	g.invincible = 120
	if g.lives <= 0 {
		g.gameOver = true
	}
}

func (g *game) movePlayer() {
	const accel = 0.7
	const drag = 0.82
	const maxSpeed = 7.0

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.vx -= accel
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.vx += accel
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.vy -= accel
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.vy += accel
	}

	targetX, targetY, pointer := pointerPosition()
	if pointer {
		g.player.vx += clamp((targetX-g.player.x)*0.08, -accel*1.8, accel*1.8)
		g.player.vy += clamp((targetY-g.player.y)*0.08, -accel*1.8, accel*1.8)
	}

	g.player.vx = clamp(g.player.vx*drag, -maxSpeed, maxSpeed)
	g.player.vy = clamp(g.player.vy*drag, -maxSpeed, maxSpeed)
	g.player.x = clamp(g.player.x+g.player.vx, 22, screenWidth-22)
	g.player.y = clamp(g.player.y+g.player.vy, 90, screenHeight-52)
}

func (g *game) tryShoot() {
	if g.shootWait > 0 {
		g.shootWait--
	}
	shooting := ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) ||
		len(ebiten.AppendTouchIDs(nil)) > 0
	if shooting && g.shootWait == 0 {
		g.bullets = append(g.bullets, body{
			x: g.player.x, y: g.player.y - 20,
			vy: -10, radius: 5,
		})
		g.shootWait = 9
	}
}

// --- ここから Update ---
func (g *game) Update() error {
	if g.gameOver {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}

	g.frame++
	if g.invincible > 0 {
		g.invincible--
	}

	// 背景の星を流す
	for i := range g.stars {
		g.stars[i].y += g.stars[i].speed
		if g.stars[i].y > screenHeight {
			g.stars[i].y = 0
			g.stars[i].x = g.rng.Float64() * screenWidth
		}
	}

	g.movePlayer()
	g.tryShoot()

	// 弾を動かす
	for i := range g.bullets {
		g.bullets[i].x += g.bullets[i].vx
		g.bullets[i].y += g.bullets[i].vy
	}
	for i := range g.enemyBullets {
		g.enemyBullets[i].x += g.enemyBullets[i].vx
		g.enemyBullets[i].y += g.enemyBullets[i].vy
	}

	// 敵の移動と射撃
	for i := range g.enemies {
		e := &g.enemies[i]
		e.y += e.vy
		e.x += math.Sin(float64(g.frame)*0.025+e.phase) * 1.2
		e.reload--
		if e.reload <= 0 && e.y > 60 {
			dx := g.player.x - e.x
			dy := g.player.y - e.y
			vx, vy := aimVelocity(dx, dy, 2.3)
			g.enemyBullets = append(g.enemyBullets, body{
				x: e.x, y: e.y + 15,
				vx:     vx,
				vy:     vy,
				radius: 7,
			})
			e.reload = float64(max(42, 125-g.wave*5) + g.rng.Intn(70))
		}
	}

	// 自分の弾 × 敵
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

	// 敵の弾 × 自分
	for i := len(g.enemyBullets) - 1; i >= 0; i-- {
		b := g.enemyBullets[i]
		offscreen := b.y > screenHeight+20 || b.x < -20 || b.x > screenWidth+20
		if offscreen {
			g.enemyBullets = append(g.enemyBullets[:i], g.enemyBullets[i+1:]...)
			continue
		}
		if g.invincible == 0 && overlaps(g.player, b) {
			g.enemyBullets = append(g.enemyBullets[:i], g.enemyBullets[i+1:]...)
			g.hurt()
		}
	}

	// 敵が下まで来たらダメージ
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

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{4, 11, 29, 255})

	for _, s := range g.stars {
		vector.DrawFilledCircle(screen, float32(s.x), float32(s.y), float32(0.8+s.speed*0.35), color.RGBA{160, 210, 255, 210}, false)
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

	// 無敵中は点滅
	if g.invincible%12 < 6 {
		hero.DrawCentered(screen, g.player.x, g.player.y, 40)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %06d", g.score), 18, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("WAVE %02d", g.wave), 205, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("LIFE %d", g.lives), 390, 18)
	if g.gameOver {
		vector.DrawFilledRect(screen, 55, 280, 370, 150, color.RGBA{5, 14, 34, 235}, false)
		ebitenutil.DebugPrintAt(screen, "MISSION FAILED\n\nTAP / SPACE TO RETRY", 155, 330)
	} else {
		ebitenutil.DebugPrintAt(screen, "MOVE: ARROWS / DRAG    FIRE: SPACE / TOUCH", 72, 690)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Space Shooter — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
