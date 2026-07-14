package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"image/color"
)

type game struct {
	f    cameralab.Frame
	room int
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.room = (g.room + 1) % 3
		g.f.Start(36)
	}
	g.f.Tick()
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill([]color.RGBA{{25, 48, 80, 255}, {58, 42, 88, 255}, {35, 72, 65, 255}}[g.room])
	x, y, w, h := cameralab.SafeRect(480, 720, 24)
	vector.StrokeRect(s, float32(x), float32(y), float32(w), float32(h), 2, color.RGBA{80, 220, 190, 255}, true)
	ebitenutil.DebugPrintAt(s, "CAMERA LAB 05 · ROOM FRAME", 105, 55)
	ebitenutil.DebugPrintAt(s, "SAFE UI AREA", 45, 100)
	ebitenutil.DebugPrintAt(s, "ROOM "+string(rune('A'+g.room)), 210, 350)
	bar := g.f.Letterbox(720)
	if bar > 0 {
		vector.DrawFilledRect(s, 0, 0, 480, float32(bar), color.Black, true)
		vector.DrawFilledRect(s, 0, float32(720-bar), 480, float32(bar), color.Black, true)
	}
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: room transition + letterbox", 105, 650)
}
func (g *game) Layout(_, _ int) (int, int) { return 480, 720 }
func main() {
	ebiten.SetWindowSize(480, 720)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
