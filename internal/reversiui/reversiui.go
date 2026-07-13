// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package reversiui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/reversi"
)

const (
	ScreenWidth  = 900
	ScreenHeight = 720
	BoardX       = 40
	BoardY       = 100
	CellSize     = 70
)

type Variant int

const (
	BoardGrid Variant = iota
	LegalMoves
	FlipStones
	PassAndScore
	CPUEvaluation
)

type Game struct {
	variant   Variant
	board     reversi.Board
	turn      reversi.Player
	legal     []reversi.Move
	lastMove  reversi.Move
	lastFlips []reversi.Move
	passCount int
	cpuWait   int
	message   string
	over      bool
	result    string
}

func New(variant Variant) *Game {
	g := &Game{variant: variant}
	g.reset()
	return g
}

func (g *Game) reset() {
	g.board = reversi.NewBoard()
	g.turn = reversi.Black
	g.legal = reversi.ValidMoves(g.board, g.turn)
	g.lastMove = reversi.Move{X: -1, Y: -1}
	g.lastFlips = nil
	g.passCount = 0
	g.cpuWait = 18
	g.over = false
	g.result = ""
	g.message = g.initialMessage()
}

func (g *Game) initialMessage() string {
	switch g.variant {
	case BoardGrid:
		return "Click a cell: x and y become board data. SPACE resets."
	case LegalMoves:
		return "Blue circles are legal moves. Click one to place BLACK."
	case FlipStones:
		return "Take turns locally. A bracketed line is flipped."
	case PassAndScore:
		return "You are BLACK. The simple CPU chooses its first legal move."
	default:
		return "You are BLACK. The CPU chooses the highest weighted position."
	}
}

func (g *Game) Update() error {
	if g.over {
		if pressedAny() {
			g.reset()
		}
		return nil
	}

	if g.variant >= PassAndScore && g.turn == reversi.White {
		if g.cpuWait > 0 {
			g.cpuWait--
			return nil
		}
		var move reversi.Move
		var ok bool
		if g.variant == CPUEvaluation {
			move, ok, _ = reversi.ChooseBest(g.board, reversi.White)
		} else {
			move, ok = reversi.ChooseFirst(g.board, reversi.White)
		}
		if !ok {
			g.passOrFinish()
			return nil
		}
		g.play(move)
		return nil
	}

	if g.variant == BoardGrid {
		if x, y, ok := press(); ok {
			if cell, inside := boardCell(x, y); inside && g.board[cell.Y][cell.X] == reversi.Empty {
				g.board[cell.Y][cell.X] = reversi.Black
				g.lastMove = cell
				g.message = fmt.Sprintf("board[%d][%d] = BLACK", cell.Y, cell.X)
			}
		}
		return nil
	}

	if x, y, ok := press(); ok {
		if cell, inside := boardCell(x, y); inside && containsMove(g.legal, cell) {
			g.play(cell)
		}
	}
	return nil
}

