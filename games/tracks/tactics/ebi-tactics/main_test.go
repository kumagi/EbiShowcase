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
		t.Fatal("FINISH UNIT must mark the selected unit as acted")
	}
}

func TestAllyMovementAnimatesBeforeAcceptingAnotherAction(t *testing.T) {
	g := newGame()
	destination := pt{0, 6}
	g.choose(destination)

	if !g.units[0].moved || g.units[0].p != destination {
		t.Fatal("choosing a reachable tile must commit the ally move")
	}
	if len(g.moves) != 1 {
		t.Fatalf("ally move must start one animation, got %d", len(g.moves))
	}

	startX, startY := g.unitScreenPosition(0, g.units[0])
	wantX := float32(ox + tile/2)
	wantY := float32(oy + 7*tile + tile/2)
	if startX != wantX || startY != wantY {
		t.Fatalf("animation must begin at the old tile: got (%v,%v), want (%v,%v)", startX, startY, wantX, wantY)
	}
	g.moves[0].tick = g.moves[0].duration / 2
	_, middleY := g.unitScreenPosition(0, g.units[0])
	endY := float32(oy + 6*tile + tile/2)
	if middleY >= startY || middleY <= endY-8 {
		t.Fatalf("animation must visibly travel between tiles with a small hop, got y=%v between %v and %v", middleY, startY, endY)
	}

	for len(g.moves) > 0 {
		g.advanceMoves()
	}
	if g.message != "MOVE COMPLETE. Attack an enemy or FINISH THIS UNIT." {
		t.Fatalf("move completion must explain the next unit-level choice, got %q", g.message)
	}
}

func TestEnemyDamageWaitsForMovementAnimation(t *testing.T) {
	g := newGame()
	g.units[2].p = pt{0, 6}
	g.units[0].acted = true
	g.units[1].acted = true
	before := g.units[0].hp

	g.enemyTurnIfDone()
	if !g.enemyPhase {
		t.Fatal("finishing both allies must begin the enemy phase")
	}
	if g.units[0].hp != before {
		t.Fatal("enemy damage must wait until enemy movement is visible")
	}
	if len(g.moves) == 0 {
		t.Fatal("enemy phase must animate enemies that change tiles")
	}

	for g.enemyPhase {
		g.advanceMoves()
	}
	if g.units[0].hp != before-2 {
		t.Fatalf("enemy strike must resolve after movement, got HP %d, want %d", g.units[0].hp, before-2)
	}
}
