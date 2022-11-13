package dlineate

import "github.com/marcuswu/dlineate/utils"

func (s *Sketch) AddPerpendicularConstraint(p1 *Element, p2 *Element) (*Constraint, error) {
	c, err := s.AddAngleConstraint(p1, p2, 90, false)
	if err != nil {
		utils.Logger.Error().Msgf("error: %s", err)
	}
	if c != nil {
		c.constraintType = Perpendicular
	}
	return c, err
}
