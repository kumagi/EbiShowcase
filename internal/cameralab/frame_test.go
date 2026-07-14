package cameralab

import "testing"

func TestFrameAndSafeRect(t *testing.T) {
	var f Frame
	f.Start(2)
	if f.Letterbox(720) == 0 {
		t.Fatal("no bars")
	}
	f.Tick()
	f.Tick()
	if f.Letterbox(720) != 0 {
		t.Fatal("bars remain")
	}
	x, y, w, h := SafeRect(480, 720, 24)
	if x != 24 || y != 24 || w != 432 || h != 672 {
		t.Fatal("safe")
	}
}
