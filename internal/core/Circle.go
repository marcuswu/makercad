package core

type Circle interface {
	Center() *Point
	Radius() float64
	IsConstruction() bool
}
