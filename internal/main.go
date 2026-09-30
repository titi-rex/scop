package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/go-gl/glfw/v3.4/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	fps   = 10
	fSize = 4
)

var (
	triangle = []float32{
		0, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}

	vertices = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom let
		-0.5, 0.5, 0.0, // top let
	}

	indices = []uint32{
		0, 1, 3, // premier triangle
		1, 2, 3, // second triangle
	}
)

func main() {
	runtime.LockOSThread()

	if err := initGlfw(); err != nil {
		fmt.Printf("end: %s\n", err.Error())
		return
	}
	defer glfw.Terminate()

	window, err := initWindow()
	if err != nil {
		fmt.Printf("end: %s\n", err.Error())
		return
	}

	program, err := initOpenGL()
	if err != nil {
		fmt.Printf("end: %s\n", err.Error())
		return
	}

	var vao, vbo, ebo uint32

	log.Println("begin")

	// vao.1
	vao = NewVAO()
	defer func() {
		gl.DeleteVertexArrays(1, &vao)
	}()

	// vbo
	vbo = NewBufferObject(gl.ARRAY_BUFFER, vertices)
	defer func() {
		gl.DeleteBuffers(1, &vbo)
	}()

	// ebo
	ebo = NewBufferObject(gl.ELEMENT_ARRAY_BUFFER, indices)
	defer func() {
		gl.DeleteBuffers(1, &ebo)
	}()

	// vao.2
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)

	// wireframe mode
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	ctr := gl.Str("transform\x00")

	for !window.ShouldClose() {
		t := time.Now() // used for fps cap

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// use program shader
		gl.UseProgram(program)

		// transformation matrice
		mat := mgl.Ident4()
		// mat = mat.Mul4(mgl.Translate3D(1, 1, 0))
		mat = mat.Mul4(mgl.HomogRotate3DZ(mgl.DegToRad(float32(glfw.GetTime()))))
		mat = mat.Mul4(mgl.Scale3D(0.5, 0.5, 0.5))

		// set matrice
		uloc := gl.GetUniformLocation(program, ctr)
		gl.UniformMatrix4fv(uloc, 1, false, &mat[0])

		// set the vao to draw
		gl.BindVertexArray(vao)
		// draw
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		// gl.DrawArrays(gl.TRIANGLES, 0, 3)

		// Push buffer to window
		window.SwapBuffers()

		// poll event
		glfw.PollEvents()

		// exit on Escape press
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			break
		}

		// cap fps
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

// func draw(cells World, window *glfw.Window, program uint32) {
// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
// 	gl.UseProgram(program)

// 	for _, rows := range cells {
// 		for _, c := range rows {
// 			c.draw()
// 		}
// 	}

// }

// func (c *Cell) draw() {
// 	if !c.alive {
// 		return
// 	}

// 	gl.BindVertexArray(c.drawable)
// 	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
// }
