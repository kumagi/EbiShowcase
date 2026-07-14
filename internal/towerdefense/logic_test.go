package towerdefense

import "testing"

func TestPathPositionAcrossCorners(t *testing.T) {
	p := NewPath([]Vec{{0, 0}, {10, 0}, {10, 10}})
	for _, tc := range []struct {
		d    float64
		want Vec
	}{{0, Vec{0, 0}}, {5, Vec{5, 0}}, {15, Vec{10, 5}}, {99, Vec{10, 10}}} {
		got := p.Position(tc.d)
		if got != tc.want {
			t.Fatalf("Position(%v)=%v, want %v", tc.d, got, tc.want)
		}
	}
}

func TestSelectFrontIgnoresBehindOutOfRangeAndDead(t *testing.T) {
	targets := []Target{{Vec{3, 0}, 20, true}, {Vec{4, 0}, 80, true}, {Vec{2, 0}, 100, false}, {Vec{8, 0}, 120, true}}
	if got := SelectFront(Vec{}, 5, targets); got != 1 {
		t.Fatalf("SelectFront=%d, want 1", got)
	}
}

func TestSelectFrontReturnsMinusOne(t *testing.T) {
	if got := SelectFront(Vec{}, 2, []Target{{Vec{4, 0}, 9, true}}); got != -1 {
		t.Fatalf("got %d", got)
	}
}
