package camera

import (
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"
)

type orthoCamera struct {
	core core
}

var _ Camera = (*orthoCamera)(nil)

// NewOrhto create a camera moving only on the world orientation
func NewOrhto(pos,center vector.Vector) orthoCamera {
	return orthoCamera{
		core: core{
			eye:    pos,
			center: center,
			up:     vector.Vector{0, 1, 0},
			speed:  2.5,
		},
	}
}

// SetCenter implements [Camera].
func (c *orthoCamera) SetCenter(center vector.Vector) {
	c.core.center = center
}

// View implements [Camera].
func (c *orthoCamera) View() matrice.Matrice {
	return c.core.view()
}

// Move implements [Camera].
func (c *orthoCamera) Move(d Direction, deltaTime float32) {
	velocity := c.core.speed * deltaTime
	switch d {
	case Forward:
		c.core.zOrthoPlus(velocity)
	case Backward:
		c.core.zOrthoMinus(velocity)
	case Right:
		c.core.xOrthoMinus(velocity)
	case Left:
		c.core.xOrthoPlus(velocity)
	case Up:
		c.core.yOrthoPlus(velocity)
	case Down:
		c.core.yOrthoMinus(velocity)
	}
	// log.Println("camera:", c.core.eye)
}
