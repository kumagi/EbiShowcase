package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 720
	cellSize     = 24
)

// マス目上の1点
type point struct {
	x, y int
}

type game struct {
	body     []point // [0] が頭
	dir      point   // いま進んでいる向き
	nextDir  point   // 次の移動で使う向き
	food     point
	frame    int
	score    int
	gameOver bool
	rng      *rand.Rand
}

func newGame() *game {
	g := &game{
		body: []point{
			{10, 14},
			{9, 14},
			{8, 14},
		},
		dir:     point{1, 0},
		nextDir: point{1, 0},
		rng:     rand.New(rand.NewSource(31)),
	}
	g.placeFood()
	return g
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
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

// 逆向きには曲がれない
func (g *game) setDir(d point) {
	if d.x+g.dir.x == 0 && d.y+g.dir.y == 0 {
		return
	}
	g.nextDir = d
}

func (g *game) placeFood() {
	for {
		p := point{
			x: g.rng.Intn(screenWidth / cellSize),
			y: 2 + g.rng.Intn(screenHeight/cellSize-4),
		}
		onBody := false
		for _, b := range g.body {
			if b == p {
				onBody = true
				break
			}
		}
		if !onBody {
			g.food = p
			return
		}
	}
}

func (g *game) readInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.setDir(point{0, -1})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.setDir(point{0, 1})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.setDir(point{-1, 0})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.setDir(point{1, 0})
	}

	// タップした方向へ曲がる
	px, py, pressed := 0, 0, false
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		px, py = ebiten.CursorPosition()
		pressed = true
	}
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		px, py = ebiten.TouchPosition(ids[0])
		pressed = true
	}
	if !pressed {
		return
	}

	head := g.body[0]
	headPX := head.x*cellSize + cellSize/2
	headPY := head.y*cellSize + cellSize/2
	dx := px - headPX
	dy := py - headPY
	if abs(dx) > abs(dy) {
		if dx < 0 {
			g.setDir(point{-1, 0})
		} else {
			g.setDir(point{1, 0})
		}
	} else {
		if dy < 0 {
			g.setDir(point{0, -1})
		} else {
			g.setDir(point{0, 1})
		}
	}
}

func (g *game) stepSnake() {
	g.dir = g.nextDir
	head := point{
		x: g.body[0].x + g.dir.x,
		y: g.body[0].y + g.dir.y,
	}

	// 壁
	maxX := screenWidth / cellSize
	maxY := screenHeight / cellSize
	if head.x < 0 || head.x >= maxX || head.y < 2 || head.y >= maxY-2 {
		g.gameOver = true
		return
	}

	// 自分の体
	for _, b := range g.body {
		if head == b {
			g.gameOver = true
			return
		}
	}

	// 頭を前に足す
	g.body = append([]point{head}, g.body...)

	if head == g.food {
		g.score++
		g.placeFood()
		return // しっぽを残して長くする
	}
	// 食べていなければしっぽを1つ消す
	g.body = g.body[:len(g.body)-1]
}

// --- ここから Update ---
func (g *game) Update() error {
	if g.gameOver {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}

	g.readInput()

	g.frame++
	// スコアが上がると少し速くなる
	wait := max(4, 10-g.score/3)
	if g.frame%wait != 0 {
		return nil
	}
	g.stepSnake()
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})

	for x := 0; x < screenWidth; x += cellSize {
		for y := cellSize * 2; y < screenHeight-cellSize*2; y += cellSize {
			vector.StrokeRect(screen, float32(x), float32(y), cellSize, cellSize, 1, color.RGBA{45, 226, 194, 20}, false)
		}
	}

	vector.DrawFilledCircle(
		screen,
		float32(g.food.x*cellSize+cellSize/2),
		float32(g.food.y*cellSize+cellSize/2),
		9,
		color.RGBA{255, 105, 79, 255},
		false,
	)

	for i, b := range g.body {
		c := color.RGBA{45, 226, 194, 255}
		if i == 0 {
			c = color.RGBA{255, 210, 72, 255}
		}
		vector.DrawFilledRect(
			screen,
			float32(b.x*cellSize+2),
			float32(b.y*cellSize+2),
			float32(cellSize-4),
			float32(cellSize-4),
			c,
			false,
		)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %02d", g.score), 205, 18)
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER\n\nTAP / SPACE TO RETRY", 160, 340)
	} else {
		ebitenutil.DebugPrintAt(screen, "ARROWS / TAP A DIRECTION", 150, 685)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
