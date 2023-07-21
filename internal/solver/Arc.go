package solver

import (
	"github.com/marcuswu/dlineate"
)

type Arc struct {
	dlineate.Element
	solver         SketchSolver
	center         *Point
	start          *Point
	end            *Point
	isConstruction bool
}

func (a *Arc) getElement() *dlineate.Element {
	return &a.Element
}

func (a *Arc) Center() *Point {
	return a.center
}

func (a *Arc) Start() *Point {
	return a.start
}

func (a *Arc) End() *Point {
	return a.end
}

func (a *Arc) IsConstruction() bool {
	return a.isConstruction
}
