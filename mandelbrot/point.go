package mandelbrot

type Point struct {
	Point      complex128
	iterations int
}

// NewPoint returns a new point at a given coordinates
func NewPoint(r, i float64) Point {
	return Point{
		Point: complex(r, i),
	}
}

// Calculate performs as many calculations as MaxIterations to determine if the point belongs to the set or not
func (m *Point) Calculate(MaxIterations int) {
	var z complex128
	point := m.Point
	iterations := m.iterations

	for real(z)*real(z)+imag(z)*imag(z) < 4 && iterations < MaxIterations {
		iterations++
		z = z*z + point
	}

	m.iterations = iterations
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
