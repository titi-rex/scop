package camera

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type yLockedCamera struct {
	core core
}

var _ Camera = (*yLockedCamera)(nil)

func NewYLocked(pos, center mgl.Vec3) yLockedCamera {
	return yLockedCamera{
		core: core{
			eye:    pos,
			center: center,
			up:     mgl.Vec3{0, 1, 0},
			speed:  2.5,
		},
	}
}

// SetCenter implements [Camera].
func (c *yLockedCamera) SetCenter(center mgl.Vec3) {
	c.core.center = center
}

// View implements [Camera].
func (c *yLockedCamera) View() mgl.Mat4 {
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
