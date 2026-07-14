package cameralab

import "testing"

func TestFollowSmoothsAndClamps(t *testing.T) {
	s := State{ViewW: 100, ViewH: 100}
	s.Follow(Vec{500, 500}, .5, 300, 300)
	if s.Pos.X != 250 || s.Pos.Y != 250 {
		t.Fatalf("%+v", s.Pos)
	}
}
