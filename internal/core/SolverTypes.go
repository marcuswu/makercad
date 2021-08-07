package core

import "libmakercad/third_party/planegcs"

type Point struct {
	planegcs.Point
	solver         SketchSolver
	x              float64
	y              float64
	isConstruction bool
}

func NewPoint(a ...interface{}) *Point {
	argc := len(a)
	if argc != 1 || argc != 3 {
		panic("Wrong number of parameters to NewPoint")
	}

	solver, ok := a[0].(SketchSolver)
	if !ok {
		panic("First parameter to NewPoint must be SketchSolver")
	}

	var x float64 = 0
	var y float64 = 0

	if argc != 1 {
		x = a[1].(float64)
		y = a[2].(float64)
	}

	p := &Point{solver: solver, x: x, y: y}
	p.Point = planegcs.NewPoint(p.x, p.y)
	return p
}

func (p *Point) GetX() float64 {
	return p.x
}

func (p *Point) GetY() float64 {
	return p.y
}

func (p *Point) SetX(x float64) {
	p.x = x
}

func (p *Point) SetY(y float64) {
	p.y = y
}

func (p *Point) Coincident(a ...interface{}) *Point {
	argc := len(a)
	if argc < 1 || argc > 1 {
		panic("Wrong number of parameters to Point.Coincident")
	}

	switch v := a[0].(type) {
	case Point:
		p.solver.system.AddConstraintP2PCoincident(p, v)
	case Line:
		p.solver.system.AddConstraintPointOnLine(p, v)
	case Arc:
		p.solver.system.AddConstraintPointOnArc(p, v)
	case Circle:
		p.solver.system.AddConstraintPointOnCircle(p, v)
	default:
		panic("Wrong type sent to Point.Coincident")
	}

	return p
}

func (p *Point) Horizontal(p2 Point) *Point {
	p.solver.system.AddConstraintHorizontal(p, p2)
	return p
}

func (p *Point) Vertical(p2 Point) *Point {
	p.solver.system.AddConstraintVertical(p, p2)
	return p
}

func (p *Point) Distance(a ...interface{}) *Point {
	argc := len(a)
	if argc < 2 || argc > 2 {
		panic("Wrong number of parameters to Point.Distance")
	}

	dist, ok := a[1].(float64)
	if !ok {
		panic("Wrong parameter type for distance in Point.Distance")
	}

	switch v := a[0].(type) {
	case Point:
		p.solver.system.AddConstraintP2PDistance(p, v, dist)
	case Line:
		p.solver.system.AddConstraintP2LDistance(p, v, dist)
	default:
		panic("Wrong type sent to Point.Distance")
	}

	return p
}

func (p *Point) Construction(c bool) *Point {
	p.isConstruction = c
	return p
}

func (p *Point) ToString() string {

}
