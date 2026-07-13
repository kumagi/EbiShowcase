package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const (
	width, height  = 480, 720
	cols, rows     = 6, 10
	cell           = 40
	boardX, boardY = 55, 104
	empty          = -1
	garbage        = 4
	phasePlayer    = 0
	phaseClear     = 1
	phaseGravity   = 2
)

var names = []string{"RED", "BLUE", "YELLOW", "GREEN"}

type point struct{ x, y int }
type spark struct {
	x, y, vx, vy, life float64
	kind               int
}
type duelRule struct {
	name                         string
	goal, missLimit, missGarbage int
	bg                           color.RGBA
}

var duels = []duelRule{
	{"SUNNY LAGOON", 6, 3, 1, color.RGBA{10, 35, 55, 255}},
	{"CURRENT CAVE", 9, 3, 2, color.RGBA{8, 48, 52, 255}},
	{"STORM PALACE", 12, 2, 2, color.RGBA{32, 18, 58, 255}},
}

var offsets = []point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
var directions = []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

type game struct {
	board                        [rows][cols]int
	x, y, rotation, timer        int
	pivotKind, childKind         int
	nextPivotKind, nextChildKind int
	round, pairs, phase          int
	chain, score, sent           int
	opponentGarbage, misses      int
	marked                       map[point]bool
	clear, over                  bool
	stageWon                     bool
	duel, best, tick, shake      int
	sparks                       []spark
	message                      string
}

func newGame() *game {
	g := &game{}
	g.seedPuzzle(0)
	return g
}

func (g *game) seedPuzzle(round int) {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.board[y][x] = empty
		}
	}
	g.round = round
	if round%2 == 0 {
		g.board[9][0], g.board[9][1], g.board[8][0] = 0, 0, 0
		g.board[7][0], g.board[6][0], g.board[5][0] = 1, 1, 1
		g.nextPivotKind, g.nextChildKind = 0, 1
	} else {
		g.board[9][5], g.board[9][4], g.board[8][5] = 3, 3, 3
		g.board[7][5], g.board[6][5], g.board[5][5] = 2, 2, 2
		g.nextPivotKind, g.nextChildKind = 3, 2
	}
	g.phase = phasePlayer
	g.marked = map[point]bool{}
	g.spawn()
	g.message = fmt.Sprintf("REEF %d: complete %s, then watch the second color fall.", round+1, names[g.pivotKind])
}

func (g *game) spawn() {
	g.pivotKind, g.childKind = g.nextPivotKind, g.nextChildKind
	// Each guided reef keeps the same rescue pair in the preview queue, so a
	// missed placement remains recoverable until the miss limit is reached.
	g.nextPivotKind, g.nextChildKind = g.pivotKind, g.childKind
	g.x, g.y, g.rotation, g.timer = cols/2, 1, 0, 0
	g.phase = phasePlayer
	if g.blocked(g.x, g.y, g.rotation) {
		g.over = true
		g.message = "The spawn area is blocked by the stack."
	}
}

func (g *game) pairCells(x, y, rotation int) [2]point {
	o := offsets[rotation%4]
	return [2]point{{x, y}, {x + o.x, y + o.y}}
}

func (g *game) blocked(x, y, rotation int) bool {
	for _, p := range g.pairCells(x, y, rotation) {
		if p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || g.board[p.y][p.x] != empty {
			return true
		}
	}
	return false
}

