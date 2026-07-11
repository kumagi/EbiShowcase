package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const width, height, tile = 480, 720, 48

type point struct{ x, y int }
type game struct {
	p, facing          point
	npcs               map[point]string
	talk               []string
	line               int
	quest, herb, clear bool
}

func newGame() *game {
	return &game{p: point{4, 10}, facing: point{0, -1}, npcs: map[point]string{{4, 3}: "elder", {2, 6}: "child", {7, 8}: "guard"}}
}
func (g *game) Update() error {
	action := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	x, y, tap := press()
	if tap && y < 180 {
		action = true
	}
	if len(g.talk) > 0 {
		if action || tap {
			g.line++
			if g.line >= len(g.talk) {
				g.talk = nil
				g.line = 0
			}
		}
		return nil
	}
	if g.clear {
		if action || tap {
			*g = *newGame()
		}
		return nil
	}
	d := point{}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		d.x = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		d.x = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		d.y = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		d.y = 1
	}
	if tap && !action {
		dx, dy := x-(g.p.x*tile+24), y-(74+g.p.y*tile+24)
		if abs(dx) > abs(dy) {
			if dx < 0 {
				d.x = -1
			} else {
				d.x = 1
			}
		} else {
			if dy < 0 {
				d.y = -1
			} else {
				d.y = 1
			}
		}
	}
	if d.x != 0 || d.y != 0 {
		g.facing = d
		n := point{g.p.x + d.x, g.p.y + d.y}
		if n.x > 0 && n.x < 9 && n.y > 0 && n.y < 12 {
			if _, blocked := g.npcs[n]; !blocked {
				g.p = n
			}
		}
	}
	if g.quest && !g.herb && g.p == (point{1, 1}) {
		g.herb = true
		g.talk = []string{"You found the shining herb!", "Take it back to the elder."}
	}
	if action {
		g.interact()
	}
	return nil
}
func (g *game) interact() {
	front := point{g.p.x + g.facing.x, g.p.y + g.facing.y}
	who, ok := g.npcs[front]
	if !ok {
		return
	}
	switch who {
	case "elder":
		if !g.quest {
			g.quest = true
			g.talk = []string{"Elder: A healing herb grows northwest.", "Please bring it back to me!"}
		} else if !g.herb {
			g.talk = []string{"Elder: The herb is in the northwest corner."}
		} else {
			g.talk = []string{"Elder: You found it! Thank you, hero!", "QUEST COMPLETE"}
			g.clear = true
		}
	case "child":
		g.talk = []string{"Child: Press Space or tap the top to talk!"}
	case "guard":
		g.talk = []string{"Guard: Flags let the same person remember events."}
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 28, 40, 255})
	oy := 74
	for y := 0; y < 12; y++ {
		for x := 0; x < 10; x++ {
			c := color.RGBA{91, 171, 94, 255}
			if x == 0 || x == 9 || y == 0 || y == 11 {
				c = color.RGBA{65, 91, 111, 255}
			}
			vector.DrawFilledRect(s, float32(x*tile), float32(oy+y*tile), tile, tile, c, false)
		}
	}
	for p, n := range g.npcs {
		c := color.RGBA{255, 211, 62, 255}
		if n == "elder" {
			c = color.RGBA{139, 86, 205, 255}
		}
		vector.DrawFilledCircle(s, float32(p.x*tile+24), float32(oy+p.y*tile+24), 15, c, false)
	}
	if g.quest && !g.herb {
		vector.DrawFilledCircle(s, 72, float32(oy+72), 10, color.RGBA{45, 225, 194, 255}, false)
	}
	px, py := float32(g.p.x*tile+24), float32(oy+g.p.y*tile+24)
	vector.DrawFilledCircle(s, px, py, 15, color.RGBA{240, 74, 90, 255}, false)
	vector.DrawFilledCircle(s, px+float32(g.facing.x*8), py+float32(g.facing.y*8), 3, color.White, false)
	status := "TALK TO THE PURPLE ELDER"
	if g.quest && !g.herb {
		status = "FIND THE HERB IN THE NORTHWEST"
	} else if g.herb && !g.clear {
		status = "RETURN TO THE ELDER"
	}
	ebitenutil.DebugPrintAt(s, status, 125, 30)
	ebitenutil.DebugPrintAt(s, "MOVE: ARROWS / TAP    TALK: SPACE / TOP TAP", 85, 690)
	if len(g.talk) > 0 {
		vector.DrawFilledRect(s, 25, 500, 430, 150, color.RGBA{5, 17, 35, 245}, false)
		vector.StrokeRect(s, 25, 500, 430, 150, 3, color.White, false)
		ebitenutil.DebugPrintAt(s, g.talk[g.line], 45, 535)
		ebitenutil.DebugPrintAt(s, "SPACE / TAP TO CONTINUE", 155, 615)
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
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Dialogue Flags — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
