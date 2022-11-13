package solver

import (
	"math"

	"github.com/marcuswu/dlineate/internal/constraint"
	el "github.com/marcuswu/dlineate/internal/element"
	"github.com/marcuswu/dlineate/utils"
)

func SolveForLine(c1 *constraint.Constraint, c2 *constraint.Constraint) SolveState {
	line, solveState := LineResult(c1, c2)

	if line == nil {
		return solveState
	}

	c1e, _ := c1.Element(line.GetID())
	c1Line := c1e.AsLine()
	c1Line.SetA(line.GetA())
	c1Line.SetB(line.GetB())
	c1Line.SetC(line.GetC())

	c2e, _ := c2.Element(line.GetID())
	c2Line := c2e.AsLine()
	c2Line.SetA(line.GetA())
	c2Line.SetB(line.GetB())
	c2Line.SetC(line.GetC())

	return solveState
}

func LineResult(c1 *constraint.Constraint, c2 *constraint.Constraint) (*el.SketchLine, SolveState) {
	/*
		There are only two possibilities:
		 * two lines and a point: The two lines must have an angle constraint between them
		 * two points and a line
	*/
	_, numLines := typeCounts(c1, c2)
	// 2 lines, 1 point -> LineFromPointLine
	var line *el.SketchLine = nil
	var solveState SolveState = NonConvergent
	if numLines == 3 {
		line, solveState = LineFromPointLine(c1, c2)
		utils.Logger.Trace().
			Str("result", solveState.String()).
			Str("line", line.String()).
			Msg("LineFromPointLine result")
	}

	// 1 line, 2 points -> LineFromPoints
	if numLines == 2 {
		line, solveState = LineFromPoints(c1, c2)
		utils.Logger.Trace().
			Str("result", solveState.String()).
			Str("line", line.String()).
			Msgf("LineFromPoints result")
	}

	if solveState == Solved {
		c1.Solved = true
		c2.Solved = true
	}

	return line, solveState
}

// MoveLineToPoint solves a constraint between a line and a point where the line needs to move
func MoveLineToPoint(c *constraint.Constraint) SolveState {
	if c.Type != constraint.Distance {
		utils.Logger.Error().
			Uint("constraint", c.GetID()).
			Msg("MoveLineToPoint constraint was not Distance type")
		return NonConvergent
	}

	var point *el.SketchPoint
	var line *el.SketchLine
	var e1Type = c.Element1.GetType()
	var e2Type = c.Element2.GetType()
	if e1Type == e2Type {
		utils.Logger.Error().
			Uint("constraint", c.GetID()).
			Msg("MoveLineToPoint did not have the correct element types")
		return NonConvergent
	}
	if e1Type == el.Point && e2Type == el.Line {
		point = c.Element1.(*el.SketchPoint)
		line = c.Element2.(*el.SketchLine)
	}
	if e2Type == el.Point && e1Type == el.Line {
		point = c.Element2.(*el.SketchPoint)
		line = c.Element1.(*el.SketchLine)
	}

	// If two points, get distance between them, translate constraint value - distance between
	// If point and line, get distance between them, translate normal to line constraint value - distance between
	dist := line.DistanceTo(point)
	translate1 := dist + c.GetValue()
	translate2 := dist - c.GetValue()

	if math.Abs(translate1) < math.Abs(translate2) {
		line.TranslateDistance(translate1)
	} else {
		line.TranslateDistance(translate2)
	}

	c.Solved = true

	return Solved
}

