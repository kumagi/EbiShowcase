package uilab

import "testing"

func TestFocusSkipsDisabled(t *testing.T) {
	f := Focus{Count: 3, Disabled: map[int]bool{1: true}}
	f.Move(1)
	if f.Index != 2 {
		t.Fatal(f.Index)
	}
	f.Move(1)
	if f.Index != 0 || !f.Activate() {
		t.Fatal("focus")
	}
}
