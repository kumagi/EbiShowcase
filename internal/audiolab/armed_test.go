package audiolab

import "testing"

func TestGateOnlyArmsAfterGesture(t *testing.T) {
	var g Gate
	if g.Arm(false) || g.Armed() {
		t.Fatal("audio armed without gesture")
	}
	if !g.Arm(true) || !g.Armed() {
		t.Fatal("gesture did not arm audio")
	}
}

func TestOneShotProducesPCM(t *testing.T) {
	for _, wave := range []Wave{Sine, Square, Noise} {
		b := OneShot(wave, 440, .1)
		if len(b) != int(.1*SampleRate)*4 {
			t.Fatalf("wave %d size %d", wave, len(b))
		}
		nonzero := false
		for _, v := range b {
			if v != 0 {
				nonzero = true
				break
			}
		}
		if !nonzero {
			t.Fatalf("wave %d is silent", wave)
		}
	}
}

func TestADSRHasAttackSustainAndRelease(t *testing.T) {
	a := ADSR{.1, .1, .5, .2}
	if a.Level(.05, .5) <= 0 || a.Level(.05, .5) >= 1 {
		t.Fatal("attack")
	}
	if a.Level(.3, .5) != .5 {
		t.Fatal("sustain")
	}
	if a.Level(.6, .5) >= .5 || a.Level(.6, .5) <= 0 {
		t.Fatal("release")
	}
	if a.Level(.8, .5) != 0 {
		t.Fatal("end")
	}
}
