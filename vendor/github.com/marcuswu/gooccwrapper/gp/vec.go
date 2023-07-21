package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_vec.h>
import "C"

type Vec struct {
	Vec C.gpVec
}

func NewVec(x float64, y float64, z float64) Vec {
	return Vec{C.gpVec_Init(C.double(x), C.double(y), C.double(z))}
}

func NewVecDir(dir Dir) Vec {
	return Vec{C.gpVec_InitDir(dir.dir)}
}

func NewVecPoints(start Pnt, end Pnt) Vec {
	return Vec{C.gpVec_InitPoints(start.Pnt, end.Pnt)}
}

func (v Vec) Free() {
	C.gpVec_Free(v.Vec)
}

func (v Vec) Multiplied(dist float64) Vec {
	return Vec{C.gpVec_Multiplied(v.Vec, C.double(dist))}
}

func (v Vec) Crossed(other Vec) Vec {
	return Vec{C.gpVec_Crossed(v.Vec, other.Vec)}
}

func (v Vec) Normalized() Vec {
	return Vec{C.gpVec_Normalized(v.Vec)}
}

func (v Vec) Magnitude() float64 {
	return float64(C.gpVec_Magnitude(v.Vec))
}
