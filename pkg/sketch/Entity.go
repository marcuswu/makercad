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
}
