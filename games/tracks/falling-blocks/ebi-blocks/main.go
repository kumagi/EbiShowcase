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
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const (
	screenW = 480
	screenH = 720
	cols    = 10
	rows    = 18
	cell    = 26
	boardX  = 110
	boardY  = 92
	noPiece = -1
)

type point struct{ x, y int }

type piece struct {
	kind, rotation int
	pos            point
}

type stageRule struct {
	name       string
	goal       int
	baseFall   int
	background color.RGBA
}

type spark struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}

var stages = [...]stageRule{
	{"SUNSET DOCK", 2, 42, color.RGBA{11, 20, 38, 255}},
	{"CORAL CAVE", 3, 34, color.RGBA{25, 15, 48, 255}},
	{"AURORA DECK", 4, 27, color.RGBA{7, 38, 50, 255}},
}

var shapeNames = [...]string{"I", "O", "T", "L", "J", "S", "Z"}
var shapeColors = [...]color.RGBA{
	{65, 190, 207, 255}, {239, 190, 62, 255}, {174, 97, 205, 255},
	{232, 142, 56, 255}, {71, 118, 209, 255}, {90, 181, 101, 255}, {224, 82, 86, 255},
}
var bases = [7][4]point{
	{{-1, 0}, {0, 0}, {1, 0}, {2, 0}},
	{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	{{-1, 0}, {0, 0}, {1, 0}, {0, 1}},
	{{-1, 0}, {0, 0}, {1, 0}, {1, 1}},
	{{-1, 0}, {0, 0}, {1, 0}, {-1, 1}},
	{{0, 0}, {1, 0}, {-1, 1}, {0, 1}},
	{{-1, 0}, {0, 0}, {0, 1}, {1, 1}},
}
var kickTests = [...]point{{0, 0}, {-1, 0}, {1, 0}, {-2, 0}, {2, 0}, {0, -1}}

type game struct {
	board                     [rows][cols]int
	active                    piece
	rng                       *rand.Rand
	bag, queue                []int
	hold                      int
	canHold                   bool
	tick, score, lines, combo int
	level                     int
	stage, stageLines, pieces int
	best                      int
	clearFlash, shake         int
	frame                     int
	sparks                    []spark
	scene                     *ebiten.Image
	message                   string
	won, lost                 bool
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(6606)), hold: noPiece, level: 1, combo: -1, scene: ebiten.NewImage(screenW, screenH)}
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = noPiece
		}
	}
	g.fillQueue()
	g.spawn(g.takeNext())
	g.message = "Clear 2 lines to cross the sunset dock!"
	return g
}

func (g *game) refillBag() {
	g.bag = g.rng.Perm(len(bases))
}

func (g *game) fillQueue() {
	for len(g.queue) < 5 {
		if len(g.bag) == 0 {
			g.refillBag()
		}
		g.queue = append(g.queue, g.bag[0])
		g.bag = g.bag[1:]
	}
}

func (g *game) takeNext() int {
	g.fillQueue()
	kind := g.queue[0]
	g.queue = g.queue[1:]
	g.fillQueue()
	return kind
}

func (g *game) spawn(kind int) {
	g.active = piece{kind: kind, pos: point{cols / 2, 0}}
	g.tick = 0
	g.canHold = true
	if g.collides(g.active) {
		g.lost = true
		g.message = "No room for the next piece."
	}
}

func (g *game) Update() error {
	g.frame++
	if g.clearFlash > 0 {
		g.clearFlash--
	}
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .08
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.won || g.lost {
		if retryPressed() {
			best := g.best
			*g = *newGame()
			g.best = best
		}
		return nil
	}

	left := inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	rotate := inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	soft := inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)
	hard := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	hold := inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyShift)
	if x, y, ok := pointerPress(); ok && y >= 620 {
		switch x / 80 {
		case 0:
			left = true
		case 1:
			right = true
		case 2:
			rotate = true
		case 3:
			soft = true
		case 4:
			hard = true
		case 5:
			hold = true
		}
	}

	if left {
		g.move(-1, 0)
	}
	if right {
		g.move(1, 0)
	}
	if rotate {
		g.rotate()
	}
	if hold {
		g.swapHold()
	}
	if hard {
		distance := 0
		for g.move(0, 1) {
			distance++
		}
		g.score += distance * 2
		g.lock()
		return nil
	}
	if soft {
		if g.move(0, 1) {
			g.score++
		} else {
			g.lock()
		}
		g.tick = 0
		return nil
	}

	g.tick++
	fallEvery := max(8, stages[g.stage].baseFall-(g.level-1)*3)
	if g.tick >= fallEvery {
		g.tick = 0
		if !g.move(0, 1) {
			g.lock()
		}
	}
	return nil
}

