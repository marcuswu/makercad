package element

import (
	"fmt"
	"math"

	"github.com/marcuswu/dlineate/utils"
)

// SketchLine represents a line in a 2D sketch in the form
// Ax + By + C = 0. A and B are represented as x and y in the BaseElement
type SketchLine struct {
	elementType     Type
	id              uint
	a               float64
	b               float64
	c               float64
	constraintLevel ConstraintLevel
}

// NewSketchLine creates a new SketchLine
func NewSketchLine(id uint, a float64, b float64, c float64) *SketchLine {
	// A & B represent a normal vector for the line. This also determines
	// the direction of the line. C represents a magnitude of the normal
	// vector to reach from origin to the line.
	l := &SketchLine{
		elementType:     Line,
		id:              id,
		a:               a,
		b:               b,
		c:               c,
		constraintLevel: FullyConstrained,
	}
	l.Normalize()
	return l
}

// GetID returns the line element identifier
func (l *SketchLine) GetID() uint { return l.id }

// SetID sets the line element identifier
func (l *SketchLine) SetID(id uint) { l.id = id }

// GetA returns A in the formula Ax + By + C = 0
func (l *SketchLine) GetA() float64 { return l.a }

// GetB returns B in the formula Ax + By + C = 0
func (l *SketchLine) GetB() float64 { return l.b }

// GetC returns c in the formula Ax + By + C = 0
func (l *SketchLine) GetC() float64 { return l.c }

// SetC set the a value for the line (Ax + Bx + C = 0)
func (l *SketchLine) SetA(a float64) { l.a = a }

// SetC set the b value for the line (Ax + Bx + C = 0)
func (l *SketchLine) SetB(b float64) { l.b = b }

// SetC set the c value for the line (Ax + Bx + C = 0)
func (l *SketchLine) SetC(c float64) { l.c = c }

// GetType returns the sketch type
func (l *SketchLine) GetType() Type { return l.elementType }

// Is returns true if the two elements are equal
func (l *SketchLine) Is(o SketchElement) bool {
	return l.id == o.GetID()
}

func (l *SketchLine) Normalize() {
	magnitude := math.Sqrt((l.a * l.a) + (l.b * l.b))
	l.a = l.a / magnitude
	l.b = l.b / magnitude
	l.c = l.c / magnitude
}

// IsEquivalent returns true if the two lines are equivalent
func (l *SketchLine) IsEquivalent(o *SketchLine) bool {
	return utils.StandardFloatCompare(l.a, o.a) == 0 &&
		utils.StandardFloatCompare(l.b, o.b) == 0 &&
		utils.StandardFloatCompare(l.c, o.c) == 0
}

// SquareDistanceTo returns the squared distance to the other element
func (l *SketchLine) SquareDistanceTo(o SketchElement) float64 {
	d := l.DistanceTo(o)

	return d * d
}

func (l *SketchLine) distanceToPoint(x float64, y float64) float64 {
	return (l.a * x) + (l.b * y) + l.c
}

// NearestPoint returns the point on the line nearest the provided point
func (l *SketchLine) NearestPoint(x float64, y float64) *SketchPoint {
	px := (l.b * ((l.b * x) - (l.a * y))) - l.a*l.c
	py := (l.a * ((l.a * y) - (l.b * x))) - l.b*l.c

	return NewSketchPoint(0, px, py)
}

// DistanceTo returns the distance to the other element
func (l *SketchLine) DistanceTo(o SketchElement) float64 {
	switch o.GetType() {
	case Line:
		// Technically I should return 0 if lines aren't parallel
		// Here I am instead comparing min distances to origin
		return l.distanceToPoint(0, 0) - o.(*SketchLine).distanceToPoint(0, 0)
	default:
		return l.distanceToPoint(o.(*SketchPoint).GetX(), o.(*SketchPoint).GetY())
	}
}

// GetOriginDistance returns the distance to the origin for this line
func (l *SketchLine) GetOriginDistance() float64 { return l.distanceToPoint(0, 0) }

// PointNearestOrigin get the point on the line nearest to the origin
func (l *SketchLine) PointNearestOrigin() *SketchPoint {
	if utils.StandardFloatCompare((l.a*l.a)+(l.b*l.b), 0) != 0 {
		l.Normalize()
	}
	return NewSketchPoint(
		0,
		-l.GetC()*l.GetA(),
		-l.GetC()*l.GetB())
}

// TranslateDistance translates the line by a distance along its normal
func (l *SketchLine) TranslateDistance(dist float64) {
	// find point nearest to origin
	l.c = l.TranslatedDistance(dist).GetC()
}

// TranslatedDistance returns the line translated by a distance along its normal
func (l *SketchLine) TranslatedDistance(dist float64) *SketchLine {
	if utils.StandardFloatCompare((l.a*l.a)+(l.b*l.b), 0) != 0 {
		l.Normalize()
	}
	return &SketchLine{Line, l.GetID(), l.GetA(), l.GetB(), l.GetC() - dist, l.constraintLevel}
}

