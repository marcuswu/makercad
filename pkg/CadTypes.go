package MakerCad

import (
	"fmt"
	"libmakercad/internal/utils"
	. "libmakercad/third_party/occt"
)

type Vector struct {
	X float64
	Y float64
	Z float64
}

func NewVector(a ...interface{}) *Vector {
	argc := len(a)

	if argc != 3 && argc != 1 && argc != 0 {
		panic("No match for overloaded function call")
	}

	if argc == 0 {
		return &Vector{0, 0, 0}
	}

	if argc == 3 {
		return &Vector{a[0].(float64), a[1].(float64), a[2].(float64)}
	}

	switch v := a[0].(type) {
	case Gp_Pnt:
		return &Vector{v.X(), v.Y(), v.Z()}
	case Gp_Vec:
		return &Vector{v.X(), v.Y(), v.Z()}
	case Gp_Dir:
		return &Vector{v.X(), v.Y(), v.Z()}
	default:
		panic("Invalid value passed to NewVector")
	}
}

func (v *Vector) ToPoint() Gp_Pnt {
	return NewGp_Pnt(v.X, v.Y, v.Z)
}

func (v *Vector) ToVector() Gp_Vec {
	return NewGp_Vec(v.X, v.Y, v.Z)
}

func (v *Vector) ToString() string {
	return fmt.Sprintf("{ x: %f, y: %f, z: %f }\n", v.X, v.Y, v.Z)
}

type PlaneParameters struct {
	Location *Vector
	X        *Vector
	Y        *Vector
}

func NewPlaneParameters(a ...interface{}) *PlaneParameters {
	argc := len(a)

	if argc != 3 && argc != 1 && argc != 0 {
		panic("No match for overloaded function call")
	}

	if argc == 0 {
		return &PlaneParameters{NewVector(), NewVector(), NewVector()}
	}

	if argc == 1 {
		coord, ok := a[0].(Gp_Ax3)
		if !ok {
			panic("Invalid parameter for NewPlaneParameters")
		}
		return &PlaneParameters{
			NewVector(coord.Location()),
			NewVector(coord.XDirection()),
			NewVector(coord.YDirection()),
		}
	}

	loc, ok1 := a[0].(*Vector)
	xDir, ok2 := a[0].(*Vector)
	yDir, ok3 := a[0].(*Vector)
	if !ok1 || !ok2 || !ok3 {
		panic("Invalid parameters for NewPlaneParameters")
	}
	return &PlaneParameters{loc, xDir, yDir}
}

func (p *PlaneParameters) ToAx3() Gp_Ax3 {
	return NewGp_Ax3(
		p.Location.ToPoint(),
		NewGp_Dir(p.X.ToVector().Crossed(p.Y.ToVector())),
		NewGp_Dir(p.X.ToVector()),
	)
}

func (p *PlaneParameters) Rotated(a ...interface{}) *PlaneParameters {
	argc := len(a)
	if argc < 1 || argc > 3 {
		panic("No match for overloaded function call")
	}
	degrees := 0.0
	axis := NewVector(0, 0, 0)
	origin := NewVector(0, 0, 0)

	if argc < 2 {
		// calculate axis
		axis = NewVector(p.X.ToVector().Crossed(p.Y.ToVector()))
	}
	if argc < 3 {
		// use current location as origin
		origin = p.Location
	}
	coordinates := p.ToAx3()
	coordinates.Rotate(NewGp_Ax1(origin.ToPoint(), NewGp_Dir(axis.ToVector())), utils.ToRadians(degrees))
	return NewPlaneParameters(coordinates)
}

func (p *PlaneParameters) Translated(dir Vector) *PlaneParameters {
	return NewPlaneParameters(p.ToAx3().Translated(dir.ToVector()))
}

func (p *PlaneParameters) ToString() string {
	return fmt.Sprintf("{ location: %s, x: %s, y: %s }\n", p.Location.ToString(), p.X.ToString(), p.Y.ToString())
}

type Edge interface {
	IsLine() bool
	IsCircle() bool
	IsEllipse() bool
	GetLine() Line
	LineLength() float64
	GetCircle() Circle
	CircleRadius() float64
}

type Point interface {
	GetX() float64
	GetY() float64
	GetZ() float64
	SetX(float64)
	SetY(float64)

	Coincident(...interface{}) Point
	Horizontal(Point) Point
	Vertical(Point) Point
	Distance(...interface{}) Point
	HorizontalDistance(Point, float64) Point
	VerticalDistance(Point, float64) Point
	Construction(bool) Point
	ToString() string
}

type Line interface {
	Horizontal() Line
	Vertical() Line
	Length(float64) Line
	MidPoint(Point) Line
	Symmetric(Point, Point) Line
	Angle(Line, float64) Line
	Construction(bool) Line
	ToString()
	GetStart() Point
	GetEnd() Point
	ToVector() Vector
}

type Arc interface {
	Diameter(float64) Arc
	Construction(bool) Arc
	StartTangent(Line) Arc
	EndTangent(Line) Arc
	GetCenter() Point
	GetStart() Point
	GetEnd() Point
}

type Circle interface {
	Diameter(float64) Circle
	Construction(bool) Circle
	Equal(Circle) Circle
	GetDiameter() float64
	GetCenter() Point
}
