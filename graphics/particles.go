package graphics

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	_ "embed"
)

//go:embed shaders/particle.vs
var ParticleShaderVS string

//go:embed shaders/particle.fs
var ParticleShaderFS string

type ParticleRenderer struct {
	program Program
	vbo     Buffer
}

type Particle struct {
	Transform mgl32.Mat3
	Frame     float32
}

type ParticleVertex struct {
	x, y   float32
	tx, ty float32
}

type ParticleSystem struct {
	vao              uint32
	ibo              Buffer
	texture          Texture
	instances        int
	capacity         int
	hFrames, vFrames int
	update           func(*Particle) bool
}

func (system *ParticleSystem) Update() {
	gl.BindBuffer(gl.ARRAY_BUFFER, system.ibo.ID)
	base := gl.MapBuffer(gl.ARRAY_BUFFER, gl.READ_WRITE)
	for i := 0; i < system.instances; {
		p := (*Particle)(unsafe.Pointer(uintptr(base) + uintptr(i)*unsafe.Sizeof(Particle{})))
		if !system.update(p) {
			system.instances -= 1
			lastParticle := (*Particle)(unsafe.Pointer(uintptr(base) + uintptr(system.instances)*unsafe.Sizeof(Particle{})))
			*p = *lastParticle
		} else {
			i++
		}
	}
	gl.UnmapBuffer(gl.ARRAY_BUFFER)
}

func (system *ParticleSystem) AppendParticle(particles ...Particle) {
	gl.BindBuffer(gl.ARRAY_BUFFER, system.ibo.ID)
	base := gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY)
	for i, particle := range particles {
		p := (*Particle)(unsafe.Pointer(uintptr(base) + uintptr(system.instances+i)*unsafe.Sizeof(Particle{})))
		*p = particle
	}
	system.instances += len(particles)
	gl.UnmapBuffer(gl.ARRAY_BUFFER)
}

func NewParticle(transform mgl32.Mat3, frame float32) Particle {
	return Particle{transform, frame}
}

func (renderer *ParticleRenderer) CreateParticleSystem(capacity int, texture Texture, hFrames, vFrames int, update func(*Particle) bool) *ParticleSystem {
	vao := uint32(0)
	gl.CreateVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vbo.ID)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, int32(unsafe.Sizeof(ParticleVertex{})), 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, int32(unsafe.Sizeof(ParticleVertex{})), 4*2)

	ibo := CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, ibo.ID)
	gl.BufferData(gl.ARRAY_BUFFER, capacity*int(unsafe.Sizeof(Particle{})), nil, gl.DYNAMIC_DRAW)
	for i := uint32(0); i < 3; i++ {
		gl.EnableVertexAttribArray(i + 2)
		gl.VertexAttribPointerWithOffset(i+2, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), uintptr(i*3*4))
		gl.VertexAttribDivisor(i+2, 1)
	}
	gl.EnableVertexAttribArray(5)
	gl.VertexAttribPointerWithOffset(5, 1, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), uintptr(9*4))
	gl.VertexAttribDivisor(5, 1)

	return &ParticleSystem{vao, ibo, texture, 0, capacity, hFrames, vFrames, update}
}

func CreateParticleRenderer() *ParticleRenderer {
	renderer := ParticleRenderer{}
	vs, err := CreateVertexShader(ParticleShaderVS)
	if err != nil {
		log.Fatal(err)
	}
	fs, err := CreateFragmentShader(ParticleShaderFS)
	if err != nil {
		log.Fatal(err)
	}
	renderer.program, err = CreateProgramVSFS(vs, fs)
	if err != nil {
		log.Fatal(err)
	}
	vs.Delete()
	fs.Delete()

	vertices := []ParticleVertex{
		{-1, 1, 0, 0},
		{1, 1, 1, 0},
		{-1, -1, 0, 1},
		{1, -1, 1, 1},
	}

	renderer.vbo = CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vbo.ID)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(ParticleVertex{})), unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	return &renderer
}

func (renderer *ParticleRenderer) Render(particles *ParticleSystem, transform mgl32.Mat3) {
	renderer.program.Bind(map[string]Uniform{
		"textureSampler":  0,
		"hFrames":         particles.hFrames,
		"vFrames":         particles.vFrames,
		"globalTransform": transform,
	})

	gl.BindVertexArray(particles.vao)

	particles.texture.Bind(0)

	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, int32(particles.instances))
}
