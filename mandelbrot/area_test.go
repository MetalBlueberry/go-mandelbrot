package mandelbrot_test

import (
	"testing"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func BenchmarkArea(b *testing.B) {
	set := mandelbrot.Area{
		HorizontalResolution: 200,
		VerticalResolution:   200,
		MaxIterations:        200,
		TopLeft:              complex(-2, 2),
		BottomRight:          complex(2, -2),
	}
	set.Init()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Calculate()
	}
}
