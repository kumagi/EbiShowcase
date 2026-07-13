package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	W       = 480
	H       = 720
	brave   = 1
	kind    = 2
	curious = 4
)

type choice struct {
	text       string
	next, flag int
}
type scene struct {
	chapter                             int
	speaker, text, expression, entrance string
	choices                             []choice
}

var story = []scene{
	{1, "MIO", "The Moon Lantern has vanished before tonight's festival! Will you search with me?", "worried", "fade", []choice{{"Run straight to the old pier.", 1, brave}, {"Ask if anyone needs help first.", 2, kind}}},
	{1, "MIO", "Bold plan! A silver ribbon points from the pier toward the market.", "surprised", "slide", []choice{{"Follow the ribbon.", 3, 0}}},
	{1, "REN", "Thank you. I saw blue sparks heading toward the market clock.", "smile", "pop", []choice{{"Remember the clue.", 3, 0}}},
	{2, "REN", "The clock is stopped at 8:15. A tiny door hides behind its pendulum.", "thinking", "slide", []choice{{"Open the secret door.", 4, curious}, {"Find the clock keeper first.", 5, kind}}},
	{2, "MIO", "Inside is a map signed by the Star Keeper. It marks the hill observatory.", "joy", "pop", []choice{{"Race up the hill.", 6, brave}}},
	{2, "KEEPER", "You asked before touching the clock. Take this observatory key—and my thanks.", "smile", "fade", []choice{{"Promise to return it.", 6, kind}}},
	{3, "MIO", "At the observatory, the lantern floats inside a ring of storm clouds!", "worried", "shake", []choice{{"Jump through and catch it.", 7, brave}, {"Study the cloud pattern.", 8, curious}}},
	{3, "REN", "You caught it! The storm is fading, but the lantern needs one final wish.", "joy", "bounce", []choice{{"Wish for everyone's festival.", -1, kind}, {"Wish for another adventure.", -1, curious}}},
	{3, "MIO", "The clouds repeat in a three-beat rhythm. We can walk through the quiet beat.", "thinking", "pulse", []choice{{"Count together and enter.", -1, curious}, {"Guide everyone through safely.", -1, kind}}},
}
var endingNames = []string{"MOON PROMISE", "HARBOR HERO", "TOWN LIGHT", "STAR SEEKER"}

type particle struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}
type game struct {
	node, shown, frames, flags, ending, endingMask, shake, flash int
	enter                                                        float64
	ended                                                        bool
	particles                                                    []particle
	rng                                                          *rand.Rand
}

var collectedEndings int

