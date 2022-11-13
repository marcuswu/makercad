package dlineate

func (s *Sketch) AddCoincidentConstraint(p1 *Element, p2 *Element) *Constraint {
	// If two points are coincident, they are the same point -- make them reference the same element
	if p1.elementType == Point && p2.elementType == Point {
		if p2.element.GetID() == 0 {
			return s.AddCoincidentConstraint(p2, p1)
		}

		p1.element = s.sketch.CombinePoints(p1.element, p2.element)
		p2.element = p1.element
		// These elements must now reference the same constraints
		for _, c := range s.eToC[p2.id] {
			c.replaceElement(p2, p1)
		}
		s.eToC[p1.id] = append(s.eToC[p1.id], s.eToC[p2.id]...)
		delete(s.eToC, p2.id)
		p2.id = p1.id
		return nil
	}
	c := s.AddDistanceConstraint(p1, p2, 0)
	c.constraintType = Coincident
	return c
}

/*
 A point is coincident with a line segment or arc if:
  * The point is coincident with the line or arc /and/
  * The distance from point to the start and end is less than the segment lenth
*/
