package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
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

type course struct {
	rx, ry  float64
	road    float64
	aiSpeed float64
	laps    int
	name    string
	bg      color.RGBA
}

var courses = []course{{175, 265, .42, 2.8, 1, "CORAL COAST", color.RGBA{18, 93, 105, 255}}, {150, 245, .34, 3.15, 2, "PEARL REEF", color.RGBA{20, 42, 83, 255}}, {195, 225, .30, 3.45, 2, "STORM CITADEL", color.RGBA{26, 31, 50, 255}}}

type spark struct{ x, y, vx, vy, life float64 }

type game struct {
	cars                                  []car
	frames                                int
	won, lost                             bool
	message                               string
	stage, totalFrames, bestFrames, shake int
	gates                                 [][2]float64
	sparks                                []spark
	audio                                 *audio.Context
	gate                                  audiolab.Gate
	pulse                                 *shaderlab.Pulse
	cam                                   cameralab.State
	badge                                 *ebiten.Image
}

func newGame() *game {
	prepareRacingArt()
	g := &game{stage: 1}
	g.audio = audiolab.Context()
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: W / 2, Y: H / 2}, ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{246, 198, 72, 255})
	g.loadCourse()
	return g
}
func (g *game) loadCourse() {
	c := courses[g.stage-1]
	g.gates = nil
	for i := 0; i < 8; i++ {
		a := -math.Pi/2 + float64(i)*math.Pi/4
		g.gates = append(g.gates, [2]float64{240 + math.Cos(a)*c.rx, 355 + math.Sin(a)*c.ry})
	}
	start := g.gates[0]
	// Begin in motion: the first thumbnail already shows a race, not two parked
	// rectangles waiting for the learner to discover the accelerator.
	g.cars = []car{{start[0] + 28, start[1] + 36, 0, 2.6, 0, 1, false}, {start[0] - 28, start[1] + 58, 0, c.aiSpeed * .82, 0, 1, true}}
	g.frames = 0
	g.won = false
	g.lost = false
	g.message = "Accelerate, steer, and pass the glowing gates in order."
}
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			if g.won && g.stage < 3 {
				g.totalFrames += g.frames
				g.stage++
				g.loadCourse()
			} else {
				best := g.bestFrames
				if g.won {
					total := g.totalFrames + g.frames
					if best == 0 || total < best {
						best = total
					}
				}
				*g = *newGame()
				g.bestFrames = best
			}
		}
		return nil
	}
	g.frames++
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
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
	tx, ty := g.gates[a.next][0], g.gates[a.next][1]
	target := math.Atan2(ty-a.y, tx-a.x) + math.Pi/2
	diff := angleDiff(target, a.angle)
	a.angle += math.Max(-.025, math.Min(.025, diff))
	a.speed = courses[g.stage-1].aiSpeed
	g.move(a)
	g.checkGate(p)
	g.checkGate(a)
	if p.lap >= courses[g.stage-1].laps {
		g.won = true
		g.message = "Course complete!"
		g.burst(p.x, p.y, 28)
	}
	if a.lap >= courses[g.stage-1].laps || g.frames > 100*60 {
		g.lost = true
		g.message = "The rival finished first."
	}
	return nil
}
func (g *game) move(c *car) {
	c.x += math.Sin(c.angle) * c.speed
	c.y -= math.Cos(c.angle) * c.speed
	if !g.onRoad(c.x, c.y) {
		c.speed *= .91
		if g.frames%8 == 0 {
			g.sparks = append(g.sparks, spark{c.x, c.y, 0, .4, 18})
		}
	}
	c.x = math.Max(25, math.Min(W-25, c.x))
	c.y = math.Max(65, math.Min(650, c.y))
}
func (g *game) onRoad(x, y float64) bool {
	c := courses[g.stage-1]
	dx, dy := (x-240)/c.rx, (y-355)/c.ry
	r := dx*dx + dy*dy
	return r > c.road && r < 1.22
}
func (g *game) checkGate(c *car) {
	q := g.gates[c.next]
	if math.Hypot(c.x-q[0], c.y-q[1]) < 48 {
		c.next = (c.next + 1) % len(g.gates)
		if c.next == 0 {
			c.lap++
			g.play(760)
		}
	}
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Square, freq, .045)).Play()
}
func (g *game) burst(x, y float64, n int) {
	g.shake = 7
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 28 + float64(i%9)})
	}
}
func angleDiff(a, b float64) float64 { d := math.Mod(a-b+math.Pi, 2*math.Pi) - math.Pi; return d }
func (g *game) Draw(s *ebiten.Image) {
	course := courses[g.stage-1]
	s.Fill(course.bg)
	drawCourseArt(s, g.stage)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.frames)*2) * 5
	}
	roadColors := []color.RGBA{{35, 45, 52, 205}, {19, 31, 54, 205}, {18, 22, 33, 212}}
	laneColors := []color.RGBA{{255, 239, 191, 230}, {137, 235, 255, 230}, {95, 224, 255, 230}}
	drawEllipse(s, 240+float32(ox), 355, float32(course.rx), float32(course.ry), 82, roadColors[g.stage-1])
	drawEllipse(s, 240+float32(ox), 355, float32(course.rx), float32(course.ry), 3, laneColors[g.stage-1])
	// Kerbs and lane markers sell speed without changing the simple oval math.
	for i := 0; i < 32; i++ {
		a := float64(i) * math.Pi / 16
		x := float32(240 + math.Cos(a)*course.rx)
		y := float32(355 + math.Sin(a)*course.ry)
		c := color.RGBA{245, 238, 215, 255}
		if i%2 == 0 {
			c = color.RGBA{232, 72, 72, 255}
		}
		vector.DrawFilledCircle(s, x, y, 6, c, true)
	}
	q := g.gates[g.cars[0].next]
	gatePulse := float32(22 + math.Sin(float64(g.frames)*.13)*2)
	// Two thin rings and one directional chevron keep the target legible while
	// leaving the generated course art visible through the center.
	vector.StrokeCircle(s, float32(q[0]), float32(q[1]), gatePulse, 4, color.RGBA{255, 218, 95, 245}, true)
	vector.StrokeCircle(s, float32(q[0]), float32(q[1]), gatePulse-6, 2, color.RGBA{92, 231, 255, 210}, true)
	dx, dy := q[0]-g.cars[0].x, q[1]-g.cars[0].y
	d := math.Max(1, math.Hypot(dx, dy))
	ux, uy := dx/d, dy/d
	vx, vy := -uy, ux
	tipX, tipY := q[0]-ux*float64(gatePulse+3), q[1]-uy*float64(gatePulse+3)
	baseX, baseY := q[0]-ux*float64(gatePulse+14), q[1]-uy*float64(gatePulse+14)
	vector.StrokeLine(s, float32(baseX+vx*7), float32(baseY+vy*7), float32(tipX), float32(tipY), 4, color.RGBA{92, 231, 255, 235}, true)
	vector.StrokeLine(s, float32(baseX-vx*7), float32(baseY-vy*7), float32(tipX), float32(tipY), 4, color.RGBA{92, 231, 255, 235}, true)
	for i, c := range g.cars {
		drawCar(s, c.x+ox, c.y, c.angle, i == 0, g.frames)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	vector.DrawFilledRect(s, 0, 0, W, 60, color.RGBA{8, 17, 31, 230}, false)
	g.drawHUD(s, course)
	g.drawEffectBadge(s)
	labels := []string{"LEFT", "GAS", "BRAKE", "RIGHT"}
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*120+5), 650, 110, 55, color.RGBA{45, 78, 113, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*120+38, 675)
	}
	if g.won {
		msg := "COURSE WIN!\n\nTAP / ENTER FOR NEXT COURSE"
		if g.stage == 3 {
			msg = "CUP COMPLETE!\n\nTAP / ENTER FOR A NEW CUP"
		}
		overlay(s, msg)
	}
	if g.lost {
		overlay(s, "RACE LOST\n\nTAP / ENTER TO RETRY")
	}
}

