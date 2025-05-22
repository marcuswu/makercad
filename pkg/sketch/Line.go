package sketch

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
	isConstruction bool
}

func (l *Line) IsConstruction() bool {
	return l.isConstruction
}

func (l *Line) SetConstruction(isConstruction bool) {
	l.isConstruction = isConstruction
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
	if l.End == nil || l.Start == nil {
		return l
	}
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

func (l *Line) UpdateFromValues() {
	l.Start.UpdateFromValues()
	l.End.UpdateFromValues()
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v to %v", l.Start.ToString(), l.End.ToString())
}