func LineFromPoints(c1 *constraint.Constraint, c2 *constraint.Constraint) (*el.SketchLine, SolveState) {
	line := c1.First().AsLine()
	if line == nil {
		line = c1.Second().AsLine()
	}

	if line == nil {
		utils.Logger.Error().
			Uint("constraint 1", c1.GetID()).
			Uint("constraint 2", c2.GetID()).
			Msg("LineFromPoints could not find the line to work with.")
		return line, NonConvergent
	}

	p1e, _ := c1.Other(line.GetID())
	p2e, _ := c2.Other(line.GetID())
	p1 := p1e.AsPoint()
	p2 := p2e.AsPoint()
	if p1 == nil || p2 == nil {
		utils.Logger.Error().
			Uint("constraint 1", c1.GetID()).
			Uint("constraint 2", c2.GetID()).
			Msg("LineFromPoints could not find the points to work with.")
		return line, NonConvergent
	}
	p1Dist := c1.Value
	p2Dist := c2.Value

	// Special case where distances are both 0
	if p1Dist == 0 && p2Dist == 0 {
		la1 := p2.Y - p1.Y                  // y' - y
		lb1 := p1.X - p2.X                  // x - x'
		lc1 := (-la1 * p1.X) - (lb1 * p1.Y) // c = -ax - by from ax + by + c = 0
		la2 := p1.Y - p2.Y                  // y' - y
		lb2 := p2.X - p1.X                  // x - x'
		lc2 := (-la2 * p1.X) - (lb2 * p1.Y) // c = -ax - by from ax + by + c = 0
		lineV := &el.Vector{X: line.GetA(), Y: line.GetB()}
		angleTo1 := lineV.AngleTo(&el.Vector{X: la1, Y: lb1})
		angleTo2 := lineV.AngleTo(&el.Vector{X: la2, Y: lb2})
		line.SetA(la1)
		line.SetB(lb1)
		line.SetC(lc1)
		if math.Abs(angleTo2) < math.Abs(angleTo1) {
			line.SetA(la2)
			line.SetB(lb2)
			line.SetC(lc2)
		}
		return line, Solved
	}

	// Rotate line to horizontal (and rotate points the same)
	// Translate line p2Dist so it lies on p2
	// The line must be tangent to the two circles defined by the two points and their distances
	// TODO: fix this check -- this is not true for external tangents!
	if p1.DistanceTo(p2) < p1Dist+p2Dist {
		utils.Logger.Error().
			Uint("constraint 1", c1.GetID()).
			Uint("constraint 2", c2.GetID()).
			Msgf("LineFromPoints determined the points were too close together.")
		return line, NonConvergent
	}

	// Math from https://en.wikipedia.org/wiki/Tangent_lines_to_circles#Analytic_geometry
	deltaR := p2Dist - p1Dist
	deltaX := p2.X - p1.X
	deltaY := p2.Y - p1.Y
	d := p1.DistanceTo(p2)
	R := deltaR / d
	X := deltaX / d
	Y := deltaY / d
	rSquared := R * R

	// Internal vs external tangents will be handled by positive or negative distance constraint values
	// Both the same sign will be external, opposing signs will be internal
	// There will be two options aside from internal or external -- plus or minus k
	// Use the one closest to the existing line angle (closest slope)
	var k float64 = 1
	tanA1 := (R * X) - (k*Y)*math.Sqrt(1.0-rSquared)
	tanB1 := (R * Y) + (k*X)*math.Sqrt(1.0-rSquared)
	tanC1 := p1Dist - ((tanA1 * p1.X) + (tanB1 * p1.Y))

	k = -1
	tanA2 := (R * X) - (k*Y)*math.Sqrt(1.0-rSquared)
	tanB2 := (R * Y) + (k*X)*math.Sqrt(1.0-rSquared)
	tanC2 := p1Dist - ((tanA2 * p1.X) + (tanB2 * p1.Y))

	k = 1
	R = (-p2Dist - p1Dist) / d
	rSquared = R * R
	tanA3 := (R * X) - (k*Y)*math.Sqrt(1.0-rSquared)
	tanB3 := (R * Y) + (k*X)*math.Sqrt(1.0-rSquared)
	tanC3 := p1Dist - ((tanA3 * p1.X) + (tanB3 * p1.Y))

	k = -1
	tanA4 := (R * X) - (k*Y)*math.Sqrt(1.0-rSquared)
	tanB4 := (R * Y) + (k*X)*math.Sqrt(1.0-rSquared)
	tanC4 := p1Dist - ((tanA4 * p1.X) + (tanB4 * p1.Y))

	origSlope := line.GetB() / line.GetA()
	externalSlope := tanB1 / tanA1
	internalSlope := tanB3 / tanA3
	externalSlopeDist := math.Abs(externalSlope - origSlope)
	internalSlopeDist := math.Abs(internalSlope - origSlope)
	mag1 := math.Sqrt(tanA1*tanA1) + (tanB1 * tanB1)
	mag2 := math.Sqrt(tanA2*tanA2) + (tanB2 * tanB2)
	c1Normal := tanC1 / mag1
	c2Normal := tanC2 / mag2
	useInternal := externalSlopeDist > internalSlopeDist
	if useInternal {
		tanA1 = tanA3
		tanB1 = tanB3
		tanC1 = tanC3
		tanA2 = tanA4
		tanB2 = tanB4
		tanC2 = tanC4
		mag1 := math.Sqrt(tanA1*tanA1) + (tanB1 * tanB1)
		mag2 := math.Sqrt(tanA2*tanA2) + (tanB2 * tanB2)
		c1Normal = tanC1 / mag1
		c2Normal = tanC2 / mag2
	}

	originalOriginDistance := line.GetOriginDistance()
	originDistance1 := math.Abs(c1Normal - originalOriginDistance)
	originDistance2 := math.Abs(c2Normal - originalOriginDistance)
	useOption1 := originDistance2 > originDistance1

	line.SetA(tanA1)
	line.SetB(tanB1)
	line.SetC(tanC1)
	if !useOption1 {
		line.SetA(tanA2)
		line.SetB(tanB2)
		line.SetC(tanC2)
	}

	return line, Solved
}

