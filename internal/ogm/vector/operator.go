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

func Normalize(v1 Vector) Vector {
	l := 1.0 / v1.Len()
	return Vector{
		v1[0] * l,
		v1[1] * l,
		v1[2] * l,
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
