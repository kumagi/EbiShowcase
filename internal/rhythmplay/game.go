// Package rhythmplay renders the shared classroom rhythm stage.
// All rules and judgments remain in rhythmcore for deterministic tests.
// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package rhythmplay

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/ogfont"
	"github.com/kumagi/EbiShowcase/internal/rhythmcore"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	lastDelta                int
	gradeTimer, shake, frame int
	parts                    []particle
	scene                    *ebiten.Image
	audioContext             *audio.Context
	audioPlayer              *audio.Player
	offset                   int
	lang, audioState         string
	silentPractice           bool
	pulse                    *shaderlab.Pulse
	cam                      cameralab.State
	badge                    *ebiten.Image
}

var (
	rhythmFontOnce sync.Once
	rhythmFontBase *opentype.Font
	rhythmFontErr  error
	rhythmFaces    = map[float64]font.Face{}
)

func rhythmFace(size float64) font.Face {
	rhythmFontOnce.Do(func() { rhythmFontBase, rhythmFontErr = opentype.Parse(ogfont.NotoSansJP) })
	if rhythmFontErr != nil {
		panic(rhythmFontErr)
	}
	if f := rhythmFaces[size]; f != nil {
		return f
	}
	f, err := opentype.NewFace(rhythmFontBase, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		panic(err)
	}
	rhythmFaces[size] = f
	return f
}

func New(cfg Config) *Game {
	if len(cfg.Songs) == 0 {
		panic("rhythmplay: at least one song is required")
	}
	lang := browserLanguage()
	state := "SOUND READY"
	if lang == "ja" {
		state = "音声準備完了"
	}
	g := &Game{cfg: cfg, menu: true, best: map[string]int{}, scene: ebiten.NewImage(W, H), audioContext: audio.NewContext(48000), offset: storedInt("ebiShowcase.rhythm.offset"), lang: lang, audioState: state}
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{46, 230, 200, 255})
	return g
}

func (g *Game) tr(en, ja string) string {
	if g.lang == "ja" {
		return ja
	}
	return en
}

