package dlineate

func (s *Sketch) AddEqualConstraint(p1 *Element, p2 *Element) *Constraint {
	c := s.AddRatioConstraint(p1, p2, 1)
	return c
}
