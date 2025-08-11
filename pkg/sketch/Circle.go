package sketch

import (
	"log"

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
	isConstruction bool
}

func (c *Circle) IsConstruction() bool {
	return c.isConstruction
}

func (c *Circle) SetConstruction(isConstruction bool) {
	c.isConstruction = isConstruction
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

func (c *Circle) UpdateFromValues() {
	values := c.Element.Values()
	c.Center.UpdateFromValues()
	c.Radius = values[2]
}

func (c *Circle) MakeEdge() *Edge {
	log.Printf("Making edge from circle %s\n", c.String())
	centerPoint := gp.NewPnt(c.Center.X, c.Center.Y, 0.0)
	radius := c.Radius
	center := gp.NewAx2(centerPoint, c.solver.CoordinateSystem().Direction(), c.solver.CoordinateSystem().XDirection())
	circle := geom.MakeCircle(center, radius)
	return &Edge{brepbuilderapi.NewMakeEdge(circle).ToTopoDSEdge()}
}

func (c *Circle) IsConnectedTo(other Entity) bool {
	switch o := other.(type) {
	case *Point:
		return c.Center.IsDistanceFrom(o, c.Radius)
	case *Line:
		return c.Center.IsDistanceFrom(o.Start, c.Radius) || c.Center.IsDistanceFrom(o.End, c.Radius)
	case *Arc:
		return c.Center.IsDistanceFrom(o.Start, c.Radius) || c.Center.IsDistanceFrom(o.End, c.Radius)
	case *Circle:
		return c.Center.IsDistanceFrom(o.Center, c.Radius+o.Radius)
	}
	return false
}
