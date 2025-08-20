package makercad

import "github.com/marcuswu/makercad/sketcher"

type Sketch struct {
	solver sketcher.SketchSolver
}

// No constructors -- Sketches should be created via MakerCad

func (s *Sketch) Solve() error {
	return s.solver.Solve()
}

func (s *Sketch) Origin() *sketcher.Point {
	return s.solver.Origin()
}

func (s *Sketch) XAxis() *sketcher.Line {
	return s.solver.XAxis()
}

func (s *Sketch) YAxis() *sketcher.Line {
	return s.solver.YAxis()
}

// Arc creates an arc clockwise from start to end around center
func (s *Sketch) Arc(centerX float64, centerY float64, startX float64, startY float64, endX float64, endY float64) *sketcher.Arc {
	return s.solver.CreateArc(centerX, centerY, startX, startY, endX, endY)
}

func (s *Sketch) Circle(centerX float64, centerY float64, diameter float64) *sketcher.Circle {
	return s.solver.CreateCircle(centerX, centerY, diameter/2.0)
}

func (s *Sketch) Line(startX float64, startY float64, endX float64, endY float64) *sketcher.Line {
	return s.solver.CreateLine(startX, startY, endX, endY)
}

func (s *Sketch) Point(x float64, y float64) *sketcher.Point {
	return s.solver.CreatePoint(x, y)
}

func (s *Sketch) OverConstrained() []string {
	return s.solver.OverConstrained()
}

func (s *Sketch) DebugGraph(file string) error {
	return s.solver.LogDebug(file)
}

func (s *Sketch) ExportImage(file string, args ...float64) error {
	return s.solver.ExportImage(file, args...)
}

func (s *Sketch) Project(edge *sketcher.Edge) sketcher.Entity {
	if edge.IsCircle() {
		return edge.GetCircle(s.solver)
	}
	if edge.IsLine() {
		return edge.GetLine(s.solver)
	}
	return nil
}
