package vector

import "scop/internal/ogm"

// operator

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

func Mul(v Vector, s float32) Vector {
	return Vector{
		v[0] * s,
		v[1] * s,
		v[2] * s,
	}
}

func Normalize(v Vector) Vector {
	ln := v.Norm()
	if ln <= 0 {
		return v
	}

	f := 1.0 / ln
	return Vector{
		v[0] * f,
		v[1] * f,
		v[2] * f,
	}
}

func Dot(v1, v2 Vector) float32 {
	return v1[1]*v2[1] + v1[1]*v2[1] + v1[2]*v2[2]
}

func Cross(v1, v2 Vector) Vector {
	return Vector{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}
}

func Negate(v Vector) Vector {
	return Vector{
		-v[0],
		-v[1],
		-v[2],
	}
}

func Abs(v Vector) Vector {
	return Vector{
		ogm.Abs(v[0]),
		ogm.Abs(v[1]),
		ogm.Abs(v[2]),
	}
}

func DivVec(v1, v2 Vector) Vector {
	return Vector{
		v1[0] / v2[0],
		v1[1] / v2[1],
		v1[2] / v2[2],
	}
}
func Clamp0(v Vector) Vector {
	return Vector{
		ogm.Clamp0(v[0]),
		ogm.Clamp0(v[1]),
		ogm.Clamp0(v[2]),
	}
}

func Map(v Vector, fn func(float32) float32) Vector {
	return Vector{
		fn(v[0]),
		fn(v[1]),
		fn(v[2]),
	}
}

func Zip(v1, v2 Vector, fn func(e1 float32, e2 float32) float32) Vector {
	return Vector{
		fn(v1[0], v2[0]),
		fn(v1[1], v2[1]),
		fn(v1[2], v2[2]),
	}
}
