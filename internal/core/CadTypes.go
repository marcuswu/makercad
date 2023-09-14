package core

import (
	"fmt"
	"libmakercad/internal/utils"

	"github.com/marcuswu/gooccwrapper/brepalgoapi"
	"github.com/marcuswu/gooccwrapper/gp"
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
	case gp.Pnt:
		return &Vector{v.X(), v.Y(), v.Z()}
	case gp.Vec:
		return &Vector{v.X(), v.Y(), v.Z()}
	case gp.Dir:
		return &Vector{v.X(), v.Y(), v.Z()}
	default:
		panic("Invalid value passed to NewVector")
	}
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

func NewPlaneParameters(a ...interface{}) *PlaneParameters {
	argc := len(a)

	if argc != 3 && argc != 1 && argc != 0 {
		panic("No match for overloaded function call")
	}

	if argc == 0 {
		return &PlaneParameters{NewVector(), NewVector(), NewVector()}
	}

	if argc == 1 {
		coord, ok := a[0].(gp.Ax3)
		if !ok {
			panic("Invalid parameter for NewPlaneParameters")
		}
		return &PlaneParameters{
			NewVector(coord.Location()),
			NewVector(NewVector(coord.YDirection()).ToVector().Crossed(NewVector(coord.XDirection()).ToVector())),
			NewVector(coord.XDirection()),
		}
	}

	loc, ok1 := a[0].(*Vector)
	normalDir, ok2 := a[0].(*Vector)
	yDir, ok3 := a[0].(*Vector)
	if !ok1 || !ok2 || !ok3 {
		panic("Invalid parameters for NewPlaneParameters")
	}
	return &PlaneParameters{loc, normalDir, yDir}
}

func (p *PlaneParameters) ToAx3() gp.Ax3 {
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
	axis := NewVector(0, 0, 0)
	origin := NewVector(0, 0, 0)

	if argc < 2 {
		// calculate axis
		axis = NewVector(p.Normal.ToVector())
	}
	if argc < 3 {
		// use current location as origin
		origin = p.Location
	}
	coordinates := p.ToAx3()
	coordinates.Rotate(gp.NewAx1(origin.ToPoint(), gp.NewDirVec(axis.ToVector())), utils.ToRadians(degrees))
	return NewPlaneParameters(coordinates)
}

func (p *PlaneParameters) Translated(dir Vector) *PlaneParameters {
	return NewPlaneParameters(p.ToAx3().Translated(dir.ToVector()))
}

func (p *PlaneParameters) ToString() string {
	return fmt.Sprintf("{ location: %s, x: %s, y: %s }\n", p.Location.ToString(), p.Normal.ToString(), p.X.ToString())
}

type MergeType int

const (
	MergeTypeNew MergeType = iota
	MergeTypeAdd
	MergeTypeRemove
	MergeTypeMax
)

type CadOperation struct {
	shape     Shape
	operation *brepalgoapi.Boolean
}

func NewCadOperation(e Shape, op *brepalgoapi.Boolean) *CadOperation {
	return &CadOperation{shape: e, operation: op}
}

func (o *CadOperation) Shape() Shape {
	shape := o.shape
	if o.operation != nil {
		shape = Shape{o.operation.Shape()}
	}

	return shape
}
