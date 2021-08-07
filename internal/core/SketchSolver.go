package core

/*
 * 2D Geometric Constraint Solver should implement this interface.
 * MakerCad object provides a way to create a new sketch.
 */
type SketchSolver interface {
	//createSketch(PlaneParameters) Sketch2D
	// LookupEntity(uint) Entity
	CreatePoint(float64, float64) *Point
	CreateLine(Entity, Entity) *Line
	CreateCircle(Entity, float64) *Circle
	CreateArc(Entity, Entity, Entity) *Arc
	CreateDistance(float64) *Distance
	// Creates fixed entities not solved for
	CreateWorkplanePoint(float64, float64)
	CreateWorkplaneLine(Entity, Entity) Entity
	CreateWorkplaneCircle(Entity, float64) Entity
	CreateWorkplaneArc(Entity, Entity, Entity)
	CreateWorkplaneDistance(float64) Entity

	PointCoincident(Entity, Entity)
	PointVerticalDistance(Entity, Entity, float64)
	PointHorizontalDistance(Entity, Entity, float64)
	PointProjectedDistance(Entity, Entity, Entity, float64)
	LineSymmetric(Entity, Entity, Entity)
	LineMidpoint(Entity, Entity)
	LineAngle(Entity, Entity, float64)
	ArcLineTangent(Entity, Entity, int)
	Distance(float64, Entity, Entity)
	Horizontal(Entity, Entity)
	Vertical(Entity, Entity)
	LineLength(float64, Entity)
	EqualCircles(Entity, Entity)
	CircleDiameter(Entity, float64)
	Solve()
	ToFace() *Face
}
