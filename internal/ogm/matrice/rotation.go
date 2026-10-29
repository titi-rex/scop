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
