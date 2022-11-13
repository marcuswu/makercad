package dlineate

import (
	"errors"
	"fmt"
	"math"

	c "github.com/marcuswu/dlineate/internal/constraint"
	el "github.com/marcuswu/dlineate/internal/element"
	"github.com/tdewolff/canvas"
)

// Type of a Constraint(Distance or Angle)
type ElementType uint

// ElementType constants
const (
	Point ElementType = iota
	Axis
	Line
	Circle
	Arc
)

func (et ElementType) String() string {
	switch et {
	case Point:
		return "Point"
	case Axis:
		return "Axis"
	case Line:
		return "Line"
	case Circle:
		return "Circle"
	case Arc:
		return "Arc"
	default:
		return fmt.Sprintf("%d", int(et))
	}
}

type Element struct {
	id          uint
	values      []float64
	elementType ElementType
	constraints []*c.Constraint
	element     el.SketchElement
	children    []*Element
	isChild     bool
	valuePass   int
}

func emptyElement() *Element {
	ec := new(Element)
	ec.values = make([]float64, 0)
	ec.constraints = make([]*c.Constraint, 0)
	ec.children = make([]*Element, 0)
	ec.isChild = false
	ec.valuePass = 0
	return ec
}

func (e *Element) valuesFromSketch(s *Sketch) error {
	switch e.elementType {
	case Point:
		p := e.element.AsPoint()
		e.values[0] = p.GetX()
		e.values[1] = p.GetY()
	case Axis:
		p := e.element.AsLine()
		e.values[0] = p.GetA()
		e.values[1] = p.GetB()
		e.values[2] = p.GetC()
	case Line:
		p1 := e.children[0].element.AsPoint()
		p2 := e.children[1].element.AsPoint()
		e.values[0] = p1.GetX()
		e.values[1] = p1.GetY()
		e.values[2] = p2.GetX()
		e.values[3] = p2.GetY()
	case Circle:
		/*
			Circle radius is determined either by
			  * a distance constraint against the Circle
			  * a coincident constraint against a Circle with the location of the center constrained
		*/
		var err error = nil
		c := e.children[0].element.AsPoint()
		e.values[0] = c.GetX()
		e.values[1] = c.GetY()
		// find distance constraint on e
		constraint, err := s.findConstraint(Distance, e)
		if err != nil {
			constraint, err = s.findConstraint(Coincident, e)
		}
		if err != nil {
			return err
		}
		e.values[2], err = e.getCircleRadius(constraint)
		if err != nil {
			return err
		}
	case Arc:
		center := e.children[0].element.AsPoint()
		start := e.children[1].element.AsPoint()
		end := e.children[2].element.AsPoint()
		e.values[0] = center.GetX()
		e.values[1] = center.GetY()
		e.values[2] = start.GetX()
		e.values[3] = start.GetY()
		e.values[4] = end.GetX()
		e.values[5] = end.GetY()
	}
	e.valuePass = s.passes

	return nil
}

func (e *Element) getCircleRadius(c *Constraint) (float64, error) {
	if e.elementType != Circle {
		return 0, errors.New("can't return radius for a non-circle")
	}
	if c.constraintType == Distance && len(c.elements) == 1 && c.elements[0].id == e.id {
		return c.dataValue, nil
	}
	if c.constraintType == Coincident {
		constraint := c.constraints[0]
		other := constraint.Element1
		if other == e.children[0].element {
			other = constraint.Element2
		}

		return other.DistanceTo(e.children[0].element.AsPoint()), nil
	}

	return 0, errors.New("Constraint type for circle radius must be Distance or Coincident")
}

func (e *Element) Values(s *Sketch) []float64 {
	if e.valuePass != s.passes {
		e.valuesFromSketch(s)
	}
	return e.values
}

func (e *Element) ConstraintLevel() el.ConstraintLevel {
	level := e.element.ConstraintLevel()
	var childLevel el.ConstraintLevel
	for _, c := range e.children {
		childLevel = c.element.ConstraintLevel()
		if childLevel < level {
			level = childLevel
		}
	}
	return level
}

