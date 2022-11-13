package solver

import (
	"math"

	"github.com/marcuswu/dlineate/internal/constraint"
	el "github.com/marcuswu/dlineate/internal/element"
	"github.com/marcuswu/dlineate/utils"
)

func SolveConstraint(c *constraint.Constraint) SolveState {
	if c.Type == constraint.Distance {
		return SolveDistanceConstraint(c)
	}
	newLine, state := SolveAngleConstraint(c, c.Element2.GetID())

	cl := c.Element2.AsLine()
	cl.SetA(newLine.GetA())
	cl.SetB(newLine.GetB())
	cl.SetC(newLine.GetC())

	return state
}

func SolveConstraints(c1 *constraint.Constraint, c2 *constraint.Constraint, solveFor el.SketchElement) SolveState {
	if solveFor.GetType() == el.Point {
		return SolveForPoint(c1, c2)
	}

	return SolveForLine(c1, c2)
}

func ConstraintResult(c1 *constraint.Constraint, c2 *constraint.Constraint, solveFor el.SketchElement) (el.SketchElement, SolveState) {
	if solveFor.GetType() == el.Point {
		return PointResult(c1, c2)
	}

	return LineResult(c1, c2)
}

// SolveConstraints solve two constraints and return the solution state
func SolveForPoint(c1 *constraint.Constraint, c2 *constraint.Constraint) SolveState {
	newP3, state := PointResult(c1, c2)

	if newP3 == nil {
		return state
	}

	c1e, _ := c1.Element(newP3.GetID())
	c1p := c1e.AsPoint()
	c1p.X = newP3.X
	c1p.Y = newP3.Y

	c2e, _ := c2.Element(newP3.GetID())
	c2p := c2e.AsPoint()
	c2p.X = newP3.X
	c2p.Y = newP3.Y

	return state
}

// PointResult returns the result of solving two constraints sharing one point
func PointResult(c1 *constraint.Constraint, c2 *constraint.Constraint) (*el.SketchPoint, SolveState) {
	numPoints, _ := typeCounts(c1, c2)
	// 4 points -> PointFromPoints
	var point *el.SketchPoint = nil
	var solveState SolveState = NonConvergent
	if numPoints == 4 {
		point, solveState = PointFromPoints(c1, c2)
	}

	// 3 points, 1 line -> PointFromPointLine
	if numPoints == 3 {
		point, solveState = PointFromPointLine(c1, c2)
	}
	// 2 points, 2 lines -> PointFromLineLine
	if numPoints == 2 {
		point, solveState = PointFromLineLine(c1, c2)
	}

	if solveState == Solved {
		c1.Solved = true
		c2.Solved = true
	}

	return point, solveState
}

// SolveDistanceConstraint solves a distance constraint and returns the solution state
func SolveDistanceConstraint(c *constraint.Constraint) SolveState {
	if c.Type != constraint.Distance {
		utils.Logger.Error().
			Uint("constraint", c.GetID()).
			Msgf("SolveDistanceConstraint: was not sent a distance constraint")
		return NonConvergent
	}

	var point *el.SketchPoint
	var other el.SketchElement
	if c.Element1.GetType() == el.Point {
		point = c.Element1.(*el.SketchPoint)
		other = c.Element2
	} else {
		point = c.Element2.(*el.SketchPoint)
		other = c.Element1
	}

	// If two points, get distance between them, translate constraint value - distance between
	// If point and line, get distance between them, translate normal to line constraint value - distance between
	trans := point.VectorTo(other)
	dist := trans.Magnitude()

	if dist == 0 && c.GetValue() > 0 {
		utils.Logger.Error().Msg("SolveDistanceConstraint: points are coincident, but they shouldn't be. Infinite solutions.")
		return NonConvergent
	}

	if dist == 0 && c.GetValue() == 0 {
		c.Solved = true
		return Solved
	}
	otherP := other.AsPoint()
	if otherP == nil {
		otherP = other.AsLine().NearestPoint(point.GetX(), point.GetY())
	}

	trans.Scaled(c.GetValue() / dist)
	newPoint := otherP.Translated(trans.GetX(), trans.GetY())
	point.X = newPoint.X
	point.Y = newPoint.Y
	c.Solved = true

	return Solved
}

