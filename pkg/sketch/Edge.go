package sketch

import (
	"slices"

	"github.com/marcuswu/gooccwrapper/brepadapter"
	"github.com/marcuswu/gooccwrapper/breptool"
	"github.com/marcuswu/gooccwrapper/gp"
	"github.com/marcuswu/gooccwrapper/topexp"
	"github.com/marcuswu/gooccwrapper/topods"
)

type Edge struct {
	Edge topods.Edge
}
type ListOfEdge []*Edge

type EdgeFilter func(*Edge) bool
type EdgeSorter func(a, b *Edge) int

func (l ListOfEdge) First(filter EdgeFilter) *Edge {
	for _, edge := range l {
		if filter(edge) {
			return edge
		}
	}
	return nil
}

func (l ListOfEdge) Matching(filter EdgeFilter) ListOfEdge {
	newList := make(ListOfEdge, 0, len(l))
	for _, edge := range l {
		if filter(edge) {
			newList = append(newList, edge)
		}
	}
	return newList
}

func (l ListOfEdge) Sort(sorter EdgeSorter) {
	slices.SortFunc(l, sorter)
}

func NewEdgeFromRef(shape topods.Shape) *Edge {
	return &Edge{topods.NewEdgeFromRef(topods.TopoDSEdge(shape.Shape))}
}

func (e *Edge) IsLine() bool {
	curve := brepadapter.NewCurve(e.Edge)
	return curve.IsLine()
}

func (e *Edge) IsCircle() bool {
	curve := brepadapter.NewCurve(e.Edge)
	return curve.IsCircle()
}

func (e *Edge) IsEllipse() bool {
	curve := brepadapter.NewCurve(e.Edge)
	return curve.IsEllipse()
}

func (e *Edge) projectPointToSketch(solver SketchSolver, point gp.Pnt) (float64, float64) {
	origin := solver.CoordinateSystem().Location()
	originDir := gp.NewDir(origin.X(), origin.Y(), origin.Z())
	pointDir := gp.NewDirVec(gp.NewVec(point.X(), point.Y(), point.Z()))
	u := solver.CoordinateSystem().XDirection()
	v := solver.CoordinateSystem().YDirection()
	x := u.Dot(pointDir) - u.Dot(originDir)
	y := v.Dot(pointDir) - u.Dot(originDir)
	return x, y
}

func (e *Edge) GetLine(solver SketchSolver) *Line {
	if !e.IsLine() {
		return nil
	}

	ex := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(&e.Edge)), topexp.Vertex)
	startX, startY := e.projectPointToSketch(solver, breptool.Pnt(topods.NewVertexFromRef(topods.TopoDSVertex(ex.Current().Shape))))
	ex.Next()
	endX, endY := e.projectPointToSketch(solver, breptool.Pnt(topods.NewVertexFromRef(topods.TopoDSVertex(ex.Current().Shape))))

	line := solver.CreateLine(startX, startY, endX, endY)
	solver.MakeFixed(line)
	return line
}

func (e *Edge) LineLength() float64 {
	if !e.IsLine() {
		return 0.0
	}

	ex := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(&e.Edge)), topexp.Vertex)
	start := breptool.Pnt(topods.NewVertexFromRef(topods.TopoDSVertex(ex.Current().Shape)))
	ex.Next()
	end := breptool.Pnt(topods.NewVertexFromRef(topods.TopoDSVertex(ex.Current().Shape)))

	return start.Distance(end)
}

func (e *Edge) GetCircle(solver SketchSolver) *Circle {
	if !e.IsCircle() {
		return nil
	}

	curve := brepadapter.NewCurve(e.Edge)
	circle := curve.ToCircle()
	centerX, centerY := e.projectPointToSketch(solver, circle.Location())
	radius := circle.Radius()

	circ := solver.CreateCircle(centerX, centerY, radius)
	solver.MakeFixed(circ)
	return circ
}

func (e *Edge) CircleRadius() float64 {
	if !e.IsCircle() {
		return 0.0
	}

	curve := brepadapter.NewCurve(e.Edge)
	circle := curve.ToCircle()

	return circle.Radius()
}
