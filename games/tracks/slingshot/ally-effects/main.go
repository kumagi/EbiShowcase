package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width, height = 480, 720
	playerRadius  = 18.0
	allyRadius    = 27.0
	enemyRadius   = 24.0
	maxHP         = 5
	timeLimit     = 50 * 60
	effectLaser   = 0
	effectHeal    = 1
	effectPulse   = 2
)

type ability struct {
	name  string
	kind  int
	power int
	color color.RGBA
}

var abilities = []ability{
	{"LASER", effectLaser, 2, color.RGBA{87, 204, 229, 255}},
	{"HEAL", effectHeal, 2, color.RGBA{91, 202, 121, 255}},
	{"PULSE", effectPulse, 1, color.RGBA{235, 173, 66, 255}},
}

type ally struct {
	x, y        float64
	ability     int
	cooldown    int
	wasTouching bool
}

type enemy struct {
	x, y     float64
	hp       int
	cooldown int
}

type game struct {
	x, y, vx, vy       float64
	allies             []ally
	enemies            []enemy
	hp, frames, events int
	damageCooldown     int
	effectTimer        int
	effectKind         int
	effectX, effectY   float64
	clear, over        bool
	message            string
}

func newGame() *game {
	return &game{
		x: 240, y: 520, hp: maxHP,
		allies:  []ally{{95, 190, 0, 0, false}, {380, 240, 1, 0, false}, {220, 405, 2, 0, false}},
		enemies: []enemy{{245, 145, 4, 0}, {95, 420, 4, 0}, {390, 455, 4, 0}},
		message: "Touch allies: LASER, HEAL, and PULSE run different data.",
	}
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.frames >= timeLimit {
		g.over = true
		g.message = "Time up. Trigger effects more directly!"
		return nil
	}
	if g.damageCooldown > 0 {
		g.damageCooldown--
	}
	if g.effectTimer > 0 {
		g.effectTimer--
	}

	ax, ay := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		ax--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		ax++
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		ay--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		ay++
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		px, py := ebiten.CursorPosition()
		ax, ay = float64(px)-g.x, float64(py)-g.y
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		px, py := ebiten.TouchPosition(ids[0])
		ax, ay = float64(px)-g.x, float64(py)-g.y
	}
	if length := math.Hypot(ax, ay); length > 0 {
		g.vx += ax / length * 0.38
		g.vy += ay / length * 0.38
	}
	speed := math.Hypot(g.vx, g.vy)
	if speed > 8.5 {
		g.vx, g.vy = g.vx*8.5/speed, g.vy*8.5/speed
	}
	g.vx *= 0.986
	g.vy *= 0.986
	g.x += g.vx
	g.y += g.vy
	g.bounceWalls()

	for i := range g.allies {
		a := &g.allies[i]
		if a.cooldown > 0 {
			a.cooldown--
		}
		nowTouching := distance(g.x, g.y, a.x, a.y) <= playerRadius+allyRadius
		if nowTouching && !a.wasTouching && a.cooldown == 0 {
			g.activate(a)
			a.cooldown = 75
		}
		a.wasTouching = nowTouching
	}

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
		if distance(g.x, g.y, e.x, e.y) <= playerRadius+enemyRadius {
			g.pushFrom(e.x, e.y)
			if g.damageCooldown == 0 {
				g.hp--
				g.damageCooldown = 70
				g.message = "The coral guard hit back. Find the HEAL ally!"
				if g.hp <= 0 {
					g.over = true
					return nil
				}
			}
		}
	}
	if alive == 0 {
		g.clear = true
		g.message = "All guards cleared by ally contact events!"
	}
	return nil
}

func (g *game) activate(a *ally) {
	data := abilities[a.ability]
	g.events++
	g.effectKind, g.effectTimer = data.kind, 28
	g.effectX, g.effectY = a.x, a.y
	switch data.kind {
	case effectLaser:
		index := g.nearestEnemy(a.x, a.y)
		if index >= 0 {
			g.enemies[index].hp -= data.power
			g.effectX, g.effectY = g.enemies[index].x, g.enemies[index].y
			g.message = "LASER event: nearest guard took 2 damage."
		}
	case effectHeal:
		before := g.hp
		g.hp = min(maxHP, g.hp+data.power)
		g.message = fmt.Sprintf("HEAL event: HP %d -> %d.", before, g.hp)
	case effectPulse:
		hits := 0
		for i := range g.enemies {
			if g.enemies[i].hp > 0 && distance(a.x, a.y, g.enemies[i].x, g.enemies[i].y) <= 235 {
				g.enemies[i].hp -= data.power
				hits++
			}
		}
		g.message = fmt.Sprintf("PULSE event: %d nearby guard(s) took damage.", hits)
	}
}

func (g *game) nearestEnemy(x, y float64) int {
	best, bestDistance := -1, math.Inf(1)
	for i, e := range g.enemies {
		if e.hp <= 0 {
			continue
		}
		d := distance(x, y, e.x, e.y)
		if d < bestDistance {
			best, bestDistance = i, d
		}
	}
	return best
}

