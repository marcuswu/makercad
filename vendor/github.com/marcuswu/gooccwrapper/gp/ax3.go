package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_ax3.h>
import "C"

type Ax3 struct {
	ax3 C.gpAx3
}

func NewAx3(origin Pnt, normal Dir, xDir Dir) Ax3 {
	return Ax3{C.gpAx3_Init(origin.Pnt, normal.dir, xDir.dir)}
}

func (a Ax3) Free() {
	C.gpAx3_Free(a.ax3)
}

func (a Ax3) XDirection() Dir {
	return Dir{C.gpAx3_XDirection(a.ax3)}
}

func (a Ax3) YDirection() Dir {
	return Dir{C.gpAx3_YDirection(a.ax3)}
}

func (a Ax3) Direction() Dir {
	return Dir{C.gpAx3_Direction(a.ax3)}
}
