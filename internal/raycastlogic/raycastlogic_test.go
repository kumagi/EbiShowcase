// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package raycastlogic

import (
	"math"
	"testing"
)

func TestCastFindsFirstWall(t *testing.T) {
	grid := [][]int{{1, 1, 1, 1}, {1, 0, 0, 1}, {1, 0, 0, 1}, {1, 1, 1, 1}}
	hit := Cast(grid, 1.5, 1.5, 0)
	if hit.MapX != 3 || hit.MapY != 1 || math.Abs(hit.Distance-1.5) > 1e-9 {
		t.Fatalf("Cast() = %+v, want wall (3,1) at 1.5", hit)
	}
}

func TestCorrectDistanceRemovesFishEye(t *testing.T) {
	got := CorrectDistance(2, math.Pi/3, math.Pi/2)
	want := 2 * math.Cos(-math.Pi/6)
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("CorrectDistance() = %v, want %v", got, want)
	}
}

func TestProjectSpriteCenterAndBehind(t *testing.T) {
	center := ProjectSprite(1, 1, 0, 4, 1, math.Pi/3)
	if math.Abs(center.ScreenX) > 1e-9 || center.Depth != 3 {
		t.Fatalf("center projection = %+v", center)
	}
	behind := ProjectSprite(1, 1, 0, 0, 1, math.Pi/3)
	if behind.Depth >= 0 {
		t.Fatalf("behind projection = %+v", behind)
	}
}
