package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
)

const (
	screenWidth  = 480
	screenHeight = 720
	basketY      = 640
	basketHalfW  = 58
)

// 落ちてくる星 1つ
type star struct {
	x, y  float64
	speed float64
}

type game struct {
	basketX  float64
	stars    []star
	frame    int
	score    int
	lives    int
	gameOver bool
	rng      *rand.Rand
}

func newGame() *game {
	return &game{
		basketX: screenWidth / 2,
		lives:   3,
		rng:     rand.New(rand.NewSource(19)),
	}
}

func restartPressed() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(ebiten.AppendTouchIDs(nil)) > 0 {
		return true
	}
	return false
}

// キーボード・マウス・タッチでカゴを動かす
func (g *game) moveBasket() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.basketX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.basketX += 5
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		g.basketX = float64(x)
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		g.basketX = float64(x)
	}

	// 画面の端から出ないようにする
	if g.basketX < 55 {
		g.basketX = 55
	}
	if g.basketX > screenWidth-55 {
		g.basketX = screenWidth - 55
	}
}

// 星がカゴに入ったか（AABB：四角どうしの重なり）
func (g *game) caught(s star) bool {
	inY := s.y > 625 && s.y < 675
	inX := s.x > g.basketX-basketHalfW && s.x < g.basketX+basketHalfW
	return inY && inX
}

// --- ここから Update ---
func (g *game) Update() error {
	if g.gameOver {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}

	g.moveBasket()

	g.frame++
	// 45 tickごとに星を1つ出す
	if g.frame%45 == 0 {
		g.stars = append(g.stars, star{
			x:     25 + g.rng.Float64()*(screenWidth-50),
			y:     -20,
			speed: 2.6 + g.rng.Float64()*2,
		})
	}

	// 残す星だけを next に集める
	next := g.stars[:0]
	for _, s := range g.stars {
		s.y += s.speed

		if g.caught(s) {
			g.score++
			continue // 取れたので消す
		}
		if s.y > screenHeight {
			g.lives--
			if g.lives <= 0 {
				g.gameOver = true
			}
			continue // 落としたので消す
		}
		next = append(next, s)
	}
	g.stars = next
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})

	// うすい星空
	for i := 0; i < 24; i++ {
		vector.DrawFilledCircle(
			screen,
			float32((i*97)%screenWidth),
			float32((i*53+g.frame/4)%screenHeight),
			1.5,
			color.RGBA{45, 226, 194, 120},
			false,
		)
	}

	for _, s := range g.stars {
		vector.DrawFilledCircle(screen, float32(s.x), float32(s.y), 15, color.RGBA{255, 210, 72, 255}, false)
		vector.DrawFilledCircle(screen, float32(s.x-5), float32(s.y-5), 4, color.White, false)
	}

	vector.DrawFilledRect(screen, float32(g.basketX-basketHalfW), basketY, 116, 38, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledRect(screen, float32(g.basketX-48), 630, 96, 12, color.RGBA{255, 105, 79, 255}, false)
	// カゴを持つ主人公
	hero.DrawBottomCentered(screen, g.basketX, basketY+8, 56)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %02d   LIVES %d", g.score, g.lives), 170, 25)
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER\n\nCLICK / TOUCH TO RETRY", 155, 320)
	} else {
		ebitenutil.DebugPrintAt(screen, "MOVE: ARROWS / DRAG", 170, 690)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Catch the Stars — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
