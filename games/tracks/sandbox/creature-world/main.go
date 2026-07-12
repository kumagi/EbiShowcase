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
	screenW   = 480
	screenH   = 720
	cycleTime = 15 * 60
	timeLimit = 60 * 60
)

type vec struct{ x, y float64 }

type creatureState int

const (
	wander creatureState = iota
	curious
	chase
	flee
	friend
)

var stateNames = [...]string{"WANDER", "CURIOUS", "CHASE", "FLEE", "FRIEND"}

type creature struct {
	kind     int // 0: day grazer, 1: night prowler
	state    creatureState
	pos, dir vec
	timer    int
}

type berry struct {
	pos       vec
	collected bool
}

type game struct {
	player                   vec
	creatures                []creature
	berries                  []berry
	rng                      *rand.Rand
	frames, hp, stock, torch int
	friends, hurtCooldown    int
	message                  string
	won, lost                bool
}

func newGame() *game {
	return &game{
		player: vec{240, 345},
		berries: []berry{
			{pos: vec{70, 165}}, {pos: vec{155, 250}}, {pos: vec{90, 440}},
			{pos: vec{385, 145}}, {pos: vec{330, 315}}, {pos: vec{390, 490}},
		},
		rng:     rand.New(rand.NewSource(7606)),
		hp:      5,
		torch:   4,
		message: "Collect berries. Day grazers become curious nearby.",
	}
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.hurtCooldown > 0 {
		g.hurtCooldown--
	}
	if g.frames >= timeLimit {
		g.lost = true
		g.message = "The survey ended before four friendships formed."
		return nil
	}

	dx, dy, feed, useTorch := readInput(g.player)
	if dx != 0 || dy != 0 {
		length := math.Hypot(dx, dy)
		g.player.x += dx / length * 3.2
		g.player.y += dy / length * 3.2
		g.player.x = clamp(g.player.x, 24, 456)
		g.player.y = clamp(g.player.y, 112, 565)
	}
	if feed {
		g.feedNearest()
	}
	if useTorch {
		g.flashTorch()
	}
	g.collectBerries()
	g.spawnCreatures()
	g.updateCreatures()
	if g.hp <= 0 {
		g.lost = true
		g.message = "The night prowlers exhausted the explorer."
	}
	if g.friends >= 4 && g.frames >= cycleTime*2 {
		g.won = true
		g.message = "Four friends observed through a full day and night!"
	}
	return nil
}

func (g *game) isDay() bool { return (g.frames/cycleTime)%2 == 0 }

func (g *game) spawnCreatures() {
	if g.frames%180 != 1 || len(g.creatures) >= 10 {
		return
	}
	if g.isDay() {
		// Grazers spawn only in the green meadow on the left.
		g.creatures = append(g.creatures, creature{
			kind: 0, state: wander,
			pos: vec{45 + g.rng.Float64()*150, 135 + g.rng.Float64()*390},
			dir: randomDirection(g.rng), timer: 90,
		})
		g.message = "DAY spawn: a grazer appeared on meadow terrain."
	} else {
		// Prowlers spawn only in the rocky night habitat on the right.
		g.creatures = append(g.creatures, creature{
			kind: 1, state: wander,
			pos: vec{300 + g.rng.Float64()*135, 135 + g.rng.Float64()*390},
			dir: randomDirection(g.rng), timer: 90,
		})
		g.message = "NIGHT spawn: a prowler emerged from rocky terrain."
	}
}

