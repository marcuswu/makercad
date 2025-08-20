package makercad

import (
	"errors"
	"math"
	"slices"

	floatUtils "github.com/marcuswu/dlineate/utils"
	"github.com/marcuswu/makercad/sketcher"
	"github.com/marcuswu/makercad/utils"
	"github.com/rs/zerolog/log"

	"github.com/marcuswu/gooccwrapper/brepadapter"
	"github.com/marcuswu/gooccwrapper/brepalgoapi"
	"github.com/marcuswu/gooccwrapper/brepbuilderapi"
	"github.com/marcuswu/gooccwrapper/brepgprop"
	"github.com/marcuswu/gooccwrapper/brepprimapi"
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

type Orientation int

const (
	Forward = iota
	Reversed
	Internal
	External
)

// Face represents a surface on a 3D [Shape]
type Face struct {
	face topods.Face
}

// ListOfFace is a filterable and sortable slice of [Face]
type ListOfFace []*Face

type FaceFilter func(*Face) bool
type FaceSorter func(a, b *Face) int

// FirstMatching returns the first Face matching the provided filter
func (l ListOfFace) FirstMatching(filter FaceFilter) *Face {
	for _, face := range l {
		if filter(face) {
			return face
		}
	}
	return nil
}

// Matching returns the Faces matching the provided filter
func (l ListOfFace) Matching(filter FaceFilter) ListOfFace {
	newList := make(ListOfFace, 0, len(l))
	for _, face := range l {
		if filter(face) {
			newList = append(newList, face)
		}
	}
	return newList
}

// Sort sorts the list of faces by the sorter
func (l ListOfFace) Sort(sorter FaceSorter) {
	slices.SortFunc(l, sorter)
}

// Planar filters out faces which are not planar
func (l ListOfFace) Planar() ListOfFace {
	return l.Matching(func(f *Face) bool { return f.IsPlanar() })
}

// AlignedWith filters out faces which do not have a parallel normal to the provided plane
func (l ListOfFace) AlignedWith(plane *sketcher.PlaneParameters) ListOfFace {
	return l.Matching(func(f *Face) bool { return f.IsAlignedWithPlane(plane) })
}

// Sort faces by the X position of their center
func (l ListOfFace) SortByX(inverse bool) {
	l.Sort(func(a, b *Face) int {
		aX := a.getCenter().X()
		bX := b.getCenter().X()
		if inverse {
			return floatUtils.StandardFloatCompare(bX, aX)
		}
		return floatUtils.StandardFloatCompare(aX, bX)
	})
}

// Sort faces by the Y position of their center
func (l ListOfFace) SortByY(inverse bool) {
	l.Sort(func(a, b *Face) int {
		aY := a.getCenter().Y()
		bY := b.getCenter().Y()
		if inverse {
			return floatUtils.StandardFloatCompare(bY, aY)
		}
		return floatUtils.StandardFloatCompare(aY, bY)
	})
}

// Sort faces by the Z position of their center
func (l ListOfFace) SortByZ(inverse bool) {
	l.Sort(func(a, b *Face) int {
		aZ := a.getCenter().Z()
		bZ := b.getCenter().Z()
		if inverse {
			return floatUtils.StandardFloatCompare(bZ, aZ)
		}
		return floatUtils.StandardFloatCompare(aZ, bZ)
	})
}

// Return the edges which are contained within this Face
func (l ListOfFace) Edges() sketcher.ListOfEdge {
	le := sketcher.ListOfEdge{}
	for _, e := range l {
		le = append(le, e.Edges()...)
	}
	return le
}

// NewFace creates a face based on the provided sketch. Ignores construction entities and any entities which do not create a Wire (like Points).
// NewFace attempts to create the face by edges ordered by connectivity. If it cannot determine order by connectivity, a non-manifold Face may be returned.
func NewFace(s *Sketch) *Face {
	brepbuilderapi.SetPrecision(0.0001)
	wires := make([]topods.Wire, 0)
	entities := s.solver.Entities()
	wired := make([]sketcher.Entity, 0, len(entities))

	if len(entities) == 0 {
		return nil
	}

	connectsToWire := func(e sketcher.Entity) bool {
		for _, ent := range wired {
			if ent.IsConnectedTo(e) {
				return true
			}
		}
		return false
	}

	addConnectedEntities := func() {
		for i := len(entities) - 1; i >= 0; i-- {
			ent := entities[i]
			// We don't care about construction geometry
			if ent.IsConstruction() {
				entities[i] = entities[len(entities)-1]
				entities = entities[:len(entities)-1]
				continue
			}
			// Come back to unconnected geometry later unless it's our first entity
			if len(wired) > 0 && !connectsToWire(ent) {
				continue
			}
			// Add the connected geometry to our wire list and remove it from available entities
			wired = append(wired, ent)
			entities[i] = entities[len(entities)-1]
			entities = entities[:len(entities)-1]
		}
	}

	// Pass through twice; if we still have entities, try adding directly
	// TODO: Think about this and test some more. There's probably a better way to check completion
	addConnectedEntities()
	addConnectedEntities()

	for _, entity := range entities {
		if entity.IsConstruction() {
			continue
		}
		wired = append(wired, entity)
	}

	for _, entity := range wired {
		if entity.IsConstruction() {
			continue
		}
		edge := entity.MakeEdge()
		if edge != nil {
			wires = append(wires, brepbuilderapi.NewMakeWireWithEdge(edge.Edge).ToTopoDSWire())
		}
	}

	log.Debug().Int("wire count", len(wires)).Msg("Making combined wire for face")
	combined := brepbuilderapi.NewMakeWire()
	for i := range wires {
		combined.AddWire(wires[i])
	}
	wire := combined.ToTopoDSWire()
	makeFace := brepbuilderapi.NewMakeFace(wire)
	topoFace := makeFace.ToTopoDSFace()

	return &Face{topoFace}
}

func (f *Face) getCenter() gp.Pnt {
	shellProps := gprop.NewGProps()
	brepgprop.SurfaceProperties(topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face)), shellProps, false, false)

	if shellProps.Mass() < utils.Confusion {
		return gp.NewPnt(0, 0, 0)
	}

	return shellProps.CenterOfMass()
}

