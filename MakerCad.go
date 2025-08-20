package makercad

import (
	"errors"

	"github.com/marcuswu/makercad/sketcher"

	"github.com/marcuswu/gooccwrapper/brep"
	"github.com/marcuswu/gooccwrapper/brepalgoapi"
	"github.com/marcuswu/gooccwrapper/brepfilletapi"
	"github.com/marcuswu/gooccwrapper/brepmesh"
	"github.com/marcuswu/gooccwrapper/brepprimapi"
	"github.com/marcuswu/gooccwrapper/gp"
	"github.com/marcuswu/gooccwrapper/stepcontrol"
	"github.com/marcuswu/gooccwrapper/stlapi"
	"github.com/marcuswu/gooccwrapper/topods"
	"github.com/marcuswu/gooccwrapper/toptools"
)

type ExportQuality int

const (
	QualityVeryLow ExportQuality = iota
	QualityLow
	QualityMedium
	QualityHigh
)

const StepExportSuccess = 1

type MakerCad struct {
	sketches    []*Sketch
	FrontPlane  *sketcher.PlaneParameters
	BackPlane   *sketcher.PlaneParameters
	TopPlane    *sketcher.PlaneParameters
	BottomPlane *sketcher.PlaneParameters
	LeftPlane   *sketcher.PlaneParameters
	RightPlane  *sketcher.PlaneParameters
}

func NewMakerCad() *MakerCad {
	return &MakerCad{
		sketches: make([]*Sketch, 0),
		FrontPlane: sketcher.NewPlaneParametersFromVectors(
			sketcher.NewVectorFromValues(0, 0, 0),
			sketcher.NewVectorFromValues(0, -1, 0),
			sketcher.NewVectorFromValues(1, 0, 0),
		),
		BackPlane: sketcher.NewPlaneParametersFromVectors(
			sketcher.NewVectorFromValues(0, 0, 0),
			sketcher.NewVectorFromValues(0, 1, 0),
			sketcher.NewVectorFromValues(-1, 0, 0),
		),
		TopPlane: sketcher.NewPlaneParametersFromVectors(
			sketcher.NewVectorFromValues(0, 0, 0),
			sketcher.NewVectorFromValues(0, 0, 1),
			sketcher.NewVectorFromValues(1, 0, 0),
		),
		BottomPlane: sketcher.NewPlaneParametersFromVectors(
			sketcher.NewVectorFromValues(0, 0, 0),
			sketcher.NewVectorFromValues(0, 0, -1),
			sketcher.NewVectorFromValues(1, 0, 0),
		),
		LeftPlane: sketcher.NewPlaneParametersFromVectors(
			sketcher.NewVectorFromValues(0, 0, 0),
			sketcher.NewVectorFromValues(-1, 0, 0),
			sketcher.NewVectorFromValues(0, -1, 0),
		),
		RightPlane: sketcher.NewPlaneParametersFromVectors(
			sketcher.NewVectorFromValues(0, 0, 0),
			sketcher.NewVectorFromValues(1, 0, 0),
			sketcher.NewVectorFromValues(0, 1, 0),
		),
	}
}

func (m *MakerCad) Sketch(planer sketcher.Planer) *Sketch {
	sketch := &Sketch{sketcher.NewDlineateSolver(planer)}
	m.sketches = append(m.sketches, sketch)
	return sketch
}

func (*MakerCad) ExportStl(filename string, shapes ListOfShape, quality ExportQuality) error {
	linear := 0.01
	angular := 0.1
	compound := topods.NewCompound()
	builder := brep.NewBuilder()
	builder.MakeCompound(compound)
	for i := range shapes {
		builder.Add(compound, shapes[i].Shape)
	}

	stlWriter := stlapi.NewWriter()
	switch quality {
	case QualityVeryLow:
		linear = 0.5
		angular = 0.8
	case QualityLow:
		linear = 0.1
		angular = 0.5
	case QualityMedium:
		linear = 0.01
		angular = 0.1
	case QualityHigh:
		linear = 0.001
		angular = 0.08
	}

	_ = brepmesh.NewIncrementalMesh(compound, linear, false, angular, true)
	if !stlWriter.Write(compound, filename) {
		return errors.New("Failed to write STL")
	}
	return nil
}

func (*MakerCad) ExportStep(filename string, shapes ListOfShape) error {
	writer := stepcontrol.NewWriter()

	for i := range shapes {
		writer.Transfer(shapes[i].Shape, stepcontrol.ManifoldSolidBrep)
	}

	retStatus := writer.Write(filename)
	if retStatus != StepExportSuccess {
		return errors.New("Failed to write STEP")
	}

	return nil
}

func (*MakerCad) MakeBox(plane sketcher.Planer, dx, dy, dz float64, centerXY bool) Shape {
	p := plane.Plane()
	origin := p.Location()
	if centerXY {
		origin.Translate(gp.NewVec(-dx/2.0, -dy/2.0, 0))
	}
	position := gp.NewAx2(origin, p.Direction(), p.XDirection())
	return Shape{brepprimapi.NewMakeBox(position, dx, dy, dz).Shape()}
}

func (*MakerCad) MakeCylinder(plane sketcher.Planer, radius, height float64) Shape {
	p := plane.Plane()
	position := gp.NewAx2(p.Location(), p.Direction(), p.XDirection())
	return Shape{brepprimapi.NewMakeCylinder(position, radius, height).Shape()}
}

func (*MakerCad) Combine(target Shape, tools ListOfShape) (*CadOperation, error) {
	operation := brepalgoapi.NewFuse().ToBooleanOperation()
	arguments := toptools.NewListOfShape()
	arguments.Append(target.Shape)

	operation.SetTools(tools.ToCascadeList())
	operation.SetArguments(arguments)
	operation.Build()

	return NewCadOperation(tools, &operation), nil
}

func (*MakerCad) Remove(target Shape, tools ListOfShape) (*CadOperation, error) {
	operation := brepalgoapi.NewCut().ToBooleanOperation()
	arguments := toptools.NewListOfShape()
	arguments.Append(target.Shape)

	operation.SetTools(tools.ToCascadeList())
	operation.SetArguments(arguments)
	operation.Build()

	return NewCadOperation(tools, &operation), nil
}

func (*MakerCad) Chamfer(target Shape, edges sketcher.ListOfEdge, depth float64) (Shape, error) {
	fillet := brepfilletapi.NewMakeChamfer(topods.TopoDSShape(target.Shape.Shape))
	for _, e := range edges {
		fillet.AddEdge(topods.TopoDSEdge(e.Edge.Edge), depth)
	}
	return Shape{fillet.Shape()}, nil
}

func (*MakerCad) Fillet(target Shape, edges sketcher.ListOfEdge, radius float64) (Shape, error) {
	fillet := brepfilletapi.NewMakeFillet(topods.TopoDSShape(target.Shape.Shape))
	for _, e := range edges {
		fillet.AddEdge(topods.TopoDSEdge(e.Edge.Edge), radius)
	}
	return Shape{fillet.Shape()}, nil
}