func rhythmLabel(screen *ebiten.Image, s string, x, y int, c color.Color, size float64) {
	text.Draw(screen, s, rhythmFace(size), x, y, c)
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
	key := fmt.Sprintf("%d/%d", g.song, g.difficulty)
	g.best[key] = max(g.best[key], storedInt("ebiShowcase.rhythm.best."+key))
	g.session = rhythmcore.NewSessionWithOffset(c, g.offset)
	g.menu = false
	g.finished = false
	g.last = ""
	g.lastDelta = 0
	g.gradeTimer = 0
	g.parts = nil
	g.prev = [4]bool{}
	if g.audioPlayer != nil {
		_ = g.audioPlayer.Close()
	}
	if g.silentPractice {
		g.audioState = g.tr("SILENT PRACTICE", "無音練習")
	} else {
		g.audioPlayer = g.audioContext.NewPlayerF32FromBytes(synthTrack(c, g.cfg.Songs[g.song].Tone))
		g.audioPlayer.Play()
		g.audioState = g.tr("SOUND ON", "音声オン")
	}
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
		g.lastDelta = r.Delta
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
		storeInt("ebiShowcase.rhythm.best."+key, g.best[key])
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
	if inpututil.IsKeyJustPressed(ebiten.KeyBracketLeft) {
		g.offset = max(-12, g.offset-1)
		storeInt("ebiShowcase.rhythm.offset", g.offset)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBracketRight) {
		g.offset = min(12, g.offset+1)
		storeInt("ebiShowcase.rhythm.offset", g.offset)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.silentPractice = !g.silentPractice
	}
	if x, y, ok := press(); ok {
		if y >= 180 && y < 180+len(g.cfg.Songs)*72 {
			g.song = min(len(g.cfg.Songs)-1, max(0, (y-180)/72))
			return
		}
		if y >= 410 && y < 460 {
			if x < W/2 {
				g.offset = max(-12, g.offset-1)
			} else {
				g.offset = min(12, g.offset+1)
			}
			storeInt("ebiShowcase.rhythm.offset", g.offset)
			return
		}
		if y >= 540 && y < 560 {
			g.silentPractice = !g.silentPractice
			return
		}
		if g.cfg.Difficulty && y >= 480 && y < 550 {
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
	g.cam.ViewW, g.cam.ViewH = float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	g.drawEffectBadge(screen)
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

func (g *Game) drawEffectBadge(screen *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frame)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(W-34, 10)
	screen.DrawImage(fx, op)
}

func (g *Game) drawMenu(s *ebiten.Image) {
	s.Fill(color.RGBA{9, 17, 35, 255})
	rhythmLabel(s, g.cfg.Title, 118, 42, color.RGBA{46, 230, 200, 255}, 21)
	rhythmLabel(s, g.cfg.Subtitle, 48, 76, color.RGBA{210, 220, 238, 255}, 14)
	rhythmLabel(s, g.tr("CHOOSE SONG  [LEFT / RIGHT]", "曲を選ぶ  [← / →]"), 106, 132, color.RGBA{255, 211, 112, 255}, 15)
	for i, song := range g.cfg.Songs {
		y := 180 + i*72
		c := color.RGBA{35, 55, 88, 255}
		if i == g.song {
			c = color.RGBA{43, 145, 151, 255}
		}
		vector.DrawFilledRect(s, 55, float32(y), 370, 58, c, false)
		rhythmLabel(s, fmt.Sprintf("%d  %s  %d BPM", i+1, song.Name, song.Easy.BPM), 86, y+35, color.White, 16)
	}
	vector.DrawFilledRect(s, 8, 410, W/2-16, 50, color.RGBA{42, 66, 101, 255}, false)
	vector.DrawFilledRect(s, W/2+8, 410, W/2-16, 50, color.RGBA{42, 66, 101, 255}, false)
	rhythmLabel(s, g.tr("OFFSET -", "調整 -"), 70, 442, color.White, 16)
	rhythmLabel(s, g.tr("OFFSET +", "調整 +"), 310, 442, color.White, 16)
	rhythmLabel(s, g.tr("TIMING ", "判定ずれ ")+fmt.Sprintf("%+d frames", g.offset), 148, 466, color.RGBA{184, 211, 233, 255}, 14)
	if g.cfg.Difficulty {
		names := []string{"EASY", "HARD"}
		for i, n := range names {
			c := color.RGBA{44, 62, 91, 255}
			if i == g.difficulty {
				c = color.RGBA{190, 119, 55, 255}
			}
			vector.DrawFilledRect(s, float32(i*W/2+8), 470, W/2-16, 58, c, false)
			rhythmLabel(s, n, i*W/2+100, 516, color.White, 18)
		}
	}
	vector.DrawFilledRect(s, 45, 570, 390, 80, color.RGBA{207, 102, 67, 255}, false)
	rhythmLabel(s, g.tr("START  [ENTER / TAP]", "スタート  [Enter / タップ]"), 128, 620, color.White, 18)
	rhythmLabel(s, g.audioState+" · "+g.tr("Tap START to enable the original beat.", "STARTをタップして自作ビートを有効にします。"), 44, 676, color.RGBA{184, 211, 233, 255}, 13)
	rhythmLabel(s, g.tr("M / tap here: SILENT PRACTICE", "M / ここをタップ：無音練習"), 115, 558, color.RGBA{184, 211, 233, 255}, 12)
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
	rhythmLabel(s, fmt.Sprintf("%s / %s", g.cfg.Title, g.cfg.Songs[g.song].Name), 88, 28, color.White, 15)
	rhythmLabel(s, fmt.Sprintf("%s %06d  %s %03d  %s %03d", g.tr("SCORE", "得点"), g.session.Score, g.tr("COMBO", "コンボ"), g.session.Combo, g.tr("BEST", "最大"), g.session.Best), 42, 54, color.RGBA{255, 211, 112, 255}, 14)
	rhythmLabel(s, fmt.Sprintf("PERFECT %02d  GOOD %02d  MISS %02d", g.session.Perfects, g.session.Goods, g.session.Misses), 94, 76, color.RGBA{210, 220, 238, 255}, 13)
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
		judge := string(g.last)
		if g.lastDelta < 0 {
			judge += g.tr("  EARLY ", "  早い ") + fmt.Sprint(-g.lastDelta)
		} else if g.lastDelta > 0 {
			judge += g.tr("  LATE ", "  遅い ") + fmt.Sprint(g.lastDelta)
		} else {
			judge += g.tr("  ON TIME", "  ぴったり")
		}
		col := color.RGBA{55, 235, 207, 255}
		if g.last == rhythmcore.Good {
			col = color.RGBA{255, 205, 88, 255}
		}
		if g.last == rhythmcore.Miss {
			col = color.RGBA{255, 116, 140, 255}
		}
		rhythmLabel(s, judge, 128, 122, col, 23)
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
		rhythmLabel(s, fmt.Sprintf("%s  RANK %s", g.tr("SONG CLEAR!", "曲クリア！"), rank), 144, 286, color.RGBA{46, 230, 200, 255}, 22)
		rhythmLabel(s, fmt.Sprintf("%s %d   BEST %d", g.tr("SCORE", "得点"), g.session.Score, max(g.best[key], g.session.Score)), 130, 326, color.White, 18)
		rhythmLabel(s, g.tr("TAP / ENTER TO CHOOSE AGAIN", "タップ / Enterで曲を選び直す"), 93, 395, color.RGBA{184, 211, 233, 255}, 15)
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
