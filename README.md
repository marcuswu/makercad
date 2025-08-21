# MakerCAD

## Getting Started

### Install dependencies
 * [OpenCascade](https://dev.opencascade.org/release) CAD Kernel - Installing a binary package is recommended (compiling takes hours).
 * [occwrapper](https://github.com/marcuswu/occwrapper) OpenCascade C wrapper to help make using OpenCascade with Go easier.

### Install MakerCAD
Start a go project

`mkdir myproject && go mod init [module-path]`

Add makercad dependency

`go get github.com/marcuswu/makercad`

Start designing!

## Using MakerCAD ##

Start by creating an instance of MakerCad:

`cad := makercad.NewMakerCad()`

### Creating Primitives ###

#### Plane ####
Planes can be used for creating sketches or placing primitives
```
		MyTopPlane := sketcher.NewPlaneParametersFromVectors(
			sketcher.NewVectorFromValues(0, 0, 0),  // Location (origin for the plane)
			sketcher.NewVectorFromValues(0, 0, 1),  // Normal (Z) direction vector
			sketcher.NewVectorFromValues(1, 0, 0),  // X Direction vector
		)
```

A plane can also be created from a face of a model

`sketcher.NewPlaneParametersFromCoordinateSystem(aFace.Plane())`

There are also some origin planes predefined in MakerCAD:
```
type MakerCad struct {
	sketches    []*Sketch
	FrontPlane  *sketcher.PlaneParameters
	BackPlane   *sketcher.PlaneParameters
	TopPlane    *sketcher.PlaneParameters
	BottomPlane *sketcher.PlaneParameters
	LeftPlane   *sketcher.PlaneParameters
	RightPlane  *sketcher.PlaneParameters
}
```

#### Rectangular Cuboid ####
Provide a plane to locate and orient the shape, its dimensions, and whether to center it on the location
Specific global coordinates can be specified by altering the plane

`block := cad.MakeBox(cad.TopPlane, width, depth, height, true | false)`

#### Cylinder ####
Provide a plane to locate and orient the shape and its radius and height. Cylinders are always centered on the plane origin

`cylinder := cad.MakeCylinder(cad.TopPlane, radius, height)`

### Boolean operations ###
Boolean operations allow combining shapes to create more complex features on your model. These operations return a CadOperation from which the resulting shape can be retrieved.

#### Union ####

`op, err = cad.Combine(targetShape, makercad.ListOfShape{tools...})`

#### Difference ####

`op, err = cad.Remove(targetShape, makercad.ListOfShape{tools...})`

### Sketching ###
Sketching allows for creation of more complex 3D shapes by drawing a 2D shape, optionally adding constraints and solving them, then extruding or revolving them to 3D shapes

#### Creating a Sketch ####
A sketch may be created on a plane or a face of a shape:

`sketch := cad.Sketch(plane | face)`

#### Defining Geometry ####

Lines are defined by start and end points

`l1 := sketch.Line(startX, startY, endX, endY)`

Arcs are defined clockwise around a center point

`arc1 := sketch.Arc(centerX, centerY, startX, startY, endX, endY)`

Circles are defined by a center point and a diameter (does not define a constraint)

`circ1 := sketch.Circle(centerX, centerY, diameter)`

#### Constraining Geometry ####
Sometimes it is not easy to determine the exact geometry when defining a sketch. In these cases, let the computer do the work. Define geometry close to what you need and specify constraints to define how the final geometry should relate.

| Method | Description |
| ------ | ----------- |
| Coincident(Entity, Entity) | Ensure one entity lies on another entity (eg a point on a line) |
| PointVerticalDistance(*Point, Entity, float64) | Ensures a point is a specific distance along the Y axis from the specified entity |
| PointHorizontalDistance(*Point, Entity, float64) | Ensures a point is a specific distance along the X axis from the specified entity |
| PointProjectedDistance(*Point, Entity, float64) | Ensures that a point's projected distance along the normal of Entity is a specific distance |
| LineMidpoint(*Line, Entity) | Ensures entity is coincident with Line and halfway between its start and end points |
| LineAngle(*Line, *Line, float64) | Ensures the angle between two lines is the specified angle (in radians) |
| ArcLineTangent(*Arc, *Line) | Ensures the specified arc and line are tangent to one another |
| Distance(Entity, Entity, float64) | Ensures the two entities are the specified distance from each other |
| HorizontalLine(*Line) | Ensures the specified line is parallel with the X axis | 
| HorizontalPoints(*Point, *Point) | Ensures the imaginary line segment between the two points specified is parallel with the X axis |
| VerticalLine(*Line) | Ensures the specified line is parallel with the X axis | 
| VerticalPoints(*Point, *Point) | Ensures the imaginary line segment between the two points specified is parallel with the X axis |
| LineLength(*Line, float64) | Ensures the specified line has the indicated length |
| Equal(Entity, Entity) | Ensures the two entities are equal (lines the same length, circles the same diameter, etc) |
| CurveDiameter(Entity, float64) | Ensures the arc or circle specified has the indicated diameter |

Many of these constraints also have convenience functions on an entity. For instance to set a line's length:

`line1.Length(10)`

Most of these convenience functions return the entity so they can be chained:

`line1.Length(10).Horizontal()`

#### Solving Constraints ####
Runs the constraint solver algorithm. Returns an error should it be unable to solve.

`err = sketch.Solve()`

#### Debugging Sketches ####
Sketches can be overconstrained or underconstrained. The OverConstrained method returns a list of constraints that overdefine the problem:

`fmt.PrintLn("Over constrained constraints: ", sketch.OverConstrained())`

An underconstrained sketch will have multiple solutions (often infinite).

A sketch can be converted to an SVG image:

`sketch.LogDebug("sketch.svg")`

The constraint solver uses a graph based approach to simplifying the system of equations that are necessary to solve a sketch. If there is an issue with a sketch, logging the graph can be helpful in determining why:

`sketch.ExportImage("sketch.dot")`

These .dot files can be visualized via GraphViz:
```
dot -Tsvg clustered.dot -o clustered.svg
```

#### Extruding or Revolving Sketches ####
First, convert the sketch into a Face:
```
face1 := makercad.NewFace(sketch)
```

Then it can be extruded or revolved
```
operation, err := face1.Extrude(distance)
```

```
operation, err := face1.Revolve(axis, angleInRadians)
```

With either of these, if combining via a boolean operation with an existing shape, it can be done in one step:
```
operation1, err := face1.ExtrudeMerging(distance, MergeType, makercad.ListOfShape{someOp.Shape()})
operation2, err := face2.ExtrudeMerging(distance, MergeType, makercad.ListOfShape{someOp.Shape()})
```

### Finding a Face or Edge ###
A Shape can return its list of Faces:
```
shape1.Faces()
```

A Face or list of Faces can return its list of Edges:
```
face1.Edges()
shape1.Faces().Edges()
```

A list of Faces or Edges can be filtered and sorted:
```
someOperation.Shape().Faces().Edges().
  IsCircle().
  Matching(func(e *sketcher.Edge) bool {
    return e.CircleRadius() == 5.8/2.
  })
```

Connecting the pieces:
```
	block := cad.MakeBox(cad.TopPlane, blockWidth, blockWidth, blockHeight, true)

	// Find the top face aligned with Z positive
	faces := block.Faces().AlignedWith(cad.TopPlane)
	faces.SortByZ(true)
	topFace := faces[0]

  sketch := cad.Sketch(topFace)
  circle := sketch.Circle(0, 0, 5)
  circle.Diameter(5)
  circle.Center.Coincident(sketch.Origin())
  err := sketch.Solve()
  if err != nil {
    // do something
  }

  face1 := makercad.NewFace(sketch)
  newBlock, err = face1.ExtrudeMerging(-2, makercad.MergeTypeRemove, makercad.ListOfShape{block})
```

### Saving Results ###
MakerCAD can export to STL or STEP:
```
exports := makercad.ListOfShape{block}
cad.ExportStl("my-model.stl", exports, makercad.QualityHigh)
cad.ExportStep("my-model.step", exports)
```
