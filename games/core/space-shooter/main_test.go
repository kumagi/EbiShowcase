package main

import (
	"math"
	"testing"
)

func TestAimVelocityNormalizesDirection(t *testing.T) {
	vx, vy := aimVelocity(3, 4, 2.5)
	if got := math.Hypot(vx, vy); math.Abs(got-2.5) > 1e-9 {
		t.Fatalf("length = %v, want 2.5", got)
	}
	if vx != 1.5 || vy != 2 {
		t.Fatalf("vector = (%v, %v), want (1.5, 2)", vx, vy)
	}
}

func TestAimVelocityZeroDirectionIsSafe(t *testing.T) {
	vx, vy := aimVelocity(0, 0, 2.5)
	if vx != 0 || vy != 0 {
		t.Fatalf("zero direction = (%v, %v), want (0, 0)", vx, vy)
	}
}
