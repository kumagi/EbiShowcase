package cameralab

import "math"

type Impact struct {
	Frames, Stop int
	Strength     float64
}

func (i *Impact) Trigger(frames, stop int, strength float64) {
	i.Frames = frames
	i.Stop = stop
	i.Strength = strength
}
func (i *Impact) Tick() {
	if i.Stop > 0 {
		i.Stop--
		return
	}
	if i.Frames > 0 {
		i.Frames--
	}
}
func (i Impact) Frozen() bool { return i.Stop > 0 }
func (i Impact) Offset(frame int) Vec {
	if i.Frames <= 0 {
		return Vec{}
	}
	scale := i.Strength * float64(i.Frames) / 12
	return Vec{math.Sin(float64(frame)*12.7) * scale, math.Cos(float64(frame)*9.1) * scale}
}
