package graphics

import (
	_ "embed"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed shaders/sprite.vs
var SpriteShaderVS string

//go:embed shaders/sprite.fs
var SpriteShaderFS string

type SpriteRenderer struct {
	program *Program
}

type SpriteVertex struct {
	pos    mgl32.Vec2
	tx, ty float32
}

type Sprite struct {
	transform           mgl32.Mat3
	texOffX, texOffY    float32
	texWidth, texHeight float32
}

type SpriteBuffer struct {
	vao      uint32
	vbo, ibo Buffer
	size     int
}

func NewSprite(transform mgl32.Mat3) Sprite {
	return Sprite{transform, 0, 0, 1, 1}
}

func NewSpriteFromAtlas(transform mgl32.Mat3, texOffX, texOffY, texWidth, texHeight float32) Sprite {
	return Sprite{transform, texOffX, texOffY, texWidth, texHeight}
}

func CreateSpriteBuffer(sprites []Sprite) SpriteBuffer {

	vao := uint32(0)
	gl.CreateVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vbo := CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.ID)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, int32(unsafe.Sizeof(SpriteVertex{})), 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, int32(unsafe.Sizeof(SpriteVertex{})), 2*4)

	ibo := CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo.ID)

	buffer := SpriteBuffer{vao, vbo, ibo, 0}

	buffer.SetContents(sprites)

	return buffer
}

func (buffer *SpriteBuffer) SetContents(sprites []Sprite) {

	indicies := spritesToIndicies(sprites, 0)
	vertices := spritesToVertices(sprites)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.vbo.ID)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(SpriteVertex{})), unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.ibo.ID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicies)*4, unsafe.Pointer(&indicies[0]), gl.STATIC_DRAW)

	buffer.size = len(indicies)
}

func spritesToVertices(sprites []Sprite) []SpriteVertex {
	vertices := make([]SpriteVertex, len(sprites)*4)

	for i, sprite := range sprites {
		copy(vertices[i*4:i*4+4], []SpriteVertex{
			{sprite.transform.Mul3x1(mgl32.Vec3{-1.0, 1.0, 1.0}).Vec2(), sprite.texOffX, sprite.texOffY},
			{sprite.transform.Mul3x1(mgl32.Vec3{1.0, 1.0, 1.0}).Vec2(), sprite.texOffX + sprite.texHeight, sprite.texOffY},
			{sprite.transform.Mul3x1(mgl32.Vec3{1.0, -1.0, 1.0}).Vec2(), sprite.texOffX + sprite.texHeight, sprite.texOffY + sprite.texHeight},
			{sprite.transform.Mul3x1(mgl32.Vec3{-1.0, -1.0, 1.0}).Vec2(), sprite.texOffX, sprite.texOffY + sprite.texHeight},
		})
	}

	return vertices
}

func spritesToIndicies(sprites []Sprite, start uint32) []uint32 {
	indicies := make([]uint32, len(sprites)*6)

	for i := range sprites {
		j := start + uint32(i*4)
		copy(indicies[i*6:i*6+6], []uint32{
			j, j + 1, j + 2,
			j, j + 2, j + 3,
		})
	}

	return indicies
}

func (buffer *SpriteBuffer) UpdateContents(start int, sprites []Sprite) {
	indicies := spritesToIndicies(sprites, uint32(start)*4)
	vertices := spritesToVertices(sprites)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.vbo.ID)
	gl.BufferSubData(gl.ARRAY_BUFFER, start*int(unsafe.Sizeof(SpriteVertex{}))*4, len(vertices)*int(unsafe.Sizeof(SpriteVertex{})), unsafe.Pointer(&vertices[0]))

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.ibo.ID)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, start*4*6, len(indicies)*4, unsafe.Pointer(&indicies[0]))
}

func CreateRenderer() *SpriteRenderer {
	fs, err := CreateFragmentShader(SpriteShaderFS)
	if err != nil {
		log.Fatal(err)
	}
	vs, err := CreateVertexShader(SpriteShaderVS)
	if err != nil {
		log.Fatal(err)
	}
	program, err := CreateProgramVSFS(vs, fs)
	if err != nil {
		log.Fatal(err)
	}
	vs.Delete()
	fs.Delete()
	return &SpriteRenderer{program}
}

func (renderer *SpriteRenderer) Render(texture Texture, sprites SpriteBuffer, transform mgl32.Mat3, aspectRatio float32) {

	texture.Bind(0)

	gl.BindVertexArray(sprites.vao)

	renderer.program.Bind(map[string]interface{}{
		"textureSampler": 0,
		"transform":      ComputeAspectRatio(aspectRatio).Mul3(transform),
	})

	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(sprites.size), gl.UNSIGNED_INT, 0)
}

func ComputeAspectRatio(aspectRatio float32) mgl32.Mat3 {
	if aspectRatio < 1 {
		return mgl32.Scale2D(1/aspectRatio, 1)
	}
	return mgl32.Scale2D(1, aspectRatio)
}

func (renderer SpriteRenderer) Delete() {
	renderer.program.Delete()
}

func (buffer SpriteBuffer) Delete() {
	buffer.vbo.Delete()
	buffer.ibo.Delete()
	gl.DeleteVertexArrays(1, &buffer.vao)
}
