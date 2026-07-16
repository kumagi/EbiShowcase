// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
// Package reversiui turns the pure Reversi rules into the five playable
// lessons. The final lesson adds responsive presentation, a CPU ladder, and
// small presentation state without hiding the testable board rules.
package reversiui

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/ogfont"
	"github.com/kumagi/EbiShowcase/internal/reversi"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Variant int

const (
	BoardGrid Variant = iota
	LegalMoves
	FlipStones
	PassAndScore
	CPUEvaluation
)

type difficulty int

const (
	Friendly difficulty = iota
	Positional
	Scout
)

const (
	flipFrames = 18
	bestPrefix = "ebiShowcase.reversi.best."
)

type flipAnimation struct {
	cell      reversi.Move
	from, to  reversi.Player
	remaining int
}

type boardView struct {
	x, y, size float32
	portrait   bool
}

type Game struct {
	variant      Variant
	board        reversi.Board
	turn         reversi.Player
	legal        []reversi.Move
	cursor       reversi.Move
	lastMove     reversi.Move
	lastPulse    int
	flips        []flipAnimation
	passCount    int
	cpuWait      int
	message      string
	over         bool
	result       string
	difficulty   difficulty
	bestMargin   int
	moves        int
	lang         string
	viewW, viewH int
	audio        *audio.Context
	gate         audiolab.Gate
	pulse        *shaderlab.Pulse
	cam          cameralab.State
	badge        *ebiten.Image
}

var (
	fontOnce sync.Once
	fontBase *opentype.Font
	fontErr  error
	faces    = map[float64]font.Face{}
	artOnce  sync.Once
	backdrop *ebiten.Image
)

//go:embed assets/reversi-championship.png
var backdropPNG []byte

func loadBackdrop() {
	artOnce.Do(func() {
		decoded, _, err := image.Decode(bytes.NewReader(backdropPNG))
		if err != nil {
			panic(err)
		}
		backdrop = ebiten.NewImageFromImage(decoded)
	})
}

func uiFace(size float64) font.Face {
	fontOnce.Do(func() { fontBase, fontErr = opentype.Parse(ogfont.NotoSansJP) })
	if fontErr != nil {
		panic(fontErr)
	}
	if face := faces[size]; face != nil {
		return face
	}
	face, err := opentype.NewFace(fontBase, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		panic(err)
	}
	faces[size] = face
	return face
}

func New(variant Variant) *Game {
	loadBackdrop()
	g := &Game{variant: variant, difficulty: Positional, lang: browserLanguage()}
	g.audio = audiolab.Context()
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{ViewW: 900, ViewH: 720}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{255, 211, 112, 255})
	g.bestMargin = storedBest(bestPrefix + g.difficultyName())
	g.reset()
	return g
}

func (g *Game) reset() {
	g.board = reversi.NewBoard()
	g.turn = reversi.Black
	g.legal = reversi.ValidMoves(g.board, g.turn)
	g.cursor = g.legal[0]
	g.lastMove = reversi.Move{X: -1, Y: -1}
	g.lastPulse = 0
	g.flips = nil
	g.passCount = 0
	g.cpuWait = 22
	g.over = false
	g.result = ""
	g.moves = 0
	g.message = g.initialMessage()
}

func (g *Game) jp() bool { return g.lang == "ja" }

func (g *Game) tr(en, ja string) string {
	if g.jp() {
		return ja
	}
	return en
}

func (g *Game) difficultyName() string {
	switch g.difficulty {
	case Friendly:
		return "FRIENDLY"
	case Scout:
		return "SCOUT"
	default:
		return "POSITION"
	}
}

func (g *Game) difficultyHint() string {
	switch g.difficulty {
	case Friendly:
		return g.tr("first legal move", "置ける場所を左上から選ぶ")
	case Scout:
		return g.tr("one reply ahead", "相手の次の手まで考える")
	default:
		return g.tr("values safe squares", "角と安全なマスを大切にする")
	}
}

