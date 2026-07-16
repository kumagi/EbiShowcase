// SPDX-License-Identifier: Apache-2.0
package lessonlogic

import "math"

// FallingOutcome is the result of advancing one falling collectible.
type FallingOutcome int

const (
	FallingKeep FallingOutcome = iota
	FallingCaught
	FallingMissed
)

// ClassifyFallingObject decides whether an already-advanced object remains,
// was caught, or fell below the screen. Catch wins if both flags are true.
func ClassifyFallingObject(caught bool, y, bottom float64) FallingOutcome {
	if caught {
		return FallingCaught
	}
	if y > bottom {
		return FallingMissed
	}
	return FallingKeep
}

// IntegrateGravity advances vertical velocity and position by one game tick.
func IntegrateGravity(position, velocity, acceleration float64) (nextPosition, nextVelocity float64) {
	nextVelocity = velocity + acceleration
	nextPosition = position + nextVelocity
	return nextPosition, nextVelocity
}

// ExitScore reports which side receives a point when a ball leaves a vertical
// play field. -1 means no score, 0 means the upper side, and 1 the lower side.
func ExitScore(y, min, max float64) int {
	if y < min {
		return 0
	}
	if y > max {
		return 1
	}
	return -1
}

// AimVelocity returns a vector of the requested length toward (dx, dy).
// A zero-length direction intentionally returns a stopped vector.
func AimVelocity(dx, dy, speed float64) (vx, vy float64) {
	distance := math.Hypot(dx, dy)
	if distance == 0 {
		return 0, 0
	}
	return dx / distance * speed, dy / distance * speed
}

// SpendLife removes one life and reports whether none remain.
func SpendLife(lives int) (nextLives int, gameOver bool) {
	nextLives = lives - 1
	return nextLives, nextLives <= 0
}

// SnakeStepInterval converts score into the number of ticks between moves.
func SnakeStepInterval(score int) int {
	return max(4, 10-score/3)
}

// AdvanceTween advances a normalized 0..1 animation progress value.
func AdvanceTween(progress, step float64) (next float64, complete bool) {
	next = progress + step
	if next >= 1 {
		return 1, true
	}
	return next, false
}

// HorizontalVelocity applies left/right input, friction, and a speed cap.
func HorizontalVelocity(vx float64, left, right bool, acceleration, friction, maxSpeed float64) float64 {
	if left {
		vx -= acceleration
	}
	if right {
		vx += acceleration
	}
	if !left && !right {
		vx *= friction
	}
	return math.Max(-maxSpeed, math.Min(maxSpeed, vx))
}

// VerticalVelocity applies a grounded jump first, then gravity and a fall cap.
func VerticalVelocity(vy float64, jump, onGround bool, jumpSpeed, gravity, maxFall float64) (nextVY float64, leftGround bool) {
	if jump && onGround {
		vy = jumpSpeed
		leftGround = true
	}
	return math.Min(vy+gravity, maxFall), leftGround
}

// EnemyMode applies hysteresis to an enemy's wander/chase state.
// Mode 0 wanders and mode 1 chases.
func EnemyMode(current int, distance, chaseBelow, wanderAbove float64) int {
	if distance < chaseBelow {
		return 1
	}
	if distance > wanderAbove {
		return 0
	}
	return current
}

// CircleHit reports whether two circles overlap.
func CircleHit(ax, ay, ar, bx, by, br float64) bool {
	return math.Hypot(ax-bx, ay-by) < ar+br
}

// Outside reports whether a point has left a rectangle plus a margin.
func Outside(x, y, minX, minY, maxX, maxY, margin float64) bool {
	return x < minX-margin || x > maxX+margin || y < minY-margin || y > maxY+margin
}
