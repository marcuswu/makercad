package core

import "github.com/marcuswu/gooccwrapper/topods"

type edger interface {
	MakeEdge() *Edge
}

type Edge struct {
	edge topods.Edge
}
