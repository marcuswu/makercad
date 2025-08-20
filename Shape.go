package makercad

import (
	"github.com/marcuswu/gooccwrapper/brepalgoapi"
	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/gp"
	"github.com/marcuswu/gooccwrapper/topexp"
	"github.com/marcuswu/gooccwrapper/topods"
	"github.com/marcuswu/gooccwrapper/toptools"
)

// Shape represents an OpenCascade topods.Shape. Usually in MakerCad it's a 3D shape.
type Shape struct {
	Shape topods.Shape
}

// Retrieve all faces part of this shape. Returns a filterable and sortable list of faces
func (s Shape) Faces() ListOfFace {
	faces := make([]*Face, 0)
	for ex := topexp.NewExplorer(s.Shape, topexp.Face); ex.More(); ex.Next() {
		faces = append(faces, &Face{topods.NewFaceFromRef(topods.TopoDSFace(ex.Current().Shape))})
	}

	return faces
}

// ListOfShape is simply a list of filterable and sortable shapes.
type ListOfShape []Shape

// ToCascadeList converts this slice of Shapes to something OpenCascade can use
func (l ListOfShape) ToCascadeList() toptools.ListOfShape {
	list := toptools.NewListOfShape()
	for i := range l {
		list.Append(l[i].Shape)
	}
	return list
}

// Remove performs a boolean subtraction of the tools from this Shape
func (s Shape) Remove(tools ListOfShape) (*CadOperation, error) {
	operation := brepalgoapi.NewCut().ToBooleanOperation()
	arguments := make(ListOfShape, 1)
	arguments[0] = s

	operation.SetTools(tools.ToCascadeList())
	operation.SetArguments(arguments.ToCascadeList())
	operation.Build()

	return NewCadOperation(tools, &operation), nil
}

// Transform performs a rigid body transform (rotate & translate) of this Shape
func (s Shape) Transform(trsf gp.Trsf) Shape {
	transform := brepbuilderapi.NewTransform(s.Shape, trsf)
	return Shape{transform.Shape()}
}
