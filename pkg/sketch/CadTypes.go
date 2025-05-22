package sketch

import (
	"fmt"

	"github.com/marcuswu/libmakercad/internal/utils"

	"github.com/marcuswu/gooccwrapper/gp"
)

type Vector struct {
	X float64
	Y float64
	Z float64
}

func NewVectorFromValues(x float64, y float64, z float64) *Vector {
	return &Vector{x, y, z}
}

type XYZer interface {
	X() float64
	Y() float64
	Z() float64
}

func NewVector(v XYZer) *Vector {
	return &Vector{v.X(), v.Y(), v.Z()}
}

func (v *Vector) ToPoint() gp.Pnt {
	return gp.NewPnt(v.X, v.Y, v.Z)
}

func (v *Vector) ToVector() gp.Vec {
	return gp.NewVec(v.X, v.Y, v.Z)
}

func (v *Vector) ToString() string {
	return fmt.Sprintf("{ x: %f, y: %f, z: %f }\n", v.X, v.Y, v.Z)
}

type PlaneParameters struct {
	Location *Vector
	Normal   *Vector
	X        *Vector
}

func NewPlaneParameters() *PlaneParameters {
	return &PlaneParameters{
		NewVectorFromValues(0, 0, 0),
		NewVectorFromValues(0, 0, 1),
		NewVectorFromValues(1, 0, 0),
	}
}

func NewPlaneParametersFromCoordinateSystem(coord gp.Ax3) *PlaneParameters {
	return &PlaneParameters{
		NewVector(coord.Location()),
		NewVector(NewVector(coord.YDirection()).ToVector().Crossed(NewVector(coord.XDirection()).ToVector())),
		NewVector(coord.XDirection()),
	}
}

func NewPlaneParametersFromVectors(loc *Vector, normal *Vector, xDir *Vector) *PlaneParameters {
	return &PlaneParameters{loc, normal, xDir}
}

func (p *PlaneParameters) Plane() gp.Ax3 {
	return gp.NewAx3(
		p.Location.ToPoint(),
		gp.NewDirVec(p.Normal.ToVector()),
		gp.NewDirVec(p.X.ToVector()),
	)
}

func (p *PlaneParameters) Rotated(a ...interface{}) *PlaneParameters {
	argc := len(a)
	if argc < 1 || argc > 3 {
		panic("No match for overloaded function call")
	}
	degrees := 0.0
	axis := NewVectorFromValues(0, 0, 0)
	origin := NewVectorFromValues(0, 0, 0)

	if argc < 2 {
		// calculate axis
		axis = NewVector(p.Normal.ToVector())
	}
	if argc < 3 {
		// use current location as origin
		origin = p.Location
	}
	coordinates := p.Plane()
	coordinates.Rotate(gp.NewAx1(origin.ToPoint(), gp.NewDirVec(axis.ToVector())), utils.ToRadians(degrees))
	return NewPlaneParametersFromCoordinateSystem(coordinates)
}

func (p *PlaneParameters) Translated(dir Vector) *PlaneParameters {
	return NewPlaneParametersFromCoordinateSystem(p.Plane().Translated(dir.ToVector()))
}

func (p *PlaneParameters) ToString() string {
	return fmt.Sprintf("{ location: %s, x: %s, y: %s }\n", p.Location.ToString(), p.Normal.ToString(), p.X.ToString())
}
