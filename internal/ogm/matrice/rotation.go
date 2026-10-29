package matrice

import (
	"scop/internal/ogm"
	"scop/internal/ogm/vector"
)

// NewRotationX create a rotation matrice that rotate by angle around the X-axis.
func NewRotationX(angle float32) *Matrice {
	sin, cos := ogm.Sincos(angle)

	return &Matrice{
		1, 0, 0, 0,
		0, cos, sin, 0,
		0, -sin, cos, 0,
		0, 0, 0, 1,
	}
}

// NewRotationY create a rotation matrice that rotate by angle around the Y-axis.
func NewRotationY(angle float32) *Matrice {
	sin, cos := ogm.Sincos(angle)

	return &Matrice{
		cos, 0, -sin, 0,
		0, 1, 0, 0,
		sin, 0, cos, 0,
		0, 0, 0, 1,
	}
}

// NewRotationZ create a rotation matrice that rotate by angle around the Z-axis.
func NewRotationZ(angle float32) *Matrice {
	sin, cos := ogm.Sincos(angle)

	return &Matrice{
		cos, sin, 0, 0,
		-sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// NewRotation create a rotation matrix that rotate around an arbitrary axis.
//
//	[[ x^2(1-c)+c, xy(1-c)-zs, xz(1-c)+ys, 0 ]]
//	[[ xy(1-c)+zs, y^2(1-c)+c, yz(1-c)-xs, 0 ]]
//	[[ xz(1-c)-ys, yz(1-c)+xs, z^2(1-c)+c, 0 ]]
//	[[ 0         , 0         , 0         , 1 ]]
func NewRotation(angle float32, axis vector.Vector) *Matrice {
	x, y, z := axis[0], axis[1], axis[2]
	sin, cos := ogm.Sincos(angle)
	k := 1 - cos

	return &Matrice{
		x*x*k + cos, x*y*k + z*sin, x*z*k - y*sin, 0,
		x*y*k - z*sin, y*y*k + cos, y*z*k + x*sin, 0,
		x*z*k + y*sin, y*z*k - x*sin, z*z*k + cos, 0,
		0, 0, 0, 1,
	}
}

// RotateX rotate the matrice by angle around the X-axis.
func (m *Matrice) RotateX(angle float32) *Matrice {
	rot := NewRotationX(angle)

	return m.Mul(rot)
}

// RotateY rotate the matrice by angle around the Y-axis.
func (m *Matrice) RotateY(angle float32) *Matrice {
	rot := NewRotationY(angle)

	return m.Mul(rot)
}

// RotateZ rotate the matrice by angle around the Z-axis.
func (m *Matrice) RotateZ(angle float32) *Matrice {
	rot := NewRotationZ(angle)

	return m.Mul(rot)
}

// RotateXYZ rotate the matrice by each angles around XYZ.
// (apply in order: Z -> Y -> X)
func (m *Matrice) RotateXYZ(angles vector.Vector) *Matrice {
	return m.
		RotateZ(angles.Z()).
		RotateY(angles.Y()).
		RotateX(angles.X())
}

// RotateY rotate the matrice by angle around axis.
func (m *Matrice) Rotate(angle float32, axis vector.Vector) *Matrice {
	rot := NewRotation(angle, axis)

	return m.Mul(rot)
}
