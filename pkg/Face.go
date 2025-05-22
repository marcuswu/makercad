package makercad

import (
	"math"

	"github.com/marcuswu/libmakercad/internal/utils"
	"github.com/marcuswu/libmakercad/pkg/sketch"

	"github.com/marcuswu/gooccwrapper/brepadapter"
	"github.com/marcuswu/gooccwrapper/brepalgoapi"
	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/brepgprop"
	"github.com/marcuswu/gooccwrapper/brepprimapi"
	"github.com/marcuswu/gooccwrapper/breptool"
	"github.com/marcuswu/gooccwrapper/breptools"
	"github.com/marcuswu/gooccwrapper/geomabs"
	"github.com/marcuswu/gooccwrapper/geomadapter"
	"github.com/marcuswu/gooccwrapper/geomlprop"
	"github.com/marcuswu/gooccwrapper/gp"
	"github.com/marcuswu/gooccwrapper/gprop"
	"github.com/marcuswu/gooccwrapper/topexp"
	"github.com/marcuswu/gooccwrapper/topods"
	"github.com/marcuswu/gooccwrapper/toptools"
)

type Face struct {
	face topods.Face
}

func NewFace(sketch *Sketch) *Face {
	wires := make([]topods.Wire, 0)
	entities := sketch.solver.Entities()
	for i := range entities {
		entity := entities[i]
		if entity.IsConstruction() {
			continue
		}
		edge := entity.MakeEdge()
		if edge != nil {
			wires = append(wires, brepbuilderapi.NewMakeWireWithEdge(edge.Edge).ToTopoDSWire())
		}
	}

	combined := brepbuilderapi.NewMakeWire()
	for i := range wires {
		combined.AddWire(wires[i])
	}

	return &Face{brepbuilderapi.NewMakeFace(combined.ToTopoDSWire()).ToTopoDSFace()}
}

func (f *Face) getCenter() gp.Pnt {
	shellProps := gprop.NewGProps()
	brepgprop.SurfaceProperties(topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face)), shellProps, false, false)

	if shellProps.Mass() < utils.Confusion {
		return gp.NewPnt(0, 0, 0)
	}

	return shellProps.CenterOfMass()
}

func (f *Face) Plane() gp.Ax3 {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return sketch.NewPlaneParameters().Plane()
	}
	normal := surface.Plane().Axis().Direction()
	location := f.getCenter()
	xDir := surface.Plane().Position().XDirection()

	facePlane := sketch.NewPlaneParametersFromVectors(
		sketch.NewVectorFromValues(location.X(), location.Y(), location.Z()),
		sketch.NewVectorFromValues(normal.X(), normal.Y(), normal.Z()),
		sketch.NewVectorFromValues(xDir.X(), xDir.Y(), xDir.Z()),
	)

	return facePlane.Plane()
}

func (f *Face) Normal() gp.Dir {
	umin, _, vmin, _ := breptools.UVBounds(f.face)
	surface := breptool.Surface(f.face)
	props := geomlprop.NewSLProps(surface, umin, vmin, 1, 0.01)
	return props.Normal()
}

func (f *Face) Revolve(axis *sketch.Line, angle float64) *CadOperation {
	list := toptools.NewListOfShape()
	return f.RevolveMerging(axis, angle, MergeTypeNew, list)
}

func (f *Face) RevolveMerging(axis *sketch.Line, angle float64, merge MergeType, list toptools.ListOfShape) *CadOperation {
	if !f.IsPlanar() {
		return nil
	}

	start := axis.Start.Convert()
	end := axis.End.Convert()
	dir := gp.NewDir(end.X()-start.X(), end.Y()-start.Y(), end.Z()-start.Z())

	ax1 := gp.NewAx1(axis.Start.Convert(), dir)
	shape := brepprimapi.NewMakeRevol(f.face, ax1, angle).Shape()
	if merge == MergeTypeNew || list.Extent() < 1 {
		return &CadOperation{Shape{shape}, nil}
	}

	operation := mergeTypeToOperation(merge)
	tools := toptools.NewListOfShape()
	tools.Append(shape)
	arguments := toptools.NewListOfShape()
	arguments.AppendList(list)

	operation.SetTools(tools)
	operation.SetArguments(arguments)
	operation.Build()

	return &CadOperation{Shape{shape}, operation}
}

func (f *Face) Extrude(distance float64) *CadOperation {
	list := toptools.NewListOfShape()
	return f.ExtrudeMerging(distance, MergeTypeNew, list)
}

func (f *Face) ExtrudeMerging(distance float64, merge MergeType, list toptools.ListOfShape) *CadOperation {
	if !f.IsPlanar() {
		return nil
	}

	coordSystem := f.Normal()

	shape := brepprimapi.NewMakePrism(f.face, gp.NewVecDir(coordSystem).Multiplied(distance)).Shape()
	if merge == MergeTypeNew || list.Extent() < 1 {
		return &CadOperation{Shape{shape}, nil}
	}

	operation := mergeTypeToOperation(merge)
	tools := toptools.NewListOfShape()
	tools.Append(shape)
	arguments := toptools.NewListOfShape()
	arguments.AppendList(list)

	operation.SetTools(tools)
	operation.SetArguments(arguments)
	operation.Build()

	return &CadOperation{Shape{shape}, operation}
}

