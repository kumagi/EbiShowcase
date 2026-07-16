package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"sync"

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
	screenW     = 480
	screenH     = 720
	allyRadius  = 19.0
	enemyRadius = 27.0
	stopSpeed   = 0.11
	friction    = 0.982
	maxTurns    = 8
)

//go:embed assets/strike-pearl-coliseum.png assets/strike-allies-atlas.png assets/strike-enemies-atlas.png assets/strike-obstacles-atlas.png
var strikeArtFS embed.FS

var (
	strikeArtOnce   sync.Once
	strikeArt       map[string]*ebiten.Image
	strikeAllies    [2]*ebiten.Image
	strikeEnemies   [3]*ebiten.Image
	strikeObstacles [3]*ebiten.Image
	strikeFace14    *text.GoTextFace
	strikeFace16    *text.GoTextFace
	strikeFace20    *text.GoTextFace
)

type vec struct{ x, y float64 }

type ally struct {
	pos      vec
	velocity vec
}

type enemy struct {
	pos      vec
	hp       int
	cooldown int
}
type spark struct{ x, y, vx, vy, life float64 }

type game struct {
	allies                              [2]ally
	enemies                             []enemy
	active, turns                       int
	dragging, moving                    bool
	dragNow, dragOrigin                 vec
	allyEffectUsed                      bool
	pulseAt                             vec
	pulseFrames                         int
	message                             string
	won, lost                           bool
	stage, totalTurns, bestTurns, shake int
	tick                                int
	pillars                             []vec
	sparks                              []spark
	audio                               *audio.Context
	gate                                audiolab.Gate
	shader                              *shaderlab.Pulse
	cam                                 cameralab.State
	badge                               *ebiten.Image
}

func newGame() *game {
	loadStrikeArt()
	g := &game{stage: 1}
	g.audio = audiolab.Context()
	g.shader = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: screenW / 2, Y: screenH / 2}, ViewW: screenW, ViewH: screenH}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{250, 210, 72, 255})
	g.loadStage()
	return g
}
func (g *game) loadStage() {
	g.allies = [2]ally{{pos: vec{125, 525}}, {pos: vec{355, 525}}}
	enemySets := [][]enemy{{{pos: vec{105, 185}, hp: 2}, {pos: vec{370, 205}, hp: 2}, {pos: vec{240, 345}, hp: 3}}, {{pos: vec{80, 160}, hp: 2}, {pos: vec{240, 220}, hp: 4}, {pos: vec{400, 160}, hp: 2}, {pos: vec{240, 430}, hp: 3}}, {{pos: vec{75, 165}, hp: 3}, {pos: vec{405, 165}, hp: 3}, {pos: vec{120, 355}, hp: 3}, {pos: vec{360, 355}, hp: 3}, {pos: vec{240, 260}, hp: 5}}}
	pillarSets := [][]vec{{{240, 205}, {145, 375}, {350, 390}}, {{160, 255}, {320, 255}, {160, 420}, {320, 420}}, {{240, 170}, {90, 270}, {390, 270}, {240, 390}}}
	g.enemies = enemySets[g.stage-1]
	g.pillars = pillarSets[g.stage-1]
	g.active = 0
	g.turns = 0
	g.dragging = false
	g.moving = false
	g.won = false
	g.lost = false
	g.message = "Drag the glowing ally backward and release."
}

