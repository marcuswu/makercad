package sketcher

import "github.com/marcuswu/gooccwrapper/gp"

type Planer interface {
	Plane() gp.Ax3
}
