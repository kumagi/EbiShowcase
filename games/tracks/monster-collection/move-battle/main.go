package main

import (
	"fmt"
	"image/color"
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
	playerHP, rivalHP int
	moves             [4]moveState
	rng               *rand.Rand
	turn              int
	queue             []action
	resolveTimer      int
	guard             bool
	orderText         string
	message           string
	won, lost         bool
}

func newGame() *game {
	g := &game{
		playerHP:  100,
		rivalHP:   125,
		rng:       rand.New(rand.NewSource(8202)),
		message:   "Choose a move. Power, accuracy, and uses all matter.",
		orderText: "TURN ORDER: waiting for your choice",
	}
	for i, data := range playerMoveData {
		g.moves[i] = moveState{data: data, uses: data.maxUses}
	}
	return g
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
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

	if g.rng.Intn(100) >= current.move.accuracy {
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
	ebitenutil.DebugPrintAt(screen, "MOVE DATA BATTLE", 183, 20)
	shownTurn := g.turn + 1
	if len(g.queue) > 0 {
		shownTurn = g.turn
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TURN %02d", shownTurn), 212, 46)
	ebitenutil.DebugPrintAt(screen, g.orderText, 80, 72)
	vector.DrawFilledRect(screen, 20, 94, 440, 350, color.RGBA{27, 53, 72, 255}, false)

	drawHP(screen, 34, 115, "YOUR / TIDEBUD", g.playerHP, 100, color.RGBA{72, 191, 213, 255})
	drawHP(screen, 266, 115, "RIVAL / GLOOMFIN", g.rivalHP, 125, color.RGBA{205, 79, 104, 255})
	vector.DrawFilledCircle(screen, 125, 275, 55, color.RGBA{72, 191, 213, 255}, false)
	vector.DrawFilledCircle(screen, 355, 245, 61, color.RGBA{173, 83, 191, 255}, false)
	vector.StrokeCircle(screen, 125, 275, 55, 4, color.White, false)
	vector.StrokeCircle(screen, 355, 245, 61, 4, color.White, false)
	ebitenutil.DebugPrintAt(screen, "SPD 12", 105, 270)
	ebitenutil.DebugPrintAt(screen, "SPD 10", 335, 240)
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
	if g.won {
		overlay(screen, "TACTICAL VICTORY!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
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
