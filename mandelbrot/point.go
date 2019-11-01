package mandelbrot

import (
	"math/cmplx"
)

type Point struct {
	Point      complex128
	iterations int
	z          complex128
}

// NewPoint returns a new point at a given coordenates
func NewPoint(r, i float64) Point {
	return Point{
		Point: complex(r, i),
	}
}

// Calculate performs as many calculations as MaxIterations to determine if the point belongs to the set or not
func (m *Point) Calculate(MaxIterations int) {
	for !m.Diverges() && m.iterations < MaxIterations {
		m.iterations++
		m.z = cmplx.Pow(m.z, 2) + m.Point
	}
}

// Diverges returns whether the points diverges from the set.
func (m *Point) Diverges() bool {
	return real(m.z)*real(m.z)+imag(m.z)*imag(m.z) > 4
}

// Iterations returns the number of performed iterations.
func (m *Point) Iterations() int {
	return m.iterations
}

//func (m *MandelbrotPoint) Color() color.Color {
//return color.RGBA{
//R: 255 - uint8(255.0*float64(m.iterations)/float64(m.MaxIterations)),
//A: 255,
//}
//}
