package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const width, height = 480, 720
const (
	none = iota
	light
	heavy
)

type game struct {
	move, frame, buffer, bufferLife, combo, completed int
	message                                           string
	clear                                             bool
}

func (g *game) Update() error {
	if g.clear {
		if any() {
			*g = game{}
		}
		return nil
	}
	in := none
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		in = light
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		in = heavy
	}
	if x, y, ok := press(); ok && y > 500 {
		if x < width/2 {
			in = light
		} else {
			in = heavy
		}
	}
	if in != none {
		g.buffer = in
		g.bufferLife = 8
	}
	if g.bufferLife > 0 {
		g.bufferLife--
	} else {
		g.buffer = none
	}
	if g.move == none && g.buffer != none {
		g.start(g.buffer)
	} else if g.move != none {
		g.frame++
		active := g.frame == 8
		cancel := g.frame >= 8 && g.frame <= 13
		if active {
			expected := light
			if g.combo == 2 {
				expected = heavy
			}
			if g.move == expected {
				g.combo++
				g.message = fmt.Sprintf("HIT %d! CANCEL WINDOW OPEN", g.combo)
			} else {
				g.combo = 0
				g.message = "Wrong move — combo reset"
			}
		}
		if cancel && g.buffer != none {
			expected := light
			if g.combo >= 2 {
				expected = heavy
			}
			if g.buffer == expected {
				g.start(g.buffer)
			}
		}
		limit := 20
		if g.move == heavy {
			limit = 28
		}
		if g.frame >= limit {
			g.move = none
			g.frame = 0
			if g.combo == 3 {
				g.completed++
				g.message = "LIGHT > LIGHT > HEAVY COMPLETE!"
				g.combo = 0
				if g.completed >= 3 {
					g.clear = true
				}
			} else if g.combo > 0 {
				g.message = "Too late — combo dropped"
				g.combo = 0
			}
		}
	}
	return nil
}
func (g *game) start(m int) { g.move = m; g.frame = 1; g.buffer = none; g.bufferLife = 0 }
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{19, 28, 44, 255})
	vector.DrawFilledRect(s, 0, 610, 480, 110, color.RGBA{55, 66, 78, 255}, false)
	vector.DrawFilledCircle(s, 170, 500, 24, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 150, 524, 40, 80, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledCircle(s, 330, 500, 24, color.RGBA{240, 75, 91, 255}, false)
	if g.move != none && g.frame >= 7 && g.frame <= 11 {
		reach := float32(80)
		c := color.RGBA{255, 210, 62, 255}
		if g.move == heavy {
			reach = 120
			c = color.RGBA{255, 105, 65, 255}
		}
		vector.DrawFilledRect(s, 190, 525, reach, 30, c, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("3-HIT COMBOS %d/3   CURRENT HITS %d", g.completed, g.combo), 110, 45)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MOVE %s  FRAME %02d   BUFFER %s (%d)", name(g.move), g.frame, name(g.buffer), g.bufferLife), 85, 80)
	ebitenutil.DebugPrintAt(s, g.message, 95, 125)
	vector.DrawFilledRect(s, 30, 500, 0, 0, color.Black, false)
	vector.DrawFilledRect(s, 25, 625, 205, 65, color.RGBA{75, 150, 240, 255}, false)
	vector.DrawFilledRect(s, 250, 625, 205, 65, color.RGBA{240, 100, 70, 255}, false)
	ebitenutil.DebugPrintAt(s, "J / Z  LIGHT", 80, 655)
	ebitenutil.DebugPrintAt(s, "K / X  HEAVY", 305, 655)
	if g.clear {
		overlay(s, "THREE COMBOS COMPLETE!\n\nTAP / SPACE TO RESET")
	}
}
func name(v int) string { return []string{"NONE", "LIGHT", "HEAVY"}[v] }
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
	ebitenutil.DebugPrintAt(s, msg, 115, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Combo and Cancel — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
