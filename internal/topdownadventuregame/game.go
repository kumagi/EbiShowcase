// Package topdownadventuregame is the shared Ebitengine presentation used by
// the eight progressively larger top-down adventure lessons.
package topdownadventuregame

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/topdownadventurelogic"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const (
	W = 480
	H = 720
)

type enemy struct {
	pos          topdownadventurelogic.Vec
	hp, cooldown int
	boss         bool
}
type node struct {
	pos  topdownadventurelogic.Vec
	kind int
	on   bool
}
type particle struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}
type game struct {
	lesson, stage, frames, attack, score, best, keys, tool, door, transition, shake, flash int
	player                                                                                 topdownadventurelogic.Fighter
	face                                                                                   topdownadventurelogic.Vec
	enemies                                                                                []enemy
	gems                                                                                   []topdownadventurelogic.Vec
	nodes                                                                                  []node
	chest, clear, over                                                                     bool
	message                                                                                string
	particles                                                                              []particle
	rng                                                                                    *rand.Rand
}

var sessionBest [9]int
var titles = [9]string{"", "EIGHT-WAY RELIC RUN", "SWORD REACH TRIAL", "HURT & RECOVERY", "THREE SEALED ROOMS", "KEY AND TREASURE", "RELIC TOOL PUZZLES", "THREE-PHASE GUARDIAN", "EBI RELIC DUNGEON"}