// GetPointFromPoints calculates where a 3rd point exists in relation to two others with
// distance constraints from the first two
func GetPointFromPoints(p1 el.SketchElement, originalP2 el.SketchElement, originalP3 el.SketchElement, p1Radius float64, p2Radius float64) (*el.SketchPoint, SolveState) {
	// Don't mutate the originals
	p2 := el.CopySketchElement(originalP2)
	p3 := el.CopySketchElement(originalP3)
	pointDistance := p1.DistanceTo(p2)
	constraintDist := p1Radius + p2Radius

	if utils.StandardFloatCompare(pointDistance, constraintDist) > 0 {
		utils.Logger.Error().
			Uint("point 1", p1.GetID()).
			Uint("point 2", p2.GetID()).
			Msg("GetPointFromPoints no solution because the points are too far apart")
		return nil, NonConvergent
	}

	if utils.StandardFloatCompare(pointDistance, constraintDist) == 0 {
		translate := p1.VectorTo(p2)
		translate.Scaled(p1Radius / translate.Magnitude())
		newP3 := el.NewSketchPoint(p3.GetID(), p1.AsPoint().X-translate.X, p1.AsPoint().Y-translate.Y)
		return newP3, Solved
	}

	// Solve for p3
	// translate to p1 (p2 and p3)
	p2.ReverseTranslateByElement(p1)
	p3.ReverseTranslateByElement(p1)
	// rotate p2 and p3 so p2 is on x axis
	angle := p2.AngleTo(&el.Vector{X: 1, Y: 0})
	p2.Rotate(angle)
	p3.Rotate(angle)
	// calculate possible p3s
	p2Dist := p2.(*el.SketchPoint).GetX()

	// https://mathworld.wolfram.com/Circle-CircleIntersection.html
	xDelta := ((-(p2Radius * p2Radius) + (p2Dist * p2Dist)) + (p1Radius * p1Radius)) / (2 * p2Dist)
	yDelta := math.Sqrt((p1Radius * p1Radius) - (xDelta * xDelta))
	p3X := xDelta
	p3Y1 := yDelta
	p3Y2 := -yDelta
	// determine which is closest to the p3 from constraint
	newP31 := el.NewSketchPoint(p3.GetID(), p3X, p3Y1)
	newP32 := el.NewSketchPoint(p3.GetID(), p3X, p3Y2)
	actualP3 := newP31
	if newP32.SquareDistanceTo(p3) < newP31.SquareDistanceTo(p3) {
		actualP3 = newP32
	}
	// unrotate actualP3
	actualP3.Rotate(-angle)
	// untranslate actualP3
	actualP3.TranslateByElement(p1)

	// return actualP3
	return actualP3, Solved
}

// PointFromPoints calculates a new p3 representing p3 moved to satisfy
// distance constraints from p1 and p2
func PointFromPoints(c1 *constraint.Constraint, c2 *constraint.Constraint) (*el.SketchPoint, SolveState) {
	p1 := c1.Element1
	p2 := c2.Element1
	p3 := c1.Element2
	p1Radius := c1.GetValue()
	p2Radius := c2.GetValue()

	switch {
	case c1.Element1.Is(c2.Element1):
		p3, p1, p2 = c1.Element1, c1.Element2, c2.Element2
	case c1.Element2.Is(c2.Element1):
		p3, p1, p2 = c1.Element2, c1.Element1, c2.Element2
	case c1.Element1.Is(c2.Element2):
		p3, p1, p2 = c1.Element1, c1.Element2, c2.Element1
	case c1.Element2.Is(c2.Element2):
		break
	}

	return GetPointFromPoints(p1, p2, p3, p1Radius, p2Radius)
}

func pointFromPointLine(originalP1 el.SketchElement, originalL2 el.SketchElement, originalP3 el.SketchElement, pointDist float64, lineDist float64) (*el.SketchPoint, SolveState) {
	p1 := el.CopySketchElement(originalP1).(*el.SketchPoint)
	l2 := el.CopySketchElement(originalL2).(*el.SketchLine)
	p3 := el.CopySketchElement(originalP3).(*el.SketchPoint)
	distanceDifference := l2.DistanceTo(p1)

	// rotate l2 to X axis
	angle := l2.AngleTo(&el.Vector{X: 1, Y: 0})
	l2.Rotate(angle)
	p1.Rotate(angle)
	p3.Rotate(angle)

	// translate l2 to X axis
	yTranslate := l2.GetC() - lineDist
	if math.Abs(p1.GetY()+yTranslate) > pointDist {
		yTranslate = l2.GetC() + lineDist
	}
	l2.Translate(0, yTranslate)
	// move p1 to Y axis
	xTranslate := p1.GetX()
	p1.Translate(-xTranslate, yTranslate)
	p3.Translate(-xTranslate, yTranslate)

	if utils.StandardFloatCompare(pointDist, math.Abs(p1.GetY())) < 0 {
		utils.Logger.Error().
			Float64("point distance", pointDist).
			Float64("p1.y", math.Abs(p1.GetY())).
			Msg("pointFromPointLine: Nonconvergent")
		return nil, NonConvergent
	}

	// Find points where circle at p1 with radius pointDist intersects with x axis
	xPos := math.Sqrt(math.Abs((pointDist * pointDist) - (p1.GetY() * p1.GetY())))
	if utils.StandardFloatCompare(distanceDifference, 0) == 0 {
		xPos = pointDist
	}

	newP31 := el.NewSketchPoint(p3.GetID(), xPos, 0)
	newP32 := el.NewSketchPoint(p3.GetID(), -xPos, 0)
	actualP3 := newP31
	if newP32.SquareDistanceTo(p3) < newP31.SquareDistanceTo(p3) {
		actualP3 = newP32
	}
	actualP3.Translate(xTranslate, -yTranslate)
	actualP3.Rotate(-angle)

	return actualP3, Solved
}

