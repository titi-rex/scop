package camera

import (
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"
)

type lockedCamera struct {
	matrice matrice.Matrice
}

var _ Camera = (*lockedCamera)(nil)

func NewLocked() lockedCamera {
	return lockedCamera{
		matrice: *matrice.NewIdent(),
	}
}

// SetCenter implements [Camera].
func (c *lockedCamera) SetCenter(_ vector.Vector) {}

// View implements [Camera].
func (c *lockedCamera) View() *matrice.Matrice { return &c.matrice }

// Move implements [Camera].
func (c *lockedCamera) Move(_ Direction, _ float32) {}
