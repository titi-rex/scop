package matrice

// Add perform matrice addition
func (m *Matrice) Add(o Matrice) *Matrice {
	*m = Add(*m, o)

	return m
}

// Mul perform matrice substraction
func (m *Matrice) Mul(o Matrice) *Matrice {
	*m = Mul(*m, o)

	return m
}

// Scale perform matrice substraction
func (m *Matrice) Scale(s float32) *Matrice {
	*m = Scale(*m, s)

	return m
}

// Inverse inverse matrice
func (m *Matrice) Inverse() *Matrice {
	*m = Inverse(*m)

	return m
}
