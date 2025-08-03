package sketch

import (
	"slices"

	"github.com/marcuswu/dlineate/utils"
	"github.com/marcuswu/gooccwrapper/brepadapter"
	"github.com/marcuswu/gooccwrapper/gcpnts"
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

func (l ListOfEdge) IsLine() ListOfEdge {
	return l.Matching(func(e *Edge) bool {
		return e.IsLine()
	})
}

func (l ListOfEdge) Length(length float64) ListOfEdge {
	return l.Matching(func(e *Edge) bool {
		return utils.FloatCompare(e.LineLength(), length, utils.StandardCompare) == 0
	})
}

func (l ListOfEdge) Parallel(dir gp.Dir) ListOfEdge {
	return l.Matching(func(e *Edge) bool {
		return e.IsParallel(dir)
	})
}

func (l ListOfEdge) Sort(sorter EdgeSorter) {
	slices.SortFunc(l, sorter)
}

func (l ListOfEdge) SortByLength(inverse bool) {
	l.Sort(func(a, b *Edge) int {
		aL := a.LineLength()
		bL := b.LineLength()
		if inverse {
			return utils.StandardFloatCompare(bL, aL)
		}
		return utils.StandardFloatCompare(aL, bL)
	})
}

func (l ListOfEdge) SortByX(inverse bool) {
	l.Sort(func(a, b *Edge) int {
		aX := a.FirstVertex().X()
		bX := b.FirstVertex().X()
		if inverse {
			return utils.StandardFloatCompare(bX, aX)
		}
		return utils.StandardFloatCompare(aX, bX)
	})
}

func (l ListOfEdge) SortByY(inverse bool) {
	l.Sort(func(a, b *Edge) int {
		aY := a.FirstVertex().Y()
		bY := b.FirstVertex().Y()
		if inverse {
			return utils.StandardFloatCompare(bY, aY)
		}
		return utils.StandardFloatCompare(aY, bY)
	})
}

func (l ListOfEdge) SortByZ(inverse bool) {
	l.Sort(func(a, b *Edge) int {
		aZ := a.FirstVertex().Z()
		bZ := b.FirstVertex().Z()
		if inverse {
			return utils.StandardFloatCompare(bZ, aZ)
		}
		return utils.StandardFloatCompare(aZ, bZ)
	})
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

func (e *Edge) FirstVertex() gp.Pnt {
	verts := e.Vertexes()
	if len(verts) < 1 {
		return gp.NewPnt(0, 0, 0)
	}

	return verts[0].ToPoint()
}

func (e *Edge) LastVertex() gp.Pnt {
	verts := e.Vertexes()
	if len(verts) < 1 {
		return gp.NewPnt(0, 0, 0)
	}

	return verts[len(verts)-1].ToPoint()
}

func (e *Edge) projectPointToSketch(solver SketchSolver, point gp.Pnt) (float64, float64) {
	origin := solver.CoordinateSystem().Location()
	originVec := gp.NewVec(origin.X(), origin.Y(), origin.Z())
	pointVec := gp.NewVec(point.X(), point.Y(), point.Z())
	xDir := solver.CoordinateSystem().XDirection()
	u := gp.NewVec(xDir.X(), xDir.Y(), xDir.Z())
	yDir := solver.CoordinateSystem().YDirection()
	v := gp.NewVec(yDir.X(), yDir.Y(), yDir.Z())
	x := u.Dot(pointVec) - u.Dot(originVec)
	y := v.Dot(pointVec) - v.Dot(originVec)
	return x, y
}

func (e *Edge) GetLine(solver SketchSolver) *Line {
	if !e.IsLine() {
		return nil
	}

	ex := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(&e.Edge)), topexp.Vertex)
	startX, startY := e.projectPointToSketch(solver, topods.NewVertexFromRef(topods.TopoDSVertex(ex.Current().Shape)).Pnt())
	ex.Next()
	endX, endY := e.projectPointToSketch(solver, topods.NewVertexFromRef(topods.TopoDSVertex(ex.Current().Shape)).Pnt())

	line := solver.CreateLine(startX, startY, endX, endY)
	line.Start.VerticalDistance(solver.XAxis(), startY)
	line.Start.HorizontalDistance(solver.YAxis(), startX)
	line.End.VerticalDistance(solver.XAxis(), endY)
	line.End.HorizontalDistance(solver.YAxis(), endX)
	solver.MakeFixed(line)
	return line
}

func (e *Edge) LineLength() float64 {
	return gcpnts.CurveLength(brepadapter.NewCurve(e.Edge))
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
	if utils.StandardFloatCompare(gp.NewVec(centerX, centerY, 0).Magnitude(), 0.0) == 0 {
		circ.Center.Coincident(solver.Origin())
	} else {
		circ.Center.VerticalDistance(solver.XAxis(), centerY)
		circ.Center.HorizontalDistance(solver.YAxis(), centerX)
	}
	circ.Diameter(radius * 2)
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

func (e *Edge) IsParallel(v gp.Dir) bool {
	if !e.IsLine() {
		return false
	}

	verts := e.Vertexes()
	if len(verts) < 2 {
		return false
	}
	first := verts[0].ToPoint()
	last := verts[len(verts)-1].ToPoint()
	dir := gp.NewDir(last.X()-first.X(), last.Y()-first.Y(), last.Z()-first.Z())

	return dir.IsParallel(v)
}

func (e *Edge) Midpoint() gp.Pnt {
	if !e.IsLine() {
		return gp.NewPnt(0, 0, 0)
	}

	verts := e.Vertexes()
	if len(verts) < 2 {
		return gp.NewPnt(0, 0, 0)
	}
	first := verts[0].ToPoint()
	last := verts[len(verts)-1].ToPoint()
	return gp.NewPnt((last.X()+first.X())/2., (last.Y()+first.Y())/2., (last.Z()+first.Z())/2.)
}

func (e *Edge) Vertexes() ListOfVertex {
	edges := make(ListOfVertex, 0)
	explorer := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(e.Edge.Edge)), topexp.Vertex)
	for ; explorer.More(); explorer.Next() {
		if explorer.Depth() > 1 {
			continue
		}
		edges = append(edges, NewVertexFromRef(explorer.Current()))
	}

	return edges
}
