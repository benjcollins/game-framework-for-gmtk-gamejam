package main

import (
	"engine/graphics"
	"fmt"
	"log"
	"runtime"

	_ "embed"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed test.png
var Texture []byte

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	check(glfw.Init())
	defer glfw.Terminate()

	runtime.LockOSThread()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Example Window", nil, nil)
	check(err)
	defer window.Destroy()

	window.MakeContextCurrent()

	err = gl.Init()
	check(err)

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))

	renderer := graphics.CreateRenderer()
	defer renderer.Delete()
	texture, err := graphics.CreateTexture(Texture)
	check(err)
	defer texture.Delete()

	spriteBuffer := graphics.CreateSpriteBuffer([]graphics.Sprite{
		graphics.NewSprite(mgl32.Translate2D(0.1, 0.4).Mul3(mgl32.Scale2D(0.2, 0.2))),
		graphics.NewSprite(mgl32.Translate2D(-0.3, -0.3).Mul3(mgl32.Scale2D(0.2, 0.2))),
		graphics.NewSpriteFromAtlas(mgl32.Scale2D(0.3, 0.3), 0, 0, 0.5, 0.5),
	})
	defer spriteBuffer.Delete()

	gl.ClearColor(0.5, 0.3, 0.8, 1.0)

	transform := mgl32.Ident3()

	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT)
		width, height := window.GetSize()
		renderer.Render(texture, spriteBuffer, transform, float32(width)/float32(height))
		window.SwapBuffers()
	}
}
