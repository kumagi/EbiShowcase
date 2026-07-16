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
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const (
	W      = 480
	H      = 720
	worldW = 3200
	ground = 520
)

type game struct {
	x, y, vx, vy, cam                   float64
	onGround, dash, highJump, won, lost bool
	facing                              int
	hp, invuln, attack, abilityFX       int
	revealed                            map[int]bool
	relics                              map[int]bool
	enemies                             []enemy
	frames                              int
	message                             string
	flash, shake, bestFrames            int
	sparks                              []spark
	audio                               *audio.Context
	gate                                audiolab.Gate
	pulse                               *shaderlab.Pulse
	camState                            cameralab.State
	badge                               *ebiten.Image
}
type spark struct{ x, y, vx, vy, life float64 }
type enemy struct {
	x, home, dir float64
	hp, maxHP    int
	hit          int
	guardian     bool
}

func newGame() *game {
	prepareDepthsArt()
	g := &game{x: 80, y: ground - 36, facing: 1, hp: 5, revealed: map[int]bool{}, relics: map[int]bool{720: true, 1450: true, 2350: true, 2980: true}, message: "Explore, fight, and reveal the enormous world room by room."}
	for _, e := range []enemy{
		{x: 330, home: 330, dir: 1, hp: 1, maxHP: 1}, {x: 790, home: 790, dir: -1, hp: 1, maxHP: 1}, {x: 1030, home: 1030, dir: -1, hp: 4, maxHP: 4, guardian: true},
		{x: 1290, home: 1290, dir: 1, hp: 1, maxHP: 1}, {x: 1640, home: 1640, dir: -1, hp: 1, maxHP: 1}, {x: 2120, home: 2120, dir: -1, hp: 4, maxHP: 4, guardian: true},
		{x: 2390, home: 2390, dir: 1, hp: 1, maxHP: 1}, {x: 2730, home: 2730, dir: -1, hp: 1, maxHP: 1}, {x: 3050, home: 3050, dir: -1, hp: 5, maxHP: 5, guardian: true},
	} {
		g.enemies = append(g.enemies, e)
	}
	g.audio = audiolab.Context()
	g.pulse = shaderlab.NewPulse()
	g.camState = cameralab.State{Pos: cameralab.Vec{X: W / 2, Y: H / 2}, ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{245, 190, 68, 255})
	return g
}
func floorAt(x float64) float64 {
	if x > 980 && x < 1250 {
		return 610
	}
	if x > 1900 && x < 2150 {
		return 585
	}
	return ground - 45*math.Sin(x/330)
}
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			best := g.bestFrames
			*g = *newGame()
			g.bestFrames = best
		}
		return nil
	}
	g.frames++
	if g.flash > 0 {
		g.flash--
	}
	if g.shake > 0 {
		g.shake--
	}
	if g.invuln > 0 {
		g.invuln--
	}
	if g.attack > 0 {
		g.attack--
	}
	if g.abilityFX > 0 {
		g.abilityFX--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	left, right, jump, attack, dash := controls()
	acc := .45
	if left {
		g.vx -= acc
		g.facing = -1
	}
	if right {
		g.vx += acc
		g.facing = 1
	}
	if attack && g.attack == 0 {
		g.attack = 18
		g.play(520)
	}
	g.vx *= .86
	if dash && g.dash {
		if right {
			g.vx = 9
		} else if left {
			g.vx = -9
		} else {
			g.vx = 9
		}
		for i := 0; i < 8; i++ {
			g.sparks = append(g.sparks, spark{g.x - float64(i)*5, g.y + 20, -.3, 0, 18})
		}
	}
	if jump && g.onGround {
		g.vy = -9
		if g.highJump {
			g.vy = -12
		}
		g.onGround = false
	}
	g.vy += .45
	g.x += g.vx
	// Each region guardian is a combat gate. Its collision rule lives in
	// Update; Draw only chooses how the same guardian state looks.
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.guardian && e.hp > 0 && g.x > e.x-62 && g.x < e.x+90 {
			g.x = e.x - 62
			if g.vx > 0 {
				g.vx = 0
			}
			g.message = "A guardian seals this region. ATTACK, retreat, then strike again."
		}
	}
	if !g.dash && g.x > 1710 {
		g.x = 1710
		g.vx = 0
		g.message = "A sealed current blocks the path. Find the dash crest."
	}
	if !g.highJump && g.x > 2520 {
		g.x = 2520
		g.vx = 0
		g.message = "The crystal ledge is too high. Find the tide wings."
	}
	g.y += g.vy
	f := floorAt(g.x)
	if g.y >= f-36 {
		g.y = f - 36
		g.vy = 0
		g.onGround = true
	}
	g.x = math.Max(20, math.Min(worldW-20, g.x))
	g.cam += (g.x - 240 - g.cam) * .1
	g.cam = math.Max(0, math.Min(worldW-W, g.cam))
	room := int(g.x / 400)
	g.revealed[room] = true
	g.updateEnemies()
	for rx := range g.relics {
		if math.Abs(g.x-float64(rx)) < 35 {
			g.play(760)
			delete(g.relics, rx)
			g.abilityFX = 90
			if rx == 1450 {
				g.dash = true
				g.message = "Dash crest found! Hold a direction and press X."
			} else if rx == 2350 {
				g.highJump = true
				g.message = "Tide wings found! Jumps now reach high ledges."
			} else {
				g.message = "Map fragment found. More of the world is recorded."
			}
			g.burst(float64(rx), g.y, 18)
		}
	}
	if len(g.relics) == 0 && g.x > 3100 {
		g.won = true
		if g.bestFrames == 0 || g.frames < g.bestFrames {
			g.bestFrames = g.frames
		}
		g.burst(g.x, g.y, 36)
	}
	for _, hx := range []float64{1080, 2050, 2790} {
		if math.Abs(g.x-hx) < 22 && g.onGround {
			g.flash = 20
			g.shake = 7
			g.x = math.Max(40, float64(int(g.x/400))*400+40)
			g.vx = 0
			g.message = "Spikes! Returned to the room entrance."
		}
	}
	if g.hp <= 0 || g.y > 690 || g.frames > 150*60 {
		g.lost = true
	}
	return nil
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .06)).Play()
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 26 + float64(i%8)})
	}
}

