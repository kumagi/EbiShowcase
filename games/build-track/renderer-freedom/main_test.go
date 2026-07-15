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
