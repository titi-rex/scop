package camera

import (
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"
)

type yLockedCamera struct {
	core core
}

var _ Camera = (*yLockedCamera)(nil)

func NewYLocked(pos, center vector.Vector) yLockedCamera {
	return yLockedCamera{
		core: core{
			eye:    pos,
			center: center,
			up:     vector.Vector{0, 1, 0},
			speed:  2.5,
		},
	}
}

// SetCenter implements [Camera].
func (c *yLockedCamera) SetCenter(center vector.Vector) {
	c.core.center = center
}

// View implements [Camera].
func (c *yLockedCamera) View() *matrice.Matrice {
	return c.core.view()
}

// Move implements [Camera].
func (c *yLockedCamera) Move(d Direction, deltaTime float32) {
	velocity := c.core.speed * deltaTime
	switch d {
	case Forward:
		c.core.zMinus(velocity)
	case Backward:
		c.core.zPlus(velocity)
	case Right:
		c.core.xPlus(velocity)
	case Left:
		c.core.xMinus(velocity)
	case Up:
	case Down:
	}
	c.core.eye[1] = 0
}
