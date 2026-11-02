package main

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	targetName = map[uint32]string{
		gl.ARRAY_BUFFER:         "Array Buffer",
		gl.ELEMENT_ARRAY_BUFFER: "Element Array Buffer",
	}
)

// initOpenGL initializes OpenGL.
func initOpenGL() error {
	// init openGL after a current context
	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to init openGL: %w", err)
	}

	// log openGL version
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	var nv int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nv)
	log.Println("Maximum vertex attributes:", nv)

	gl.Enable((gl.DEPTH_TEST))
	log.Println("Enable Depth Test")

	return nil
}

// NewBufferObject create a BO for target and fill it with data
// if data is != nil
// Set usage to [gl.STATIC_DRAW]
func NewBufferObject[T uint32 | int32 | float32](target uint32, data []T) uint32 {
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

func UnbindVAO() {
	gl.BindVertexArray(0)
}
