package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/uilab"
	"image/color"
)

type game struct{ ja bool }

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.ja = !g.ja
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 15, 38, 255})
	f, _ := uilab.Face(map[bool]string{true: "ja", false: "en"}[g.ja], 20)
	words := "make small games together"
	if g.ja {
		words = "小さなゲームを一緒に作ろう"
	}
	lines := uilab.Wrap(words, 10, g.ja)
	ebitenutil.DebugPrintAt(s, "UI LAB 02 · MEASURE / ALIGN / WRAP", 65, 30)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: JA / EN", 145, 55)
	vector.StrokeRect(s, 80, 180, 320, 260, 2, color.RGBA{75, 180, 230, 255}, true)
	for i, line := range lines {
		w := text.Advance(line, f)
		op := &text.DrawOptions{}
		op.GeoM.Translate(240-w/2, float64(230+i*45))
		text.Draw(s, line, f, op)
	}
	ebitenutil.DebugPrintAt(s, "Measure then align; select CJK-safe line breaks before drawing.", 54, 560)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
