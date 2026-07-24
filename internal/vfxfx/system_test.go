package vfxfx

import (
	"image/color"
	"testing"
)

func TestSystemCapsOptionalParticlesWithoutAffectingCallers(t *testing.T) {
	system := System{MaxParts: 5}

	system.Burst(10, 20, 40, 3, color.White, true)

	if got := system.Count(); got != 5 {
		t.Fatalf("Count() = %d, want budget cap 5", got)
	}
}

func TestRingGrowsInsteadOfDriftingSideways(t *testing.T) {
	var system System

	system.Ring(30, 40, 2, color.White)

	if len(system.Parts) != 1 {
		t.Fatalf("rings = %d, want 1", len(system.Parts))
	}
	ring := system.Parts[0]
	if ring.VX != 0 || ring.VY != 0 {
		t.Fatalf("ring velocity = (%.2f, %.2f), want stationary center", ring.VX, ring.VY)
	}
	if ring.EndScale <= ring.Scale {
		t.Fatalf("ring scale %.2f → %.2f, want growth", ring.Scale, ring.EndScale)
	}
}

func TestUpdateExpiresParticlesAndNeverLeavesNegativeFlash(t *testing.T) {
	system := System{
		Flash: 0.03,
		Parts: []Particle{{Life: 1, Max: 1, Scale: 1}},
	}

	system.Update()

	if system.Flash != 0 {
		t.Fatalf("Flash = %.2f, want 0", system.Flash)
	}
	if system.Count() != 0 {
		t.Fatalf("Count() = %d, want expired particle removed", system.Count())
	}
}
