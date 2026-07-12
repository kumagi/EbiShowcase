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

// --- 画面の大きさ（ゲーム内のマス目） ---
const (
	screenWidth  = 480
	screenHeight = 720
	startSeconds = 30
)

// game はこのゲームの全部の数字を入れた箱です。
type game struct {
	circleX    float64 // 丸の中心 X
	circleY    float64 // 丸の中心 Y
	radius     float64 // 丸の半径
	score      int     // 得点
	framesLeft int     // 残りフレーム（60フレーム ≒ 1秒）
	started    bool    // スタートしたか
	rng        *rand.Rand
}

func newGame() *game {
	g := &game{
		radius:     38,
		framesLeft: startSeconds * 60,
		rng:        rand.New(rand.NewSource(11)),
	}
	g.moveTarget()
	return g
}

// 丸を別の場所へ移す
func (g *game) moveTarget() {
	g.circleX = 60 + g.rng.Float64()*(screenWidth-120)
	g.circleY = 130 + g.rng.Float64()*(screenHeight-240)
}

// マウスまたはタッチで「いま押した位置」を返す
func pressedPosition() (x, y int, pressed bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y = ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y = ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}

// --- ここから Update（数字を進める） ---
func (g *game) Update() error {
	px, py, pressed := pressedPosition()

	// まだ始まっていない → 押したらスタート
	if !g.started {
		if pressed {
			g.started = true
		}
		return nil
	}

	// 時間切れ → 押したら最初から
	if g.framesLeft <= 0 {
		if pressed {
			g.score = 0
			g.framesLeft = startSeconds * 60
			g.started = false
			g.moveTarget()
		}
		return nil
	}

	// タイマーを1フレーム減らす
	g.framesLeft--

	// 押した瞬間だけ、丸との距離を調べる
	if pressed {
		dx := float64(px) - g.circleX
		dy := float64(py) - g.circleY
		distance := math.Hypot(dx, dy)
		if distance <= g.radius {
			g.score++
			// 当たるほど丸を少し小さくする
			g.radius = math.Max(22, 38-float64(g.score)/2)
			g.moveTarget()
		}
	}
	return nil
}

// --- ここから Draw（画面に描く） ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})

	// 背景のうすい円
	for i := 0; i < 12; i++ {
		vector.StrokeCircle(
			screen,
			float32(screenWidth/2),
			float32(screenHeight/2),
			float32(55+i*38),
			1,
			color.RGBA{45, 226, 194, 25},
			false,
		)
	}

	// 光る丸（外側のぼかし → 本体 → ハイライト）
	vector.DrawFilledCircle(screen, float32(g.circleX), float32(g.circleY), float32(g.radius+8), color.RGBA{45, 226, 194, 45}, false)
	vector.DrawFilledCircle(screen, float32(g.circleX), float32(g.circleY), float32(g.radius), color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledCircle(screen, float32(g.circleX-10), float32(g.circleY-10), float32(g.radius/4), color.RGBA{230, 255, 249, 220}, false)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %02d", g.score), 24, 24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME  %02d", max(0, g.framesLeft/60)), 380, 24)

	if !g.started {
		ebitenutil.DebugPrintAt(screen, "TAP THE TARGET\n\nCLICK / TOUCH TO START", 145, 320)
	} else if g.framesLeft <= 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME UP!  SCORE %d\n\nTAP TO RETRY", g.score), 160, 320)
	}
}

// --- Layout（画面の大きさ） ---
func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tap the Target — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
