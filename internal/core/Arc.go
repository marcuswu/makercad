package core

type Arc interface {
	Center() *Point
	Start() *Point
	End() *Point
	IsConstruction() bool
}
