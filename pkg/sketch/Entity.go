package sketch

import (
	"fmt"

	"github.com/marcuswu/dlineate"
)

type Entity interface {
	edger
	fmt.Stringer
	getElement() *dlineate.Element
	UpdateFromValues()
	IsConstruction() bool
	SetConstruction(bool)
}

func AsArc(e Entity) *Arc {
	c, ok := e.(*Arc)
	if !ok {
		return nil
	}
	return c
}

func AsCircle(e Entity) *Circle {
	c, ok := e.(*Circle)
	if !ok {
		return nil
	}
	return c
}

func AsLine(e Entity) *Line {
	c, ok := e.(*Line)
	if !ok {
		return nil
	}
	return c
}

func AsPoint(e Entity) *Point {
	c, ok := e.(*Point)
	if !ok {
		return nil
	}
	return c
}
