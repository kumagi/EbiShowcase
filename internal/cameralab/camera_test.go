package cameralab

import "testing"

func TestCoordinateRoundTrip(t *testing.T) {
	s := State{Pos: Vec{50, 70}, ViewW: 480, ViewH: 720}
	p := Vec{121, 214}
	if got := s.ScreenToWorld(s.WorldToScreen(p)); got != p {
		t.Fatalf("%+v", got)
	}
}
