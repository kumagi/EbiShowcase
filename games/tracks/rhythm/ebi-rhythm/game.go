// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/rhythmcore"
	"github.com/kumagi/EbiShowcase/internal/rhythmplay"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const (
	W = 480
	H = 720
)

type config struct {
	Title, Subtitle string
	Songs           []rhythmplay.Song
}
type sparkle struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}
type game struct {
	cfg                                        config
	song, difficulty, frame, gradeTimer, shake int
	menu, finished                             bool
	session                                    *rhythmcore.Session
	prev                                       [4]bool
	last                                       rhythmcore.Grade
	lastDelta                                  int
	best                                       map[string]int
	parts                                      []sparkle
	audio                                      *audio.Context
	player                                     *audio.Player
}

func newGame(cfg config) *game {
	if len(cfg.Songs) == 0 {
		panic("ebi-rhythm: at least one song is required")
	}
	return &game{cfg: cfg, menu: true, best: map[string]int{}, audio: audio.NewContext(48000)}
}

func (g *game) chart() rhythmcore.Chart {
	s := g.cfg.Songs[g.song]
	if g.difficulty == 1 {
		return s.Hard
	}
	return s.Easy
}

func (g *game) start() {
	g.session = rhythmcore.NewSession(g.chart())
	g.menu, g.finished = false, false
	g.last, g.gradeTimer, g.shake = "", 0, 0
	g.prev, g.parts = [4]bool{}, nil
	if g.player != nil {
		_ = g.player.Close()
	}
	g.player = g.audio.NewPlayerF32FromBytes(synthTrack(g.chart(), g.cfg.Songs[g.song].Tone))
	g.player.Play()
}

func (g *game) Update() error {
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
		if enterPressed() {
			g.menu, g.finished = true, false
		}
		return nil
	}
	now := [4]bool{}
	keys := [...]ebiten.Key{ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK}
	for i, key := range keys {
		now[i] = ebiten.IsKeyPressed(key)
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 120 {
			now[min(3, max(0, x/(W/4)))] = true
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 120 {
			now[min(3, max(0, x/(W/4)))] = true
		}
	}
	inputs := make([]rhythmcore.Input, 0, 4)
	for i := range now {
		if now[i] != g.prev[i] {
			inputs = append(inputs, rhythmcore.Input{Lane: i, Down: now[i]})
		}
	}
	g.prev = now
	for _, result := range g.session.Step(inputs) {
		g.last, g.lastDelta, g.gradeTimer = result.Grade, result.Delta, 34
		n, c := 5, color.RGBA{255, 116, 140, 255}
		if result.Grade == rhythmcore.Good {
			n, c = 12, color.RGBA{255, 210, 91, 255}
		}
		if result.Grade == rhythmcore.Perfect {
			n, c, g.shake = 22, color.RGBA{95, 245, 234, 255}, 5
		}
		g.burst(result.Lane, n, c)
	}
	if g.session.Finished() {
		g.finished = true
		key := fmt.Sprintf("%d/%d", g.song, g.difficulty)
		g.best[key] = max(g.best[key], g.session.Score)
	}
	return nil
}

func (g *game) updateMenu() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.song = (g.song + len(g.cfg.Songs) - 1) % len(g.cfg.Songs)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.song = (g.song + 1) % len(g.cfg.Songs)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.difficulty = 1 - g.difficulty
	}
	if x, y, ok := justPressed(); ok {
		if y >= 208 && y < 388 {
			g.song = min(len(g.cfg.Songs)-1, max(0, (y-208)/60))
			return
		}
		if y >= 432 && y < 522 {
			if x < W/2 {
				g.difficulty = 0
			} else {
				g.difficulty = 1
			}
			return
		}
		if y >= 574 {
			g.start()
			return
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.start()
	}
}

