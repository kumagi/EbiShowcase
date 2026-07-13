package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW = 480
	screenH = 720
)

type moveData struct {
	name                     string
	power, accuracy, maxUses int
	priority                 int
	effect                   int // 0 damage, 1 guard, 2 heal
}

type moveState struct {
	data moveData
	uses int
}

type action struct {
	actor int // 0 player, 1 rival
	move  moveData
}

// battleFX is visual state only. HP, accuracy, guard, and turn order remain in
// resolveNext; this layer reads the result and makes it visible over time.
type battleFX struct {
	active        bool
	actor, effect int
	timer         int
	hit           bool
}

type spark struct {
	x, y, vx, vy, life float64
	c                  color.RGBA
}

var playerMoveData = [...]moveData{
	{name: "JET PINCH", power: 18, accuracy: 100, maxUses: 6, priority: 1},
	{name: "REEF BURST", power: 42, accuracy: 70, maxUses: 3, priority: 0},
	{name: "SHELL GUARD", accuracy: 100, maxUses: 3, priority: 3, effect: 1},
	{name: "TIDE MEND", accuracy: 100, maxUses: 2, priority: 2, effect: 2},
}

var rivalMoves = [...]moveData{
	{name: "SHADOW NIP", power: 17, accuracy: 95, priority: 1},
	{name: "ABYSS ROAR", power: 34, accuracy: 75, priority: 0},
}

type game struct {
	playerHP, rivalHP           int
	shownPlayerHP, shownRivalHP int
	moves                       [4]moveState
	rng                         *rand.Rand
	turn                        int
	queue                       []action
	resolveTimer                int
	guard                       bool
	orderText                   string
	message                     string
	won, lost                   bool
	tick, shake                 int
	flashActor                  int
	flashTimer                  int
	fx                          battleFX
	sparks                      []spark
}

func newGame() *game {
	g := &game{
		playerHP:      100,
		rivalHP:       125,
		shownPlayerHP: 100,
		shownRivalHP:  125,
		rng:           rand.New(rand.NewSource(8202)),
		message:       "Choose a move. Power, accuracy, and uses all matter.",
		orderText:     "TURN ORDER: waiting for your choice",
		flashActor:    -1,
	}
	for i, data := range playerMoveData {
		g.moves[i] = moveState{data: data, uses: data.maxUses}
	}
	return g
}

func (g *game) Update() error {
	g.tick++
	g.updateFX()
	if g.won || g.lost {
		if g.fx.active {
			return nil
		}
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	if g.fx.active {
		return nil
	}
	if len(g.queue) > 0 {
		g.resolveTimer--
		if g.resolveTimer <= 0 {
			g.resolveNext()
		}
		return nil
	}

	choice := -1
	for i, key := range [...]ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if inpututil.IsKeyJustPressed(key) {
			choice = i
		}
	}
	if x, y, ok := pressPosition(); ok && y >= 470 {
		row := (y - 470) / 95
		col := x / 240
		choice = row*2 + col
	}
	if choice >= 0 && choice < len(g.moves) {
		g.startTurn(choice)
	}
	return nil
}

func (g *game) startTurn(choice int) {
	selected := &g.moves[choice]
	if selected.uses == 0 {
		g.message = selected.data.name + " has no uses left."
		return
	}
	selected.uses--
	g.turn++
	g.guard = false
	rivalMove := rivalMoves[0]
	if g.turn%3 == 0 {
		rivalMove = rivalMoves[1]
	}
	g.queue = []action{{actor: 0, move: selected.data}, {actor: 1, move: rivalMove}}
	// Higher priority acts first. A tie uses creature speed: player 12, rival 10.
	sort.SliceStable(g.queue, func(i, j int) bool {
		if g.queue[i].move.priority != g.queue[j].move.priority {
			return g.queue[i].move.priority > g.queue[j].move.priority
		}
		return speedFor(g.queue[i].actor) > speedFor(g.queue[j].actor)
	})
	g.orderText = fmt.Sprintf("ORDER: %s -> %s", g.queue[0].move.name, g.queue[1].move.name)
	g.message = fmt.Sprintf("Turn %d queued. Resolve the first action...", g.turn)
	g.resolveTimer = 35
}

func speedFor(actor int) int {
	if actor == 0 {
		return 12
	}
	return 10
}

