package polymark

import (
	"encoding/json"
	"image"
	"image/color"
	"math"
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
	if got, want := len(s.P[0].PS), 32; got != want {
		t.Fatalf("large scribe circle has wrong points: got=%d, want=%d", got, want)
	}

	pen.Scribe = 0.1
	s = pen.Circle(s, polygon.Point{X: 13, Y: -11}, 2)
	if len(s.P) != 2 {
		t.Fatalf("got %d lines, wanted 1", len(s.P))
	}
	if got, want := len(s.P[1].PS), 320; got != want {
		t.Fatalf("small scribe circle has wrong points: got=%d, want=%d", got, want)
	}
}

func TestLine(t *testing.T) {
	pen := &Pen{Scribe: 1}
	s := pen.Line(nil, []polygon.Point{
		polygon.Point{0, 1},
		polygon.Point{30, 1},
	}, 7, true, true)
	s.Union()
	got := display(s)
	want := []string{
		"..##################################...",
		".##................................##..",
		"##..................................##.",
		"#....................................#.",
		"#....................................#.",
		"##..................................##.",
		".##................................##..",
		"..##################################...",
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
	s.Union()
	got := display(s)
	want := []string{
		"....######.....",
		"...#..######...",
		"..#.##....#.#..",
		".#.#.......#.#.",
		".##........#.#.",
		"##..........##.",
		"##############.",
		"##############.",
		"##.............",
		"##.............",
		"#.#............",
		".##.........##.",
		".#.#.......#.#.",
		"..#.#....##.#..",
		"...######..#...",
		".....######....",
		"...............",
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
	pt := polygon.Point{6, 5}
	// Rotate the glyph (polygon rotation is clockwise = +ve, but
	// because the fonts are +ve y down the page, when we reverse
	// y to render) we get a clockwise rotation of 45 degrees
	// with this function.
	s = s.Transform(pt, pt, math.Pi/4, .5)
	got = display(s)
	want = []string{
		"...####..",
		".#####.#.",
		".#.#..##.",
		"##..#..#.",
		"##...###.",
		"##....##.",
		"#.#......",
		".####....",
		".........",
	}
	if len(got) != len(want) {
		t.Errorf("incorrect number of lines got=%d want=%d", len(got), len(want))
	}
	for i, line := range got {
		t.Logf("[%2d]  got=%q", i, line)
		if line != want[i] {
			t.Errorf("[%2d] want=%q", i, want[i])
		}
	}
	// Reverse the Y direction of the glyph.
	pen.Reflect = true
	s = pen.Text(nil, 1, 1, .6, 0, font, "p")
	got = display(s)
	want = []string{
		"##.............",
		"##.............",
		"##.............",
		"##.............",
		"##.............",
		"##.............",
		"##.............",
		"##.######......",
		"###########....",
		"#####....###...",
		"###.......###..",
		"##.........##..",
		"##.........#.#.",
		"##..........##.",
		"##..........##.",
		"##..........##.",
		"##..........##.",
		"##.........#.#.",
		"##.........##..",
		"###.......###..",
		"####....####...",
		"###########....",
		"##..######.....",
		"...............",
	}
	if len(got) != len(want) {
		t.Errorf("incorrect number of lines got=%d want=%d", len(got), len(want))
	}
	for i, line := range got {
		t.Logf("[%2d]  got=%q", i, line)
		if i < len(want) && line != want[i] {
			t.Errorf("[%2d] want=%q", i, want[i])
		}
	}

}

func TestAlign(t *testing.T) {
	ts := []struct {
		a Alignment
		s string
	}{
		{AlignBelow + AlignLeft, "TL"},
		{AlignBelow + AlignCenter, "TC"},
		{AlignBelow + AlignRight, "TR"},
		{AlignMiddle + AlignLeft, "ML"},
		{AlignMiddle + AlignCenter, "MC"},
		{AlignMiddle + AlignRight, "MR"},
		{AlignAbove + AlignLeft, "BL"},
		{AlignAbove + AlignCenter, "BC"},
		{AlignAbove + AlignRight, "BR"},
	}
	pen := &Pen{Scribe: 1}
	font, err := hershey.New("futural")
	if err != nil {
		t.Fatalf("unable to load font: %v", err)
	}
	var s *polygon.Shapes
	for i, v := range ts {
		x := float64(60 * (i % 3))
		y := float64(15 * (i / 3))
		s = pen.Text(s, x, y, .3, v.a, font, v.s)
	}
	got := display(s)
	failed := len(got) != 76 || got[74] != ".#############........##############........#############...........#######.#................#############........##............##" || got[42] != "##....##..##....##......##...............##....##..##....##......#.#...........#..........##....##..##....##......##.........#...." || got[17] != "........##..........##..............................#...........##...........##.......................##..........##.........#...."
	for i, line := range got {
		if failed {
			t.Errorf("%3d> %q", i, line)
		} else {
			t.Logf("%3d> %q", i, line)
		}
	}
}

func TestSpiral(t *testing.T) {
	pen := Pen{
		Scribe: .1,
	}

	from := polygon.Point{0, 0}
	to := polygon.Point{2, 4}
	pt := polygon.Point{0.5, 5}

	s, err := pen.Spiral(nil, from, to, pt, 0.5, true, true, true, 0)
	if err != nil {
		t.Fatalf("failed with simple arc: %v", err)
	}
	s.Union()
	if len(s.P) != 1 {
		t.Errorf("expect one polygon, got %d", len(s.P))
		d, err := json.Marshal(*s)
		t.Errorf("s = %s (%v)", string(d), err)
	}
	ex := []string{
		"##..",
		"###.",
		"..##",
		"..##",
		"..##",
		"....",
	}
	for i, line := range display(s) {
		t.Logf("%2d: %q", i, line)
		if line != ex[i] {
			t.Errorf("  wanted: %q", ex[i])
		}
	}

	if _, err = pen.Spiral(nil, from, from, pt, 0.5, true, true, true, 0); err == nil {
		t.Error("spiral for no distance should not be possible")
	}

	s, err = pen.Spiral(nil, from, from, pt, 0.25, true, true, true, 1)
	if err != nil {
		t.Fatalf("failed with simple arc: %v", err)
	}
	s.Union()
	if len(s.P) != 2 {
		t.Errorf("expect one polygon, got %d", len(s.P))
		d, _ := json.Marshal(*s)
		t.Errorf("s = %s", string(d))
	}
	ex = []string{
		"...#####....",
		"..###..###..",
		".##.....###.",
		"##.......##.",
		"##........#.",
		"#.........#.",
		"#.........#.",
		"##.......##.",
		".##......##.",
		".###...###..",
		"..#######...",
		"............",
	}
	for i, line := range display(s) {
		t.Logf("%2d: %q", i, line)
		if line != ex[i] {
			t.Errorf("  wanted: %q", ex[i])
		}
	}
}
