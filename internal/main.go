package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
)

const (
	width  = 500
	height = 500
	fps    = 10
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()
	cells := makeWorld()

	for !window.ShouldClose() {
		t := time.Now()
		gCount = 0

		for _, rows := range cells {
			for _, c := range rows {
				c.checkState(cells)
			}
		}

		draw(cells, window, program)

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
		if gCount == 0 {
			break
		}
	}
}

func draw(cells World, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for _, rows := range cells {
		for _, c := range rows {
			c.draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}
