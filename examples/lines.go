// Program lines draws the outline of a polygon line with --n lines.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"zappem.net/pub/graphics/hershey"
	"zappem.net/pub/graphics/polymark"
	"zappem.net/pub/graphics/raster"
	"zappem.net/pub/math/polygon"
)

var (
	n      = flag.Int("n", 6, "line count of polygon")
	m      = flag.Int("m", 7, "number of overlapping polygons")
	width  = flag.Int("width", 500, "width of image")
	height = flag.Int("height", 500, "height of image")
	dest   = flag.String("dest", "image.png", "name out output image")
	wide   = flag.Float64("wide", 10, "pixel width of lines")
	weight = flag.Float64("weight", 2, "edge line weight")
	fact   = flag.Float64("fact", 2, "how much to spread polygons")
	mid    = flag.Bool("mid-cap", true, "round mid-points of lines")
	end    = flag.Bool("end-cap", true, "round end-points of lines")
	fill   = flag.Bool("fill", false, "fill the interior of the shape")
	ids    = flag.Bool("ids", false, "show polygon sequence number")
	fn     = flag.String("font", "futural", "hershey font name")
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

	if *fill {
		var holes []int
		for i, p := range poly.P {
			if p.Hole {
				holes = append(holes, i)
			}
		}
		for i, p := range poly.P {
			if p.Hole {
				continue
			}
			lines, err := poly.Slice(i, 2, holes...)
			if err != nil {
				log.Fatalf("slice failed: %v", err)
			}
			col := color.RGBA{0xb0, 0xa0, 0xf0, 0xff}
			for _, line := range lines {
				raster.LineTo(rast, true, line.From.X, line.From.Y, line.To.X, line.To.Y, 1)
				rast.Render(im, 0, 0, col)
				rast.Reset()
			}
		}
	}

	if *ids {
		font, err := hershey.New(*fn)
		if err != nil {
			log.Fatalf("failed to load font: %v", err)
		}
		var s *polygon.Shapes
		tPen := &polymark.Pen{
			Scribe: .5,
		}
		for i, p := range poly.P {
			align := polymark.AlignMiddle | polymark.AlignCenter
			x, y := 0.5*(p.MinX+p.MaxX), 0.5*(p.MinY+p.MaxY)
			if !p.Hole {
				align = polymark.AlignBelow | polymark.AlignCenter
				y = p.MinY
			}
			s = tPen.Text(s, x, y, .15, align, font, fmt.Sprint(i))
		}
		s.Union()
		col := color.RGBA{0x0, 0x0, 0x0, 0xff}
		for _, p := range s.P {
			from := p.PS[len(p.PS)-1]
			for _, to := range p.PS {
				raster.LineTo(rast, true, from.X, from.Y, to.X, to.Y, *weight)
				from = to
			}
			rast.Render(im, 0, 0, col)
			rast.Reset()
		}
	}

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
