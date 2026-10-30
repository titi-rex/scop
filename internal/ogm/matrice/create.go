package matrice

import (
	"scop/internal/ogm"
	"scop/internal/ogm/vector"
)

// NewIdent create a Identity matrice
func NewIdent() *Matrice {
	return &Matrice{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// NewZero create a Zero matrice
func NewZero() *Matrice {
	return &Matrice{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
}

// NewPerspective create a Perspective matrice
func NewPerspective(fovY, aspect, near, far float32) *Matrice {
	nfDiff := near - far
	f := 1 / ogm.Tan(fovY/2)

	return &Matrice{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (near + far) / nfDiff, -1,
		0, 0, (2. * far * near) / nfDiff, 0,
	}
}

// LookAt
// see [mgl32.LookAtV]
func LookAt(eye, center, up vector.Vector) *Matrice {
	f := vector.Sub(center, eye)
	f.Normalize()
	s := vector.Cross(f, up)
	s.Normalize()
	u := vector.Cross(s, f)

	M := &Matrice{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		0, 0, 0, 1,
	}

	return M.Translate(vector.Negate(eye))
}
