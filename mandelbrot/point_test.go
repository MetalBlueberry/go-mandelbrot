package mandelbrot_test

import (
	"testing"

	. "github.com/metalblueberry/mandelbrot/mandelbrot"
)

type testMandelbrotPointCases struct {
	point      *Point
	iterations int
	diverges   bool
}

func TestMandelbrotPoint(t *testing.T) {
	tests := []testMandelbrotPointCases{
		testMandelbrotPointCases{
			point:      NewPoint(1, 0),
			iterations: 100,
			diverges:   true,
		},
		testMandelbrotPointCases{
			point:      NewPoint(-1, 0),
			iterations: 100,
			diverges:   false,
		},
		testMandelbrotPointCases{
			point:      NewPoint(-0.5, 0.5),
			iterations: 100,
			diverges:   false,
		},
		testMandelbrotPointCases{
			point:      NewPoint(0, 1),
			iterations: 100,
			diverges:   true,
		},
	}

	for i, test := range tests {
		t.Log(test)
		test.point.Calculate(test.iterations)
		if test.point.Diverges() != test.diverges {
			t.Errorf("Test %d failed, Point %f diverges %t expected %t", i, test.point.Point, test.point.Diverges(), test.diverges)
		}
	}
}
