package main

import "testing"

func TestNewGameHasVisibleStartingPosition(t *testing.T) {
	g := newGame()
	if g.x <= 0 || g.y <= 0 {
		t.Fatalf("newGame position = (%d, %d), want positive coordinates", g.x, g.y)
	}
}
