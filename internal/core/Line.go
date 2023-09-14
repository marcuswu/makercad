package core

import (
	"fmt"

	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/geom"

	"github.com/marcuswu/dlineate"
)

type Line struct {
	dlineate.Element
	solver         SketchSolver
	Start          *Point
	End            *Point
	IsConstruction bool
}

func (l *Line) isConstruction() bool {
	return l.IsConstruction
}

func (l *Line) getElement() *dlineate.Element {
	return &l.Element
}

func (l *Line) Horizontal() *Line {
	l.solver.HorizontalLine(l)

	return l
}

func (l *Line) Vertical() *Line {
	l.solver.VerticalLine(l)

	return l
}

func (l *Line) Length(length float64) *Line {
	l.solver.Distance(l.Start, l.End, length)

	return l
}

func (l *Line) Midpoint(point *Point) *Line {
	l.solver.LineMidpoint(l, point)

	return l
}

func (l *Line) Angle(other *Line, angle float64) *Line {
	l.solver.LineAngle(l, other, angle)

	return l
}

func (l *Line) MakeEdge() *Edge {
	segment := geom.MakeSegment(l.Start.Convert(), l.End.Convert())
	return &Edge{brepbuilderapi.NewMakeEdge(segment).ToTopoDSEdge()}
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v to %v", l.Start, l.End)
}
