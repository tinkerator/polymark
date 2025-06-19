// Program lines draws the outline of a polygon line with --n lines.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"zappem.net/pub/graphics/polymark"
	"zappem.net/pub/graphics/raster"
	"zappem.net/pub/math/polygon"
)

var (
	n      = flag.Int("n", 7, "line count of polygon")
	m      = flag.Int("m", 1, "number of overlapping polygons")
	width  = flag.Int("width", 500, "width of image")
	height = flag.Int("height", 500, "height of image")
	dest   = flag.String("dest", "image.png", "name out output image")
	wide   = flag.Float64("wide", 10, "pixel width of lines")
	weight = flag.Float64("weight", 2, "edge line weight")
	fact   = flag.Float64("fact", 2, "how much to spread polygons")
	mid    = flag.Bool("mid-cap", true, "round mid-points of lines")
	end    = flag.Bool("end-cap", true, "round end-points of lines")
)

func main() {
	flag.Parse()

	if *n <= 1 || *m <= 0 {
		log.Fatalf("both --n=%d and --m=%d must be >1 and positive respectively", n, m)
	}

	w := float64(*width)
	h := float64(*height)
	r := w
	if r > h {
		r = h
	}
	r = r / (2.1 * float64(*m))
	o := r * *fact * float64(*m-1) / float64(*m)

	pen := &polymark.Pen{
		Scribe: *wide,
	}

	angO := 2 * math.Pi / float64(*m)
	angR := 2 * math.Pi / float64(*n)
	var poly *polygon.Shapes
	for i := 0; i < *m; i++ {
		ang := float64(i) * angO
		x0, y0 := 0.5*w+o*math.Cos(ang), 0.5*h+o*math.Sin(ang)
		var pts []polygon.Point
		for j := 0; j <= *n; j++ {
			theta := ang + float64(j)*angR
			x, y := x0+r*math.Cos(theta), y0+r*math.Sin(theta)
			pts = append(pts, polygon.Point{x, y})
		}
		poly = pen.Line(poly, pts, *wide, *mid, *end)
	}
	poly.Union()

	im := image.NewRGBA(image.Rect(0, 0, *width, *height))
	draw.Draw(im, im.Bounds(), &image.Uniform{color.RGBA{0xff, 0xff, 0xff, 0xff}}, image.ZP, draw.Src)

	rast := raster.NewRasterizer()
	for _, p := range poly.P {
		col := color.RGBA{0xff, 0, 0, 0xff}
		if p.Hole {
			col = color.RGBA{0, 0, 0xff, 0xff}
		}
		from := p.PS[len(p.PS)-1]
		for _, to := range p.PS {
			raster.LineTo(rast, true, from.X, from.Y, to.X, to.Y, *weight)
			from = to
		}
		rast.Render(im, 0, 0, col)
		rast.Reset()
	}

	f, err := os.Create(*dest)
	if err != nil {
		log.Fatalf("failed to create %q: %v", *dest, err)
	}
	defer f.Close()
	png.Encode(f, im)
	log.Printf("wrote result to %q", *dest)
}
