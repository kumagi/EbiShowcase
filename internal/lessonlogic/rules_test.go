// SPDX-License-Identifier: Apache-2.0
package lessonlogic

import "testing"

func TestPointInCircle(t *testing.T) {
	testCases := []struct {
		name       string
		x, y       float64
		wantInside bool
	}{
		{name: "center is inside", x: 100, y: 80, wantInside: true},
		{name: "point one pixel inside is inside", x: 129, y: 80, wantInside: true},
		{name: "point on edge is inside", x: 130, y: 80, wantInside: true},
		{name: "point one pixel outside is outside", x: 131, y: 80, wantInside: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotInside := PointInCircle(tc.x, tc.y, 100, 80, 30)

			if gotInside != tc.wantInside {
				t.Fatalf("PointInCircle(%v, %v) = %v, want %v", tc.x, tc.y, gotInside, tc.wantInside)
			}
		})
	}
}

func TestBouncedPosition(t *testing.T) {
	testCases := []struct {
		name                    string
		position, speed         float64
		wantPosition, wantSpeed float64
	}{
		{name: "movement within bounds keeps direction", position: 40, speed: 3, wantPosition: 43, wantSpeed: 3},
		{name: "crossing right edge clamps and reverses", position: 99, speed: 3, wantPosition: 100, wantSpeed: -3},
		{name: "crossing left edge clamps and reverses", position: 1, speed: -3, wantPosition: 0, wantSpeed: 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotPosition, gotSpeed := BouncedPosition(tc.position, tc.speed, 0, 100)

			if gotPosition != tc.wantPosition {
				t.Errorf("position = %v, want %v", gotPosition, tc.wantPosition)
			}
			if gotSpeed != tc.wantSpeed {
				t.Errorf("speed = %v, want %v", gotSpeed, tc.wantSpeed)
			}
		})
	}
}

func TestTimingScore(t *testing.T) {
	testCases := []struct {
		name       string
		distance   float64
		wantPoints int
		wantLabel  string
	}{
		{name: "center scores perfect", distance: 0, wantPoints: 100, wantLabel: "PERFECT +100"},
		{name: "perfect outer edge scores perfect", distance: 8, wantPoints: 100, wantLabel: "PERFECT +100"},
		{name: "just beyond perfect scores great", distance: 8.1, wantPoints: 50, wantLabel: "GREAT +50"},
		{name: "great outer edge scores great", distance: 28, wantPoints: 50, wantLabel: "GREAT +50"},
		{name: "good outer edge scores good", distance: 55, wantPoints: 10, wantLabel: "GOOD +10"},
		{name: "just beyond good scores miss", distance: 55.1, wantPoints: 0, wantLabel: "MISS"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotPoints, gotLabel := TimingScore(tc.distance)

			if gotPoints != tc.wantPoints {
				t.Errorf("points = %d, want %d", gotPoints, tc.wantPoints)
			}
			if gotLabel != tc.wantLabel {
				t.Errorf("label = %q, want %q", gotLabel, tc.wantLabel)
			}
		})
	}
}

func TestAdvanceGauge(t *testing.T) {
	testCases := []struct {
		name         string
		gauge, speed int
		wantGauge    int
		wantReady    bool
	}{
		{name: "ordinary charge stays below ready", gauge: 400, speed: 12, wantGauge: 412, wantReady: false},
		{name: "one point below limit stays below ready", gauge: 987, speed: 12, wantGauge: 999, wantReady: false},
		{name: "reaching limit resets and becomes ready", gauge: 988, speed: 12, wantGauge: 0, wantReady: true},
		{name: "passing limit resets and becomes ready", gauge: 995, speed: 12, wantGauge: 0, wantReady: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotGauge, gotReady := AdvanceGauge(tc.gauge, tc.speed, 1000)

			if gotGauge != tc.wantGauge {
				t.Errorf("gauge = %d, want %d", gotGauge, tc.wantGauge)
			}
			if gotReady != tc.wantReady {
				t.Errorf("ready = %v, want %v", gotReady, tc.wantReady)
			}
		})
	}
}
