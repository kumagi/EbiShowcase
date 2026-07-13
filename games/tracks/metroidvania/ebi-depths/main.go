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
	W      = 480
	H      = 720
	worldW = 3200
	ground = 520
)

type game struct {
	x, y, vx, vy, cam         float64
	onGround, dash, won, lost bool
	revealed                  map[int]bool
	relics                    map[int]bool
	frames                    int
	message                   string
}

func newGame() *game {
	return &game{x: 80, y: ground - 36, revealed: map[int]bool{}, relics: map[int]bool{720: true, 1450: true, 2350: true, 2980: true}, message: "Explore the huge world. The map reveals one room at a time."}
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
			*g = *newGame()
		}
		return nil
	}
	g.frames++
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
	}
	if jump && g.onGround {
		g.vy = -9
		g.onGround = false
	}
	g.vy += .45
	g.x += g.vx
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
			delete(g.relics, rx)
			if rx == 1450 {
				g.dash = true
				g.message = "Dash crest found! Hold a direction and press X."
			} else {
				g.message = "Map fragment found. More of the world is recorded."
			}
		}
	}
	if len(g.relics) == 0 && g.x > 3100 {
		g.won = true
	}
	if g.y > 690 || g.frames > 150*60 {
		g.lost = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 18, 33, 255})
	for x := int(g.cam/40) * 40; x < int(g.cam)+W+80; x += 40 {
		sx := float32(float64(x) - g.cam)
		f := float32(floorAt(float64(x)))
		vector.DrawFilledRect(s, sx, f, 42, H-f, color.RGBA{42, 67, 73, 255}, false)
	}
	for rx := range g.relics {
		sx := float32(float64(rx) - g.cam)
		if sx > -30 && sx < W+30 {
			vector.DrawFilledCircle(s, sx, float32(floorAt(float64(rx))-60), 14, color.RGBA{245, 190, 68, 255}, false)
		}
	}
	px, py := float32(g.x-g.cam), float32(g.y)
	vector.DrawFilledRect(s, px-12, py, 24, 36, color.RGBA{231, 91, 77, 255}, false)
	vector.DrawFilledRect(s, 0, 0, W, 80, color.RGBA{5, 11, 24, 235}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("WORLD %04d/%d  ROOMS %d/8  RELICS LEFT %d  DASH %v", int(g.x), worldW, len(g.revealed), len(g.relics), g.dash), 38, 18)
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
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Depths")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
