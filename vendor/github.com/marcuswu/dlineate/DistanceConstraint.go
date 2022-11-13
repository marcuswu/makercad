package dlineate

import (
	ic "github.com/marcuswu/dlineate/internal/constraint"
	"github.com/marcuswu/dlineate/utils"
)

func DistanceConstraint(p1 *Element, p2 *Element) *Constraint {
	constraint := emptyConstraint()
	constraint.elements = append(constraint.elements, p1)
	if p2 != nil {
		constraint.elements = append(constraint.elements, p2)
	}
	constraint.constraintType = Distance
	constraint.state = Resolved

	return constraint
}

func (s *Sketch) addDistanceConstraint(p1 *Element, p2 *Element, v float64) *ic.Constraint {
	switch p1.elementType {
	case Point:
		if p2.elementType != Point {
			return s.addDistanceConstraint(p2, p1, v)
		}

		return s.sketch.AddConstraint(ic.Distance, p1.element, p2.element, v)
	case Circle:
		if p2 == nil {
			// If p2 is nil, we're setting the circle radius
			// This is more of a placeholder for being able to fulfill other constraints as there is no
			// element to constrain to a distance from the center
			// Add a constraint to pkg/Sketch (not translatable to internal solver)
			return nil
		}
		return nil
	case Axis:
		fallthrough
	case Line:
		if p2 == nil {
			utils.Logger.Debug().Msgf(
				"Adding distance constraint for line %d. Translating to distance constraint between points %d and %d",
				p1.element.GetID(),
				p1.children[0].element.GetID(),
				p1.children[1].element.GetID(),
			)
			return s.sketch.AddConstraint(ic.Distance, p1.children[0].element, p1.children[1].element, v)
		}
		isCircle := p2.elementType == Circle
		isArc := p2.elementType == Arc
		if isArc || isCircle {
			return s.addDistanceConstraint(p2, p1, v)
		}
		return s.sketch.AddConstraint(ic.Distance, p1.element, p2.element, v)
	case Arc:
		if p2 == nil {
			// Add a constraint to pkg/Sketch (not translatable to internal solver)
			// If p2 is nil, we're setting the arc radius, so distance to start or end works
			return nil
		}
		// If p2 is not nil, we need to know the arc's radius is constrained
		return nil
	}
	return nil
}

func (s *Sketch) AddDistanceConstraint(p1 *Element, p2 *Element, v float64) *Constraint {
	c := DistanceConstraint(p1, p2)

	constraint := s.addDistanceConstraint(p1, p2, v)
	if constraint != nil {
		utils.Logger.Debug().Msgf("AddDistanceConstraint: added constraint id %d", constraint.GetID())
		p1.constraints = append(p1.constraints, constraint)
		if p2 != nil {
			p2.constraints = append(p2.constraints, constraint)
		}
		c.constraints = append(c.constraints, constraint)
	} else {
		if p2 != nil {
			c.state = Unresolved
		}
		c.dataValue = v
	}
	s.constraints = append(s.constraints, c)

	// This is might be wrong unless p1.element is always in constraint c -- for a line this is not true
	// check to see how eToC is used!
	s.eToC[p1.id] = append(s.eToC[p1.id], c)
	if p2 != nil {
		s.eToC[p2.id] = append(s.eToC[p2.id], c)
	}

	s.resolveDistanceConstraint(c)

	return c
}

func (s *Sketch) resolveCurveDistance(e1 *Element, e2 *Element, c *Constraint) bool {
	/*
		To resolve a curve's radius, we need either:
		1. A resolved distance constraint on the pkg/Element curve
		2. A solved center point and a solved element constrained to a distance from the pkg/Element curve
	*/
	if c.state == Resolved || c.state == Solved {
		return true
	}
	if e1 == nil {
		return false
	}
	eRadius, ok := s.resolveCurveRadius(e1)
	if !ok {
		return false
	}

	utils.Logger.Debug().
		Float64("center x", e1.values[0]).
		Float64("center y", e1.values[1]).
		Float64("radius", eRadius).
		Msg("Resolved curve radius")
	var constraint *ic.Constraint = nil
	if e2 != nil {
		constraint = s.sketch.AddConstraint(ic.Distance, e1.element, e2.element, eRadius+c.dataValue)
	}
	utils.Logger.Debug().
		Uint("constraint", constraint.GetID()).
		Msgf("resolveDistanceConstraint: added constraint")
	if constraint != nil {
		e1.constraints = append(e1.constraints, constraint)
		c.constraints = append(c.constraints, constraint)
	}
	s.constraints = append(s.constraints, c)
	if c.state != Solved {
		c.state = Resolved
	}

	return c.state == Resolved || c.state == Solved
}

func (s *Sketch) resolveDistanceConstraint(c *Constraint) bool {
	p1 := c.elements[0]
	var p2 *Element = nil
	if len(c.elements) > 1 {
		p2 = c.elements[1]
	}
	if len(c.constraints) > 0 {
		c.state = Resolved
		return true
	}
	if c.state == Resolved || c.state == Solved {
		return true
	}

	if s.resolveCurveDistance(p1, p2, c) {
		return c.state == Resolved || c.state == Solved
	}

	if s.resolveCurveDistance(p2, p1, c) {
		return c.state == Resolved || c.state == Solved
	}

	return c.state == Resolved || c.state == Solved
}