func (g *Game) initialMessage() string {
	switch g.variant {
	case BoardGrid:
		return g.tr("Click a cell to write BLACK. R resets.", "マスをクリックして黒い石を置こう。Rでやり直し。")
	case LegalMoves:
		return g.tr("Blue dots are legal moves. Click one.", "青い点が置ける場所。1つ選んで置こう。")
	case FlipStones:
		return g.tr("Take turns locally. Bracketed stones flip.", "交代で置こう。はさんだ石がひっくり返る。")
	case PassAndScore:
		return g.tr("You are BLACK. CPU takes the first legal move.", "あなたは黒。CPUは最初に見つけた手を選ぶ。")
	default:
		return g.tr("Choose a difficulty, then play BLACK.", "難しさを選んで、黒の石で対戦しよう。")
	}
}

func (g *Game) Update() error {
	if g.lastPulse > 0 {
		g.lastPulse--
	}
	for i := range g.flips {
		g.flips[i].remaining--
	}
	for i := len(g.flips) - 1; i >= 0; i-- {
		if g.flips[i].remaining <= 0 {
			g.flips = append(g.flips[:i], g.flips[i+1:]...)
		}
	}

	if g.variant == CPUEvaluation {
		if keyDifficulty, ok := keyboardDifficulty(); ok {
			g.setDifficulty(keyDifficulty)
		}
		if x, y, ok := press(); ok && g.difficultyAt(x, y) >= 0 {
			g.setDifficulty(difficulty(g.difficultyAt(x, y)))
			return nil
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || g.restartTouched() {
		g.reset()
		return nil
	}
	if g.over {
		if anyPressed() {
			g.reset()
		}
		return nil
	}

	if g.variant >= PassAndScore && g.turn == reversi.White {
		if g.cpuWait > 0 {
			g.cpuWait--
			return nil
		}
		move, ok := g.cpuMove()
		if !ok {
			g.passOrFinish()
			return nil
		}
		g.play(move)
		return nil
	}

	if g.variant == BoardGrid {
		if x, y, ok := press(); ok {
			if cell, inside := g.cellAt(x, y); inside && g.board[cell.Y][cell.X] == reversi.Empty {
				g.board[cell.Y][cell.X] = reversi.Black
				g.lastMove, g.lastPulse = cell, 16
				g.message = fmt.Sprintf("board[%d][%d] = BLACK", cell.Y, cell.X)
			}
		}
		return nil
	}

	g.moveCursor()
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if containsMove(g.legal, g.cursor) {
			g.play(g.cursor)
		}
		return nil
	}
	if x, y, ok := press(); ok {
		if cell, inside := g.cellAt(x, y); inside && containsMove(g.legal, cell) {
			g.play(cell)
		}
	}
	return nil
}

func (g *Game) setDifficulty(next difficulty) {
	if next < Friendly || next > Scout || g.difficulty == next {
		return
	}
	g.difficulty = next
	g.bestMargin = storedBest(bestPrefix + g.difficultyName())
	g.reset()
}

func (g *Game) cpuMove() (reversi.Move, bool) {
	if g.variant != CPUEvaluation {
		return reversi.ChooseFirst(g.board, reversi.White)
	}
	switch g.difficulty {
	case Friendly:
		return reversi.ChooseFirst(g.board, reversi.White)
	case Scout:
		move, ok, _ := reversi.ChooseLookahead(g.board, reversi.White)
		return move, ok
	default:
		move, ok, _ := reversi.ChooseBest(g.board, reversi.White)
		return move, ok
	}
}

func (g *Game) moveCursor() {
	dx, dy := 0, 0
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dx = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dx = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		dy = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		dy = 1
	}
	if dx != 0 || dy != 0 {
		g.cursor.X = min(reversi.Size-1, max(0, g.cursor.X+dx))
		g.cursor.Y = min(reversi.Size-1, max(0, g.cursor.Y+dy))
	}
}

func (g *Game) play(move reversi.Move) {
	player := g.turn
	captured := reversi.Apply(&g.board, player, move)
	if len(captured) == 0 {
		return
	}
	g.playSE(420 + float64(len(captured))*55)
	for i, cell := range captured {
		g.flips = append(g.flips, flipAnimation{cell: cell, from: reversi.Opponent(player), to: player, remaining: flipFrames + i*2})
	}
	g.lastMove, g.lastPulse = move, 22
	g.moves++
	g.passCount = 0
	g.turn = reversi.Opponent(g.turn)
	g.cpuWait = 22
	g.refreshTurn()
}

