// SPDX-License-Identifier: Apache-2.0
package lessonlogic

import (
	"math"
	"testing"
)

const floatTolerance = 1e-9

func TestClassifyFallingObject(t *testing.T) {
	testCases := []struct {
		name   string
		caught bool
		y      float64
		want   FallingOutcome
	}{
		{name: "visible object stays in play", y: 719, want: FallingKeep},
		{name: "object on bottom edge stays in play", y: 720, want: FallingKeep},
		{name: "object below bottom is missed", y: 721, want: FallingMissed},
		{name: "caught object is collected", caught: true, y: 500, want: FallingCaught},
		{name: "catch takes priority below bottom", caught: true, y: 721, want: FallingCaught},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ClassifyFallingObject(tc.caught, tc.y, 720)

			if got != tc.want {
				t.Fatalf("outcome = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntegrateGravity(t *testing.T) {
	testCases := []struct {
		name                        string
		position, velocity, gravity float64
		wantPosition, wantVelocity  float64
	}{
		{name: "upward velocity slows before position advances", position: 100, velocity: -7.4, gravity: 0.42, wantPosition: 93.02, wantVelocity: -6.98},
		{name: "downward velocity grows before position advances", position: 100, velocity: 2, gravity: 0.42, wantPosition: 102.42, wantVelocity: 2.42},
		{name: "zero gravity preserves velocity for one tick", position: 100, velocity: 2, gravity: 0, wantPosition: 102, wantVelocity: 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotPosition, gotVelocity := IntegrateGravity(tc.position, tc.velocity, tc.gravity)

			if math.Abs(gotPosition-tc.wantPosition) > floatTolerance {
				t.Errorf("position = %v, want %v", gotPosition, tc.wantPosition)
			}
			if math.Abs(gotVelocity-tc.wantVelocity) > floatTolerance {
				t.Errorf("velocity = %v, want %v", gotVelocity, tc.wantVelocity)
			}
		})
	}
}

func TestExitScore(t *testing.T) {
	testCases := []struct {
		name string
		y    float64
		want int
	}{
		{name: "ball inside field gives no score", y: 360, want: -1},
		{name: "ball on upper boundary gives no score", y: -20, want: -1},
		{name: "ball above upper boundary scores for player", y: -20.1, want: 0},
		{name: "ball on lower boundary gives no score", y: 740, want: -1},
		{name: "ball below lower boundary scores for CPU", y: 740.1, want: 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ExitScore(tc.y, -20, 740)

			if got != tc.want {
				t.Fatalf("score side = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestAimVelocity(t *testing.T) {
	testCases := []struct {
		name          string
		dx, dy, speed float64
		wantX, wantY  float64
	}{
		{name: "three four five direction scales to requested speed", dx: 3, dy: 4, speed: 10, wantX: 6, wantY: 8},
		{name: "negative direction preserves left and upward signs", dx: -3, dy: -4, speed: 10, wantX: -6, wantY: -8},
		{name: "zero direction produces stopped bullet", speed: 10, wantX: 0, wantY: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotX, gotY := AimVelocity(tc.dx, tc.dy, tc.speed)

			if math.Abs(gotX-tc.wantX) > floatTolerance {
				t.Errorf("x velocity = %v, want %v", gotX, tc.wantX)
			}
			if math.Abs(gotY-tc.wantY) > floatTolerance {
				t.Errorf("y velocity = %v, want %v", gotY, tc.wantY)
			}
		})
	}
}

func TestSpendLife(t *testing.T) {
	testCases := []struct {
		name         string
		lives        int
		wantLives    int
		wantGameOver bool
	}{
		{name: "spending one of three lives continues game", lives: 3, wantLives: 2, wantGameOver: false},
		{name: "spending penultimate life leaves one life", lives: 2, wantLives: 1, wantGameOver: false},
		{name: "spending last life ends game", lives: 1, wantLives: 0, wantGameOver: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotLives, gotGameOver := SpendLife(tc.lives)

			if gotLives != tc.wantLives {
				t.Errorf("lives = %d, want %d", gotLives, tc.wantLives)
			}
			if gotGameOver != tc.wantGameOver {
				t.Errorf("gameOver = %v, want %v", gotGameOver, tc.wantGameOver)
			}
		})
	}
}

func TestSnakeStepInterval(t *testing.T) {
	testCases := []struct {
		name         string
		score        int
		wantInterval int
	}{
		{name: "zero score uses ten tick interval", score: 0, wantInterval: 10},
		{name: "two points keep ten tick interval", score: 2, wantInterval: 10},
		{name: "three points reduce interval by one", score: 3, wantInterval: 9},
		{name: "eighteen points reach four tick floor", score: 18, wantInterval: 4},
		{name: "high score stays at four tick floor", score: 99, wantInterval: 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotInterval := SnakeStepInterval(tc.score)

			if gotInterval != tc.wantInterval {
				t.Fatalf("interval = %d, want %d", gotInterval, tc.wantInterval)
			}
		})
	}
}

func TestAdvanceTween(t *testing.T) {
	testCases := []struct {
		name         string
		progress     float64
		wantProgress float64
		wantComplete bool
	}{
		{name: "step below end stays incomplete", progress: 0.28, wantProgress: 0.42, wantComplete: false},
		{name: "step ending just below one stays incomplete", progress: 0.85, wantProgress: 0.99, wantComplete: false},
		{name: "step past end clamps and completes", progress: 0.98, wantProgress: 1, wantComplete: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotProgress, gotComplete := AdvanceTween(tc.progress, 0.14)

			if math.Abs(gotProgress-tc.wantProgress) > floatTolerance {
				t.Errorf("progress = %v, want %v", gotProgress, tc.wantProgress)
			}
			if gotComplete != tc.wantComplete {
				t.Errorf("complete = %v, want %v", gotComplete, tc.wantComplete)
			}
		})
	}
}

func TestHorizontalVelocity(t *testing.T) {
	testCases := []struct {
		name        string
		vx          float64
		left, right bool
		want        float64
	}{
		{name: "left input accelerates left", left: true, want: -0.65},
		{name: "right input accelerates right", right: true, want: 0.65},
		{name: "no input applies friction", vx: 5, want: 3.9},
		{name: "acceleration stops at speed cap", vx: 5.8, right: true, want: 6},
		{name: "opposite inputs cancel from rest", left: true, right: true, want: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := HorizontalVelocity(tc.vx, tc.left, tc.right, 0.65, 0.78, 6)

			if math.Abs(got-tc.want) > floatTolerance {
				t.Fatalf("velocity = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestVerticalVelocity(t *testing.T) {
	testCases := []struct {
		name           string
		vy             float64
		jump, grounded bool
		wantVY         float64
		wantLeftGround bool
	}{
		{name: "falling applies gravity", vy: 2, wantVY: 2.62, wantLeftGround: false},
		{name: "grounded jump leaves ground before gravity", jump: true, grounded: true, wantVY: -11.88, wantLeftGround: true},
		{name: "air jump is ignored while gravity continues", vy: 2, jump: true, wantVY: 2.62, wantLeftGround: false},
		{name: "large downward speed stops at fall cap", vy: 13.8, wantVY: 14, wantLeftGround: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotVY, gotLeftGround := VerticalVelocity(tc.vy, tc.jump, tc.grounded, -12.5, 0.62, 14)

			if math.Abs(gotVY-tc.wantVY) > floatTolerance {
				t.Errorf("vertical velocity = %v, want %v", gotVY, tc.wantVY)
			}
			if gotLeftGround != tc.wantLeftGround {
				t.Errorf("leftGround = %v, want %v", gotLeftGround, tc.wantLeftGround)
			}
		})
	}
}

func TestEnemyMode(t *testing.T) {
	testCases := []struct {
		name     string
		current  int
		distance float64
		want     int
	}{
		{name: "distance below chase boundary starts chase", distance: 164.9, want: 1},
		{name: "exact chase boundary keeps wander mode", distance: 165, want: 0},
		{name: "middle distance keeps wander mode", distance: 200, want: 0},
		{name: "middle distance keeps chase mode", current: 1, distance: 200, want: 1},
		{name: "exact wander boundary keeps chase mode", current: 1, distance: 230, want: 1},
		{name: "distance above wander boundary returns to wander", current: 1, distance: 230.1, want: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := EnemyMode(tc.current, tc.distance, 165, 230)

			if got != tc.want {
				t.Fatalf("mode = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestCircleHit(t *testing.T) {
	testCases := []struct {
		name string
		bx   float64
		want bool
	}{
		{name: "overlapping circles hit", bx: 9, want: true},
		{name: "circles touching at edges do not hit", bx: 10, want: false},
		{name: "separated circles do not hit", bx: 10.1, want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := CircleHit(0, 0, 4, tc.bx, 0, 6)

			if got != tc.want {
				t.Fatalf("hit = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestOutside(t *testing.T) {
	testCases := []struct {
		name string
		x, y float64
		want bool
	}{
		{name: "point inside field is not outside", x: 240, y: 360, want: false},
		{name: "point on margin is not outside", x: -30, y: 100, want: false},
		{name: "point beyond margin is outside", x: -30.1, y: 100, want: true},
		{name: "point beyond right margin is outside", x: 510.1, y: 100, want: true},
		{name: "point beyond top margin is outside", x: 100, y: -30.1, want: true},
		{name: "point beyond bottom margin is outside", x: 100, y: 750.1, want: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Outside(tc.x, tc.y, 0, 0, 480, 720, 30)

			if got != tc.want {
				t.Fatalf("outside = %v, want %v", got, tc.want)
			}
		})
	}
}
