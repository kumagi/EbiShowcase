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

const (
	screenW = 480
	screenH = 720
	groundY = 650
	gravity = 0.42
	flapV   = -7.4
	pipeW   = 76
	gapH    = 178
)

type pipe struct {
	x, gapY float64
	passed  bool
}

type game struct {
	birdY, velocity float64
	pipes           []pipe
	score, best     int
	frame           int
	started, over   bool
	rng             *rand.Rand
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(7))}
	g.reset()
	return g
}

func (g *game) reset() {
	g.birdY, g.velocity = 320, 0
	g.pipes = []pipe{{x: 560, gapY: 300}, {x: 830, gapY: 390}, {x: 1100, gapY: 265}}
	g.score, g.frame, g.started, g.over = 0, 0, false, false
}

func justPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Update() error {
	g.frame++
	pressed := justPressed()
	if g.over {
		if pressed {
			g.reset()
		}
		return nil
	}
	if !g.started {
		g.birdY = 320 + math.Sin(float64(g.frame)*0.06)*7
		if pressed {
			g.started, g.velocity = true, flapV
		}
		return nil
	}
	if pressed {
		g.velocity = flapV
	}
	g.velocity += gravity
	g.birdY += g.velocity

	for i := range g.pipes {
		g.pipes[i].x -= 2.8
		if !g.pipes[i].passed && g.pipes[i].x+pipeW < 128 {
			g.pipes[i].passed = true
			g.score++
		}
	}
	if g.pipes[0].x+pipeW < -10 {
		g.pipes = append(g.pipes[1:], pipe{x: g.pipes[len(g.pipes)-1].x + 270, gapY: 225 + g.rng.Float64()*230})
	}
	if g.birdY+17 >= groundY || g.birdY-17 <= 0 || g.hitPipe() {
		g.over = true
		if g.score > g.best {
			g.best = g.score
		}
	}
	return nil
}

func (g *game) hitPipe() bool {
	const bx, br = 128.0, 15.0
	for _, p := range g.pipes {
		if bx+br > p.x && bx-br < p.x+pipeW && (g.birdY-br < p.gapY-gapH/2 || g.birdY+br > p.gapY+gapH/2) {
			return true
		}
	}
	return false
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{9, 21, 43, 255})
	// Far skyline and drifting stars make the canvas feel alive without assets.
	for i := 0; i < 26; i++ {
		x := float32((i*83-g.frame/5)%540 - 30)
		y := float32(35 + (i*47)%390)
		vector.DrawFilledCircle(s, x, y, float32(1+i%2), color.RGBA{100, 230, 225, 150}, false)
	}
	for i := 0; i < 10; i++ {
		x := float32(i*62 - (g.frame/3)%62)
		h := float32(55 + (i*29)%105)
		vector.DrawFilledRect(s, x, groundY-h, 48, h, color.RGBA{14, 39, 67, 255}, false)
	}
	for _, p := range g.pipes {
		g.drawPipe(s, p)
	}
	vector.DrawFilledRect(s, 0, groundY, screenW, screenH-groundY, color.RGBA{27, 203, 159, 255}, false)
	for x := -40 + (g.frame*2)%40; x < screenW; x += 40 {
		vector.DrawFilledRect(s, float32(x), groundY+13, 22, 5, color.RGBA{9, 98, 91, 255}, false)
	}
	g.drawBird(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE  %02d", g.score), 198, 28)
	if !g.started {
		g.panel(s, "TAP TO FLY", "SPACE / CLICK / TOUCH")
	}
	if g.over {
		g.panel(s, "GAME OVER", fmt.Sprintf("SCORE %d   BEST %d\n\nTAP TO RETRY", g.score, g.best))
	}
}

func (g *game) drawPipe(s *ebiten.Image, p pipe) {
	green, light, dark := color.RGBA{22, 196, 145, 255}, color.RGBA{70, 236, 170, 255}, color.RGBA{6, 103, 91, 255}
	topH := float32(p.gapY - gapH/2)
	bottomY := float32(p.gapY + gapH/2)
	vector.DrawFilledRect(s, float32(p.x), 0, pipeW, topH, green, false)
	vector.DrawFilledRect(s, float32(p.x-7), topH-28, pipeW+14, 28, green, false)
	vector.DrawFilledRect(s, float32(p.x), bottomY, pipeW, groundY-bottomY, green, false)
	vector.DrawFilledRect(s, float32(p.x-7), bottomY, pipeW+14, 28, green, false)
	vector.DrawFilledRect(s, float32(p.x+9), 0, 8, topH-28, light, false)
	vector.DrawFilledRect(s, float32(p.x+pipeW-10), 0, 10, topH, dark, false)
	vector.DrawFilledRect(s, float32(p.x+9), bottomY+28, 8, groundY-bottomY-28, light, false)
	vector.DrawFilledRect(s, float32(p.x+pipeW-10), bottomY, 10, groundY-bottomY, dark, false)
}

func (g *game) drawBird(s *ebiten.Image) {
	y := float32(g.birdY)
	wingY := y + float32(math.Sin(float64(g.frame)*0.5)*5)
	vector.DrawFilledCircle(s, 128, y, 18, color.RGBA{255, 210, 72, 255}, false)
	vector.DrawFilledCircle(s, 119, wingY, 10, color.RGBA{255, 126, 75, 255}, false)
	vector.DrawFilledCircle(s, 137, y-6, 6, color.White, false)
	vector.DrawFilledCircle(s, 139, y-6, 2.5, color.RGBA{9, 21, 43, 255}, false)
	vector.DrawFilledRect(s, 144, y-1, 16, 6, color.RGBA{255, 93, 77, 255}, false)
}

func (g *game) panel(s *ebiten.Image, title, detail string) {
	vector.DrawFilledRect(s, 68, 250, 344, 166, color.RGBA{5, 16, 34, 225}, false)
	vector.StrokeRect(s, 68, 250, 344, 166, 3, color.RGBA{45, 226, 194, 255}, false)
	ebitenutil.DebugPrintAt(s, title, 192-len(title)*2, 286)
	ebitenutil.DebugPrintAt(s, detail, 137, 338)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Flight — Ebitengine WASM")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