func drawCar(s *ebiten.Image, x, y, angle float64, player bool, frame int) {
	// Exhaust light and shadow are derived from the current car state. Neither
	// changes speed or position; the generated vehicle is only a projection.
	rearX := x - math.Sin(angle)*30
	rearY := y + math.Cos(angle)*30
	glow := color.RGBA{85, 225, 255, 175}
	if player {
		glow = color.RGBA{255, 203, 91, 185}
	}
	vector.DrawFilledCircle(s, float32(x+3), float32(y+5), 23, color.RGBA{2, 7, 15, 105}, true)
	vector.DrawFilledCircle(s, float32(rearX), float32(rearY), float32(7+math.Sin(float64(frame)*.3)*2), glow, true)
	name := "rival"
	if player {
		name = "player"
	}
	drawRaceCar(s, name, x, y, angle, 82)
	if player {
		pulse := float32(34 + math.Sin(float64(frame)*.16)*3)
		vector.StrokeCircle(s, float32(x), float32(y), pulse, 2, color.RGBA{255, 226, 112, 190}, true)
	}
}
func (g *game) drawHUD(s *ebiten.Image, c course) {
	place := 2
	playerProgress := g.cars[0].lap*len(g.gates) + g.cars[0].next
	rivalProgress := g.cars[1].lap*len(g.gates) + g.cars[1].next
	if playerProgress >= rivalProgress {
		place = 1
	}
	label := fmt.Sprintf("YOU %d/2  COURSE %d/3 %s  LAP %d/%d  SPEED %.1f", place, g.stage, c.name, g.cars[0].lap+1, c.laps, math.Abs(g.cars[0].speed))
	if face, err := uilab.Face("en", 14); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(25, 8)
		text.Draw(s, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(s, label, 25, 20)
}
func (g *game) drawEffectBadge(s *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frames)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(W-34, 32)
	s.DrawImage(fx, op)
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
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 640 {
			left = left || x < 120
			gas = gas || x >= 120 && x < 240
			brake = brake || x >= 240 && x < 360
			right = right || x >= 360
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
