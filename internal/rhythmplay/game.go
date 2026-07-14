// Package rhythmplay renders the shared classroom rhythm stage.
// All rules and judgments remain in rhythmcore for deterministic tests.
// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package rhythmplay

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/rhythmcore"
)

const W, H = 480, 720

type Song struct {
	Name       string
	Tone       float64
	Easy, Hard rhythmcore.Chart
}

type Config struct {
	Title, Subtitle string
	Songs           []Song
	Difficulty      bool
}

type particle struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}

type Game struct {
	cfg                      Config
	song, difficulty         int
	menu, finished           bool
	session                  *rhythmcore.Session
	prev                     [4]bool
	best                     map[string]int
	last                     rhythmcore.Grade
	gradeTimer, shake, frame int
	parts                    []particle
	scene                    *ebiten.Image
	audioContext             *audio.Context
	audioPlayer              *audio.Player
}

func New(cfg Config) *Game {
	if len(cfg.Songs) == 0 {
		panic("rhythmplay: at least one song is required")
	}
	return &Game{cfg: cfg, menu: true, best: map[string]int{}, scene: ebiten.NewImage(W, H), audioContext: audio.NewContext(48000)}
}

func (g *Game) chart() rhythmcore.Chart {
	s := g.cfg.Songs[g.song]
	if g.difficulty == 1 && len(s.Hard.Notes) > 0 {
		return s.Hard
	}
	return s.Easy
}

func (g *Game) start() {
	c := g.chart()
	g.session = rhythmcore.NewSession(c)
	g.menu = false
	g.finished = false
	g.last = ""
	g.gradeTimer = 0
	g.parts = nil
	g.prev = [4]bool{}
	if g.audioPlayer != nil {
		_ = g.audioPlayer.Close()
	}
	g.audioPlayer = g.audioContext.NewPlayerF32FromBytes(synthTrack(c, g.cfg.Songs[g.song].Tone))
	g.audioPlayer.Play()
}

// synthTrack makes a tiny original backing track directly from numbers. No
// downloaded song is needed: every beat gets a decaying click and alternating
// beats get two soft bass pitches. This also keeps the example Apache-2.0 safe.
func synthTrack(chart rhythmcore.Chart, tone float64) []byte {
	const sampleRate = 48000
	lastFrame := 360
	for _, note := range chart.Notes {
		lastFrame = max(lastFrame, note.At+note.Duration+180)
	}
	seconds := float64(lastFrame) / 60
	samples := int(seconds * sampleRate)
	data := make([]byte, samples*8) // stereo float32
	beatSeconds := 60 / float64(max(1, chart.BPM))
	if tone <= 0 {
		tone = 110
	}
	for i := 0; i < samples; i++ {
		t := float64(i) / sampleRate
		beat := int(t / beatSeconds)
		local := t - float64(beat)*beatSeconds
		v := 0.0
		if local < .07 {
			v += math.Sin(2*math.Pi*880*local) * math.Exp(-local*42) * .34
		}
		bassFreq := tone
		if beat%2 == 1 {
			bassFreq *= 1.25
		}
		v += math.Sin(2*math.Pi*bassFreq*t) * .075
		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(data[i*8:], bits)
		binary.LittleEndian.PutUint32(data[i*8+4:], bits)
	}
	return data
}

