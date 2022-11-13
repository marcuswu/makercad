package utils

import "math"

// StandardCompare is the tolerance for a standard float64 comparison
const StandardCompare = 0.000001

// FloatCompare returns 0 if floats are equal, -1 if a < b, 1 if a > b using a tolerance
func FloatCompare(a float64, b float64, tol float64) int {
	if math.Abs(a-b) < tol {
		return 0
	}
	if a-b < 0-tol {
		return -1
	}
	return 1
}

// StandardFloatCompare returns 0 if floats are equal, -1 if a < b, 1 if a > b using a tolerance
func StandardFloatCompare(a float64, b float64) int {
	return FloatCompare(a, b, StandardCompare)
}
