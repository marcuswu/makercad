package sketch

import (
	"slices"

	"github.com/marcuswu/gooccwrapper/gp"
	"github.com/marcuswu/gooccwrapper/topods"
)

type Vertex struct {
	Vertex topods.Vertex
}
type ListOfVertex []*Vertex

type VertexFilter func(*Vertex) bool
type VertexSorter func(a, b *Vertex) int

func NewVertexFromRef(shape topods.Shape) *Vertex {
	return &Vertex{topods.NewVertexFromRef(topods.TopoDSVertex(shape.Shape))}
}

func (l ListOfVertex) First(filter VertexFilter) *Vertex {
	for _, edge := range l {
		if filter(edge) {
			return edge
		}
	}
	return nil
}

func (l ListOfVertex) Matching(filter VertexFilter) ListOfVertex {
	newList := make(ListOfVertex, 0, len(l))
	for _, edge := range l {
		if filter(edge) {
			newList = append(newList, edge)
		}
	}
	return newList
}

func (l ListOfVertex) Sort(sorter VertexSorter) {
	slices.SortFunc(l, sorter)
}

func (v Vertex) ToPoint() gp.Pnt {
	return v.Vertex.Pnt()
}
