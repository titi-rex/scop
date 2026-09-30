package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// initOpenGL initializes OpenGL and returns an initialized program.
func initOpenGL() (uint32, error) {
	// init openGL after a current context
	if err := gl.Init(); err != nil {
		return 0, fmt.Errorf("failed to init openGL: %w", err)
	}

	// log openGL version
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	// compile vertex shader
	vertexShader, err := compileShader(gl.VERTEX_SHADER, "vertouille", vertexShaderRaw)
	if err != nil {
		return 0, fmt.Errorf("failed to init openGL: %w", err)
	}

	// compile fragment shader
	fragmentShader, err := compileShader(gl.FRAGMENT_SHADER, "fragouille", fragmentShaderRaw)
	if err != nil {
		return 0, fmt.Errorf("failed to init openGL: %w", err)
	}

	// init program object
	prog, err := initProgram(vertexShader, fragmentShader)
	if err != nil {
		return 0, fmt.Errorf("failed to init openGL: %w", err)
	}

	return prog, nil
}

// create and init a program object
func initProgram(shaders ...uint32) (uint32, error) {
	// create a program object
	id := gl.CreateProgram()

	// attache shaders to program object
	for _, s := range shaders {
		gl.AttachShader(id, s)
	}

	// link a program object
	gl.LinkProgram(id)

	// retrieve the program status and logs
	var status, logLen int32
	gl.GetProgramiv(id, gl.LINK_STATUS, &status)
	gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &logLen)

	// print log if it exist
	if logLen > 0 {
		msg := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(id, logLen, nil, gl.Str(msg))
		log.Println("program link logs:", msg)
	}

	// cleanup shaders objects
	for _, s := range shaders {
		// gl.DetachShader(id, s)
		gl.DeleteShader(s)
	}

	// check link status
	if status == gl.FALSE {
		return 0, fmt.Errorf("failed to link program")
	}

	return id, nil
}

type bufferObjectCn struct {
	target uint32
	usage  uint32
}

var targetName = map[uint32]string{
	gl.ARRAY_BUFFER: "Array Buffer",
	gl.ELEMENT_ARRAY_BUFFER: "Element Array Buffer",
}

func NewBufferObject[T uint32 | float32](target uint32, data []T) uint32 {
	var id uint32

	log.Println("create new Buffer Object for target:", targetName[target])

	gl.GenBuffers(1, &id)
	gl.BindBuffer(target, id)
	if data != nil {
		gl.BufferData(target, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
	}

	return id
}


func NewVAO() uint32 {
	// tell how to use the data in gl.ARRAY_BUFFER
	// vertex array object, store the format of vertex data and buffer object
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao) // only one target to bind (gl.ARRAY_BUFFER)
	// enable atribute (disable by default)
	gl.EnableVertexAttribArray(0)

	// set the attributes for the vbo
	// AFTER vbo/ebo creation
	// gl.VertexAttribPointer( // bind attribute index to buffer currently bound to gl.ARRAY_BUFFER!!
	// 	0,        // attribute 0. Arbitrary, but must match shader's layout.
	// 	3,        // number of componant
	// 	gl.FLOAT, // type
	// 	false,    // normalized?
	// 	0,        // stride, offset between element
	// 	nil,      // array buffer offset
	// )

	return vao
}
