package main

import "testing"

func TestActionAddsOneStar(t *testing.T) {
	round := NewRound()
	round.Step(Input{Action: true})
	if round.Score != 10 {
		t.Fatalf("score = %d, want 10", round.Score)
	}
}

func TestTimeLimitEndsRound(t *testing.T) {
	round := NewRound()
	round.Frames = RoundFrames - 1
	round.Step(Input{})
	if !round.Over {
		t.Fatal("round must end")
	}
}

func TestRestartMakesFreshRound(t *testing.T) {
	round := Round{Score: 40, Frames: 900, Over: true}
	round.Step(Input{Restart: true})
	if round != (Round{}) {
		t.Fatalf("restart = %+v", round)
	}
}
