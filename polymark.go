// Package polymark is a set of convenience functions to generate
// polygon shapes that represent common geometric shapes and text.
package polymark

import (
	"errors"
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

	// Reflect reverses the Y component of the font as implemented
	// by (*Pen).Text(). The default hershey fonts have coordinate
	// Y increasing down to the lower edges of the Glyphs. This
	// attribute of the Pen causes the Y component of the font to
	// be negated, increasing Y at the upper edges of the Glyph.
	Reflect bool
}

// circle constructs an approximate circle polygon with points
// rotationally offset by theta.
func (pen *Pen) circle(s *polygon.Shapes, pt polygon.Point, r, theta float64) *polygon.Shapes {
	n := math.Floor(4 * r / pen.Scribe)
	if n < 4 {
		n = 4
	}
	n *= 4 // want a multiple of 4 for symmetry
	ang := 2 * math.Pi / n
	var pts []polygon.Point
	for i := 0.0; i < n; i++ {
		pts = append(pts, polygon.Point{
			X: pt.X + r*math.Cos(theta+i*ang),
			Y: pt.Y + r*math.Sin(theta+i*ang),
		})
	}
	return s.Builder(pts...)
}

// Circle constructs an approximate circle polygon.
func (pen *Pen) Circle(s *polygon.Shapes, pt polygon.Point, r float64) *polygon.Shapes {
	return pen.circle(s, pt, r, 0)
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
			theta := 0.0
			if len(pts) > 1 {
				theta = -math.Atan2(pts[1].Y-pt.Y, pts[1].X-pt.X)
			}
			if endCap {
				working = pen.circle(working, pt, half, theta)
			}
			last = pt
			continue
		}
		dX, dY := pt.X-last.X, pt.Y-last.Y
		theta := -math.Atan2(dY, dX)
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
			working = pen.circle(working, pt, half, theta)
		}
	}
	for _, p := range working.P {
		s = s.Builder(p.PS...)
	}
	return s
}

const twoPi = 2.0 * math.Pi

// ErrNoSolution indicates that no solution is possible.
var ErrNoSolution = errors.New("no solution")

// spiral constructs a list of points that follow a spiral path from
// from to to, around pt in rotational direction cc (true = counter
// clockwise). The number of full rotations around pt is captured in
// the winding number. If either the from or to points are equal to
// pt, no solution is viable and an error will be returned.
func spiral(width float64, from, to, pt polygon.Point, dir bool, winding uint) (pts []polygon.Point, err error) {
	var u0, u1 polygon.Point

	// Unit vector required to determine angles of start and
	// end around center of spiral, pt.
	if u0, err = pt.Unit(from); err != nil {
		return
	}
	if u1, err = pt.Unit(to); err != nil {
		return
	}
	if winding == 0 {
		if diff := from.AddX(to, -1); diff.Dot(diff) < polygon.Zeroish2 {
			err = ErrNoSolution
			return
		}
	}

	d0 := from.AddX(pt, -1)
	r0 := u0.Dot(d0)
	d1 := to.AddX(pt, -1)
	r1 := u1.Dot(d1)
	th0 := math.Atan2(u0.Y, u0.X)
	th1 := math.Atan2(u1.Y, u1.X)

	dTh := th1 - th0
	delta := float64(winding) * twoPi
	if dir {
		if dTh < 0 {
			dTh += twoPi
		}
		delta += dTh
	} else {
		if dTh > 0 {
			dTh -= twoPi
		}
		delta = dTh - delta
	}
	r := r0
	if r1 > r {
		r = r1
	}
	n := math.Floor(4 * r / width)
	if n < 4 {
		n = 4
	}
	n *= 4 // want a multiple of 4 for no less precision than circle.
	dA := delta / n
	seg := int(math.Round(n))
	step := (r1 - r0) / n
	for i := 0; i < seg; i++ {
		inc := float64(i)
		r = r0 + inc*step
		ang := th0 + dA*inc
		pts = append(pts, polygon.Point{
			pt.X + r*math.Cos(ang),
			pt.Y + r*math.Sin(ang),
		})
	}
	pts = append(pts, to)
	return
}

// Spiral returns s augmented with a width spiral polygon outline. The
// spiral has the winding number in the dir (counter-clockwise)
// direction from from, to to around a central pt.
func (pen *Pen) Spiral(s *polygon.Shapes, from, to, pt polygon.Point, width float64, dir, endCap, midCap bool, winding uint) (*polygon.Shapes, error) {
	pts, err := spiral(width, from, to, pt, dir, winding)
	if err != nil {
		return s, err
	}
	return pen.Line(s, pts, width, midCap, endCap), nil
}

// Alignment holds the horizontal and vertical alignment for rendering
// text.
type Alignment int

// AlignLeft, AlignCenter, AlignRight specify horizontal alignment.
// AlignMiddle, AlignAbove, AlignBelow specify vertical alignment.
const (
	AlignLeft   Alignment = 0
	AlignCenter Alignment = 1
	AlignRight  Alignment = 2
	AlignMiddle Alignment = 0
	AlignAbove  Alignment = 4
	AlignBelow  Alignment = 8
)

// Text renders some text as a series of polygon outlines. For scale
// >= 1.0 the enclosed polygon will have width scale*pen.Scribe, and
// the rendered font will also be scaled.  A scale of < 1.0 renders
// the characters at native size of pen.Scribe per pixel but with less
// and less of a width for the lines.
func (pen *Pen) Text(s *polygon.Shapes, x, y, scale float64, a Alignment, font *hershey.Font, text string) *polygon.Shapes {
	gl, xL, xR := font.Text(text)
	xScale := pen.Scribe * scale
	wScale := xScale * 1.8
	if scale <= 1.0 {
		xScale = pen.Scribe
	}
	yScale := xScale
	if pen.Reflect {
		yScale = -yScale
	}
	var x0, y0 float64
	trX := func(x int) float64 {
		return x0 + float64(x)*xScale
	}
	trY := func(y int) float64 {
		return y0 + float64(y)*yScale
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

	return s
}
