package main

import "testing"

func TestQuestTargetsShareOneWalkableRoute(t *testing.T) {
	type point struct{ x, y int }
	start := point{1, 10}
	want := []point{{2, 9}, {8, 2}, {8, 1}, {5, 6}}
	seen := map[point]bool{start: true}
	queue := []point{start}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, d := range []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := point{p.x + d.x, p.y + d.y}
			if passableTile(n.x, n.y) && !seen[n] {
				seen[n] = true
				queue = append(queue, n)
			}
		}
	}
	for _, target := range want {
		if !seen[target] {
			t.Errorf("quest target %+v is not reachable from %+v", target, start)
		}
	}
}

func TestBattleActionPulsePeaksAtImpact(t *testing.T) {
	if got := battleActionPulse(36); got != 0 {
		t.Fatalf("pulse at action start = %v, want 0", got)
	}
	if got := battleActionPulse(18); got != 1 {
		t.Fatalf("pulse at impact = %v, want 1", got)
	}
	if got := battleActionPulse(0); got != 0 {
		t.Fatalf("pulse at action end = %v, want 0", got)
	}
}

func TestPartyAttackDamage(t *testing.T) {
	if got := partyAttackDamage(0, 0, false); got != 10 {
		t.Fatalf("solo damage = %d, want 10", got)
	}
	if got := partyAttackDamage(0, 0, true); got != 15 {
		t.Fatalf("party damage = %d, want 15", got)
	}
	if got := partyAttackDamage(1, 0, true); got != 7 {
		t.Fatalf("shielded knight damage = %d, want 7", got)
	}
	if got := partyAttackDamage(1, 1, true); got != 15 {
		t.Fatalf("open knight damage = %d, want 15", got)
	}
}

func TestEnemyAttackIntentAndGuard(t *testing.T) {
	if got := enemyAttackDamage(2, 0, false); got != 15 {
		t.Fatalf("normal boss hit = %d, want 15", got)
	}
	if got := enemyAttackDamage(2, 1, false); got != 21 {
		t.Fatalf("heavy boss hit = %d, want 21", got)
	}
	if got := enemyAttackDamage(2, 1, true); got != 11 {
		t.Fatalf("guarded heavy boss hit = %d, want 11", got)
	}
	if got := enemyAttackDamage(0, 2, false); got != 4 {
		t.Fatalf("quick slime hit = %d, want 4", got)
	}
}
