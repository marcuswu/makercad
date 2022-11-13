package dlineate

import (
	"errors"

	"github.com/marcuswu/dlineate/utils"
)

func TangentConstraint(p1 *Element, p2 *Element) *Constraint {
	constraint := emptyConstraint()
	constraint.elements = append(constraint.elements, p1)
	constraint.elements = append(constraint.elements, p2)
	// constraint.elements = append(constraint.elements, p3)
	constraint.constraintType = Tangent
	constraint.state = Unresolved

	return constraint
}

func (s *Sketch) AddTangentConstraint(p1 *Element, p2 *Element) (*Constraint, error) {
	var line, curve, err = orderParams(p1, p2)

	if err != nil {
		utils.Logger.Error().Msg("Tangent constraint had incorrect parameters")
		return nil, err
	}

	c := TangentConstraint(line, curve)
	s.eToC[p1.id] = append(s.eToC[p1.id], c)
	s.eToC[p2.id] = append(s.eToC[p2.id], c)
	s.constraints = append(s.constraints, c)
	// s.eToC[p3.element.GetID()] = append(s.eToC[p3.element.GetID()], c)

	s.resolveTangentConstraint(c)

	return c, nil
}

func orderParams(p1 *Element, p2 *Element) (*Element, *Element, error) {
	var line /*point,*/, curve *Element

	switch Line {
	case p1.elementType:
		line = p1
	default:
		line = p2
	}

	switch true {
	case p1.elementType == Circle || p1.elementType == Arc:
		curve = p1
	default:
		curve = p2
	}

	if line == curve {
		return p1, p2, errors.New("incorrect element types for tangent constraint")
	}

	return line, curve, nil
}

func (s *Sketch) resolveTangentConstraint(c *Constraint) bool {
	radius, ok := s.resolveCurveRadius(c.elements[1])
	if ok {
		utils.Logger.Debug().
			Str("element 1", c.elements[0].String()).
			Str("element 2", c.elements[1].children[0].String()).
			Msg("addDistanceConstraint")
		constraint := s.addDistanceConstraint(c.elements[0], c.elements[1].children[0], radius)
		utils.Logger.Debug().
			Uint("constraint", constraint.GetID()).
			Msg("resolveTangentConstraint: added constraint")
		c.elements[0].constraints = append(c.elements[0].constraints, constraint)
		c.elements[1].constraints = append(c.elements[1].constraints, constraint)
		c.constraints = append(c.constraints, constraint)
		c.state = Resolved
	}

	return c.state == Resolved
}