func (g *Game) playSE(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .05)).Play()
}

func (g *Game) refreshTurn() {
	if reversi.Full(g.board) {
		g.finish()
		return
	}
	g.legal = reversi.ValidMoves(g.board, g.turn)
	if len(g.legal) == 0 {
		g.passOrFinish()
		return
	}
	g.cursor = g.legal[0]
	if g.turn == reversi.Black {
		g.message = g.tr(fmt.Sprintf("Your turn: %d legal moves. Arrows + Space also work.", len(g.legal)), fmt.Sprintf("あなたの番：置ける場所は%d個。矢印＋Spaceでも置ける。", len(g.legal)))
	}
}

func (g *Game) passOrFinish() {
	g.passCount++
	if g.passCount >= 2 {
		g.finish()
		return
	}
	passed := g.turn
	g.turn = reversi.Opponent(g.turn)
	g.legal = reversi.ValidMoves(g.board, g.turn)
	g.cursor = reversi.Move{X: -1, Y: -1}
	if len(g.legal) > 0 {
		g.cursor = g.legal[0]
	}
	g.cpuWait = 22
	g.message = g.tr(fmt.Sprintf("%s has no move: PASS", playerName(passed)), fmt.Sprintf("%sは置けないのでパス", playerNameJP(passed)))
}

func (g *Game) finish() {
	g.over = true
	black, white := reversi.Count(g.board, reversi.Black), reversi.Count(g.board, reversi.White)
	if black > white {
		margin := black - white
		if margin > g.bestMargin {
			g.bestMargin = margin
			storeBest(bestPrefix+g.difficultyName(), margin)
		}
		g.result = g.tr(fmt.Sprintf("BLACK WINS %d–%d", black, white), fmt.Sprintf("黒の勝ち %d–%d", black, white))
	} else if white > black {
		g.result = g.tr(fmt.Sprintf("WHITE WINS %d–%d", white, black), fmt.Sprintf("白の勝ち %d–%d", white, black))
	} else {
		g.result = g.tr(fmt.Sprintf("DRAW %d–%d", black, white), fmt.Sprintf("引き分け %d–%d", black, white))
	}
	g.message = g.tr("Enter, R, or tap REPLAY for another match.", "Enter、R、またはREPLAYで次の対戦へ。")
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 20, 37, 255})
	v := g.view()
	w, h := g.dimensions()
	g.drawBackdrop(screen, w, h, v)
	g.drawEffectBadge(screen)
	cell := v.size / reversi.Size

	g.drawTop(screen, w, v)
	vector.DrawFilledRect(screen, v.x-17, v.y-11, v.size+34, v.size+34, color.NRGBA{3, 8, 18, 205}, false)
	vector.DrawFilledRect(screen, v.x-9, v.y-9, v.size+18, v.size+18, color.RGBA{172, 113, 43, 255}, false)
	vector.DrawFilledRect(screen, v.x-5, v.y-5, v.size+10, v.size+10, color.RGBA{13, 51, 65, 255}, false)
	for y := 0; y < reversi.Size; y++ {
		for x := 0; x < reversi.Size; x++ {
			cx, cy := v.x+float32(x)*cell, v.y+float32(y)*cell
			shade := uint8(38)
			if (x+y)%2 == 0 {
				shade = 48
			}
			vector.DrawFilledRect(screen, cx, cy, cell, cell, color.RGBA{10, shade, 61, 255}, false)
			vector.StrokeRect(screen, cx, cy, cell, cell, maxf(1, cell*.018), color.NRGBA{129, 220, 213, 135}, false)
			move := reversi.Move{X: x, Y: y}
			if containsMove(g.legal, move) && !g.over && g.variant != BoardGrid {
				vector.DrawFilledCircle(screen, cx+cell/2, cy+cell/2, maxf(5, cell*.13), color.RGBA{76, 213, 239, 180}, false)
			}
			if g.turn == reversi.Black && g.cursor == move && !g.over {
				vector.StrokeRect(screen, cx+3, cy+3, cell-6, cell-6, maxf(2, cell*.045), color.RGBA{255, 220, 112, 255}, false)
			}
			if stone := g.board[y][x]; stone != reversi.Empty {
				stoneColor, radius := g.stoneFrame(move, stone, cell*.39)
				vector.DrawFilledCircle(screen, cx+cell/2+cell*.035, cy+cell/2+cell*.09, radius*.91, color.RGBA{3, 9, 14, 105}, true)
				vector.DrawFilledCircle(screen, cx+cell/2, cy+cell/2, radius, stoneColor, false)
				vector.StrokeCircle(screen, cx+cell/2, cy+cell/2, radius, maxf(1, cell*.025), color.RGBA{255, 255, 255, 125}, false)
				vector.DrawFilledCircle(screen, cx+cell*.39, cy+cell*.36, maxf(1.5, radius*.12), color.RGBA{255, 255, 255, 95}, true)
			}
			if g.lastMove == move && g.lastPulse > 0 {
				pulse := cell*.33 + float32(22-g.lastPulse)*cell*.014
				vector.StrokeCircle(screen, cx+cell/2, cy+cell/2, pulse, maxf(1, cell*.025), color.RGBA{255, 207, 83, 230}, false)
			}
		}
	}
	g.drawCoordinates(screen, v, cell)
	g.drawInfo(screen, w, h, v)
	if g.over {
		g.drawResult(screen, v)
	}
}

