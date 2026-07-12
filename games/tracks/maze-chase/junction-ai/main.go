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
	screenW   = 480
	screenH   = 720
	cols      = 11
	rows      = 13
	tileSize  = 36
	mazeX     = 42
	mazeY     = 92
	goal      = 8
	timeLimit = 60 * 60
)

var maze = [...]string{
	"###########", "#.........#", "#.###.###.#", "#.........#",
	"###.#.#.###", "#.........#", "#.###.###.#", "#.........#",
	"###.#.#.###", "#.........#", "#.###.###.#", "#.........#", "###########",
}

type point struct{ x, y int }

var directions = [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

type candidate struct {
	dir      point
	distance int
}

type personality struct {
	name  string
	ahead int
	speed int
	color color.RGBA
}

type enemy struct {
	position    point
	direction   point
	personality personality
	tick        int
	target      point
	candidates  []candidate
	chosen      point
}

type game struct {
	player            point
	playerDirection   point
	playerCooldown    int
	enemies           []enemy
	chips             map[point]bool
	collected, frames int
	message           string
	won, lost         bool
}

func newGame() *game {
	g := &game{
		player:          point{5, 7},
		playerDirection: point{1, 0},
		chips:           map[point]bool{},
		message:         "Read both targets, then change corridors before they close in.",
		enemies: []enemy{
			{position: point{1, 1}, direction: point{1, 0}, personality: personality{"HUNTER", 0, 14, color.RGBA{244, 93, 93, 255}}},
			{position: point{9, 11}, direction: point{-1, 0}, personality: personality{"AMBUSHER", 3, 17, color.RGBA{96, 214, 232, 255}}},
		},
	}
	for _, p := range []point{{3, 1}, {7, 1}, {1, 3}, {9, 3}, {3, 5}, {7, 5}, {1, 9}, {9, 9}, {3, 11}, {7, 11}} {
		g.chips[p] = true
	}
	for i := range g.enemies {
		g.previewChoice(&g.enemies[i])
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
	g.frames++
	if g.frames >= timeLimit {
		g.lost = true
		g.message = "Time up. Study the target markers and try another route."
		return nil
	}
	if g.playerCooldown > 0 {
		g.playerCooldown--
	}
	if g.playerCooldown == 0 {
		if d, ok := directionInput(); ok {
			next := add(g.player, d)
			if passable(next) {
				g.player, g.playerDirection, g.playerCooldown = next, d, 7
				g.collectChip()
			}
		}
	}
	for i := range g.enemies {
		e := &g.enemies[i]
		e.tick++
		if e.tick >= e.personality.speed {
			e.tick = 0
			g.moveEnemy(e)
		}
		if e.position == g.player {
			g.lost = true
			g.message = e.personality.name + " predicted your route."
		}
	}
	return nil
}

func (g *game) moveEnemy(e *enemy) {
	g.previewChoice(e)
	e.position = add(e.position, e.chosen)
	e.direction = e.chosen
	g.previewChoice(e)
}

func (g *game) previewChoice(e *enemy) {
	e.target = g.player
	if e.personality.ahead > 0 {
		e.target = point{g.player.x + g.playerDirection.x*e.personality.ahead, g.player.y + g.playerDirection.y*e.personality.ahead}
	}
	legal := make([]point, 0, 4)
	for _, d := range directions {
		if passable(add(e.position, d)) {
			legal = append(legal, d)
		}
	}
	reverse := point{-e.direction.x, -e.direction.y}
	if len(legal) > 1 {
		forward := legal[:0]
		for _, d := range legal {
			if d != reverse {
				forward = append(forward, d)
			}
		}
		if len(forward) > 0 {
			legal = forward
		}
	}
	e.candidates = e.candidates[:0]
	bestDistance := 1 << 30
	for _, d := range legal {
		distance := manhattan(add(e.position, d), e.target)
		e.candidates = append(e.candidates, candidate{d, distance})
		if distance < bestDistance {
			bestDistance, e.chosen = distance, d
		}
	}
}

func (g *game) collectChip() {
	if !g.chips[g.player] {
		return
	}
	delete(g.chips, g.player)
	g.collected++
	g.message = fmt.Sprintf("Data chip collected: %d/%d", g.collected, goal)
	if g.collected >= goal {
		g.won = true
		g.message = "You read two personalities and escaped with eight chips!"
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 17, 31, 255})
	seconds := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, "JUNCTION AI", 197, 16)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CHIPS %d/%d   TIME %02d", g.collected, goal, seconds), 169, 40)
	ebitenutil.DebugPrintAt(screen, g.message, 42, 66)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(mazeX+x*tileSize), float32(mazeY+y*tileSize)
			c := color.RGBA{20, 38, 53, 255}
			if maze[y][x] == '#' {
				c = color.RGBA{47, 73, 117, 255}
			}
			vector.DrawFilledRect(screen, px+1, py+1, tileSize-2, tileSize-2, c, false)
			if g.chips[point{x, y}] {
				vector.DrawFilledRect(screen, px+13, py+13, 10, 10, color.RGBA{250, 205, 82, 255}, false)
			}
		}
	}
	for i := range g.enemies {
		g.drawEnemyThinking(screen, &g.enemies[i])
	}
	px, py := tileCenter(g.player)
	vector.DrawFilledCircle(screen, px, py, 12, color.RGBA{249, 224, 83, 255}, false)
	vector.DrawFilledCircle(screen, px+float32(g.playerDirection.x*6), py+float32(g.playerDirection.y*6), 3, color.RGBA{8, 18, 35, 255}, false)
	for i := range g.enemies {
		e := &g.enemies[i]
		ex, ey := tileCenter(e.position)
		vector.DrawFilledCircle(screen, ex, ey, 13, e.personality.color, false)
		vector.DrawFilledRect(screen, ex-10, ey+6, 20, 8, e.personality.color, false)
	}
	g.drawInspector(screen)
	drawControls(screen)
	if g.won || g.lost {
		g.drawResult(screen)
	}
}

