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
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed campfire.png
var Texture []byte

//go:embed particle.png
var FireTexture []byte

//go:embed font.ttf
var fontData []byte

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
	gl.BlendEquation(gl.FUNC_ADD)

	f, err := truetype.Parse(fontData)
	check(err)

	pathRenderer := graphics.CreatePathRenderer()
	defer pathRenderer.Delete()

	path := GenerateCharacterPath(f, 'O')
	pathBuffer := path.ToBuffer()
	defer pathBuffer.Delete()

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	width, height := window.GetSize()
	aspectRatio := graphics.ComputeAspectRatio(float32(width) / float32(height))

	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		aspectRatio = graphics.ComputeAspectRatio(float32(width) / float32(height))
	})

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press {
			return
		}
		path := GenerateCharacterPath(f, rune(key))
		pathBuffer.Delete()
		pathBuffer = path.ToBuffer()
	})

	for !window.ShouldClose() {
		glfw.PollEvents()

		transform := aspectRatio

		gl.Clear(gl.COLOR_BUFFER_BIT)
		pathRenderer.Fill(pathBuffer, transform, mgl32.Vec4{0, 0, 0, 1})
		pathRenderer.Stroke(pathBuffer, transform, mgl32.Vec4{0, 0, 0, 1}, 0.005)
		window.SwapBuffers()
	}
}

func GenerateCharacterPath(f *truetype.Font, r rune) graphics.Path {
	i := f.Index(r)

	buffer := truetype.GlyphBuf{}

	buffer.Load(f, 50, i, font.HintingNone)

	path := graphics.Path{}

	getPoint := func(i int) mgl32.Vec2 {
		return mgl32.Vec2{float32(buffer.Points[i].X) / 100.0, float32(buffer.Points[i].Y) / 100.0}
	}

	path.MoveTo(mgl32.Vec2{0, 0})

	start := 0
	for _, end := range buffer.Ends {
		path.MoveTo(getPoint(start))
		finished := false
		for i := start + 1; i < end; i++ {
			if buffer.Points[i].Flags&1 == 0 {
				fmt.Println("control")
				if i+1 >= len(buffer.Points) {
					path.QuadraticTo(getPoint(i), getPoint(start))
					finished = true
				} else if buffer.Points[i+1].Flags&1 == 0 {
					path.QuadraticTo(getPoint(i), getPoint(i).Add(getPoint(i+1)).Mul(0.5))
				} else {
					path.QuadraticTo(getPoint(i), getPoint(i+1))
					i++
				}
			} else {
				fmt.Println("other")
				path.LineTo(getPoint(i))
			}
		}
		if !finished {
			path.LineTo(getPoint(start))
		}
		start = end
	}

	return path
}
