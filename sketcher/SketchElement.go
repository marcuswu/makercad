package sketcher

// A SketchElement is any of Arc, Circle, Line, or Point
type SketchElement interface {
	Arc | Circle | Line | Point
}