func newGame() *game {
	return &game{endingMask: collectedEndings, enter: 0, rng: rand.New(rand.NewSource(81))}
}
func (g *game) Update() error {
	if g.ended {
		if retry() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.enter < 1 {
		g.enter += .045
	}
	s := story[g.node]
	runes := []rune(s.text)
	if g.shown < len(runes) && g.frames%2 == 0 {
		g.shown++
	}
	pick := -1
	pressed := false
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		pick = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		pick = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		pressed = true
	}
	if _, y, ok := press(); ok {
		pressed = true
		if y >= 555 {
			pick = min(1, (y-555)/62)
		}
	}
	if pressed && g.shown < len(runes) {
		g.shown = len(runes)
		pick = -1
	} else if pressed && len(s.choices) == 1 {
		pick = 0
	}
	if pick >= 0 && pick < len(s.choices) && g.shown >= len(runes) {
		g.selectChoice(s.choices[pick])
	}
	g.updateParticles()
	if g.shake > 0 {
		g.shake--
	}
	if g.flash > 0 {
		g.flash--
	}
	return nil
}
func (g *game) selectChoice(c choice) {
	g.flags |= c.flag
	g.flash = 5
	for i := 0; i < 18; i++ {
		a := float64(i) * .65
		g.particles = append(g.particles, particle{240, 580, math.Cos(a) * 2.5, math.Sin(a)*2.5 - 1, 32, color.RGBA{255, 194, 76, 255}})
	}
	if c.next < 0 {
		g.finish()
		return
	}
	g.node = c.next
	g.shown = 0
	g.frames = 0
	g.enter = 0
	if story[g.node].entrance == "shake" {
		g.shake = 18
	}
}
func (g *game) finish() {
	switch {
	case g.flags == (brave | kind | curious):
		g.ending = 0
	case g.flags&brave != 0 && g.flags&kind == 0:
		g.ending = 1
	case g.flags&kind != 0 && g.flags&curious == 0:
		g.ending = 2
	default:
		g.ending = 3
	}
	collectedEndings |= 1 << g.ending
	g.endingMask = collectedEndings
	g.ended = true
	g.flash = 18
}
func (g *game) updateParticles() {
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .04
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}
}
func (g *game) Draw(screen *ebiten.Image) {
	s := story[g.node]
	background(screen, s.chapter, g.frames)
	dx, dy := 0, 0
	if g.shake > 0 {
		dx = g.rng.Intn(9) - 4
		dy = g.rng.Intn(9) - 4
	}
	stage := ebiten.NewImage(W, H)
	drawPortrait(stage, s, g.enter, g.frames)
	for _, p := range g.particles {
		vector.DrawFilledCircle(stage, float32(p.x), float32(p.y), 3, p.c, false)
	}
	vector.DrawFilledRect(stage, 18, 438, 444, 264, color.RGBA{5, 12, 29, 238}, false)
	vector.StrokeRect(stage, 18, 438, 444, 264, 3, color.RGBA{244, 188, 78, 255}, false)
	ebitenutil.DebugPrintAt(stage, fmt.Sprintf("CHAPTER %d / 3", s.chapter), 350, 454)
	ebitenutil.DebugPrintAt(stage, s.speaker, 40, 466)
	drawWrapped(stage, string([]rune(s.text)[:min(g.shown, len([]rune(s.text)))]), 40, 494, 49)
	if g.shown >= len([]rune(s.text)) {
		for i, ch := range s.choices {
			y := 555 + i*62
			vector.DrawFilledRect(stage, 36, float32(y), 408, 50, color.RGBA{54, 72 + uint8(i*10), 112, 255}, false)
			vector.StrokeRect(stage, 36, float32(y), 408, 50, 2, color.RGBA{155, 190, 230, 255}, false)
			ebitenutil.DebugPrintAt(stage, fmt.Sprintf("[%d] %s", i+1, ch.text), 52, y+19)
		}
	} else {
		ebitenutil.DebugPrintAt(stage, "SPACE / TAP: show the whole line", 123, 672)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	screen.DrawImage(stage, op)
	if g.flash > 0 {
		vector.DrawFilledRect(screen, 0, 0, W, H, color.RGBA{255, 255, 255, 55}, false)
	}
	if g.ended {
		g.drawEnding(screen)
	}
}
func background(s *ebiten.Image, chapter, frames int) {
	colors := []color.RGBA{{32, 39, 78, 255}, {42, 29, 68, 255}, {18, 31, 61, 255}}
	s.Fill(colors[chapter-1])
	for i := 0; i < 22; i++ {
		x := float32((i*73 + chapter*31) % W)
		y := float32((i*41 + frames/5) % 410)
		a := uint8(100 + (i*7)%120)
		vector.DrawFilledCircle(s, x, y, float32(1+i%3), color.RGBA{255, 230, 170, a}, false)
	}
	if chapter == 1 {
		vector.DrawFilledRect(s, 0, 350, W, 88, color.RGBA{22, 65, 87, 255}, false)
	} else if chapter == 2 {
		vector.DrawFilledCircle(s, 240, 230, 120, color.RGBA{67, 57, 104, 255}, false)
	} else {
		vector.StrokeCircle(s, 240, 230, 125, 8, color.RGBA{141, 187, 229, 130}, false)
	}
}
func drawPortrait(s *ebiten.Image, sc scene, t float64, frames int) {
	x := 240.0
	if sc.entrance == "slide" {
		x = 540 - 300*t
	} else if sc.speaker == "REN" || sc.speaker == "KEEPER" {
		x = 120
	}
	y := 260.0 + math.Sin(float64(frames)*.08)*3
	scale := 1.0
	if sc.entrance == "pop" {
		scale = .7 + .3*ease(t)
	}
	if sc.entrance == "bounce" {
		y -= math.Sin(t*math.Pi) * 35
	}
	alpha := uint8(min(255, int(t*255)))
	body := color.RGBA{225, 104, 143, alpha}
	if sc.speaker == "REN" {
		body = color.RGBA{86, 155, 190, alpha}
	} else if sc.speaker == "KEEPER" {
		body = color.RGBA{151, 116, 190, alpha}
	}
	vector.DrawFilledRect(s, float32(x-72*scale), float32(y-10), float32(144*scale), float32(190*scale), body, false)
	vector.DrawFilledCircle(s, float32(x), float32(y-45), float32(66*scale), color.RGBA{247, 204, 170, alpha}, false)
	eyeY := float32(y - 54)
	blink := frames%170 > 160
	if blink {
		vector.StrokeLine(s, float32(x-28), eyeY, float32(x-16), eyeY, 3, color.RGBA{35, 31, 50, alpha}, false)
		vector.StrokeLine(s, float32(x+16), eyeY, float32(x+28), eyeY, 3, color.RGBA{35, 31, 50, alpha}, false)
	} else {
		vector.DrawFilledCircle(s, float32(x-22), eyeY, 4, color.RGBA{35, 31, 50, alpha}, false)
		vector.DrawFilledCircle(s, float32(x+22), eyeY, 4, color.RGBA{35, 31, 50, alpha}, false)
	}
	mouth := float32(8)
	if sc.expression == "worried" {
		mouth = -6
	}
	vector.StrokeLine(s, float32(x-10), float32(y-22), float32(x+10), float32(y-22+float64(mouth)), 3, color.RGBA{100, 43, 61, alpha}, false)
	ebitenutil.DebugPrintAt(s, sc.speaker+" / "+sc.expression, int(x)-50, int(y+120))
}
func (g *game) drawEnding(s *ebiten.Image) {
	vector.DrawFilledRect(s, 25, 175, 430, 370, color.RGBA{4, 10, 25, 248}, false)
	vector.StrokeRect(s, 25, 175, 430, 370, 4, color.RGBA{255, 198, 78, 255}, false)
	ebitenutil.DebugPrintAt(s, "ENDING UNLOCKED", 174, 210)
	ebitenutil.DebugPrintAt(s, endingNames[g.ending], 175, 255)
	ebitenutil.DebugPrintAt(s, "ENDING GALLERY", 178, 320)
	for i, name := range endingNames {
		mark := "LOCKED"
		if g.endingMask&(1<<i) != 0 {
			mark = "FOUND"
		}
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d  %-14s %s", i+1, name, mark), 95, 355+i*28)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d / 4 endings found", bitCount(g.endingMask)), 169, 480)
	ebitenutil.DebugPrintAt(s, "TAP / ENTER: TRY ANOTHER ROUTE", 116, 520)
}
func drawWrapped(s *ebiten.Image, t string, x, y, w int) {
	r := []rune(t)
	for len(r) > 0 {
		n := min(w, len(r))
		ebitenutil.DebugPrintAt(s, string(r[:n]), x, y)
		r = r[n:]
		y += 20
	}
}
func ease(t float64) float64 { return 1 - math.Pow(1-t, 3) }
func bitCount(v int) int {
	n := 0
	for v > 0 {
		n += v & 1
		v >>= 1
	}
	return n
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
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Moon Lantern Stories")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
