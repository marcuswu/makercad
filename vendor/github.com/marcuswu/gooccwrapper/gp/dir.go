package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_dir.h>
import "C"

type GPDir C.gpDir

type Dir struct {
	dir C.gpDir
}

func NewDir(x float64, y float64, z float64) Dir {
	return Dir{C.gpDir_Init(C.double(x), C.double(y), C.double(z))}
}

func NewDirVec(vector Vec) Dir {
	return Dir{C.gpDir_InitVec(vector.Vec)}
}

func NewDirFromRef(ref GPDir) Dir {
	return Dir{C.gpDir(ref)}
}

func (d Dir) Free() {
	C.gpDir_Free(d.dir)
}
