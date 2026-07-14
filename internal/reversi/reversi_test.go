// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package reversi

import "testing"

func TestInitialBoardHasFourLegalMoves(t *testing.T) {
	board := NewBoard()
	moves := ValidMoves(board, Black)
	if len(moves) != 4 {
		t.Fatalf("initial legal moves = %d, want 4", len(moves))
	}
}

func TestApplyFlipsTheCapturedLine(t *testing.T) {
	board := NewBoard()
	flipped := Apply(&board, Black, Move{X: 2, Y: 3})
	if len(flipped) != 1 || board[3][3] != Black || board[3][2] != Black {
		t.Fatalf("move did not flip expected line: flipped=%v board=%v", flipped, board)
	}
}

func TestCornerWeightBeatsRiskySquare(t *testing.T) {
	var board Board
	board[0][0] = Black
	board[0][1] = White
	if Evaluate(board, Black) != 140 {
		t.Fatalf("evaluation = %d, want 140", Evaluate(board, Black))
	}
}

func TestChooseBestPrefersAHighValueCorner(t *testing.T) {
	var board Board
	board[0][1] = White
	board[0][2] = Black
	move, ok, _ := ChooseBest(board, Black)
	if !ok || move != (Move{X: 0, Y: 0}) {
		t.Fatalf("best move = %#v, ok=%v; want corner", move, ok)
	}
}

func TestChooseLookaheadIsLegalAndDeterministic(t *testing.T) {
	board := NewBoard()
	first, ok, score := ChooseLookahead(board, Black)
	if !ok || len(Captures(board, Black, first)) == 0 {
		t.Fatalf("lookahead returned illegal move %#v, ok=%v", first, ok)
	}
	second, ok2, score2 := ChooseLookahead(board, Black)
	if !ok2 || second != first || score2 != score {
		t.Fatalf("lookahead must be deterministic: first=%#v/%d second=%#v/%d", first, score, second, score2)
	}
}
