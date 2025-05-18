package polymark

import (
	"image"
	"image/color"
	"testing"

	"zappem.net/pub/graphics/hershey"
	"zappem.net/pub/math/polygon"
)

func line(im *image.Gray16, from, to polygon.Point) {
	x0, y0, x1, y1 := int(from.X+.5), int(from.Y+.5), int(to.X+.5), int(to.Y+.5)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	dx, dy := sx*(x1-x0), -sy*(y1-y0)
	er := dx + dy
	x, y := x0, y0
	for {
		im.Set(x, y, color.Gray{1})
		e2 := 2 * er
		if e2 >= dy {
			if x == x1 {
				break
			}
			er += dy
			x += sx
		}
		if e2 <= dx {
			if y == y1 {
				break
			}
			er += dx
			y += sy
		}
	}
}

// ASCII display of polygon.Shapes.
func display(s *polygon.Shapes) []string {
	ll, tr := s.BB()
	dx := (tr.X - ll.X) + 2
	dy := (tr.Y - ll.Y) + 2
	width := int(dx)
	height := int(dy)
	im := image.NewGray16(image.Rect(0, 0, width, height))
	for _, p := range s.P {
		var old polygon.Point
		for i, pt := range p.PS {
			if i == 0 {
				old = p.PS[len(p.PS)-1]
			}
			line(im, old.AddX(ll, -1), pt.AddX(ll, -1))
			old = pt
		}
	}
	var lines []string
	for i := 0; i < height; i++ {
		cs := make([]byte, width)
		for j := 0; j < width; j++ {
			c := '.'
			if im.Gray16At(j, i).Y != 0 {
				c = '#'
			}
			cs[j] = byte(c)
		}
		lines = append(lines, string(cs))
	}
	return lines
}

func TestCircle(t *testing.T) {
	pen := &Pen{Scribe: 1.0}
	s := pen.Circle(nil, polygon.Point{X: 13, Y: 11}, 2)
	if len(s.P) != 1 {
		t.Fatalf("got %d lines, wanted 1", len(s.P))
	}
	if got, want := len(s.P[0].PS), 8; got != want {
		t.Fatalf("large scribe circle has wrong points: got=%d, want=%d", got, want)
	}

	pen.Scribe = 0.1
	s = pen.Circle(s, polygon.Point{X: 13, Y: -11}, 2)
	if len(s.P) != 2 {
		t.Fatalf("got %d lines, wanted 1", len(s.P))
	}
	if got, want := len(s.P[1].PS), 80; got != want {
		t.Fatalf("small scribe circle has wrong points: got=%d, want=%d", got, want)
	}
}

func TestLine(t *testing.T) {
	pen := &Pen{Scribe: 1}
	s := pen.Line(nil, []polygon.Point{
		polygon.Point{0, 1},
		polygon.Point{30, 1},
	}, 7, true, true)
	got := display(s)
	want := []string{
		"..#################################....",
		".#.................................##..",
		"#....................................#.",
		"#....................................#.",
		"#....................................#.",
		"#....................................#.",
		".##.................................#..",
		"...#################################...",
		".......................................",
	}
	if len(got) != len(want) {
		t.Fatalf("incorrect number of lines got=%d want=%d", len(got), len(want))
	}
	for i, line := range got {
		t.Logf("[%2d]  got=%q", i, line)
		if line != want[i] {
			t.Errorf("[%2d] want=%q", i, want[i])
		}
	}
}

func TestText(t *testing.T) {
	pen := &Pen{Scribe: 1}
	font, err := hershey.New("rowmans")
	if err != nil {
		t.Fatalf("unable to load font: %v", err)
	}
	s := pen.Text(nil, 1, 1, .6, 0, font, "e")
	got := display(s)
	want := []string{
		"....#####.....",
		"...#.#######..",
		"..###.....###.",
		".##........#.#",
		".##.........##",
		"##..........##",
		"##############",
		"##############",
		"##............",
		"##............",
		".##...........",
		".##.........##",
		".#.#.......#.#",
		"..##.....##.#.",
		"...######.##..",
		".....#####....",
	}
	if len(got) != len(want) {
		t.Fatalf("incorrect number of lines got=%d want=%d", len(got), len(want))
	}
	for i, line := range got {
		t.Logf("[%2d]  got=%q", i, line)
		if line != want[i] {
			t.Errorf("[%2d] want=%q", i, want[i])
		}
	}
}
