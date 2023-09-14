package core

import (
	"github.com/marcuswu/dlineate"
	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/geom"
	"github.com/marcuswu/gooccwrapper/gp"
)

type Circle struct {
	dlineate.Element
	solver         SketchSolver
	Center         *Point
	Radius         float64
	IsConstruction bool
}

func (c *Circle) isConstruction() bool {
	return c.IsConstruction
}

func (c *Circle) getElement() *dlineate.Element {
	return &c.Element
}

func (c *Circle) Diameter(d float64) *Circle {
	c.solver.CurveDiameter(c, d)

	return c
}

func (c *Circle) Equal(other *Circle) *Circle {
	c.solver.Equal(c, other)

	return c
}

func (c *Circle) MakeEdge() *Edge {
	center := gp.NewPnt(c.Center.X, c.Center.Y, 0.0)
	radius := c.Radius
	circle := geom.MakeCircle(center, radius)
	return &Edge{brepbuilderapi.NewMakeEdge(circle).ToTopoDSEdge()}
}
