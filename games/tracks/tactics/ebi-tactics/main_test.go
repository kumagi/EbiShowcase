package main

import "testing"

func TestNewGameCanBeCreatedAgainForRetry(t *testing.T) {
	first := newGame()
	second := newGame()
	if first.audio == nil || first.audio != second.audio {
		t.Fatal("retry must reuse the single Ebitengine audio context")
	}
}

func TestWaitEndsTheSelectedUnitsAction(t *testing.T) {
	g := newGame()
	g.selected = 0
	g.units[0].moved = true
	g.waitUnit()
	if !g.units[0].acted {
		t.Fatal("END ACTION must mark the selected unit as acted")
	}
}