func (g *game) updateCreatures() {
	for i := range g.creatures {
		c := &g.creatures[i]
		d := distance(c.pos, g.player)
		if c.kind == 0 {
			switch c.state {
			case wander:
				if d < 115 {
					c.state = curious
				}
			case curious:
				if d > 170 {
					c.state = wander
				}
			case friend:
				c.dir = direction(c.pos, g.player)
				c.pos.x += c.dir.x * 1.25
				c.pos.y += c.dir.y * 1.25
				continue
			}
		} else {
			switch c.state {
			case flee:
				c.timer--
				c.dir = direction(g.player, c.pos)
				if c.timer <= 0 {
					c.state = wander
				}
			case wander:
				if d < 210 {
					c.state = chase
				}
			case chase:
				c.dir = direction(c.pos, g.player)
				if d > 260 {
					c.state = wander
				}
				if d < 25 && g.hurtCooldown == 0 {
					g.hp--
					g.hurtCooldown = 75
					c.state = flee
					c.timer = 90
					g.message = "A chasing prowler touched you: HP -1!"
				}
			}
		}

		if c.state == wander {
			c.timer--
			if c.timer <= 0 {
				c.dir = randomDirection(g.rng)
				c.timer = 60 + g.rng.Intn(100)
			}
		}
		speed := 0.75
		if c.state == curious {
			c.dir = direction(c.pos, g.player)
			speed = 1.0
		}
		if c.state == chase || c.state == flee {
			speed = 1.5
		}
		c.pos.x = clamp(c.pos.x+c.dir.x*speed, 25, 455)
		c.pos.y = clamp(c.pos.y+c.dir.y*speed, 112, 565)
	}
}

func (g *game) collectBerries() {
	for i := range g.berries {
		b := &g.berries[i]
		if !b.collected && distance(b.pos, g.player) < 22 {
			b.collected = true
			g.stock++
			g.message = fmt.Sprintf("Berry collected. Food stock %d.", g.stock)
		}
	}
}

func (g *game) feedNearest() {
	if g.stock == 0 {
		g.message = "No berries: walk over a pink berry bush first."
		return
	}
	best, bestDistance := -1, 70.0
	for i := range g.creatures {
		c := &g.creatures[i]
		if c.kind == 0 && c.state == curious {
			if d := distance(c.pos, g.player); d < bestDistance {
				best, bestDistance = i, d
			}
		}
	}
	if best < 0 {
		g.message = "No curious grazer close enough to feed."
		return
	}
	g.stock--
	g.friends++
	g.creatures[best].state = friend
	g.message = fmt.Sprintf("Grazer befriended! %d/4 now follows you.", g.friends)
}

func (g *game) flashTorch() {
	if g.torch == 0 {
		g.message = "The torch has no flashes left."
		return
	}
	g.torch--
	scared := 0
	for i := range g.creatures {
		c := &g.creatures[i]
		if c.kind == 1 && distance(c.pos, g.player) < 180 {
			c.state = flee
			c.timer = 150
			scared++
		}
	}
	g.message = fmt.Sprintf("Torch flash: %d prowler(s) changed to FLEE.", scared)
}

