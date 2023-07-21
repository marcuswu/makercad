package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_pln.h>
import "C"

type GPPln C.gpPln

type Pln struct {
	pln C.gpPln
}

func NewPln() Pln {
	return Pln{C.gpPln_Init()}
}

func NewPlnAx3(coord Ax3) Pln {
	return Pln{C.gpPln_InitAx3(coord.ax3)}
}

func NewPlnPntDir(origin Pnt, dir Dir) Pln {
	return Pln{C.gpPln_InitPntDir(origin.Pnt, dir.dir)}
}

func NewPlnFromRef(ref GPPln) Pln {
	return Pln{C.gpPln(ref)}
}

func (p Pln) Axis() Ax1 {
	return Ax1{C.gpPln_Axis(p.pln)}
}

func (p Pln) Position() Ax3 {
	return Ax3{C.gpPln_Position(p.pln)}
}

func (p Pln) ContainsPoint(point Pnt, tolerance float64) bool {
	return bool(C.gpPln_ContainsPoint(p.pln, point.Pnt, C.double(tolerance)))
}

func (p Pln) Free() {
	C.gpPln_Free(p.pln)
}
