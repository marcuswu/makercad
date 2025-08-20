package sketcher

import "github.com/marcuswu/gooccwrapper/gp"

// SketchSolver is implemented by the 2D Geometric Constraint Solvers MakerCad supports.
// MakerCad object provides a way to create a new sketch.
type SketchSolver interface {
	CoordinateSystem() gp.Ax3
	Origin() *Point
	XAxis() *Line
	YAxis() *Line
	CreatePoint(x float64, y float64) *Point
	CreateLine(startX float64, startY float64, endX float64, endY float64) *Line
	CreateCircle(centerX float64, centerY float64, radius float64) *Circle
	CreateArc(centerX float64, centerY float64, startX float64, startY float64, endX float64, endY float64) *Arc

	Coincident(Entity, Entity)
	PointVerticalDistance(*Point, Entity, float64)
	PointHorizontalDistance(*Point, Entity, float64)
	PointProjectedDistance(*Point, Entity, float64)
	LineMidpoint(*Line, Entity)
	LineAngle(*Line, *Line, float64)
	ArcLineTangent(*Arc, *Line)
	Distance(Entity, Entity, float64)
	HorizontalLine(*Line)
	HorizontalPoints(*Point, *Point)
	VerticalLine(*Line)
	VerticalPoints(*Point, *Point)
	LineLength(*Line, float64)
	Equal(Entity, Entity)
	CurveDiameter(Entity, float64)
	MakeFixed(Entity)
	Transform() gp.Trsf
	Solve() error
	OverConstrained() []string
	Entities() []Entity
	LogDebug(string) error
	ExportImage(string, ...float64) error
}
