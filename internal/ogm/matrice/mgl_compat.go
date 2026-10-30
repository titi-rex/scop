package matrice

import "github.com/go-gl/mathgl/mgl32"

func MglToMat(m mgl32.Mat4) Matrice {
	mat := Matrice{}

	copy(mat[:], m[:])

	return mat
}

func MatToMgl(m Matrice) mgl32.Mat4 {
	mat := mgl32.Mat4{}

	copy(mat[:], m[:])

	return mat
}
