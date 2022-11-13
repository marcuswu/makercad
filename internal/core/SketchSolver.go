package core

/*
 * 2D Geometric Constraint Solver should implement this interface.
 * MakerCad object provides a way to create a new sketch.
 */
type SketchSolver interface {
	//createSketch(PlaneParameters) Sketch2D
	// LookupEntity(uint) Entity
	CreatePoint(float64, float64) *Point
	CreateLine(*Point, *Point) *Line
	CreateCircle(*Point, float64) *Circle
	CreateArc(*Point, *Point, *Point) *Arc
	// CreateDistance(float64) *Distance
	// Creates fixed entities not solved for
	CreateWorkplanePoint(float64, float64) *Point
	CreateWorkplaneLine(*Point, *Point) *Line
	CreateWorkplaneCircle(*Point, float64) *Circle
	CreateWorkplaneArc(*Point, *Point, *Point) *Arc

	Coincident(Entity, Entity)
	PointVerticalDistance(*Point, Entity, float64)
	PointHorizontalDistance(*Point, Entity, float64)
	PointProjectedDistance(*Point, Entity, Entity, float64)
	LineMidpoint(*Line, Entity)
	LineAngle(*Line, *Line, float64)
	ArcLineTangent(*Arc, *Line)
	Distance(Entity, Entity, float64)
	Horizontal(Entity, Entity)
	Vertical(Entity, Entity)
	LineLength(*Line, float64)
	Equal(Entity, Entity)
	CurveDiameter(Entity, float64)
	Solve()
	ToFace() *Face
}
