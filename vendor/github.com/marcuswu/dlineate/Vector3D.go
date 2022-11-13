package dlineate

import "math"

type Vector3D struct {
	X float64
	Y float64
	Z float64
}

func NewVector(x float64, y float64, z float64) *Vector3D {
	v := new(Vector3D)
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// Dot product with another vector
func (v *Vector3D) Dot(u *Vector3D) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

// SquareMagnitude returns the squared magnitude of the vector
func (v *Vector3D) SquareMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Magnitude returns the magnitude of the vector
func (v *Vector3D) Magnitude() float64 {
	return math.Sqrt(v.SquareMagnitude())
}

// UnitVector returns a unit vector with the same direction
func (v *Vector3D) UnitVector() (*Vector3D, bool) {
	mag := v.Magnitude()
	if mag == 0 {
		return nil, false
	}
	return &Vector3D{v.X / mag, v.Y / mag, v.Z / mag}, true
}
