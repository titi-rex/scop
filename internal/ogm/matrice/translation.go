package matrice

import "scop/internal/ogm/vector"

// NewTranslation create a translation matrice that move  by x, y, z.
//
//	[[1, 0, 0, Tx]]
//	[[0, 1, 0, Ty]]
//	[[0, 0, 1, Tz]]
//	[[0, 0, 0, 1 ]]
func NewTranslation(x, y, z float32) *Matrice {
	return &Matrice{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	}
}

// Translate translate the matrice by vec.
func (m *Matrice) Translate(vec vector.Vector) *Matrice {
	tr := NewTranslation(vec.X(), vec.Y(), vec.Z())

	return m.Mul(tr)
}
