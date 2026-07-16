// Copyright 2026 Ebi Showcase contributors
// SPDX-License-Identifier: Apache-2.0

package main

import "testing"

func TestAutoplayPushesBoxOntoGoal(t *testing.T) {
	g := newGame()
	for i := 0; i < 800 && g.clears == 0; i++ {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
	}
	if g.clears == 0 {
		t.Fatal("autoplay never completed the warehouse puzzle")
	}
}

func TestBlockedCellsIncludeOuterWall(t *testing.T) {
	for _, c := range []cell{{0, 3}, {6, 3}, {3, 0}, {3, 6}, {2, 2}} {
		if !blocked(c) {
			t.Fatalf("expected %#v to be blocked", c)
		}
	}
}

func TestCompletedMoveDoesNotDrawPreviousCellDuringPause(t *testing.T) {
	g := newGame()
	g.pause = 0
	g.beginMove(move{1, 0})
	for g.moveFrame > 0 {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
	}
	px, py, _, _ := g.positions()
	if px != float64(g.player.x) || py != float64(g.player.y) {
		t.Fatalf("draw position after move = (%.1f, %.1f), state = %#v", px, py, g.player)
	}
}
