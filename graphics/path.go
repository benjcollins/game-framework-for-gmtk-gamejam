package graphics

import (
	"log"
	"math"
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
		averageNormal := averageNormals(path.lastNormal, newNormal)
		path.points = append(path.points, PathVertex{path.lastPoint, averageNormal}, PathVertex{path.lastPoint, averageNormal})
	} else {
		path.points = append(path.points, PathVertex{path.lastPoint, newNormal})
	}
	path.lastPoint = newPoint
	path.lastNormal = newNormal
}

func averageNormals(n1, n2 mgl32.Vec2) mgl32.Vec2 {
	normal := n1.Add(n2).Normalize()
	factor := 1.0 / normal.Dot(n1)
	return normal.Mul(factor)
}

func normalInArc(center, point mgl32.Vec2) mgl32.Vec2 {
	return point.Sub(center).Normalize()
}

// func (path *Path) ArcTo(newPoint mgl32.Vec2, theta float64) {

// 	if !path.lastNormal.ApproxEqual(mgl32.Vec2{0, 0}) {
// 		averageNormal := averageNormals(path.lastNormal, normalInArc(center, path.lastPoint))
// 		path.points = append(path.points, PathVertex{path.lastPoint, averageNormal}, PathVertex{path.lastPoint, averageNormal})
// 	} else {
// 		path.points = append(path.points, PathVertex{path.lastPoint, normalInArc(center, path.lastPoint)})
// 	}

// 	n := 16
// 	centerToOld := path.lastPoint.Sub(center)
// 	startTheta := math.Atan(float64(centerToOld.Y()/centerToOld.X())) + math.Pi

// 	radius := center.Sub(path.lastPoint).Len()

// 	for i := 1; i < n; i++ {
// 		turn := startTheta + theta/float64(n)*float64(i)
// 		point := mgl32.Vec2{float32(math.Cos(turn)), float32(math.Sin(turn))}.Mul(radius).Add(center)
// 		normal := normalInArc(center, point)
// 		path.points = append(path.points, PathVertex{point, normal}, PathVertex{point, normal})
// 	}

// 	point := mgl32.Vec2{float32(math.Cos(startTheta + theta)), float32(math.Sin(startTheta + theta))}.Mul(radius).Add(center)
// 	path.lastPoint = point
// 	path.lastNormal = normalInArc(center, point)
// }

func (path *Path) ArcTo(newPoint mgl32.Vec2, theta float64) {

	oldToNew := newPoint.Sub(path.lastPoint)
	radius := 0.5 * oldToNew.Len() / float32(math.Sin(theta))
	center := path.lastPoint.Add(newPoint).Mul(0.5).Add(mgl32.Vec2{oldToNew.Y(), -oldToNew.X()}.Normalize().Mul(float32(math.Sqrt(float64(radius*radius - 0.25*oldToNew.Len()*oldToNew.Len())))))

	if !path.lastNormal.ApproxEqual(mgl32.Vec2{0, 0}) {
		averageNormal := averageNormals(path.lastNormal, normalInArc(center, path.lastPoint))
		path.points = append(path.points, PathVertex{path.lastPoint, averageNormal}, PathVertex{path.lastPoint, averageNormal})
	} else {
		path.points = append(path.points, PathVertex{path.lastPoint, normalInArc(center, path.lastPoint)})
	}

	n := 16
	centerToOld := path.lastPoint.Sub(center)
	startTheta := math.Atan(float64(centerToOld.Y()/centerToOld.X())) + math.Pi

	for i := 1; i < n; i++ {
		turn := startTheta - theta/float64(n)*float64(i)
		point := mgl32.Vec2{float32(math.Cos(turn)), float32(math.Sin(turn))}.Mul(radius).Add(center)
		normal := normalInArc(center, point)
		path.points = append(path.points, PathVertex{point, normal}, PathVertex{point, normal})
	}

	path.lastPoint = newPoint
	path.lastNormal = normalInArc(center, newPoint)
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

func (renderer PathRenderer) Delete() {
	renderer.fillProgram.Delete()
	renderer.strokeProgram.Delete()
}

func (buffer PathBuffer) Delete() {
	gl.DeleteVertexArrays(1, &buffer.vao)
	buffer.vbo.Delete()
}
