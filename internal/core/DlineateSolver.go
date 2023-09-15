package core

import (
	"github.com/marcuswu/dlineate"
	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/gp"
	"github.com/marcuswu/gooccwrapper/topods"
)

type DlineateSolver struct {
	system           *dlineate.Sketch
	entities         []Entity
	coordinateSystem gp.Ax3
}

func NewDlineateSolverFromPlane(plane PlaneParameters) *DlineateSolver {
	return &DlineateSolver{dlineate.NewSketch(), make([]Entity, 0), plane.ToAx3()}
}

func NewDlineateSolverFromFace(face *Face) *DlineateSolver {
	return &DlineateSolver{dlineate.NewSketch(), make([]Entity, 0), face.Plane()}
}

func (s *DlineateSolver) CreatePoint(x float64, y float64) *Point {
	entity := &Point{Element: *s.system.AddPoint(x, y), solver: s, X: x, Y: y, IsConstruction: false}
	s.entities = append(s.entities, entity)
	return entity
}

func (s *DlineateSolver) CreateLine(p1 *Point, p2 *Point) *Line {
	entity := &Line{Element: *s.system.AddLine(p1.X, p1.Y, p2.X, p2.Y), solver: s, Start: p1, End: p2, IsConstruction: false}
	s.entities = append(s.entities, entity)
	return entity
}

func (s *DlineateSolver) CreateCircle(p *Point, r float64) *Circle {
	entity := &Circle{Element: *s.system.AddCircle(p.X, p.Y, r), solver: s, Center: p, Radius: r, IsConstruction: false}
	s.entities = append(s.entities, entity)
	return entity
}

func (s *DlineateSolver) CreateArc(center *Point, start *Point, end *Point) *Arc {
	entity := &Arc{
		Element:        *s.system.AddArc(center.X, center.Y, start.X, start.Y, end.X, end.Y),
		solver:         s,
		Center:         center,
		Start:          start,
		End:            end,
		IsConstruction: false,
	}
	s.entities = append(s.entities, entity)
	return entity
}

/*
func (s *DlineateSolver) CreateDistance(float64) *Distance {
}

func (s *DlineateSolver) CreateWorkplanePoint(x float64, y float64) *Point {
}

func (s *DlineateSolver) CreateWorkplaneLine(*Point, *Point) *Line {
}

func (s *DlineateSolver) CreateWorkplaneCircle(*Point, float64) *Circle {
}

func (s *DlineateSolver) CreateWorkplaneArc(*Point, *Point, *Point) *Arc {
}
*/

func (s *DlineateSolver) Coincident(e1 Entity, e2 Entity) {
	_, isE1Point := e1.(*Point)
	_, isE2Point := e2.(*Point)
	if isE1Point && isE2Point {
		s.system.AddCoincidentConstraint(e1.getElement(), e2.getElement())
		return
	}

	s.system.AddDistanceConstraint(e1.getElement(), e2.getElement(), 0)
}

func (s *DlineateSolver) PointVerticalDistance(p *Point, e Entity, d float64) {
	pe, ok := e.(*Point)
	if !ok {
		pe = s.CreatePoint(0, 0)
		pe.IsConstruction = true
	}
	cl := s.CreateLine(p, pe)
	cl.IsConstruction = true
	s.system.AddVerticalConstraint(cl.getElement())
	s.system.AddDistanceConstraint(pe.getElement(), e.getElement(), 0)
	s.system.AddDistanceConstraint(p.getElement(), e.getElement(), d)
}

func (s *DlineateSolver) PointHorizontalDistance(p *Point, e Entity, d float64) {
	pe, ok := e.(*Point)
	if !ok {
		pe = s.CreatePoint(0, 0)
		pe.IsConstruction = true
	}
	cl := s.CreateLine(p, pe)
	cl.IsConstruction = true
	s.system.AddHorizontalConstraint(cl.getElement())
	s.system.AddDistanceConstraint(pe.getElement(), e.getElement(), 0)
	s.system.AddDistanceConstraint(p.getElement(), e.getElement(), d)
}

