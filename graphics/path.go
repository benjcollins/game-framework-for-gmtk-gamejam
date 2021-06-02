package graphics

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	_ "embed"
)

//go:embed shaders/path.vs
var PathVS string

//go:embed shaders/stroke.fs
var StrokeFS string

//go:embed shaders/fill.fs
var FillFS string

//go:embed shaders/stroke.gs
var StrokeGS string

//go:embed shaders/fill.gs
var FillGS string

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
		averageNormal := averageNormals(path.lastNormal, newNormal)
		path.points = append(path.points, PathVertex{path.lastPoint, averageNormal}, PathVertex{path.lastPoint, averageNormal})
	} else {
		path.points = append(path.points, PathVertex{path.lastPoint, newNormal})
	}
	path.lastPoint = newPoint
	path.lastNormal = newNormal
}

func normalInQuadratic(start, control, end mgl32.Vec2, t float32) mgl32.Vec2 {
	tangent := control.Sub(start).Mul(2 * (1 - t)).Add(end.Sub(control).Mul(2 * t))
	return perp(tangent).Normalize()
}

func pointInQuadratic(start, control, end mgl32.Vec2, t float32) mgl32.Vec2 {
	return start.Mul((1 - t) * (1 - t)).Add(control.Mul(2 * (1 - t) * t)).Add(end.Mul(t * t))
}

func (path *Path) QuadraticTo(controlPoint, newPoint mgl32.Vec2) {

	normal := normalInQuadratic(path.lastPoint, controlPoint, newPoint, 0.0)
	if !path.lastNormal.ApproxEqual(mgl32.Vec2{0, 0}) {
		averageNormal := averageNormals(path.lastNormal, normal)
		path.points = append(path.points, PathVertex{path.lastPoint, averageNormal}, PathVertex{path.lastPoint, averageNormal})
	} else {
		path.points = append(path.points, PathVertex{path.lastPoint, normal})
	}

	n := 4

	for i := 1; i < n; i++ {
		point := pointInQuadratic(path.lastPoint, controlPoint, newPoint, float32(i)/float32(n))
		normal := normalInQuadratic(path.lastPoint, controlPoint, newPoint, float32(i)/float32(n))
		path.points = append(path.points, PathVertex{point, normal}, PathVertex{point, normal})
	}

	path.lastPoint = newPoint
	path.lastNormal = normalInQuadratic(path.lastPoint, controlPoint, newPoint, 1.0)
}

func averageNormals(n1, n2 mgl32.Vec2) mgl32.Vec2 {
	normal := n1.Add(n2).Normalize()
	factor := 1.0 / normal.Dot(n1)
	return normal.Mul(factor)
}

func normalInArc(center, point mgl32.Vec2) mgl32.Vec2 {
	return point.Sub(center).Normalize()
}

func (path *Path) MoveTo(newPoint mgl32.Vec2) {
	if !path.lastNormal.ApproxEqual(mgl32.Vec2{0, 0}) {
		path.points = append(path.points, PathVertex{path.lastPoint, path.lastNormal})
	}
	path.lastPoint = newPoint
	path.lastNormal = mgl32.Vec2{0, 0}
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
	vs, err := CreateVertexShader(PathVS)
	check(err)
	fs, err := CreateFragmentShader(FillFS)
	check(err)
	gs, err := CreateGeometryShader(FillGS)
	check(err)
	program, err := CreateProgramVSGSFS(vs, gs, fs)
	check(err)
	renderer.fillProgram = program

	gs, err = CreateGeometryShader(StrokeGS)
	check(err)
	fs, err = CreateFragmentShader(StrokeFS)
	check(err)
	program, err = CreateProgramVSGSFS(vs, gs, fs)
	check(err)
	renderer.strokeProgram = program

	return renderer
}

func (renderer *PathRenderer) Fill(path PathBuffer, transform mgl32.Mat3, color mgl32.Vec4) {
	gl.Enable(gl.STENCIL_TEST)
	gl.Clear(gl.STENCIL_BUFFER_BIT)
	gl.StencilOp(gl.INVERT, gl.INVERT, gl.INVERT)
	gl.StencilFunc(gl.ALWAYS, 0, 1)
	gl.ColorMask(false, false, false, false)

	gl.BindVertexArray(path.vao)
	renderer.fillProgram.Bind(map[string]Uniform{
		"transform": transform,
		"color":     color,
		"origin":    mgl32.Vec2{0, 0},
	})
	gl.DrawArrays(gl.LINES, 0, int32(path.size))

	gl.ColorMask(true, true, true, true)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.KEEP)
	gl.StencilFunc(gl.EQUAL, 0xFF, 0xFF)
	gl.DrawArrays(gl.LINES, 0, int32(path.size))

	gl.Disable(gl.STENCIL_TEST)
}

func (renderer *PathRenderer) Stroke(path PathBuffer, transform mgl32.Mat3, color mgl32.Vec4, width float32) {
	gl.BindVertexArray(path.vao)
	renderer.strokeProgram.Bind(map[string]Uniform{
		"transform": transform,
		"color":     color,
		"width":     width,
		"sides":     0,
		"threshold": float32(0.004),
	})
	gl.DrawArrays(gl.LINES, 0, int32(path.size))
}

func (renderer PathRenderer) Delete() {
	renderer.fillProgram.Delete()
	renderer.strokeProgram.Delete()
}

func (buffer PathBuffer) Delete() {
	gl.DeleteVertexArrays(1, &buffer.vao)
	buffer.vbo.Delete()
}
