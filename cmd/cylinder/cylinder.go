package main

import (
	makercad "github.com/marcuswu/makercad"
)

func main() {
	cad := makercad.NewMakerCad()
	sketch := cad.Sketch(cad.FrontPlane)

	circle := sketch.Circle(0, 0, 5)
	circle.Diameter(10)
	circle.Center.Coincident(sketch.Origin())

	sketch.Solve()
	face := makercad.NewFace(sketch)
	cylinderOp, _ := face.Extrude(10)

	exports := make(makercad.ListOfShape, 0, 1)
	exports = append(exports, cylinderOp.Shape())

	cad.ExportStl("cylinder.stl", exports, makercad.QualityHigh)
}