func (g *Game) Update() error {
	g.frame++
	for i := len(g.parts) - 1; i >= 0; i-- {
		p := &g.parts[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .08
		p.life--
		if p.life <= 0 {
			g.parts = append(g.parts[:i], g.parts[i+1:]...)
		}
	}
	if g.gradeTimer > 0 {
		g.gradeTimer--
	}
	if g.shake > 0 {
		g.shake--
	}
	if g.menu {
		g.updateMenu()
		return nil
	}
	if g.finished {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || pointerJustPressed() {
			g.menu = true
			g.finished = false
		}
		return nil
	}
	lanes := g.session.Chart.Lanes
	now := [4]bool{}
	keys := laneKeys(lanes)
	for i := 0; i < lanes; i++ {
		now[i] = ebiten.IsKeyPressed(keys[i])
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, _ := ebiten.TouchPosition(id)
		lane := min(lanes-1, max(0, x*lanes/W))
		now[lane] = true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		lane := min(lanes-1, max(0, x*lanes/W))
		now[lane] = true
	}
	inputs := []rhythmcore.Input{}
	for i := 0; i < lanes; i++ {
		if now[i] != g.prev[i] {
			inputs = append(inputs, rhythmcore.Input{Lane: i, Down: now[i]})
		}
	}
	g.prev = now
	for _, r := range g.session.Step(inputs) {
		g.last = r.Grade
		g.gradeTimer = 32
		if r.Grade == rhythmcore.Perfect {
			g.shake = 6
			g.burst(r.Lane, 18, color.RGBA{55, 235, 207, 255})
		} else if r.Grade == rhythmcore.Good {
			g.burst(r.Lane, 8, color.RGBA{255, 205, 88, 255})
		} else {
			g.burst(r.Lane, 3, color.RGBA{155, 164, 188, 255})
		}
	}
	if g.session.Finished() {
		g.finished = true
		key := fmt.Sprintf("%d/%d", g.song, g.difficulty)
		g.best[key] = max(g.best[key], g.session.Score)
	}
	return nil
}

func (g *Game) updateMenu() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.song = (g.song + len(g.cfg.Songs) - 1) % len(g.cfg.Songs)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.song = (g.song + 1) % len(g.cfg.Songs)
	}
	if g.cfg.Difficulty && (inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyX)) {
		g.difficulty = 1 - g.difficulty
	}
	if x, y, ok := press(); ok {
		if y >= 180 && y < 180+len(g.cfg.Songs)*72 {
			g.song = min(len(g.cfg.Songs)-1, max(0, (y-180)/72))
			return
		}
		if g.cfg.Difficulty && y >= 470 && y < 540 {
			g.difficulty = min(1, x/(W/2))
			return
		}
		if y >= 560 {
			g.start()
			return
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.start()
	}
}