func (g *game) Update() error {
	g.tick++
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .05
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.clear || g.over || g.stageWon {
		if retryPressed() {
			if g.stageWon {
				n := newGame()
				n.duel = g.duel + 1
				n.best = g.best
				n.seedPuzzle(0)
				*g = *n
			} else {
				best := g.best
				*g = *newGame()
				g.best = best
			}
		}
		return nil
	}
	if g.phase != phasePlayer {
		g.updateResolution()
		return nil
	}

	left := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	right := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	turn := inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	drop := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if x, y, ok := pressPosition(); ok && y >= 590 {
		switch {
		case x < 120:
			left = true
		case x < 240:
			turn = true
		case x < 360:
			right = true
		default:
			drop = true
		}
	}
	if left && !g.blocked(g.x-1, g.y, g.rotation) {
		g.x--
	}
	if right && !g.blocked(g.x+1, g.y, g.rotation) {
		g.x++
	}
	if turn {
		g.rotate()
	}
	if drop {
		for !g.blocked(g.x, g.y+1, g.rotation) {
			g.y++
		}
		g.lockPair()
		return nil
	}
	g.timer++
	if g.timer >= 85 {
		g.timer = 0
		if !g.blocked(g.x, g.y+1, g.rotation) {
			g.y++
		} else {
			g.lockPair()
		}
	}
	return nil
}

func (g *game) rotate() {
	next := (g.rotation + 1) % 4
	for _, kick := range []int{0, -1, 1} {
		if !g.blocked(g.x+kick, g.y, next) {
			g.x, g.rotation = g.x+kick, next
			return
		}
	}
}

func (g *game) lockPair() {
	cells := g.pairCells(g.x, g.y, g.rotation)
	g.board[cells[0].y][cells[0].x] = g.pivotKind
	g.board[cells[1].y][cells[1].x] = g.childKind
	g.pairs++
	g.chain = 0
	g.marked = findGroups(g.board)
	if len(g.marked) == 0 {
		g.misses++
		g.dropIncomingGarbage()
		if g.misses >= duels[g.duel].missLimit {
			g.over = true
			g.message = fmt.Sprintf("%d no-chain moves let the rival overwhelm the reef.", duels[g.duel].missLimit)
			return
		}
		g.message = fmt.Sprintf("No group. Rival dropped garbage! Misses %d/%d.", g.misses, duels[g.duel].missLimit)
		g.spawn()
		return
	}
	g.chain = 1
	g.phase = phaseClear
	g.timer = 32
	g.message = "Groups found: simultaneous clear is queued."
}

func (g *game) updateResolution() {
	g.timer--
	if g.timer > 0 {
		return
	}
	switch g.phase {
	case phaseClear:
		burst := g.marked
		cleared := g.clearMarkedAndAdjacentGarbage()
		g.score += cleared * 10 * g.chain
		g.shake = 4 + g.chain*2
		for p := range burst {
			for i := 0; i < 5; i++ {
				a := float64(i) * math.Pi * .4
				g.sparks = append(g.sparks, spark{float64(boardX + p.x*cell + cell/2), float64(boardY + p.y*cell + cell/2), math.Cos(a) * (1 + float64(g.chain)*.3), math.Sin(a) * (1 + float64(g.chain)*.3), 26, (p.x + p.y) % 4})
			}
		}
		g.phase = phaseGravity
		g.timer = 28
		g.message = fmt.Sprintf("CHAIN %d cleared %d: +%d. Gravity next.", g.chain, cleared, cleared*10*g.chain)
	case phaseGravity:
		g.compactColumns()
		g.marked = findGroups(g.board)
		if len(g.marked) > 0 {
			g.chain++
			g.phase = phaseClear
			g.timer = 32
			g.message = fmt.Sprintf("Rescan found another group: CHAIN %d!", g.chain)
			return
		}
		g.finishResolution()
	}
}

func findGroups(board [rows][cols]int) map[point]bool {
	visited := map[point]bool{}
	marked := map[point]bool{}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			start := point{x, y}
			kind := board[y][x]
			if kind == empty || kind == garbage || visited[start] {
				continue
			}
			queue := []point{start}
			visited[start] = true
			group := []point{}
			for len(queue) > 0 {
				p := queue[0]
				queue = queue[1:]
				group = append(group, p)
				for _, d := range directions {
					n := point{p.x + d.x, p.y + d.y}
					if n.x < 0 || n.x >= cols || n.y < 0 || n.y >= rows || visited[n] || board[n.y][n.x] != kind {
						continue
					}
					visited[n] = true
					queue = append(queue, n)
				}
			}
			if len(group) >= 4 {
				for _, p := range group {
					marked[p] = true
				}
			}
		}
	}
	return marked
}

