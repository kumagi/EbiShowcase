package main

import "testing"

func TestFreshProgressStartsAtFirstDataStage(t *testing.T) {
	p := NewProgress()
	if p.Stage != 0 || p.Clear {
		t.Fatal(p)
	}
}
func TestPuzzleTargetAdvancesOneStage(t *testing.T) {
	p := NewProgress()
	for range Stages[0].Target {
		p.Step(Input{Solve: true})
	}
	if p.Stage != 1 || p.Presses != 0 {
		t.Fatal(p)
	}
}
func TestLastDataStageClearsInsteadOfOverflowing(t *testing.T) {
	p := Progress{Stage: 2}
	for range Stages[2].Target {
		p.Step(Input{Solve: true})
	}
	if !p.Clear || p.Stage != 2 {
		t.Fatal(p)
	}
}
func TestRestartMakesFreshProgress(t *testing.T) {
	p := Progress{Stage: 2, Clear: true}
	p.Step(Input{Restart: true})
	if p != (Progress{}) {
		t.Fatal(p)
	}
}
