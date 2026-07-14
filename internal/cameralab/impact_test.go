package cameralab

import "testing"

func TestImpactDeterministicAndRecovers(t *testing.T) {
	var i Impact
	i.Trigger(3, 2, 5)
	if !i.Frozen() {
		t.Fatal("stop")
	}
	a, b := i.Offset(7), i.Offset(7)
	if a != b {
		t.Fatal("nondeterministic")
	}
	i.Tick()
	i.Tick()
	if i.Frozen() {
		t.Fatal("stop remained")
	}
	i.Tick()
	i.Tick()
	i.Tick()
	if i.Offset(7) != (Vec{}) {
		t.Fatal("shake remained")
	}
}
