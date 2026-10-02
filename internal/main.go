package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"scop/internal/camera"
	"scop/internal/parser"
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
	cam = camera.NewBasic(mgl.Vec3{0, 0, -9}, mgl.Vec3{0, 0, -1})
)

const (
	uVarModel      = "model"
	uVarView       = "view"
	uVarProjection = "projection"
)

func main() {
	// _, err := parser.NewObj("./resources/42.obj")
	// if err != nil {
	// 	fmt.Printf("error loading obj: %s\n", err.Error())
	// 	return
	// }
	// return

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

	err = initOpenGL()
	if err != nil {
		fmt.Printf("end: %s\n", err.Error())
		return
	}

	shader, err := initShader()
	if err != nil {
		fmt.Printf("end: %s\n", err.Error())
		return
	}

	// init vao
	vao := NewVAO()
	defer func() {
		gl.DeleteVertexArrays(1, &vao)
	}()

	// load object from file
	obj, err := parser.NewObj("./resources/42.obj")
	if err != nil {
		fmt.Printf("error loading obj: %s\n", err.Error())
		return
	}

	// load obj vertice into vbo
	vbo := NewBufferObject(gl.ARRAY_BUFFER, obj.Vertices)
	defer func() {
		gl.DeleteBuffers(1, &vbo)
	}()

	// load obj face into ebo
	ebo := NewBufferObject(gl.ELEMENT_ARRAY_BUFFER, obj.Indices)
	defer func() {
		gl.DeleteBuffers(1, &ebo)
	}()

	// vao.2
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)
	// gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
	// gl.EnableVertexAttribArray(1)

	// wireframe mode
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

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
		shader.Use()

		// compute view
		viewMat = cam.View()

		// compute projection
		projMat = mgl.Perspective(mgl.DegToRad(float32(fov)), wWidth/wHeight, 0.1, 100)

		// set matrice
		err = errors.Join(
			shader.SetUVMat4(uVarModel, modelMat),
			shader.SetUVMat4(uVarView, viewMat),
			shader.SetUVMat4(uVarProjection, projMat),
		)
		if err != nil {
			log.Println("bad code!", err.Error())
		}

		// set the vao to draw
		gl.BindVertexArray(vao)
		// draw
		gl.DrawElements(gl.TRIANGLES, obj.ICount, gl.UNSIGNED_INT, nil)

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

func initShader() (Shader, error) {
	// create shader
	shader := NewShader()

	// compile source
	err := errors.Join(
		shader.CompileSource(gl.VERTEX_SHADER, "vertouille", vertexShaderRaw),
		shader.CompileSource(gl.FRAGMENT_SHADER, "fragouille", fragmentShaderRaw),
	)
	if err != nil {
		return Shader{}, err
	}

	// link shader
	err = shader.Link()
	if err != nil {
		return Shader{}, err
	}

	// shader variable location
	err = shader.SearchUVAll(uVarModel, uVarView, uVarProjection)
	if err != nil {
		return Shader{}, err
	}

	return shader, nil
}