// Translated returns a line translated by an x and y value
func (l *SketchLine) Translated(tx float64, ty float64) *SketchLine {
	if utils.StandardFloatCompare((l.a*l.a)+(l.b*l.b), 0) != 0 {
		l.Normalize()
	}
	pointOnLine := Vector{l.GetA() * -l.GetC(), l.GetB() * -l.GetC()}
	pointOnLine.Translate(tx, ty)
	newC := (-l.GetA() * pointOnLine.GetX()) - (l.GetB() * pointOnLine.GetY())
	// If (A, B) is a unit vector normal to the line,
	// C is the magnitude of the vector to the line,
	// and (tx, ty) is a vector to translate the line,
	// then the dot product of the vectors is the change to C to move the line by tx, ty
	// newC := l.GetC() + (l.GetA() * tx) + (l.GetB() * ty)
	return &SketchLine{Line, l.GetID(), l.GetA(), l.GetB(), newC, l.constraintLevel}
}

// Translate translates the location of this line by an x and y distance
func (l *SketchLine) Translate(tx float64, ty float64) {
	l.c = l.Translated(tx, ty).GetC()
}

// TranslateByElement translates the location of this line by another element
func (l *SketchLine) TranslateByElement(e SketchElement) {
	var point *SketchPoint
	if e.GetType() == Line {
		point = e.(*SketchLine).PointNearestOrigin()
	} else {
		point = e.(*SketchPoint)
	}
	l.Translate(point.GetX(), point.GetY())
}

// ReverseTranslateByElement translates the location of this line by the inverse of another element
func (l *SketchLine) ReverseTranslateByElement(e SketchElement) {
	var point *SketchPoint
	if e.GetType() == Line {
		point = e.(*SketchLine).PointNearestOrigin()
	} else {
		point = e.(*SketchPoint)
	}
	l.Translate(-point.GetX(), -point.GetY())
}

// GetSlope returns the slope of the line (Ax + By + C = 0)
func (l *SketchLine) GetSlope() float64 {
	return -l.GetA() / l.GetB()
}

// AngleTo returns the angle to another vector in radians
func (l *SketchLine) AngleTo(u *Vector) float64 {
	// point [0, -C / B] - point[-C / A, 0]
	lv := &Vector{l.GetB(), -l.GetA()}
	return lv.AngleTo(u)
}

// AngleToLine returns the angle the line needs to rotate to be equivalent to to another line in radians
func (l *SketchLine) AngleToLine(o *SketchLine) float64 {
	lv := &Vector{l.GetB(), -l.GetA()}
	ov := &Vector{o.GetB(), -o.GetA()}
	return lv.AngleTo(ov)
}

// Rotated returns a line representing this line rotated around the origin by angle radians
func (l *SketchLine) Rotated(angle float64) *SketchLine {
	// create vectors with points from the line (x and y intercepts)
	l.Normalize()
	n := &Vector{l.GetA(), l.GetB()}
	n.Rotate(angle)
	return NewSketchLine(l.GetID(), n.GetX(), n.GetY(), l.GetC())
}

// Rotate returns a line representing this line rotated around the origin by angle radians
func (l *SketchLine) Rotate(angle float64) {
	rotated := l.Rotated(angle)
	l.a = rotated.GetA()
	l.b = rotated.GetB()
	l.c = rotated.GetC()
}

// Intersection returns the intersection of two lines
func (l *SketchLine) Intersection(l2 *SketchLine) Vector {
	y := ((l.a * l2.c) - (l.c * l2.a)) / ((l.b * l2.a) - (l.a * l2.b))
	var x float64 = 0.0
	if utils.StandardFloatCompare(l2.a, 0) == 0 {
		x = ((l.b * y) + l.c) / -l.a
	} else {
		x = ((l2.b * y) + l2.c) / -l2.a
	}

	return Vector{x, y}
}

// VectorTo returns a Vector to SketchElement o
func (l *SketchLine) VectorTo(o SketchElement) *Vector {
	var point *SketchPoint
	var myPoint *SketchPoint
	if utils.StandardFloatCompare((l.a*l.a)+(l.b*l.b), 0) != 0 {
		l.Normalize()
	}
	if o.GetType() == Point {
		point = o.(*SketchPoint)
		myPoint = l.NearestPoint(point.GetX(), point.GetY())
	} else {
		oline := o.AsLine()
		if utils.StandardFloatCompare((oline.a*oline.a)+(oline.b*oline.b), 0) != 0 {
			oline.Normalize()
		}
		point = NewSketchPoint(0, oline.a*oline.c, oline.b*oline.c)
		myPoint = NewSketchPoint(0, l.a*l.c, l.b*l.c)
	}

	return &Vector{myPoint.GetX() - point.GetX(), myPoint.GetY() - point.GetY()}
}

// AsPoint returns a SketchElement as a *SketchPoint or nil
func (l *SketchLine) AsPoint() *SketchPoint {
	return nil
}

// AsLine returns a SketchElement as a *SketchLine or nil
func (l *SketchLine) AsLine() *SketchLine {
	return l
}

func (l *SketchLine) ConstraintLevel() ConstraintLevel {
	return l.constraintLevel
}

func (l *SketchLine) SetConstraintLevel(cl ConstraintLevel) {
	l.constraintLevel = cl
}

func (l *SketchLine) String() string {
	return fmt.Sprintf("Line(%d) %fx + %fy + %f = 0", l.id, l.a, l.b, l.c)
}

func (l *SketchLine) ToGraphViz(cId int) string {
	return toGraphViz(l, cId)
}
