package camera

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type lockedCamera struct {
	matrice mgl.Mat4
}

var _ Camera = (*lockedCamera)(nil)

func NewLocked() lockedCamera {
	return lockedCamera{
		matrice: mgl.Ident4(),
	}
}

// SetCenter implements [Camera].
func (c *lockedCamera) SetCenter(_ mgl.Vec3) {}

// View implements [Camera].
func (c *lockedCamera) View() mgl.Mat4 { return c.matrice }

// Move implements [Camera].
func (c *lockedCamera) Move(_ Direction, _ float32) {}
