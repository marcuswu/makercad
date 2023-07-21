package core

type Line interface {
	Start() *Point
	End() *Point
	IsConstruction() bool
}