// PointFromPointLine construct a point from a point and a line. c2 must contain the line.
func PointFromPointLine(c1 *constraint.Constraint, c2 *constraint.Constraint) (*el.SketchPoint, SolveState) {
	p1 := c1.Element1
	l2 := c2.Element1
	p3 := c1.Element2
	pointDist := c1.GetValue()
	lineDist := c2.GetValue()

	switch {
	case c1.Element1.Is(c2.Element1):
		p3 = c1.Element1
		p1 = c1.Element2
		l2 = c2.Element2
	case c1.Element2.Is(c2.Element1):
		p3 = c1.Element2
		p1 = c1.Element1
		l2 = c2.Element2
	case c1.Element1.Is(c2.Element2):
		p3 = c1.Element1
		p1 = c1.Element2
		l2 = c2.Element1
	case c1.Element2.Is(c2.Element2):
		break
	}

	if p1.GetType() == el.Line && l2.GetType() == el.Point {
		p1, l2 = l2, p1
		pointDist, lineDist = lineDist, pointDist
	}

	return pointFromPointLine(p1, l2, p3, pointDist, lineDist)
}

func pointFromLineLine(l1 *el.SketchLine, l2 *el.SketchLine, p3 *el.SketchPoint, line1Dist float64, line2Dist float64) (*el.SketchPoint, SolveState) {
	sameSlope := utils.StandardFloatCompare(l1.GetA(), l2.GetA()) == 0 && utils.StandardFloatCompare(l1.GetB(), l2.GetB()) == 0
	// If l1 and l2 are parallel, and line distances aren't what is passed in, there is no solution
	if sameSlope &&
		utils.StandardFloatCompare(line1Dist-line2Dist, l1.DistanceTo(l2)) != 0 {
		utils.Logger.Error().
			Uint("line 1", l1.GetID()).
			Uint("line 2", l2.GetID()).
			Msg("pointFromLineLine no solution to find a point because the lines are parallel")
		return nil, NonConvergent
	}

	// If l1 & l2 are parallel and it's solvable, there are infinite solutions
	// Choose the one closest to the current point location
	if sameSlope {
		translate := l1.VectorTo(p3)
		translate.Scaled((p3.DistanceTo(l1) - line1Dist) / translate.Magnitude())
		return el.NewSketchPoint(p3.GetID(), p3.X+translate.X, p3.Y+translate.Y), Solved
	}
	// Translate l1 line1Dist
	line1TranslatePos := l1.TranslatedDistance(line1Dist)
	line1TranslateNeg := l1.TranslatedDistance(-line1Dist)
	// Translate l2 line2Dist
	line2TranslatedPos := l2.TranslatedDistance(line2Dist)
	line2TranslatedNeg := l2.TranslatedDistance(-line2Dist)

	// If line1 and line2 are the same line,
	intersect1 := el.SketchPointFromVector(p3.GetID(), line1TranslatePos.Intersection(line2TranslatedPos))
	intersect2 := el.SketchPointFromVector(p3.GetID(), line1TranslatePos.Intersection(line2TranslatedNeg))
	intersect3 := el.SketchPointFromVector(p3.GetID(), line1TranslateNeg.Intersection(line2TranslatedPos))
	intersect4 := el.SketchPointFromVector(p3.GetID(), line1TranslateNeg.Intersection(line2TranslatedNeg))

	// Return closest intersection point
	closest := intersect1
	dist := p3.DistanceTo(intersect1)
	if next := p3.DistanceTo(intersect2); next < dist {
		dist = next
		closest = intersect2
	}
	if next := p3.DistanceTo(intersect3); next < dist {
		dist = next
		closest = intersect3
	}
	if next := p3.DistanceTo(intersect4); next < dist {
		closest = intersect4
	}

	return closest, Solved
}

// PointFromLineLine construct a point from two lines. c2 must contain the point.
func PointFromLineLine(c1 *constraint.Constraint, c2 *constraint.Constraint) (*el.SketchPoint, SolveState) {
	l1 := c1.Element1
	l2 := c2.Element1
	p3 := c1.Element2
	line1Dist := c1.GetValue()
	line2Dist := c2.GetValue()

	switch {
	case c1.Element1.Is(c2.Element1):
		p3 = c1.Element1
		l1 = c1.Element2
		l2 = c2.Element2
	case c1.Element2.Is(c2.Element1):
		p3 = c1.Element2
		l1 = c1.Element1
		l2 = c2.Element2
	case c1.Element1.Is(c2.Element2):
		p3 = c1.Element1
		l1 = c1.Element2
		l2 = c2.Element1
	case c1.Element2.Is(c2.Element2):
		break
	}

	return pointFromLineLine(l1.AsLine(), l2.AsLine(), p3.AsPoint(), line1Dist, line2Dist)
}
