package makercad

import "github.com/marcuswu/gooccwrapper/brepalgoapi"

type MergeType int

const (
	MergeTypeNew MergeType = iota
	MergeTypeAdd
	MergeTypeRemove
	MergeTypeMax
)

type CadOperation struct {
	shapes    ListOfShape
	operation *brepalgoapi.Boolean
}

func NewCadOperation(e ListOfShape, op *brepalgoapi.Boolean) *CadOperation {
	return &CadOperation{shapes: e, operation: op}
}

func (o *CadOperation) Shape() Shape {
	shape := Shape{}
	if len(o.shapes) > 0 {
		shape = o.shapes[0]
	}
	if o.operation != nil {
		shape = Shape{o.operation.Shape()}
	}

	return shape
}
