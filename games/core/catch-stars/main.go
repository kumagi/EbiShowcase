package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type star struct{ x, y, speed float64 }
type game struct {
	basketX             float64
	stars               []star
	frame, score, lives int
	rng                 *rand.Rand
	over                bool
}

func newGame() *game { return &game{basketX: width / 2, lives: 3, rng: rand.New(rand.NewSource(19))} }

func (g *game) Update() error {
	if g.over {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || len(ebiten.AppendTouchIDs(nil)) > 0 {
			*g = *newGame()
		}
		return nil
	}
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
	if g.basketX < 55 {
		g.basketX = 55
	}
	if g.basketX > width-55 {
		g.basketX = width - 55
	}
	g.frame++
	if g.frame%45 == 0 {
		g.stars = append(g.stars, star{x: 25 + g.rng.Float64()*(width-50), y: -20, speed: 2.6 + g.rng.Float64()*2})
	}
	next := g.stars[:0]
	for _, s := range g.stars {
		s.y += s.speed
		if s.y > 625 && s.y < 675 && s.x > g.basketX-58 && s.x < g.basketX+58 {
			g.score++
			continue
		}
		if s.y > height {
			g.lives--
			if g.lives <= 0 {
				g.over = true
			}
			continue
		}
		next = append(next, s)
	}
	g.stars = next
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})
	for i := 0; i < 24; i++ {
		vector.DrawFilledCircle(screen, float32((i*97)%width), float32((i*53+g.frame/4)%height), 1.5, color.RGBA{45, 226, 194, 120}, false)
	}
	for _, s := range g.stars {
		vector.DrawFilledCircle(screen, float32(s.x), float32(s.y), 15, color.RGBA{255, 210, 72, 255}, false)
		vector.DrawFilledCircle(screen, float32(s.x-5), float32(s.y-5), 4, color.White, false)
	}
	vector.DrawFilledRect(screen, float32(g.basketX-58), 640, 116, 38, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledRect(screen, float32(g.basketX-48), 630, 96, 12, color.RGBA{255, 105, 79, 255}, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %02d   LIVES %d", g.score, g.lives), 170, 25)
	if g.over {
		ebitenutil.DebugPrintAt(screen, "GAME OVER\n\nCLICK / TOUCH TO RETRY", 155, 320)
	} else {
		ebitenutil.DebugPrintAt(screen, "MOVE: ARROWS / DRAG", 170, 690)
	}
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Catch the Stars — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
