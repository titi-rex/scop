package vector

import "scop/internal/ogm"

// ToRad convert the vector value to Radian
func ToRad(v Vector) Vector {
	return Vector{
		ogm.Rad(v[0]),
		ogm.Rad(v[1]),
		ogm.Rad(v[2]),
	}
}
