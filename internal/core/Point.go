package core

import (
	"fmt"

	"github.com/marcuswu/dlineate"
)

type Point struct {
	dlineate.Element
	solver         SketchSolver
	X              float64
	Y              float64
	IsConstruction bool
}

func (p *Point) getElement() *dlineate.Element {
	return &p.Element
}

func (p *Point) Coincident(other Entity) *Point {
	p.solver.Coincident(p, other)

	return p
}

func (p *Point) Horizontal(p2 *Point) *Point {
	p.solver.Horizontal(p, p2)
	return p
}

func (p *Point) Vertical(p2 *Point) *Point {
	p.solver.Vertical(p, p2)
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
