package core

import "github.com/marcuswu/dlineate"

type Arc struct {
	dlineate.Element
	solver         SketchSolver
	Center         *Point
	Start          *Point
	End            *Point
	IsConstruction bool
}