func (g *game) drawEnemyThinking(screen *ebiten.Image, e *enemy) {
	ex, ey := tileCenter(e.position)
	for _, c := range e.candidates {
		nx, ny := tileCenter(add(e.position, c.dir))
		lineColor := color.RGBA{120, 139, 157, 150}
		if c.dir == e.chosen {
			lineColor = color.RGBA{249, 224, 83, 230}
		}
		vector.StrokeLine(screen, ex, ey, nx, ny, 3, lineColor, false)
	}
	tx := float32(mazeX + e.target.x*tileSize + tileSize/2)
	ty := float32(mazeY + e.target.y*tileSize + tileSize/2)
	if tx >= mazeX && tx < mazeX+cols*tileSize && ty >= mazeY && ty < mazeY+rows*tileSize {
		vector.StrokeCircle(screen, tx, ty, 9, 2, e.personality.color, false)
	}
}

func (g *game) drawInspector(screen *ebiten.Image) {
	for i := range g.enemies {
		e := &g.enemies[i]
		y := 570 + i*25
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s target(%d,%d) ahead=%d", e.personality.name, e.target.x, e.target.y, e.personality.ahead), 42, y)
		text := "candidates "
		for _, c := range e.candidates {
			text += fmt.Sprintf("%s:%d ", directionName(c.dir), c.distance)
		}
		text += "-> " + directionName(e.chosen)
		ebitenutil.DebugPrintAt(screen, text, 42, y+12)
	}
}

func (g *game) drawResult(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 48, 275, 384, 140, color.RGBA{7, 17, 31, 238}, false)
	title := "ROUTE CAUGHT"
	if g.won {
		title = "AI OUTSMARTED!"
	}
	ebitenutil.DebugPrintAt(screen, title, 185, 305)
	ebitenutil.DebugPrintAt(screen, g.message, 70, 335)
	ebitenutil.DebugPrintAt(screen, "PRESS ENTER / SPACE / TAP TO RETRY", 99, 375)
}

func directionInput() (point, bool) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return point{-1, 0}, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		return point{0, -1}, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		return point{0, 1}, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return point{1, 0}, true
	}
	x, y, pressed := pointer()
	if pressed {
		for i, box := range controlBoxes {
			if x >= box.x && x < box.x+box.w && y >= box.y && y < box.y+box.h {
				return directions[i], true
			}
		}
	}
	return point{}, false
}

type rect struct{ x, y, w, h int }

var controlBoxes = [...]rect{{42, 655, 88, 48}, {143, 625, 88, 48}, {244, 655, 88, 48}, {345, 655, 88, 48}}

func pointer() (int, int, bool) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	ids := ebiten.TouchIDs()
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}

func drawControls(screen *ebiten.Image) {
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT"}
	for i, box := range controlBoxes {
		vector.DrawFilledRect(screen, float32(box.x), float32(box.y), float32(box.w), float32(box.h), color.RGBA{32, 57, 76, 255}, false)
		ebitenutil.DebugPrintAt(screen, labels[i], box.x+25, box.y+19)
	}
}

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func passable(p point) bool {
	return p.x >= 0 && p.x < cols && p.y >= 0 && p.y < rows && maze[p.y][p.x] != '#'
}
func add(a, b point) point { return point{a.x + b.x, a.y + b.y} }
func manhattan(a, b point) int {
	dx, dy := a.x-b.x, a.y-b.y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
func tileCenter(p point) (float32, float32) {
	return float32(mazeX + p.x*tileSize + tileSize/2), float32(mazeY + p.y*tileSize + tileSize/2)
}
func directionName(d point) string {
	if d == (point{-1, 0}) {
		return "L"
	}
	if d == (point{1, 0}) {
		return "R"
	}
	if d == (point{0, -1}) {
		return "U"
	}
	if d == (point{0, 1}) {
		return "D"
	}
	return "-"
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Junction AI - Ebi Showcase")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
