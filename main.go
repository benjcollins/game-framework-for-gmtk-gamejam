package main

import (
	"engine/graphics"
	"fmt"
	"log"
	"math/rand"
	"runtime"

	_ "embed"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed campfire.png
var Texture []byte

//go:embed particle.png
var FireTexture []byte

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
	// glfw.WindowHint(glfw.Samples, 16)

	window, err := glfw.CreateWindow(800, 600, "Example Window", nil, nil)
	check(err)
	defer window.Destroy()

	window.MakeContextCurrent()

	err = gl.Init()
	check(err)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))

	spriteRenderer := graphics.CreateRenderer()
	defer spriteRenderer.Delete()
	texture, err := graphics.CreateTexture(Texture)
	check(err)
	defer texture.Delete()

	sprites := graphics.CreateSpriteBuffer([]graphics.Sprite{
		graphics.NewSprite(mgl32.Scale2D(0.5, 0.5)),
	}, texture)
	defer sprites.Delete()

	particleTexture, err := graphics.CreateTexture(FireTexture)
	check(err)
	defer particleTexture.Delete()
	particleRenderer := graphics.CreateParticleRenderer()
	particleSystem := particleRenderer.CreateParticleSystem(1000, particleTexture, 4, 4, func(p *graphics.Particle) bool {
		p.Frame += 0.0005
		p.Transform = mgl32.Translate2D(0, 0.00005).Mul3(p.Transform)
		return p.Frame < 6
	})

	pathRenderer := graphics.CreatePathRenderer()
	path := graphics.Path{}
	path.MoveTo(mgl32.Vec2{0, 0})
	path.LineTo(mgl32.Vec2{0.2, 0})
	// path.LineTo(mgl32.Vec2{0.5, 0.5})
	// path.LineTo(mgl32.Vec2{0.6, 0.3})
	// path.LineTo(mgl32.Vec2{0.8, 0.6})
	pathBuffer := path.ToBuffer()

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	width, height := window.GetSize()
	aspectRatio := graphics.ComputeAspectRatio(float32(width) / float32(height))

	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		aspectRatio = graphics.ComputeAspectRatio(float32(width) / float32(height))
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action != glfw.Press {
			return
		}
		width, height := window.GetSize()
		x, y := window.GetCursorPos()
		point := mgl32.Vec3{float32(x/float64(width)*4.0 - 2.0), float32(y/-float64(height)*4.0 + 2.0), 1.0}
		path.LineTo(aspectRatio.Inv().Mul3x1(point).Vec2())
		pathBuffer = path.ToBuffer()
	})

	for !window.ShouldClose() {
		glfw.PollEvents()

		if rand.Float32() > 0.99 {
			transform := mgl32.Scale2D(0.1, 0.1)
			transform = mgl32.Translate2D(0.2*rand.Float32()-0.1, 0.1+0.2*rand.Float32()-0.1).Mul3(transform)
			particleSystem.AppendParticle(graphics.NewParticle(transform, 0.0))
		}

		particleSystem.Update()

		transform := aspectRatio

		gl.Clear(gl.COLOR_BUFFER_BIT)
		spriteRenderer.Render(sprites, transform)
		particleRenderer.Render(particleSystem, transform)
		pathRenderer.Stroke(pathBuffer, transform, mgl32.Vec4{0, 0, 0, 1}, 0.02)
		window.SwapBuffers()
	}
}
