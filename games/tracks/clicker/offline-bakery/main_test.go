package main

import "testing"

func TestMilestoneFor(t *testing.T) {
	tests := []struct {
		total float64
		want  int
	}{{0, 0}, {499.9, 0}, {500, 1}, {4999.9, 1}, {5000, 2}, {24999.9, 2}, {25000, 3}}
	for _, tt := range tests {
		if got := milestoneFor(tt.total); got != tt.want {
			t.Fatalf("milestoneFor(%v) = %d, want %d", tt.total, got, tt.want)
		}
	}
}

func TestProductionAndNextTarget(t *testing.T) {
	g := game{ovens: 2, mixers: 1, shops: 1, cost: 90, mixerCost: 70, shopCost: 900}
	if got, want := g.rate(), 99.0; got != want {
		t.Fatalf("rate = %v, want %v", got, want)
	}
	if got, want := g.tapPower(), 2.5; got != want {
		t.Fatalf("tapPower = %v, want %v", got, want)
	}
	name, cost := g.nextTarget()
	if name != "MIXER" || cost != 70 {
		t.Fatalf("nextTarget = (%q, %v), want MIXER 70", name, cost)
	}
}
