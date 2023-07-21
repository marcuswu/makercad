package kernel

import (
	"github.com/marcuswu/gooccwrapper/gp"
	"github.com/marcuswu/libmakercad/internal/solver"
)

// Takes geometry defined in core, converts to edges, then to wires, then to a face

func ArctoEdge(a solver.Arc) Edge {
	// convert center point
	// convert point 1
	// convert point 2
	_ = gp.NewAx2(gp.NewPnt(0, 0, 0), gp.NewDir(0, 0, 0))
	return Edge{}
}
