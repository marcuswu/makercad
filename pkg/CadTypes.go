package MakerCad

// TODO: These either need to be eliminated or updated
type Edge interface {
	IsLine() bool
	IsCircle() bool
	IsEllipse() bool
	GetLine() Line
	LineLength() float64
	GetCircle() Circle
	CircleRadius() float64
}

type Point interface {
	GetX() float64
	GetY() float64
	GetZ() float64
	SetX(float64)
	SetY(float64)

	Coincident(...interface{}) Point
	Horizontal(Point) Point
	Vertical(Point) Point
	Distance(...interface{}) Point
	HorizontalDistance(Point, float64) Point
	VerticalDistance(Point, float64) Point
	Construction(bool) Point
	ToString() string
}

type Line interface {
	Horizontal() Line
	Vertical() Line
	Length(float64) Line
	MidPoint(Point) Line
	Symmetric(Point, Point) Line
	Angle(Line, float64) Line
	Construction(bool) Line
	ToString()
	GetStart() Point
	GetEnd() Point
	ToVector() Vector
}

type Arc interface {
	Diameter(float64) Arc
	Construction(bool) Arc
	StartTangent(Line) Arc
	EndTangent(Line) Arc
	GetCenter() Point
	GetStart() Point
	GetEnd() Point
}

type Circle interface {
	Diameter(float64) Circle
	Construction(bool) Circle
	Equal(Circle) Circle
	GetDiameter() float64
	GetCenter() Point
}
