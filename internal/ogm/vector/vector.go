package vector

import (
	"scop/internal/ogm"

	"github.com/go-gl/mathgl/mgl32"
)

type Vector [3]float32

func (v *Vector) X() float32 { return v[0] }
func (v *Vector) Y() float32 { return v[1] }
func (v *Vector) Z() float32 { return v[2] }

func (v1 Vector) Len() float32 {
	return ogm.Sqrt(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])
}

func MglToVec(v mgl32.Vec3) Vector {
	return Vector{
		v[0],
		v[1],
		v[2],
	}
}

func VecToMgl(v Vector) mgl32.Vec3 {
	return mgl32.Vec3{
		v[0],
		v[1],
		v[2],
	}
}