func (g *game) Update() error {
	g.tick++
	if g.won || g.lost {
		if retryPressed() {
			if g.won && g.stage < 3 {
				g.totalTurns += g.turns
				g.stage++
				g.loadStage()
			} else {
				best := g.bestTurns
				if g.won {
					total := g.totalTurns + g.turns
					if best == 0 || total < best {
						best = total
					}
				}
				*g = *newGame()
				g.bestTurns = best
			}
		}
		return nil
	}
	if g.pulseFrames > 0 {
		g.pulseFrames--
	}
	if g.shake > 0 {
		g.shake--
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
	for i := range g.enemies {
		if g.enemies[i].cooldown > 0 {
			g.enemies[i].cooldown--
		}
	}

	if !g.moving {
		g.updateAim()
		return nil
	}

	a := &g.allies[g.active]
	a.pos.x += a.velocity.x
	a.pos.y += a.velocity.y
	g.bounceWalls(a)
	g.bouncePillars(a)
	g.hitEnemies(a)
	g.hitAlly(a)
	a.velocity.x *= friction
	a.velocity.y *= friction
	if math.Hypot(a.velocity.x, a.velocity.y) < stopSpeed {
		a.velocity = vec{}
		g.moving = false
		g.endTurn()
	}
	return nil
}

func (g *game) updateAim() {
	a := &g.allies[g.active]
	if !g.dragging {
		x, y, ok := pressPosition()
		if ok && distance(vec{float64(x), float64(y)}, a.pos) <= allyRadius+18 {
			g.dragging = true
			g.dragOrigin = a.pos
			g.dragNow = vec{float64(x), float64(y)}
		}
		return
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.dragNow = vec{float64(x), float64(y)}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		g.dragNow = vec{float64(x), float64(y)}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) || len(inpututil.AppendJustReleasedTouchIDs(nil)) > 0 {
		dragged := dragPosition(g.dragOrigin, g.dragNow)
		pull := vec{g.dragOrigin.x - dragged.x, g.dragOrigin.y - dragged.y}
		length := math.Hypot(pull.x, pull.y)
		g.dragging = false
		if length < 14 {
			g.message = "Pull farther before release."
			return
		}
		a.pos = dragged
		a.velocity = vec{pull.x * 0.105, pull.y * 0.105}
		g.play(380)
		g.turns++
		g.moving = true
		g.allyEffectUsed = false
		g.message = "Moving: bounce, attack, then slow to rest."
	}
}

func (g *game) bounceWalls(a *ally) {
	if a.pos.x < 28 {
		a.pos.x = 28
		a.velocity.x = math.Abs(a.velocity.x) * 0.88
	}
	if a.pos.x > 452 {
		a.pos.x = 452
		a.velocity.x = -math.Abs(a.velocity.x) * 0.88
	}
	if a.pos.y < 105 {
		a.pos.y = 105
		a.velocity.y = math.Abs(a.velocity.y) * 0.88
	}
	if a.pos.y > 575 {
		a.pos.y = 575
		a.velocity.y = -math.Abs(a.velocity.y) * 0.88
	}
}

func (g *game) bouncePillars(a *ally) {
	for _, pillar := range g.pillars {
		g.reflectCircle(a, pillar, 24)
	}
}

func (g *game) reflectCircle(a *ally, center vec, otherRadius float64) float64 {
	dx, dy := a.pos.x-center.x, a.pos.y-center.y
	d := math.Hypot(dx, dy)
	minimum := allyRadius + otherRadius
	if d >= minimum {
		return 0
	}
	if d == 0 {
		dx, d = 1, 1
	}
	nx, ny := dx/d, dy/d
	a.pos = vec{center.x + nx*minimum, center.y + ny*minimum}
	dot := a.velocity.x*nx + a.velocity.y*ny
	impact := math.Abs(dot)
	if dot < 0 {
		a.velocity.x -= 1.85 * dot * nx
		a.velocity.y -= 1.85 * dot * ny
	}
	return impact
}

func (g *game) hitEnemies(a *ally) {
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp <= 0 {
			continue
		}
		impact := g.reflectCircle(a, e.pos, enemyRadius)
		if impact >= 1.2 && e.cooldown == 0 {
			e.hp--
			g.play(760)
			e.cooldown = 22
			g.shake = 4
			g.burst(e.pos.x, e.pos.y, 10)
			g.message = fmt.Sprintf("Direct contact! Enemy HP %d.", e.hp)
		}
	}
	g.checkWin()
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 26 + float64(i%8)})
	}
}

func (g *game) hitAlly(a *ally) {
	if g.allyEffectUsed {
		return
	}
	other := &g.allies[1-g.active]
	if g.reflectCircle(a, other.pos, allyRadius) == 0 {
		return
	}
	g.allyEffectUsed = true
	g.pulseAt = other.pos
	g.pulseFrames = 30
	hits := 0
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp > 0 && distance(e.pos, other.pos) <= 165 {
			e.hp--
			hits++
		}
	}
	g.message = fmt.Sprintf("ALLY WAVE! %d enemy hit(s).", hits)
	g.play(560)
	g.checkWin()
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .07)).Play()
}

