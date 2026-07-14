package cameralab

type Frame struct{ Transition int }

func (f *Frame) Start(frames int) { f.Transition = frames }
func (f *Frame) Tick() {
	if f.Transition > 0 {
		f.Transition--
	}
}
func (f Frame) Letterbox(height float64) float64 {
	if f.Transition <= 0 {
		return 0
	}
	return height * .08
}
func SafeRect(w, h, margin float64) (float64, float64, float64, float64) {
	return margin, margin, w - margin*2, h - margin*2
}
