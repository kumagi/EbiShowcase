package main

import "testing"

func TestAttackAnimationVisitsAllEightFramesInOrder(t *testing.T) {
	for _, kind := range []int{attackJab, attackHeavy} {
		spec := moveSpecs[kind]
		seen := [8]bool{}
		previous := -1
		for tick := 0; tick < spec.total; tick++ {
			frame := attackAnimationFrame(fighter{attackKind: kind, attackTick: tick})
			if frame < previous {
				t.Fatalf("%s animation moved backward: frame %d after %d", spec.name, frame, previous)
			}
			if frame < 0 || frame >= len(seen) {
				t.Fatalf("%s animation frame %d is outside the eight-frame sheet", spec.name, frame)
			}
			seen[frame] = true
			previous = frame
		}
		for frame, visited := range seen {
			if !visited {
				t.Fatalf("%s never displays animation frame %d", spec.name, frame)
			}
		}
	}
}

func TestAttacksHaveReadableStartupAndWhiffRecovery(t *testing.T) {
	for _, kind := range []int{attackJab, attackHeavy} {
		spec := moveSpecs[kind]
		activeTicks := spec.activeEnd - spec.activeStart
		recoveryTicks := spec.total - spec.activeEnd
		if spec.activeStart < 6 {
			t.Fatalf("%s startup = %d, want a readable anticipation", spec.name, spec.activeStart)
		}
		if recoveryTicks <= activeTicks*2 {
			t.Fatalf("%s recovery = %d, want a punishable whiff after %d active ticks", spec.name, recoveryTicks, activeTicks)
		}
	}
}

func TestAttackLocksOutButtonMashingUntilRecoveryEnds(t *testing.T) {
	f := fighter{hp: 100, guardHP: 100}
	g := &game{}
	g.startAttack(&f, attackJab)
	if f.canAct() {
		t.Fatal("fighter can act during jab startup")
	}
	g.startAttack(&f, attackHeavy)
	if f.attackKind != attackJab {
		t.Fatal("a second button press replaced an attack before recovery ended")
	}
	for f.attackKind != attackNone {
		g.advanceAttack(&f)
	}
	if !f.canAct() {
		t.Fatal("fighter cannot act after the full recovery finished")
	}
}

func TestHeavyIsTheGuardBreakingCommitment(t *testing.T) {
	jab, heavy := moveSpecs[attackJab], moveSpecs[attackHeavy]
	if heavy.guardDamage <= jab.guardDamage*2 {
		t.Fatalf("heavy guard damage = %d, want more than twice jab's %d", heavy.guardDamage, jab.guardDamage)
	}
	if heavy.damage <= jab.damage*2 {
		t.Fatalf("heavy damage = %d, want more than twice jab's %d", heavy.damage, jab.damage)
	}
	if heavy.total <= jab.total {
		t.Fatalf("heavy total = %d, want a longer commitment than jab's %d", heavy.total, jab.total)
	}
}
