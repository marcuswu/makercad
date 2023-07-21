package core

import (
	"fmt"
)

type Point interface {
	X() float64
	Y() float64
	IsConstruction() bool
}

func (p *Point) ToString() string {
	return fmt.Sprintf("%f, %f", p.X, p.Y)
}