func LineFromPointLine(c1 *constraint.Constraint, c2 *constraint.Constraint) (*el.SketchLine, SolveState) {
	var targetLine *el.SketchLine
	var point *el.SketchPoint
	distC := c1
	angleC := c2
	if c1.Type == constraint.Angle {
		angleC = c1
		distC = c2
	}

	targetLine = distC.First().AsLine()
	point = distC.Second().AsPoint()
	if targetLine == nil {
		targetLine = distC.Second().AsLine()
		point = distC.First().AsPoint()
	}

	// Solve angle
	newLine, state := SolveAngleConstraint(angleC, targetLine.GetID())

	// Translate to distC.Value from the point
	dist1 := newLine.DistanceTo(point) - distC.Value
	dist2 := newLine.DistanceTo(point) + distC.Value
	line1 := newLine.TranslatedDistance(dist1)
	line2 := newLine.TranslatedDistance(dist2)

	line1Distance := targetLine.DistanceTo(line1)
	line2Distance := targetLine.DistanceTo(line2)

	if math.Abs(line1Distance) < math.Abs(line2Distance) {
		return line1, state
	}

	return line2, state
}

// SolveAngleConstraint solve an angle constraint between two lines
func SolveAngleConstraint(c *constraint.Constraint, e uint) (*el.SketchLine, SolveState) {
	if c.Type != constraint.Angle {
		utils.Logger.Error().
			Uint("constraint", c.GetID()).
			Msgf("SolveAngleConstraint was not sent an angle constraint")
		return nil, NonConvergent
	}

	l1 := c.Element1.(*el.SketchLine)
	l2 := c.Element2.(*el.SketchLine)
	desired := c.Value
	if l1.GetID() == e {
		l1, l2 = l2, l1
	}

	angle1 := l2.AngleToLine(l1)
	rotate1 := angle1 + desired
	rotate2 := desired + angle1
	reverseRotate1 := angle1 - desired
	reverseRotate2 := desired - angle1

	lines := []*el.SketchLine{
		l2.Rotated(rotate1),
		l2.Rotated(rotate2),
		l2.Rotated(reverseRotate1),
		l2.Rotated(reverseRotate2),
	}

	var newLine *el.SketchLine = nil
	for _, line := range lines {
		if newLine == nil || math.Abs(line.AngleToLine(l2)) < math.Abs(newLine.AngleToLine(l2)) {
			newLine = line
		}
	}

	c.Solved = true
	return newLine, Solved
}