func Run(lesson int) {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle(titles[lesson])
	if err := ebiten.RunGame(newGame(lesson)); err != nil {
		panic(err)
	}
}
func newGame(lesson int) *game {
	g := &game{lesson: lesson, face: topdownadventurelogic.Vec{X: 1}, player: topdownadventurelogic.Fighter{Pos: topdownadventurelogic.Vec{X: 75, Y: 320}, HP: 5}, best: sessionBest[lesson], rng: rand.New(rand.NewSource(int64(100 + lesson)))}
	g.setup()
	return g
}
func (g *game) setup() {
	g.player.Pos = topdownadventurelogic.Vec{X: 75, Y: 320}
	g.enemies = nil
	g.gems = nil
	g.nodes = nil
	g.transition = 0
	switch g.lesson {
	case 1:
		g.gems = []topdownadventurelogic.Vec{{120, 160}, {230, 130}, {370, 190}, {150, 390}, {280, 340}, {400, 450}}
		g.message = "Move in eight directions and collect all six star relics."
	case 2:
		g.spawnEnemies(4, false)
		g.message = "Face a target, then use the visible sword rectangle."
	case 3:
		g.spawnEnemies(5, false)
		g.message = "Take a hit, watch the flash, then counter during recovery."
	case 4:
		g.startRoom()
	case 5:
		g.gems = []topdownadventurelogic.Vec{{145, 185}}
		g.message = "Collect the key, open the blue door, then reach the chest."
	case 6:
		g.nodes = []node{{topdownadventurelogic.Vec{120, 180}, 0, false}, {topdownadventurelogic.Vec{360, 190}, 1, false}, {topdownadventurelogic.Vec{240, 430}, 2, false}}
		g.message = "Match SWORD, GUST, and LAMP to the three relic seals."
	case 7:
		g.spawnBoss()
		g.message = "Read GUARD, DASH, and STORM phases; strike after moving."
	case 8:
		g.message = "ROOM 1: find the key and open the treasure gate."
		g.gems = []topdownadventurelogic.Vec{{145, 185}}
	}
}
func (g *game) spawnEnemies(n int, boss bool) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		hp := 1
		if g.lesson >= 3 {
			hp = 2
		}
		g.enemies = append(g.enemies, enemy{pos: topdownadventurelogic.Vec{X: 270 + math.Cos(a)*130, Y: 315 + math.Sin(a)*150}, hp: hp, boss: boss})
	}
}
func (g *game) spawnBoss() {
	g.enemies = []enemy{{pos: topdownadventurelogic.Vec{X: 345, Y: 285}, hp: 12, boss: true}}
}
func (g *game) startRoom() {
	g.enemies = nil
	g.spawnEnemies(2+g.stage, false)
	g.message = fmt.Sprintf("ROOM %d/3: defeat every crawler to unseal the exit.", g.stage+1)
}
func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame(g.lesson)
		}
		return nil
	}
	g.frames++
	g.readMovement()
	if attackPressed() && g.attack == 0 {
		g.attack = 18
	}
	if toolPressed() {
		g.tool = (g.tool + 1) % 3
		g.message = "Tool selected: " + []string{"SWORD", "GUST", "LAMP"}[g.tool]
	}
	if g.attack > 0 {
		g.attack--
		if g.attack == 10 {
			g.resolveAction()
		}
	}
	g.player.Tick()
	g.player.Pos.X = max(48, min(432, g.player.Pos.X))
	g.player.Pos.Y = max(108, min(520, g.player.Pos.Y))
	g.collect()
	g.updateEnemies()
	g.updateLesson()
	g.updateParticles()
	if g.shake > 0 {
		g.shake--
	}
	if g.flash > 0 {
		g.flash--
	}
	return nil
}
func (g *game) readMovement() {
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		dy++
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 590 {
			switch min(5, x/80) {
			case 0:
				dx--
			case 1:
				dy--
			case 2:
				dy++
			case 3:
				dx++
			}
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 590 {
			switch min(5, x/80) {
			case 0:
				dx--
			case 1:
				dy--
			case 2:
				dy++
			case 3:
				dx++
			}
		}
	}
	v := topdownadventurelogic.Normalize(topdownadventurelogic.Vec{X: dx, Y: dy})
	if v.X != 0 || v.Y != 0 {
		g.face = v
		if g.attack == 0 {
			g.player.Pos.X += v.X * 3
			g.player.Pos.Y += v.Y * 3
		}
	}
}
func (g *game) resolveAction() {
	box := topdownadventurelogic.AttackBox(g.player.Pos, g.face, 58, 44)
	hit := false
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp <= 0 || e.cooldown > 0 {
			continue
		}
		r := topdownadventurelogic.Rect{X: e.pos.X - 18, Y: e.pos.Y - 18, W: 36, H: 36}
		if box.Intersects(r) {
			e.hp--
			e.cooldown = 24
			hit = true
			g.score += 100
			d := topdownadventurelogic.Normalize(topdownadventurelogic.Vec{X: e.pos.X - g.player.Pos.X, Y: e.pos.Y - g.player.Pos.Y})
			e.pos.X += d.X * 22
			e.pos.Y += d.Y * 22
			g.burst(e.pos, color.RGBA{255, 118, 102, 255}, 12)
		}
	}
	for i := range g.nodes {
		n := &g.nodes[i]
		if n.on || n.kind != g.tool {
			continue
		}
		if math.Hypot(n.pos.X-g.player.Pos.X, n.pos.Y-g.player.Pos.Y) < 82 {
			n.on = true
			hit = true
			g.score += 180
			g.burst(n.pos, color.RGBA{80, 230, 203, 255}, 18)
			g.message = []string{"Sword cut the vine seal!", "Gust spun the wind seal!", "Lamp revealed the shadow seal!"}[n.kind]
		}
	}
	if hit {
		g.shake = 5
		g.flash = 3
	} else {
		g.message = "The action missed. Face the target and check the reach."
	}
}
func (g *game) collect() {
	for i := len(g.gems) - 1; i >= 0; i-- {
		p := g.gems[i]
		if math.Hypot(p.X-g.player.Pos.X, p.Y-g.player.Pos.Y) < 28 {
			g.gems = append(g.gems[:i], g.gems[i+1:]...)
			g.keys++
			g.score += 80
			g.burst(p, color.RGBA{255, 211, 82, 255}, 14)
			g.message = "Relic collected!"
		}
	}
}
func (g *game) updateEnemies() {
	alive := 0
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp <= 0 {
			continue
		}
		alive++
		if e.cooldown > 0 {
			e.cooldown--
		}
		speed := .6 + float64(g.lesson)*.06
		if e.boss {
			phase := topdownadventurelogic.PhaseForHP(e.hp, 12)
			speed = []float64{.65, 1.15, 1.45, 0}[phase]
			if phase == topdownadventurelogic.BossDash && g.frames%110 < 22 {
				speed = 2.8
			}
			if phase == topdownadventurelogic.BossStorm && g.frames%75 < 25 {
				speed = 2
			}
			if phase == topdownadventurelogic.BossStorm && g.frames%90 == 0 && math.Hypot(e.pos.X-g.player.Pos.X, e.pos.Y-g.player.Pos.Y) < 155 {
				if g.player.Hurt(1, e.pos, 75) {
					g.shake, g.flash = 14, 12
					g.burst(g.player.Pos, color.RGBA{147, 112, 255, 255}, 22)
					g.message = "STORM ring hit! Move outside its warning circle."
				}
			}
		}
		d := topdownadventurelogic.Normalize(topdownadventurelogic.Vec{X: g.player.Pos.X - e.pos.X, Y: g.player.Pos.Y - e.pos.Y})
		e.pos.X += d.X * speed
		e.pos.Y += d.Y * speed
		if g.lesson >= 3 && math.Hypot(e.pos.X-g.player.Pos.X, e.pos.Y-g.player.Pos.Y) < 34 {
			if g.player.Hurt(1, e.pos, 75) {
				g.shake = 12
				g.flash = 12
				g.burst(g.player.Pos, color.RGBA{255, 70, 96, 255}, 18)
				g.message = "Hit! Flashing means invulnerable—move to safety."
			}
		}
	}
	if alive == 0 && g.lesson >= 2 && g.lesson != 5 && g.lesson != 6 {
		g.transition++
	}
}
func (g *game) updateLesson() {
	switch g.lesson {
	case 1:
		if len(g.gems) == 0 {
			g.win("Eight-way movement mastered!")
		}
	case 2, 3:
		if g.transition > 35 {
			g.win("All targets cleared. Retry for a faster score!")
		}
	case 4:
		if g.transition > 35 {
			g.stage++
			g.transition = 0
			if g.stage >= 3 {
				g.win("All three rooms unsealed!")
			} else {
				g.startRoom()
			}
		}
	case 5:
		g.updateKeyDoor(false)
	case 6:
		if allNodes(g.nodes) {
			g.win("All three tool seals solved!")
		}
	case 7:
		if g.transition > 45 {
			g.win("The three-phase guardian is defeated!")
		}
	case 8:
		g.updateDungeon()
	}
	if g.player.HP <= 0 {
		g.over = true
		g.message = "The expedition ended. Tap or Enter to retry."
	}
}
func (g *game) updateKeyDoor(dungeon bool) {
	if g.keys > 0 && g.door == 0 && g.player.Pos.X > 275 {
		g.door = 1
		g.keys--
		g.flash = 6
		g.message = "The key opened the blue door!"
	}
	if g.door == 0 && g.player.Pos.X > 285 {
		g.player.Pos.X = 285
	}
	if g.door == 1 && g.player.Pos.X > 392 {
		g.chest = true
		g.score += 400
		g.burst(topdownadventurelogic.Vec{410, 320}, color.RGBA{255, 211, 82, 255}, 24)
		if dungeon {
			g.advanceDungeon()
		} else {
			g.win("Treasure opened—key and door state connected!")
		}
	}
}
func (g *game) updateDungeon() {
	switch g.stage {
	case 0:
		g.updateKeyDoor(true)
	case 1:
		if g.transition > 35 {
			g.advanceDungeon()
		}
	case 2:
		if allNodes(g.nodes) {
			g.advanceDungeon()
		}
	case 3:
		if g.transition > 45 {
			g.win("The Relic Dungeon is safe again!")
		}
	}
}
func (g *game) advanceDungeon() {
	g.stage++
	g.transition = 0
	g.keys = 0
	g.door = 0
	g.chest = false
	g.enemies = nil
	g.gems = nil
	g.nodes = nil
	g.player.Pos = topdownadventurelogic.Vec{X: 70, Y: 320}
	switch g.stage {
	case 1:
		g.spawnEnemies(4, false)
		g.message = "ROOM 2: defeat every crawler to open the seal."
	case 2:
		g.nodes = []node{{topdownadventurelogic.Vec{120, 180}, 0, false}, {topdownadventurelogic.Vec{360, 190}, 1, false}, {topdownadventurelogic.Vec{240, 430}, 2, false}}
		g.message = "ROOM 3: match each tool to its colored seal."
	case 3:
		g.spawnBoss()
		g.message = "FINAL ROOM: read all three guardian phases."
	}
}
func (g *game) win(message string) {
	g.clear = true
	g.message = message
	g.score += g.player.HP*100 + max(0, 1800-g.frames/3)
	if g.score > sessionBest[g.lesson] {
		sessionBest[g.lesson] = g.score
	}
	g.best = sessionBest[g.lesson]
}
func allNodes(n []node) bool {
	if len(n) == 0 {
		return false
	}
	for _, v := range n {
		if !v.on {
			return false
		}
	}
	return true
}
func (g *game) burst(p topdownadventurelogic.Vec, c color.RGBA, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * .7
		g.particles = append(g.particles, particle{p.X, p.Y, math.Cos(a) * 2.5, math.Sin(a)*2.5 - 1, 32, c})
	}
}
func (g *game) updateParticles() {
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .05
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}
}
func (g *game) Draw(screen *ebiten.Image) {
	palette := []color.RGBA{{}, {18, 55, 65, 255}, {46, 42, 71, 255}, {70, 38, 53, 255}, {33, 56, 78, 255}, {52, 54, 43, 255}, {35, 48, 75, 255}, {65, 31, 52, 255}, {23, 48, 62, 255}}
	screen.Fill(palette[g.lesson])
	dx, dy := 0, 0
	if g.shake > 0 {
		dx = g.rng.Intn(7) - 3
		dy = g.rng.Intn(7) - 3
	}
	world := ebiten.NewImage(W, H)
	g.drawArena(world)
	g.drawActors(world)
	g.drawUI(world)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	screen.DrawImage(world, op)
	if g.flash > 0 {
		vector.DrawFilledRect(screen, 0, 0, W, H, color.RGBA{255, 255, 255, 60}, false)
	}
	if g.clear || g.over {
		vector.DrawFilledRect(screen, 35, 225, 410, 235, color.RGBA{4, 11, 27, 246}, false)
		title := "TRIAL COMPLETE!"
		if g.over {
			title = "EXPEDITION LOST"
		}
		ebitenutil.DebugPrintAt(screen, title, 165, 270)
		ebitenutil.DebugPrintAt(screen, g.message, 65, 315)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %05d  BEST %05d", g.score, g.best), 145, 355)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO PLAY AGAIN", 135, 415)
	}
}
func (g *game) drawArena(s *ebiten.Image) {
	vector.DrawFilledRect(s, 30, 90, 420, 450, color.RGBA{12, 25, 40, 210}, false)
	for x := 48; x < 450; x += 48 {
		vector.StrokeLine(s, float32(x), 90, float32(x), 540, 1, color.RGBA{255, 255, 255, 14}, false)
	}
	for y := 108; y < 540; y += 48 {
		vector.StrokeLine(s, 30, float32(y), 450, float32(y), 1, color.RGBA{255, 255, 255, 14}, false)
	}
	for _, p := range g.gems {
		trackatlas.DrawCentered(s, "power-star", p.X, p.Y, 26)
	}
	if g.lesson == 5 || (g.lesson == 8 && g.stage == 0) {
		c := color.RGBA{63, 104, 165, 255}
		if g.door > 0 {
			c = color.RGBA{45, 73, 91, 90}
		}
		vector.DrawFilledRect(s, 300, 90, 24, 450, c, false)
		trackatlas.DrawCentered(s, "tile-crate", 410, 320, 44)
	}
	for _, n := range g.nodes {
		c := []color.RGBA{{224, 91, 91, 255}, {87, 202, 214, 255}, {242, 203, 79, 255}}[n.kind]
		if n.on {
			c = color.RGBA{65, 230, 166, 255}
		}
		vector.DrawFilledCircle(s, float32(n.pos.X), float32(n.pos.Y), 28, c, false)
		ebitenutil.DebugPrintAt(s, []string{"SWORD", "GUST", "LAMP"}[n.kind], int(n.pos.X)-20, int(n.pos.Y)-4)
	}
	for _, p := range g.particles {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), 3, p.c, false)
	}
}
func (g *game) drawActors(s *ebiten.Image) {
	if g.attack > 0 {
		b := topdownadventurelogic.AttackBox(g.player.Pos, g.face, 58, 44)
		vector.DrawFilledRect(s, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), color.RGBA{255, 220, 105, 90}, false)
	}
	bob := math.Sin(float64(g.frames)*.2) * 2
	alpha := float32(1)
	if g.player.Invulnerable > 0 && g.player.Invulnerable%10 < 5 {
		alpha = .25
	}
	trackatlas.DrawTinted(s, "hero", g.player.Pos.X, g.player.Pos.Y+bob, 34, 1, 1, 1, alpha)
	for _, e := range g.enemies {
		if e.hp <= 0 {
			continue
		}
		sprite := "slug"
		size := 34.0
		if e.boss {
			sprite = "boss-crab"
			size = 66
		}
		trackatlas.DrawCentered(s, sprite, e.pos.X, e.pos.Y+math.Sin(float64(g.frames)*.15+e.pos.X)*3, size)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("HP%d", e.hp), int(e.pos.X)-12, int(e.pos.Y)-35)
		if e.boss {
			phase := topdownadventurelogic.PhaseForHP(e.hp, 12)
			if phase == topdownadventurelogic.BossStorm {
				radius := float32(70 + (g.frames%90)*85/90)
				vector.StrokeCircle(s, float32(e.pos.X), float32(e.pos.Y), radius, 4, color.RGBA{170, 125, 255, 170}, false)
			}
			ebitenutil.DebugPrintAt(s, []string{"GUARD", "DASH", "STORM", "DOWN"}[phase], int(e.pos.X)-24, int(e.pos.Y)-52)
		}
	}
}
func (g *game) drawUI(s *ebiten.Image) {
	ebitenutil.DebugPrintAt(s, titles[g.lesson], 130, 16)
	room := g.stage + 1
	if g.lesson < 4 {
		room = 1
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HP %d  SCORE %05d  ROOM %d  TOOL %s", g.player.HP, g.score, room, []string{"SWORD", "GUST", "LAMP"}[g.tool]), 95, 44)
	ebitenutil.DebugPrintAt(s, g.message, 36, 68)
	labels := []string{"LEFT", "UP", "DOWN", "RIGHT", "ATTACK", "TOOL"}
	for i, l := range labels {
		c := color.RGBA{48, 82, 120, 255}
		if i == 4 {
			c = color.RGBA{190, 78, 74, 255}
		} else if i == 5 {
			c = color.RGBA{178, 137, 50, 255}
		}
		vector.DrawFilledRect(s, float32(i*80+3), 590, 74, 70, c, false)
		ebitenutil.DebugPrintAt(s, l, i*80+14, 621)
	}
	ebitenutil.DebugPrintAt(s, "WASD/arrows move · X/Space attack · Q changes tool", 70, 685)
}
func attackPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return y >= 590 && x >= 320 && x < 400
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 590 && x >= 320 && x < 400 {
			return true
		}
	}
	return false
}
func toolPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return y >= 590 && x >= 400
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 590 && x >= 400 {
			return true
		}
	}
	return false
}
func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(int, int) (int, int) { return W, H }
