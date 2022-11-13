package dlineate

type Workplane struct {
	origin *Vector3D
	xDir   *Vector3D
	yDir   *Vector3D
}

func NewWorkPlane(origin *Vector3D, xDir *Vector3D, yDir *Vector3D) *Workplane {
	wp := new(Workplane)
	wp.origin = origin
	wp.xDir = xDir
	wp.yDir = yDir
	return wp
}
