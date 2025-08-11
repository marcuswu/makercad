package sketch

import (
	"fmt"

	"github.com/marcuswu/dlineate"
	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/geom"
	"github.com/marcuswu/gooccwrapper/gp"
)

type Arc struct {
	dlineate.Element
	solver         SketchSolver
	Center         *Point
	Start          *Point
	End            *Point
	isConstruction bool
}

func (a *Arc) IsConstruction() bool {
	return a.isConstruction
}

func (a *Arc) SetConstruction(isConstruction bool) {
	a.isConstruction = isConstruction
}

func (a *Arc) getElement() *dlineate.Element {
	return &a.Element
}

func (a *Arc) Diameter(d float64) *Arc {
	a.solver.CurveDiameter(a, d)

	return a
}

func (a *Arc) Tangent(l *Line) *Arc {
	a.solver.ArcLineTangent(a, l)

	return a
}

func (a *Arc) MakeEdge() *Edge {
	centerPoint := a.Center.Convert()
	normalDir := a.solver.CoordinateSystem().Direction()
	xDir := a.solver.CoordinateSystem().XDirection()
	center := gp.NewAx2(centerPoint, normalDir, xDir)
	start := a.Start.Convert()
	end := a.End.Convert()
	radius := gp.NewVecPoints(centerPoint, start).Magnitude()
	circle := gp.NewCirc(center, radius)
	arc := geom.MakeArc(circle, start, end, true)
	return &Edge{brepbuilderapi.NewMakeEdge(arc).ToTopoDSEdge()}
}

func (a *Arc) UpdateFromValues() {
	a.Center.UpdateFromValues()
	a.Start.UpdateFromValues()
	a.End.UpdateFromValues()
}

func (a *Arc) String() string {
	return fmt.Sprintf("%v to %v around %v", a.Start.String(), a.End.String(), a.Center.String())
}

func (a *Arc) IsConnectedTo(other Entity) bool {
	switch o := other.(type) {
	case *Point:
		return o.IsConnectedTo(a.Start) || o.IsConnectedTo(a.End)
	case *Line:
		return o.Start.IsConnectedTo(a.Start) ||
			o.Start.IsConnectedTo(a.End) ||
			o.End.IsConnectedTo(a.Start) ||
			o.End.IsConnectedTo(a.End)
	case *Arc:
		return o.Start.IsConnectedTo(a.Start) ||
			o.Start.IsConnectedTo(a.End) ||
			o.End.IsConnectedTo(a.Start) ||
			o.End.IsConnectedTo(a.End)
	case *Circle:
		return o.IsConnectedTo(a.Start) || o.IsConnectedTo(a.End)
	}
	return false
}
