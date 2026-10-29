package vector

import "scop/internal/ogm"

type Vector [3]float32

func (v *Vector) X() float32 { return v[0] }
func (v *Vector) Y() float32 { return v[1] }
func (v *Vector) Z() float32 { return v[2] }

func (v1 Vector) Len() float32 {
	return ogm.Sqrt(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])
}
