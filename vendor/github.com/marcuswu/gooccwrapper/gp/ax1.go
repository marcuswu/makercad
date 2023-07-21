package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_ax1.h>
import "C"

type Ax1 struct {
	Ax1 C.gpAx1
}

func NewAx1(origin Pnt, dir Dir) Ax1 {
	return Ax1{C.gpAx1_Init(origin.Pnt, dir.dir)}
}

func (a Ax1) Direction() Dir {
	return Dir{C.gpAx1_Direction(a.Ax1)}
}

func (a Ax1) Free() {
	C.gpAx1_Free(a.Ax1)
}
