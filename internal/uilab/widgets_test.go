package uilab

import "testing"

func TestAlignmentStillKeepsPanelCoordinates(t *testing.T) {
	if X(10, 4, Right) != 6 {
		t.Fatal("right edge")
	}
}
