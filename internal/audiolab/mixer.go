package audiolab

// Mixer is audio policy only. Platform player calls consume these values later.
type Mixer struct {
	BGM, SE float64
	Paused  bool
	duck    int
}

func NewMixer() Mixer { return Mixer{BGM: 1, SE: 1} }
func (m *Mixer) TriggerImportantSE(frames int) {
	if frames > m.duck {
		m.duck = frames
	}
}
func (m *Mixer) Tick() {
	if m.Paused {
		return
	}
	if m.duck > 0 {
		m.duck--
	}
}
func (m Mixer) BGMGain() float64 {
	if m.Paused {
		return 0
	}
	if m.duck > 0 {
		return m.BGM * .35
	}
	return m.BGM
}
func (m Mixer) SEGain() float64 {
	if m.Paused {
		return 0
	}
	return m.SE
}
