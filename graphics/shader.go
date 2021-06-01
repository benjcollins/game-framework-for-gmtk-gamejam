package graphics

import (
	"errors"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	shaderID uint32
}

type VertexShader struct {
	Shader
}

type GeometryShader struct {
	Shader
}

type FragmentShader struct {
	Shader
}

func CreateVertexShader(source string) (VertexShader, error) {
	shader, err := CreateShader(source, gl.VERTEX_SHADER)
	return VertexShader{shader}, err
}

func CreateGeometryShader(source string) (GeometryShader, error) {
	shader, err := CreateShader(source, gl.GEOMETRY_SHADER)
	return GeometryShader{shader}, err
}

func CreateFragmentShader(source string) (FragmentShader, error) {
	shader, err := CreateShader(source, gl.FRAGMENT_SHADER)
	return FragmentShader{shader}, err
}

func CreateShader(source string, ty uint32) (Shader, error) {
	shaderID := gl.CreateShader(ty)
	csource, free := gl.Strs(source)
	length := int32(len(source))
	gl.ShaderSource(shaderID, 1, csource, &length)
	free()
	gl.CompileShader(shaderID)

	status := int32(0)
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &status)
	if status == gl.TRUE {
		return Shader{shaderID}, nil
	}
	logLength := int32(0)
	gl.GetShaderiv(shaderID, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength)+1)
	gl.GetShaderInfoLog(shaderID, logLength, nil, gl.Str(log))
	return Shader{}, errors.New(log)
}

func (shader Shader) Delete() {
	gl.DeleteShader(shader.shaderID)
}
