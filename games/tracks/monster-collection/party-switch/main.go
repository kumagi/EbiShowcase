package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width, height = 480, 720
	tide          = 0
	flame         = 1
	leaf          = 2
	maxTurns      = 12
	timeLimit     = 60 * 60
)

var typeNames = []string{"TIDE", "FLAME", "LEAF"}
var typeColors = []color.RGBA{{70, 164, 218, 255}, {230, 100, 69, 255}, {91, 183, 103, 255}}

// Rows are attack types, columns are defender types.
var matchup = [3][3]float64{
	{1, 2, 0.5},
	{0.5, 1, 2},
	{2, 0.5, 1},
}

type partner struct {
	name string
	kind int
	hp   int
}

type game struct {
	party         [3]partner
	active        int
	enemyKinds    [3]int
	wave, enemyHP int
	turn, frames  int
	mustSwitch    bool
	clear, over   bool
	message       string
}

func newGame() *game {
	return &game{
		party:      [3]partner{{"TIDEBUD", tide, 34}, {"EMBERIMP", flame, 34}, {"MOSSHELL", leaf, 34}},
		enemyKinds: [3]int{leaf, tide, flame},
		enemyHP:    36,
		message:    "Enemy is LEAF. Switch to a strong front partner, then attack.",
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
		g.message = "Time up. Read the enemy type before choosing the front."
		return nil
	}
	if (inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeySpace)) && !g.mustSwitch {
		g.attack()
	}
	for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if inpututil.IsKeyJustPressed(key) {
			g.switchTo(i)
		}
	}
	if x, y, ok := pressPosition(); ok {
		if y >= 410 && y < 535 {
			g.switchTo(min(2, x/160))
		} else if y >= 565 {
			if !g.mustSwitch {
				g.attack()
			}
		}
	}
	return nil
}

func (g *game) attack() {
	front := &g.party[g.active]
	enemyKind := g.enemyKinds[g.wave]
	multiplier := matchup[front.kind][enemyKind]
	damage := int(12 * multiplier)
	g.enemyHP = max(0, g.enemyHP-damage)
	g.turn++
	if g.enemyHP == 0 {
		g.wave++
		if g.wave == len(g.enemyKinds) {
			g.clear = true
			g.message = "All three wardens defeated by a well-switched party!"
			return
		}
		g.enemyHP = 36
		g.message = fmt.Sprintf("Wave cleared with x%.1f! Next enemy is %s.", multiplier, typeNames[g.enemyKinds[g.wave]])
		g.checkTurnLimit()
		return
	}
	g.message = fmt.Sprintf("%s attacked x%.1f for %d. Enemy counterattacks.", front.name, multiplier, damage)
	g.enemyAttack()
	g.checkTurnLimit()
}

func (g *game) switchTo(index int) {
	if index < 0 || index >= len(g.party) {
		return
	}
	if g.party[index].hp <= 0 {
		g.message = g.party[index].name + " has fainted and cannot take the front."
		return
	}
	if index == g.active && !g.mustSwitch {
		g.message = g.party[index].name + " is already in front."
		return
	}
	oldName := g.party[g.active].name
	g.active = index
	if g.mustSwitch {
		g.mustSwitch = false
		g.message = fmt.Sprintf("Forced switch: %s takes the front. No extra turn spent.", g.party[index].name)
		return
	}
	// A planned switch is a battle action, so the rival gets this turn.
	g.turn++
	g.message = fmt.Sprintf("%s -> %s. Switching spent the turn!", oldName, g.party[index].name)
	g.enemyAttack()
	g.checkTurnLimit()
}

func (g *game) enemyAttack() {
	front := &g.party[g.active]
	enemyKind := g.enemyKinds[g.wave]
	multiplier := matchup[enemyKind][front.kind]
	damage := int(8 * multiplier)
	front.hp = max(0, front.hp-damage)
	g.message += fmt.Sprintf(" Counter %s x%.1f dealt %d.", typeNames[enemyKind], multiplier, damage)
	if front.hp == 0 {
		if g.aliveCount() == 0 {
			g.over = true
			g.message = "Every party member fainted. Use favorable matchups."
			return
		}
		g.mustSwitch = true
		g.message += " Front fainted: choose a living reserve!"
	}
}

