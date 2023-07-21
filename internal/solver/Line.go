package solver

import (
	"fmt"

	"github.com/marcuswu/dlineate"
)

type Line struct {
	dlineate.Element
	solver         SketchSolver
	start          *Point
	end            *Point
	isConstruction bool
}

func (l *Line) getElement() *dlineate.Element {
	return &l.Element
}

func (l *Line) Start() *Point {
	return l.start
}

func (l *Line) End() *Point {
	return l.end
}

func (l *Line) IsConstruction() bool {
	return l.isConstruction
}

func (l *Line) Horizontal() *Line {
	l.solver.Horizontal(l.Start(), l.End())

	return l
}

func (l *Line) Vertical() *Line {
	l.solver.Vertical(l.Start(), l.End())

	return l
}

func (l *Line) Length(length float64) *Line {
	l.solver.Distance(l.Start(), l.End(), length)

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

// func (l *Line) MakeEdge() occt.TopoDS_Edge {
// 	return nil
// }

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v to %v", l.Start(), l.End())
}
