package makercad

import "github.com/marcuswu/makercad/sketcher"

// Sketch represents a 2D sketch on a face or plane. Sketches can be solved for a set of constraints. Sketches are created via an instance of [MakerCad]
type Sketch struct {
	solver sketcher.SketchSolver
}

// Solve will attempt to solve the sketch based on the established constraints
func (s *Sketch) Solve() error {
	return s.solver.Solve()
}

// Origin returns the origin element for the sketch
func (s *Sketch) Origin() *sketcher.Point {
	return s.solver.Origin()
}

// XAxis returns the X axis element for the sketch (represented as a [Line])
func (s *Sketch) XAxis() *sketcher.Line {
	return s.solver.XAxis()
}

// YAxis returns the Y axis element for the sketch (represented as a [Line])
func (s *Sketch) YAxis() *sketcher.Line {
	return s.solver.YAxis()
}

// Arc creates an arc clockwise from start to end around center
func (s *Sketch) Arc(centerX float64, centerY float64, startX float64, startY float64, endX float64, endY float64) *sketcher.Arc {
	return s.solver.CreateArc(centerX, centerY, startX, startY, endX, endY)
}

// Circle creates a circle around center with the specified diameter (does not automatically create a diameter constraint)
func (s *Sketch) Circle(centerX float64, centerY float64, diameter float64) *sketcher.Circle {
	return s.solver.CreateCircle(centerX, centerY, diameter/2.0)
}

// Line creates a line through the specified points. Automatically creates start and end points and sets them coincident to the line
func (s *Sketch) Line(startX float64, startY float64, endX float64, endY float64) *sketcher.Line {
	return s.solver.CreateLine(startX, startY, endX, endY)
}

// Point creates a point at the specified location
func (s *Sketch) Point(x float64, y float64) *sketcher.Point {
	return s.solver.CreatePoint(x, y)
}

// OverConstrained returns a string representation of conflicting constraints
func (s *Sketch) OverConstrained() []string {
	return s.solver.OverConstrained()
}

// DebugGraph outputs a GraphViz formatted graph representing the current sketch graph
func (s *Sketch) DebugGraph(file string) error {
	return s.solver.LogDebug(file)
}

// ExportImage writes an image representing the current sketch to the specified file
func (s *Sketch) ExportImage(file string, args ...float64) error {
	return s.solver.ExportImage(file, args...)
}

// Project projects an edge to the current sketch. Note that curves are not projected, just points. A circle projected between non-normal surfaces will not become an ellipse.
func (s *Sketch) Project(edge *sketcher.Edge) sketcher.Entity {
	if edge.IsCircle() {
		return edge.GetCircle(s.solver)
	}
	if edge.IsLine() {
		return edge.GetLine(s.solver)
	}
	return nil
}
