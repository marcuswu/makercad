package solver

import (
	"github.com/marcuswu/dlineate"
)

type Circle struct {
	dlineate.Element
	solver         SketchSolver
	center         *Point
	radius         float64
	isConstruction bool
}

func (c *Circle) getElement() *dlineate.Element {
	return &c.Element
}

func (c *Circle) Center() *Point {
	return c.center
}

func (c *Circle) Radius() float64 {
	return c.radius
}

func (c *Circle) IsConstruction() bool {
	return c.isConstruction
}
