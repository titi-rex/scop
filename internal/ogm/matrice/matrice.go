package matrice

import "scop/internal/ogm/vector"

type Matrice [4 * 4]float32

func (m *Matrice) Add(o *Matrice) *Matrice {
	*m = *Add(m, o)

	return m
}

func (m *Matrice) Mul(o *Matrice) *Matrice {
	*m = *Mul(m, o)

	return m
}

func (m *Matrice) RotateX(angle float32) *Matrice {
	rot := NewRotationX(angle)

	return m.Mul(rot)
}

func (m *Matrice) RotateY(angle float32) *Matrice {
	rot := NewRotationY(angle)

	return m.Mul(rot)
}

func (m *Matrice) RotateZ(angle float32) *Matrice {
	rot := NewRotationZ(angle)

	return m.Mul(rot)
}
func (m *Matrice) RotateXYZ(v vector.Vector) *Matrice {
	return m.RotateX(v.X()).RotateY(v.Y()).RotateZ(v.Z())
}

func (m *Matrice) Translate(v vector.Vector) *Matrice {
	tr := NewTranslation(v.X(), v.Y(), v.Z())

	return m.Mul(tr)
}

// func
// Rotation
// Trans
// Perspective
// LookAt
//
