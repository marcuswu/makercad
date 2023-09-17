package makercad

import (
	"errors"
	"libmakercad/pkg/sketch"

	"github.com/marcuswu/gooccwrapper/brep"
	"github.com/marcuswu/gooccwrapper/brepmesh"
	"github.com/marcuswu/gooccwrapper/stepcontrol"
	"github.com/marcuswu/gooccwrapper/stlapi"
	"github.com/marcuswu/gooccwrapper/topods"
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
	FrontPlane  *sketch.PlaneParameters
	BackPlane   *sketch.PlaneParameters
	TopPlane    *sketch.PlaneParameters
	BottomPlane *sketch.PlaneParameters
	LeftPlane   *sketch.PlaneParameters
	RightPlane  *sketch.PlaneParameters
}

func NewMakerCad() *MakerCad {
	return &MakerCad{
		sketches: make([]*Sketch, 0),
		FrontPlane: sketch.NewPlaneParametersFromVectors(
			sketch.NewVectorFromValues(0, 0, 0),
			sketch.NewVectorFromValues(0, -1, 0),
			sketch.NewVectorFromValues(1, 0, 0),
		),
		BackPlane: sketch.NewPlaneParametersFromVectors(
			sketch.NewVectorFromValues(0, 0, 0),
			sketch.NewVectorFromValues(0, 1, 0),
			sketch.NewVectorFromValues(-1, 0, 0),
		),
		TopPlane: sketch.NewPlaneParametersFromVectors(
			sketch.NewVectorFromValues(0, 0, 0),
			sketch.NewVectorFromValues(0, 0, 1),
			sketch.NewVectorFromValues(1, 0, 0),
		),
		BottomPlane: sketch.NewPlaneParametersFromVectors(
			sketch.NewVectorFromValues(0, 0, 0),
			sketch.NewVectorFromValues(0, 0, -1),
			sketch.NewVectorFromValues(1, 0, 0),
		),
		LeftPlane: sketch.NewPlaneParametersFromVectors(
			sketch.NewVectorFromValues(0, 0, 0),
			sketch.NewVectorFromValues(-1, 0, 0),
			sketch.NewVectorFromValues(0, -1, 0),
		),
		RightPlane: sketch.NewPlaneParametersFromVectors(
			sketch.NewVectorFromValues(0, 0, 0),
			sketch.NewVectorFromValues(1, 0, 0),
			sketch.NewVectorFromValues(0, 1, 0),
		),
	}
}

func (m *MakerCad) Sketch(planer sketch.Planer) *Sketch {
	sketch := &Sketch{sketch.NewDlineateSolver(planer)}
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
