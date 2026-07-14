package starter

type State struct {
	Room   int
	HasKey bool
}

func EnterNext(s State) State {
	if s.Room < 3 {
		s.Room++
	}
	return s
}
