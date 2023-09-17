package sketch

import "github.com/marcuswu/dlineate"

type Entity interface {
	edger
	getElement() *dlineate.Element
	IsConstruction() bool
}
