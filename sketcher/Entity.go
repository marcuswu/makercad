package sketcher

import (
	"fmt"

	"github.com/marcuswu/dlineate"
)

// Entity is any element that can be used in a [Sketch]
// Entities can be converted to [Edge]s
// Entities can be string formatted
// Entities can have their values updated from a [Sketch]
type Entity interface {
	edger
	fmt.Stringer
	getElement() *dlineate.Element
	UpdateFromValues()
	IsConstruction() bool
	SetConstruction(bool)
	IsConnectedTo(Entity) bool
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
