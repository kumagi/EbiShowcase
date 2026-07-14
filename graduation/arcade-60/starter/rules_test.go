package starter

import "testing"

func TestAddStar(t *testing.T) {
	if AddStar(20) != 30 {
		t.Fatal("a star is worth 10")
	}
}