func (g *game) resolveNext() {
	current := g.queue[0]
	g.queue = g.queue[1:]
	if current.actor == 1 && g.rivalHP <= 0 {
		g.queue = nil
		return
	}
	if current.actor == 0 && g.playerHP <= 0 {
		g.queue = nil
		return
	}

	hit := g.rng.Intn(100) < current.move.accuracy
	if !hit {
		g.message = current.move.name + " missed! Accuracy roll failed."
	} else {
		switch current.move.effect {
		case 1:
			g.guard = true
			g.message = "SHELL GUARD: the next rival hit is halved."
		case 2:
			before := g.playerHP
			g.playerHP = min(100, g.playerHP+25)
			g.message = fmt.Sprintf("TIDE MEND restored %d HP.", g.playerHP-before)
		default:
			damage := current.move.power
			if current.actor == 0 {
				g.rivalHP = max(0, g.rivalHP-damage)
				g.message = fmt.Sprintf("%s hit for %d. Rival HP %d.", current.move.name, damage, g.rivalHP)
			} else {
				if g.guard {
					damage = (damage + 1) / 2
					g.guard = false
				}
				g.playerHP = max(0, g.playerHP-damage)
				g.message = fmt.Sprintf("Rival %s dealt %d. Your HP %d.", current.move.name, damage, g.playerHP)
			}
		}
	}
	g.startFX(current.actor, current.move.effect, hit)

	if g.rivalHP <= 0 {
		g.won = true
		g.queue = nil
		g.message = "The rival creature is defeated!"
		return
	}
	if g.playerHP <= 0 {
		g.lost = true
		g.queue = nil
		g.message = "Your creature cannot battle."
		return
	}
	if len(g.queue) == 0 {
		g.guard = false
		g.orderText = "TURN ORDER: choose the next move"
		if g.noUsesLeft() {
			g.lost = true
			g.message = "Every move is out of uses. Plan the limited moves!"
		}
		return
	}
	g.resolveTimer = 42
}

func (g *game) startFX(actor, effect int, hit bool) {
	g.fx = battleFX{active: true, actor: actor, effect: effect, timer: 32, hit: hit}
}

func (g *game) updateFX() {
	if g.shake > 0 {
		g.shake--
	}
	if g.flashTimer > 0 {
		g.flashTimer--
		if g.flashTimer == 0 {
			g.flashActor = -1
		}
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .06
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if !g.fx.active {
		return
	}
	g.fx.timer--
	if g.fx.timer == 15 && g.fx.hit {
		target := 1 - g.fx.actor
		if g.fx.effect == 1 || g.fx.effect == 2 {
			target = g.fx.actor
		}
		g.flashActor, g.flashTimer = target, 7
		// Rules already produced the result; reveal those values on the visual
		// contact frame so the HP bars and impact agree without moving rules here.
		g.shownPlayerHP, g.shownRivalHP = g.playerHP, g.rivalHP
		g.shake = 4
		if g.fx.effect == 0 {
			g.shake = 8
		}
		x, y := creaturePosition(target)
		c := color.RGBA{255, 196, 66, 255}
		if g.fx.effect == 2 {
			c = color.RGBA{92, 235, 183, 255}
		}
		for i := 0; i < 16; i++ {
			a := float64(i) * math.Pi / 8
			g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * (2 + float64(i%3)), math.Sin(a) * (2 + float64(i%3)), 26, c})
		}
	}
	if g.fx.timer <= 0 {
		g.fx.active = false
	}
}

func creaturePosition(actor int) (float64, float64) {
	if actor == 0 {
		return 125, 275
	}
	return 355, 245
}

