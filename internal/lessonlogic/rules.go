// Package lessonlogic contains small, deterministic game rules.
//
// The package deliberately does not import Ebitengine. Games can read input in
// Update, pass plain numbers to these functions, and test the returned values
// without opening a window.
package lessonlogic

import "math"

// PointInCircle reports whether a point is inside or on a circle.
func PointInCircle(pointX, pointY, centerX, centerY, radius float64) bool {
	dx := pointX - centerX
	dy := pointY - centerY
	return math.Hypot(dx, dy) <= radius
}

// BouncedPosition advances a marker by one tick and reflects it at either end.
// The returned position is always between min and max.
func BouncedPosition(position, speed, min, max float64) (nextPosition, nextSpeed float64) {
	nextPosition = position + speed
	nextSpeed = speed
	if nextPosition < min || nextPosition > max {
		nextSpeed = -nextSpeed
		nextPosition = math.Max(min, math.Min(max, nextPosition))
	}
	return nextPosition, nextSpeed
}

// TimingScore converts distance from the center of a timing meter into points
// and a label shown to the player.
func TimingScore(distance float64) (points int, label string) {
	switch {
	case distance <= 8:
		return 100, "PERFECT +100"
	case distance <= 28:
		return 50, "GREAT +50"
	case distance <= 55:
		return 10, "GOOD +10"
	default:
		return 0, "MISS"
	}
}

// AdvanceGauge adds speed to an action gauge for one tick. When it reaches the
// limit, the gauge resets and ready becomes true.
func AdvanceGauge(gauge, speed, limit int) (nextGauge int, ready bool) {
	nextGauge = gauge + speed
	if nextGauge >= limit {
		return 0, true
	}
	return nextGauge, false
}
