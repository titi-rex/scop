package shaders

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	UniVarModel      = "model"
	UniVarView       = "view"
	UniVarProjection = "projection"
)

func InitShader() (Shader, error) {
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
	err = shader.SearchUniVarAll(
		UniVarModel,
		UniVarView,
		UniVarProjection,
	)
	if err != nil {
		return Shader{}, err
	}

	return shader, nil
}
