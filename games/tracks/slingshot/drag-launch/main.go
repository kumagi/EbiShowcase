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

const width, height = 480, 720

var targets = [][2]float64{{110, 178}, {370, 238}, {240, 120}}

type game struct {
	x, y, vx, vy     float64
	pullX, pullY     float64
	dragging, moving bool
	touchID          ebiten.TouchID
	touching         bool
	target, shots    int
	clear, over      bool
	message          string
}

func newGame() *game {
	g := &game{}
	g.resetBall()
	g.message = "Drag from the pearl, pull back, then release!"
	return g
}

func (g *game) resetBall() {
	g.x, g.y = 240, 560
	g.vx, g.vy = 0, 0
	g.pullX, g.pullY = 0, 110
	g.dragging, g.moving, g.touching = false, false, false
}

func (g *game) beginDrag(px, py float64) {
	dx, dy := px-g.x, py-g.y
	if dx*dx+dy*dy <= 34*34 {
		g.dragging = true
	}
}

func (g *game) setPull(px, py float64) {
	dx, dy := px-g.x, py-g.y
	length := math.Hypot(dx, dy)
	if length > 160 {
		dx, dy = dx*160/length, dy*160/length
	}
	g.pullX, g.pullY = dx, dy
}

func (g *game) launch() {
	if g.moving || math.Hypot(g.pullX, g.pullY) < 24 {
		return
	}
	// The launch vector points from the dragged pointer back to the ball.
	g.vx = -g.pullX * 0.092
	g.vy = -g.pullY * 0.092
	g.moving, g.dragging = true, false
	g.shots++
	g.message = "Launch velocity = -drag vector"
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	if g.moving {
		g.x += g.vx
		g.y += g.vy
		g.vx *= 0.992
		g.vy *= 0.992
		tx, ty := targets[g.target][0], targets[g.target][1]
		dx, dy := g.x-tx, g.y-ty
		if dx*dx+dy*dy <= 47*47 {
			g.target++
			if g.target >= len(targets) {
				g.clear = true
				g.message = "All three signal rings reached!"
			} else {
				g.message = "Hit! Aim for the next signal ring."
				g.resetBall()
			}
			return nil
		}
		if g.x < -24 || g.x > width+24 || g.y < -24 || g.y > height+24 || math.Hypot(g.vx, g.vy) < 0.32 {
			if g.shots >= 6 {
				g.over = true
				g.message = "Six launches used. Read the opposite arrow!"
			} else {
				g.message = "Missed — pull farther or change the angle."
				g.resetBall()
			}
		}
		return nil
	}

	// Mouse drag.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		px, py := ebiten.CursorPosition()
		g.beginDrag(float64(px), float64(py))
	}
	if g.dragging && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		px, py := ebiten.CursorPosition()
		g.setPull(float64(px), float64(py))
	}
	if g.dragging && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.launch()
	}

	// Touch drag uses the same pointer-to-ball vector.
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		px, py := ebiten.TouchPosition(ids[0])
		g.beginDrag(float64(px), float64(py))
		if g.dragging {
			g.touchID, g.touching = ids[0], true
		}
	}
	if g.dragging && g.touching {
		released := false
		for _, id := range inpututil.AppendJustReleasedTouchIDs(nil) {
			if id == g.touchID {
				released = true
			}
		}
		if released {
			g.touching = false
			g.launch()
		} else {
			px, py := ebiten.TouchPosition(g.touchID)
			g.setPull(float64(px), float64(py))
		}
	}

	// Keyboard alternative: arrows move the virtual pull point, Space fires.
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.pullX -= 2.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.pullX += 2.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.pullY -= 2.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.pullY += 2.5
	}
	length := math.Hypot(g.pullX, g.pullY)
	if length > 160 {
		g.pullX, g.pullY = g.pullX*160/length, g.pullY*160/length
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.launch()
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{11, 23, 43, 255})
	ebitenutil.DebugPrintAt(screen, "EBI VECTOR LAUNCH", 178, 24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("RINGS %d/3   SHOTS %d/6", g.target, g.shots), 162, 52)
	for i, p := range targets {
		c := color.RGBA{53, 80, 111, 255}
		if i == g.target {
			c = color.RGBA{239, 174, 62, 255}
		}
		vector.StrokeCircle(screen, float32(p[0]), float32(p[1]), 35, 6, c, true)
		if i < g.target {
			vector.DrawFilledCircle(screen, float32(p[0]), float32(p[1]), 13, color.RGBA{80, 198, 143, 255}, true)
		}
	}
	if !g.moving && !g.clear && !g.over {
		px, py := float32(g.x+g.pullX), float32(g.y+g.pullY)
		vector.StrokeLine(screen, float32(g.x), float32(g.y), px, py, 5, color.RGBA{220, 102, 89, 255}, true)
		vector.StrokeLine(screen, float32(g.x), float32(g.y), float32(g.x-g.pullX), float32(g.y-g.pullY), 4, color.RGBA{245, 205, 78, 255}, true)
		vector.DrawFilledCircle(screen, px, py, 8, color.RGBA{220, 102, 89, 255}, true)
	}
	vector.DrawFilledCircle(screen, float32(g.x), float32(g.y), 19, color.RGBA{105, 207, 229, 255}, true)
	vector.StrokeCircle(screen, float32(g.x), float32(g.y), 19, 3, color.White, true)
	ebitenutil.DebugPrintAt(screen, g.message, 65, 624)
	ebitenutil.DebugPrintAt(screen, "DRAG PEARL / ARROWS AIM / SPACE LAUNCH", 86, 672)
	if g.clear || g.over {
		title := "VECTOR LAUNCH CLEAR!"
		if g.over {
			title = "OUT OF LAUNCHES!"
		}
		vector.DrawFilledRect(screen, 42, 276, 396, 154, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 150, 322)
		ebitenutil.DebugPrintAt(screen, g.message, 92, 355)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 394)
	}
}

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Vector Launch — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
