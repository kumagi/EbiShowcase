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

type choice struct {
	text string
	next int
}
type scene struct {
	speaker, text, expression, enter string
	choices                          []choice
}

var script = []scene{{"MIO", "The lantern festival starts tonight. Will you help me?", "smile", "fade", []choice{{"Of course.", 1}, {"What should I do?", 2}}}, {"MIO", "Then meet me by the river. I saved the brightest lantern for us.", "joy", "slide", []choice{{"I will be there.", 3}, {"Let's invite everyone.", 4}}}, {"MIO", "Write one wish on each lantern, then let the current carry it.", "explain", "pop", []choice{{"A wish for us.", 1}, {"A wish for the town.", 4}}}, {"MIO", "You came! The river looks like a road of stars.", "blush", "fade", []choice{{"END: promise", 5}}}, {"MIO", "Good idea. A shared light can still make a special memory.", "laugh", "shake", []choice{{"END: friendship", 5}}}, {"", "Story complete. Your choices selected a route.", "", "fade", nil}}

type game struct {
	node   int
	shown  int
	fade   float64
	frames int
	ended  bool
}

func newGame() *game { return &game{fade: 0} }
func (g *game) Update() error {
	if g.ended {
		if retry() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.fade < 1 {
		g.fade += .04
	}
	s := script[g.node]
	if g.shown < len([]rune(s.text)) && g.frames%2 == 0 {
		g.shown++
	}
	pick := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		pick = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		pick = 1
	}
	if _, y, ok := press(); ok && y >= 560 {
		pick = min(1, (y-560)/65)
	}
	if pick >= 0 && pick < len(s.choices) && g.shown >= len([]rune(s.text)) {
		g.node = s.choices[pick].next
		g.shown = 0
		g.fade = 0
		g.frames = 0
		if len(script[g.node].choices) == 0 {
			g.ended = true
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.shown = len([]rune(s.text))
	}
	return nil
}
func (g *game) Draw(scr *ebiten.Image) {
	scr.Fill(color.RGBA{24, 26, 52, 255})
	s := script[g.node]
	vector.DrawFilledCircle(scr, 100, 120, 90, color.RGBA{255, 185, 112, uint8(110 * g.fade)}, false)
	vector.DrawFilledCircle(scr, 390, 180, 130, color.RGBA{91, 105, 176, uint8(100 * g.fade)}, false)
	x := float32(240)
	if s.enter == "slide" {
		x = 480 - float32(240*g.fade)
	}
	c := color.RGBA{230, 118, 142, uint8(255 * g.fade)}
	vector.DrawFilledRect(scr, x-75, 160, 150, 260, c, false)
	vector.DrawFilledCircle(scr, x, 145, 74, color.RGBA{248, 205, 172, uint8(255 * g.fade)}, false)
	ebitenutil.DebugPrintAt(scr, "MIO  "+s.expression, int(x)-42, 280)
	vector.DrawFilledRect(scr, 20, 445, 440, 250, color.RGBA{6, 12, 28, 235}, false)
	vector.StrokeRect(scr, 20, 445, 440, 250, 3, color.RGBA{244, 188, 78, 255}, false)
	ebitenutil.DebugPrintAt(scr, s.speaker, 42, 470)
	drawWrapped(scr, string([]rune(s.text)[:min(g.shown, len([]rune(s.text)))]), 42, 500, 48)
	for i, ch := range s.choices {
		y := 560 + i*65
		vector.DrawFilledRect(scr, 40, float32(y), 400, 52, color.RGBA{54, 72, 112, 255}, false)
		ebitenutil.DebugPrintAt(scr, fmt.Sprintf("[%d] %s", i+1, ch.text), 58, y+20)
	}
	if g.shown < len([]rune(s.text)) {
		ebitenutil.DebugPrintAt(scr, "SPACE / tap choice after text finishes", 102, 680)
	}
	if g.ended {
		overlay(scr, "ROUTE COMPLETE\n\nTAP / ENTER TO RESTART")
	}
}
func drawWrapped(s *ebiten.Image, t string, x, y, w int) {
	r := []rune(t)
	for len(r) > 0 {
		n := min(w, len(r))
		ebitenutil.DebugPrintAt(s, string(r[:n]), x, y)
		r = r[n:]
		y += 22
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
func retry() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 50, 260, 380, 150, color.RGBA{4, 10, 24, 245}, false)
	ebitenutil.DebugPrintAt(s, t, 130, 320)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Dialogue Stage")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
