package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"

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

//go:embed assets/falling-blocks-cargo-tower.png
var cargoTowerPNG []byte

//go:embed assets/falling-blocks-tile-atlas.png
var tileAtlasPNG []byte
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
	message                   string
	won, lost                 bool
	audio                     *audio.Context
	gate                      audiolab.Gate
	pulse                     *shaderlab.Pulse
	cam                       cameralab.State
	badge                     *ebiten.Image
	background                *ebiten.Image
	tileArt                   [7]*ebiten.Image
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(6606)), hold: noPiece, level: 1, combo: -1}
	g.loadGeneratedArt()
	g.audio = audiolab.Context()
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: screenW / 2, Y: screenH / 2}, ViewW: screenW, ViewH: screenH}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{239, 190, 62, 255})
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = noPiece
		}
	}
	g.seedBoard(0)
	g.fillQueue()
	g.spawn(g.takeNext())
	g.message = "Clear 2 lines to cross the sunset dock!"
	return g
}

// seedBoard makes each stage begin as a small, readable recovery puzzle instead
// of an empty waiting room. The holes never form a completed row by themselves.
func (g *game) seedBoard(stage int) {
	patterns := [][]string{
		{"111111....", "..22..333."},
		{"44..555.66", "4...5...6."},
		{"777.11.222", ".7..1...2."},
	}
	for dy, row := range patterns[stage%len(patterns)] {
		y := rows - 1 - dy
		for x, ch := range row {
			if ch != '.' {
				g.board[y][x] = int(ch-'1') % len(bases)
			}
		}
	}
}

func mustDecodePNG(data []byte) *ebiten.Image {
	source, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(source)
}

func (g *game) loadGeneratedArt() {
	g.background = mustDecodePNG(cargoTowerPNG)
	atlas := mustDecodePNG(tileAtlasPNG)
	panelW := atlas.Bounds().Dx() / len(g.tileArt)
	const cropTop, cropBottom = 115, 320
	for i := range g.tileArt {
		panel := atlas.SubImage(image.Rect(i*panelW, cropTop, (i+1)*panelW, cropBottom))
		g.tileArt[i] = ebiten.NewImageFromImage(panel)
	}
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
		g.play(440 + float64(cleared)*120)
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
		g.play(180)
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
		g.seedBoard(g.stage)
		g.message = fmt.Sprintf("STAGE %d! %s: clear %d lines.", g.stage+1, stages[g.stage].name, stages[g.stage].goal)
		g.clearFlash = 45
	}
	g.spawn(g.takeNext())
}

func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Square, freq, .055)).Play()
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
	// The presentation buffer is local: Draw projects game state but never
	// mutates a cache stored in game.
	scene := ebiten.NewImage(screenW, screenH)
	g.drawScene(scene)
	op := &ebiten.DrawImageOptions{}
	if g.shake > 0 {
		op.GeoM.Translate(float64((g.frame%3-1)*3), float64(((g.frame/2)%3-1)*2))
	}
	screen.DrawImage(scene, op)
}

