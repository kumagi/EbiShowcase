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

const (
	screenW = 480
	screenH = 720
	radius  = 18.0
)

type vec struct{ x, y float64 }

var frictionValues = [...]float64{0.990, 0.978, 0.960}
var frictionNames = [...]string{"ICE", "NORMAL", "SAND"}
var targets = [...]vec{{115, 185}, {365, 245}, {235, 365}}

type game struct {
	pos, velocity      vec
	dragStart, dragNow vec
	dragging, moving   bool
	frictionIndex      int
	stopThreshold      float64
	shots, hits        int
	message            string
	won, lost          bool
}

func newGame() *game {
	g := &game{frictionIndex: 1, stopThreshold: 0.10}
	g.resetPuck()
	g.message = "Choose friction, then drag from the puck."
	return g
}

func (g *game) resetPuck() {
	g.pos = vec{240, 520}
	g.velocity = vec{}
	g.dragging = false
	g.moving = false
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}

	if !g.moving {
		g.updateControls()
		g.updateDrag()
		return nil
	}

	g.pos.x += g.velocity.x
	g.pos.y += g.velocity.y
	if g.pos.x < 38 {
		g.pos.x = 38
		g.velocity.x = math.Abs(g.velocity.x) * 0.82
	}
	if g.pos.x > 442 {
		g.pos.x = 442
		g.velocity.x = -math.Abs(g.velocity.x) * 0.82
	}
	if g.pos.y < 105 {
		g.pos.y = 105
		g.velocity.y = math.Abs(g.velocity.y) * 0.82
	}
	if g.pos.y > 550 {
		g.pos.y = 550
		g.velocity.y = -math.Abs(g.velocity.y) * 0.82
	}

	friction := frictionValues[g.frictionIndex]
	g.velocity.x *= friction
	g.velocity.y *= friction
	if speed(g.velocity) < g.stopThreshold {
		g.velocity = vec{} // snap to exact zero instead of drifting forever
		g.moving = false
		g.finishShot()
	}
	return nil
}

func (g *game) updateControls() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.frictionIndex = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.frictionIndex = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.frictionIndex = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.changeThreshold(-0.02)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.changeThreshold(0.02)
	}
	if x, y, ok := pressPosition(); ok && y >= 590 {
		switch {
		case x < 105:
			g.frictionIndex = 0
		case x < 205:
			g.frictionIndex = 1
		case x < 305:
			g.frictionIndex = 2
		case x < 390:
			g.changeThreshold(-0.02)
		default:
			g.changeThreshold(0.02)
		}
		g.message = "Settings changed. Predict the stopping place!"
	}
}

func (g *game) changeThreshold(delta float64) {
	g.stopThreshold = math.Max(0.04, math.Min(0.24, g.stopThreshold+delta))
}

func (g *game) updateDrag() {
	if !g.dragging {
		x, y, ok := pressPosition()
		if ok && y < 580 && distance(vec{float64(x), float64(y)}, g.pos) <= radius+18 {
			g.dragging = true
			g.dragStart = g.pos
			g.dragNow = vec{float64(x), float64(y)}
		}
		return
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.dragNow = vec{float64(x), float64(y)}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		g.dragNow = vec{float64(x), float64(y)}
	}
	mouseReleased := inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
	touchReleased := len(inpututil.AppendJustReleasedTouchIDs(nil)) > 0
	if mouseReleased || touchReleased {
		pull := vec{g.dragStart.x - g.dragNow.x, g.dragStart.y - g.dragNow.y}
		length := speed(pull)
		if length > 150 {
			pull.x *= 150 / length
			pull.y *= 150 / length
		}
		g.dragging = false
		if length < 12 {
			g.message = "Pull farther to launch."
			return
		}
		g.velocity = vec{pull.x * 0.095, pull.y * 0.095}
		g.moving = true
		g.shots++
		g.message = "Moving: velocity shrinks every frame."
	}
}

func (g *game) finishShot() {
	target := targets[g.hits]
	d := distance(g.pos, target)
	if d <= 54 {
		g.hits++
		g.message = fmt.Sprintf("Perfect stop! Target %d/3 cleared.", g.hits)
		if g.hits == len(targets) {
			g.won = true
			return
		}
	} else {
		g.message = fmt.Sprintf("Stopped %.0f px from target. Adjust friction!", d)
	}
	if g.shots >= 6 {
		g.lost = true
		return
	}
	g.resetPuck()
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{11, 23, 42, 255})
	ebitenutil.DebugPrintAt(screen, "FRICTION STOP CHALLENGE", 153, 22)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TARGETS %d/3   SHOTS %d/6   SPEED %.2f", g.hits, g.shots, speed(g.velocity)), 94, 49)
	ebitenutil.DebugPrintAt(screen, g.message, 76, 73)
	vector.DrawFilledRect(screen, 20, 88, 440, 485, color.RGBA{30, 56, 74, 255}, false)
	vector.StrokeRect(screen, 20, 88, 440, 485, 4, color.RGBA{104, 150, 166, 255}, false)

	for i, target := range targets {
		c := color.RGBA{245, 190, 67, 255}
		if i < g.hits {
			c = color.RGBA{87, 184, 113, 255}
		}
		vector.DrawFilledCircle(screen, float32(target.x), float32(target.y), 52, color.RGBA{c.R, c.G, c.B, 42}, false)
		vector.StrokeCircle(screen, float32(target.x), float32(target.y), 52, 4, c, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i+1), int(target.x)-3, int(target.y)-6)
	}
	if g.dragging {
		vector.StrokeLine(screen, float32(g.pos.x), float32(g.pos.y), float32(g.dragNow.x), float32(g.dragNow.y), 5, color.RGBA{244, 181, 65, 255}, false)
		launchX := g.pos.x + (g.pos.x - g.dragNow.x)
		launchY := g.pos.y + (g.pos.y - g.dragNow.y)
		vector.StrokeLine(screen, float32(g.pos.x), float32(g.pos.y), float32(launchX), float32(launchY), 2, color.White, false)
	}
	vector.DrawFilledCircle(screen, float32(g.pos.x), float32(g.pos.y), radius, color.RGBA{236, 99, 84, 255}, false)
	vector.StrokeCircle(screen, float32(g.pos.x), float32(g.pos.y), radius, 3, color.White, false)

	labels := [...]string{"ICE .990", "NORMAL .978", "SAND .960", "STOP -", "STOP +"}
	starts := [...]int{5, 105, 205, 305, 390}
	widths := [...]int{95, 95, 95, 80, 85}
	for i, label := range labels {
		c := color.RGBA{52, 83, 119, 255}
		if i == g.frictionIndex {
			c = color.RGBA{193, 123, 57, 255}
		}
		vector.DrawFilledRect(screen, float32(starts[i]), 590, float32(widths[i]), 62, c, false)
		ebitenutil.DebugPrintAt(screen, label, starts[i]+10, 615)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FRICTION %.3f  |  STOP BELOW %.2f  (keys 1-3, -/+)", frictionValues[g.frictionIndex], g.stopThreshold), 62, 672)
	if g.won {
		overlay(screen, "THREE PERFECT STOPS!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "OUT OF SHOTS\n\nTAP / ENTER TO RETRY")
	}
}

func speed(v vec) float64       { return math.Hypot(v.x, v.y) }
func distance(a, b vec) float64 { return math.Hypot(a.x-b.x, a.y-b.y) }

func pressPosition() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}

func retryPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	_, _, ok := pressPosition()
	return ok
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 45, 270, 390, 150, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 45, 270, 390, 150, 4, color.RGBA{243, 189, 70, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 116, 325)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Friction Stop Challenge — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
