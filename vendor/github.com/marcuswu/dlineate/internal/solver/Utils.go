package solver

import (
	"fmt"

	"github.com/marcuswu/dlineate/internal/constraint"
	el "github.com/marcuswu/dlineate/internal/element"
)

// SolveState The state of the sketch graph
type SolveState uint

// SolveState constants
const (
	None SolveState = iota
	OverConstrained
	NonConvergent
	Solved
)

func (ss SolveState) String() string {
	switch ss {
	case OverConstrained:
		return "over constrained"
	case NonConvergent:
		return "non-convergent"
	case Solved:
		return "solved"
	default:
		return fmt.Sprintf("%d", int(ss))
	}
}

func typeCounts(c1 *constraint.Constraint, c2 *constraint.Constraint) (int, int) {
	numPoints := 0
	numLines := 0
	elements := []el.SketchElement{c1.Element1, c1.Element2, c2.Element1, c2.Element2}

	for _, element := range elements {
		if element.GetType() == el.Point {
			numPoints++
		} else {
			numLines++
		}
	}

	return numPoints, numLines
}
