package main

import "testing"

func TestFreshProgressStartsAtFirstDataStage(t *testing.T) {
	p := NewProgress()
	if p.Stage != 0 || p.Clear {
		t.Fatalf("fresh = %+v", p)
	}
}
func TestPuzzleTargetAdvancesOneStage(t *testing.T) {
	p := NewProgress()
	for range Stages[0].Target {
		p.Step(Input{Solve: true})
	}
	if p.Stage != 1 || p.Presses != 0 {
		t.Fatalf("after first puzzle = %+v", p)
	}
}
func TestLastDataStageClearsInsteadOfOverflowing(t *testing.T) {
	p := Progress{Stage: len(Stages) - 1}
	for range Stages[len(Stages)-1].Target {
		p.Step(Input{Solve: true})
	}
	if !p.Clear || p.Stage != len(Stages)-1 {
		t.Fatalf("last = %+v", p)
	}
}
func TestRestartMakesFreshProgress(t *testing.T) {
	p := Progress{Stage: 2, Presses: 2, Clear: true}
	p.Step(Input{Restart: true})
	if p != (Progress{}) {
		t.Fatal(p)
	}
}
