package sketch

import (
	"github.com/marcuswu/dlineate"
	"github.com/marcuswu/gooccwrapper/gp"
)

type DlineateSolver struct {
	system           *dlineate.Sketch
	entities         []Entity
	coordinateSystem gp.Ax3
	origin           *Point
	XAxis            *Line
	YAxis            *Line
}

func NewDlineateSolver(planer Planer) *DlineateSolver {
	solver := &DlineateSolver{dlineate.NewSketch(), make([]Entity, 0), planer.Plane(), nil, nil, nil}
	solver.origin = &Point{Element: *solver.system.Origin, solver: solver, X: 0, Y: 0, isConstruction: true}
	solver.XAxis = &Line{Element: *solver.system.XAxis, solver: solver, Start: nil, End: nil, isConstruction: true}
	solver.YAxis = &Line{Element: *solver.system.YAxis, solver: solver, Start: nil, End: nil, isConstruction: true}
	return solver
}

func (s *DlineateSolver) Entities() []Entity {
	return s.entities
}

func (s *DlineateSolver) Origin() *Point {
	return s.origin
}

func (s *DlineateSolver) CreatePoint(x float64, y float64) *Point {
	entity := &Point{Element: *s.system.AddPoint(x, y), solver: s, X: x, Y: y, isConstruction: false}
	s.entities = append(s.entities, entity)
	return entity
}

func (s *DlineateSolver) PointFromRef(ref *dlineate.Element) *Point {
	return &Point{Element: *ref, solver: s, X: ref.Values()[0], Y: ref.Values()[1], isConstruction: false}
}

func (s *DlineateSolver) CreateLine(p1X float64, p1Y float64, p2X float64, p2Y float64) *Line {
	entity := &Line{Element: *s.system.AddLine(p1X, p1Y, p2X, p2Y), solver: s, isConstruction: false}
	entity.Start = s.PointFromRef(entity.Element.Start())
	entity.End = s.PointFromRef(entity.Element.End())
	s.entities = append(s.entities, entity)
	return entity
}

func (s *DlineateSolver) CreateCircle(centerX float64, centerY float64, r float64) *Circle {
	entity := &Circle{Element: *s.system.AddCircle(centerX, centerY, r), solver: s, Radius: r, isConstruction: false}
	entity.Center = s.PointFromRef(entity.Element.Center())
	s.entities = append(s.entities, entity)
	return entity
}

func (s *DlineateSolver) CreateArc(centerX float64, centerY float64, startX float64, startY float64, endX float64, endY float64) *Arc {
	entity := &Arc{
		Element:        *s.system.AddArc(centerX, centerY, startX, startY, endX, endY),
		solver:         s,
		isConstruction: false,
	}
	entity.Start = s.PointFromRef(entity.Element.Start())
	entity.Center = s.PointFromRef(entity.Element.Center())
	entity.End = s.PointFromRef(entity.Element.End())
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
		s.system.AddCoincidentConstraint(pe.getElement(), e.getElement())
		pe.isConstruction = true
	}
	cl := s.CreateLine(p.X, p.Y, pe.X, pe.Y)
	cl.isConstruction = true
	s.system.AddCoincidentConstraint(p.getElement(), cl.getElement())
	s.system.AddCoincidentConstraint(pe.getElement(), cl.getElement())
	s.system.AddVerticalConstraint(cl.getElement())
	s.system.AddDistanceConstraint(p.getElement(), e.getElement(), d)
}

func (s *DlineateSolver) PointHorizontalDistance(p *Point, e Entity, d float64) {
	pe, ok := e.(*Point)
	if !ok {
		pe = s.CreatePoint(0, 0)
		s.system.AddCoincidentConstraint(pe.getElement(), e.getElement())
		pe.isConstruction = true
	}
	cl := s.CreateLine(p.X, p.Y, pe.X, pe.Y)
	cl.isConstruction = true
	s.system.AddCoincidentConstraint(p.getElement(), cl.getElement())
	s.system.AddCoincidentConstraint(pe.getElement(), cl.getElement())
	s.system.AddHorizontalConstraint(cl.getElement())
	s.system.AddDistanceConstraint(p.getElement(), e.getElement(), d)
}

func (s *DlineateSolver) PointProjectedDistance(p *Point, e Entity, d float64) {
	pe, ok := e.(*Point)
	if !ok {
		pe = s.CreatePoint(0, 0)
		s.system.AddCoincidentConstraint(pe.getElement(), e.getElement())
		pe.isConstruction = true
	}
	cl := s.CreateLine(p.X, p.Y, pe.X, pe.Y)
	cl.isConstruction = true
	s.system.AddCoincidentConstraint(p.getElement(), cl.getElement())
	s.system.AddCoincidentConstraint(pe.getElement(), cl.getElement())
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
	hl := s.CreateLine(p1.X, p1.Y, p2.X, p2.Y)
	hl.isConstruction = true
	s.system.AddCoincidentConstraint(hl.getElement(), p1.getElement())
	s.system.AddCoincidentConstraint(hl.getElement(), p2.getElement())
	s.system.AddHorizontalConstraint(hl.getElement())
}

func (s *DlineateSolver) VerticalLine(l *Line) {
	s.system.AddVerticalConstraint(l.getElement())
}

func (s *DlineateSolver) VerticalPoints(p1 *Point, p2 *Point) {
	vl := s.CreateLine(p1.X, p1.Y, p2.X, p2.Y)
	vl.isConstruction = true
	s.system.AddCoincidentConstraint(vl.getElement(), p1.getElement())
	s.system.AddCoincidentConstraint(vl.getElement(), p2.getElement())
	s.system.AddVerticalConstraint(vl.getElement())
}

func (s *DlineateSolver) LineLength(l *Line, d float64) {
	s.system.AddDistanceConstraint(l.Start.getElement(), l.End.getElement(), d)
}

func (s *DlineateSolver) Equal(e1 Entity, e2 Entity) {
	s.system.AddEqualConstraint(e1.getElement(), e2.getElement())
}

func (s *DlineateSolver) CurveDiameter(e Entity, d float64) {
	if a, ok := e.(*Arc); ok {
		s.system.AddDistanceConstraint(a.Center.getElement(), a.Start.getElement(), d/2)
		return
	}
	c, ok := e.(*Circle)
	if !ok {
		return
	}
	// This requires that the circle
	s.system.AddDistanceConstraint(c.getElement(), nil, d/2)
}

func (s *DlineateSolver) CoordinateSystem() gp.Ax3 {
	return s.coordinateSystem
}

func (s *DlineateSolver) MakeFixed(e Entity) {
	s.system.MakeFixed(e.getElement())
}

func (s *DlineateSolver) Transform() gp.Trsf {
	defaultCoords := gp.NewAx3(gp.NewPnt(0, 0, 0), gp.NewDir(0, 0, 1), gp.NewDir(1, 0, 0))
	transform := gp.NewTrsf()
	transform.SetTransformation(s.coordinateSystem, defaultCoords)
	return transform
}

func (s *DlineateSolver) Solve() error {
	err := s.system.Solve()
	for _, e := range s.entities {
		e.UpdateFromValues()
	}
	return err
}

func (s *DlineateSolver) OverConstrained() []string {
	constraints := s.system.ConflictingConstraints()
	ret := make([]string, 0, len(constraints))
	for _, c := range constraints {
		ret = append(ret, c.String())
	}
	return ret
}

func (s *DlineateSolver) LogDebug(file string) error {
	return s.system.ExportGraphViz(file)
}
