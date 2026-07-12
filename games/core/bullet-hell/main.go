package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
)

const (
	screenWidth  = 480
	screenHeight = 720
	bossX        = 240.0
	bossY        = 128.0
	bossMaxHP    = 320
)

type bullet struct {
	x, y   float64
	vx, vy float64
	r      float64
	enemy  bool // true = 敵弾
}

type game struct {
	playerX    float64
	playerY    float64
	vx, vy     float64
	bullets    []bullet
	frame      int
	bossHP     int
	life       int
	invincible int
	cleared    bool
	gameOver   bool
}

func newGame() *game {
	return &game{
		playerX: 240,
		playerY: 620,
		bossHP:  bossMaxHP,
		life:    5,
	}
}

func clamp(v, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, v))
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

func overlay(screen *ebiten.Image, msg string) {
	vector.DrawFilledRect(screen, 55, 280, 370, 150, color.RGBA{6, 12, 35, 240}, false)
	ebitenutil.DebugPrintAt(screen, msg, 145, 330)
}

func (g *game) movePlayer() {
	ax, ay := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		ax--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		ax++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		ay--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		ay++
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		ax = float64(x) - g.playerX
		ay = float64(y) - g.playerY
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		ax = float64(x) - g.playerX
		ay = float64(y) - g.playerY
	}

	if ax != 0 || ay != 0 {
		length := math.Hypot(ax, ay)
		g.vx += ax / length * 0.8
		g.vy += ay / length * 0.8
	}

	g.vx *= 0.8
	g.vy *= 0.8
	speed := math.Hypot(g.vx, g.vy)
	if speed > 5.2 {
		g.vx *= 5.2 / speed
		g.vy *= 5.2 / speed
	}

	g.playerX = clamp(g.playerX+g.vx, 18, 462)
	g.playerY = clamp(g.playerY+g.vy, 90, 690)
}

// 円形弾幕
func (g *game) spawnRing() {
	count := 18
	spin := float64(g.frame) * 0.017
	speed := 2.0 + float64((g.frame/38)%3)*0.35
	for i := 0; i < count; i++ {
		a := spin + float64(i)*math.Pi*2/float64(count)
		g.bullets = append(g.bullets, bullet{
			x: bossX, y: bossY,
			vx: math.Cos(a) * speed,
			vy: math.Sin(a) * speed,
			r:  6, enemy: true,
		})
	}
}

// 自分めがけた扇状弾
func (g *game) spawnAimFan() {
	base := math.Atan2(g.playerY-bossY, g.playerX-bossX)
	for i := -3; i <= 3; i++ {
		a := base + float64(i)*0.11
		g.bullets = append(g.bullets, bullet{
			x: bossX, y: bossY,
			vx: math.Cos(a) * 3.0,
			vy: math.Sin(a) * 3.0,
			r:  7, enemy: true,
		})
	}
}

// --- ここから Update ---
func (g *game) Update() error {
	if g.cleared || g.gameOver {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}

	g.frame++
	if g.invincible > 0 {
		g.invincible--
	}

	g.movePlayer()

	// 自動で自分の弾を撃つ
	if g.frame%7 == 0 {
		g.bullets = append(g.bullets,
			bullet{g.playerX - 7, g.playerY - 15, -0.25, -10, 4, false},
			bullet{g.playerX + 7, g.playerY - 15, 0.25, -10, 4, false},
		)
	}

	// ボスの弾幕パターン
	if g.frame%38 == 0 {
		g.spawnRing()
	}
	if g.frame%91 == 0 {
		g.spawnAimFan()
	}

	// 弾を動かして当たりを見る
	for i := len(g.bullets) - 1; i >= 0; i-- {
		b := &g.bullets[i]
		b.x += b.vx
		b.y += b.vy

		remove := b.x < -30 || b.x > 510 || b.y < -30 || b.y > 750

		// 自分の弾がボスに当たった
		if !b.enemy && math.Hypot(b.x-bossX, b.y-bossY) < 48 {
			g.bossHP--
			remove = true
			if g.bossHP <= 0 {
				g.cleared = true
			}
		}

		// 敵弾が自分に当たった
		if b.enemy && g.invincible == 0 && math.Hypot(b.x-g.playerX, b.y-g.playerY) < b.r+4 {
			g.life--
			g.invincible = 90
			remove = true
			if g.life <= 0 {
				g.gameOver = true
			}
		}

		if remove {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
		}
	}
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{5, 8, 24, 255})

	for i := 0; i < 36; i++ {
		x := float32((i*83 + g.frame/3) % screenWidth)
		y := float32((i*137 + g.frame) % screenHeight)
		vector.DrawFilledCircle(screen, x, y, 1, color.RGBA{150, 170, 230, 170}, false)
	}

	vector.DrawFilledCircle(screen, bossX, bossY, 42, color.RGBA{143, 55, 208, 255}, false)
	vector.StrokeCircle(screen, bossX, bossY, 52, 4, color.RGBA{246, 80, 133, 220}, false)
	vector.DrawFilledCircle(screen, 226, 122, 5, color.RGBA{255, 224, 90, 255}, false)
	vector.DrawFilledCircle(screen, 254, 122, 5, color.RGBA{255, 224, 90, 255}, false)

	// ボス HP バー
	vector.DrawFilledRect(screen, 55, 55, 370, 14, color.RGBA{44, 48, 72, 255}, false)
	hpW := float32(370 * max(g.bossHP, 0) / bossMaxHP)
	vector.DrawFilledRect(screen, 55, 55, hpW, 14, color.RGBA{247, 72, 113, 255}, false)

	for _, b := range g.bullets {
		c := color.RGBA{255, 218, 72, 255}
		if b.enemy {
			c = color.RGBA{255, 76, 145, 255}
			vector.StrokeCircle(screen, float32(b.x), float32(b.y), float32(b.r+2), 1, color.RGBA{139, 105, 255, 220}, false)
		}
		vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), float32(b.r), c, false)
	}

	if g.invincible%10 < 5 {
		hero.DrawCentered(screen, g.playerX, g.playerY, 38)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BOSS %03d     LIFE %d     BULLETS %03d", max(g.bossHP, 0), g.life, len(g.bullets)), 70, 20)
	ebitenutil.DebugPrintAt(screen, "MOVE: WASD / ARROWS / DRAG    AUTO FIRE", 96, 688)
	if g.cleared {
		overlay(screen, "BOSS DEFEATED!\n\nTAP / SPACE TO FIGHT AGAIN")
	} else if g.gameOver {
		overlay(screen, "CONTINUE?\n\nTAP / SPACE TO RETRY")
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bullet Hell Boss — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
