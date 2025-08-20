package sketcher

import (
	"fmt"

	"github.com/marcuswu/dlineate"
	"github.com/marcuswu/dlineate/utils"
	"github.com/marcuswu/gooccwrapper/gp"
)

type Point struct {
	dlineate.Element
	solver         SketchSolver
	X              float64
	Y              float64
	isConstruction bool
	converted      gp.Pnt
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

// IsConstruction returns whether this element is construction geometry
func (p *Point) IsConstruction() bool {
	return p.isConstruction
}

// SetConstruction sets the element as construction geometry
func (p *Point) SetConstruction(isConstruction bool) {
	p.isConstruction = isConstruction
}

func (p *Point) getElement() *dlineate.Element {
	return &p.Element
}

// Coincident creates a constraint placing this point on the specified entity
func (p *Point) Coincident(other Entity) *Point {
	p.solver.Coincident(p, other)

	return p
}

// Horizontal creates a constraint placing this point along the X axis from the specified point
func (p *Point) Horizontal(p2 *Point) *Point {
	p.solver.HorizontalPoints(p, p2)
	return p
}

// Vertical creates a constraint placing this point along the Y axis from the specified point
func (p *Point) Vertical(p2 *Point) *Point {
	p.solver.VerticalPoints(p, p2)
	return p
}

// Distance creates a constraint placing this point the specified distance away from the provided entity
func (p *Point) Distance(other Entity, distance float64) *Point {
	p.solver.Distance(p, other, distance)

	return p
}

// Horizontal Distance creates a constraint placing this point the specified distance along the X axis away from the provided entity
func (p *Point) HorizontalDistance(other Entity, distance float64) *Point {
	p.solver.PointHorizontalDistance(p, other, distance)

	return p
}

// Horizontal Distance creates a constraint placing this point the specified distance along the Y axis away from the provided entity
func (p *Point) VerticalDistance(other Entity, distance float64) *Point {
	p.solver.PointVerticalDistance(p, other, distance)

	return p
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

// Convert translates this point to an OpenCascade element
func (p *Point) Convert() gp.Pnt {
	if p.converted.Pnt == nil {
		p.converted = gp.NewPnt(p.X, p.Y, 0).Transformed(p.solver.Transform())
	}
	return p.converted
}

// UpdateFromValues updates the element's center, start, and end based on the current sketch values
// Automatically called when the sketch is solved
func (p *Point) UpdateFromValues() {
	values := p.Element.Values()
	p.X = values[0]
	p.Y = values[1]
}

// MakeEdge generates an edge from the sketch element. Usually this is handled by MakerCad.
func (p *Point) MakeEdge() *Edge {
	return nil
}

// IsDistanceFrom returns whether or not this point is a specified distance from another point
func (p *Point) IsDistanceFrom(other *Point, dist float64) bool {
	return utils.StandardFloatCompare(p.DistanceBetweenPoints(&other.Element), dist) == 0
}

// IsConnectedTo returns whether this entity is connected to the supplied entity
func (p *Point) IsConnectedTo(other Entity) bool {
	switch o := other.(type) {
	case *Point:
		return utils.StandardFloatCompare(p.X, o.X) == 0 && utils.StandardFloatCompare(p.Y, o.Y) == 0
	case *Line:
		return o.IsConnectedTo(p)
	case *Arc:
		return o.IsConnectedTo(p)
	case *Circle:
		return o.IsConnectedTo(p)
	}
	return false
}
