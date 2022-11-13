package dlineate

import (
	"fmt"

	c "github.com/marcuswu/dlineate/internal/constraint"
	"github.com/marcuswu/dlineate/utils"
)

// Type of a Constraint(Distance or Angle)
type ConstraintType uint

// ElementType constants
const (
	Coincident ConstraintType = iota
	Distance
	Angle
	Perpendicular
	Parallel
	Tangent

	// Two pass constraints
	Ratio
	Midpoint
)

func (t ConstraintType) String() string {
	switch t {
	case Coincident:
		return "Coincident"
	case Distance:
		return "Distance"
	case Angle:
		return "Angle"
	case Perpendicular:
		return "Perpendicular"
	case Parallel:
		return "Parallel"
	case Tangent:
		return "Tangent"
	case Ratio:
		return "Ratio"
	case Midpoint:
		return "Midpoint"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

type ConstraintState uint

const (
	Unresolved ConstraintState = iota
	Resolved
	Solved
)

func (t ConstraintState) String() string {
	switch t {
	case Unresolved:
		return "Unresolved"
	case Resolved:
		return "Resolved"
	case Solved:
		return "Solved"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

type Constraint struct {
	constraints    []*c.Constraint
	elements       []*Element
	constraintType ConstraintType
	state          ConstraintState
	dataValue      float64
}

func emptyConstraint() *Constraint {
	ec := new(Constraint)
	ec.constraints = make([]*c.Constraint, 0)
	ec.elements = make([]*Element, 0)
	ec.state = Unresolved
	return ec
}

func (c *Constraint) replaceElement(from *Element, to *Element) {
	for i, e := range c.elements {
		if e.id == from.id {
			c.elements[i] = to
		}
	}
}

func (c *Constraint) checkSolved() bool {
	// solved := true
	// if len(c.constraints) == 0 {
	// 	solved = true
	// }
	solved := c.state >= Resolved
	for _, constraint := range c.constraints {
		utils.Logger.Trace().
			Uint("constraint", constraint.GetID()).
			Bool("state", constraint.Solved).
			Msg("Constraint solve state")
		solved = solved && constraint.Solved
	}
	if solved {
		c.state = Solved
	}
	utils.Logger.Debug().
		Str("type", c.constraintType.String()).
		Str("state", c.state.String()).
		Msg("Constraint solve state")

	return c.state == Solved
}

/*

One Pass Constraints
-------------
Distance constraint -- line segment, between elements, radius
Coincident constraint -- points, point & line, point & curve, line & curve
Angle -- two lines
Perpendicular -- two lines
Parallel -- two lines

Two Pass Constraints
-------------
Equal constraint -- 2nd pass constraint
Distance ratio constraint -- 2nd pass constraint
Midpoint -- 2nd pass constraint (equal distances to either end of the line or arc)
Tangent -- line and curve
Symmetric -- TODO

*/
