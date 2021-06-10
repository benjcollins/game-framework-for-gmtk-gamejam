package graphics

import (
	"errors"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Uniform interface{}

type Program struct {
	ID       uint32
	uniforms map[string]int32
}

func CreateProgramVSFS(vs VertexShader, fs FragmentShader) (Program, error) {
	return CreateProgram([]Shader{vs.Shader, fs.Shader})
}

func CreateProgramVSGSFS(vs VertexShader, gs GeometryShader, fs FragmentShader) (Program, error) {
	return CreateProgram([]Shader{vs.Shader, gs.Shader, fs.Shader})
}

func CreateProgram(shaders []Shader) (Program, error) {
	programID := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(programID, shader.shaderID)
	}
	gl.LinkProgram(programID)
	for _, shader := range shaders {
		gl.DetachShader(programID, shader.shaderID)
	}
	status := int32(0)
	gl.GetProgramiv(programID, gl.LINK_STATUS, &status)
	if status == gl.TRUE {
		return Program{programID, make(map[string]int32)}, nil
	}
	logLength := int32(0)
	gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength)+1)
	gl.GetProgramInfoLog(programID, logLength, nil, gl.Str(log))
	return Program{}, errors.New(log)
}

func (program *Program) Bind(uniforms map[string]Uniform) {
	for name := range uniforms {
		if _, ok := program.uniforms[name]; !ok {
			location := gl.GetUniformLocation(program.ID, gl.Str(name+"\x00"))
			if location == -1 {
				log.Fatalf("invalid uniform name '%s'", name)
			}
			program.uniforms[name] = location
		}
	}
	gl.UseProgram(program.ID)
	for name, data := range uniforms {
		location := program.uniforms[name]
		switch data := data.(type) {
		case int:
			gl.Uniform1i(location, int32(data))
		case float32:
			gl.Uniform1f(location, data)
		case mgl32.Vec2:
			gl.Uniform2fv(location, 1, &data[0])
		case mgl32.Vec3:
			gl.Uniform3fv(location, 1, &data[0])
		case mgl32.Vec4:
			gl.Uniform4fv(location, 1, &data[0])
		case mgl32.Mat3:
			gl.UniformMatrix3fv(location, 1, false, &data[0])
		}
	}
}

func (program Program) Delete() {
	gl.DeleteProgram(program.ID)
}