func (s *DlineateSolver) PointProjectedDistance(p *Point, e Entity, d float64) {
	pe, ok := e.(*Point)
	if !ok {
		pe = s.CreatePoint(0, 0)
		pe.IsConstruction = true
	}
	cl := s.CreateLine(p, pe)
	cl.IsConstruction = true
	s.system.AddDistanceConstraint(pe.getElement(), e.getElement(), 0)
	s.system.AddPerpendicularConstraint(cl.getElement(), e.getElement())
	s.system.AddDistanceConstraint(p.getElement(), e.getElement(), d)
}

func (s *DlineateSolver) LineMidpoint(l *Line, e Entity) {
	s.system.AddMidpointConstraint(e.getElement(), l.getElement())
}

func (s *DlineateSolver) LineAngle(l1 *Line, l2 *Line, d float64) {
	s.system.AddAngleConstraint(l1.getElement(), l2.getElement(), d, false)
}

func (s *DlineateSolver) ArcLineTangent(a *Arc, l *Line) {
	s.system.AddTangentConstraint(a.getElement(), l.getElement())
}

func (s *DlineateSolver) Distance(e1 Entity, e2 Entity, d float64) {
	s.system.AddDistanceConstraint(e1.getElement(), e2.getElement(), d)
}

func (s *DlineateSolver) HorizontalLine(l *Line) {
	s.system.AddHorizontalConstraint(l.getElement())
}

func (s *DlineateSolver) HorizontalPoints(p1 *Point, p2 *Point) {
	hl := s.CreateLine(p1, p2)
	hl.IsConstruction = true
	s.system.AddHorizontalConstraint(hl.getElement())
}

func (s *DlineateSolver) VerticalLine(l *Line) {
	s.system.AddVerticalConstraint(l.getElement())
}

func (s *DlineateSolver) VerticalPoints(p1 *Point, p2 *Point) {
	hl := s.CreateLine(p1, p2)
	hl.IsConstruction = true
	s.system.AddVerticalConstraint(hl.getElement())
}

func (s *DlineateSolver) LineLength(l *Line, d float64) {
	s.system.AddDistanceConstraint(l.Start.getElement(), l.End.getElement(), d)
}

func (s *DlineateSolver) Equal(e1 Entity, e2 Entity) {
	s.system.AddEqualConstraint(e1.getElement(), e2.getElement())
}

func (s *DlineateSolver) CurveDiameter(e Entity, d float64) {
	if a, ok := e.(*Arc); ok {
		s.system.AddDistanceConstraint(a.Center.getElement(), a.Start.getElement(), d)
		return
	}
	c, ok := e.(*Circle)
	if !ok {
		return
	}
	s.system.AddDistanceConstraint(c.Center.getElement(), nil, d)
}

func (s *DlineateSolver) CoordinateSystem() gp.Ax3 {
	return s.coordinateSystem
}

func (s *DlineateSolver) Transform() gp.Trsf {
	defaultCoords := gp.NewAx3(gp.NewPnt(0, 0, 0), gp.NewDir(0, 0, 1), gp.NewDir(1, 0, 0))
	transform := gp.NewTrsf()
	transform.SetTransformation(s.coordinateSystem, defaultCoords)
	return transform
}

func (s *DlineateSolver) Solve() {
	s.system.Solve()
}

func (s *DlineateSolver) ToFace() *Face {
	wires := make([]topods.Wire, 0)
	for i := range s.entities {
		entity := s.entities[i]
		if entity.isConstruction() {
			continue
		}
		wires = append(wires, brepbuilderapi.NewMakeWireWithEdge(entity.MakeEdge().edge).ToTopoDSWire())
	}

	combined := brepbuilderapi.NewMakeWire()
	for i := range wires {
		combined.AddWire(wires[i])
	}

	return &Face{brepbuilderapi.NewMakeFace(combined.ToTopoDSWire()).ToTopoDSFace()}
}
