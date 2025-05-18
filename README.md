# polymark - render shapes and text as polygons

## Overview

For line graphics and plotters/lasers/cnc devices, it is convenient to
be able to render shapes as outlines. This package converts desired
shapes into such outlines in the form of
[`polygon.Shapes`](https://zappem.net/pub/math/polygon).

For now, the package contains a simple set of tests that can be run as
follows:

```
$ git clone https://github.com/tinkerator/polymark.git
$ cd polymark
$ go test -v
=== RUN   TestCircle
--- PASS: TestCircle (0.00s)
=== RUN   TestLine
    polymark_test.go:119: [ 0]  got="..#################################...."
    polymark_test.go:119: [ 1]  got=".#.................................##.."
    polymark_test.go:119: [ 2]  got="#....................................#."
    polymark_test.go:119: [ 3]  got="#....................................#."
    polymark_test.go:119: [ 4]  got="#....................................#."
    polymark_test.go:119: [ 5]  got="#....................................#."
    polymark_test.go:119: [ 6]  got=".##.................................#.."
    polymark_test.go:119: [ 7]  got="...#################################..."
    polymark_test.go:119: [ 8]  got="......................................."
--- PASS: TestLine (0.00s)
=== RUN   TestText
    polymark_test.go:156: [ 0]  got="....#####....."
    polymark_test.go:156: [ 1]  got="...#.#######.."
    polymark_test.go:156: [ 2]  got="..###.....###."
    polymark_test.go:156: [ 3]  got=".##........#.#"
    polymark_test.go:156: [ 4]  got=".##.........##"
    polymark_test.go:156: [ 5]  got="##..........##"
    polymark_test.go:156: [ 6]  got="##############"
    polymark_test.go:156: [ 7]  got="##############"
    polymark_test.go:156: [ 8]  got="##............"
    polymark_test.go:156: [ 9]  got="##............"
    polymark_test.go:156: [10]  got=".##..........."
    polymark_test.go:156: [11]  got=".##.........##"
    polymark_test.go:156: [12]  got=".#.#.......#.#"
    polymark_test.go:156: [13]  got="..##.....##.#."
    polymark_test.go:156: [14]  got="...######.##.."
    polymark_test.go:156: [15]  got=".....#####...."
--- PASS: TestText (0.00s)
PASS
ok      zappem.net/pub/graphics/polymark        0.003s
```

## Reporting bugs

The `polymark` package has been developed purely out of self-interest
and offers no guarantee of fixes/support. That being said, if you
would like to suggest a feature addition or suggest a fix, please use
the [bug tracker](https://github.com/tinkerator/polymark/issues).

## License information

See the [LICENSE](LICENSE) file: the same BSD 3-clause license as that
used by [golang](https://golang.org/LICENSE) itself.
