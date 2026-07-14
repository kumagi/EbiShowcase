package uilab

import "testing"

func TestAlignAndWrap(t *testing.T) {
	if X(100, 40, Center) != 80 || X(100, 40, Right) != 60 {
		t.Fatal("alignment")
	}
	en := Wrap("make small games together", 10, false)
	if len(en) < 2 {
		t.Fatal(en)
	}
	ja := Wrap("小さなゲームを一緒に作ろう", 5, true)
	if len(ja) < 2 {
		t.Fatal(ja)
	}
}
