package main

import (
	"engine/graphics"
	"log"
	"runtime"

	_ "embed"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed assets/font.json
var fontData []byte

//go:embed assets/font.png
var fontTexture []byte

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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	font, err := graphics.LoadFont(fontTexture, fontData)
	check(err)

	textRenderer := graphics.CreateTextRenderer()

	text := graphics.CreateText("Hello, World!", font)

	width, height := window.GetSize()
	aspectRatio := graphics.ComputeAspectRatio(float32(width) / float32(height))

	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		aspectRatio = graphics.ComputeAspectRatio(float32(width) / float32(height))
	})

	for !window.ShouldClose() {
		glfw.PollEvents()

		transform := aspectRatio.Mul3(mgl32.Scale2D(0.001, 0.001))

		gl.Clear(gl.COLOR_BUFFER_BIT)

		textRenderer.Render(text, transform)

		window.SwapBuffers()
	}
}
