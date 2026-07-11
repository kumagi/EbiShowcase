package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strings"
)

const width, height = 480, 720

type game struct {
	history               []string
	age, projectile, hits int
	message               string
	clear                 bool
}

func (g *game) Update() error {
	if g.clear {
		if any() {
			*g = game{}
		}
		return nil
	}
	token := ""
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		token = "D"
	}
	if (inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)) && (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)) {
		token = "DR"
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		token = "R"
	}
	punch := inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if x, y, ok := press(); ok && y > 560 {
		zone := min(3, x/(width/4))
		if zone == 3 {
			punch = true
		} else {
			token = []string{"D", "DR", "R"}[zone]
		}
	}
	if token != "" {
		g.history = append(g.history, token)
		if len(g.history) > 6 {
			g.history = g.history[len(g.history)-6:]
		}
		g.age = 30
	}
	if g.age > 0 {
		g.age--
	} else {
		g.history = nil
	}
	if punch {
		if ends(g.history, []string{"D", "DR", "R"}) {
			g.projectile = 35
			g.message = "QUARTER CIRCLE + PUNCH! EBI WAVE!"
			g.history = nil
		} else {
			g.message = "Command missed. Try D, DR, R, Punch."
			g.history = nil
		}
	}
	if g.projectile > 0 {
		g.projectile--
		if g.projectile == 5 {
			g.hits++
			if g.hits >= 5 {
				g.clear = true
			}
		}
	}
	return nil
}
func ends(h, w []string) bool {
	if len(h) < len(w) {
		return false
	}
	start := len(h) - len(w)
	for i := range w {
		if h[start+i] != w[i] {
			return false
		}
	}
	return true
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{19, 28, 44, 255})
	vector.DrawFilledRect(s, 0, 590, 480, 130, color.RGBA{55, 66, 78, 255}, false)
	vector.DrawFilledCircle(s, 105, 490, 25, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 85, 515, 40, 70, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledCircle(s, 390, 490, 25, color.RGBA{240, 75, 91, 255}, false)
	if g.projectile > 0 {
		x := float32(125 + (35-g.projectile)*7)
		vector.DrawFilledCircle(s, x, 530, 18, color.RGBA{75, 170, 255, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI WAVES %d/5", g.hits), 190, 40)
	ebitenutil.DebugPrintAt(s, "COMMAND: DOWN > DOWN-RIGHT > RIGHT + PUNCH", 90, 80)
	hist := "EMPTY"
	if len(g.history) > 0 {
		hist = strings.Join(g.history, " > ")
	}
	ebitenutil.DebugPrintAt(s, "INPUT HISTORY: "+hist, 110, 125)
	ebitenutil.DebugPrintAt(s, g.message, 70, 165)
	labels := []string{"DOWN", "DOWN-RIGHT", "RIGHT", "PUNCH"}
	for i, l := range labels {
		x := float32(5 + i*119)
		vector.DrawFilledRect(s, x, 615, 113, 70, color.RGBA{60 + uint8(i)*25, 135, 210, 255}, false)
		ebitenutil.DebugPrintAt(s, l, int(x)+18, 648)
	}
	if g.clear {
		overlay(s, "COMMAND MASTERED!\n\nTAP / SPACE TO RESET")
	}
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
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Command Fighter — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
