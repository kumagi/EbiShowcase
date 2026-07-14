package starter

import "testing"

func TestEnterNextStopsAtThirdRoom(t *testing.T) {
	if got := EnterNext(State{Room: 3}); got.Room != 3 {
		t.Fatal(got)
	}
}