func (g *game) burst(lane, count int, c color.RGBA) {
	x := float64(lane*120 + 60)
	for i := 0; i < count; i++ {
		a := float64(i) * 2.399
		speed := 1.8 + float64(i%5)*.45
		g.parts = append(g.parts, sparkle{x, 573, math.Cos(a) * speed, math.Sin(a)*speed - 1.2, 22 + i%15, c})
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.menu {
		g.drawMenu(screen)
		return
	}
	scene := ebiten.NewImage(W, H)
	g.drawPlay(scene)
	op := &ebiten.DrawImageOptions{}
	if g.shake > 0 {
		op.GeoM.Translate(float64((g.frame%3-1)*3), float64(((g.frame/2)%3-1)*2))
	}
	screen.DrawImage(scene, op)
}

func stageName(song int) string {
	return []string{"stage-sunrise", "stage-neon", "stage-tempest"}[song%3]
}
func laneNote(lane int) string {
	return []string{"note-coral", "note-gold", "note-cyan", "note-violet"}[lane%4]
}

func (g *game) drawMenu(s *ebiten.Image) {
	drawCover(s, art(stageName(g.song)))
	vector.DrawFilledRect(s, 0, 0, W, H, color.RGBA{2, 8, 24, 105}, true)
	vector.DrawFilledRect(s, 0, 0, W, 124, color.RGBA{3, 10, 30, 220}, true)
	vector.DrawFilledRect(s, 16, 138, 448, 420, color.RGBA{5, 13, 36, 205}, true)
	drawContain(s, art("artist-ready"), 238, 96, 240, 465, .94)
	label(s, g.cfg.Title, 24, 22, 26, color.RGBA{255, 239, 187, 255})
	label(s, g.cfg.Subtitle, 25, 60, 13, color.RGBA{210, 237, 250, 255})
	label(s, "SELECT A LIVE SET", 28, 157, 14, color.RGBA{96, 246, 232, 255})
	for i, song := range g.cfg.Songs {
		y := 208 + i*60
		bg := color.RGBA{15, 28, 59, 220}
		edge := color.RGBA{110, 156, 195, 130}
		if i == g.song {
			bg = color.RGBA{18, 89, 109, 238}
			edge = color.RGBA{255, 221, 122, 255}
		}
		vector.DrawFilledRect(s, 25, float32(y), 265, 48, bg, true)
		vector.StrokeRect(s, 25, float32(y), 265, 48, 2, edge, true)
		label(s, fmt.Sprintf("0%d  %s", i+1, song.Name), 39, float64(y+12), 14, color.White)
		label(s, fmt.Sprintf("%d BPM", []int{song.Easy.BPM, song.Hard.BPM}[g.difficulty]), 207, float64(y+31), 10, color.RGBA{194, 231, 246, 255})
	}
	label(s, "CHART", 27, 404, 12, color.RGBA{96, 246, 232, 255})
	for i, name := range []string{"EASY", "HARD"} {
		x := float64(i*224 + 16)
		bg := color.RGBA{8, 18, 45, 218}
		if i == g.difficulty {
			bg = color.RGBA{68, 38, 98, 235}
		}
		vector.DrawFilledRect(s, float32(x), 432, 216, 88, bg, true)
		vector.StrokeRect(s, float32(x), 432, 216, 88, 2, color.RGBA{255, 219, 127, 210}, true)
		drawContain(s, art([]string{"difficulty-easy", "difficulty-hard"}[i]), x+10, 435, 72, 76, 1)
		label(s, name, x+85, 451, 18, color.White)
		stars := 1 + i*2
		label(s, fmt.Sprintf("%d PEARL CHART", stars), x+85, 480, 10, color.RGBA{210, 237, 250, 255})
	}
	vector.DrawFilledRect(s, 38, 574, 404, 92, color.RGBA{218, 77, 112, 240}, true)
	vector.StrokeRect(s, 38, 574, 404, 92, 3, color.RGBA{255, 232, 160, 255}, true)
	label(s, "BEGIN PERFORMANCE", 100, 595, 23, color.White)
	label(s, "ENTER / TAP  •  D F J K", 130, 633, 11, color.RGBA{255, 239, 187, 255})
}

func (g *game) drawPlay(s *ebiten.Image) {
	drawCover(s, art(stageName(g.song)))
	vector.DrawFilledRect(s, 0, 0, W, H, color.RGBA{2, 6, 20, 45}, true)
	artist := "artist-play"
	if g.finished || g.progress() > .86 {
		artist = "artist-encore"
	}
	drawContain(s, art(artist), 235, 118, 245, 405, .62)
	// The highway leaves the performer visible, but has enough contrast to read at phone size.
	vector.DrawFilledRect(s, 18, 118, 444, 510, color.RGBA{3, 8, 28, 192}, true)
	vector.StrokeRect(s, 18, 118, 444, 510, 2, color.RGBA{126, 231, 239, 190}, true)
	for i := 0; i < 4; i++ {
		x := float32(20 + i*110)
		vector.DrawFilledRect(s, x, 122, 108, 501, color.RGBA{uint8(15 + i*5), 25, 57, 165}, true)
		vector.StrokeLine(s, x, 122, x, 623, 1, color.RGBA{159, 224, 239, 100}, true)
	}
	// Progress and difficulty are generated ornaments carrying live state.
	drawContain(s, art("progress-rail"), 76, 70, 330, 42, 1)
	vector.DrawFilledRect(s, 122, 89, float32(236*g.progress()), 6, color.RGBA{91, 245, 232, 255}, true)
	vector.DrawFilledCircle(s, float32(122+236*g.progress()), 92, 7, color.RGBA{255, 223, 120, 255}, true)
	drawContain(s, art([]string{"difficulty-easy", "difficulty-hard"}[g.difficulty]), 10, 8, 70, 68, 1)
	label(s, g.cfg.Songs[g.song].Name, 89, 15, 17, color.White)
	label(s, fmt.Sprintf("SCORE %06d   COMBO %03d", g.session.Score, g.session.Combo), 89, 43, 13, color.RGBA{255, 230, 153, 255})
	judgeY := float32(568)
	drawContain(s, art("judgment-rail"), 16, float64(judgeY-31), 448, 62, 1)
	for i := 0; i < 4; i++ {
		vector.DrawFilledCircle(s, float32(i*110+74), judgeY, 19, color.NRGBA{63, 221, 226, 78}, true)
	}
	speed := float32(4.15)
	for i, n := range g.session.Chart.Notes {
		if g.session.Resolved(i) {
			continue
		}
		x := float64(n.Lane*110 + 45)
		y := float64(judgeY - float32(n.At-g.session.Frame)*speed)
		if y < 92 || y > 635 {
			continue
		}
		if n.Kind == rhythmcore.Hold {
			end := float64(judgeY - float32(n.At+n.Duration-g.session.Frame)*speed)
			top := min(y, end)
			height := math.Abs(y - end)
			vector.DrawFilledRect(s, float32(x+34), float32(top), 24, float32(max(16.0, height)), color.RGBA{255, 197, 86, 185}, true)
			drawContain(s, art("note-hold"), x+10, y-30, 72, 62, 1)
		} else if n.Kind == rhythmcore.Roll {
			drawContain(s, art("note-roll"), x+9, y-35, 74, 70, 1)
			label(s, fmt.Sprintf("%d/%d", g.session.RollHits(i), n.Need), x+33, y-7, 10, color.White)
		} else {
			drawContain(s, art(laneNote(n.Lane)), x+10, y-30, 70, 60, 1)
		}
	}
	for _, p := range g.parts {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), 3, color.NRGBA{p.c.R, p.c.G, p.c.B, uint8(min(255, p.life*10))}, true)
	}
	if g.gradeTimer > 0 {
		judge := string(g.last)
		if g.lastDelta < 0 {
			judge += fmt.Sprintf("  EARLY %d", -g.lastDelta)
		} else if g.lastDelta > 0 {
			judge += fmt.Sprintf("  LATE %d", g.lastDelta)
		} else {
			judge += "  ON BEAT"
		}
		c := color.RGBA{95, 245, 234, 255}
		if g.last == rhythmcore.Good {
			c = color.RGBA{255, 217, 104, 255}
		}
		if g.last == rhythmcore.Miss {
			c = color.RGBA{255, 112, 146, 255}
		}
		label(s, judge, 132, 130, 22, c)
	}
	for i, key := range []string{"D", "F", "J", "K"} {
		x := float32(i * 120)
		vector.DrawFilledRect(s, x, 630, 118, 90, color.RGBA{8, 20, 49, 235}, true)
		vector.StrokeRect(s, x, 630, 118, 90, 2, color.RGBA{90, 231, 227, 170}, true)
		label(s, key, float64(x+53), 651, 20, color.White)
	}
	if g.finished {
		g.drawResult(s)
	}
}

