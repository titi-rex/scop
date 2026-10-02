package main

import (
	"fmt"
	"log"
	"runtime"
	"scop/internal/camera"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/go-gl/glfw/v3.4/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	fps       = 31
	fSize     = 4
	radius    = float32(10)
	baseSpeed = 2.5
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

	cam = camera.NewBasic(mgl.Vec3{0, 0, -3}, mgl.Vec3{0, 0, -1})
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

	log.Println("enable Depth Test")
	gl.Enable((gl.DEPTH_TEST))

	var vao, vbo, ebo uint32

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
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// shader variable location
	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	// model matrice
	modelMat := mgl.HomogRotate3D(70, mgl.Vec3{1, 0, 0})

	var viewMat, projMat mgl.Mat4

	// time
	var lastFrame float64

	q := false
	for !window.ShouldClose() && !q {
		t := time.Now() // used for fps cap
		currentFrame := glfw.GetTime()
		deltaTime := float32(currentFrame - lastFrame)
		lastFrame = currentFrame

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// use program shader
		gl.UseProgram(program)

		// compute view
		viewMat = cam.View()

		// compute projection
		projMat = mgl.Perspective(mgl.DegToRad(float32(fov)), wWidth/wHeight, 0.1, 100)

		// set matrice
		gl.UniformMatrix4fv(modelLoc, 1, false, &modelMat[0])
		gl.UniformMatrix4fv(viewLoc, 1, false, &viewMat[0])
		gl.UniformMatrix4fv(projLoc, 1, false, &projMat[0])

		// set the vao to draw
		gl.BindVertexArray(vao)
		// draw
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		// gl.DrawArrays(gl.TRIANGLES, 0, 3)

		// Push buffer to window
		window.SwapBuffers()

		// poll event
		glfw.PollEvents()

		// do smtg on key press
		doKeyPress(glfw.KeyEscape, func() { q = true })

		// Z
		doKeyPress(glfw.KeyW, func() { cam.Move(camera.Forward, float32(deltaTime)) })
		doKeyPress(glfw.KeyS, func() { cam.Move(camera.Backward, deltaTime) })

		// X
		doKeyPress(glfw.KeyD, func() { cam.Move(camera.Right, deltaTime) })
		doKeyPress(glfw.KeyA, func() { cam.Move(camera.Left, deltaTime) })

		// Y
		doKeyPress(glfw.KeyH, func() { cam.Move(camera.Up, deltaTime) })
		doKeyPress(glfw.KeyN, func() { cam.Move(camera.Down, deltaTime) })

		// cap fps
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func doKeyPress(k glfw.Key, do func()) {
	if glfw.GetCurrentContext().GetKey(k) == glfw.Press {
		// log.Println("key pressed", k)
		do()
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
// // transformation matrice
// mat := mgl.Ident4()
// // mat = mat.Mul4(mgl.Translate3D(1, 1, 0))
// mat = mat.Mul4(mgl.HomogRotate3DZ(mgl.DegToRad(float32(glfw.GetTime()))))
// mat = mat.Mul4(mgl.Scale3D(0.5, 0.5, 0.5))

// rotatating view