func (g *game) Draw(screen *ebiten.Image) {
	day := g.isDay()
	background := color.RGBA{39, 91, 76, 255}
	if !day {
		background = color.RGBA{19, 37, 62, 255}
	}
	screen.Fill(color.RGBA{8, 18, 33, 255})
	phase := "DAY"
	if !day {
		phase = "NIGHT"
	}
	phaseLeft := (cycleTime - g.frames%cycleTime + 59) / 60
	ebitenutil.DebugPrintAt(screen, "LIVING WORLD FIELD STUDY", 156, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s %02ds   FRIENDS %d/4   HP %d   BERRIES %d   TORCH %d", phase, phaseLeft, g.friends, g.hp, g.stock, g.torch), 59, 45)
	ebitenutil.DebugPrintAt(screen, g.message, 62, 70)
	vector.DrawFilledRect(screen, 15, 92, 450, 490, background, false)
	vector.DrawFilledRect(screen, 15, 92, 220, 490, color.RGBA{66, 132, 78, 120}, false)
	vector.DrawFilledRect(screen, 235, 92, 230, 490, color.RGBA{80, 82, 103, 120}, false)
	ebitenutil.DebugPrintAt(screen, "MEADOW / DAY GRAZER SPAWN", 25, 102)
	ebitenutil.DebugPrintAt(screen, "ROCK / NIGHT PROWLER SPAWN", 267, 102)

	for _, b := range g.berries {
		if !b.collected {
			vector.DrawFilledCircle(screen, float32(b.pos.x), float32(b.pos.y), 8, color.RGBA{238, 99, 157, 255}, false)
		}
	}
	for _, c := range g.creatures {
		cc := color.RGBA{235, 200, 83, 255}
		if c.kind == 1 {
			cc = color.RGBA{166, 91, 205, 255}
		}
		if c.state == friend {
			cc = color.RGBA{88, 207, 147, 255}
		}
		vector.DrawFilledCircle(screen, float32(c.pos.x), float32(c.pos.y), 13, cc, false)
		vector.StrokeCircle(screen, float32(c.pos.x), float32(c.pos.y), 13, 2, color.White, false)
		ebitenutil.DebugPrintAt(screen, stateNames[c.state], int(c.pos.x)-22, int(c.pos.y)-28)
	}
	playerColor := color.RGBA{71, 184, 218, 255}
	if g.hurtCooldown > 0 && g.hurtCooldown%10 < 5 {
		playerColor = color.RGBA{238, 93, 83, 255}
	}
	vector.DrawFilledCircle(screen, float32(g.player.x), float32(g.player.y), 14, playerColor, false)
	vector.StrokeCircle(screen, float32(g.player.x), float32(g.player.y), 14, 3, color.White, false)

	labels := [...]string{"LEFT", "RIGHT", "UP", "DOWN", "FEED", "TORCH"}
	for i, label := range labels {
		c := color.RGBA{52, 83, 119, 255}
		if i >= 4 {
			c = color.RGBA{184, 111, 61, 255}
		}
		vector.DrawFilledRect(screen, float32(i*80+3), 602, 74, 66, c, false)
		ebitenutil.DebugPrintAt(screen, label, i*80+18, 631)
	}
	ebitenutil.DebugPrintAt(screen, "WASD / ARROWS   E: FEED   SPACE: TORCH", 91, 687)
	if g.won {
		overlay(screen, "ECOSYSTEM STUDY CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "FIELD STUDY FAILED\n\nTAP / ENTER TO RETRY")
	}
}

func readInput(player vec) (float64, float64, bool, bool) {
	dx, dy := 0.0, 0.0
	feed := inpututil.IsKeyJustPressed(ebiten.KeyE)
	torch := inpututil.IsKeyJustPressed(ebiten.KeySpace)
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
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 600 {
			if x < 320 {
				dx, dy, _, _ = actionForButton(x)
			}
		} else {
			dx, dy = float64(x)-player.x, float64(y)-player.y
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 600 {
			if x < 320 {
				dx, dy, _, _ = actionForButton(x)
			}
		} else {
			dx, dy = float64(x)-player.x, float64(y)-player.y
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 600 && x >= 320 {
			_, _, feed, torch = actionForButton(x)
		}
	}
	justTouches := inpututil.AppendJustPressedTouchIDs(nil)
	if len(justTouches) > 0 {
		x, y := ebiten.TouchPosition(justTouches[0])
		if y >= 600 && x >= 320 {
			_, _, feed, torch = actionForButton(x)
		}
	}
	return dx, dy, feed, torch
}

func actionForButton(x int) (float64, float64, bool, bool) {
	switch x / 80 {
	case 0:
		return -1, 0, false, false
	case 1:
		return 1, 0, false, false
	case 2:
		return 0, -1, false, false
	case 3:
		return 0, 1, false, false
	case 4:
		return 0, 0, true, false
	default:
		return 0, 0, false, true
	}
}

func randomDirection(rng *rand.Rand) vec {
	a := rng.Float64() * math.Pi * 2
	return vec{math.Cos(a), math.Sin(a)}
}

func direction(from, to vec) vec {
	dx, dy := to.x-from.x, to.y-from.y
	d := math.Hypot(dx, dy)
	if d == 0 {
		return vec{}
	}
	return vec{dx / d, dy / d}
}

func distance(a, b vec) float64       { return math.Hypot(a.x-b.x, a.y-b.y) }
func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 105, 327)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Living World Field Study — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
