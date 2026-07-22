package main

import "testing"

func TestExchangeRequiresOverlappingActionBoxes(t *testing.T) {
	if !exchangeOverlaps(160, 240, strike, guard) {
		t.Fatal("close strike and guard should overlap")
	}
	if exchangeOverlaps(100, 340, strike, guard) {
		t.Fatal("distant strike and guard must whiff")
	}
}

func TestPriorityRunsAfterContact(t *testing.T) {
	if !wins(guard, strike) || !wins(throwAct, guard) || !wins(strike, throwAct) {
		t.Fatal("strike/guard/throw priority table is incomplete")
	}
}
