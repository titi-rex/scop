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
	radius    = float32(10)
	baseSpeed = 2.5
)

var (
	cam = camera.NewBasic(mgl.Vec3{2, 0, -8}, mgl.Vec3{0, 0, 0})
	// cam = camera.NewOrhto(mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 0, 0})
	// cam = camera.NewLocked()
	loopStatus = false
	deltaTime = float32(0)
)

const (
	uVarModel      = "model"
	uVarView       = "view"
	uVarProjection = "projection"
)

func main() {
	// obj1, err := parser.NewObj("./resources/42.obj")
	// if err != nil {
	// 	fmt.Printf("error loading obj: %s\n", err.Error())
	// 	return
	// }
	// log.Println("n sommet", len(obj1.Vertices) / 3)
	// log.Println("n idx", len(obj1.Indices))
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
	gl.BindVertexArray(0)

	// wireframe mode
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// time
	var lastFrame float64
	var ir float64

	for !window.ShouldClose() && !loopStatus {
		t := time.Now() // used for fps cap
		currentFrame := glfw.GetTime()
		deltaTime = float32(currentFrame - lastFrame)
		lastFrame = currentFrame

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// use program shader
		shader.Use()

		// compute projection
		projMat := mgl.Perspective(mgl.DegToRad(float32(fov)), wWidth/wHeight, 0.1, 100)

		// m1 := obj.Matrice.Mul4(mgl.Translate3D(.5,0,0))
		// m1 = m1.Mul4(mgl.HomogRotate3D(float32(glfw.GetTime() * Rad(-90.0)), mgl.Vec3{0,1,0}))
		ir += 0.5
		m1 := RotateY(float32(-90 + ir))

		// set matrice
		err = errors.Join(
			shader.SetUVMat4(uVarModel, m1),
			shader.SetUVMat4(uVarView, cam.View()),
			shader.SetUVMat4(uVarProjection, projMat),
		)
		if err != nil {
			log.Println("bad code!", err.Error())
			return
		}

		// set the vao to draw
		gl.BindVertexArray(vao)
		// draw
		gl.DrawElements(gl.TRIANGLES, obj.IndicesCount, gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)

		// Push buffer to window
		window.SwapBuffers()

		// poll event
		glfw.PollEvents()


		// cap fps
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
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
	err = shader.SearchUVAll(
		uVarModel,
		uVarView,
		uVarProjection,
	)
	if err != nil {
		return Shader{}, err
	}

	return shader, nil
}