func mergeTypeToOperation(merge MergeType) *brepalgoapi.Boolean {
	var boolOp brepalgoapi.Boolean
	switch merge {
	case MergeTypeAdd:
		boolOp = brepalgoapi.NewFuse().ToBooleanOperation()
	case MergeTypeRemove:
		boolOp = brepalgoapi.NewCut().ToBooleanOperation()
	default:
		return nil
	}
	return &boolOp
}

func (f *Face) GetFace() topods.Face {
	return f.face
}

func (f *Face) HasEdge(edge topods.Edge) bool {
	explorer := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face)), topexp.Edge)
	for ; explorer.More(); explorer.Next() {
		if explorer.Current().IsEqual(topods.TopoDSShape(edge.Edge)) {
			return true
		}
	}
	return false
}

func (f *Face) IsAlignedWithFace(other *Face) bool {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return false
	}
	otherSurface := brepadapter.NewSurface(other.face)
	if otherSurface.Type() != geomabs.Plane {
		return false
	}
	normal := surface.Plane().Axis().Direction()
	otherNormal := otherSurface.Plane().Axis().Direction()
	return normal.IsParallel(otherNormal)
}

func (f *Face) IsAlignedWithPlane(plane *sketch.PlaneParameters) bool {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return false
	}
	normal := surface.Plane().Axis().Direction()
	planeNormal := gp.NewDirVec(gp.NewVec(plane.Normal.X, plane.Normal.Y, plane.Normal.Z))
	return planeNormal.IsParallel(normal)
}

func (f *Face) IsConical() bool {
	surf := geomadapter.NewSurface(breptool.Surface(f.face))
	defer surf.Free()
	return surf.IsConical()
}

func (f *Face) IsCylindrical() bool {
	surf := geomadapter.NewSurface(breptool.Surface(f.face))
	defer surf.Free()
	return surf.IsCylindrical()
}

func (f *Face) IsPlanar() bool {
	surf := geomadapter.NewSurface(breptool.Surface(f.face))
	defer surf.Free()
	return surf.IsPlanar()
}

func (f *Face) IsNormalAngle(other *Face, angle float64, tolerance float64) bool {
	oShape := topods.NewShapeFromRef(topods.TopoDSShape(other.face.Face))
	shape := topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face))
	return math.Abs(
		oShape.Location().Transformation().Rotation().Multiplied(
			shape.Location().Transformation().Rotation().Inverted(),
		).RotationAngle()-angle,
	) < tolerance
}

func (f *Face) IsOnPlane(plane *sketch.PlaneParameters) bool {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return false
	}
	normal := gp.NewVecDir(surface.Plane().Axis().Direction()).Normalized()
	planeNormal := plane.Normal.ToVector().Normalized()

	sameDirection := planeNormal.IsEqual(normal)

	planeOrigin := plane.Location.ToPoint()
	planePln := gp.NewPlnPntDir(planeOrigin, gp.NewDirVec(planeNormal))
	surfaceOrigin := surface.Plane().Position().Location()
	onPlane := planePln.ContainsPoint(surfaceOrigin)

	return sameDirection && onPlane
}

func (f *Face) IsOpposingNormal(other *Face) bool {
	shape := topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face))
	oShape := topods.NewShapeFromRef(topods.TopoDSShape(other.face.Face))
	return shape.Location().Transformation().Rotation().Inverted().IsEqual(
		oShape.Location().Transformation().Rotation())
}

func (f *Face) IsInDirection(x float64, y float64, z float64) bool {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return false
	}
	normal := surface.Plane().Axis().Direction()
	direction := gp.NewDirVec(gp.NewVec(x, y, z))

	return direction.IsParallel(normal)
}

func (f *Face) DistanceAlong(x float64, y float64, z float64) float64 {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return 0
	}

	location := surface.Plane().Position().Location()
	fromOriginVector := gp.NewVecPoints(gp.NewPnt(0, 0, 0), location)
	if fromOriginVector.Magnitude() < gp.Resolution() {
		return 0.0
	}
	fromOrigin := gp.NewDirVec(fromOriginVector)
	direction := gp.NewDirVec(gp.NewVec(x, y, x))

	return direction.Dot(fromOrigin)
}

func (f *Face) DistanceFrom(x float64, y float64, z float64) float64 {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return 0
	}

	location := surface.Plane().Position().Location()

	return location.Distance(gp.NewPnt(x, y, z))
}

func (f *Face) Edges() []sketch.Edge {
	edges := make([]sketch.Edge, 0)
	explorer := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face)), topexp.Edge)
	for ; explorer.More(); explorer.Next() {
		if explorer.Depth() > 1 {
			continue
		}
		edges = append(edges, *sketch.NewEdgeFromRef(explorer.Current()))
	}

	return edges
}
