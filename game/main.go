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

// --- 授業で触る数字 ---
const (
	screenWidth  = 480
	screenHeight = 720
	groundY      = 650
	gravity      = 0.42 // 毎フレーム、速さに足す量（加速度）
	flapSpeed    = -7.4 // はばたいたときの上向きの速さ
	birdX        = 128.0
	birdRadius   = 15.0
	pipeWidth    = 76.0
	gapHeight    = 178.0
)

// パイプ 1本（すきまの中心 Y を持つ）
type pipe struct {
	x      float64
	gapY   float64
	scored bool // もう得点したか
}

type game struct {
	birdY    float64
	velocity float64 // 速さ（マイナス＝上、プラス＝下）
	pipes    []pipe
	score    int
	best     int
	frame    int
	started  bool
	gameOver bool
	rng      *rand.Rand
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(7))}
	g.reset()
	return g
}

func (g *game) reset() {
	g.birdY = 320
	g.velocity = 0
	g.pipes = []pipe{
		{x: 560, gapY: 300},
		{x: 830, gapY: 390},
		{x: 1100, gapY: 265},
	}
	g.score = 0
	g.frame = 0
	g.started = false
	g.gameOver = false
}

func justPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
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

func (g *game) hitPipe() bool {
	for _, p := range g.pipes {
		inX := birdX+birdRadius > p.x && birdX-birdRadius < p.x+pipeWidth
		if !inX {
			continue
		}
		gapTop := p.gapY - gapHeight/2
		gapBottom := p.gapY + gapHeight/2
		hitTop := g.birdY-birdRadius < gapTop
		hitBottom := g.birdY+birdRadius > gapBottom
		if hitTop || hitBottom {
			return true
		}
	}
	return false
}

// --- ここから Update（数字を進める） ---
func (g *game) Update() error {
	g.frame++
	pressed := justPressed()

	if g.gameOver {
		if pressed {
			g.reset()
		}
		return nil
	}

	// スタート前は少し上下にゆらす
	if !g.started {
		g.birdY = 320 + math.Sin(float64(g.frame)*0.06)*7
		if pressed {
			g.started = true
			g.velocity = flapSpeed
		}
		return nil
	}

	// 押した瞬間、上向きの速さにする
	if pressed {
		g.velocity = flapSpeed
	}

	// 1. 速さに重力を足す（加速度）
	g.velocity += gravity
	// 2. 速さで位置を動かす
	g.birdY += g.velocity

	// パイプを左へ流す
	for i := range g.pipes {
		g.pipes[i].x -= 2.8
		if !g.pipes[i].scored && g.pipes[i].x+pipeWidth < birdX {
			g.pipes[i].scored = true
			g.score++
		}
	}

	// 画面外のパイプを消して、右に新しいパイプを足す
	if g.pipes[0].x+pipeWidth < -10 {
		lastX := g.pipes[len(g.pipes)-1].x
		newPipe := pipe{
			x:    lastX + 270,
			gapY: 225 + g.rng.Float64()*230,
		}
		g.pipes = append(g.pipes[1:], newPipe)
	}

	// 地面・天井・パイプに当たったらゲームオーバー
	hitGround := g.birdY+17 >= groundY
	hitCeiling := g.birdY-17 <= 0
	if hitGround || hitCeiling || g.hitPipe() {
		g.gameOver = true
		if g.score > g.best {
			g.best = g.score
		}
	}
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)
	for _, p := range g.pipes {
		g.drawPipe(screen, p)
	}
	g.drawGround(screen)
	g.drawBird(screen)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE  %02d", g.score), 198, 28)
	if !g.started {
		g.drawPanel(screen, "TAP TO FLY", "SPACE / CLICK / TOUCH")
	}
	if g.gameOver {
		g.drawPanel(screen, "GAME OVER", fmt.Sprintf("SCORE %d   BEST %d\n\nTAP TO RETRY", g.score, g.best))
	}
}

func (g *game) drawBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 21, 43, 255})
	for i := 0; i < 26; i++ {
		x := float32((i*83-g.frame/5)%540 - 30)
		y := float32(35 + (i*47)%390)
		vector.DrawFilledCircle(screen, x, y, float32(1+i%2), color.RGBA{100, 230, 225, 150}, false)
	}
	for i := 0; i < 10; i++ {
		x := float32(i*62 - (g.frame/3)%62)
		h := float32(55 + (i*29)%105)
		vector.DrawFilledRect(screen, x, groundY-h, 48, h, color.RGBA{14, 39, 67, 255}, false)
	}
}

func (g *game) drawGround(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, groundY, screenWidth, screenHeight-groundY, color.RGBA{27, 203, 159, 255}, false)
	for x := -40 + (g.frame*2)%40; x < screenWidth; x += 40 {
		vector.DrawFilledRect(screen, float32(x), groundY+13, 22, 5, color.RGBA{9, 98, 91, 255}, false)
	}
}

func (g *game) drawPipe(screen *ebiten.Image, p pipe) {
	green := color.RGBA{22, 196, 145, 255}
	light := color.RGBA{70, 236, 170, 255}
	dark := color.RGBA{6, 103, 91, 255}
	topH := float32(p.gapY - gapHeight/2)
	bottomY := float32(p.gapY + gapHeight/2)

	vector.DrawFilledRect(screen, float32(p.x), 0, pipeWidth, topH, green, false)
	vector.DrawFilledRect(screen, float32(p.x-7), topH-28, pipeWidth+14, 28, green, false)
	vector.DrawFilledRect(screen, float32(p.x), bottomY, pipeWidth, groundY-bottomY, green, false)
	vector.DrawFilledRect(screen, float32(p.x-7), bottomY, pipeWidth+14, 28, green, false)
	vector.DrawFilledRect(screen, float32(p.x+9), 0, 8, topH-28, light, false)
	vector.DrawFilledRect(screen, float32(p.x+pipeWidth-10), 0, 10, topH, dark, false)
	vector.DrawFilledRect(screen, float32(p.x+9), bottomY+28, 8, groundY-bottomY-28, light, false)
	vector.DrawFilledRect(screen, float32(p.x+pipeWidth-10), bottomY, 10, groundY-bottomY, dark, false)
}

func (g *game) drawBird(screen *ebiten.Image) {
	// オリジナル主人公「Ebi Boy」を鳥の代わりに描く
	hero.DrawCentered(screen, birdX, g.birdY, 44)
}

func (g *game) drawPanel(screen *ebiten.Image, title, detail string) {
	vector.DrawFilledRect(screen, 68, 250, 344, 166, color.RGBA{5, 16, 34, 225}, false)
	vector.StrokeRect(screen, 68, 250, 344, 166, 3, color.RGBA{45, 226, 194, 255}, false)
	ebitenutil.DebugPrintAt(screen, title, 192-len(title)*2, 286)
	ebitenutil.DebugPrintAt(screen, detail, 137, 338)
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebi Flight — Ebitengine WASM")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
