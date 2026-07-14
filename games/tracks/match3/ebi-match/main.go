package main

import (
	"fmt"
	"image/color"
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
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const (
	screenWidth  = 480
	screenHeight = 720
	cols         = 6
	rows         = 7
	cell         = 64
	boardX       = 48
	boardY       = 150
	empty        = -1
	specialNone  = 0
	specialRow   = 1
	specialCol   = 2
	specialColor = 3
	specialArea  = 4
)

var pieceColors = []color.RGBA{
	{239, 93, 87, 255}, {73, 161, 230, 255}, {244, 184, 64, 255},
	{105, 194, 119, 255}, {177, 94, 218, 255},
}

type point struct{ x, y int }

// stage is data, not a rule. A designer can make another level by changing
// these values without rewriting Update or the match scanner.
type stage struct {
	name        string
	moves       int
	targetScore int
	seed        int64
	bonusKind   int // -1 means every gem is worth the same.
	chainBoost  int // Added to the cascade multiplier.
	board       [rows][cols]int
}

var firstStage = stage{
	name:        "CORAL COVE",
	moves:       12,
	targetScore: 650,
	seed:        6106,
	bonusKind:   -1,
	board: [rows][cols]int{
		{0, 1, 2, 3, 4, 0},
		{1, 2, 0, 4, 0, 1},
		{2, 0, 3, 0, 1, 2},
		{3, 0, 0, 1, 0, 3},
		{4, 4, 1, 2, 3, 4},
		{0, 1, 2, 3, 4, 0},
		{1, 2, 3, 4, 0, 1},
	},
}

var stages = []stage{
	firstStage,
	{name: "TIDE TEMPLE", moves: 10, targetScore: 900, seed: 9123, bonusKind: 1, board: [rows][cols]int{
		{0, 2, 1, 4, 3, 0}, {3, 1, 4, 0, 2, 1}, {2, 4, 0, 3, 1, 4}, {1, 3, 2, 4, 0, 2},
		{4, 0, 3, 1, 2, 3}, {0, 2, 4, 3, 1, 0}, {3, 1, 0, 2, 4, 1},
	}},
	{name: "STARLIGHT REEF", moves: 9, targetScore: 1250, seed: 15721, bonusKind: -1, chainBoost: 1, board: [rows][cols]int{
		{4, 1, 3, 0, 2, 4}, {0, 3, 1, 4, 2, 0}, {2, 4, 0, 1, 3, 2}, {1, 0, 4, 3, 2, 1},
		{3, 2, 1, 0, 4, 3}, {4, 1, 3, 2, 0, 4}, {0, 3, 2, 4, 1, 0},
	}},
}

type faller struct {
	kind, special int
	x, fromY, toY int
}

type particle struct {
	x, y, vx, vy, life float64
	kind               int
}

type game struct {
	level            stage
	board            [rows][cols]int
	specials         [rows][cols]int
	rng              *rand.Rand
	cursor, selected point
	hasSelection     bool
	moves, score     int
	combo            int
	message          string
	won, lost        bool
	stageIndex       int
	totalScore       int
	bestScore        int
	tick, shake      int
	particles        []particle
	forge            map[point]int
	forgeAnchor      point
	forgeAnchorValid bool

	// Juice: animate swaps, flash matched cells, then lerp falls. Input is
	// ignored while one of these visual states is active.
	busy      bool
	swapping  bool
	swapBack  bool
	swapFrom  point
	swapTo    point
	swapT     float64
	flash     map[point]bool
	flashLeft int
	falling   bool
	progress  float64
	fallers   []faller
	pending   bool // continue cascade after anim
	audio     *audio.Context
	gate      audiolab.Gate
	pulse     *shaderlab.Pulse
	cam       cameralab.State
	badge     *ebiten.Image
}

func newGame(level stage) *game {
	g := &game{level: level, rng: rand.New(rand.NewSource(level.seed))}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: screenWidth / 2, Y: screenHeight / 2}, ViewW: screenWidth, ViewH: screenHeight}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{245, 190, 69, 255})
	g.board = level.board
	g.moves = level.moves
	g.cursor = point{2, 3}
	g.message = "Swap neighbors. Match 3 or more!"
	return g
}

