package main

import "testing"

func patternedBoard() [rows][cols]int {
	var board [rows][cols]int
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			board[y][x] = (x + y*2) % len(pieceColors)
		}
	}
	return board
}

func TestSpecialForMatch(t *testing.T) {
	tests := []struct {
		name string
		want int
		make func(*[rows][cols]int) point
	}{
		{"four horizontal", specialRow, func(b *[rows][cols]int) point {
			for x := 1; x <= 4; x++ {
				b[3][x] = 0
			}
			return point{2, 3}
		}},
		{"four vertical", specialCol, func(b *[rows][cols]int) point {
			for y := 1; y <= 4; y++ {
				b[y][3] = 1
			}
			return point{3, 3}
		}},
		{"five", specialColor, func(b *[rows][cols]int) point {
			for x := 0; x <= 4; x++ {
				b[3][x] = 2
			}
			return point{2, 3}
		}},
		{"cross", specialArea, func(b *[rows][cols]int) point {
			for x := 1; x <= 3; x++ {
				b[3][x] = 3
			}
			for y := 2; y <= 4; y++ {
				b[y][2] = 3
			}
			return point{2, 3}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := patternedBoard()
			anchor := tt.make(&board)
			if got := specialForMatch(board, anchor); got != tt.want {
				t.Fatalf("specialForMatch() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSpecialsChainIntoEachOther(t *testing.T) {
	g := &game{board: patternedBoard()}
	g.specials[3][1] = specialRow
	g.specials[3][4] = specialCol
	got := g.expandSpecials(map[point]bool{{1, 3}: true})
	if len(got) != cols+rows-1 {
		t.Fatalf("chain cleared %d cells, want %d", len(got), cols+rows-1)
	}
}

func TestFirstStageOffersASpecialSwap(t *testing.T) {
	board := firstStage.board
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			for _, d := range [...]point{{1, 0}, {0, 1}} {
				n := point{x + d.x, y + d.y}
				if n.x >= cols || n.y >= rows {
					continue
				}
				p := point{x, y}
				board[p.y][p.x], board[n.y][n.x] = board[n.y][n.x], board[p.y][p.x]
				matches := scan(board)
				if matches[n] && specialForMatch(board, n) != specialNone || matches[p] && specialForMatch(board, p) != specialNone {
					return
				}
				board = firstStage.board
			}
		}
	}
	t.Fatal("first stage has no swap that demonstrates a special piece")
}
