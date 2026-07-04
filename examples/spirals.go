// Program spirals generates a selection of spiral polygon shapes.
// The program can output in polygon.Shapes json format, or as a PNG.
// To convert to an SVG, see the examples/outline.go tool of the
package main

import (
	"encoding/json"
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
	width = flag.Float64("width", 5, "point width of lines")
	poly  = flag.String("poly", "", "named json output file")
	img   = flag.String("png", "", "named png output file")
)

func main() {
	flag.Parse()

	var ps *polygon.Shapes
	dx, dy := 100., 100.
	for i := 0; i < 5; i++ {
		cx := float64(i*2+1) * dx / 2
		pen := polymark.Pen{
			Scribe: float64(5 - i),
		}
		a := math.Pi * float64(i) / 4
		c, s := math.Cos(a), math.Sin(a)
		for j := 0; j < 5; j++ {
			cy := float64(j*2+1) * dy / 2
			pt := polygon.Point{cx, cy}
			from := polygon.Point{cx + .4*c*dy, cy + .4*s*dx}
			to := polygon.Point{cx, cy - .4*dx/float64(5-j)}
			dir := i&1 == 0
			winding := uint(4 - j)

			working, err := pen.Spiral(nil, from, to, pt, *width, dir, true, true, winding)
			if err != nil {
				log.Fatalf("[%d,%d] spiral generation failed: %v", i, j, err)
			}
			working.Union()
			ps = ps.Include(working.P...)
		}
	}

	if *poly != "" {
		out, err := os.Create(*poly)
		if err != nil {
			log.Fatalf("Unable to generate polygon.Shapes output %q: %v", *poly, err)
		}
		defer out.Close()
		enc := json.NewEncoder(out)
		if err := enc.Encode(ps); err != nil {
			log.Fatalf("failed to json encode %q: %v", *poly, err)
		}
	} else if *img != "" {
		im := image.NewRGBA(image.Rect(0, 0, 500, 500))
		draw.Draw(im, im.Bounds(), &image.Uniform{color.RGBA{0xff, 0xff, 0xff, 0xff}}, image.ZP, draw.Src)

		rast := raster.NewRasterizer()
		for _, p := range ps.P {
			col := color.RGBA{0xff, 0, 0, 0xff}
			if p.Hole {
				col = color.RGBA{0, 0, 0xff, 0xff}
			}
			from := p.PS[len(p.PS)-1]
			for _, to := range p.PS {
				raster.LineTo(rast, true, from.X, 500-from.Y, to.X, 500-to.Y, 1)
				from = to
			}
			rast.Render(im, 0, 0, col)
			rast.Reset()
		}

		f, err := os.Create(*img)
		if err != nil {
			log.Fatalf("failed to create %q: %v", *img, err)
		}
		defer f.Close()
		png.Encode(f, im)
		return
	} else {
		log.Fatal("supply --poly=<name> or --png=<name> for output")
	}
}
