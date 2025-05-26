package makercad

import (
	"github.com/marcuswu/gooccwrapper/topexp"
	"github.com/marcuswu/gooccwrapper/topods"
	"github.com/marcuswu/gooccwrapper/toptools"
)

type Shape struct {
	Shape topods.Shape
}

func (s Shape) Faces() ListOfFace {
	faces := make([]*Face, 0)
	for ex := topexp.NewExplorer(s.Shape, topexp.Face); ex.More(); ex.Next() {
		if ex.Depth() > 1 {
			continue
		}

		faces = append(faces, &Face{topods.NewFaceFromRef(topods.TopoDSFace(ex.Current().Shape))})
	}

	return faces
}

type ListOfShape []Shape

func (l ListOfShape) ToCascadeList() toptools.ListOfShape {
	list := toptools.NewListOfShape()
	for i := range l {
		list.Append(l[i].Shape)
	}
	return list
}