func newRun(best int) *game {
	g := newGame(stages[0])
	g.bestScore = best
	return g
}

func (g *game) Update() error {
	g.tick++
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .07
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}
	if g.won || g.lost {
		if retryPressed() {
			if g.won && g.stageIndex < len(stages)-1 {
				n := newGame(stages[g.stageIndex+1])
				n.stageIndex, n.totalScore, n.bestScore = g.stageIndex+1, g.totalScore+g.score, g.bestScore
				*g = *n
			} else {
				*g = *newRun(g.bestScore)
			}
		}
		return nil
	}

	if g.busy {
		if g.swapping {
			g.updateSwap()
			return nil
		}
		if g.flashLeft > 0 {
			g.flashLeft--
			if g.flashLeft == 0 {
				g.clear(g.flash)
				g.flash = nil
				g.beginFall()
			}
			return nil
		}
		if g.falling {
			g.progress += 0.1
			if g.progress >= 1 {
				g.falling = false
				g.progress = 0
				g.fallers = nil
				g.refillEmpties()
				g.busy = false
				if g.pending {
					g.resolveMatches()
				}
			}
			return nil
		}
		g.busy = false
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.cursor.x > 0 {
		g.cursor.x--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.cursor.x < cols-1 {
		g.cursor.x++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && g.cursor.y > 0 {
		g.cursor.y--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && g.cursor.y < rows-1 {
		g.cursor.y++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.choose(g.cursor)
	}
	if p, ok := pressedCell(); ok {
		g.cursor = p
		g.choose(p)
	}
	return nil
}

func (g *game) choose(p point) {
	if !g.hasSelection {
		g.selected = p
		g.hasSelection = true
		g.message = "Now choose a neighbor."
		return
	}
	if p == g.selected {
		g.hasSelection = false
		g.message = "Selection canceled."
		return
	}
	if abs(p.x-g.selected.x)+abs(p.y-g.selected.y) != 1 {
		g.selected = p
		g.message = "Choose a touching neighbor."
		return
	}

	a := g.selected
	g.hasSelection = false
	g.busy = true
	g.swapping = true
	g.swapBack = false
	g.swapFrom = a
	g.swapTo = p
	g.swapT = 0
	g.message = "Swapping..."
}

func (g *game) updateSwap() {
	g.swapT += 0.12
	if g.swapT < 1 {
		return
	}
	if g.swapBack {
		g.swapping = false
		g.busy = false
		g.swapT = 0
		g.message = "No line—swap returned. Try again!"
		return
	}

	test := g.board
	a, b := g.swapFrom, g.swapTo
	test[a.y][a.x], test[b.y][b.x] = test[b.y][b.x], test[a.y][a.x]
	matches := scan(test)
	if len(matches) == 0 {
		g.swapBack = true
		g.swapT = 0
		g.message = "No match—returning the pieces."
		return
	}

	g.board = test
	g.play(540)
	g.specials[a.y][a.x], g.specials[b.y][b.x] = g.specials[b.y][b.x], g.specials[a.y][a.x]
	g.swapping = false
	g.swapT = 0
	g.moves--
	g.combo = 0
	if matches[b] {
		g.forgeAnchor = b
		g.forgeAnchorValid = true
	} else if matches[a] {
		g.forgeAnchor = a
		g.forgeAnchorValid = true
	}
	g.resolveMatches()
}

func (g *game) resolveMatches() {
	matches := scan(g.board)
	if len(matches) == 0 {
		g.pending = false
		g.forgeAnchorValid = false
		if !hasValidSwap(g.board) {
			g.makePlayableBoard()
		}
		if g.score >= g.level.targetScore {
			g.won = true
			final := g.totalScore + g.score + g.moves*50
			if final > g.bestScore {
				g.bestScore = final
			}
		} else if g.moves == 0 {
			g.lost = true
		}
		return
	}
	g.forge = nil
	forgedName := ""
	if g.forgeAnchorValid && matches[g.forgeAnchor] && g.specials[g.forgeAnchor.y][g.forgeAnchor.x] == specialNone {
		if special := specialForMatch(g.board, g.forgeAnchor); special != specialNone {
			g.forge = map[point]int{g.forgeAnchor: special}
			forgedName = specialName(special)
		}
	}
	g.forgeAnchorValid = false
	matches = g.expandSpecials(matches)
	g.combo++
	multiplier := g.combo + g.level.chainBoost
	gain := len(matches) * 10 * multiplier
	if g.level.bonusKind >= 0 {
		for p := range matches {
			if g.board[p.y][p.x] == g.level.bonusKind {
				gain += 20
			}
		}
	}
	g.score += gain
	g.play(420 + float64(min(g.combo, 4))*90)
	if g.combo > 1 {
		g.shake = 5 + g.combo
	}
	for p := range matches {
		for i := 0; i < 4; i++ {
			a := float64(i)*math.Pi/2 + float64(p.x+p.y)
			g.particles = append(g.particles, particle{float64(boardX + p.x*cell + cell/2), float64(boardY + p.y*cell + cell/2), math.Cos(a) * (1 + float64(g.combo)*.2), math.Sin(a) * (1 + float64(g.combo)*.2), 24, g.board[p.y][p.x]})
		}
	}
	if forgedName != "" {
		g.message = fmt.Sprintf("%s forged! %d cleared", forgedName, len(matches)-1)
	} else {
		g.message = fmt.Sprintf("%d pieces! Chain x%d", len(matches), g.combo)
	}
	g.flash = matches
	g.flashLeft = 10
	g.busy = true
	g.pending = true
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .06)).Play()
}

func (g *game) beginFall() {
	old := g.board
	oldSpecials := g.specials
	g.fallOnly()
	g.fallers = nil
	for x := 0; x < cols; x++ {
		src := []int{}
		for y := 0; y < rows; y++ {
			if old[y][x] != empty {
				src = append(src, y)
			}
		}
		si := 0
		for y := 0; y < rows; y++ {
			if g.board[y][x] == empty {
				continue
			}
			from := src[si]
			si++
			if from != y {
				g.fallers = append(g.fallers, faller{kind: g.board[y][x], special: oldSpecials[from][x], x: x, fromY: from, toY: y})
			}
		}
	}
	if len(g.fallers) == 0 {
		g.refillEmpties()
		g.busy = false
		if g.pending {
			g.resolveMatches()
		}
		return
	}
	g.falling = true
	g.progress = 0
}

func (g *game) fallOnly() {
	for x := 0; x < cols; x++ {
		write := rows - 1
		for y := rows - 1; y >= 0; y-- {
			if g.board[y][x] != empty {
				g.board[write][x] = g.board[y][x]
				g.specials[write][x] = g.specials[y][x]
				if write != y {
					g.board[y][x] = empty
					g.specials[y][x] = specialNone
				}
				write--
			}
		}
		for write >= 0 {
			g.board[write][x] = empty
			g.specials[write][x] = specialNone
			write--
		}
	}
}

func (g *game) refillEmpties() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if g.board[y][x] == empty {
				g.board[y][x] = g.rng.Intn(len(pieceColors))
				g.specials[y][x] = specialNone
			}
		}
	}
}