func (g *game) aliveCount() int {
	count := 0
	for _, member := range g.party {
		if member.hp > 0 {
			count++
		}
	}
	return count
}

func (g *game) checkTurnLimit() {
	if !g.clear && !g.over && g.turn >= maxTurns {
		g.over = true
		g.message = "Twelve turns passed. Switch less and attack with x2 matchups."
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{11, 22, 41, 255})
	ebitenutil.DebugPrintAt(screen, "FRONT & RESERVE BATTLE", 164, 18)
	shownWave := min(g.wave, len(g.enemyKinds)-1)
	enemyKind := g.enemyKinds[shownWave]
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("WAVE %d/3   TURN %d/%d   TIME %02d", shownWave+1, g.turn, maxTurns, max(0, 60-g.frames/60)), 111, 45)
	front := g.party[g.active]

	vector.DrawFilledCircle(screen, 122, 190, 60, typeColors[front.kind], false)
	vector.StrokeCircle(screen, 122, 190, 60, 4, color.White, false)
	ebitenutil.DebugPrintAt(screen, "FRONT", 104, 184)
	ebitenutil.DebugPrintAt(screen, front.name, 90, 263)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s HP %d/34", typeNames[front.kind], front.hp), 80, 287)

	vector.DrawFilledCircle(screen, 358, 190, 60, typeColors[enemyKind], false)
	vector.StrokeCircle(screen, 358, 190, 60, 4, color.White, false)
	ebitenutil.DebugPrintAt(screen, "WARDEN", 336, 184)
	ebitenutil.DebugPrintAt(screen, typeNames[enemyKind], 340, 263)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %d/36", g.enemyHP), 333, 287)

	strong := (enemyKind + 2) % 3
	// The table cycle is TIDE > FLAME > LEAF > TIDE; find the x2 row directly.
	for kind := 0; kind < 3; kind++ {
		if matchup[kind][enemyKind] == 2 {
			strong = kind
		}
	}
	vector.DrawFilledRect(screen, 36, 315, 408, 68, color.RGBA{30, 51, 74, 255}, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FORECAST: %s attacks %s at x2.0", typeNames[strong], typeNames[enemyKind]), 103, 333)
	ebitenutil.DebugPrintAt(screen, "A planned switch spends one turn; forced switch does not.", 61, 360)

	for i, member := range g.party {
		x := i * 160
		fill := color.RGBA{37, 58, 79, 255}
		if i == g.active {
			fill = color.RGBA{63, 91, 119, 255}
		}
		if member.hp == 0 {
			fill = color.RGBA{49, 49, 57, 255}
		}
		vector.DrawFilledRect(screen, float32(x+7), 410, 146, 125, fill, false)
		vector.StrokeRect(screen, float32(x+7), 410, 146, 125, 3, typeColors[member.kind], false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %s", i+1, member.name), x+30, 429)
		ebitenutil.DebugPrintAt(screen, typeNames[member.kind], x+59, 459)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %02d/34", member.hp), x+48, 488)
		status := "RESERVE"
		if i == g.active {
			status = "FRONT"
		}
		if member.hp == 0 {
			status = "FAINTED"
		}
		ebitenutil.DebugPrintAt(screen, status, x+53, 513)
	}
	ebitenutil.DebugPrintAt(screen, g.message, 45, 548)
	buttonColor := color.RGBA{220, 113, 69, 255}
	label := "ATTACK [A / SPACE]"
	if g.mustSwitch {
		buttonColor = color.RGBA{89, 83, 112, 255}
		label = "FORCED SWITCH: CHOOSE A LIVING CARD"
	}
	vector.DrawFilledRect(screen, 30, 578, 420, 82, buttonColor, false)
	ebitenutil.DebugPrintAt(screen, label, 240-len(label)*3, 615)
	if g.clear || g.over {
		title := "PARTY STRATEGY CLEAR!"
		if g.over {
			title = "PARTY DEFEATED!"
		}
		vector.DrawFilledRect(screen, 40, 270, 400, 160, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 151, 315)
		ebitenutil.DebugPrintAt(screen, g.message, 65, 350)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 390)
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
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Front & Reserve Battle — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