func (g *Game) drawBackdrop(screen *ebiten.Image, w, h int, v boardView) {
	if backdrop != nil {
		b := backdrop.Bounds()
		sx, sy := float64(w)/float64(b.Dx()), float64(h)/float64(b.Dy())
		scale := math.Max(sx, sy)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate((float64(w)-float64(b.Dx())*scale)/2, (float64(h)-float64(b.Dy())*scale)/2)
		op.Filter = ebiten.FilterLinear
		screen.DrawImage(backdrop, op)
	}
	vector.DrawFilledRect(screen, 0, 0, float32(w), 106, color.NRGBA{2, 9, 23, 208}, false)
	if g.variant == CPUEvaluation {
		drawLabel(screen, g.tr("CHAMPIONSHIP TABLE", "チャンピオンシップ"), max(18, w-230), h-24, 10, color.RGBA{255, 211, 112, 110})
	}
}

func (g *Game) drawEffectBadge(screen *ebiten.Image) {
	if g.variant != CPUEvaluation || g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.moves)*.2) {
		return
	}
	w, _ := g.dimensions()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w-36), 12)
	screen.DrawImage(fx, op)
}

func (g *Game) drawTop(screen *ebiten.Image, w int, v boardView) {
	drawLabel(screen, g.tr("EBI REVERSI", "えびリバーシ"), 18, 25, 18, color.RGBA{245, 250, 255, 255})
	if g.variant != CPUEvaluation {
		drawLabel(screen, g.title(), 18, 47, 12, color.RGBA{144, 214, 190, 255})
		return
	}
	buttonW := float32(min(125, max(76, (w-58)/4)))
	for i, d := range []difficulty{Friendly, Positional, Scout} {
		x := float32(16) + float32(i)*(buttonW+6)
		selected := d == g.difficulty
		drawButton(screen, x, 38, buttonW, 35, []string{"1 FRIEND", "2 POSITION", "3 SCOUT"}[i], selected)
	}
	drawButton(screen, float32(w)-94, 38, 78, 35, "REPLAY", false)
	if v.portrait {
		drawLabel(screen, g.tr("CPU: "+g.difficultyHint(), "CPU："+g.difficultyHint()), 18, 96, 12, color.RGBA{184, 211, 233, 255})
	}
}

