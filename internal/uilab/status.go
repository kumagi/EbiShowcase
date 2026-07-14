package uilab

type Status struct {
	Text   string
	Frames int
}

func (s *Status) Set(text string, frames int) { s.Text = text; s.Frames = frames }
func (s *Status) Tick() {
	if s.Frames > 0 {
		s.Frames--
	}
	if s.Frames == 0 {
		s.Text = ""
	}
}

type Scroll struct{ Offset, Max int }

func (s *Scroll) Move(delta int) {
	s.Offset += delta
	if s.Offset < 0 {
		s.Offset = 0
	}
	if s.Offset > s.Max {
		s.Offset = s.Max
	}
}
