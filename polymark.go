// Package polymark is a set of convenience functions to generate
// polygon shapes that represent common geometric shapes and text.
package polymark

import (
	"math"

	"zappem.net/pub/graphics/hershey"
	"zappem.net/pub/math/polygon"
)

// Pen holds the drawing implement.
type Pen struct {
	// Scribe is the width of the smallest achievable mark.  This
	// is essentially the minimum width of any outlines
	// represented by polygons.
	Scribe float64
}

// Circle constructs an approximate circle polygon.
func (pen *Pen) Circle(s *polygon.Shapes, pt polygon.Point, r float64) *polygon.Shapes {
	n := math.Floor(4 * r / pen.Scribe)
	if n < 8 {
		n = 8
	}
	ang := 2 * math.Pi / n
	var pts []polygon.Point
	for i := 0.0; i < n; i++ {
		pts = append(pts, polygon.Point{
			X: pt.X + r*math.Cos(i*ang),
			Y: pt.Y + r*math.Sin(i*ang),
		})
	}
	return s.Builder(pts...)
}

// Line constructs the outline of a series of straight line segments
// of a specified width. The corners of the line are rounded if midCap
// is true, and the endCap value determines if the ends of the line
// are rounded.
func (pen *Pen) Line(s *polygon.Shapes, pts []polygon.Point, width float64, midCap, endCap bool) *polygon.Shapes {
	var last polygon.Point
	var working *polygon.Shapes
	half := width / 2
	for i, pt := range pts {
		if i == 0 {
			if endCap {
				working = pen.Circle(working, pt, half)
			}
			last = pt
			continue
		}
		dX, dY := pt.X-last.X, pt.Y-last.Y
		r := math.Sqrt(dX*dX + dY*dY)
		dX, dY = half*dX/r, half*dY/r
		working = working.Builder(
			polygon.Point{X: last.X + dY, Y: last.Y - dX},
			polygon.Point{X: pt.X + dY, Y: pt.Y - dX},
			polygon.Point{X: pt.X - dY, Y: pt.Y + dX},
			polygon.Point{X: last.X - dY, Y: last.Y + dX},
		)
		last = pt
		if final := i == len(pts)-1; (midCap && !final) || (endCap && final) {
			working = pen.Circle(working, pt, half)
		}
	}
	working.Union()
	for _, p := range working.P {
		s = s.Builder(p.PS...)
	}
	s.Union()
	return s
}

// Alignment holds the horizontal and vertical alignment for rendering
// text.
type Alignment int

// AlignLeft, AlignCenter, AlignRight specify horizontal alignment.
// AlignMiddle, AlignAbove, AlignBelow specify vertical alignment.
const (
	AlignLeft   Alignment = 0
	AlignCenter           = 1
	AlignRight            = 2
	AlignMiddle           = 0
	AlignAbove            = 4
	AlignBelow            = 8
)

const ()

// Text renders some text as a series of polygon outlines. For scale
// >= 1.0 the enclosed polygon will have width scale*pen.Scribe, and
// the rendered font will also be scaled.  A scale of < 1.0 renders
// the characters at native size of pen.Scribe per pixel but with less
// and less of a width for the lines.
func (pen *Pen) Text(s *polygon.Shapes, x, y, scale float64, a Alignment, font *hershey.Font, text string) *polygon.Shapes {
	gl, xL, xR := font.Text(text)
	wScale := pen.Scribe * scale
	fScale := wScale
	if scale < 1.0 {
		fScale = pen.Scribe
	}

	var x0, y0 float64
	trX := func(x int) float64 {
		return x0 + float64(x)*fScale
	}
	trY := func(y int) float64 {
		return y0 + float64(y)*fScale
	}

	switch a & 3 {
	case AlignLeft:
		x0 = x
	case AlignCenter:
		x0 = x - trX(xR-xL)/2
	case AlignRight:
		x0 = x - trX(xR)
	}
	switch a & ^3 {
	case AlignAbove:
		y0 = y - trY(gl.Top)
	case AlignMiddle:
		y0 = y
	case AlignBelow:
		y0 = y - trY(gl.Bottom)
	}

	for _, line := range gl.Strokes {
		if len(line) == 0 {
			continue
		}
		var pts []polygon.Point
		for _, pt := range line {
			to := polygon.Point{
				X: trX(pt[0]),
				Y: trY(pt[1]),
			}
			pts = append(pts, to)
		}
		s = pen.Line(s, pts, wScale, true, true)
	}
	s.Union()

	return s
}