func (g *Game) drawInfo(screen *ebiten.Image, w, h int, v boardView) {
	black, white := reversi.Count(g.board, reversi.Black), reversi.Count(g.board, reversi.White)
	if v.portrait {
		drawLabel(screen, fmt.Sprintf("● %02d    ○ %02d    %s", black, white, g.tr("TURN", "手番")), 18, 119, 13, color.RGBA{245, 250, 255, 255})
		g.drawScoreMap(screen, v, true)
		return
	}
	x := int(v.x + v.size + 24)
	y := int(v.y + 12)
	drawLabel(screen, fmt.Sprintf("● BLACK %02d", black), x, y, 16, color.RGBA{245, 250, 255, 255})
	drawLabel(screen, fmt.Sprintf("○ WHITE %02d", white), x, y+28, 16, color.RGBA{245, 250, 255, 255})
	drawLabel(screen, g.tr("TURN "+playerName(g.turn), "手番 "+playerNameJP(g.turn)), x, y+58, 14, color.RGBA{184, 211, 233, 255})
	if g.variant == CPUEvaluation {
		drawLabel(screen, "CPU "+g.difficultyName(), x, y+88, 14, color.RGBA{255, 211, 112, 255})
		drawLabel(screen, g.difficultyHint(), x, y+108, 12, color.RGBA{184, 211, 233, 255})
		drawLabel(screen, fmt.Sprintf("MAP EVAL %d", reversi.Evaluate(g.board, reversi.Black)), x, y+136, 13, color.RGBA{245, 250, 255, 255})
		g.drawScoreMap(screen, boardView{x: float32(x), y: float32(y + 154)}, false)
		drawLabel(screen, g.tr("BEST WIN +", "最高勝ち差 +")+fmt.Sprint(g.bestMargin), x, y+360, 13, color.RGBA{255, 211, 112, 255})
	} else {
		drawLabel(screen, fmt.Sprintf("LEGAL %d", len(g.legal)), x, y+92, 14, color.RGBA{76, 213, 239, 255})
	}
	drawLabel(screen, g.message, int(v.x), int(v.y+v.size+29), 13, color.RGBA{220, 232, 242, 255})
	_ = h
}

func (g *Game) drawScoreMap(screen *ebiten.Image, v boardView, portrait bool) {
	size := float32(22)
	if portrait {
		size = maxf(11, minf(16, v.size/23))
	}
	x, y := v.x, v.y
	if portrait {
		x = (float32(g.viewW) - size*reversi.Size) / 2
		y = v.y + v.size + 34
	}
	for row := 0; row < reversi.Size; row++ {
		for col := 0; col < reversi.Size; col++ {
			value := reversi.ScoreMap[row][col]
			shade := uint8(62)
			if value >= 100 {
				shade = 240
			} else if value < 0 {
				shade = 90
			} else if value >= 15 {
				shade = 160
			}
			vector.DrawFilledRect(screen, x+float32(col)*size, y+float32(row)*size, size-1, size-1, color.RGBA{shade, 185, 126, 255}, false)
			drawLabel(screen, fmt.Sprint(value), int(x+float32(col)*size+2), int(y+float32(row)*size+size*.68), maxf(7, size*.45), color.RGBA{13, 28, 32, 255})
		}
	}
	if portrait {
		drawLabel(screen, g.tr("corner +120 / next to corner −40", "角 +120 / 角のとなり −40"), int(x), int(y+size*reversi.Size+18), 10, color.RGBA{184, 211, 233, 255})
	}
}

func (g *Game) drawCoordinates(screen *ebiten.Image, v boardView, cell float32) {
	for i := 0; i < reversi.Size; i++ {
		drawLabel(screen, fmt.Sprint(i+1), int(v.x+float32(i)*cell+cell*.45), int(v.y+v.size+14), maxf(9, cell*.22), color.RGBA{220, 232, 242, 255})
		drawLabel(screen, string(rune('A'+i)), int(v.x-16), int(v.y+float32(i)*cell+cell*.58), maxf(9, cell*.22), color.RGBA{220, 232, 242, 255})
	}
}

func (g *Game) drawResult(screen *ebiten.Image, v boardView) {
	w, h := v.size*.74, v.size*.25
	x, y := v.x+(v.size-w)/2, v.y+(v.size-h)/2
	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{7, 18, 31, 235}, false)
	vector.StrokeRect(screen, x, y, w, h, 2, color.RGBA{255, 211, 112, 255}, false)
	drawLabel(screen, g.result, int(x+18), int(y+h*.42), maxf(15, v.size*.045), color.RGBA{245, 250, 255, 255})
	drawLabel(screen, g.tr("REPLAY: Enter / R / tap", "もう一度：Enter / R / タップ"), int(x+18), int(y+h*.72), maxf(11, v.size*.03), color.RGBA{184, 211, 233, 255})
}

