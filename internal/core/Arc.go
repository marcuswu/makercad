package core

import (
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
	IsConstruction bool
}

func (a *Arc) isConstruction() bool {
	return a.IsConstruction
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
	// centerPnt := a.Center.Convert()
	centerPoint := gp.NewPnt(a.Center.X, a.Center.Y, 0.0)
	center := gp.NewAx2(centerPoint, a.solver.CoordinateSystem().Direction())
	start := gp.NewPnt(a.Start.X, a.Start.Y, 0.0)
	end := gp.NewPnt(a.End.X, a.End.Y, 0.0)
	radius := gp.NewVecPoints(start, end).Magnitude()
	circle := gp.NewCirc(center, radius)
	arc := geom.MakeArc(circle, start, end, true)
	return &Edge{brepbuilderapi.NewMakeEdge(arc).ToTopoDSEdge()}
}