func (g *game) clearMarkedAndAdjacentGarbage() int {
	toClear := map[point]bool{}
	for p := range g.marked {
		toClear[p] = true
		for _, d := range directions {
			n := point{p.x + d.x, p.y + d.y}
			if n.x >= 0 && n.x < cols && n.y >= 0 && n.y < rows && g.board[n.y][n.x] == garbage {
				toClear[n] = true
			}
		}
	}
	for p := range toClear {
		g.board[p.y][p.x] = empty
	}
	g.marked = map[point]bool{}
	return len(toClear)
}

func (g *game) compactColumns() {
	for x := 0; x < cols; x++ {
		write := rows - 1
		for read := rows - 1; read >= 0; read-- {
			if g.board[read][x] == empty {
				continue
			}
			g.board[write][x] = g.board[read][x]
			if write != read {
				g.board[read][x] = empty
			}
			write--
		}
		for write >= 0 {
			g.board[write][x] = empty
			write--
		}
	}
}

func (g *game) finishResolution() {
	allClearBonus := 0
	if g.boardEmpty() {
		allClearBonus = 2
	}
	sent := g.chain*2 + allClearBonus
	g.sent += sent
	g.opponentGarbage += sent
	g.message = fmt.Sprintf("Settled: chain %d sent %d garbage", g.chain, sent)
	if allClearBonus > 0 {
		g.message += " (ALL CLEAR +2)!"
	}
	if g.opponentGarbage >= duels[g.duel].goal {
		runScore := g.score + g.sent*50 - g.misses*25
		if runScore > g.best {
			g.best = runScore
		}
		if g.duel == len(duels)-1 {
			g.clear = true
			g.message = "All three rival reefs overflowed!"
		} else {
			g.stageWon = true
			g.message = "Rival overflowed! The next reef is waiting."
		}
		return
	}
	if g.boardEmpty() {
		g.seedPuzzle(g.round + 1)
		return
	}
	g.spawn()
}

func (g *game) boardEmpty() bool {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if g.board[y][x] != empty {
				return false
			}
		}
	}
	return true
}

