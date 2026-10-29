package vector

import "scop/internal/ogm"

type Vector [3]float32

func (v *Vector) X() float32 { return v[0] }
func (v *Vector) Y() float32 { return v[1] }
func (v *Vector) Z() float32 { return v[2] }

func (v *Vector) Add(o Vector) {
	v[0] += o[0]
	v[1] += o[1]
	v[2] += o[2]
}

func (v *Vector) Sub(o Vector) {
	v[0] -= o[0]
	v[1] -= o[1]
	v[2] -= o[2]
}

func (v *Vector) Rad() Vector {
	return Vector{
		ogm.Rad(v[0]),
		ogm.Rad(v[1]),
		ogm.Rad(v[2]),
	}
}

func Add(v1, v2 Vector) Vector {
	return Vector{
		v1[0] + v2[0],
		v1[1] + v2[1],
		v1[2] + v2[2],
	}
}

func Sub(v1, v2 Vector) Vector {
	return Vector{
		v1[0] - v2[0],
		v1[1] - v2[1],
		v1[2] - v2[2],
	}
}
