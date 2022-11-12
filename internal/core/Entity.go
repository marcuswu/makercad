package core

import "github.com/marcuswu/dlineate"

type Entity interface {
	isConstruction() bool
	setConstruction(bool) Entity
	getElement() *dlineate.Element
}
