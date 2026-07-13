// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
// Package reversi contains the rules and evaluation helpers used by the
// playable Reversi lessons.
package reversi

const Size = 8

type Player int

const (
	Empty Player = 0
	Black Player = 1
	White Player = -1
)

type Move struct {
	X, Y int
}

type Board [Size][Size]Player

// ScoreMap rewards corners and stable edges while discouraging risky squares
// next to an empty corner.  It is deliberately visible to learners.
var ScoreMap = [Size][Size]int{
	{120, -20, 20, 5, 5, 20, -20, 120},
	{-20, -40, -5, -5, -5, -5, -40, -20},
	{20, -5, 15, 3, 3, 15, -5, 20},
	{5, -5, 3, 3, 3, 3, -5, 5},
	{5, -5, 3, 3, 3, 3, -5, 5},
	{20, -5, 15, 3, 3, 15, -5, 20},
	{-20, -40, -5, -5, -5, -5, -40, -20},
	{120, -20, 20, 5, 5, 20, -20, 120},
}

var directions = [...]Move{
	{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1},
	{X: -1, Y: 0}, {X: 1, Y: 0},
	{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1},
}

func NewBoard() Board {
	var board Board
	board[3][3], board[4][4] = White, White
	board[3][4], board[4][3] = Black, Black
	return board
}

func Opponent(player Player) Player {
	if player == Black {
		return White
	}
	return Black
}

func InBounds(x, y int) bool {
	return x >= 0 && x < Size && y >= 0 && y < Size
}

func Captures(board Board, player Player, move Move) []Move {
	if player == Empty || !InBounds(move.X, move.Y) || board[move.Y][move.X] != Empty {
		return nil
	}
	opponent := Opponent(player)
	var captured []Move
	for _, direction := range directions {
		line := make([]Move, 0, Size)
		x, y := move.X+direction.X, move.Y+direction.Y
		for InBounds(x, y) && board[y][x] == opponent {
			line = append(line, Move{X: x, Y: y})
			x += direction.X
			y += direction.Y
		}
		if len(line) > 0 && InBounds(x, y) && board[y][x] == player {
			captured = append(captured, line...)
		}
	}
	return captured
}

func ValidMoves(board Board, player Player) []Move {
	moves := make([]Move, 0, 20)
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			move := Move{X: x, Y: y}
			if len(Captures(board, player, move)) > 0 {
				moves = append(moves, move)
			}
		}
	}
	return moves
}

func Apply(board *Board, player Player, move Move) []Move {
	if board == nil {
		return nil
	}
	captured := Captures(*board, player, move)
	if len(captured) == 0 {
		return nil
	}
	board[move.Y][move.X] = player
	for _, cell := range captured {
		board[cell.Y][cell.X] = player
	}
	return captured
}

func Count(board Board, player Player) int {
	count := 0
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			if board[y][x] == player {
				count++
			}
		}
	}
	return count
}

func Full(board Board) bool {
	return Count(board, Empty) == 0
}

// Evaluate sums position weights from the requested player's point of view.
// With black=+1 and white=-1 this is the same idea as stone*ScoreMap summed
// over all 64 cells, then oriented toward player.
func Evaluate(board Board, player Player) int {
	if player == Empty {
		return 0
	}
	value := 0
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			if board[y][x] == player {
				value += ScoreMap[y][x]
			} else if board[y][x] == Opponent(player) {
				value -= ScoreMap[y][x]
			}
		}
	}
	return value
}

// ChooseFirst is the deliberately simple CPU used before the evaluation
// lesson: it makes the first legal move in stable row-major order.
func ChooseFirst(board Board, player Player) (Move, bool) {
	moves := ValidMoves(board, player)
	if len(moves) == 0 {
		return Move{}, false
	}
	return moves[0], true
}

// ChooseBest applies every legal move on a copy and keeps the highest
// position-map score. Ties stay deterministic, making the lesson testable.
func ChooseBest(board Board, player Player) (Move, bool, int) {
	moves := ValidMoves(board, player)
	if len(moves) == 0 {
		return Move{}, false, 0
	}
	best := moves[0]
	bestScore := -1 << 30
	for _, move := range moves {
		candidate := board
		Apply(&candidate, player, move)
		score := Evaluate(candidate, player)
		if score > bestScore {
			best, bestScore = move, score
		}
	}
	return best, true, bestScore
}
