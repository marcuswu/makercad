package gp

// #cgo LDFLAGS: -loccwrapper
// #include <occwrapper/gp_pnt.h>
import "C"

type GPPnt C.gpPnt

type Pnt struct {
	Pnt C.gpPnt
}

func NewPnt(x float64, y float64, z float64) Pnt {
	var ret Pnt
	ret.Pnt = C.gpPnt_Init(C.double(x), C.double(y), C.double(z))
	return ret
}

func NewPntFromRef(ref GPPnt) Pnt {
	return Pnt{C.gpPnt(ref)}
}

func (p Pnt) Free() {
	C.gpPnt_Free(p.Pnt)
}

func (p Pnt) SetCoord(x float64, y float64, z float64) {
	C.gpPnt_SetCoord(p.Pnt, C.double(x), C.double(y), C.double(z))
}

func (p Pnt) SetX(x float64) {
	C.gpPnt_SetX(p.Pnt, C.double(x))
}

func (p Pnt) SetY(y float64) {
	C.gpPnt_SetY(p.Pnt, C.double(y))
}

func (p Pnt) SetZ(z float64) {
	C.gpPnt_SetZ(p.Pnt, C.double(z))
}

func (p Pnt) X() float64 {
	return float64(C.gpPnt_X(p.Pnt))
}

func (p Pnt) Y() float64 {
	return float64(C.gpPnt_Y(p.Pnt))
}

func (p Pnt) Z() float64 {
	return float64(C.gpPnt_Z(p.Pnt))
}

func (p Pnt) IsEqual(o Pnt, tol float64) bool {
	return bool(C.gpPnt_IsEqual(p.Pnt, o.Pnt, C.double(tol)))
}

func (p Pnt) Distance(o Pnt) float64 {
	return float64(C.gpPnt_Distance(p.Pnt, o.Pnt))
}

func (p Pnt) SquareDistance(o Pnt) float64 {
	return float64(C.gpPnt_SquareDistance(p.Pnt, o.Pnt))
}

func (p Pnt) MirrorCenterPoint(center Pnt) {
	C.gpPnt_MirrorCenterPoint(p.Pnt, center.Pnt)
}

func (p Pnt) MirroredCenterPoint(center Pnt) Pnt {
	return Pnt{C.gpPnt_MirroredCenterPoint(p.Pnt, center.Pnt)}
}

func (p Pnt) MirrorAxis(axis Ax1) {
	C.gpPnt_MirrorAxis(p.Pnt, axis.Ax1)
}

func (p Pnt) MirroredAxis(axis Ax1) Pnt {
	return Pnt{C.gpPnt_MirroredAxis(p.Pnt, axis.Ax1)}
}

func (p Pnt) Rotate(axis Ax1, angle float64) {
	C.gpPnt_Rotate(p.Pnt, axis.Ax1, C.double(angle))
}

func (p Pnt) Rotated(axis Ax1, angle float64) Pnt {
	return Pnt{C.gpPnt_Rotated(p.Pnt, axis.Ax1, C.double(angle))}
}

func (p Pnt) Scale(other Pnt, scale float64) {
	C.gpPnt_Scale(p.Pnt, other.Pnt, C.double(scale))
}

func (p Pnt) Scaled(other Pnt, scale float64) Pnt {
	return Pnt{C.gpPnt_Scaled(p.Pnt, other.Pnt, C.double(scale))}
}

func (p Pnt) Transform(trans Trsf) {
	C.gpPnt_Transform(p.Pnt, trans.trsf)
}

func (p Pnt) Transformed(trans Trsf) Pnt {
	return Pnt{C.gpPnt_Transformed(p.Pnt, trans.trsf)}
}

func (p Pnt) Translate(vector Vec) {
	C.gpPnt_Translate(p.Pnt, vector.Vec)
}

func (p Pnt) Translated(vector Vec) Pnt {
	return Pnt{C.gpPnt_Translated(p.Pnt, vector.Vec)}
}

func (p Pnt) TranslatePoints(pnt1 Pnt, pnt2 Pnt) {
	C.gpPnt_TranslatePoints(p.Pnt, pnt1.Pnt, pnt2.Pnt)
}

func (p Pnt) TranslatedPoints(pnt1 Pnt, pnt2 Pnt) Pnt {
	return Pnt{C.gpPnt_TranslatedPoints(p.Pnt, pnt1.Pnt, pnt2.Pnt)}
}
