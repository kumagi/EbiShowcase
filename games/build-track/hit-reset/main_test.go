package main

import "testing"

func TestHit(t *testing.T) {
	if !hit(10, 10, 10, 10, 5) {
		t.Fatal("centre should hit")
	}
	if hit(20, 10, 10, 10, 5) {
		t.Fatal("far point should miss")
	}
}

func TestMoveTargetAdvancesDeterministically(t *testing.T) {
	g := newGame()
	g.moveTarget()
	if g.targetX != targets[1][0] || g.targetY != targets[1][1] {
		t.Fatalf("got (%v,%v), want next target", g.targetX, g.targetY)
	}
}
