package core

import "github.com/marcuswu/dlineate"

type Entity interface {
	edger
	getElement() *dlineate.Element
	isConstruction() bool
}