// Plane returns the plane this Face is on
func (f *Face) Plane() gp.Ax3 {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return sketcher.NewPlaneParameters().Plane()
	}
	normal := f.Normal()
	location := f.getCenter()
	xDir := surface.Plane().Position().XDirection()

	facePlane := sketcher.NewPlaneParametersFromVectors(
		sketcher.NewVectorFromValues(location.X(), location.Y(), location.Z()),
		sketcher.NewVectorFromValues(normal.X(), normal.Y(), normal.Z()),
		sketcher.NewVectorFromValues(xDir.X(), xDir.Y(), xDir.Z()),
	)

	return facePlane.Plane()
}

// Normal returns the normal direction of the Face
func (f *Face) Normal() gp.Dir {
	umin, _, vmin, _ := breptools.UVBounds(f.face)
	surface := f.face.Surface()
	props := geomlprop.NewSLProps(surface, umin, vmin, 1, 0.01)
	normal := props.Normal()
	if f.face.Orientation() == Reversed {
		normal = gp.NewDir(-normal.X(), -normal.Y(), -normal.Z())
	}
	return normal
}

// Revolve creates a 3D shape by revolving this face around the provided axis by the specified angle in radians. Creates a new Shape.
func (f *Face) Revolve(axis *sketcher.Line, angle float64) (*CadOperation, error) {
	list := toptools.NewListOfShape()
	return f.RevolveMerging(axis, angle, MergeTypeNew, list)
}

// RevolveMerging creates a 3D shape by revolving this face around the provided axis by the specified angle
// in radians and performs the specified boolean operation with the shapes in the provided list
func (f *Face) RevolveMerging(axis *sketcher.Line, angle float64, merge MergeType, list toptools.ListOfShape) (*CadOperation, error) {
	if !f.IsPlanar() {
		return nil, errors.New("cannot revolve non-planar face")
	}

	start := axis.Start.Convert()
	end := axis.End.Convert()
	if start.Distance(end) == 0 {
		return nil, errors.New("revolve axis must have non-zero length")
	}
	dir := gp.NewDir(end.X()-start.X(), end.Y()-start.Y(), end.Z()-start.Z())

	ax1 := gp.NewAx1(axis.Start.Convert(), dir)
	shape := brepprimapi.NewMakeRevol(f.face, ax1, angle).Shape()
	if merge == MergeTypeNew || list.Extent() < 1 {
		return &CadOperation{[]Shape{{shape}}, nil}, nil
	}

	operation := mergeTypeToOperation(merge)
	tools := toptools.NewListOfShape()
	tools.Append(shape)
	arguments := toptools.NewListOfShape()
	arguments.AppendList(list)

	operation.SetTools(tools)
	operation.SetArguments(arguments)
	operation.Build()

	return &CadOperation{[]Shape{{shape}}, operation}, nil
}

// Extrude creates a prism using this Face along its normal by distance
func (f *Face) Extrude(distance float64) *CadOperation {
	list := toptools.NewListOfShape()
	return f.ExtrudeMerging(distance, MergeTypeNew, list)
}

