package starter

import "testing"

func TestNextStage(t *testing.T) {
	if NextStage(2, true) != 3 || NextStage(3, true) != 3 {
		t.Fatal("stage limit")
	}
}
