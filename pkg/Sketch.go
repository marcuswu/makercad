package makercad

import "github.com/marcuswu/libmakercad/pkg/sketch"

type Sketch struct {
	solver sketch.SketchSolver
}

// No constructors -- Sketches should be created via MakerCad

func (s *Sketch) Solve() error {
	return s.solver.Solve()
}

func (s *Sketch) Origin() *sketch.Point {
	return s.solver.Origin()
}

func (s *Sketch) Arc(centerX float64, centerY float64, startX float64, startY float64, endX float64, endY float64) *sketch.Arc {
	return s.solver.CreateArc(centerX, centerY, startX, startY, endX, endY)
}

func (s *Sketch) Circle(centerX float64, centerY float64, diameter float64) *sketch.Circle {
	return s.solver.CreateCircle(centerX, centerY, diameter/2.0)
}

func (s *Sketch) Line(startX float64, startY float64, endX float64, endY float64) *sketch.Line {
	return s.solver.CreateLine(startX, startY, endX, endY)
}

func (s *Sketch) Point(x float64, y float64) *sketch.Point {
	return s.solver.CreatePoint(x, y)
}

func (s *Sketch) OverConstrained() []string {
	return s.solver.OverConstrained()
}

func (s *Sketch) DebugGraph(file string) error {
	return s.solver.LogDebug(file)
}

func (s *Sketch) Project(edge *sketch.Edge) sketch.Entity {
	if edge.IsCircle() {
		return edge.GetCircle(s.solver)
	}
	if edge.IsLine() {
		return edge.GetLine(s.solver)
	}
	return nil
}
