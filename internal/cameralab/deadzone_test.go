package cameralab

import "testing"

func TestDeadZoneAndLookAhead(t *testing.T) {
	s := State{Pos: Vec{100, 100}}
	s.FollowDeadZone(Vec{110, 100}, Vec{0, 0}, 20, 20, .5)
	if s.Pos.X != 100 {
		t.Fatal("moved inside zone")
	}
	s.FollowDeadZone(Vec{110, 100}, Vec{40, 0}, 20, 20, .5)
	if s.Pos.X <= 100 {
		t.Fatal("lookahead")
	}
}
