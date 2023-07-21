package solver

import "github.com/marcuswu/dlineate"

type Entity interface {
	getElement() *dlineate.Element
}
