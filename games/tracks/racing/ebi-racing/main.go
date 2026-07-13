package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const (
	W = 480
	H = 720
)

type car struct {
	x, y, angle, speed float64
	lap, next          int
	ai                 bool
}

var gates = [][2]float64{{240, 90}, {400, 220}, {380, 500}, {240, 620}, {90, 500}, {70, 220}}

type game struct {
	cars      []car
	frames    int
	won, lost bool
	message   string
}

func newGame() *game {
	return &game{cars: []car{{240, 125, 0, 0, 0, 1, false}, {220, 145, 0, 0, 0, 1, true}}, message: "Accelerate, steer, and pass the glowing gates in order."}
}
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	p := &g.cars[0]
	gas, brake, left, right := controls()
	if gas {
		p.speed += .045
	}
	if brake {
		p.speed -= .06
	}
	p.speed *= .985
	p.speed = math.Max(-1.5, math.Min(4.8, p.speed))
	turn := .035 * (.4 + math.Abs(p.speed)/4.8)
	if left {
		p.angle -= turn
	}
	if right {
		p.angle += turn
	}
	g.move(p)
	a := &g.cars[1]
	tx, ty := gates[a.next][0], gates[a.next][1]
	target := math.Atan2(ty-a.y, tx-a.x) + math.Pi/2
	diff := angleDiff(target, a.angle)
	a.angle += math.Max(-.025, math.Min(.025, diff))
	a.speed = 3.0
	g.move(a)
	g.checkGate(p)
	g.checkGate(a)
	if p.lap >= 2 {
		g.won = true
		g.message = "Two laps complete!"
	}
	if a.lap >= 2 || g.frames > 100*60 {
		g.lost = true
		g.message = "The rival finished first."
	}
	return nil
}
func (g *game) move(c *car) {
	c.x += math.Sin(c.angle) * c.speed
	c.y -= math.Cos(c.angle) * c.speed
	if !onRoad(c.x, c.y) {
		c.speed *= .91
	}
	c.x = math.Max(25, math.Min(W-25, c.x))
	c.y = math.Max(65, math.Min(650, c.y))
}
func onRoad(x, y float64) bool {
	dx, dy := (x-240)/175, (y-355)/265
	r := dx*dx + dy*dy
	return r > .38 && r < 1.22
}
func (g *game) checkGate(c *car) {
	q := gates[c.next]
	if math.Hypot(c.x-q[0], c.y-q[1]) < 48 {
		c.next = (c.next + 1) % len(gates)
		if c.next == 0 {
			c.lap++
		}
	}
}
func angleDiff(a, b float64) float64 { d := math.Mod(a-b+math.Pi, 2*math.Pi) - math.Pi; return d }
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{42, 112, 67, 255})
	drawEllipse(s, 240, 355, 178, 268, 90, color.RGBA{68, 72, 84, 255})
	drawEllipse(s, 240, 355, 178, 268, 3, color.RGBA{255, 255, 255, 255})
	q := gates[g.cars[0].next]
	vector.StrokeCircle(s, float32(q[0]), float32(q[1]), 32, 6, color.RGBA{246, 198, 72, 255}, false)
	for i, c := range g.cars {
		col := color.RGBA{235, 91, 76, 255}
		if i == 1 {
			col = color.RGBA{76, 166, 232, 255}
		}
		op := &ebiten.DrawImageOptions{}
		img := ebiten.NewImage(26, 40)
		img.Fill(col)
		op.GeoM.Translate(-13, -20)
		op.GeoM.Rotate(c.angle)
		op.GeoM.Translate(c.x, c.y)
		s.DrawImage(img, op)
	}
	vector.DrawFilledRect(s, 0, 0, W, 60, color.RGBA{8, 17, 31, 230}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LAP %d/2  SPEED %.1f  NEXT GATE %d  RIVAL LAP %d", g.cars[0].lap+1, math.Abs(g.cars[0].speed), g.cars[0].next+1, g.cars[1].lap+1), 55, 20)
	labels := []string{"LEFT", "GAS", "BRAKE", "RIGHT"}
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*120+5), 650, 110, 55, color.RGBA{45, 78, 113, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*120+38, 675)
	}
	if g.won {
		overlay(s, "RACE WIN!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(s, "RACE LOST\n\nTAP / ENTER TO RETRY")
	}
}

func drawEllipse(s *ebiten.Image, cx, cy, rx, ry, width float32, c color.RGBA) {
	const segments = 72
	for i := 0; i < segments; i++ {
		a := float64(i) * 2 * math.Pi / segments
		b := float64(i+1) * 2 * math.Pi / segments
		vector.StrokeLine(s, cx+rx*float32(math.Cos(a)), cy+ry*float32(math.Sin(a)), cx+rx*float32(math.Cos(b)), cy+ry*float32(math.Sin(b)), width, c, false)
	}
}
func controls() (bool, bool, bool, bool) {
	gas := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	brake := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 640 {
			left = x < 120
			gas = x >= 120 && x < 240
			brake = x >= 240 && x < 360
			right = x >= 360
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 640 {
			left = x < 120
			gas = x >= 120 && x < 240
			brake = x >= 240 && x < 360
			right = x >= 360
		}
	}
	return gas, brake, left, right
}
func retry() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 50, 270, 380, 150, color.RGBA{4, 10, 24, 245}, false)
	ebitenutil.DebugPrintAt(s, t, 155, 330)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Circuit")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
