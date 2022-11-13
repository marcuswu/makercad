package element

import (
	"fmt"
	"math"
)

// Vector represents a 2D vector
type Vector struct {
	X float64
	Y float64
}

// GetX return the x value of the vector
func (v *Vector) GetX() float64 {
	return v.X
}

// GetY return the y value of the vector
func (v *Vector) GetY() float64 {
	return v.Y
}

// Dot product with another vector
func (v *Vector) Dot(u *Vector) float64 {
	return v.X*u.X + v.Y*u.Y
}

// SquareMagnitude returns the squared magnitude of the vector
func (v *Vector) SquareMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Magnitude returns the magnitude of the vector
func (v *Vector) Magnitude() float64 {
	return math.Sqrt(v.SquareMagnitude())
}

// AngleTo returns the angle to another vector in radians
// https://stackoverflow.com/a/21484228
// With this math, counter clockwise is positive
func (v *Vector) AngleTo(u *Vector) float64 {
	angle := math.Atan2(u.Y, u.X) - math.Atan2(v.Y, v.X)
	if angle > math.Pi {
		angle -= 2 * math.Pi
	}
	if angle <= -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}

// Rotated returns a vector representing this vector rotated around the origin by angle radians
func (v *Vector) Rotated(angle float64) Vector {
	sinAngle := math.Sin(angle)
	cosAngle := math.Cos(angle)

	newX := v.X*cosAngle - v.Y*sinAngle
	newY := v.X*sinAngle + v.Y*cosAngle

	return Vector{newX, newY}
}

// Rotate rotates the vector around the origin by angle radians
func (v *Vector) Rotate(angle float64) {
	rotated := v.Rotated(angle)

	v.X = rotated.X
	v.Y = rotated.Y
}

// Translated returns a vector representing this vector by an x and y distance
func (v *Vector) Translated(dx float64, dy float64) Vector {
	return Vector{v.X + dx, v.Y + dy}
}

// Translate translates the vectory by an x and y distance
func (v *Vector) Translate(dx float64, dy float64) {
	translated := v.Translated(dx, dy)

	v.X = translated.X
	v.Y = translated.Y
}

// UnitVector returns a unit vector with the same direction
func (v *Vector) UnitVector() (*Vector, bool) {
	mag := v.Magnitude()
	if mag == 0 {
		return nil, false
	}
	return &Vector{v.X / mag, v.Y / mag}, true
}

// Scaled multiplies this vector by a magnitude
func (v *Vector) Scaled(scale float64) {
	v.X *= scale
	v.Y *= scale
}

func (v *Vector) String() string {
	return fmt.Sprintf("Vector((0,0),(%f,%f))", v.X, v.Y)
}
