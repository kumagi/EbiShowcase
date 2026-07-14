package uilab

import "testing"

func TestStatusAndScrollClamp(t *testing.T) {
	var s Status
	s.Set("Saved", 2)
	s.Tick()
	if s.Text == "" {
		t.Fatal("early clear")
	}
	s.Tick()
	if s.Text != "" {
		t.Fatal("not clear")
	}
	sc := Scroll{Max: 3}
	sc.Move(9)
	if sc.Offset != 3 {
		t.Fatal(sc.Offset)
	}
	sc.Move(-9)
	if sc.Offset != 0 {
		t.Fatal(sc.Offset)
	}
}
