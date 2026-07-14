package audiolab

import (
	"encoding/binary"
	"math"
)

const SampleRate = 48000

type Wave int

const (
	Sine Wave = iota
	Square
	Noise
)

type ADSR struct{ Attack, Decay, Sustain, Release float64 }

// Level returns the envelope amplitude at time t for a sound whose held part
// ends at hold. It is pure math, reusable by every waveform family.
func (a ADSR) Level(t, hold float64) float64 {
	if t < 0 {
		return 0
	}
	if t < a.Attack {
		return t / a.Attack
	}
	if t < a.Attack+a.Decay {
		return 1 - (1-a.Sustain)*(t-a.Attack)/a.Decay
	}
	if t < hold {
		return a.Sustain
	}
	if t < hold+a.Release {
		return a.Sustain * (1 - (t-hold)/a.Release)
	}
	return 0
}
func Family(name string) (Wave, float64, ADSR) {
	switch name {
	case "magic":
		return Sine, 660, ADSR{.05, .18, .65, .45}
	case "hit":
		return Noise, 180, ADSR{.005, .04, .1, .10}
	default:
		return Square, 440, ADSR{.01, .08, .45, .20}
	}
}

func OneShot(w Wave, hz, duration float64) []byte {
	n := int(duration * SampleRate)
	b := make([]byte, n*4)
	for i := 0; i < n; i++ {
		t := float64(i) / SampleRate
		v := math.Sin(2 * math.Pi * hz * t)
		if w == Square {
			if v >= 0 {
				v = 1
			} else {
				v = -1
			}
		}
		if w == Noise {
			v = math.Sin(float64(i*i) * 12.9898)
		}
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(float32(v*math.Exp(-t*18)*.24)))
	}
	return b
}
