package camera

import (
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"
)

type core struct {
	eye    *vector.Vector
	center *vector.Vector
	up     *vector.Vector

	speed float32
}

func (c *core) view() matrice.Matrice {
	return matrice.LookAt(*c.eye, *c.eye.Sub(*c.center), *c.up)
}

// xPlus increase the core X position
func (c *core) xPlus(velocity float32) {
	c.eye.Sub(*c.center.Cross(*c.up).Normalize().Mul(velocity))
}

// xMinus decrease the core X position
func (c *core) xMinus(velocity float32) {
	c.eye.Add(*c.center.Cross(*c.up).Normalize().Mul(velocity))
}

// zPlus increase the core Z position
func (c *core) zPlus(velocity float32) {
	c.eye.Add(*c.center.Mul(velocity))
}

// zMinus decrease the core Z position
func (c *core) zMinus(velocity float32) {
	c.eye.Sub(*c.center.Mul(velocity))
}

// xOrthoPlus increase the core Y position
// regardless the core orientation
func (c *core) xOrthoPlus(velocity float32) {
	c.eye.Add(vector.Vector{velocity, 0, 0})
}

// xOrthoMinus decrease the core Y position
// regardless the core orientation
func (c *core) xOrthoMinus(velocity float32) {
	c.eye.Sub(vector.Vector{velocity, 0, 0})
}

// yOrthoPlus increase the core Y position
// regardless the core orientation
func (c *core) yOrthoPlus(velocity float32) {
	c.eye.Add(vector.Vector{0, velocity, 0})
}

// yOrthoMinus decrease the core Y position
// regardless the core orientation
func (c *core) yOrthoMinus(velocity float32) {
	c.eye.Sub(vector.Vector{0, velocity, 0})
}

// zOrthoPlus increase the core Y position
// regardless the core orientation
func (c *core) zOrthoPlus(velocity float32) {
	c.eye.Add(vector.Vector{0, 0, velocity})
}

// zOrthoMinus decrease the core Y position
// regardless the core orientation
func (c *core) zOrthoMinus(velocity float32) {
	c.eye.Sub(vector.Vector{0, 0, velocity})
}
