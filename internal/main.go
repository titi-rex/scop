package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"scop/internal/camera"
	"scop/internal/ogm"
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"
	"scop/internal/parser"
	"scop/internal/shaders"
	"scop/internal/texture"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
)

const (
	fps       = 31
	radius    = float32(10)
	baseSpeed = 2.5
)

var (
	cam        = camera.NewOrhto(vector.Vector{0, 0, 3}, vector.Vector{0, 0, 1})
	loopStatus = false
	deltaTime  = float32(0)
)

func main() {
	runtime.LockOSThread()

	if len(os.Args) != 2 {
		fmt.Println("usage: ./scop file.obj")
		return
	}

	// load object from file
	obj, err := parser.NewObj(os.Args[1])
	if err != nil {
		fmt.Printf("error loading obj: %s\n", err.Error())
		return
	}

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

	shader, err := shaders.InitShader()
	if err != nil {
		fmt.Printf("end: %s\n", err.Error())
		return
	}

	// init vao
	vao := NewVAO()
	defer func() {
		gl.DeleteVertexArrays(1, &vao)
	}()

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
	// attribute for vertex
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	// attributes for texture
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)
	UnbindVAO()

	// wireframe mode
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// load image
	r, err := texture.NewImg("./resources/wall.jpg")
	if err != nil {
		fmt.Printf("stop: texture: %s\n", err.Error())
		return
	}

	// create texture
	tid := texture.NewTexture(r)

	// set texture univar
	err = errors.Join(
		shader.SetUniVarUInt(shaders.UniVarTexture1, 0),
		// shader.SetUniVarUInt(shaders.UniVarMode, 1),
	)
	if err != nil {
		fmt.Printf("stop: univar textures : %s\n", err.Error())
		return
	}

	// time
	var lastFrame float64

	// model
	const rotationSpeed = 1.25
	const moveSpeed = 0.05
	angle := vector.Vector{0, 0, 0}
	pos := vector.Vector{}
	autoRotate := vector.Vector{}

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
		projMat := matrice.NewPerspective(ogm.Rad(float32(fov)), wWidth/wHeight, 0.1, 10)

		// compute autorotation
		angle.Add(autoRotate)

		// compute model matrice

		m0 := matrice.NewIdent().
			Translate(pos).
			RotateZ(ogm.Rad(angle.Z())).
			RotateY(ogm.Rad(angle.Y())).
			RotateX(ogm.Rad(angle.X()))

		// set matrice
		err = errors.Join(
			shader.SetUniVarMatrice(shaders.UniVarModel, m0),
			shader.SetUniVarMatrice(shaders.UniVarView, cam.View()),
			shader.SetUniVarMatrice(shaders.UniVarProjection, projMat),
		)
		if err != nil {
			log.Println("bad code!", err.Error())
			return
		}

		// draw
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, tid)
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, obj.IndicesCount, gl.UNSIGNED_INT, nil)
		UnbindVAO()
		window.SwapBuffers()

		// handle event
		glfw.PollEvents()
		doKeyPress(glfw.KeyEscape, func() { loopStatus = true })
		doKeyPress(glfw.KeyQ, func() { loopStatus = true })
		doKeyPress(glfw.KeyR, func() {
			angle = vector.Vector{}
			autoRotate = vector.Vector{}
			pos = vector.Vector{}
		})
		doKeyPress(glfw.KeyUp, func() { angle[0] += rotationSpeed })
		doKeyPress(glfw.KeyDown, func() { angle[0] -= rotationSpeed })
		doKeyPress(glfw.KeyRight, func() { angle[1] += rotationSpeed })
		doKeyPress(glfw.KeyLeft, func() { angle[1] -= rotationSpeed })
		doKeyPress(glfw.KeyZ, func() { angle[2] += rotationSpeed })
		doKeyPress(glfw.KeyX, func() { angle[2] -= rotationSpeed })
		doKeyPress(glfw.Key1, func() { autoRotate[0] += rotationSpeed })
		doKeyPress(glfw.Key2, func() { autoRotate[1] += rotationSpeed })
		doKeyPress(glfw.Key3, func() { autoRotate[2] += rotationSpeed })
		doKeyPress(glfw.KeySpace, func() { autoRotate = vector.Vector{} })
		doKeyPress(glfw.KeyD, func() { pos[0] += moveSpeed })
		doKeyPress(glfw.KeyA, func() { pos[0] -= moveSpeed })
		doKeyPress(glfw.KeyV, func() { pos[1] += moveSpeed })
		doKeyPress(glfw.KeyC, func() { pos[1] -= moveSpeed })
		doKeyPress(glfw.KeyS, func() { pos[2] += moveSpeed })
		doKeyPress(glfw.KeyW, func() { pos[2] -= moveSpeed })

		// cap fps
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func doKeyPress(k glfw.Key, do func()) {
	if glfw.GetCurrentContext().GetKey(k) == glfw.Press {
		do()
	}
}
