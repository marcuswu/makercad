package sketch

import (
	"fmt"

	"github.com/marcuswu/dlineate"
	"github.com/marcuswu/gooccwrapper/gp"
)

type Point struct {
	dlineate.Element
	solver         SketchSolver
	X              float64
	Y              float64
	isConstruction bool
}

func (p *Point) GetX() float64 {
	return p.X
}

func (p *Point) GetY() float64 {
	return p.Y
}

func (p *Point) GetZ() float64 {
	return 0.0
}

func (p *Point) SetX(x float64) {
	p.X = x
}

func (p *Point) SetY(y float64) {
	p.Y = y
}

func (p *Point) IsConstruction() bool {
	return p.isConstruction
}

func (p *Point) getElement() *dlineate.Element {
	return &p.Element
}

func (p *Point) Coincident(other Entity) *Point {
	p.solver.Coincident(p, other)

	return p
}

func (p *Point) Horizontal(p2 *Point) *Point {
	p.solver.HorizontalPoints(p, p2)
	return p
}

func (p *Point) Vertical(p2 *Point) *Point {
	p.solver.VerticalPoints(p, p2)
	return p
}

func (p *Point) Distance(other Entity, distance float64) *Point {
	p.solver.Distance(p, other, distance)

	return p
}

func (p *Point) HorizontalDistance(other Entity, distance float64) *Point {
	p.solver.PointHorizontalDistance(p, other, distance)

	return p
}

func (p *Point) VerticalDistance(other Entity, distance float64) *Point {
	p.solver.PointVerticalDistance(p, other, distance)

	return p
}

func (p *Point) ToString() string {
	return fmt.Sprintf("%f, %f", p.X, p.Y)
}

func (p *Point) Convert() gp.Pnt {
	return gp.NewPnt(p.X, p.Y, 0).Transformed(p.solver.Transform())
}

func (p *Point) MakeEdge() *Edge {
	return nil
}
