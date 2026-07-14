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
	W      = 480
	H      = 720
	worldW = 3200
	ground = 520
)

type game struct {
	x, y, vx, vy, cam                   float64
	onGround, dash, highJump, won, lost bool
	revealed                            map[int]bool
	relics                              map[int]bool
	frames                              int
	message                             string
	flash, shake, bestFrames            int
	sparks                              []spark
	audio                               *audio.Context
	gate                                audiolab.Gate
	pulse                               *shaderlab.Pulse
	camState                            cameralab.State
	badge                               *ebiten.Image
}
type spark struct{ x, y, vx, vy, life float64 }

func newGame() *game {
	g := &game{x: 80, y: ground - 36, revealed: map[int]bool{}, relics: map[int]bool{720: true, 1450: true, 2350: true, 2980: true}, message: "Explore the huge world. The map reveals one room at a time."}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.camState = cameralab.State{Pos: cameralab.Vec{X: W / 2, Y: H / 2}, ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{245, 190, 68, 255})
	return g
}
func floorAt(x float64) float64 {
	if x > 980 && x < 1250 {
		return 610
	}
	if x > 1900 && x < 2150 {
		return 585
	}
	return ground - 45*math.Sin(x/330)
}
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			best := g.bestFrames
			*g = *newGame()
			g.bestFrames = best
		}
		return nil
	}
	g.frames++
	if g.flash > 0 {
		g.flash--
	}
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
	left, right, jump, dash := controls()
	acc := .45
	if left {
		g.vx -= acc
	}
	if right {
		g.vx += acc
	}
	g.vx *= .86
	if dash && g.dash {
		if right {
			g.vx = 9
		} else if left {
			g.vx = -9
		} else {
			g.vx = 9
		}
		for i := 0; i < 8; i++ {
			g.sparks = append(g.sparks, spark{g.x - float64(i)*5, g.y + 20, -.3, 0, 18})
		}
	}
	if jump && g.onGround {
		g.vy = -9
		if g.highJump {
			g.vy = -12
		}
		g.onGround = false
	}
	g.vy += .45
	g.x += g.vx
	if !g.dash && g.x > 1710 {
		g.x = 1710
		g.vx = 0
		g.message = "A sealed current blocks the path. Find the dash crest."
	}
	if !g.highJump && g.x > 2520 {
		g.x = 2520
		g.vx = 0
		g.message = "The crystal ledge is too high. Find the tide wings."
	}
	g.y += g.vy
	f := floorAt(g.x)
	if g.y >= f-36 {
		g.y = f - 36
		g.vy = 0
		g.onGround = true
	}
	g.x = math.Max(20, math.Min(worldW-20, g.x))
	g.cam += (g.x - 240 - g.cam) * .1
	g.cam = math.Max(0, math.Min(worldW-W, g.cam))
	room := int(g.x / 400)
	g.revealed[room] = true
	for rx := range g.relics {
		if math.Abs(g.x-float64(rx)) < 35 {
			g.play(760)
			delete(g.relics, rx)
			if rx == 1450 {
				g.dash = true
				g.message = "Dash crest found! Hold a direction and press X."
			} else if rx == 2350 {
				g.highJump = true
				g.message = "Tide wings found! Jumps now reach high ledges."
			} else {
				g.message = "Map fragment found. More of the world is recorded."
			}
			g.burst(float64(rx), g.y, 18)
		}
	}
	if len(g.relics) == 0 && g.x > 3100 {
		g.won = true
		if g.bestFrames == 0 || g.frames < g.bestFrames {
			g.bestFrames = g.frames
		}
		g.burst(g.x, g.y, 36)
	}
	for _, hx := range []float64{1080, 2050, 2790} {
		if math.Abs(g.x-hx) < 22 && g.onGround {
			g.flash = 20
			g.shake = 7
			g.x = math.Max(40, float64(int(g.x/400))*400+40)
			g.vx = 0
			g.message = "Spikes! Returned to the room entrance."
		}
	}
	if g.y > 690 || g.frames > 150*60 {
		g.lost = true
	}
	return nil
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .06)).Play()
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 26 + float64(i%8)})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	region := minInt(2, int(g.x/1100))
	bgs := []color.RGBA{{10, 18, 33, 255}, {23, 31, 53, 255}, {47, 24, 48, 255}}
	s.Fill(bgs[region])
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.frames)*2) * 5
	}
	for x := int(g.cam/40) * 40; x < int(g.cam)+W+80; x += 40 {
		sx := float32(float64(x) - g.cam + ox)
		f := float32(floorAt(float64(x)))
		vector.DrawFilledRect(s, sx, f, 42, H-f, color.RGBA{42, 67, 73, 255}, false)
	}
	for _, hx := range []float64{1080, 2050, 2790} {
		sx := float32(hx - g.cam + ox)
		if sx > -30 && sx < W+30 {
			for i := 0; i < 3; i++ {
				vector.DrawFilledRect(s, sx+float32(i*12-18), float32(floorAt(hx)-15), 8, 15, color.RGBA{245, 90, 90, 255}, false)
			}
		}
	}
	for rx := range g.relics {
		sx := float32(float64(rx) - g.cam)
		if sx > -30 && sx < W+30 {
			vector.DrawFilledCircle(s, sx, float32(floorAt(float64(rx))-60), 14, color.RGBA{245, 190, 68, 255}, false)
		}
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x-g.cam+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	px, py := float32(g.x-g.cam+ox), float32(g.y)
	pc := color.RGBA{231, 91, 77, 255}
	if g.flash%4 < 2 && g.flash > 0 {
		pc = color.RGBA{255, 255, 255, 255}
	}
	bob := float32(0)
	if math.Abs(g.vx) > .3 && g.onGround {
		bob = float32(math.Sin(float64(g.frames)*.3) * 3)
	}
	vector.DrawFilledRect(s, px-12, py+bob, 24, 36-bob, pc, false)
	vector.DrawFilledRect(s, 0, 0, W, 80, color.RGBA{5, 11, 24, 235}, false)
	g.drawHUD(s, region)
	g.drawEffectBadge(s)
	ebitenutil.DebugPrintAt(s, g.message, 25, 46)
	for i := 0; i < 8; i++ {
		c := color.RGBA{39, 48, 64, 255}
		if g.revealed[i] {
			c = color.RGBA{84, 151, 145, 255}
		}
		vector.DrawFilledRect(s, float32(40+i*50), 90, 42, 20, c, false)
	}
	labels := []string{"LEFT", "JUMP", "DASH", "RIGHT"}
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*120+5), 650, 110, 55, color.RGBA{45, 78, 113, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*120+35, 675)
	}
	if g.won {
		overlay(s, "THE DEPTHS MAPPED!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(s, "EXPLORATION FAILED\n\nTAP / ENTER TO RETRY")
	}
}
func (g *game) drawHUD(s *ebiten.Image, region int) {
	label := fmt.Sprintf("REGION %d/3 WORLD %04d ROOMS %d/8 RELICS %d DASH %v WINGS %v BEST %.1f", region+1, int(g.x), len(g.revealed), len(g.relics), g.dash, g.highJump, float64(g.bestFrames)/60)
	if face, err := uilab.Face("en", 13); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(10, 6)
		text.Draw(s, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(s, label, 10, 18)
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
	op.GeoM.Translate(W-34, 52)
	s.DrawImage(fx, op)
}
func controls() (bool, bool, bool, bool) {
	left := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyZ)
	dash := inpututil.IsKeyJustPressed(ebiten.KeyX)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 640 {
			left = x < 120
			jump = x >= 120 && x < 240
			dash = x >= 240 && x < 360
			right = x >= 360
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 640 {
			left = x < 120
			jump = x >= 120 && x < 240
			dash = x >= 240 && x < 360
			right = x >= 360
		}
	}
	return left, right, jump, dash
}
func press() (int, int, bool) {
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
func retry() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 45, 270, 390, 150, color.RGBA{4, 10, 24, 245}, false)
	ebitenutil.DebugPrintAt(s, t, 110, 330)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Depths")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
