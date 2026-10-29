package camera

import (
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"
)

type Camera interface {
	// Return the view matrice.
	View() matrice.Matrice
	// Set the center (direction of look)
	SetCenter(center vector.Vector)
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
