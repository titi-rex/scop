package shaders

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"scop/internal/ogm/matrice"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	ErrLocNotFound    = errors.New("uniform location not found")
	ErrLocNotRegister = errors.New("uniform location not searched")
)

//go:embed vertex.glsl
var vertexShaderRaw string

//go:embed fragment.glsl
var fragmentShaderRaw string

// Shader is the shader program.
type Shader struct {
	id      uint32
	shaders map[string]uint32
	uniLoc  map[string]int32
}

// NewShader return a [Shader] without any source.
func NewShader() Shader {
	return Shader{
		id:      gl.CreateProgram(),
		shaders: make(map[string]uint32),
		uniLoc:  make(map[string]int32),
	}
}

// CompileSource add a named source shader to the shader program and
// compile it.
func (s *Shader) CompileSource(shaderType uint32, name string, source string) error {
	id := gl.CreateShader(shaderType)

	// convert source to C-string
	csources, free := gl.Strs(source + "\x00")

	// load (replace) the source inside shader
	gl.ShaderSource(id, 1, csources, nil)

	// free C-string
	free()

	// compile the shader
	gl.CompileShader(id)

	// check compilation result and logs
	var status int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		// output log
		var logLen int32

		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLen)
		if logLen > 0 {
			msg := strings.Repeat("\x00", int(logLen+1))
			gl.GetShaderInfoLog(id, logLen, nil, gl.Str(msg))
			log.Println(name, "shader compilation logs:", msg)
		}

		return fmt.Errorf("failed to compile shader %s", name)
	}

	s.shaders[name] = id

	return nil
}

// Link link the shader program
func (s *Shader) Link() error {
	// attach shaders to program
	for _, shader := range s.shaders {
		gl.AttachShader(s.id, shader)
	}
	// cleanup shaders objects
	defer func() {
		for _, shader := range s.shaders {
			gl.DetachShader(s.id, shader)
			gl.DeleteShader(shader)
		}
	}()

	// link shader program
	gl.LinkProgram(s.id)

	// check link status
	var status int32
	gl.GetProgramiv(s.id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(s.id, gl.INFO_LOG_LENGTH, &logLen)
		if logLen > 0 {
			msg := strings.Repeat("\x00", int(logLen+1))
			gl.GetProgramInfoLog(s.id, logLen, nil, gl.Str(msg))
			log.Println("program link logs:", msg)
		}
		return fmt.Errorf("failed to link program")
	}

	return nil
}

// Use set the shader to be used by OpenGL
func (s *Shader) Use() {
	gl.UseProgram(s.id)
}

// SearchUniVar register a var name as a Uniform Variable.
// Check if name is found.
func (s *Shader) SearchUniVar(name string) error {
	loc := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	if loc == -1 {
		return fmt.Errorf("%s: %w", name, ErrLocNotFound)
	}

	s.uniLoc[name] = loc
	log.Println("register uvar", name, "at", loc)

	return nil
}

// SearchUniVarAll try to register alls names as Uniform Variable
func (s *Shader) SearchUniVarAll(names ...string) error {
	var errs []error
	for _, name := range names {
		errs = append(errs, s.SearchUniVar(name))
	}

	err := errors.Join(errs...)
	if err != nil {
		return err
	}

	return nil
}

// SetUniVarMatrice set a UV name to the value from m
func (s *Shader) SetUniVarMatrice(name string, m *matrice.Matrice) error {
	loc, ok := s.uniLoc[name]
	if !ok {
		return ErrLocNotRegister
	}

	gl.UniformMatrix4fv(loc, 1, false, m.Data())

	return nil
}

func (s *Shader) SetUniVarUInt(name string, v uint32) error {
	loc, ok := s.uniLoc[name]
	if !ok {
		return ErrLocNotRegister
	}

	gl.Uniform1ui(loc, v)

	return nil
}

func (s *Shader) SetUniVarInt(name string, v int32) error {
	loc, ok := s.uniLoc[name]
	if !ok {
		return ErrLocNotRegister
	}

	gl.Uniform1i(loc, v)

	return nil
}

func (s *Shader) SetUniVarFloat(name string, v float32) error {
	loc, ok := s.uniLoc[name]
	if !ok {
		return ErrLocNotRegister
	}

	gl.Uniform1f(loc, v)

	return nil
}
