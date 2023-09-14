package core

import (
	"github.com/marcuswu/gooccwrapper/topexp"
	"github.com/marcuswu/gooccwrapper/topods"
)

type Shape struct {
	shape topods.Shape
}

func (s *Shape) Faces() []Face {
	faces := make([]Face, 0)
	for ex := topexp.NewExplorer(s.shape, topexp.Face); ex.More(); ex.Next() {
		if ex.Depth() > 1 {
			continue
		}

		faces = append(faces, Face{topods.NewFaceFromRef(topods.TopoDSFace(ex.Current().Shape))})
	}

	return faces
}
