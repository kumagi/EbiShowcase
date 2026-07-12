package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width, height = 480, 720
	cols, rows    = 10, 14
	cell          = 34
	ox, oy        = 70, 85
	fallEvery     = 38
	goalPieces    = 14
)

type point struct{ x, y int }

type shape struct {
	name   string
	blocks [4]point
	color  color.RGBA
}

var shapes = []shape{
	{"I", [4]point{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, color.RGBA{61, 190, 219, 255}},
	{"O", [4]point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, color.RGBA{246, 190, 55, 255}},
	{"T", [4]point{{0, 0}, {1, 0}, {2, 0}, {1, 1}}, color.RGBA{178, 99, 219, 255}},
	{"S", [4]point{{1, 0}, {2, 0}, {0, 1}, {1, 1}}, color.RGBA{89, 191, 111, 255}},
	{"Z", [4]point{{0, 0}, {1, 0}, {1, 1}, {2, 1}}, color.RGBA{234, 86, 83, 255}},
	{"J", [4]point{{0, 0}, {0, 1}, {1, 1}, {2, 1}}, color.RGBA{69, 119, 224, 255}},
	{"L", [4]point{{2, 0}, {0, 1}, {1, 1}, {2, 1}}, color.RGBA{239, 148, 53, 255}},
}

type game struct {
	board                      [rows][cols]int
	rng                        *rand.Rand
	bag                        []int
	kind, held                 int
	x, y, rotation, timer      int
	pieces, lines              int
	canHold, cleared, gameOver bool
	message                    string
}

func newGame() *game {
	g := &game{rng: rand.New(rand.NewSource(6505)), held: -1}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = -1
		}
	}
	// A small uneven reef makes the ghost's predicted landing easy to see.
	for _, p := range []point{{0, 13}, {1, 13}, {4, 13}, {5, 13}, {8, 13}, {9, 13}, {0, 12}, {9, 12}} {
		g.board[p.y][p.x] = 7
	}
	g.kind = g.drawFromBag()
	g.spawn()
	g.message = "Use the ghost, and HOLD a piece when it helps."
	return g
}

func (g *game) refillBag() {
	g.bag = []int{0, 1, 2, 3, 4, 5, 6}
	g.rng.Shuffle(len(g.bag), func(i, j int) { g.bag[i], g.bag[j] = g.bag[j], g.bag[i] })
}

func (g *game) drawFromBag() int {
	if len(g.bag) == 0 {
		g.refillBag()
	}
	kind := g.bag[0]
	g.bag = g.bag[1:]
	return kind
}

func (g *game) spawn() {
	g.x, g.y, g.rotation, g.timer = 3, 0, 0, 0
	g.canHold = true
	if !g.canPlace(g.kind, g.rotation, g.x, g.y) {
		g.gameOver = true
		g.message = "No room for the next shape: the reef topped out."
	}
}