func (g *game) drawResult(s *ebiten.Image) {
	vector.DrawFilledRect(s, 24, 175, 432, 350, color.RGBA{3, 9, 28, 244}, true)
	vector.StrokeRect(s, 24, 175, 432, 350, 3, color.RGBA{255, 222, 118, 255}, true)
	drawContain(s, art("artist-encore"), 278, 200, 155, 270, 1)
	rank := "B"
	if g.session.Misses == 0 {
		rank = "A"
	}
	if g.session.Misses == 0 && g.session.Goods == 0 {
		rank = "S"
	}
	label(s, "OCEAN ENCORE", 51, 206, 14, color.RGBA{94, 245, 232, 255})
	label(s, "SONG CLEAR", 51, 237, 28, color.White)
	label(s, "RANK "+rank, 52, 292, 40, color.RGBA{255, 222, 118, 255})
	label(s, fmt.Sprintf("SCORE  %06d", g.session.Score), 52, 358, 16, color.White)
	label(s, fmt.Sprintf("PERFECT %02d  GOOD %02d  MISS %02d", g.session.Perfects, g.session.Goods, g.session.Misses), 52, 392, 12, color.RGBA{210, 237, 250, 255})
	label(s, "TAP / ENTER TO RETURN", 52, 475, 13, color.RGBA{255, 163, 190, 255})
}

