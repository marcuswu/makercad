package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_ax2.h>
import "C"

type Ax2 struct {
	Ax2 C.gpAx2
}

func NewAx2(origin Pnt, dir Dir) Ax2 {
	return Ax2{C.gpAx2_Init(origin.Pnt, dir.dir)}
}

func (a Ax2) Free() {
	C.gpAx2_Free(a.Ax2)
}
