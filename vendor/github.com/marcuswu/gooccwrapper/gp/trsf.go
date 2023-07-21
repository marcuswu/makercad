package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_trsf.h>
import "C"

type Trsf struct {
	trsf C.gpTrsf
}

func NewTrsf() Trsf {
	return Trsf{C.gpTrsf_Init()}
}

func (t Trsf) SetTransformation(from Ax3, to Ax3) {
	C.gpTrsf_SetTransformation(t.trsf, from.ax3, to.ax3)
}

func (t Trsf) Free() {
	C.gpTrsf_Free(t.trsf)
}
