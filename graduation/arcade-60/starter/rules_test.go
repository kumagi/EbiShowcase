package main

import "testing"

// TODO 4 in rules.go: a single action is worth one star.
func TestActionAddsOneStar(t *testing.T) {
	round := NewRound()
	round.Step(Input{Action: true})
	if round.Score != 10 {
		t.Fatalf("score after one action = %d, want 10", round.Score)
	}
}

// TODO 5 in rules.go: exactly 60 seconds ends the round.
func TestTimeLimitEndsRound(t *testing.T) {
	round := NewRound()
	round.Ticks = RoundTicks - 1
	round.Step(Input{})
	if !round.Over {
		t.Fatal("the round must end at RoundTicks")
	}
}

// TODO 3 in rules.go: a restart removes the old score and time.
func TestRestartMakesFreshRound(t *testing.T) {
	round := Round{Score: 40, Ticks: 900, Over: true}
	round.Step(Input{Restart: true})
	if round.Score != 0 || round.Ticks != 0 || round.Over {
		t.Fatalf("restart = %+v, want a fresh round", round)
	}
}
