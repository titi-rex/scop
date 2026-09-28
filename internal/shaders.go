package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//go:embed shaders/vertex.glsl
var vertexShaderRaw string

//go:embed shaders/fragment.glsl
var fragmentShaderRaw string

func compileShader(shaderType uint32, name string, sources ...string) (uint32, error) {
	// create shader object
	id := gl.CreateShader(shaderType)

	// convert source to C-string
	csources, free := gl.Strs(sources...)
	// load (replace) the source inside shader
	gl.ShaderSource(id, 1, csources, nil)
	// free C-string
	free()
	// compile the shader
	gl.CompileShader(id)

	// check compilation result and logs
	var status, logLen int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)
	gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLen)

	// output log
	if logLen > 0 {
		msg := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(id, logLen, nil, gl.Str(msg))
		log.Println(name, "shader compilation logs:", msg)
	}

	if status == gl.FALSE {
		return 0, fmt.Errorf("failed to compile shader %s", name)
	}

	return id, nil
}

// from got
//
// vertex.glsl
// #version 410
// in vec3 vp;
// void main() {
//     gl_Position = vec4(vp, 1.0);
// }
//
// fragment.glsl
// #version 410
// out vec4 frag_colour;
// void main() {
//     frag_colour = vec4(49, 18, 1, 1);
// }
//
//
// from ct
//
// vertex.glsl
// #version 330 core
// layout(location = 0) in vec3 vertexPosition;
// in vec3 vp;
// void main() {
//     gl_Position.xyz = vertexPosition;
//     gl_Position.w = 1.0;
// }
//
// fragment.glsl
// #version 330 core
// out vec3 color;
// void main(){
//   color = vec3(1,0,0);
// }