func scan(board [rows][cols]int) map[point]bool {
	found := map[point]bool{}
	for y := 0; y < rows; y++ {
		start := 0
		for x := 1; x <= cols; x++ {
			if x == cols || board[y][x] != board[y][start] {
				if board[y][start] != empty && x-start >= 3 {
					for i := start; i < x; i++ {
						found[point{i, y}] = true
					}
				}
				start = x
			}
		}
	}
	for x := 0; x < cols; x++ {
		start := 0
		for y := 1; y <= rows; y++ {
			if y == rows || board[y][x] != board[start][x] {
				if board[start][x] != empty && y-start >= 3 {
					for i := start; i < y; i++ {
						found[point{x, i}] = true
					}
				}
				start = y
			}
		}
	}
	return found
}

func runAt(board [rows][cols]int, p point, dx, dy int) int {
	kind := board[p.y][p.x]
	if kind == empty {
		return 0
	}
	count := 1
	for x, y := p.x-dx, p.y-dy; x >= 0 && x < cols && y >= 0 && y < rows && board[y][x] == kind; x, y = x-dx, y-dy {
		count++
	}
	for x, y := p.x+dx, p.y+dy; x >= 0 && x < cols && y >= 0 && y < rows && board[y][x] == kind; x, y = x+dx, y+dy {
		count++
	}
	return count
}

