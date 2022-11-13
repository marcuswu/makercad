package dlineate

import (
	"errors"
	"math"

	ic "github.com/marcuswu/dlineate/internal/constraint"
)

func AngleConstraint(p1 *Element, p2 *Element) *Constraint {
	constraint := emptyConstraint()
	constraint.elements = append(constraint.elements, p1)
	constraint.elements = append(constraint.elements, p2)
	constraint.constraintType = Angle
	constraint.state = Resolved

	return constraint
}

// AddAngleConstraint adds a constraint between the lines p1 and p2 where the counter-clockwise angle
// from p1 to p2 is the positive direction
func (s *Sketch) AddAngleConstraint(p1 *Element, p2 *Element, v float64, useSupplementary bool) (*Constraint, error) {
	c := AngleConstraint(p1, p2)

	if (p1.elementType != Line && p1.elementType != Axis) || (p2.elementType != Line && p2.elementType != Axis) {
		return nil, errors.New("incorrect element types for angle constraint")
	}

	radians := v / 180 * math.Pi
	radiansAlt := math.Pi - math.Abs(radians)

	if useSupplementary {
		// if useSupplementary || math.Abs(math.Abs(currentAngle)-math.Abs(radiansAlt)) < math.Abs(math.Abs(currentAngle)-math.Abs(radians)) {
		radians = radiansAlt
	}

	constraint := s.sketch.AddConstraint(ic.Angle, p1.element, p2.element, radians)
	p1.constraints = append(p1.constraints, constraint)
	p2.constraints = append(p2.constraints, constraint)
	c.constraints = append(c.constraints, constraint)
	s.constraints = append(s.constraints, c)
	s.eToC[p1.id] = append(s.eToC[p1.id], c)
	s.eToC[p2.id] = append(s.eToC[p2.id], c)

	return c, nil
}