func (g *Game) play(move reversi.Move) {
	g.lastMove = move
	g.lastFlips = reversi.Apply(&g.board, g.turn, move)
	g.passCount = 0
	g.turn = reversi.Opponent(g.turn)
	g.cpuWait = 18
	g.refreshTurn()
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
	if g.variant == LegalMoves {
		g.message = fmt.Sprintf("%s turn: click a blue legal move (%d available)", playerName(g.turn), len(g.legal))
	} else if g.variant == FlipStones {
		g.message = fmt.Sprintf("%s turn: choose a move, then watch %d stones flip", playerName(g.turn), len(g.legal))
	} else if g.turn == reversi.Black {
		g.message = fmt.Sprintf("Your turn: %d legal moves", len(g.legal))
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
	g.cpuWait = 18
	g.message = fmt.Sprintf("%s has no move: PASS → %s", playerName(passed), playerName(g.turn))
}

func (g *Game) finish() {
	g.over = true
	black, white := reversi.Count(g.board, reversi.Black), reversi.Count(g.board, reversi.White)
	switch {
	case black > white:
		g.result = fmt.Sprintf("BLACK wins %d - %d", black, white)
	case white > black:
		g.result = fmt.Sprintf("WHITE wins %d - %d", white, black)
	default:
		g.result = fmt.Sprintf("DRAW %d - %d", black, white)
	}
	g.message = "GAME OVER — SPACE or TAP to restart"
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 20, 37, 255})
	ebitenutil.DebugPrintAt(screen, g.title(), 40, 25)
	ebitenutil.DebugPrintAt(screen, g.message, 40, 52)

	vector.DrawFilledRect(screen, BoardX-8, BoardY-8, CellSize*8+16, CellSize*8+16, color.RGBA{28, 93, 74, 255}, false)
	for y := 0; y < reversi.Size; y++ {
		for x := 0; x < reversi.Size; x++ {
			cx, cy := float32(BoardX+x*CellSize), float32(BoardY+y*CellSize)
			vector.DrawFilledRect(screen, cx, cy, CellSize, CellSize, color.RGBA{28, 117, 82, 255}, false)
			vector.StrokeRect(screen, cx, cy, CellSize, CellSize, 1, color.RGBA{142, 211, 154, 170}, false)
			cell := g.board[y][x]
			if containsMove(g.legal, reversi.Move{X: x, Y: y}) && !g.over && g.variant != BoardGrid {
				vector.DrawFilledCircle(screen, cx+CellSize/2, cy+CellSize/2, 10, color.RGBA{76, 213, 239, 150}, false)
			}
			if cell != reversi.Empty {
				stoneColor := color.RGBA{237, 243, 234, 255}
				if cell == reversi.Black {
					stoneColor = color.RGBA{18, 24, 35, 255}
				}
				vector.DrawFilledCircle(screen, cx+CellSize/2, cy+CellSize/2, 27, stoneColor, false)
				vector.StrokeCircle(screen, cx+CellSize/2, cy+CellSize/2, 27, 2, color.RGBA{255, 255, 255, 100}, false)
			}
		}
	}
	for i := 0; i < reversi.Size; i++ {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i+1), BoardX+CellSize*i+31, BoardY+CellSize*8+6)
		ebitenutil.DebugPrintAt(screen, string(rune('A'+i)), BoardX-22, BoardY+CellSize*i+31)
	}

	black, white := reversi.Count(g.board, reversi.Black), reversi.Count(g.board, reversi.White)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BLACK %02d     WHITE %02d", black, white), 635, 105)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TURN  %s", playerName(g.turn)), 635, 130)
	if g.variant == CPUEvaluation {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MAP EVAL (BLACK) %d", reversi.Evaluate(g.board, reversi.Black)), 635, 165)
		drawScoreMap(screen)
		ebitenutil.DebugPrintAt(screen, "corner = +120", 635, 385)
		ebitenutil.DebugPrintAt(screen, "near corner = -40", 635, 405)
	} else if g.variant == PassAndScore {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("legal moves %d", len(g.legal)), 635, 165)
		ebitenutil.DebugPrintAt(screen, "CPU: first legal move", 635, 190)
	} else if g.variant == BoardGrid {
		ebitenutil.DebugPrintAt(screen, "x = column, y = row", 635, 165)
		ebitenutil.DebugPrintAt(screen, "Click to write one stone", 635, 190)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("legal moves %d", len(g.legal)), 635, 165)
		ebitenutil.DebugPrintAt(screen, "blue dot = legal", 635, 190)
	}
	if g.over {
		vector.DrawFilledRect(screen, BoardX+80, BoardY+240, 400, 120, color.RGBA{7, 18, 31, 235}, false)
		ebitenutil.DebugPrintAt(screen, g.result, BoardX+135, BoardY+280)
		ebitenutil.DebugPrintAt(screen, g.message, BoardX+105, BoardY+315)
	}
}

func (g *Game) title() string {
	switch g.variant {
	case BoardGrid:
		return "REVERSI 01 / BOARD GRID"
	case LegalMoves:
		return "REVERSI 02 / LEGAL MOVES"
	case FlipStones:
		return "REVERSI 03 / FLIP THE LINE"
	case PassAndScore:
		return "REVERSI 04 / PASS & SCORE"
	default:
		return "REVERSI 05 / CPU EVALUATION"
	}
}

func drawScoreMap(screen *ebiten.Image) {
	const mapX, mapY, size = 635, 225, 24
	for y := 0; y < reversi.Size; y++ {
		for x := 0; x < reversi.Size; x++ {
			value := reversi.ScoreMap[y][x]
			shade := uint8(62)
			if value >= 100 {
				shade = 240
			} else if value < 0 {
				shade = 90
			} else if value >= 15 {
				shade = 160
			}
			vector.DrawFilledRect(screen, mapX+float32(x*size), mapY+float32(y*size), size-2, size-2, color.RGBA{shade, 185, 126, 255}, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", value), mapX+x*size+3, mapY+y*size+7)
		}
	}
}

func playerName(player reversi.Player) string {
	if player == reversi.Black {
		return "BLACK"
	}
	return "WHITE"
}

func containsMove(moves []reversi.Move, target reversi.Move) bool {
	for _, move := range moves {
		if move == target {
			return true
		}
	}
	return false
}

func boardCell(x, y int) (reversi.Move, bool) {
	cell := reversi.Move{X: int(math.Floor(float64(x-BoardX) / CellSize)), Y: int(math.Floor(float64(y-BoardY) / CellSize))}
	return cell, x >= BoardX && y >= BoardY && cell.X < reversi.Size && cell.Y < reversi.Size
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

func pressedAny() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *Game) Layout(_, _ int) (int, int) { return ScreenWidth, ScreenHeight }

func Run(variant Variant) {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Ebi Reversi — Ebitengine")
	if err := ebiten.RunGame(New(variant)); err != nil {
		panic(err)
	}
}
