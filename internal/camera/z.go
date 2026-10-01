package camera

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type cameraZ struct {
	eye    mgl.Vec3
	center mgl.Vec3
	up     mgl.Vec3
	view   mgl.Mat4

	speed float32
}

// NewCameraZ create a Camera with {0,1,0} as up
func NewCameraZ(pos, center mgl.Vec3) cameraZ {
	return cameraZ{
		eye:    pos,
		center: center,
		up:     mgl.Vec3{0, 1, 0},
		speed:  2.5,
	}
}

func (c *cameraZ) SetCenter(center mgl.Vec3) {
	c.center = center
}

func (c *cameraZ) View() mgl.Mat4 {
	c.view = mgl.LookAtV(c.eye, c.eye.Sub(c.center), c.up)
	return c.view
}

func (c *cameraZ) Move(d Direction, deltaTime float32) {
	velocity := c.speed * deltaTime
	switch d {
	case Forward:
		c.zMinus(velocity)
	case Backward:
		c.zPlus(velocity)
	case Right:
		c.xPlus(velocity)
	case Left:
		c.xMinus(velocity)
	case Up:
		c.yPlus(velocity)
	case Down:
		c.yMinus(velocity)
	}
}

func (c *cameraZ) xPlus(velocity float32) {
	c.eye = c.eye.Sub(c.center.Cross(c.up).Normalize().Mul(velocity))
}

func (c *cameraZ) xMinus(velocity float32) {
	c.eye = c.eye.Add(c.center.Cross(c.up).Normalize().Mul(velocity))
}

func (c *cameraZ) yPlus(velocity float32) {
	c.eye = c.eye.Add(mgl.Vec3{0, velocity, 0})
}

func (c *cameraZ) yMinus(velocity float32) {
	c.eye = c.eye.Sub(mgl.Vec3{0, velocity, 0})
}

func (c *cameraZ) zPlus(velocity float32) {
	c.eye = c.eye.Add(c.center.Mul(velocity))
}

func (c *cameraZ) zMinus(velocity float32) {
	c.eye = c.eye.Sub(c.center.Mul(velocity))
}