func (g *game) progress() float64 {
	last := 360
	for _, n := range g.session.Chart.Notes {
		last = max(last, n.At+n.Duration+20)
	}
	return min(1, float64(g.session.Frame)/float64(last))
}

func label(dst *ebiten.Image, value string, x, y, size float64, c color.Color) {
	face, err := uilab.Face("en", size)
	if err != nil {
		ebitenutil.DebugPrintAt(dst, value, int(x), int(y))
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, value, face, op)
}

func justPressed() (int, int, bool) {
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
func enterPressed() bool {
	_, _, p := justPressed()
	return p || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
}
func (g *game) Layout(int, int) (int, int) { return W, H }

func synthTrack(chart rhythmcore.Chart, tone float64) []byte {
	const rate = 48000
	last := 360
	for _, n := range chart.Notes {
		last = max(last, n.At+n.Duration+180)
	}
	samples := int(float64(last) / 60 * rate)
	data := make([]byte, samples*8)
	beatSeconds := 60 / float64(max(1, chart.BPM))
	for i := 0; i < samples; i++ {
		t := float64(i) / rate
		beat := int(t / beatSeconds)
		local := t - float64(beat)*beatSeconds
		v := math.Sin(2*math.Pi*tone*t) * .07
		if beat%2 == 1 {
			v += math.Sin(2*math.Pi*tone*1.5*t) * .035
		}
		if local < .075 {
			v += math.Sin(2*math.Pi*920*local) * math.Exp(-local*40) * .32
		}
		bits := math.Float32bits(float32(v))
		binary.LittleEndian.PutUint32(data[i*8:], bits)
		binary.LittleEndian.PutUint32(data[i*8+4:], bits)
	}
	return data
}

func run(cfg config) {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle(cfg.Title + " — Ebitengine")
	if err := ebiten.RunGame(newGame(cfg)); err != nil {
		panic(err)
	}
}
