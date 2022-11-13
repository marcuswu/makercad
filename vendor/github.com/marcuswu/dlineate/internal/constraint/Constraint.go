package constraint

import (
	"fmt"
	"math"

	el "github.com/marcuswu/dlineate/internal/element"
	"github.com/marcuswu/dlineate/utils"
	"github.com/rs/zerolog"
)

// Type of a Constraint(Distance or Angle)
type Type uint

// ConstraintType constants
const (
	Distance Type = iota
	Angle
)

func (t Type) String() string {
	switch t {
	case Distance:
		return "Distance"
	case Angle:
		return "Angle"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

// Constraint interface representing a constraint
/*type Constraint interface {
	SetID(uint)
	GetID() uint
	GetType() ConstraintType
	GetValue() float64
	UpdateValue(float64)
	HasElementID(uint)
	First() SketchElement
	Second() SketchElement
	EquationCount() uint
	ValueCount() uint
	FillValues([]float64)
	CheckSolution([]float64, float64)
	Equals(Constraint) bool

	Calculate()
	Check()
}*/

// Constraint Represents a 2D constraint
type Constraint struct {
	id       uint
	Type     Type
	Value    float64
	Element1 el.SketchElement
	Element2 el.SketchElement
	Solved   bool
}

// GetID returns the constraint identifier
func (c *Constraint) GetID() uint {
	return c.id
}

// GetValue returns the constraint's value
func (c *Constraint) GetValue() float64 {
	return c.Value
}

// UpdateValue updates the constraint's value
func (c *Constraint) UpdateValue(v float64) {
	c.Value = v
}

// HasElementID returns whether an element with the passed ID
// exists in this constraint
func (c *Constraint) HasElementID(eID uint) bool {
	return c.Element1.GetID() == eID || c.Element2.GetID() == eID
}

func (c *Constraint) HasElements(ids ...uint) bool {
	for _, id := range ids {
		if id != c.Element1.GetID() && id != c.Element2.GetID() {
			return false
		}
	}

	return true
}

// First returns the first element in the constraint
func (c *Constraint) First() el.SketchElement {
	return c.Element1
}

// Second returns the second element in the constraint
func (c *Constraint) Second() el.SketchElement {
	return c.Element2
}

func (c *Constraint) ElementIDs() []uint {
	return []uint{c.Element1.GetID(), c.Element2.GetID()}
}

func (c *Constraint) Element(this uint) (el.SketchElement, bool) {
	if this == c.Element1.GetID() {
		return c.Element1, true
	}
	return c.Element2, this == c.Element2.GetID()
}

func (c *Constraint) Other(this uint) (el.SketchElement, bool) {
	if this == c.Element1.GetID() {
		return c.Element2, true
	}
	return c.Element1, this == c.Element2.GetID()
}

func (c *Constraint) Shared(o *Constraint) (el.SketchElement, bool) {
	if o.HasElementID(c.Element1.GetID()) {
		return c.Element1, true
	}
	if o.HasElementID(c.Element2.GetID()) {
		return c.Element2, true
	}

	return nil, false
}

func (c *Constraint) IsMet() bool {
	current := c.Element1.DistanceTo(c.Element2)
	if c.Type == Angle {
		current = c.Element1.AsLine().AngleToLine(c.Element2.AsLine())
	}

	comparison := utils.StandardFloatCompare(math.Abs(current), math.Abs(c.Value))
	if comparison != 0 {
		utils.Logger.Trace().
			Float64("value 1", math.Abs(current)).
			Float64("value 2", math.Abs(c.Value)).
			Msgf("Comparing values")
		c.Solved = false
	} else {
		c.Solved = true
	}

	return c.Solved
}

func (c *Constraint) String() string {
	units := ""
	if c.Type == Angle {
		units = " rad"
	}
	return fmt.Sprintf("Constraint(%d) type: %v, e1: %d, e2: %d, v: %f%s", c.GetID(), c.Type, c.Element1.GetID(), c.Element2.GetID(), c.Value, units)
}

func (c *Constraint) ToGraphViz(cId int) string {
	if cId < 0 {
		return fmt.Sprintf("\t%d -- %d [label=\"%v (%d)\"]\n", c.Element1.GetID(), c.Element2.GetID(), c.Type, c.id)
	}
	return fmt.Sprintf("\t\"%d-%d\" -- \"%d-%d\" [label=\"%v (%d)\"]\n", cId, c.Element1.GetID(), cId, c.Element2.GetID(), c.Type, c.id)
}

// Equals returns whether two constraints are equal
func (c *Constraint) Equals(o Constraint) bool {
	return c.id == o.GetID()
}

// NewConstraint creates a new constraint
func NewConstraint(id uint, constraintType Type, a el.SketchElement, b el.SketchElement, v float64, solved bool) *Constraint {
	return &Constraint{
		id:       id,
		Type:     constraintType,
		Value:    v,
		Element1: a,
		Element2: b,
		Solved:   false,
	}
}

// CopyConstraint creates a deep copy of a Constraint
func CopyConstraint(c *Constraint) *Constraint {
	return NewConstraint(
		c.GetID(),
		c.Type,
		el.CopySketchElement(c.Element1),
		el.CopySketchElement(c.Element2),
		c.Value,
		c.Solved,
	)
}

type ConstraintList []*Constraint

func (cl ConstraintList) Len() int           { return len(cl) }
func (cl ConstraintList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
func (cl ConstraintList) Less(i, j int) bool { return cl[i].id < cl[j].id }

func (l ConstraintList) MarshalZerologArray(a *zerolog.Array) {
	for _, c := range l {
		a.Uint(c.GetID())
	}
}
