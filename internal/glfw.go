package main

import (
	"fmt"
	"log"
	"math"
	"scop/internal/camera"
	"scop/internal/ogm"
	"scop/internal/ogm/vector"

	"github.com/go-gl/glfw/v3.4/glfw"
)

const (
	wWidth              = 800
	wHeight             = 600
	sensitivity float64 = 0.05
)

var (
	curLastX float64 = wWidth / 2
	curLastY float64 = wHeight / 2

	yaw   float64 = -90
	pitch float64 = 0
	fov   float64 = 45

	firstTime = true
)

// initGlfw initializes glfw.
func initGlfw() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to init glfw: %w", err)
	}

	log.Println("glfw version", glfw.GetVersionString())

	// glfw hints
	// glfw.WindowHint(glfw.Samples, 4)             // antialiasing
	glfw.WindowHint(glfw.Resizable, glfw.False)  // resizable
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // old: 4
	glfw.WindowHint(glfw.ContextVersionMinor, 1) // old: 1

	// openGL hints
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)    // MacOs compat
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // we don't want the old openGL

	return nil
}

// initWindow create a window to use.
func initWindow() (*glfw.Window, error) {
	// create window
	window, err := glfw.CreateWindow(wWidth, wHeight, "window", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	// set the window's context as current context
	window.MakeContextCurrent()

	// set the window input mode
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// window.SetCursorPosCallback(CursorPosHandler)
	// window.SetScrollCallback(ScrollHandler)
	// window.SetKeyCallback(KeyHandler)

	return window, nil
}

func CursorPosHandler(w *glfw.Window, xpos, ypos float64) {
	if firstTime {
		curLastX = xpos
		curLastY = ypos
		firstTime = false
	}
	xoff := xpos - curLastX
	yoff := ypos - curLastY
	curLastX = xpos
	curLastY = ypos
	xoff *= sensitivity
	yoff *= sensitivity

	yaw += xoff
	pitch += yoff

	if pitch > 89.0 {
		pitch = 89.0
	}
	if pitch < -89.0 {
		pitch = -89.0
	}

	front := vector.Vector{
		float32(math.Cos(ogm.Rad(yaw)) * math.Cos(ogm.Rad(pitch))),
		float32(math.Sin(ogm.Rad(pitch))),
		float32(math.Sin(ogm.Rad(yaw)) * math.Cos(ogm.Rad(pitch))),
	}
	cam.SetCenter(vector.Normalize(front))
}

func ScrollHandler(w *glfw.Window, xoffset, yoffset float64) {
	if fov >= 1.0 && fov <= 45.0 {
		fov -= yoffset
	}
	if fov < 1.0 {
		fov = 1.0
	}
	if fov > 45.0 {
		fov = 45.0
	}
}

func KeyHandler(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// if action != glfw.Press {
	// 	return
	// }

	switch key {
	case glfw.KeyEscape:
		loopStatus = true

	// X
	case glfw.KeyW:
		cam.Move(camera.Forward, float32(deltaTime))
	case glfw.KeyS:
		cam.Move(camera.Backward, deltaTime)

	// Y
	case glfw.KeyD:
		cam.Move(camera.Right, deltaTime)
	case glfw.KeyA:
		cam.Move(camera.Left, deltaTime)

	// Z
	case glfw.KeyH:
		cam.Move(camera.Up, deltaTime)
	case glfw.KeyN:
		cam.Move(camera.Down, deltaTime)
	}
}
