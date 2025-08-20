package sketcher

import "github.com/marcuswu/gooccwrapper/gp"

// Planer is any type which can produce a gp.Ax3 plane (PlaneParameters and Face for example)
type Planer interface {
	Plane() gp.Ax3
}
