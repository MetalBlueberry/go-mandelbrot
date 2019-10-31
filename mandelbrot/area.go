package mandelbrot

type Area struct {
	HorizontalResolution int
	VerticalResolution   int
	TopLeft              complex128
	BottomRight          complex128
	MaxIterations        int
	Points               []Point
}

// NewAreaCentered creates a mandelbrot.Area with squared shape centered area at x,y of width = 2*area
func NewAreaCentered(Resolution, MaxIterations int, x, y, area float64) *Area {
	return &Area{
		HorizontalResolution: Resolution,
		VerticalResolution:   Resolution,
		MaxIterations:        MaxIterations,
		TopLeft:              complex(x-area, y+area),
		BottomRight:          complex(x+area, y-area),
		Points:               make([]Point, Resolution*Resolution),
	}
}

func (a *Area) Init() {
	for x := 0; x < a.HorizontalResolution; x++ {
		for y := 0; y < a.VerticalResolution; y++ {
			point := NewPoint(a.getNumber(x, y))
			a.SetPoint(x, y, point)
		}
	}
}

func (a *Area) Calculate() (progress chan int) {
	progress = make(chan int)
	go func() {
		defer close(progress)
		for i, pixel := range a.Points {
			pixel.Calculate(a.MaxIterations)
			a.Points[i] = pixel

			select {
			case progress <- i:
			default:
			}
		}
	}()
	return
}

func (a *Area) IndexFor(x, y int) int {
	return x + y*a.HorizontalResolution
}

func (a *Area) ForIndex(i int) (x, y int) {
	y = i / a.VerticalResolution
	x = i % a.HorizontalResolution
	return x, y
}

// SetPoint changes the address of the point located at x,y
func (a *Area) SetPoint(x, y int, p Point) {
	a.Points[a.IndexFor(x, y)] = p
}

// GetPoint changes the address of the point located at x,y
func (a *Area) GetPoint(x, y int) Point {
	return a.Points[a.IndexFor(x, y)]
}

// getNumber gives the real and imaginary parts for the complex number located at the given x,y coordenates in the given resolution.
func (a *Area) getNumber(x, y int) (r, i float64) {
	TopLeftReal := real(a.TopLeft)
	TopLeftImag := imag(a.TopLeft)
	BottomRightReal := real(a.BottomRight)
	BottomRightImag := imag(a.BottomRight)

	r = TopLeftReal + (float64(x)/float64(a.HorizontalResolution))*(BottomRightReal-TopLeftReal)
	i = TopLeftImag + (float64(y)/float64(a.VerticalResolution))*(BottomRightImag-TopLeftImag)
	return r, i
}