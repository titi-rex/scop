package camera

import (
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"

)

type basicCamera struct {
	core core
}

var _ Camera = (*basicCamera)(nil)

// NewBasic create a nice Camera to start with
func NewBasic(pos, center vector.Vector) basicCamera {
	return basicCamera{
		core: core{
			eye:    pos,
			center: center,
			up:     vector.Vector{0, 1, 0},
			speed:  2.5,
		},
	}
}

// SetCenter implements [Camera].
func (c *basicCamera) SetCenter(center vector.Vector) {
	c.core.center = center
}

// View implements [Camera].
func (c *basicCamera) View() *matrice.Matrice {
	return c.core.view()
}

// Move implements [Camera].
func (c *basicCamera) Move(d Direction, deltaTime float32) {
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
		c.core.yOrthoPlus(velocity)
	case Down:
		c.core.yOrthoMinus(velocity)
	}
}
