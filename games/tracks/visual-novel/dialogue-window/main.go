package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	W = 480
	H = 720
)

var lines = []struct{ speaker, text string }{{"MIO", "The festival begins when the moon lantern shines."}, {"REN", "Then let us search the harbor before sunset."}, {"MIO", "Good! Our story can begin now."}}

type game struct {
	line, shown, frames int
	clear               bool
}

func (g *game) press() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Update() error {
	if g.clear {
		if g.press() {
			*g = game{}
		}
		return nil
	}
	g.frames++
	r := []rune(lines[g.line].text)
	if g.shown < len(r) && g.frames%2 == 0 {
		g.shown++
	}
	if g.press() {
		if g.shown < len(r) {
			g.shown = len(r)
		} else {
			g.line++
			g.shown = 0
			if g.line >= len(lines) {
				g.line = len(lines) - 1
				g.clear = true
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{28, 35, 72, 255})
	for i := 0; i < 25; i++ {
		vector.DrawFilledCircle(s, float32((i*71)%W), float32(40+(i*43)%390), 2, color.RGBA{255, 226, 165, 190}, false)
	}
	vector.DrawFilledRect(s, 18, 410, 444, 265, color.RGBA{5, 12, 29, 240}, false)
	vector.StrokeRect(s, 18, 410, 444, 265, 3, color.RGBA{244, 188, 78, 255}, false)
	v := lines[g.line]
	ebitenutil.DebugPrintAt(s, "TYPEWRITER WINDOW", 168, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("LINE %d/3", g.line+1), 384, 430)
	ebitenutil.DebugPrintAt(s, v.speaker, 40, 445)
	r := []rune(v.text)
	ebitenutil.DebugPrintAt(s, string(r[:min(g.shown, len(r))]), 40, 490)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: REVEAL, THEN ADVANCE", 102, 635)
	if g.clear {
		vector.DrawFilledRect(s, 45, 255, 390, 150, color.RGBA{4, 10, 24, 245}, false)
		ebitenutil.DebugPrintAt(s, "SCENE COMPLETE!", 174, 305)
		ebitenutil.DebugPrintAt(s, "TAP TO REPLAY", 188, 355)
	}
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
