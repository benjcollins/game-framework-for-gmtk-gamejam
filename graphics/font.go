package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func GenerateCharacterPath(f *truetype.Font, r rune) Path {
	i := f.Index(r)

	buffer := truetype.GlyphBuf{}

	buffer.Load(f, 50, i, font.HintingNone)

	path := Path{}

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
