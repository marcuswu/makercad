package core

import "github.com/marcuswu/dlineate"

type Circle struct {
	dlineate.Element
	solver         SketchSolver
	Center         *Point
	Radius         float64
	IsConstruction bool
}