func specialForMatch(board [rows][cols]int, p point) int {
	horizontal := runAt(board, p, 1, 0)
	vertical := runAt(board, p, 0, 1)
	if horizontal >= 3 && vertical >= 3 {
		return specialArea
	}
	if horizontal >= 5 || vertical >= 5 {
		return specialColor
	}
	if horizontal >= 4 {
		return specialRow
	}
	if vertical >= 4 {
		return specialCol
	}
	return specialNone
}

func specialName(special int) string {
	switch special {
	case specialRow:
		return "ROW ROCKET"
	case specialCol:
		return "COLUMN ROCKET"
	case specialColor:
		return "COLOR WAVE"
	case specialArea:
		return "AREA BOMB"
	default:
		return "SPECIAL"
	}
}

func (g *game) expandSpecials(matches map[point]bool) map[point]bool {
	for {
		changed := false
		for p := range matches {
			special := g.specials[p.y][p.x]
			if special == specialNone {
				continue
			}
			add := func(q point) {
				if !matches[q] {
					matches[q] = true
					changed = true
				}
			}
			switch special {
			case specialRow:
				for x := 0; x < cols; x++ {
					add(point{x, p.y})
				}
			case specialCol:
				for y := 0; y < rows; y++ {
					add(point{p.x, y})
				}
			case specialColor:
				kind := g.board[p.y][p.x]
				for y := 0; y < rows; y++ {
					for x := 0; x < cols; x++ {
						if g.board[y][x] == kind {
							add(point{x, y})
						}
					}
				}
			case specialArea:
				for y := max(0, p.y-1); y <= min(rows-1, p.y+1); y++ {
					for x := max(0, p.x-1); x <= min(cols-1, p.x+1); x++ {
						add(point{x, y})
					}
				}
			}
			// Mark this special as consumed for this resolution pass. Any newly
			// reached special remains available and will be expanded next.
			g.specials[p.y][p.x] = specialNone
		}
		if !changed {
			return matches
		}
	}
}

func (g *game) clear(matches map[point]bool) {
	for p := range matches {
		if special, keep := g.forge[p]; keep {
			g.specials[p.y][p.x] = special
			continue
		}
		g.board[p.y][p.x] = empty
		g.specials[p.y][p.x] = specialNone
	}
	g.forge = nil
}

func (g *game) fallAndRefill() {
	g.fallOnly()
	g.refillEmpties()
}

// makePlayableBoard is the safety net used by a complete level. Random refills
// can occasionally leave no useful swap, so generate a settled, playable board.
func (g *game) makePlayableBoard() {
	for {
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				for {
					kind := g.rng.Intn(len(pieceColors))
					if x >= 2 && g.board[y][x-1] == kind && g.board[y][x-2] == kind {
						continue
					}
					if y >= 2 && g.board[y-1][x] == kind && g.board[y-2][x] == kind {
						continue
					}
					g.board[y][x] = kind
					g.specials[y][x] = specialNone
					break
				}
			}
		}
		if hasValidSwap(g.board) {
			return
		}
	}
}

func hasValidSwap(board [rows][cols]int) bool {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			for _, d := range [...]point{{1, 0}, {0, 1}} {
				nx, ny := x+d.x, y+d.y
				if nx >= cols || ny >= rows {
					continue
				}
				board[y][x], board[ny][nx] = board[ny][nx], board[y][x]
				if len(scan(board)) > 0 {
					return true
				}
				board[y][x], board[ny][nx] = board[ny][nx], board[y][x]
			}
		}
	}
	return false
}

