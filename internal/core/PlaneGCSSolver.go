package core

import "libmakercad/third_party/planegcs"

type PlaneGCSSolver struct {
	system planegcs.System
}

func NewPlaneGCSSolver() *PlaneGCSSolver {
	return &PlaneGCSSolver{planegcs.NewSystem()}
}

func (s *PlaneGCSSolver) CreatePoint(x float64, y float64) *Point {
	return NewPoint(x, y)
}
