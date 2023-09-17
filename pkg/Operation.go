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
	shape     Shape
	operation *brepalgoapi.Boolean
}

func NewCadOperation(e Shape, op *brepalgoapi.Boolean) *CadOperation {
	return &CadOperation{shape: e, operation: op}
}

func (o *CadOperation) Shape() Shape {
	shape := o.shape
	if o.operation != nil {
		shape = Shape{o.operation.Shape()}
	}

	return shape
}