func (g *game) Draw(screen *ebiten.Image) {
	backgrounds := []color.RGBA{{15, 25, 45, 255}, {13, 42, 55, 255}, {30, 18, 57, 255}}
	screen.Fill(backgrounds[g.stageIndex])
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2.4) * float64(g.shake)
	}
	vector.DrawFilledCircle(screen, float32(70+g.stageIndex*155), 80, 95, color.RGBA{40, 100, 120, 35}, true)
	g.drawHUD(screen)
	g.drawEffectBadge(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("REEF %d/3  MOVES %02d  SCORE %04d/%04d", g.stageIndex+1, g.moves, g.score, g.level.targetScore), 75, 67)
	rules := []string{"CLASSIC: plan every swap", "TIDE: blue gems +20", "STORM: every chain starts x2"}
	ebitenutil.DebugPrintAt(screen, rules[g.stageIndex], 145, 113)
	barWidth := float32(360) * float32(g.score) / float32(g.level.targetScore)
	if barWidth > 360 {
		barWidth = 360
	}
	vector.DrawFilledRect(screen, 60, 94, 360, 14, color.RGBA{45, 61, 86, 255}, false)
	vector.DrawFilledRect(screen, 60, 94, barWidth, 14, color.RGBA{245, 190, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, g.message, 112, 132)

	animating := map[point]bool{}
	if g.swapping {
		animating[g.swapFrom] = true
		animating[g.swapTo] = true
	}
	if g.falling {
		for _, f := range g.fallers {
			animating[point{f.x, f.toY}] = true
			py := float32(boardY) + float32(float64(f.fromY)+float64(f.toY-f.fromY)*g.progress)*cell
			px := float64(boardX+f.x*cell+cell/2) + ox
			drawPiece(screen, f.kind, f.special, px, float64(py)+cell/2, cell-6, 1)
		}
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px := float32(boardX + x*cell)
			py := float32(boardY + y*cell)
			if g.board[y][x] == empty {
				vector.StrokeRect(screen, px+8, py+8, cell-16, cell-16, 2, color.RGBA{60, 75, 100, 255}, false)
				continue
			}
			if animating[point{x, y}] {
				continue
			}
			if g.flashLeft > 0 && g.flash[point{x, y}] {
				pulse := 1 + math.Sin(float64(g.tick)*.5)*.12
				trackatlas.DrawTinted(screen, trackatlas.Gem(g.board[y][x]), float64(px+cell/2)+ox, float64(py+cell/2), float64(cell-6)*pulse, 2.4, 2.4, 2.4, 1)
			} else {
				bob := math.Sin(float64(g.tick)*.06+float64(x+y)) * .8
				trackatlas.Draw(screen, trackatlas.Gem(g.board[y][x]), float64(px+3)+ox, float64(py+3)+bob, float64(cell-6))
			}
			special := g.specials[y][x]
			if forged, ok := g.forge[point{x, y}]; ok {
				special = forged
			}
			drawSpecial(screen, special, float64(px+cell/2)+ox, float64(py+cell/2), cell-6, 1)
			if g.hasSelection && g.selected == (point{x, y}) {
				vector.StrokeRect(screen, px+2, py+2, cell-4, cell-4, 6, color.White, false)
			}
			if !g.busy && g.cursor == (point{x, y}) {
				vector.StrokeRect(screen, px+7, py+7, cell-14, cell-14, 3, color.RGBA{25, 32, 48, 255}, false)
			}
		}
	}
	if g.swapping {
		t := smoothstep(g.swapT)
		if g.swapBack {
			t = 1 - t
		}
		a, b := g.swapFrom, g.swapTo
		ax, ay := cellCenter(a)
		bx, by := cellCenter(b)
		drawPiece(screen, g.board[a.y][a.x], g.specials[a.y][a.x], ax+(bx-ax)*t+ox, ay+(by-ay)*t, cell-6, 1)
		drawPiece(screen, g.board[b.y][b.x], g.specials[b.y][b.x], bx+(ax-bx)*t+ox, by+(ay-by)*t, cell-6, 1)
	}
	for _, p := range g.particles {
		c := pieceColors[p.kind%len(pieceColors)]
		vector.DrawFilledCircle(screen, float32(p.x+ox), float32(p.y), float32(2+p.life/10), c, true)
	}
	ebitenutil.DebugPrintAt(screen, "4: ROCKET   5: WAVE   L/T: BOMB", 91, 615)
	ebitenutil.DebugPrintAt(screen, "Tap two neighbors  |  Arrows + Space", 90, 641)
	ebitenutil.DebugPrintAt(screen, "Reach the gold score before moves run out", 74, 667)
	if g.won {
		next := "NEXT REEF"
		if g.stageIndex == len(stages)-1 {
			next = "NEW RUN"
		}
		overlay(screen, fmt.Sprintf("STAGE CLEAR!  BONUS %d\nBEST RUN %d\n\nTAP / SPACE: %s", g.moves*50, g.bestScore, next))
	}
	if g.lost {
		overlay(screen, "OUT OF MOVES\n\nTAP / SPACE TO RETRY")
	}
}

