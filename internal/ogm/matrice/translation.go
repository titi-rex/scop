package matrice

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
