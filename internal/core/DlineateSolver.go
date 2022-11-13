package core

import "github.com/marcuswu/dlineate"

type DlineateSolver struct {
	system *dlineate.Sketch
}

func NewDlineateSolver() *DlineateSolver {
	return &DlineateSolver{dlineate.NewSketch()}
}

func (s *DlineateSolver) CreatePoint(x float64, y float64) *Point {
	return &Point{Element: *s.system.AddPoint(x, y), X: x, Y: y, IsConstruction: false}
}

func (s *DlineateSolver) CreateLine(p1 *Point, p2 *Point) *Line {
	return nil
}

func (s *DlineateSolver) ToFace() *Face {
	return nil
}
