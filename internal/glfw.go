package main

import (
	"fmt"
	"log"

	"github.com/go-gl/glfw/v3.4/glfw"
)

const (
	wWidth  = 500
	wHeight = 500
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

	return window, nil
}
