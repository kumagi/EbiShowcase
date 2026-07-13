package lessonlogic

import "testing"

func TestPointInCircle(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		wantInside bool
	}{
		{name: "center", x: 100, y: 80, wantInside: true},
		{name: "just inside", x: 129, y: 80, wantInside: true},
		{name: "on the edge", x: 130, y: 80, wantInside: true},
		{name: "just outside", x: 131, y: 80, wantInside: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PointInCircle(tt.x, tt.y, 100, 80, 30)
			if got != tt.wantInside {
				t.Fatalf("PointInCircle(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.wantInside)
			}
		})
	}
}

func TestBouncedPosition(t *testing.T) {
	tests := []struct {
		name                    string
		position, speed         float64
		wantPosition, wantSpeed float64
	}{
		{name: "moves right", position: 40, speed: 3, wantPosition: 43, wantSpeed: 3},
		{name: "bounces at right edge", position: 99, speed: 3, wantPosition: 100, wantSpeed: -3},
		{name: "bounces at left edge", position: 1, speed: -3, wantPosition: 0, wantSpeed: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			position, speed := BouncedPosition(tt.position, tt.speed, 0, 100)
			if position != tt.wantPosition || speed != tt.wantSpeed {
				t.Fatalf("BouncedPosition() = (%v, %v), want (%v, %v)", position, speed, tt.wantPosition, tt.wantSpeed)
			}
		})
	}
}

func TestTimingScore(t *testing.T) {
	tests := []struct {
		distance   float64
		wantPoints int
		wantLabel  string
	}{
		{distance: 0, wantPoints: 100, wantLabel: "PERFECT +100"},
		{distance: 8, wantPoints: 100, wantLabel: "PERFECT +100"},
		{distance: 8.1, wantPoints: 50, wantLabel: "GREAT +50"},
		{distance: 28, wantPoints: 50, wantLabel: "GREAT +50"},
		{distance: 55, wantPoints: 10, wantLabel: "GOOD +10"},
		{distance: 55.1, wantPoints: 0, wantLabel: "MISS"},
	}

	for _, tt := range tests {
		points, label := TimingScore(tt.distance)
		if points != tt.wantPoints || label != tt.wantLabel {
			t.Errorf("TimingScore(%v) = (%d, %q), want (%d, %q)", tt.distance, points, label, tt.wantPoints, tt.wantLabel)
		}
	}
}

func TestAdvanceGauge(t *testing.T) {
	tests := []struct {
		name         string
		gauge, speed int
		wantGauge    int
		wantReady    bool
	}{
		{name: "charges", gauge: 400, speed: 12, wantGauge: 412, wantReady: false},
		{name: "one short", gauge: 987, speed: 12, wantGauge: 999, wantReady: false},
		{name: "exactly ready", gauge: 988, speed: 12, wantGauge: 0, wantReady: true},
		{name: "passes limit", gauge: 995, speed: 12, wantGauge: 0, wantReady: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gauge, ready := AdvanceGauge(tt.gauge, tt.speed, 1000)
			if gauge != tt.wantGauge || ready != tt.wantReady {
				t.Fatalf("AdvanceGauge() = (%d, %v), want (%d, %v)", gauge, ready, tt.wantGauge, tt.wantReady)
			}
		})
	}
}