func (g *game) drawHUD(screen *ebiten.Image) {
	label := "EBI MATCH / " + g.level.name
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(145, 16)
		text.Draw(screen, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(screen, label, 145, 28)
}

func (g *game) drawEffectBadge(screen *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.tick)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth-36, 14)
	screen.DrawImage(fx, op)
}

func cellCenter(p point) (float64, float64) {
	return float64(boardX + p.x*cell + cell/2), float64(boardY + p.y*cell + cell/2)
}

func smoothstep(t float64) float64 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return t * t * (3 - 2*t)
}

func drawPiece(screen *ebiten.Image, kind, special int, x, y, size, alpha float64) {
	trackatlas.DrawCentered(screen, trackatlas.Gem(kind), x, y, float64(size))
	drawSpecial(screen, special, x, y, size, alpha)
}

func drawSpecial(screen *ebiten.Image, special int, x, y, size, alpha float64) {
	if special == specialNone {
		return
	}
	a := uint8(255 * alpha)
	white := color.RGBA{255, 255, 255, a}
	gold := color.RGBA{255, 214, 82, a}
	fx, fy := float32(x), float32(y)
	r := float32(size * 0.31)
	switch special {
	case specialRow:
		vector.DrawFilledRect(screen, fx-r, fy-5, r*2, 10, color.RGBA{25, 31, 55, a}, false)
		vector.StrokeLine(screen, fx-r, fy, fx+r, fy, 4, white, false)
		vector.StrokeCircle(screen, fx, fy, 7, 3, gold, false)
	case specialCol:
		vector.DrawFilledRect(screen, fx-5, fy-r, 10, r*2, color.RGBA{25, 31, 55, a}, false)
		vector.StrokeLine(screen, fx, fy-r, fx, fy+r, 4, white, false)
		vector.StrokeCircle(screen, fx, fy, 7, 3, gold, false)
	case specialColor:
		vector.DrawFilledCircle(screen, fx, fy, r*.72, color.RGBA{20, 26, 48, 220}, false)
		vector.StrokeCircle(screen, fx, fy, r*.8, 4, white, false)
		vector.StrokeCircle(screen, fx, fy, r*.45, 3, gold, false)
	case specialArea:
		vector.DrawFilledCircle(screen, fx, fy, r*.72, color.RGBA{25, 31, 55, 235}, false)
		vector.StrokeCircle(screen, fx, fy, r*.8, 4, gold, false)
		vector.StrokeLine(screen, fx-r*.6, fy, fx+r*.6, fy, 3, white, false)
		vector.StrokeLine(screen, fx, fy-r*.6, fx, fy+r*.6, 3, white, false)
	}
}

func pressedCell() (point, bool) {
	x, y, ok := pointerPress()
	if !ok || x < boardX || y < boardY || x >= boardX+cols*cell || y >= boardY+rows*cell {
		return point{}, false
	}
	return point{(x - boardX) / cell, (y - boardY) / cell}, true
}

func pointerPress() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	touches := inpututil.AppendJustPressedTouchIDs(nil)
	if len(touches) > 0 {
		x, y := ebiten.TouchPosition(touches[0])
		return x, y, true
	}
	return 0, 0, false
}

func retryPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) {
		return true
	}
	_, _, ok := pointerPress()
	return ok
}

func overlay(screen *ebiten.Image, message string) {
	vector.DrawFilledRect(screen, 50, 278, 380, 160, color.RGBA{5, 16, 34, 245}, false)
	vector.StrokeRect(screen, 50, 278, 380, 160, 4, color.RGBA{245, 190, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, message, 115, 330)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (g *game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebi Match — Ebitengine")
	if err := ebiten.RunGame(newGame(firstStage)); err != nil {
		panic(err)
	}
}
