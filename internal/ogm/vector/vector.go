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

func (v *Vector) Data() [3]float32 {
	return [3]float32(*v)
}

// WARNING: this function can panic
func LoadIntoSlice(s []float32, v Vector, idx uint32) {
	s[idx] = v[0]
	s[idx+1] = v[1]
	s[idx+2] = v[2]
}

// WARNING: this function can panic
func LoadFromSlice(s []float32, idx uint32) Vector {
	return Vector{
		s[idx],
		s[idx+1],
		s[idx+2],
	}
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
