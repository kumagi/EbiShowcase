package uilab

type Focus struct {
	Index, Count int
	Disabled     map[int]bool
}

func (f *Focus) Move(delta int) {
	if f.Count < 1 {
		return
	}
	for n := 0; n < f.Count; n++ {
		f.Index = (f.Index + delta + f.Count) % f.Count
		if !f.Disabled[f.Index] {
			return
		}
	}
}
func (f Focus) Activate() bool { return f.Count > 0 && !f.Disabled[f.Index] }
