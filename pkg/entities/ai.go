package entities

import (
	"math"
)

// Distance calculates Euclidean distance between two points
func Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

// NextStepTowards returns the next (dx, dy) towards target with obstacle checking
func NextStepTowards(fromX, fromY, toX, toY int, isBlocked func(x, y int) bool) (int, int) {
	dx := toX - fromX
	dy := toY - fromY

	dist := Distance(fromX, fromY, toX, toY)
	if dist <= 1.0 {
		return 0, 0 // Already adjacent
	}

	stepX := 0
	stepY := 0

	if dx > 0 {
		stepX = 1
	} else if dx < 0 {
		stepX = -1
	}

	if dy > 0 {
		stepY = 1
	} else if dy < 0 {
		stepY = -1
	}

	// Try diagonal or primary axis
	if stepX != 0 && stepY != 0 {
		if !isBlocked(fromX+stepX, fromY+stepY) {
			return stepX, stepY
		}
	}

	// Try X first if distance in X is greater
	if math.Abs(float64(dx)) >= math.Abs(float64(dy)) {
		if stepX != 0 && !isBlocked(fromX+stepX, fromY) {
			return stepX, 0
		}
		if stepY != 0 && !isBlocked(fromX, fromY+stepY) {
			return 0, stepY
		}
	} else {
		if stepY != 0 && !isBlocked(fromX, fromY+stepY) {
			return 0, stepY
		}
		if stepX != 0 && !isBlocked(fromX+stepX, fromY) {
			return stepX, 0
		}
	}

	return 0, 0
}