func (e *Element) minMaxXY() (float64, float64, float64, float64) {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	maxX := math.MaxFloat64 * -1
	maxY := math.MaxFloat64 * -1

	switch e.elementType {
	case Point:
		if e.values[0] < minX {
			minX = e.values[0]
		}
		if e.values[0] > maxX {
			maxX = e.values[0]
		}
		if e.values[1] < minX {
			minY = e.values[1]
		}
		if e.values[1] > maxX {
			maxY = e.values[1]
		}
	case Line:
		if e.values[0] < minX {
			minX = e.values[0]
		}
		if e.values[0] > maxX {
			maxX = e.values[0]
		}
		if e.values[1] < minX {
			minY = e.values[1]
		}
		if e.values[1] > maxX {
			maxY = e.values[1]
		}
		if e.values[2] < minX {
			minX = e.values[2]
		}
		if e.values[2] > maxX {
			maxX = e.values[2]
		}
		if e.values[3] < minX {
			minY = e.values[3]
		}
		if e.values[3] > maxX {
			maxY = e.values[3]
		}
	case Circle:
		size := e.values[2]
		if e.values[0]-size < minX {
			minX = e.values[0] - size
		}
		if e.values[0]+size > maxX {
			maxX = e.values[0] + size
		}
		if e.values[1]-size < minY {
			minY = e.values[1] - size
		}
		if e.values[1]+size > maxY {
			maxY = e.values[1] + size
		}
	case Arc:
		if e.values[0] < minX {
			minX = e.values[0]
		}
		if e.values[0] > maxX {
			maxX = e.values[0]
		}
		if e.values[1] < minY {
			minY = e.values[1]
		}
		if e.values[1] > maxY {
			maxY = e.values[1]
		}
		if e.values[2] < minX {
			minX = e.values[2]
		}
		if e.values[2] > maxX {
			maxX = e.values[2]
		}
		if e.values[3] < minY {
			minY = e.values[3]
		}
		if e.values[3] > maxY {
			maxY = e.values[3]
		}
		if e.values[4] < minX {
			minX = e.values[4]
		}
		if e.values[4] > maxX {
			maxX = e.values[4]
		}
		if e.values[5] < minY {
			minY = e.values[5]
		}
		if e.values[5] > maxY {
			maxY = e.values[5]
		}
	}
	return minX, minY, maxX, maxY
}

func (e *Element) DrawToSVG(s *Sketch, ctx *canvas.Context, mult float64) {
	ctx.StrokeColor = canvas.Blue
	if e.elementType != Axis && e.ConstraintLevel() == el.FullyConstrained {
		ctx.StrokeColor = canvas.Black
	}
	if e.elementType != Axis && e.ConstraintLevel() == el.OverConstrained {
		ctx.StrokeColor = canvas.Red
	}
	ctx.StrokeWidth = 0.5
	switch e.elementType {
	case Point:
		// May want to draw a small filled circle
		ctx.MoveTo(e.values[0]*mult+0.5, e.values[1]*mult)
		ctx.Arc(0.5, 0.5, 0, 0, 360)
	case Line:
		x1 := e.values[0] * mult
		y1 := e.values[1] * mult
		ctx.MoveTo(x1, y1)
		x2 := e.values[2] * mult
		y2 := e.values[3] * mult
		ctx.LineTo(x2, y2)
	case Circle:
		cx := e.values[0] * mult
		cy := e.values[1] * mult
		// find distance constraint on e
		r := e.values[2] * mult
		ctx.MoveTo(cx, cy)
		ctx.Arc(r, r, 0, 0, 360)
	case Arc:
		cx := e.values[0] * mult
		cy := e.values[1] * mult
		sx := e.values[2] * mult
		sy := e.values[3] * mult
		ex := e.values[4] * mult
		ey := e.values[5] * mult
		r := math.Sqrt(math.Pow(sx-cx, 2) + math.Pow(sy-cy, 2))
		svx := sx - cx
		svy := sy - cy
		evx := ex - cx
		evy := ey - cy
		theta0 := math.Atan2(svx, svy)
		theta1 := math.Atan2(evx, evy)
		dot := evx*svx + evy*svy
		det := evx*svy - evy*svx
		angle := math.Atan2(det, dot)
		large := false
		if angle > math.Pi {
			large = true
		}

		sweep := theta1 < theta0
		ctx.MoveTo(sx, sy)
		ctx.ArcTo(r, r, angle, large, sweep, ex, ey)
	}
	ctx.Stroke()
	e.valuePass = s.passes
}

func (e *Element) Center() *Element {
	if e.elementType != Arc && e.elementType != Circle {
		return nil
	}
	return e.children[0]
}

func (e *Element) Start() *Element {
	if e.elementType == Arc {
		return e.children[1]
	}
	if e.elementType != Line {
		return nil
	}
	return e.children[0]
}

func (e *Element) End() *Element {
	if e.elementType == Arc {
		return e.children[2]
	}
	if e.elementType != Line {
		return nil
	}
	return e.children[1]
}

func (e *Element) String() string {
	return fmt.Sprintf("Element type %v, internal element: %v", e.elementType, e.element)
}
