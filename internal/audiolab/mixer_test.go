package audiolab

import "testing"

func TestMixerDucksAndPauses(t *testing.T) {
	m := NewMixer()
	m.TriggerImportantSE(2)
	if m.BGMGain() >= 1 {
		t.Fatal("not ducked")
	}
	m.Tick()
	m.Tick()
	if m.BGMGain() != 1 {
		t.Fatal("did not recover")
	}
	m.Paused = true
	if m.BGMGain() != 0 || m.SEGain() != 0 {
		t.Fatal("pause")
	}
}
