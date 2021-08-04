package utils

import "math"

func ToRadians(degrees float64) float64 {
	return math.Pi * (degrees / 180.0)
}

func ToDegrees(radians float64) float64 {
	return 180 * (radians / math.Pi)
}