func cells(kind, rotation int) [4]point {
	result := shapes[kind].blocks
	for turn := 0; turn < rotation%4; turn++ {
		for i, p := range result {
			result[i] = point{-p.y, p.x}
		}
	}
	minX, minY := result[0].x, result[0].y
	for _, p := range result[1:] {
		if p.x < minX {
			minX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
	}
	for i := range result {
		result[i].x -= minX
		result[i].y -= minY
	}
	return result
}

func (g *game) canPlace(kind, rotation, x, y int) bool {
	for _, p := range cells(kind, rotation) {
		px, py := x+p.x, y+p.y
		if px < 0 || px >= cols || py < 0 || py >= rows || g.board[py][px] >= 0 {
			return false
		}
	}
	return true
}

func (g *game) ghostY() int {
	y := g.y
	for g.canPlace(g.kind, g.rotation, g.x, y+1) {
		y++
	}
	return y
}

func (g *game) hold() {
	if !g.canHold {
		g.message = "HOLD can be used only once before this piece locks."
		return
	}
	if g.held < 0 {
		g.held = g.kind
		g.kind = g.drawFromBag()
	} else {
		g.held, g.kind = g.kind, g.held
	}
	g.x, g.y, g.rotation, g.timer = 3, 0, 0, 0
	g.canHold = false
	g.message = "Swapped HOLD. The ghost was predicted again."
	if !g.canPlace(g.kind, g.rotation, g.x, g.y) {
		g.gameOver = true
	}
}

func (g *game) rotate() {
	next := (g.rotation + 1) % 4
	for _, kick := range []int{0, -1, 1} {
		if g.canPlace(g.kind, next, g.x+kick, g.y) {
			g.rotation, g.x = next, g.x+kick
			return
		}
	}
}

func (g *game) lock() {
	for _, p := range cells(g.kind, g.rotation) {
		g.board[g.y+p.y][g.x+p.x] = g.kind
	}
	g.pieces++
	cleared := g.clearLines()
	g.lines += cleared
	if g.pieces >= goalPieces {
		g.cleared = true
		g.message = "Fourteen pieces planned with BAG, HOLD, and GHOST!"
		return
	}
	if cleared > 0 {
		g.message = fmt.Sprintf("Cleared %d line(s). A fresh ghost is ready.", cleared)
	} else {
		g.message = "Locked on the ghost. Read the next landing."
	}
	g.kind = g.drawFromBag()
	g.spawn()
}

func (g *game) clearLines() int {
	write, count := rows-1, 0
	for read := rows - 1; read >= 0; read-- {
		full := true
		for x := 0; x < cols; x++ {
			if g.board[read][x] < 0 {
				full = false
				break
			}
		}
		if full {
			count++
			continue
		}
		g.board[write] = g.board[read]
		write--
	}
	for write >= 0 {
		for x := 0; x < cols; x++ {
			g.board[write][x] = -1
		}
		write--
	}
	return count
}

func (g *game) Update() error {
	if g.cleared || g.gameOver {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	left := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	turn := inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	drop := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	hold := inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyH)
	if x, y, ok := pressPosition(); ok {
		if y >= 565 && y < 622 {
			hold = true
		} else if y >= 632 {
			switch {
			case x < 120:
				left = true
			case x < 240:
				turn = true
			case x < 360:
				drop = true
			default:
				right = true
			}
		}
	}
	if hold {
		g.hold()
		return nil
	}
	if left && g.canPlace(g.kind, g.rotation, g.x-1, g.y) {
		g.x--
	}
	if right && g.canPlace(g.kind, g.rotation, g.x+1, g.y) {
		g.x++
	}
	if turn {
		g.rotate()
	}
	if drop {
		g.y = g.ghostY()
		g.lock()
		return nil
	}
	g.timer++
	if g.timer >= fallEvery {
		g.timer = 0
		if g.canPlace(g.kind, g.rotation, g.x, g.y+1) {
			g.y++
		} else {
			g.lock()
		}
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 18, 35, 255})
	ebitenutil.DebugPrintAt(screen, "BAG / HOLD / GHOST REEF", 154, 18)
	held := "EMPTY"
	if g.held >= 0 {
		held = shapes[g.held].name
	}
	next := "REFILL"
	if len(g.bag) > 0 {
		next = shapes[g.bag[0]].name
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PLACED %02d/%d  LINES %02d  HOLD %s  NEXT %s", g.pieces, goalPieces, g.lines, held, next), 74, 43)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BAG LEFT %d: %s", len(g.bag), g.bagNames()), 92, 64)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.StrokeRect(screen, px, py, cell-1, cell-1, 1, color.RGBA{43, 63, 91, 255}, false)
			if g.board[y][x] >= 0 {
				c := color.RGBA{74, 112, 139, 255}
				if g.board[y][x] < len(shapes) {
					c = shapes[g.board[y][x]].color
				}
				drawBlock(screen, px, py, c)
			}
		}
	}
	if !g.cleared && !g.gameOver {
		gy := g.ghostY()
		for _, p := range cells(g.kind, g.rotation) {
			px := float32(ox + (g.x+p.x)*cell)
			py := float32(oy + (gy+p.y)*cell)
			vector.StrokeRect(screen, px+4, py+4, cell-9, cell-9, 3, color.RGBA{210, 229, 240, 175}, false)
		}
		for _, p := range cells(g.kind, g.rotation) {
			drawBlock(screen, float32(ox+(g.x+p.x)*cell), float32(oy+(g.y+p.y)*cell), shapes[g.kind].color)
		}
	}
	ebitenutil.DebugPrintAt(screen, g.message, 57, 544)
	button(screen, 22, 570, 436, 48, "HOLD [C / H]", color.RGBA{99, 74, 139, 255})
	button(screen, 8, 632, 108, 62, "LEFT", color.RGBA{51, 89, 126, 255})
	button(screen, 126, 632, 108, 62, "TURN", color.RGBA{51, 89, 126, 255})
	button(screen, 244, 632, 108, 62, "DROP", color.RGBA{230, 168, 60, 255})
	button(screen, 362, 632, 108, 62, "RIGHT", color.RGBA{51, 89, 126, 255})
	if g.cleared || g.gameOver {
		title := "FOURTEEN PIECES COMPLETE!"
		if g.gameOver {
			title = "REEF TOPPED OUT!"
		}
		vector.DrawFilledRect(screen, 42, 270, 396, 160, color.RGBA{4, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 148, 314)
		ebitenutil.DebugPrintAt(screen, g.message, 75, 347)
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE TO RETRY", 148, 388)
	}
}

func (g *game) bagNames() string {
	result := ""
	for _, kind := range g.bag {
		result += shapes[kind].name
	}
	if result == "" {
		return "—"
	}
	return result
}

func drawBlock(screen *ebiten.Image, x, y float32, c color.Color) {
	vector.DrawFilledRect(screen, x+3, y+3, cell-7, cell-7, c, false)
	vector.StrokeRect(screen, x+5, y+5, cell-11, cell-11, 2, color.RGBA{255, 255, 255, 115}, false)
}

func button(screen *ebiten.Image, x, y, w, h int, label string, fill color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), fill, false)
	ebitenutil.DebugPrintAt(screen, label, x+w/2-len(label)*3, y+h/2-4)
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
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Bag, Hold & Ghost Reef — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
