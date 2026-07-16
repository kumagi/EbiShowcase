// SPDX-License-Identifier: Apache-2.0
package lessonlogic

import (
	"math"
	"testing"
)

func TestClassifyFallingObject(t *testing.T) {
	tests := []struct {
		name   string
		caught bool
		y      float64
		want   FallingOutcome
	}{
		{name: "still visible", y: 719, want: FallingKeep},
		{name: "exactly at bottom", y: 720, want: FallingKeep},
		{name: "below bottom", y: 721, want: FallingMissed},
		{name: "caught", caught: true, y: 500, want: FallingCaught},
		{name: "catch wins at bottom", caught: true, y: 721, want: FallingCaught},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClassifyFallingObject(tt.caught, tt.y, 720); got != tt.want {
				t.Fatalf("ClassifyFallingObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegrateGravity(t *testing.T) {
	tests := []struct {
		name                        string
		position, velocity, gravity float64
		wantPosition, wantVelocity  float64
	}{
		{name: "after flap", position: 100, velocity: -7.4, gravity: 0.42, wantPosition: 93.02, wantVelocity: -6.98},
		{name: "falling", position: 100, velocity: 2, gravity: 0.42, wantPosition: 102.42, wantVelocity: 2.42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			position, velocity := IntegrateGravity(tt.position, tt.velocity, tt.gravity)
			if math.Abs(position-tt.wantPosition) > 1e-9 || math.Abs(velocity-tt.wantVelocity) > 1e-9 {
				t.Fatalf("IntegrateGravity() = (%v, %v), want (%v, %v)", position, velocity, tt.wantPosition, tt.wantVelocity)
			}
		})
	}
}

func TestExitScore(t *testing.T) {
	tests := []struct {
		name string
		y    float64
		want int
	}{
		{name: "inside", y: 360, want: -1},
		{name: "upper boundary stays", y: -20, want: -1},
		{name: "leaves top", y: -20.1, want: 0},
		{name: "lower boundary stays", y: 740, want: -1},
		{name: "leaves bottom", y: 740.1, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitScore(tt.y, -20, 740); got != tt.want {
				t.Fatalf("ExitScore(%v) = %d, want %d", tt.y, got, tt.want)
			}
		})
	}
}

func TestAimVelocity(t *testing.T) {
	tests := []struct {
		name          string
		dx, dy, speed float64
		wantX, wantY  float64
	}{
		{name: "three four five", dx: 3, dy: 4, speed: 10, wantX: 6, wantY: 8},
		{name: "zero direction", speed: 10, wantX: 0, wantY: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vx, vy := AimVelocity(tt.dx, tt.dy, tt.speed)
			if math.Abs(vx-tt.wantX) > 1e-9 || math.Abs(vy-tt.wantY) > 1e-9 {
				t.Fatalf("AimVelocity() = (%v, %v), want (%v, %v)", vx, vy, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestSpendLife(t *testing.T) {
	tests := []struct {
		lives    int
		want     int
		gameOver bool
	}{
		{lives: 3, want: 2, gameOver: false},
		{lives: 1, want: 0, gameOver: true},
	}
	for _, tt := range tests {
		got, over := SpendLife(tt.lives)
		if got != tt.want || over != tt.gameOver {
			t.Fatalf("SpendLife(%d) = (%d, %v), want (%d, %v)", tt.lives, got, over, tt.want, tt.gameOver)
		}
	}
}

func TestSnakeStepInterval(t *testing.T) {
	tests := []struct{ score, want int }{{0, 10}, {3, 9}, {18, 4}, {99, 4}}
	for _, tt := range tests {
		if got := SnakeStepInterval(tt.score); got != tt.want {
			t.Errorf("SnakeStepInterval(%d) = %d, want %d", tt.score, got, tt.want)
		}
	}
}

func TestAdvanceTween(t *testing.T) {
	tests := []struct {
		name         string
		progress     float64
		wantProgress float64
		wantComplete bool
	}{
		{name: "in progress", progress: 0.28, wantProgress: 0.42},
		{name: "reaches end", progress: 0.98, wantProgress: 1, wantComplete: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, complete := AdvanceTween(tt.progress, 0.14)
			if math.Abs(got-tt.wantProgress) > 1e-9 || complete != tt.wantComplete {
				t.Fatalf("AdvanceTween() = (%v, %v), want (%v, %v)", got, complete, tt.wantProgress, tt.wantComplete)
			}
		})
	}
}

func TestHorizontalVelocity(t *testing.T) {
	tests := []struct {
		name        string
		vx          float64
		left, right bool
		want        float64
	}{
		{name: "accelerate left", left: true, want: -0.65},
		{name: "accelerate right", right: true, want: 0.65},
		{name: "friction", vx: 5, want: 3.9},
		{name: "cap", vx: 5.8, right: true, want: 6},
		{name: "opposites cancel", left: true, right: true, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HorizontalVelocity(tt.vx, tt.left, tt.right, 0.65, 0.78, 6); math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("HorizontalVelocity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerticalVelocity(t *testing.T) {
	tests := []struct {
		name           string
		vy             float64
		jump, grounded bool
		wantVY         float64
		wantLeftGround bool
	}{
		{name: "gravity", vy: 2, wantVY: 2.62},
		{name: "grounded jump", jump: true, grounded: true, wantVY: -11.88, wantLeftGround: true},
		{name: "air jump ignored", vy: 2, jump: true, wantVY: 2.62},
		{name: "fall cap", vy: 13.8, wantVY: 14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, leftGround := VerticalVelocity(tt.vy, tt.jump, tt.grounded, -12.5, 0.62, 14)
			if math.Abs(got-tt.wantVY) > 1e-9 || leftGround != tt.wantLeftGround {
				t.Fatalf("VerticalVelocity() = (%v, %v), want (%v, %v)", got, leftGround, tt.wantVY, tt.wantLeftGround)
			}
		})
	}
}

func TestEnemyMode(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		distance float64
		want     int
	}{
		{name: "starts chase", distance: 164.9, want: 1},
		{name: "keeps wander in band", distance: 200, want: 0},
		{name: "keeps chase in band", current: 1, distance: 200, want: 1},
		{name: "returns to wander", current: 1, distance: 230.1, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnemyMode(tt.current, tt.distance, 165, 230); got != tt.want {
				t.Fatalf("EnemyMode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCircleHit(t *testing.T) {
	tests := []struct {
		name string
		bx   float64
		want bool
	}{{name: "overlap", bx: 9, want: true}, {name: "touch is not overlap", bx: 10}, {name: "separate", bx: 10.1}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CircleHit(0, 0, 4, tt.bx, 0, 6); got != tt.want {
				t.Fatalf("CircleHit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutside(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		want bool
	}{{name: "inside", x: 240, y: 360}, {name: "on margin", x: -30, y: 100}, {name: "past margin", x: -30.1, y: 100, want: true}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Outside(tt.x, tt.y, 0, 0, 480, 720, 30); got != tt.want {
				t.Fatalf("Outside() = %v, want %v", got, tt.want)
			}
		})
	}
}
