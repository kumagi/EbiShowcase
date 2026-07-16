package main

import "testing"

func TestBlastStopsAtHardAndBreakableWalls(t *testing.T) {
	g := &game{power: 4, soft: map[point]bool{{3, 1}: true}}
	blast := g.blast(point{1, 1})

	for _, want := range []point{{1, 1}, {2, 1}, {3, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}} {
		if !pointIn(blast, want) {
			t.Fatalf("blast does not contain %v: %v", want, blast)
		}
	}
	if pointIn(blast, point{4, 1}) {
		t.Fatal("blast passed through a breakable wall")
	}
	if pointIn(blast, point{2, 2}) {
		t.Fatal("blast passed through a hard pillar")
	}
}

func TestTriggerBombsCreatesChainReaction(t *testing.T) {
	bombs := []bomb{{at: point{6, 7}, timer: 30}, {at: point{6, 5}, timer: 45}}
	triggered := triggerBombs([]point{{4, 7}, {5, 7}, {6, 7}}, bombs)
	if triggered != 1 {
		t.Fatalf("triggered = %d, want 1", triggered)
	}
	if bombs[0].timer != 0 {
		t.Fatalf("reached bomb timer = %d, want 0", bombs[0].timer)
	}
	if bombs[1].timer != 45 {
		t.Fatalf("unreached bomb timer = %d, want 45", bombs[1].timer)
	}
}
