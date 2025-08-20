package sketcher

import (
	"fmt"

	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/geom"
	"github.com/rs/zerolog/log"

	"github.com/marcuswu/dlineate"
)

// Line represents a 2D line on a sketch plane
type Line struct {
	dlineate.Element
	solver         SketchSolver
	Start          *Point
	End            *Point
	isConstruction bool
}

// IsConstruction returns whether this element is construction geometry
func (l *Line) IsConstruction() bool {
	return l.isConstruction
}

// SetConstruction sets the element as construction geometry
func (l *Line) SetConstruction(isConstruction bool) {
	l.isConstruction = isConstruction
}

func (l *Line) getElement() *dlineate.Element {
	return &l.Element
}

// Horizontal creates a constraint specifying the line's angle to parallel with the sketch's X axis
func (l *Line) Horizontal() *Line {
	l.solver.HorizontalLine(l)

	return l
}

// Vertical creates a constraint specifying the line's angle to parallel with the sketch's Y axis
func (l *Line) Vertical() *Line {
	l.solver.VerticalLine(l)

	return l
}

// Length creates a constraint specifying the distance between its start and end
func (l *Line) Length(length float64) *Line {
	if l.End == nil || l.Start == nil {
		return l
	}
	l.solver.Distance(l.Start, l.End, length)

	return l
}

// Midpoint creates a constraint ensuring that the specified point lies halfway between the line's start and end.
func (l *Line) Midpoint(point *Point) *Line {
	l.solver.LineMidpoint(l, point)

	return l
}

// Angle creates a constraint ensuring the angle between this line and the provided line is at the specified angle in radians.
func (l *Line) Angle(other *Line, angle float64) *Line {
	l.solver.LineAngle(l, other, angle)

	return l
}

// MakeEdge generates an edge from the sketch element. Usually this is handled by MakerCad.
func (l *Line) MakeEdge() *Edge {
	if l.Start.ID() == l.End.ID() {
		return nil
	}
	start, end := l.Start.Convert(), l.End.Convert()
	if start.Distance(end) == 0 {
		return nil
	}
	log.Debug().Str("Line", l.String()).Msg("Making edge")
	segment := geom.MakeSegment(start, end)
	return &Edge{brepbuilderapi.NewMakeEdge(segment).ToTopoDSEdge()}
}

// UpdateFromValues updates the element's center, start, and end based on the current sketch values
// Automatically called when the sketch is solved
func (l *Line) UpdateFromValues() {
	l.Start.UpdateFromValues()
	l.End.UpdateFromValues()
}

func (l *Line) String() string {
	return fmt.Sprintf("%v to %v", l.Start.String(), l.End.String())
}

// IsConnectedTo returns whether this entity is connected to the supplied entity
func (l *Line) IsConnectedTo(other Entity) bool {
	switch o := other.(type) {
	case *Point:
		return o.IsConnectedTo(l.Start) || o.IsConnectedTo(l.End)
	case *Arc:
		return o.Start.IsConnectedTo(l.Start) ||
			o.Start.IsConnectedTo(l.End) ||
			o.End.IsConnectedTo(l.Start) ||
			o.End.IsConnectedTo(l.End)
	case *Line:
		return o.Start.IsConnectedTo(l.Start) ||
			o.Start.IsConnectedTo(l.End) ||
			o.End.IsConnectedTo(l.Start) ||
			o.End.IsConnectedTo(l.End)
	case *Circle:
		return o.IsConnectedTo(l)
	}
	return false
}
