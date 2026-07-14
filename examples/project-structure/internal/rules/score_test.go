package rules

import "testing"

func TestCollect(t *testing.T) {
	if score, gems := Collect(20, 2); score != 30 || gems != 3 {
		t.Fatalf("got %d, %d", score, gems)
	}
}
