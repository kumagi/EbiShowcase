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

func TestMissionValidationAndGrade(t *testing.T) {
	m := Mission{Name: "test", Grid: [][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}}, StartX: 1.5, StartY: 1.5, KeyX: 1.5, KeyY: 1.5, ExitX: 1.5, ExitY: 1.5, GoalTime: 60}
	if !ValidateMission(m) {
		t.Fatal("valid mission was rejected")
	}
	m.Grid[1] = []int{1, 1}
	if ValidateMission(m) {
		t.Fatal("ragged mission was accepted")
	}
	if got := Grade(40, 0, 3, 2); got != "S" {
		t.Fatalf("Grade() = %q, want S", got)
	}
	if got := Grade(90, 4, 10, 2); got != "C" {
		t.Fatalf("Grade() = %q, want C", got)
	}
}