func (g *game) bounceWalls() {
	if g.x < playerRadius+12 {
		g.x, g.vx = playerRadius+12, math.Abs(g.vx)*0.82
	}
	if g.x > width-playerRadius-12 {
		g.x, g.vx = width-playerRadius-12, -math.Abs(g.vx)*0.82
	}
	if g.y < 98+playerRadius {
		g.y, g.vy = 98+playerRadius, math.Abs(g.vy)*0.82
	}
	if g.y > 570-playerRadius {
		g.y, g.vy = 570-playerRadius, -math.Abs(g.vy)*0.82
	}
}

func (g *game) pushFrom(x, y float64) {
	dx, dy := g.x-x, g.y-y
	d := math.Hypot(dx, dy)
	if d == 0 {
		dx, dy, d = 1, 0, 1
	}
	nx, ny := dx/d, dy/d
	g.x, g.y = x+nx*(playerRadius+enemyRadius), y+ny*(playerRadius+enemyRadius)
	dot := g.vx*nx + g.vy*ny
	if dot < 0 {
		g.vx -= 1.7 * dot * nx
		g.vy -= 1.7 * dot * ny
	}
}

func distance(ax, ay, bx, by float64) float64 { return math.Hypot(ax-bx, ay-by) }

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 20, 38, 255})
	ebitenutil.DebugPrintAt(screen, "ALLY EFFECT RELAY", 181, 20)
	seconds := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME %02d   HP %d/%d   EVENTS %d", seconds, g.hp, maxHP, g.events), 142, 47)
	ebitenutil.DebugPrintAt(screen, g.message, 60, 73)
	vector.StrokeRect(screen, 10, 98, 460, 472, 4, color.RGBA{52, 84, 116, 255}, false)

	for _, a := range g.allies {
		data := abilities[a.ability]
		vector.DrawFilledCircle(screen, float32(a.x), float32(a.y), allyRadius, data.color, false)
		vector.StrokeCircle(screen, float32(a.x), float32(a.y), allyRadius, 3, color.White, false)
		ebitenutil.DebugPrintAt(screen, data.name, int(a.x)-len(data.name)*3, int(a.y)-5)
		if a.cooldown > 0 {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", (a.cooldown+14)/15), int(a.x)-3, int(a.y)+34)
		}
	}
	for _, e := range g.enemies {
		if e.hp <= 0 {
			continue
		}
		vector.DrawFilledCircle(screen, float32(e.x), float32(e.y), enemyRadius, color.RGBA{222, 82, 84, 255}, false)
		vector.StrokeCircle(screen, float32(e.x), float32(e.y), enemyRadius, 3, color.White, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP%d", max(0, e.hp)), int(e.x)-10, int(e.y)-5)
	}
	g.drawEffect(screen)
	playerColor := color.RGBA{240, 126, 83, 255}
	if g.damageCooldown > 0 {
		playerColor = color.RGBA{255, 207, 93, 255}
	}
	vector.DrawFilledCircle(screen, float32(g.x), float32(g.y), playerRadius, playerColor, false)
	vector.StrokeCircle(screen, float32(g.x), float32(g.y), playerRadius, 3, color.White, false)
	vector.StrokeLine(screen, float32(g.x), float32(g.y), float32(g.x+g.vx*6), float32(g.y+g.vy*6), 3, color.RGBA{255, 220, 112, 255}, false)
	ebitenutil.DebugPrintAt(screen, "HOLD/TAP A DIRECTION OR USE ARROWS / WASD", 72, 620)
	ebitenutil.DebugPrintAt(screen, "LEAVE AN ALLY, THEN TOUCH AGAIN TO RE-TRIGGER", 66, 649)
	if g.clear || g.over {
		title := "ALLY RELAY CLEAR!"
		if g.over {
			title = "RELAY FAILED!"
		}
		vector.DrawFilledRect(screen, 40, 270, 400, 160, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 166, 315)
		ebitenutil.DebugPrintAt(screen, g.message, 82, 350)
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE TO RETRY", 148, 390)
	}
}

func (g *game) drawEffect(screen *ebiten.Image) {
	if g.effectTimer <= 0 {
		return
	}
	switch g.effectKind {
	case effectLaser:
		// Laser starts at the LASER ally and ends at the stored target.
		a := g.allies[0]
		vector.StrokeLine(screen, float32(a.x), float32(a.y), float32(g.effectX), float32(g.effectY), 7, abilities[0].color, false)
	case effectHeal:
		r := float32(30 + (28-g.effectTimer)*2)
		vector.StrokeCircle(screen, float32(g.allies[1].x), float32(g.allies[1].y), r, 5, abilities[1].color, false)
	case effectPulse:
		r := float32(35 + (28-g.effectTimer)*8)
		vector.StrokeCircle(screen, float32(g.allies[2].x), float32(g.allies[2].y), r, 5, abilities[2].color, false)
	}
}

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ally Effect Relay — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