func (g *game) blocks(p piece) [4]point {
	result := bases[p.kind]
	if p.kind != 1 { // the square looks identical in every rotation
		for turn := 0; turn < p.rotation; turn++ {
			for i, b := range result {
				result[i] = point{-b.y, b.x}
			}
		}
	}
	for i := range result {
		result[i].x += p.pos.x
		result[i].y += p.pos.y
	}
	return result
}

func (g *game) collides(p piece) bool {
	for _, b := range g.blocks(p) {
		if b.x < 0 || b.x >= cols || b.y < 0 || b.y >= rows || g.board[b.y][b.x] != noPiece {
			return true
		}
	}
	return false
}

func (g *game) move(dx, dy int) bool {
	next := g.active
	next.pos.x += dx
	next.pos.y += dy
	if g.collides(next) {
		return false
	}
	g.active = next
	return true
}

func (g *game) rotate() {
	next := g.active
	next.rotation = (next.rotation + 1) % 4
	for _, kick := range kickTests {
		candidate := next
		candidate.pos.x += kick.x
		candidate.pos.y += kick.y
		if !g.collides(candidate) {
			g.active = candidate
			g.message = fmt.Sprintf("Rotated with kick (%+d,%+d).", kick.x, kick.y)
			return
		}
	}
	g.message = "Rotation blocked at every kick position."
}

func (g *game) swapHold() {
	if !g.canHold {
		g.message = "HOLD can be used once per falling piece."
		return
	}
	old := g.active.kind
	if g.hold == noPiece {
		g.hold = old
		g.active = piece{kind: g.takeNext(), pos: point{cols / 2, 0}}
	} else {
		old, g.hold = g.hold, old
		g.active = piece{kind: old, pos: point{cols / 2, 0}}
	}
	g.canHold = false
	g.message = "Piece moved through HOLD."
	if g.collides(g.active) {
		g.lost = true
	}
}

func (g *game) ghost() piece {
	ghost := g.active
	for {
		next := ghost
		next.pos.y++
		if g.collides(next) {
			return ghost
		}
		ghost = next
	}
}

func (g *game) lock() {
	for _, b := range g.blocks(g.active) {
		g.board[b.y][b.x] = g.active.kind
	}
	cleared := g.clearLines()
	g.pieces++
	if cleared > 0 {
		g.combo++
		points := [...]int{0, 100, 300, 500, 800}
		g.score += points[cleared]*g.level + g.combo*50
		g.lines += cleared
		g.stageLines += cleared
		g.level = 1 + g.lines/5
		g.clearFlash = 20
		g.shake = 10 + cleared*3
		g.burst(cleared)
		g.message = fmt.Sprintf("%d line(s)! Combo %d, level %d.", cleared, g.combo+1, g.level)
	} else {
		g.combo = -1
		g.message = "Piece locked. Build a full row."
	}
	if g.score > g.best {
		g.best = g.score
	}
	if g.stageLines >= stages[g.stage].goal {
		if g.stage == len(stages)-1 {
			g.won = true
			if g.score > g.best {
				g.best = g.score
			}
			return
		}
		g.stage++
		g.stageLines = 0
		for y := range g.board {
			for x := range g.board[y] {
				g.board[y][x] = noPiece
			}
		}
		g.message = fmt.Sprintf("STAGE %d! %s: clear %d lines.", g.stage+1, stages[g.stage].name, stages[g.stage].goal)
		g.clearFlash = 45
	}
	g.spawn(g.takeNext())
}

func (g *game) burst(lines int) {
	for i := 0; i < 18+lines*12; i++ {
		a := g.rng.Float64() * math.Pi * 2
		sp := 1.2 + g.rng.Float64()*3.4
		g.sparks = append(g.sparks, spark{x: boardX + cols*cell/2, y: boardY + rows*cell - 20,
			vx: math.Cos(a) * sp, vy: math.Sin(a)*sp - 1, life: 24 + g.rng.Intn(22), c: shapeColors[g.rng.Intn(len(shapeColors))]})
	}
}

func (g *game) clearLines() int {
	write, cleared := rows-1, 0
	for read := rows - 1; read >= 0; read-- {
		full := true
		for x := 0; x < cols; x++ {
			if g.board[read][x] == noPiece {
				full = false
				break
			}
		}
		if full {
			cleared++
			continue
		}
		g.board[write] = g.board[read]
		write--
	}
	for write >= 0 {
		for x := 0; x < cols; x++ {
			g.board[write][x] = noPiece
		}
		write--
	}
	return cleared
}