func (g *game) noUsesLeft() bool {
	for _, move := range g.moves {
		if move.uses > 0 {
			return false
		}
	}
	return true
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 21, 40, 255})
	shakeX := 0.0
	if g.shake > 0 {
		shakeX = math.Sin(float64(g.tick)*2.5) * float64(g.shake)
	}
	ebitenutil.DebugPrintAt(screen, "MOVE DATA BATTLE", 183, 20)
	shownTurn := g.turn + 1
	if len(g.queue) > 0 {
		shownTurn = g.turn
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TURN %02d", shownTurn), 212, 46)
	ebitenutil.DebugPrintAt(screen, g.orderText, 80, 72)
	vector.DrawFilledRect(screen, 20, 94, 440, 350, color.RGBA{27, 53, 72, 255}, false)

	drawHP(screen, 34, 115, "YOUR / TIDEBUD", g.shownPlayerHP, 100, color.RGBA{72, 191, 213, 255})
	drawHP(screen, 266, 115, "RIVAL / GLOOMFIN", g.shownRivalHP, 125, color.RGBA{205, 79, 104, 255})
	playerX, rivalX := 125.0, 355.0
	if g.fx.active && g.fx.effect == 0 {
		progress := 1 - math.Abs(float64(g.fx.timer-15))/17
		progress = max(0, progress)
		if g.fx.actor == 0 {
			playerX += progress * 55
		} else {
			rivalX -= progress * 55
		}
	}
	playerColor := color.RGBA{72, 191, 213, 255}
	rivalColor := color.RGBA{173, 83, 191, 255}
	if g.flashActor == 0 {
		playerColor = color.RGBA{255, 245, 180, 255}
	}
	if g.flashActor == 1 {
		rivalColor = color.RGBA{255, 235, 170, 255}
	}
	pulse := float32(math.Sin(float64(g.tick)*.1) * 2)
	vector.DrawFilledCircle(screen, float32(playerX+shakeX), 275, 55+pulse, playerColor, false)
	vector.DrawFilledCircle(screen, float32(rivalX+shakeX), 245, 61-pulse, rivalColor, false)
	vector.StrokeCircle(screen, float32(playerX+shakeX), 275, 55+pulse, 4, color.White, false)
	vector.StrokeCircle(screen, float32(rivalX+shakeX), 245, 61-pulse, 4, color.White, false)
	ebitenutil.DebugPrintAt(screen, "SPD 12", int(playerX+shakeX)-20, 270)
	ebitenutil.DebugPrintAt(screen, "SPD 10", int(rivalX+shakeX)-20, 240)
	if g.fx.active && g.fx.effect == 1 {
		x, y := creaturePosition(g.fx.actor)
		r := float32(62 + math.Sin(float64(g.tick)*.35)*7)
		vector.StrokeCircle(screen, float32(x+shakeX), float32(y), r, 6, color.RGBA{105, 220, 255, 255}, true)
	}
	if g.fx.active && !g.fx.hit && g.fx.timer < 20 {
		target := 1 - g.fx.actor
		x, y := creaturePosition(target)
		vector.StrokeCircle(screen, float32(x+shakeX), float32(y), 76, 3, color.RGBA{180, 195, 215, 180}, true)
		ebitenutil.DebugPrintAt(screen, "MISS", int(x)-14, int(y)-90)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(screen, float32(p.x+shakeX), float32(p.y), float32(2+p.life/9), p.c, true)
	}
	ebitenutil.DebugPrintAt(screen, g.message, 48, 410)

	for i, move := range g.moves {
		row, col := i/2, i%2
		x, y := col*240+6, row*95+470
		c := color.RGBA{51, 84, 122, 255}
		if move.uses == 0 {
			c = color.RGBA{55, 59, 70, 255}
		}
		vector.DrawFilledRect(screen, float32(x), float32(y), 228, 86, c, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%d] %s", i+1, move.data.name), x+12, y+14)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PWR %02d  ACC %03d%%", move.data.power, move.data.accuracy), x+12, y+38)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("USES %d/%d  PRIORITY %d", move.uses, move.data.maxUses, move.data.priority), x+12, y+61)
	}
	ebitenutil.DebugPrintAt(screen, "Keyboard 1-4 / click / tap a move", 103, 681)
	if g.won && !g.fx.active {
		overlay(screen, "TACTICAL VICTORY!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost && !g.fx.active {
		overlay(screen, "BATTLE LOST\n\nTAP / ENTER TO RETRY")
	}
}

func drawHP(screen *ebiten.Image, x, y int, name string, hp, maxHP int, c color.RGBA) {
	ebitenutil.DebugPrintAt(screen, name, x, y)
	vector.DrawFilledRect(screen, float32(x), float32(y+25), 180, 16, color.RGBA{44, 55, 71, 255}, false)
	vector.DrawFilledRect(screen, float32(x), float32(y+25), float32(180*hp/maxHP), 16, c, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %d/%d", hp, maxHP), x+55, y+48)
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
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 115, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Move Data Battle — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
