package graphics

import "github.com/go-gl/gl/v4.1-core/gl"

type Buffer struct {
	ID uint32
}

func CreateBuffer() Buffer {
	buffer := Buffer{}
	gl.CreateBuffers(1, &buffer.ID)
	return buffer
}

func (buffer Buffer) Delete() {
	gl.DeleteBuffers(1, &buffer.ID)
}