func (g *Game) burst(lane, n int, c color.RGBA) {
	lanes := g.session.Chart.Lanes
	x := float64((lane*W + W/2) / lanes)
	for i := 0; i < n; i++ {
		a := float64(i) * 2.399
		speed := 1.5 + float64(i%5)*.45
		g.parts = append(g.parts, particle{x, 565, math.Cos(a) * speed, math.Sin(a)*speed - 1.5, 22 + i%14, c})
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.menu {
		g.drawMenu(screen)
		return
	}
	g.scene.Clear()
	g.drawPlay(g.scene)
	op := &ebiten.DrawImageOptions{}
	if g.shake > 0 {
		op.GeoM.Translate(float64((g.frame%3-1)*3), float64(((g.frame/2)%3-1)*2))
	}
	screen.DrawImage(g.scene, op)
}

func (g *Game) drawMenu(s *ebiten.Image) {
	s.Fill(color.RGBA{9, 17, 35, 255})
	ebitenutil.DebugPrintAt(s, g.cfg.Title, 160, 35)
	ebitenutil.DebugPrintAt(s, g.cfg.Subtitle, 70, 68)
	ebitenutil.DebugPrintAt(s, "CHOOSE SONG  [LEFT / RIGHT]", 126, 125)
	for i, song := range g.cfg.Songs {
		y := 180 + i*72
		c := color.RGBA{35, 55, 88, 255}
		if i == g.song {
			c = color.RGBA{43, 145, 151, 255}
		}
		vector.DrawFilledRect(s, 55, float32(y), 370, 58, c, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d  %s  %d BPM", i+1, song.Name, song.Easy.BPM), 90, y+23)
	}
	if g.cfg.Difficulty {
		names := []string{"EASY", "HARD"}
		for i, n := range names {
			c := color.RGBA{44, 62, 91, 255}
			if i == g.difficulty {
				c = color.RGBA{190, 119, 55, 255}
			}
			vector.DrawFilledRect(s, float32(i*W/2+8), 470, W/2-16, 58, c, false)
			ebitenutil.DebugPrintAt(s, n, i*W/2+100, 492)
		}
	}
	vector.DrawFilledRect(s, 45, 570, 390, 80, color.RGBA{207, 102, 67, 255}, false)
	ebitenutil.DebugPrintAt(s, "START  [ENTER / TAP]", 147, 605)
	ebitenutil.DebugPrintAt(s, "START creates an original beat with Go code.", 102, 680)
}

func (g *Game) drawPlay(s *ebiten.Image) {
	c := g.session.Chart
	backgrounds := [...]color.RGBA{{8, 15, 32, 255}, {24, 12, 39, 255}, {6, 31, 39, 255}}
	s.Fill(backgrounds[g.song%len(backgrounds)])
	framesPerBeat := max(1, 3600/max(1, c.BPM))
	beatPhase := float64(g.session.Frame%framesPerBeat) / float64(framesPerBeat)
	pulseR := float32(10 + (1-beatPhase)*34)
	pulseA := uint8(35 + (1-beatPhase)*100)
	vector.StrokeCircle(s, W/2, 135, pulseR, 4, color.RGBA{55, 235, 207, pulseA}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s / %s", g.cfg.Title, g.cfg.Songs[g.song].Name), 105, 22)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %06d  COMBO %03d  BEST COMBO %03d", g.session.Score, g.session.Combo, g.session.Best), 78, 48)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PERFECT %02d  GOOD %02d  MISS %02d", g.session.Perfects, g.session.Goods, g.session.Misses), 113, 70)
	laneW := float32(W / c.Lanes)
	for i := 0; i < c.Lanes; i++ {
		shade := color.RGBA{18 + uint8(i%2)*8, 28, 56, 255}
		vector.DrawFilledRect(s, float32(i)*laneW, 90, laneW-2, 530, shade, false)
	}
	judgeY := float32(560)
	vector.DrawFilledRect(s, 0, judgeY-3, W, 6, color.RGBA{255, 205, 88, 255}, false)
	speed := float32(4.2)
	for i, n := range c.Notes {
		if g.session.Resolved(i) {
			continue
		}
		x := float32(n.Lane)*laneW + laneW/2
		y := judgeY - float32(n.At-g.session.Frame)*speed
		if y < 70 || y > 650 {
			continue
		}
		nc := color.RGBA{54, 222, 199, 255}
		if n.Kind == rhythmcore.Hold {
			nc = color.RGBA{255, 183, 73, 255}
			endY := judgeY - float32(n.At+n.Duration-g.session.Frame)*speed
			top := min(y, endY)
			height := float32(math.Abs(float64(y - endY)))
			vector.DrawFilledRect(s, x-12, top, 24, max(10, height), color.RGBA{216, 137, 56, 180}, false)
		} else if n.Kind == rhythmcore.Roll {
			nc = color.RGBA{228, 95, 145, 255}
			vector.StrokeCircle(s, x, y, 24, 4, nc, false)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d/%d", g.session.RollHits(i), n.Need), int(x)-12, int(y)-5)
		}
		vector.DrawFilledCircle(s, x, y, 17, nc, false)
	}
	for _, p := range g.parts {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), 3, color.RGBA{p.c.R, p.c.G, p.c.B, uint8(min(255, p.life*10))}, false)
	}
	if g.gradeTimer > 0 {
		ebitenutil.DebugPrintAt(s, string(g.last), 210, 110)
	}
	labels := laneLabels(c.Lanes)
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i)*laneW+4, 630, laneW-8, 64, color.RGBA{43, 72, 112, 255}, false)
		ebitenutil.DebugPrintAt(s, l, int(float32(i)*laneW+laneW/2)-12, 655)
	}
	if g.finished {
		vector.DrawFilledRect(s, 35, 245, 410, 210, color.RGBA{5, 12, 28, 246}, false)
		rank := "C"
		if g.session.Misses == 0 {
			rank = "A"
		}
		if g.session.Misses == 0 && g.session.Goods == 0 {
			rank = "S"
		}
		key := fmt.Sprintf("%d/%d", g.song, g.difficulty)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("SONG CLEAR!  RANK %s", rank), 160, 282)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %d   BEST %d", g.session.Score, max(g.best[key], g.session.Score)), 145, 320)
		ebitenutil.DebugPrintAt(s, "TAP / ENTER TO CHOOSE AGAIN", 125, 390)
	}
}

func laneKeys(n int) []ebiten.Key {
	if n == 1 {
		return []ebiten.Key{ebiten.KeySpace}
	}
	if n == 2 {
		return []ebiten.Key{ebiten.KeyD, ebiten.KeyK}
	}
	return []ebiten.Key{ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK}
}
func laneLabels(n int) []string {
	if n == 1 {
		return []string{"SPACE"}
	}
	if n == 2 {
		return []string{"D", "K"}
	}
	return []string{"D", "F", "J", "K"}
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
func pointerJustPressed() bool             { _, _, ok := press(); return ok }
func (g *Game) Layout(int, int) (int, int) { return W, H }

func Run(cfg Config) {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle(cfg.Title + " — Ebitengine")
	if err := ebiten.RunGame(New(cfg)); err != nil {
		panic(err)
	}
}
