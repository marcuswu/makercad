package sketch

import "github.com/marcuswu/gooccwrapper/gp"

/*
 * 2D Geometric Constraint Solver should implement this interface.
 * MakerCad object provides a way to create a new sketch.
 */
type SketchSolver interface {
	//createSketch(PlaneParameters) Sketch2D
	CoordinateSystem() gp.Ax3
	Origin() *Point
	XAxis() *Line
	YAxis() *Line
	// LookupEntity(uint) Entity
	CreatePoint(x float64, y float64) *Point
	CreateLine(startX float64, startY float64, endX float64, endY float64) *Line
	CreateCircle(centerX float64, centerY float64, radius float64) *Circle
	CreateArc(centerX float64, centerY float64, startX float64, startY float64, endX float64, endY float64) *Arc
	// CreateDistance(float64) *Distance
	/*
		// These create fixed entities not solved for
		CreateWorkplanePoint(float64, float64) *Point
		CreateWorkplaneLine(*Point, *Point) *Line
		CreateWorkplaneCircle(*Point, float64) *Circle
		CreateWorkplaneArc(*Point, *Point, *Point) *Arc
	*/

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
}
