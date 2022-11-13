package element

import (
	"fmt"
	"math"
)

// SketchPoint represents a point in a 2D Sketch
type SketchPoint struct {
	Vector
	elementType     Type
	id              uint
	constraintLevel ConstraintLevel
}

// SetID sets the id of the element
func (p *SketchPoint) SetID(id uint) {
	p.id = id
}

// GetID gets the id of the element
func (p *SketchPoint) GetID() uint {
	return p.id
}

// GetX gets the x value of the point
func (p *SketchPoint) GetX() float64 { return p.X }

// GetY gets the x value of the point
func (p *SketchPoint) GetY() float64 { return p.Y }

// GetType gets the type of the element
func (p *SketchPoint) GetType() Type {
	return p.elementType
}

// TranslateByElement translates coordinates by another element's coordinates
func (p *SketchPoint) TranslateByElement(e SketchElement) {
	var point *SketchPoint
	if e.GetType() == Point {
		point = e.(*SketchPoint)
	} else {
		point = e.(*SketchLine).PointNearestOrigin()
	}

	p.Translate(point.GetX(), point.GetY())
}

// ReverseTranslateByElement translates coordinates by the inverse of another element's coordinates
func (p *SketchPoint) ReverseTranslateByElement(e SketchElement) {
	var point *SketchPoint
	if e.GetType() == Point {
		point = e.(*SketchPoint)
	} else {
		point = e.(*SketchLine).PointNearestOrigin()
	}

	p.Translate(-point.GetX(), -point.GetY())
}

// Is returns true if the two elements are equal
func (p *SketchPoint) Is(o SketchElement) bool {
	return p.id == o.GetID()
}

// SquareDistanceTo returns the squared distance to the other element
func (p *SketchPoint) SquareDistanceTo(o SketchElement) float64 {
	if o.GetType() == Line {
		d := o.(*SketchLine).DistanceTo(p)
		return d * d
	}
	a := p.X - o.(*SketchPoint).GetX()
	b := p.Y - o.(*SketchPoint).GetY()

	return (a * a) + (b * b)
}

// DistanceTo returns the distance to the other element
func (p *SketchPoint) DistanceTo(o SketchElement) float64 {
	return math.Sqrt(p.SquareDistanceTo(o))
}

// SketchPoint represents a point in a 2D sketch

// NewSketchPoint creates a new SketchPoint
func NewSketchPoint(id uint, x float64, y float64) *SketchPoint {
	return &SketchPoint{
		Vector:          Vector{x, y},
		elementType:     Point,
		id:              id,
		constraintLevel: FullyConstrained,
	}
}

func SketchPointFromVector(id uint, v Vector) *SketchPoint {
	return &SketchPoint{
		Vector:          v,
		elementType:     Point,
		id:              id,
		constraintLevel: FullyConstrained,
	}
}

// VectorTo returns a Vector to SketchElement o
func (p *SketchPoint) VectorTo(o SketchElement) *Vector {
	var point *SketchPoint
	if o.GetType() == Point {
		point = o.(*SketchPoint)
	} else {
		point = o.(*SketchLine).NearestPoint(p.GetX(), p.GetY())
	}

	return &Vector{p.GetX() - point.GetX(), p.GetY() - point.GetY()}
}

// AsPoint returns a SketchElement as a *SketchPoint or nil
func (p *SketchPoint) AsPoint() *SketchPoint {
	return p
}

// AsLine returns a SketchElement as a *SketchLine or nil
func (p *SketchPoint) AsLine() *SketchLine {
	return nil
}

func (p *SketchPoint) ConstraintLevel() ConstraintLevel {
	return p.constraintLevel
}

func (p *SketchPoint) SetConstraintLevel(cl ConstraintLevel) {
	p.constraintLevel = cl
}

func (p *SketchPoint) String() string {
	return fmt.Sprintf("Point(%d) (%f, %f)", p.GetID(), p.X, p.Y)
}

func (p *SketchPoint) ToGraphViz(cId int) string {
	return toGraphViz(p, cId)
}
