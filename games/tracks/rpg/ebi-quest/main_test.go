package main

import "testing"

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
