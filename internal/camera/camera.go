package camera

import mgl "github.com/go-gl/mathgl/mgl32"

type Camera interface {
	// Return the view matrice.
	View() mgl.Mat4
	// Set the center (direction of look)
	SetCenter(center mgl.Vec3)
	// Move the camera in a direction
	Move(d Direction, deltaTime float32)
}

type Direction int

const (
	Forward Direction = iota
	Backward
	Right
	Left
	Up
	Down
)