func (g *game) Draw(screen *ebiten.Image) {
	g.scene.Clear()
	g.drawScene(g.scene)
	op := &ebiten.DrawImageOptions{}
	if g.shake > 0 {
		op.GeoM.Translate(float64((g.frame%3-1)*3), float64(((g.frame/2)%3-1)*2))
	}
	screen.DrawImage(g.scene, op)
}

func (g *game) drawScene(screen *ebiten.Image) {
	screen.Fill(stages[g.stage].background)
	ebitenutil.DebugPrintAt(screen, "EBI BLOCKS", 202, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("STAGE %d/3 %-11s  LINES %d/%d", g.stage+1, stages[g.stage].name, g.stageLines, stages[g.stage].goal), 98, 42)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %05d  BEST %05d  LEVEL %d", g.score, g.best, g.level), 120, 58)
	ebitenutil.DebugPrintAt(screen, g.message, 76, 76)

	vector.DrawFilledRect(screen, boardX-3, boardY-3, cols*cell+6, rows*cell+6, color.RGBA{35, 48, 70, 255}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{57, 72, 97, 255}, false)
			if g.board[y][x] != noPiece {
				drawCell(screen, x, y, shapeColors[g.board[y][x]], 255)
			}
		}
	}
	if !g.won && !g.lost {
		for _, b := range g.blocks(g.ghost()) {
			drawCell(screen, b.x, b.y, shapeColors[g.active.kind], 55)
		}
		for _, b := range g.blocks(g.active) {
			drawCell(screen, b.x, b.y, shapeColors[g.active.kind], 255)
		}
		// A breathing highlight provides in-between frames even while the piece waits.
		pulse := uint8(25 + 20*(1+math.Sin(float64(g.frame)*.14)))
		for _, b := range g.blocks(g.active) {
			drawCell(screen, b.x, b.y, color.RGBA{255, 255, 255, 255}, pulse)
		}
	}
	for _, p := range g.sparks {
		a := uint8(min(255, p.life*9))
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 2.5, color.RGBA{p.c.R, p.c.G, p.c.B, a}, false)
	}
	if g.clearFlash > 0 {
		a := uint8(g.clearFlash * 3)
		if a > 80 {
			a = 80
		}
		vector.DrawFilledRect(screen, boardX, boardY, cols*cell, rows*cell, color.RGBA{255, 245, 180, a}, false)
	}

	drawSide(screen, 8, 110, "HOLD", g.hold)
	next := noPiece
	if len(g.queue) > 0 {
		next = g.queue[0]
	}
	drawSide(screen, 382, 110, "NEXT", next)
	ebitenutil.DebugPrintAt(screen, "7-BAG", 399, 205)
	for i := 0; i < min(4, len(g.queue)); i++ {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d  %s", i+1, shapeNames[g.queue[i]]), 404, 229+i*22)
	}

	labels := [...]string{"LEFT", "RIGHT", "TURN", "DOWN", "DROP", "HOLD"}
	for i, label := range labels {
		c := color.RGBA{55, 88, 130, 255}
		if i == 2 || i == 5 {
			c = color.RGBA{194, 116, 67, 255}
		}
		vector.DrawFilledRect(screen, float32(i*80+3), 620, 74, 70, c, false)
		ebitenutil.DebugPrintAt(screen, label, i*80+18, 650)
	}
	if g.won {
		overlay(screen, fmt.Sprintf("ALL 3 STAGES CLEAR!\nSCORE %d  PIECES %d\n\nTAP / ENTER TO PLAY AGAIN", g.score, g.pieces))
	}
	if g.lost {
		overlay(screen, "STACK REACHED THE TOP\n\nTAP / ENTER TO RETRY")
	}
}

func drawCell(screen *ebiten.Image, x, y int, c color.RGBA, alpha uint8) {
	px, py := float32(boardX+x*cell), float32(boardY+y*cell)
	a := float32(alpha) / 255
	r, g, b := float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	trackatlas.DrawTinted(screen, "block-cell", float64(px+cell/2), float64(py+cell/2), float64(cell-4), r, g, b, a)
}

func drawSide(screen *ebiten.Image, x, y int, title string, kind int) {
	vector.DrawFilledRect(screen, float32(x), float32(y), 90, 78, color.RGBA{31, 45, 68, 255}, false)
	ebitenutil.DebugPrintAt(screen, title, x+26, y+12)
	value := "—"
	if kind != noPiece {
		value = shapeNames[kind]
	}
	ebitenutil.DebugPrintAt(screen, value, x+40, y+45)
}

func pointerPress() (int, int, bool) {
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
	_, _, ok := pointerPress()
	return ok
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 160, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 160, 4, color.RGBA{241, 185, 65, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 116, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Blocks — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
