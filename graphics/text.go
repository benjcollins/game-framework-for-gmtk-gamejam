package graphics

import (
	"encoding/json"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	_ "embed"
)

//go:embed shaders/text.vs
var TextVS string

//go:embed shaders/text.fs
var TextFS string

type TextRenderer struct {
	program Program
}

type Font struct {
	texture Texture
	data    FontData
}

type FontData struct {
	Name          string
	Size          int
	Bold, Italic  bool
	Width, Height int
	Characters    map[string]Character
}

type Character struct {
	X, Y             int
	Width, Height    int
	OriginX, OriginY int
	Advance          int
}

type Text struct {
	texture       Texture
	vao, vbo, ibo uint32
	font          Font
	length        int
}

type TextVertex struct {
	pos      mgl32.Vec2
	texCoord mgl32.Vec2
}

func LoadFont(textureData []byte, fontDataJSON []byte) (Font, error) {
	fontData := FontData{}
	err := json.Unmarshal(fontDataJSON, &fontData)
	if err != nil {
		return Font{}, err
	}
	texture, err := TextureFromPNG(textureData)
	if err != nil {
		return Font{}, err
	}
	return Font{texture, fontData}, nil
}

func CreateTextRenderer() TextRenderer {
	renderer := TextRenderer{}
	vs, err := CreateVertexShader(TextVS)
	if err != nil {
		log.Fatal(err)
	}
	fs, err := CreateFragmentShader(TextFS)
	if err != nil {
		log.Fatal(err)
	}
	renderer.program, err = CreateProgramVSFS(vs, fs)
	if err != nil {
		log.Fatal(err)
	}
	return renderer
}

func CreateText(str string, font Font) Text {
	text := Text{font: font}

	vertices := make([]TextVertex, len(str)*4)
	indicies := make([]uint32, len(str)*6)

	left := float32(0)

	for i := range str {

		charData := font.data.Characters[str[i:i+1]]

		right := left + float32(charData.Width)
		bottom := float32(0.0)
		top := float32(charData.Height)

		leftTex := float32(charData.X) / float32(font.texture.width)
		rightTex := float32(charData.X+charData.Width) / float32(font.texture.width)
		topTex := float32(charData.Y) / float32(font.texture.height)
		bottomTex := float32(charData.Y+charData.Height) / float32(font.texture.height)

		copy(vertices[i*4:(i+1)*4], []TextVertex{
			{mgl32.Vec2{right, top}, mgl32.Vec2{rightTex, topTex}},
			{mgl32.Vec2{right, bottom}, mgl32.Vec2{rightTex, bottomTex}},
			{mgl32.Vec2{left, bottom}, mgl32.Vec2{leftTex, bottomTex}},
			{mgl32.Vec2{left, top}, mgl32.Vec2{leftTex, topTex}},
		})

		j := uint32(i * 4)
		copy(indicies[i*6:(i+1)*6], []uint32{
			j + 0, j + 1, j + 2,
			j + 0, j + 2, j + 3,
		})

		left += float32(charData.Advance)
	}

	text.length = len(str)

	gl.CreateVertexArrays(1, &text.vao)
	gl.BindVertexArray(text.vao)

	gl.CreateBuffers(1, &text.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, text.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(TextVertex{}))*len(vertices), unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, int32(unsafe.Sizeof(TextVertex{})), unsafe.Offsetof(TextVertex{}.pos))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, int32(unsafe.Sizeof(TextVertex{})), unsafe.Offsetof(TextVertex{}.texCoord))

	gl.CreateBuffers(1, &text.ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, text.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indicies), unsafe.Pointer(&indicies[0]), gl.STATIC_DRAW)

	return text
}

func (renderer *TextRenderer) Render(text Text, transform mgl32.Mat3) {
	gl.BindVertexArray(text.vao)
	renderer.program.Bind(map[string]Uniform{
		"textureSampler": 0,
		"transform":      transform,
	})
	text.font.texture.Bind(0)
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(text.length)*6, gl.UNSIGNED_INT, 0)
}
