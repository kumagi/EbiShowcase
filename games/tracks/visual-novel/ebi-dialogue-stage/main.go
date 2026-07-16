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
	"github.com/kumagi/EbiShowcase/internal/mobileart"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
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
	audio                                                        *audio.Context
	gate                                                         audiolab.Gate
	pulse                                                        *shaderlab.Pulse
	cam                                                          cameralab.State
	badge                                                        *ebiten.Image
}

var collectedEndings int

func newGame() *game {
	mobileart.Preload()
	g := &game{endingMask: collectedEndings, enter: 0}
	g.audio = audiolab.Context()
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: W / 2, Y: H / 2}, ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{255, 198, 78, 255})
	return g
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
	g.play(620)
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
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .055)).Play()
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
	backgrounds := [...]string{"visual-novel-harbor", "visual-novel-clock", "visual-novel-observatory"}
	mobileart.DrawCover(screen, backgrounds[s.chapter-1], 0, 0, W, H)
	// Chapter grading changes the mood without changing the shared location.
	// It is a Draw-only interpretation of the story state.
	chapterTint(screen, s.chapter)
	dx, dy := 0, 0
	if g.shake > 0 {
		// Camera shake is derived only from Update-owned state. Draw never
		// consumes randomness, so skipped Draw calls cannot change the story.
		dx = (g.frames*17%9 - 4) * min(1, g.shake)
		dy = (g.frames*29%7 - 3) * min(1, g.shake)
	}
	stage := ebiten.NewImage(W, H)
	g.drawTitle(stage)
	g.drawEffectBadge(stage)
	drawPortrait(stage, s, g.enter, g.frames)
	for _, p := range g.particles {
		vector.DrawFilledCircle(stage, float32(p.x), float32(p.y), 3, p.c, false)
	}
	// Contemporary mobile dialogue UI: a glass panel, bright speaker tab and
	// generous touch targets. The artwork remains visible above the panel.
	vector.DrawFilledRect(stage, 16, 442, 448, 260, color.RGBA{3, 11, 27, 228}, false)
	vector.StrokeRect(stage, 16, 442, 448, 260, 2, color.RGBA{115, 220, 225, 210}, false)
	vector.DrawFilledRect(stage, 30, 426, 148, 42, color.RGBA{12, 76, 91, 245}, false)
	vector.StrokeRect(stage, 30, 426, 148, 42, 2, color.RGBA{246, 206, 121, 255}, false)
	drawLabel(stage, s.speaker, 46, 436, 18, color.White)
	drawLabel(stage, fmt.Sprintf("CHAPTER %d / 3", s.chapter), 344, 452, 11, color.RGBA{176, 229, 234, 255})
	drawWrappedFace(stage, string([]rune(s.text)[:min(g.shown, len([]rune(s.text)))]), 38, 486, 405, 17)
	if g.shown >= len([]rune(s.text)) {
		for i, ch := range s.choices {
			y := 555 + i*62
			vector.DrawFilledRect(stage, 34, float32(y), 412, 50, color.RGBA{18, 67 + uint8(i*8), 88, 245}, false)
			vector.StrokeRect(stage, 34, float32(y), 412, 50, 2, color.RGBA{120, 224, 225, 230}, false)
			vector.DrawFilledCircle(stage, 56, float32(y+25), 14, color.RGBA{237, 190, 100, 255}, true)
			drawLabel(stage, fmt.Sprintf("%d", i+1), 52, float64(y+16), 13, color.RGBA{18, 36, 48, 255})
			drawLabel(stage, ch.text, 82, float64(y+15), 14, color.White)
		}
	} else {
		drawLabel(stage, "TAP TO REVEAL  ◆", 304, 672, 10, color.RGBA{157, 220, 224, 255})
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
func (g *game) drawTitle(s *ebiten.Image) {
	if face, err := uilab.Face("en", 14); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(16, 14)
		text.Draw(s, "EBI DIALOGUE STAGE", face, op)
		return
	}
	ebitenutil.DebugPrintAt(s, "EBI DIALOGUE STAGE", 16, 26)
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
	op.GeoM.Translate(W-34, 10)
	s.DrawImage(fx, op)
}
func chapterTint(s *ebiten.Image, chapter int) {
	tints := []color.RGBA{{8, 33, 57, 34}, {58, 23, 72, 38}, {6, 28, 65, 18}}
	vector.DrawFilledRect(s, 0, 0, W, H, tints[chapter-1], false)
	vector.DrawFilledRect(s, 0, 0, W, 54, color.RGBA{2, 10, 25, 145}, false)
}
func drawPortrait(s *ebiten.Image, sc scene, t float64, frames int) {
	x := 36.0
	if sc.entrance == "slide" {
		x = 310 - 274*t
	}
	y := 48.0 + math.Sin(float64(frames)*.08)*2
	scale := 1.0
	if sc.entrance == "pop" {
		scale = .7 + .3*ease(t)
	}
	if sc.entrance == "bounce" {
		y -= math.Sin(t*math.Pi) * 22
	}
	name := "navigator"
	if sc.speaker == "MIO" {
		switch sc.expression {
		case "worried":
			name = "navigator-worried"
		case "surprised", "thinking":
			name = "navigator-surprised"
		case "joy", "smile":
			name = "navigator-joy"
		}
	} else if sc.speaker == "REN" {
		name = "researcher"
	} else if sc.speaker == "KEEPER" {
		name = "keeper"
	}
	w, h := 410*scale, 450*scale
	mobileart.DrawContainAlpha(s, name, x+(410-w)/2, y+(450-h), w, h, sc.speaker == "REN", float32(min(1, t)))
	// Non-MIO cast members currently use a single premium pose, so a restrained
	// thought mark supports their state without pretending the portrait changed.
	if sc.speaker != "MIO" && (sc.expression == "worried" || sc.expression == "thinking") {
		vector.DrawFilledCircle(s, 412, 93, 27, color.RGBA{5, 22, 42, 210}, true)
		drawLabel(s, "…", 400, 76, 24, color.RGBA{182, 234, 238, 255})
	}
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
func drawLabel(s *ebiten.Image, value string, x, y float64, size float64, clr color.Color) {
	face, err := uilab.Face("en", size)
	if err != nil {
		ebitenutil.DebugPrintAt(s, value, int(x), int(y))
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(s, value, face, op)
}
func drawWrappedFace(s *ebiten.Image, value string, x, y, maxWidth, size float64) {
	face, err := uilab.Face("en", size)
	if err != nil {
		ebitenutil.DebugPrintAt(s, value, int(x), int(y))
		return
	}
	words, line := []rune(value), []rune{}
	for len(words) > 0 {
		line = append(line, words[0])
		words = words[1:]
		width, _ := text.Measure(string(line), face, 0)
		if width > maxWidth && len(line) > 1 {
			last := line[len(line)-1]
			line = line[:len(line)-1]
			drawLabel(s, string(line), x, y, size, color.White)
			line = []rune{last}
			y += size + 8
		}
	}
	if len(line) > 0 {
		drawLabel(s, string(line), x, y, size, color.White)
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