// ExtrudeMerging creates a prism using this Face along its normal by distance using the specified boolean operation
// to merge the result with the list of provided shapes
func (f *Face) ExtrudeMerging(distance float64, merge MergeType, list toptools.ListOfShape) *CadOperation {
	if !f.IsPlanar() {
		return nil
	}

	coordSystem := f.Normal()

	shape := brepprimapi.NewMakePrism(f.face, gp.NewVecDir(coordSystem).Multiplied(distance)).Shape()
	if merge == MergeTypeNew || list.Extent() < 1 {
		return &CadOperation{[]Shape{{shape}}, nil}
	}

	operation := mergeTypeToOperation(merge)
	tools := toptools.NewListOfShape()
	tools.Append(shape)
	arguments := toptools.NewListOfShape()
	arguments.AppendList(list)

	operation.SetTools(tools)
	operation.SetArguments(arguments)
	operation.Build()

	return &CadOperation{[]Shape{{shape}}, operation}
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

// GetFace returns the OpenCascade face reference
func (f *Face) GetFace() topods.Face {
	return f.face
}

// AsShape returns the 2D face as a Shape object
func (f *Face) AsShape() *Shape {
	return &Shape{topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face))}
}

// Mirror mirrors this face across the specified plane. This may flip the normal
func (f *Face) Mirror(plane *sketcher.PlaneParameters) (*Face, error) {
	coord := plane.Ax2()
	gptrans := gp.NewTrsf()
	gptrans.SetMirrorAx2(coord)
	trans := brepbuilderapi.NewTransform(f.AsShape().Shape, gptrans)
	return &Face{topods.NewFaceFromRef(topods.TopoDSFace(trans.Shape().Shape))}, nil
}

// HasEdge returns whether this Face contains the specified edge
func (f *Face) HasEdge(edge topods.Edge) bool {
	explorer := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face)), topexp.Edge)
	for ; explorer.More(); explorer.Next() {
		if explorer.Current().IsEqual(topods.TopoDSShape(edge.Edge)) {
			return true
		}
	}
	return false
}

// IsAlignedWithFace returns wither this Face and the provided Face have equivalent normals
func (f *Face) IsAlignedWithFace(other *Face) bool {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return false
	}
	otherSurface := brepadapter.NewSurface(other.face)
	if otherSurface.Type() != geomabs.Plane {
		return false
	}
	normal := f.Normal()
	otherNormal := other.Normal()

	return normal.IsEqual(otherNormal)
}

// IsAlignedWithPlane returns wither this Face and the provided plane have equivalent normals
func (f *Face) IsAlignedWithPlane(plane *sketcher.PlaneParameters) bool {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return false
	}
	normal := f.Normal()
	planeNormal := gp.NewDirVec(gp.NewVec(plane.Normal.X, plane.Normal.Y, plane.Normal.Z))
	return planeNormal.IsEqual(normal)
}

// IsConical returns whether this face is a conical face
func (f *Face) IsConical() bool {
	surf := geomadapter.NewSurface(f.face.Surface())
	defer surf.Free()
	return surf.IsConical()
}

// IsCylindrical returns whether this face is a cylindrical face
func (f *Face) IsCylindrical() bool {
	surf := geomadapter.NewSurface(f.face.Surface())
	defer surf.Free()
	return surf.IsCylindrical()
}

// IsPlanar returns whether this face is a planar face
func (f *Face) IsPlanar() bool {
	surf := geomadapter.NewSurface(f.face.Surface())
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

// IsOnPlane returns whether this face exists on the provided plane
func (f *Face) IsOnPlane(plane *sketcher.PlaneParameters) bool {
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

// IsOpposingNormal returns whether this face and the provided face have opposing normals
func (f *Face) IsOpposingNormal(other *Face) bool {
	shape := topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face))
	oShape := topods.NewShapeFromRef(topods.TopoDSShape(other.face.Face))
	return shape.Location().Transformation().Rotation().Inverted().IsEqual(
		oShape.Location().Transformation().Rotation())
}

// IsInDirection returns whether this face's normal and the provided vector are aligned
func (f *Face) IsInDirection(x float64, y float64, z float64) bool {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return false
	}
	normal := surface.Plane().Axis().Direction()
	direction := gp.NewDirVec(gp.NewVec(x, y, z))

	return direction.IsParallel(normal)
}

// DistanceAlong returns the distance from global origin to this face's origin along the provided vector
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
	direction := gp.NewDirVec(gp.NewVec(x, y, z))

	return direction.Dot(fromOrigin)
}

// DistanceFrom returns the distance of the face origin from a specified location
func (f *Face) DistanceFrom(x float64, y float64, z float64) float64 {
	surface := brepadapter.NewSurface(f.face)
	if surface.Type() != geomabs.Plane {
		return 0
	}

	location := surface.Plane().Position().Location()

	return location.Distance(gp.NewPnt(x, y, z))
}

// Edges returns the list of edges that make up the face
func (f *Face) Edges() sketcher.ListOfEdge {
	edges := make(sketcher.ListOfEdge, 0)
	explorer := topexp.NewExplorer(topods.NewShapeFromRef(topods.TopoDSShape(f.face.Face)), topexp.Edge)
	for ; explorer.More(); explorer.Next() {
		if explorer.Depth() > 1 {
			continue
		}
		edges = append(edges, sketcher.NewEdgeFromRef(explorer.Current()))
	}

	return edges
}
