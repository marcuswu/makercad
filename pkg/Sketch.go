package MakerCad

import . "libmakercad/internal/core"

// TODO: Use SketchSolver interface -- provide external facing api
type Sketch struct {
	sketch SketchSolver
	// origin Point
}

// No constructors -- Sketches should be created via MakerCad

func Solve() {

}

// TODO: Origin, SketchSolver getters
// TODO: arc, circle, line, point methods
