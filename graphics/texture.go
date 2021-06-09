package graphics

import (
	"bytes"
	"fmt"
	"image"
	"unsafe"

	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	textureID     uint32
	width, height int
}

func TextureFromRGBA(rgba *image.NRGBA) Texture {
	width := int32(rgba.Rect.Dx())
	height := int32(rgba.Rect.Dy())

	textureID := uint32(0)
	gl.CreateTextures(gl.TEXTURE_2D, 1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	borderColor := [4]float32{0.0, 0.0, 0.0, 0.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&rgba.Pix[0]))

	return Texture{textureID, int(width), int(height)}
}

func TextureFromPNG(data []byte) (Texture, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return Texture{}, err
	}

	rgba, ok := img.(*image.NRGBA)
	if !ok {
		return Texture{}, fmt.Errorf("image must be a png image")
	}

	return TextureFromRGBA(rgba), nil
}

func (texture Texture) Bind(n uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + n)
	gl.BindTexture(gl.TEXTURE_2D, texture.textureID)
}

func (texture Texture) Delete() {
	gl.DeleteTextures(1, &texture.textureID)
}
