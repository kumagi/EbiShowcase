package rhythmcore

import "testing"

func TestGradeWindows(t *testing.T) {
	for _, tc := range []struct {
		d    int
		want Grade
	}{{0, Perfect}, {4, Perfect}, {5, Good}, {9, Good}, {10, Miss}, {-3, Perfect}} {
		if got := GradeDelta(tc.d); got != tc.want {
			t.Fatalf("GradeDelta(%d)=%s want %s", tc.d, got, tc.want)
		}
	}
}

func TestTapComboAndMiss(t *testing.T) {
	s := NewSession(Chart{Lanes: 1, Notes: []Note{{Lane: 0, At: 2}, {Lane: 0, At: 8}}})
	for s.Frame < 2 {
		s.Step(nil)
	}
	r := s.Step([]Input{{Lane: 0, Down: true}})
	if len(r) != 1 || r[0].Grade != Perfect || s.Combo != 1 {
		t.Fatalf("first tap: %#v combo=%d", r, s.Combo)
	}
	for s.Frame <= 18 {
		s.Step(nil)
	}
	if s.Misses != 1 || s.Combo != 0 {
		t.Fatalf("miss=%d combo=%d", s.Misses, s.Combo)
	}
}

func TestHoldRequiresReleaseAfterEnd(t *testing.T) {
	s := NewSession(Chart{Lanes: 1, Notes: []Note{{Lane: 0, At: 2, Kind: Hold, Duration: 5}}})
	for s.Frame < 2 {
		s.Step(nil)
	}
	s.Step([]Input{{Lane: 0, Down: true}})
	for s.Frame <= 7 {
		s.Step(nil)
	}
	if s.Perfects != 1 {
		t.Fatalf("hold perfects=%d", s.Perfects)
	}

	s = NewSession(Chart{Lanes: 1, Notes: []Note{{Lane: 0, At: 2, Kind: Hold, Duration: 5}}})
	for s.Frame < 2 {
		s.Step(nil)
	}
	s.Step([]Input{{Lane: 0, Down: true}})
	s.Step([]Input{{Lane: 0, Down: false}})
	if s.Misses != 1 {
		t.Fatalf("early release misses=%d", s.Misses)
	}
}

func TestRollCountsRepeatedPresses(t *testing.T) {
	s := NewSession(Chart{Lanes: 1, Notes: []Note{{Lane: 0, At: 2, Kind: Roll, Duration: 6, Need: 4}}})
	for s.Frame < 2 {
		s.Step(nil)
	}
	for i := 0; i < 4; i++ {
		s.Step([]Input{{Lane: 0, Down: true}, {Lane: 0, Down: false}})
	}
	for s.Frame <= 9 {
		s.Step(nil)
	}
	if s.Perfects != 1 {
		t.Fatalf("roll perfects=%d hits=%d", s.Perfects, s.RollHits(0))
	}
}