func (g *game) checkWin() {
	for _, e := range g.enemies {
		if e.hp > 0 {
			return
		}
	}
	g.won = true
	g.moving = false
	g.message = "Every reef guardian is defeated!"
}

func (g *game) endTurn() {
	if g.won {
		return
	}
	if g.turns >= maxTurns {
		g.lost = true
		g.message = "No turns left. Plan more ally waves!"
		return
	}
	g.active = 1 - g.active
	g.message = fmt.Sprintf("Turn ended at rest. Ally %d is ready.", g.active+1)
}

func (g *game) Draw(screen *ebiten.Image) {
	drawStrikeCover(screen, strikeArt["arena"], 0, 0, screenW, screenH)
	washes := []color.RGBA{{10, 106, 118, 10}, {82, 37, 128, 24}, {154, 28, 57, 30}}
	vector.DrawFilledRect(screen, 0, 0, screenW, screenH, washes[g.stage-1], false)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2.4) * 5
	}
	vector.DrawFilledRect(screen, 10, 8, 460, 78, color.RGBA{4, 15, 31, 226}, true)
	vector.StrokeRect(screen, 10, 8, 460, 78, 2, color.RGBA{255, 225, 157, 165}, true)
	g.drawTitle(screen)
	drawCenteredStrikeLabel(screen, fmt.Sprintf("STAGE %d/3   •   TURN %d/%d   •   TARGETS %d   •   BEST %d", g.stage, g.turns, maxTurns, g.alive(), g.bestTurns), 240, 35, strikeFace14, color.White)
	drawCenteredStrikeLabel(screen, g.message, 240, 59, strikeFace14, color.RGBA{140, 239, 247, 255})
	vector.StrokeRect(screen, 18, 92, 444, 500, 3, color.RGBA{103, 218, 231, 105}, true)

	for _, p := range g.pillars {
		vector.DrawFilledCircle(screen, float32(p.x+ox), float32(p.y+18), 30, color.RGBA{2, 13, 27, 85}, true)
		drawStrikeSprite(screen, strikeObstacles[g.stage-1], p.x+ox, p.y, 76, 0, 1, false)
	}
	for _, e := range g.enemies {
		if e.hp <= 0 {
			continue
		}
		size := []float64{76, 84, 98}[g.stage-1]
		bob := math.Sin(float64(g.tick)*.09+e.pos.x) * 2.5
		vector.DrawFilledCircle(screen, float32(e.pos.x+ox), float32(e.pos.y+size*.27), float32(size*.32), color.RGBA{2, 12, 28, 95}, true)
		drawStrikeSprite(screen, strikeEnemies[g.stage-1], e.pos.x+ox, e.pos.y+bob, size, 0, 1, e.cooldown > 0)
		vector.DrawFilledRect(screen, float32(e.pos.x-24+ox), float32(e.pos.y+size*.42), 48, 21, color.RGBA{4, 16, 31, 220}, true)
		vector.StrokeRect(screen, float32(e.pos.x-24+ox), float32(e.pos.y+size*.42), 48, 21, 1, color.RGBA{255, 123, 133, 175}, true)
		drawCenteredStrikeLabel(screen, fmt.Sprintf("HP %d", e.hp), e.pos.x+ox, e.pos.y+size*.42+2, strikeFace14, color.White)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(screen, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	if g.pulseFrames > 0 {
		r := float32(45 + (30-g.pulseFrames)*4)
		vector.StrokeCircle(screen, float32(g.pulseAt.x), float32(g.pulseAt.y), r, 5, color.RGBA{250, 210, 72, 210}, false)
	}
	for i, a := range g.allies {
		drawPos := a.pos
		if g.dragging && i == g.active {
			drawPos = dragPosition(g.dragOrigin, g.dragNow)
		}
		angle := math.Sin(float64(g.tick)*.08+float64(i)) * .035
		if g.moving && i == g.active {
			angle = math.Atan2(a.velocity.y, a.velocity.x) + float64(g.tick)*.16
		}
		vector.DrawFilledCircle(screen, float32(drawPos.x), float32(drawPos.y+18), 29, color.RGBA{2, 12, 28, 90}, true)
		drawStrikeSprite(screen, strikeAllies[i], drawPos.x, drawPos.y, 76, angle, 1, false)
		if i == g.active && !g.moving {
			vector.StrokeCircle(screen, float32(drawPos.x), float32(drawPos.y), 42+float32(math.Sin(float64(g.tick)*.12)*3), 4, color.RGBA{255, 222, 105, 235}, true)
			vector.StrokeCircle(screen, float32(drawPos.x), float32(drawPos.y), 49, 1, color.RGBA{117, 240, 255, 145}, true)
		}
	}
	if g.dragging {
		dragged := dragPosition(g.dragOrigin, g.dragNow)
		pull := vec{g.dragOrigin.x - dragged.x, g.dragOrigin.y - dragged.y}
		vector.StrokeLine(screen, float32(g.dragOrigin.x), float32(g.dragOrigin.y), float32(dragged.x), float32(dragged.y), 5, color.RGBA{246, 184, 64, 255}, false)
		vector.StrokeLine(screen, float32(dragged.x), float32(dragged.y), float32(dragged.x+pull.x), float32(dragged.y+pull.y), 2, color.White, false)
		for i := 1; i <= 6; i++ {
			t := float64(i) * .55
			vector.DrawFilledCircle(screen, float32(dragged.x+pull.x*.105*t), float32(dragged.y+pull.y*.105*t), 3, color.RGBA{255, 255, 255, 150}, true)
		}
		// A translucent landing reticle makes the slingshot payoff obvious.
		landX := dragged.x + pull.x*.72
		landY := dragged.y + pull.y*.72
		vector.StrokeCircle(screen, float32(landX), float32(landY), 24, 3, color.RGBA{255, 229, 126, 170}, true)
	}
	vector.DrawFilledRect(screen, 15, 605, 450, 76, color.RGBA{4, 15, 31, 220}, true)
	vector.StrokeRect(screen, 15, 605, 450, 76, 2, color.RGBA{255, 226, 150, 140}, true)
	drawCenteredStrikeLabel(screen, "DRAG THE GLOWING HERO BACKWARD — RELEASE TO STRIKE", 240, 620, strikeFace14, color.RGBA{255, 240, 202, 255})
	drawCenteredStrikeLabel(screen, "HIT YOUR PARTNER FOR A 165px ALLY WAVE", 240, 645, strikeFace14, color.RGBA{126, 242, 255, 255})
	drawCenteredStrikeLabel(screen, fmt.Sprintf("FRICTION %.3f   •   STOP %.2f", friction, stopSpeed), 240, 668, strikeFace14, color.RGBA{206, 220, 232, 255})
	if g.tick < 150 && g.stage == 1 {
		alpha := uint8(235)
		if g.tick > 105 {
			alpha = uint8(max(0, 235-(g.tick-105)*5))
		}
		vector.DrawFilledRect(screen, 48, 96, 384, 62, color.RGBA{4, 17, 33, alpha}, true)
		vector.StrokeRect(screen, 48, 96, 384, 62, 3, color.RGBA{255, 225, 126, alpha}, true)
		drawCenteredStrikeLabel(screen, "PULL EBI BACK  •  HIT TARGETS  •  BOUNCE INTO MOMO", 240, 111, strikeFace14, color.RGBA{255, 244, 210, alpha})
		drawCenteredStrikeLabel(screen, "CLEAR THE REEF IN 8 TURNS", 240, 135, strikeFace16, color.RGBA{128, 242, 255, alpha})
	}
	g.drawEffectBadge(screen)
	if g.won {
		msg := "STAGE CLEAR!\n\nTAP / ENTER FOR NEXT STAGE"
		if g.stage == 3 {
			msg = "REEF RESCUED!\n\nTAP / ENTER FOR A NEW RUN"
		}
		overlay(screen, msg)
	}
	if g.lost {
		overlay(screen, "OUT OF TURNS\n\nTAP / ENTER TO RETRY")
	}
}

func loadStrikeArt() {
	strikeArtOnce.Do(func() {
		strikeArt = make(map[string]*ebiten.Image, 4)
		for key, filename := range map[string]string{
			"arena":     "strike-pearl-coliseum.png",
			"allies":    "strike-allies-atlas.png",
			"enemies":   "strike-enemies-atlas.png",
			"obstacles": "strike-obstacles-atlas.png",
		} {
			data, err := strikeArtFS.ReadFile("assets/" + filename)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			strikeArt[key] = ebiten.NewImageFromImage(decoded)
		}
		allyAtlas := strikeArt["allies"]
		for i := range strikeAllies {
			strikeAllies[i] = ebiten.NewImageFromImage(allyAtlas.SubImage(image.Rect(i*512, 0, (i+1)*512, 512)))
		}
		enemyAtlas := strikeArt["enemies"]
		obstacleAtlas := strikeArt["obstacles"]
		for i := range strikeEnemies {
			strikeEnemies[i] = ebiten.NewImageFromImage(enemyAtlas.SubImage(image.Rect(i*512, 0, (i+1)*512, 512)))
			strikeObstacles[i] = ebiten.NewImageFromImage(obstacleAtlas.SubImage(image.Rect(i*512, 0, (i+1)*512, 512)))
		}
		strikeFace14, _ = uilab.Face("en", 14)
		strikeFace16, _ = uilab.Face("en", 16)
		strikeFace20, _ = uilab.Face("en", 20)
	})
}

func drawStrikeCover(dst, img *ebiten.Image, x, y, w, h float64) {
	b := img.Bounds()
	scale := math.Max(w/float64(b.Dx()), h/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Min.X), -float64(b.Min.Y))
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-float64(b.Dx())*scale)/2, y+(h-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawStrikeSprite(dst, img *ebiten.Image, centerX, centerY, size, angle float64, alpha float32, flash bool) {
	b := img.Bounds()
	scale := size / float64(max(b.Dx(), b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Min.X)-float64(b.Dx())/2, -float64(b.Min.Y)-float64(b.Dy())/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(centerX, centerY)
	if flash {
		op.ColorScale.Scale(1, .25, .25, alpha)
	} else {
		op.ColorScale.ScaleAlpha(alpha)
	}
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawStrikeLabel(dst *ebiten.Image, label string, x, y float64, face *text.GoTextFace, c color.Color) {
	if face == nil {
		ebitenutil.DebugPrintAt(dst, label, int(x), int(y))
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, label, face, op)
}

func drawCenteredStrikeLabel(dst *ebiten.Image, label string, centerX, y float64, face *text.GoTextFace, c color.Color) {
	if face == nil {
		ebitenutil.DebugPrintAt(dst, label, int(centerX)-len(label)*3, int(y))
		return
	}
	w, _ := text.Measure(label, face, 0)
	drawStrikeLabel(dst, label, centerX-w/2, y, face, c)
}

func (g *game) drawTitle(screen *ebiten.Image) {
	const label = "EBI STRIKE / REEF RESCUE"
	drawCenteredStrikeLabel(screen, label, 240, 13, strikeFace16, color.RGBA{255, 234, 188, 255})
}

func (g *game) drawEffectBadge(screen *ebiten.Image) {
	if g.shader == nil || !g.shader.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.shader.Draw(fx, g.badge, float32(g.turns)*.2) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenW-34, 12)
	screen.DrawImage(fx, op)
}

func (g *game) alive() int {
	n := 0
	for _, e := range g.enemies {
		if e.hp > 0 {
			n++
		}
	}
	return n
}

func distance(a, b vec) float64 { return math.Hypot(a.x-b.x, a.y-b.y) }

func dragPosition(origin, pointer vec) vec {
	dx, dy := pointer.x-origin.x, pointer.y-origin.y
	if length := math.Hypot(dx, dy); length > 145 {
		dx, dy = dx*145/length, dy*145/length
	}
	return vec{
		x: min(452.0, max(28.0, origin.x+dx)),
		y: min(570.0, max(100.0, origin.y+dy)),
	}
}

func pressPosition() (int, int, bool) {
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

func retryPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	_, _, ok := pressPosition()
	return ok
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 45, 270, 390, 150, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 45, 270, 390, 150, 4, color.RGBA{244, 189, 68, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 120, 326)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Strike — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
