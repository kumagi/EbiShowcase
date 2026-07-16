package main

import "testing"

func TestHandCardsKeepTheirWidthWhenCardsArePlayed(t *testing.T) {
	_, three, _ := handLayout(3)
	_, two, _ := handLayout(2)
	_, one, _ := handLayout(1)
	_, five, _ := handLayout(5)
	if five != three || three != two || two != one {
		t.Fatalf("card width changed with hand size: five=%v three=%v two=%v one=%v", five, three, two, one)
	}
}

func TestHandCardAtUsesCenteredFixedCards(t *testing.T) {
	start, cardWidth, step := handLayout(2)
	for i := 0; i < 2; i++ {
		x := int(start + float32(i)*step + cardWidth/2)
		if got := handCardAt(x, 2); got != i {
			t.Fatalf("handCardAt(%d, 2) = %d, want %d", x, got, i)
		}
	}
}
