package main

import (
	"errors"
	"fmt"
	"log"
	"os"
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
	cam = camera.NewOrhto(mgl.Vec3{0, 0, 3}, mgl.Vec3{0, 0, 1})
	loopStatus = false
	deltaTime  = float32(0)
)

const (
	uVarModel      = "model"
	uVarView       = "view"
	uVarProjection = "projection"
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
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT);
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT);
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);

	// os.Open("./resources/wall.png")


	// time
	var lastFrame float64

	// model
	const rotationSpeed = 1.5
	const moveSpeed = 0.05
	angle := mgl.Vec3{0, -90, 0}
	pos := mgl.Vec3{}

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

		// compute model matrice
		rx := mgl.HomogRotate3DX(Rad(angle.X()))
		ry := mgl.HomogRotate3DY(Rad(angle.Y()))
		rz := mgl.HomogRotate3DZ(Rad(angle.Z()))
		tr := mgl.Translate3D(pos.X(), pos.Y(), pos.Z())
		m1 := tr.Mul4(rz.Mul4(ry.Mul4(rx)))

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

		// draw
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, obj.IndicesCount, gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)
		window.SwapBuffers()

		// handle event
		glfw.PollEvents()
		doKeyPress(glfw.KeyEscape, func() { loopStatus = true })
		doKeyPress(glfw.KeyQ, func() { loopStatus = true })
		doKeyPress(glfw.KeyR, func() {
			angle = mgl.Vec3{0, -90, 0}
			pos = mgl.Vec3{}
		})
		doKeyPress(glfw.KeyX, func() { angle[0] += rotationSpeed })
		doKeyPress(glfw.KeyY, func() { angle[1] += rotationSpeed })
		doKeyPress(glfw.KeyZ, func() { angle[2] += rotationSpeed })
		doKeyPress(glfw.KeyD, func() { pos[0] += moveSpeed })
		doKeyPress(glfw.KeyA, func() { pos[0] -= moveSpeed })
		doKeyPress(glfw.KeySpace, func() { pos[1] += moveSpeed })
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