func (g *game) dropIncomingGarbage() {
	for n := 0; n < duels[g.duel].missGarbage; n++ {
		column := (g.pairs*3 + 1 + n*2) % cols
		placed := false
		for y := rows - 1; y >= 0; y-- {
			if g.board[y][column] == empty {
				g.board[y][column] = garbage
				placed = true
				break
			}
		}
		if !placed {
			g.over = true
			g.message = "Incoming garbage blocked the whole column."
			return
		}
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	rule := duels[g.duel]
	screen.Fill(rule.bg)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2.1) * float64(g.shake)
	}
	vector.DrawFilledCircle(screen, float32(65+g.duel*160), 75, 85, color.RGBA{60, 120, 130, 35}, true)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("EBI CHAIN / DUEL %d/3 %s", g.duel+1, rule.name), 130, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %05d CHAIN %d SENT %d MISSES %d/%d BEST %d", g.score, g.chain, g.sent, g.misses, rule.missLimit, g.best), 55, 44)
	ebitenutil.DebugPrintAt(screen, g.message, 48, 70)
	vector.DrawFilledRect(screen, boardX-4, boardY-4, cols*cell+8, rows*cell+8, color.RGBA{32, 47, 69, 255}, false)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			vector.StrokeRect(screen, px, py, cell, cell, 1, color.RGBA{62, 78, 102, 255}, false)
			kind := g.board[y][x]
			if kind == empty {
				continue
			}
			drawPieceAt(screen, x, y, kind, "", ox, math.Sin(float64(g.tick)*.07+float64(x+y))*.6)
			if g.marked[point{x, y}] {
				pulse := float32(2 + math.Sin(float64(g.tick)*.5)*3)
				vector.StrokeCircle(screen, px+cell/2+float32(ox), py+cell/2, cell/2-pulse, 4, color.White, false)
			}
		}
	}
	if g.phase == phasePlayer && !g.clear && !g.over {
		cells := g.pairCells(g.x, g.y, g.rotation)
		bob := math.Sin(float64(g.tick)*.12) * 1.5
		drawPieceAt(screen, cells[0].x, cells[0].y, g.pivotKind, "P", ox, bob)
		drawPieceAt(screen, cells[1].x, cells[1].y, g.childKind, "C", ox, bob)
	}

	vector.DrawFilledRect(screen, 330, 110, 120, 88, color.RGBA{30, 52, 73, 255}, false)
	ebitenutil.DebugPrintAt(screen, "NEXT PAIR", 354, 124)
	trackatlas.DrawCentered(screen, trackatlas.Gem(g.nextPivotKind), 365, 166, 26)
	trackatlas.DrawCentered(screen, trackatlas.Gem(g.nextChildKind), 410, 166, 26)
	ebitenutil.DebugPrintAt(screen, "RIVAL GARBAGE", 343, 225)
	for i := 0; i < rule.goal; i++ {
		row, col := i/3, i%3
		cx, cy := float32(350+col*38), float32(260+row*38)
		if i < g.opponentGarbage {
			trackatlas.DrawCentered(screen, "gem-trash", float64(cx), float64(cy), 26)
		} else {
			vector.DrawFilledCircle(screen, cx, cy, 14, color.RGBA{45, 56, 70, 255}, false)
		}
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%02d/%02d", min(g.opponentGarbage, rule.goal), rule.goal), 370, 421)
	ebitenutil.DebugPrintAt(screen, "2 per chain + ALL CLEAR 2", 324, 458)
	phaseNames := []string{"AIM", "CLEAR ALL", "GRAVITY + RESCAN"}
	ebitenutil.DebugPrintAt(screen, "PHASE: "+phaseNames[g.phase], 331, 493)

	labels := []string{"LEFT", "TURN", "RIGHT", "DROP"}
	for i, label := range labels {
		vector.DrawFilledRect(screen, float32(i*120+5), 590, 110, 72, color.RGBA{51, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, i*120+37, 622)
	}
	ebitenutil.DebugPrintAt(screen, "Arrows / A,W,D / Space / mouse / touch", 102, 684)
	for _, p := range g.sparks {
		c := []color.RGBA{{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255}, {105, 194, 119, 255}}[p.kind%4]
		vector.DrawFilledCircle(screen, float32(p.x+ox), float32(p.y), float32(2+p.life/10), c, true)
	}
	if g.clear {
		overlay(screen, fmt.Sprintf("REEF CHAMPION! BEST %d\n\nTAP / ENTER: NEW RUN", g.best))
	} else if g.stageWon {
		overlay(screen, fmt.Sprintf("DUEL %d CLEAR! BEST %d\n\nTAP / ENTER: NEXT REEF", g.duel+1, g.best))
	} else if g.over {
		overlay(screen, "YOUR REEF LOST THE DUEL\n\nTAP / ENTER TO RETRY")
	}
}

func drawPiece(screen *ebiten.Image, x, y, kind int, label string) {
	drawPieceAt(screen, x, y, kind, label, 0, 0)
}

func drawPieceAt(screen *ebiten.Image, x, y, kind int, label string, ox, oy float64) {
	cx := float32(boardX + x*cell + cell/2)
	cy := float32(boardY + y*cell + cell/2)
	sprite := trackatlas.Gem(kind)
	if kind == garbage {
		sprite = "gem-trash"
	}
	trackatlas.DrawCentered(screen, sprite, float64(cx)+ox, float64(cy)+oy, float64(cell-6))
	if kind == garbage {
		ebitenutil.DebugPrintAt(screen, "X", int(cx)-3, int(cy)-5)
	} else if label != "" {
		ebitenutil.DebugPrintAt(screen, label, int(cx)-3, int(cy)-5)
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

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 96, 328)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Chain Reef Duel — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
