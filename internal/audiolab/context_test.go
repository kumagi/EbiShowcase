package audiolab

import "testing"

func TestContextIsReusedAcrossGameRestarts(t *testing.T) {
	first := Context()
	second := Context()
	if first == nil {
		t.Fatal("Context() returned nil")
	}
	if first != second {
		t.Fatal("Context() created a second Ebitengine audio context")
	}
}
