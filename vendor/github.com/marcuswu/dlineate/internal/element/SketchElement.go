package element

import "fmt"

// Type of a SketchElement (Point or Line)
type Type uint

// SolveState constants
const (
	Point Type = iota
	Line
)

func (t Type) String() string {
	switch t {
	case Point:
		return "Point"
	case Line:
		return "Line"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

type ConstraintLevel uint

const (
	OverConstrained ConstraintLevel = iota
	UnderConstrained
	FullyConstrained
)

func (cl ConstraintLevel) String() string {
	switch cl {
	case OverConstrained:
		return "over constrained"
	case UnderConstrained:
		return "under constrained"
	case FullyConstrained:
		return "fully constrained"
	default:
		return fmt.Sprintf("%d", int(cl))
	}
}

// SketchElement A 2D element within a Sketch
type SketchElement interface {
	fmt.Stringer
	SetID(uint)
	GetID() uint
	GetType() Type
	AngleTo(*Vector) float64
	Translate(tx float64, ty float64)
	TranslateByElement(SketchElement)
	ReverseTranslateByElement(SketchElement)
	Rotate(tx float64)
	Is(SketchElement) bool
	SquareDistanceTo(SketchElement) float64
	DistanceTo(SketchElement) float64
	VectorTo(SketchElement) *Vector
	AsPoint() *SketchPoint
	AsLine() *SketchLine
	ConstraintLevel() ConstraintLevel
	SetConstraintLevel(ConstraintLevel)
	ToGraphViz(cId int) string
}

// List is a list of SketchElements
type List []SketchElement

func (e List) Len() int           { return len(e) }
func (e List) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e List) Less(i, j int) bool { return e[i].GetID() < e[j].GetID() }

// CopySketchElement creates a deep copy of a SketchElement
func CopySketchElement(e SketchElement) SketchElement {
	var n SketchElement
	if e.GetType() == Point {
		p := e.(*SketchPoint)
		n = NewSketchPoint(e.GetID(), p.GetX(), p.GetY())
		n.SetConstraintLevel(e.ConstraintLevel())
		return n
	}
	l := e.(*SketchLine)
	n = NewSketchLine(l.GetID(), l.GetA(), l.GetB(), l.GetC())
	n.SetConstraintLevel(e.ConstraintLevel())
	return n
}

func toGraphViz(e SketchElement, cId int) string {
	if cId < 0 {
		return fmt.Sprintf("\t%d\n", e.GetID())
	}
	return fmt.Sprintf("\"%d-%d\" [label=%d]\n", cId, e.GetID(), e.GetID())
}

// IdentityMap is a map of id to SketchElement
type IdentityMap = map[uint]SketchElement
