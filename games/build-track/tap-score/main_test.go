package main

import "testing"

func TestAddScore(t *testing.T) {
	if got := addScore(2); got != 3 {
		t.Fatalf("addScore(2) = %d, want 3", got)
	}
}