func drawButton(screen *ebiten.Image, x, y, w, h float32, label string, selected bool) {
	fill := color.RGBA{39, 65, 90, 255}
	if selected {
		fill = color.RGBA{45, 183, 157, 255}
	}
	vector.DrawFilledRect(screen, x, y, w, h, fill, false)
	vector.StrokeRect(screen, x, y, w, h, 1, color.RGBA{183, 230, 220, 255}, false)
	drawLabel(screen, label, int(x+8), int(y+h*.64), 11, color.RGBA{255, 255, 255, 255})
}

func drawLabel(screen *ebiten.Image, s string, x, y int, size float32, c color.Color) {
	text.Draw(screen, s, uiFace(float64(size)), x, y, c)
}

func (g *Game) stoneFrame(cell reversi.Move, stone reversi.Player, base float32) (color.RGBA, float32) {
	col := stoneColor(stone)
	for _, animation := range g.flips {
		if animation.cell != cell {
			continue
		}
		progress := 1 - float32(animation.remaining)/float32(flipFrames)
		if progress < .5 {
			col = stoneColor(animation.from)
		} else {
			col = stoneColor(animation.to)
		}
		return col, base * float32(.24+math.Abs(float64(progress-.5))*1.52)
	}
	return col, base
}

func stoneColor(player reversi.Player) color.RGBA {
	if player == reversi.Black {
		return color.RGBA{18, 24, 35, 255}
	}
	return color.RGBA{237, 243, 234, 255}
}

func (g *Game) view() boardView {
	w, h := g.dimensions()
	if w >= h {
		size := minf(float32(h-165), float32(w)*.60)
		return boardView{x: 34, y: 108, size: size}
	}
	size := minf(float32(w-80), float32(h-320))
	return boardView{x: (float32(w) - size) / 2, y: 128, size: size, portrait: true}
}

func (g *Game) dimensions() (int, int) {
	w, h := g.viewW, g.viewH
	if w <= 0 {
		w = 900
	}
	if h <= 0 {
		h = 720
	}
	return w, h
}

func (g *Game) cellAt(x, y int) (reversi.Move, bool) {
	v := g.view()
	cell := v.size / reversi.Size
	fx, fy := float32(x), float32(y)
	move := reversi.Move{X: int((fx - v.x) / cell), Y: int((fy - v.y) / cell)}
	return move, fx >= v.x && fy >= v.y && fx < v.x+v.size && fy < v.y+v.size && reversi.InBounds(move.X, move.Y)
}

func (g *Game) difficultyAt(x, y int) int {
	if g.variant != CPUEvaluation || y < 38 || y > 73 {
		return -1
	}
	w, _ := g.dimensions()
	buttonW := min(125, max(76, (w-58)/4))
	for i := 0; i < 3; i++ {
		left := 16 + i*(buttonW+6)
		if x >= left && x <= left+buttonW {
			return i
		}
	}
	return -1
}

func (g *Game) restartTouched() bool {
	x, y, ok := press()
	if !ok || y < 38 || y > 73 {
		return false
	}
	w, _ := g.dimensions()
	return x >= w-94 && x <= w-16
}

func keyboardDifficulty() (difficulty, bool) {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		return Friendly, true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		return Positional, true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		return Scout, true
	}
	return 0, false
}

func press() (int, int, bool) {
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

func anyPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func containsMove(moves []reversi.Move, target reversi.Move) bool {
	for _, move := range moves {
		if move == target {
			return true
		}
	}
	return false
}

func playerName(player reversi.Player) string {
	if player == reversi.Black {
		return "BLACK"
	}
	return "WHITE"
}
func playerNameJP(player reversi.Player) string {
	if player == reversi.Black {
		return "黒"
	}
	return "白"
}
func (g *Game) title() string {
	return fmt.Sprintf("REVERSI %02d / %s", g.variant+1, []string{"BOARD", "LEGAL", "FLIP", "PASS", "CPU"}[g.variant])
}
func maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
func minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	g.viewW, g.viewH = 480, 720
	return g.viewW, g.viewH
}

func Run(variant Variant) {
	ebiten.SetWindowTitle("Ebi Reversi — Ebitengine")
	if err := ebiten.RunGame(New(variant)); err != nil {
		panic(err)
	}
}