func (g *game) updateEnemies() {
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp <= 0 {
			continue
		}
		if e.hit > 0 {
			e.hit--
		}
		if !e.guardian {
			e.x += e.dir * .38
			if math.Abs(e.x-e.home) > 64 {
				e.dir *= -1
			}
		}
		dx := e.x - g.x
		inFront := dx*float64(g.facing) > -8
		if g.attack >= 7 && g.attack <= 14 && math.Abs(dx) < 82 && inFront && e.hit == 0 {
			e.hp--
			e.hit = 16
			g.shake = 4
			g.play(720)
			g.burst(e.x, floorAt(e.x)-42, 12)
			if e.hp <= 0 {
				g.message = "Guardian broken! The route ahead is safe to explore."
			}
		}
		contactRange := 34.0
		if e.guardian {
			contactRange = 66
		}
		if math.Abs(dx) < contactRange && g.invuln == 0 && e.hit == 0 {
			g.hp--
			g.invuln = 75
			g.flash = 16
			g.shake = 7
			g.vx = -float64(g.facing) * 6
			g.message = "Hit! Step back, then use ATTACK when the enemy enters range."
		}
	}
}

func (g *game) Draw(s *ebiten.Image) {
	region := minInt(2, int(g.x/1100))
	bgs := []color.RGBA{{10, 18, 33, 255}, {23, 31, 53, 255}, {47, 24, 48, 255}}
	s.Fill(bgs[region])
	regionStart := float64(region * 1100)
	drawDepthsRegion(s, []string{"gardens", "abyss", "sanctum"}[region], (g.x-regionStart)/1100)
	vector.DrawFilledRect(s, 0, 80, W, 42, color.RGBA{3, 8, 20, 65}, false)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.frames)*2) * 5
	}
	for x := int(g.cam/40) * 40; x < int(g.cam)+W+80; x += 40 {
		sx := float32(float64(x) - g.cam + ox)
		f := float32(floorAt(float64(x)))
		// Collision remains the same floorAt number. Overlapping generated ledges
		// turn that boundary into region-specific rock and ruins.
		drawDepthsPlatform(s, []string{"gardens", "abyss", "sanctum"}[region], float64(sx)+20, float64(f)-2, 96)
		vector.StrokeLine(s, sx, f, sx+42, f, 3, []color.RGBA{{93, 230, 206, 220}, {196, 123, 235, 210}, {240, 91, 147, 210}}[region], true)
	}
	for _, hx := range []float64{1080, 2050, 2790} {
		sx := float32(hx - g.cam + ox)
		if sx > -30 && sx < W+30 {
			for i := 0; i < 3; i++ {
				vector.DrawFilledRect(s, sx+float32(i*12-18), float32(floorAt(hx)-15), 8, 15, color.RGBA{245, 90, 90, 255}, false)
			}
		}
	}
	for rx := range g.relics {
		sx := float32(float64(rx) - g.cam)
		if sx > -30 && sx < W+30 {
			ry := float32(floorAt(float64(rx)) - 65 + math.Sin(float64(g.frames)*.09+float64(rx))*6)
			vector.DrawFilledCircle(s, sx, ry, 30, color.RGBA{120, 235, 255, 42}, true)
			drawDepthsCharacter(s, "spirit", float64(sx), float64(ry), 72, 72, false)
		}
	}
	for _, e := range g.enemies {
		if e.hp <= 0 {
			continue
		}
		sx := e.x - g.cam + ox
		if sx < -100 || sx > W+100 {
			continue
		}
		ey := floorAt(e.x) - 40
		name, w, h := "beetle", 76.0, 66.0
		if e.guardian {
			name, w, h = "guardian", 130, 142
			ey = floorAt(e.x) - 68
		}
		if e.hit > 0 {
			vector.DrawFilledCircle(s, float32(sx), float32(ey), float32(w*.35), color.RGBA{255, 230, 185, 90}, true)
		}
		vector.DrawFilledCircle(s, float32(sx), float32(floorAt(e.x)), float32(w*.27), color.RGBA{2, 5, 14, 120}, true)
		drawDepthsCharacter(s, name, sx, ey, w, h, e.dir > 0)
		vector.DrawFilledRect(s, float32(sx-w*.28), float32(ey+h*.38), float32(w*.56), 5, color.RGBA{8, 10, 20, 220}, false)
		vector.DrawFilledRect(s, float32(sx-w*.28), float32(ey+h*.38), float32(w*.56)*float32(e.hp)/float32(e.maxHP), 5, color.RGBA{235, 82, 114, 255}, false)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x-g.cam+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	px, py := float32(g.x-g.cam+ox), float32(g.y)
	bob := float32(0)
	if math.Abs(g.vx) > .3 && g.onGround {
		bob = float32(math.Sin(float64(g.frames)*.3) * 3)
	}
	// Large animated hero silhouette: readable even when the home card is tiny.
	shadowW := float32(28 + math.Abs(g.vx)*1.5)
	vector.DrawFilledCircle(s, px, float32(floorAt(g.x)-2), shadowW, color.RGBA{2, 6, 14, 90}, true)
	if g.flash > 0 {
		vector.DrawFilledCircle(s, px, py+18, 34, color.RGBA{255, 110, 100, 90}, true)
	}
	drawDepthsCharacter(s, "tenjiroh", float64(px), float64(py+bob+12), 92, 92, g.facing < 0)
	if g.attack > 0 {
		dir := float32(g.facing)
		phase := float32(18-g.attack) / 18
		cx := px + dir*(34+phase*22)
		vector.StrokeCircle(s, cx, py+8, 25+phase*20, 6, color.RGBA{100, 239, 225, uint8(220 * (1 - phase))}, true)
	}
	if g.abilityFX > 0 {
		vector.DrawFilledCircle(s, W/2, 325, 118, color.RGBA{111, 225, 255, uint8(min(150, g.abilityFX*2))}, true)
		drawDepthsCharacter(s, "spirit", W/2, 315, 235, 235, false)
		ebitenutil.DebugPrintAt(s, "ABILITY AWAKENED", 174, 430)
	}
	vector.DrawFilledRect(s, 0, 0, W, 80, color.RGBA{5, 11, 24, 235}, false)
	g.drawHUD(s, region)
	g.drawEffectBadge(s)
	ebitenutil.DebugPrintAt(s, g.message, 25, 46)
	for i := 0; i < 8; i++ {
		c := color.RGBA{39, 48, 64, 255}
		if g.revealed[i] {
			c = color.RGBA{84, 151, 145, 255}
		}
		vector.DrawFilledRect(s, float32(40+i*50), 90, 42, 20, c, false)
	}
	// Ability crests make the long-term goals visible from the first second.
	for i, unlocked := range []bool{g.dash, g.highJump} {
		x := float32(330 + i*58)
		c := color.RGBA{45, 53, 72, 255}
		if unlocked {
			c = []color.RGBA{{76, 211, 214, 255}, {220, 126, 231, 255}}[i]
		}
		vector.DrawFilledCircle(s, x, 132, 20, c, true)
		glyph := []string{"DASH", "WING"}[i]
		ebitenutil.DebugPrintAt(s, glyph, int(x)-16, 128)
	}
	labels := []string{"LEFT", "JUMP", "ATTACK", "DASH", "RIGHT"}
	vector.DrawFilledRect(s, 0, 638, W, 82, color.RGBA{3, 8, 18, 175}, false)
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*96+3), 650, 90, 55, color.RGBA{45, 78, 113, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*96+20, 675)
	}
	if g.won {
		overlay(s, "THE DEPTHS MAPPED!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(s, "EXPLORATION FAILED\n\nTAP / ENTER TO RETRY")
	}
}
func (g *game) drawHUD(s *ebiten.Image, region int) {
	label := fmt.Sprintf("REGION %d/3 WORLD %04d ROOMS %d/8 RELICS %d DASH %v WINGS %v BEST %.1f", region+1, int(g.x), len(g.revealed), len(g.relics), g.dash, g.highJump, float64(g.bestFrames)/60)
	if face, err := uilab.Face("en", 13); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(10, 6)
		text.Draw(s, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(s, label, 10, 18)
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
	op.GeoM.Translate(W-34, 52)
	s.DrawImage(fx, op)
}
func controls() (bool, bool, bool, bool, bool) {
	left := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	jump := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyZ)
	attack := inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyK)
	dash := inpututil.IsKeyJustPressed(ebiten.KeyX)
	buttonAt := func(x, y int) int {
		if y < 640 || x < 0 || x >= W {
			return -1
		}
		return x / 96
	}
	// A mouse can hold a direction. Jump, attack, and dash fire only on the
	// initial press, so holding a button cannot repeat an action every tick.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		switch buttonAt(x, y) {
		case 0:
			left = true
		case 4:
			right = true
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		switch buttonAt(x, y) {
		case 1:
			jump = true
		case 2:
			attack = true
		case 3:
			dash = true
		}
	}
	// Inspect every active touch so a thumb can hold LEFT/RIGHT while a second
	// thumb taps JUMP, ATTACK, or DASH.
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		switch buttonAt(x, y) {
		case 0:
			left = true
		case 4:
			right = true
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		switch buttonAt(x, y) {
		case 1:
			jump = true
		case 2:
			attack = true
		case 3:
			dash = true
		}
	}
	return left, right, jump, attack, dash
}
func retry() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 45, 270, 390, 150, color.RGBA{4, 10, 24, 245}, false)
	ebitenutil.DebugPrintAt(s, t, 110, 330)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Depths")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