func (g *game) drawScene(screen *ebiten.Image) {
	screen.Fill(stages[g.stage].background)
	g.drawGeneratedBackground(screen)
	drawBlockBackdrop(screen, g.stage, g.frame)
	g.drawTitle(screen)
	g.drawEffectBadge(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("STAGE %d/3 %-11s  LINES %d/%d", g.stage+1, stages[g.stage].name, g.stageLines, stages[g.stage].goal), 98, 42)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %05d  BEST %05d  LEVEL %d", g.score, g.best, g.level), 120, 58)
	ebitenutil.DebugPrintAt(screen, g.message, 76, 76)

	vector.DrawFilledRect(screen, boardX-12, boardY-12, cols*cell+24, rows*cell+24, color.RGBA{3, 8, 20, 175}, false)
	vector.DrawFilledRect(screen, boardX-5, boardY-5, cols*cell+10, rows*cell+10, color.RGBA{35, 48, 70, 255}, false)
	vector.StrokeRect(screen, boardX-5, boardY-5, cols*cell+10, rows*cell+10, 3, color.RGBA{113, 219, 226, 170}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{57, 72, 97, 255}, false)
			if g.board[y][x] != noPiece {
				g.drawCell(screen, x, y, g.board[y][x], 255, false)
			}
		}
	}
	if !g.won && !g.lost {
		for _, b := range g.blocks(g.ghost()) {
			g.drawCell(screen, b.x, b.y, g.active.kind, 55, false)
		}
		for _, b := range g.blocks(g.active) {
			g.drawCell(screen, b.x, b.y, g.active.kind, 255, false)
		}
		// A breathing highlight provides in-between frames even while the piece waits.
		pulse := uint8(25 + 20*(1+math.Sin(float64(g.frame)*.14)))
		for _, b := range g.blocks(g.active) {
			g.drawCell(screen, b.x, b.y, g.active.kind, pulse, true)
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

	g.drawSide(screen, 8, 110, "HOLD", g.hold)
	next := noPiece
	if len(g.queue) > 0 {
		next = g.queue[0]
	}
	g.drawSide(screen, 382, 110, "NEXT", next)
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

func (g *game) drawGeneratedBackground(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	tints := [...]struct{ r, gr, b float32 }{{.88, .94, 1}, {.92, .76, 1}, {.72, 1, .92}}
	t := tints[g.stage]
	op.ColorScale.Scale(t.r, t.gr, t.b, .94)
	screen.DrawImage(g.background, op)
	vector.DrawFilledRect(screen, 0, 0, screenW, 86, color.RGBA{2, 8, 22, 120}, false)
	vector.DrawFilledRect(screen, 0, 604, screenW, 116, color.RGBA{2, 8, 22, 145}, false)
}

func drawBlockBackdrop(screen *ebiten.Image, stage, frame int) {
	// Only subtle live scan lines sit on top of the generated cargo tower.
	for y := 90; y < 590; y += 28 {
		vector.DrawFilledRect(screen, 0, float32(y), screenW, 1, color.RGBA{103, 202, 220, 18}, false)
	}
	beamX := float32((frame*2+stage*43)%620 - 70)
	vector.StrokeLine(screen, beamX, 80, beamX-120, 590, 10, color.RGBA{78, 196, 221, 14}, true)
}

func (g *game) drawTitle(screen *ebiten.Image) {
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(202, 6)
		text.Draw(screen, "EBI BLOCKS", face, op)
		return
	}
	ebitenutil.DebugPrintAt(screen, "EBI BLOCKS", 202, 18)
}

func (g *game) drawEffectBadge(screen *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frame)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenW-30, 12)
	screen.DrawImage(fx, op)
}

func (g *game) drawCell(screen *ebiten.Image, x, y, kind int, alpha uint8, highlight bool) {
	px, py := float32(boardX+x*cell), float32(boardY+y*cell)
	art := g.tileArt[kind%len(g.tileArt)]
	op := &ebiten.DrawImageOptions{}
	drawSize := float64(cell - 3)
	op.GeoM.Scale(drawSize/float64(art.Bounds().Dx()), drawSize/float64(art.Bounds().Dy()))
	op.GeoM.Translate(float64(px)+(float64(cell)-drawSize)/2, float64(py)+(float64(cell)-drawSize)/2)
	a := float32(alpha) / 255
	if highlight {
		op.ColorScale.Scale(2, 2, 2, a)
	} else {
		op.ColorScale.ScaleAlpha(a)
	}
	screen.DrawImage(art, op)
}

func (g *game) drawSide(screen *ebiten.Image, x, y int, title string, kind int) {
	vector.DrawFilledRect(screen, float32(x), float32(y), 90, 92, color.RGBA{5, 15, 34, 215}, false)
	vector.StrokeRect(screen, float32(x), float32(y), 90, 92, 2, color.RGBA{102, 211, 230, 120}, false)
	ebitenutil.DebugPrintAt(screen, title, x+26, y+12)
	if kind != noPiece {
		art := g.tileArt[kind%len(g.tileArt)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(44/float64(art.Bounds().Dx()), 44/float64(art.Bounds().Dy()))
		op.GeoM.Translate(float64(x+23), float64(y+34))
		screen.DrawImage(art, op)
		ebitenutil.DebugPrintAt(screen, shapeNames[kind], x+69, y+68)
	} else {
		ebitenutil.DebugPrintAt(screen, "—", x+40, y+52)
	}
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
