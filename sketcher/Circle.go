package sketcher

import (
	"log"

	"github.com/marcuswu/dlineate"
	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/geom"
	"github.com/marcuswu/gooccwrapper/gp"
)

// Circle represents a 2D circle on a sketch plane
type Circle struct {
	dlineate.Element
	solver         SketchSolver
	Center         *Point
	Radius         float64
	isConstruction bool
}

// IsConstruction returns whether this element is construction geometry
func (c *Circle) IsConstruction() bool {
	return c.isConstruction
}

// SetConstruction sets the element as construction geometry
func (c *Circle) SetConstruction(isConstruction bool) {
	c.isConstruction = isConstruction
}

func (c *Circle) getElement() *dlineate.Element {
	return &c.Element
}

// Diameter creates a constraint specifying the circle's diameter
func (c *Circle) Diameter(d float64) *Circle {
	c.solver.CurveDiameter(c, d)

	return c
}

// Equal returns whether this element is equal to another
func (c *Circle) Equal(other *Circle) *Circle {
	c.solver.Equal(c, other)

	return c
}

// UpdateFromValues updates the element's center, start, and end based on the current sketch values
// Automatically called when the sketch is solved
func (c *Circle) UpdateFromValues() {
	values := c.Element.Values()
	c.Center.UpdateFromValues()
	c.Radius = values[2]
}

// MakeEdge generates an edge from the sketch element. Usually this is handled by MakerCad.
func (c *Circle) MakeEdge() *Edge {
	log.Printf("Making edge from circle %s\n", c.String())
	centerPoint := gp.NewPnt(c.Center.X, c.Center.Y, 0.0)
	radius := c.Radius
	center := gp.NewAx2(centerPoint, c.solver.CoordinateSystem().Direction(), c.solver.CoordinateSystem().XDirection())
	circle := geom.MakeCircle(center, radius)
	return &Edge{brepbuilderapi.NewMakeEdge(circle).ToTopoDSEdge()}
}

// IsConnectedTo returns whether this entity is connected to the supplied entity
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
