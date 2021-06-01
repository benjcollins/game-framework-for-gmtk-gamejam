package graphics

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	_ "embed"
)

//go:embed shaders/path.vs
var PathShaderVS string

//go:embed shaders/path.fs
var PathShaderFS string

//go:embed shaders/stroke.gs
var StrokeShaderGS string

type PathRenderer struct {
	strokeProgram Program
	fillProgram   Program
}

type Path struct {
	points     []PathVertex
	lastPoint  mgl32.Vec2
	lastNormal mgl32.Vec2
}

type PathVertex struct {
	pos    mgl32.Vec2
	normal mgl32.Vec2
}

type PathBuffer struct {
	vbo  Buffer
	vao  uint32
	size int
}

func perp(vec mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{vec.Y(), -vec.X()}
}

func (path *Path) LineTo(newPoint mgl32.Vec2) {
	newNormal := perp(newPoint.Sub(path.lastPoint)).Normalize()
	if !path.lastNormal.ApproxEqual(mgl32.Vec2{0, 0}) {
		averageNormal := newNormal.Add(path.lastNormal).Normalize()
		thing := 1.0 / averageNormal.Dot(path.lastNormal)
		path.points = append(path.points, PathVertex{path.lastPoint, averageNormal.Mul(thing)})
		path.points = append(path.points, PathVertex{path.lastPoint, averageNormal.Mul(thing)})
	} else {
		path.points = append(path.points, PathVertex{path.lastPoint, newNormal})
	}
	path.lastPoint = newPoint
	path.lastNormal = newNormal
}

func (path *Path) MoveTo(newPoint mgl32.Vec2) {
	if !path.lastNormal.ApproxEqual(mgl32.Vec2{0, 0}) {
		path.points = append(path.points, PathVertex{path.lastPoint, path.lastNormal})
	}
	path.lastPoint = newPoint
}

func (path *Path) ToBuffer() PathBuffer {
	points := path.points
	if !path.lastNormal.ApproxEqual(mgl32.Vec2{0, 0}) {
		points = append(path.points, PathVertex{path.lastPoint, path.lastNormal})
	}

	vao := uint32(0)
	gl.CreateVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vbo := CreateBuffer()
	vbo.Bind()
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*int(unsafe.Sizeof(PathVertex{})), unsafe.Pointer(&points[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, int32(unsafe.Sizeof(PathVertex{})), unsafe.Offsetof(PathVertex{}.pos))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, int32(unsafe.Sizeof(PathVertex{})), unsafe.Offsetof(PathVertex{}.normal))

	return PathBuffer{vbo, vao, len(points)}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreatePathRenderer() PathRenderer {
	renderer := PathRenderer{}
	vs, err := CreateVertexShader(PathShaderVS)
	check(err)
	fs, err := CreateFragmentShader(PathShaderFS)
	check(err)
	program, err := CreateProgramVSFS(vs, fs)
	check(err)
	renderer.fillProgram = program

	gs, err := CreateGeometryShader(StrokeShaderGS)
	check(err)
	program, err = CreateProgramVSGSFS(vs, gs, fs)
	check(err)
	renderer.strokeProgram = program

	return renderer
}

func (renderer *PathRenderer) Fill(path PathBuffer, transform mgl32.Mat3, color mgl32.Vec4) {
	gl.BindVertexArray(path.vao)
	renderer.fillProgram.Bind(map[string]Uniform{
		"transform": transform,
	})
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(path.size))
}

func (renderer *PathRenderer) Stroke(path PathBuffer, transform mgl32.Mat3, color mgl32.Vec4, width float32) {
	gl.BindVertexArray(path.vao)
	renderer.strokeProgram.Bind(map[string]Uniform{
		"transform": transform,
		"color":     color,
		"width":     width,
	})
	gl.DrawArrays(gl.LINES, 0, int32(path.size))
}
