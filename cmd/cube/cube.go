package main

import makercad "github.com/marcuswu/makercad"

func main() {

	cad := makercad.NewMakerCad()
	sketch := cad.Sketch(cad.TopPlane)

	line1 := sketch.Line(0.0, 0.0, 5.0, 0)
	line1.Length(10).Horizontal()

	line2 := sketch.Line(5.0, 0.0, 5.0, 6.0)
	line2.Length(10).Vertical()

	line3 := sketch.Line(5.0, 5.0, 0.0, 6.0)
	line3.Horizontal()

	line4 := sketch.Line(0.0, 5.0, 0.0, 1.0)
	line4.Vertical()

	line1.End.Coincident(line2.Start)
	line2.End.Coincident(line3.Start)
	line3.End.Coincident(line4.Start)
	line4.End.Coincident(line1.Start)

	// TODO: fix this to work right (dlineate issue)
	// line1Midpoint := sketch.Point(0.0, 0.0).Vertical(sketch.Origin())
	// line2Midpoint := sketch.Point(0.0, 0.0).Horizontal(sketch.Origin())

	// line1.Midpoint(line1Midpoint)
	// line2.Midpoint(line2Midpoint)
	line1.Start.Coincident(sketch.Origin())

	sketch.Solve()
	face := makercad.NewFace(sketch)
	cubeOp := face.Extrude(10)

	exports := make(makercad.ListOfShape, 0, 1)
	exports = append(exports, cubeOp.Shape())

	cad.ExportStl("cube.stl", exports, makercad.QualityHigh)
}
