package mandelbrot_test

import (
	"testing"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func BenchmarkArea(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		set := mandelbrot.Area{
			HorizontalResolution: 200,
			VerticalResolution:   200,
			MaxIterations:        200,
			TopLeft:              complex(-2, 2),
			BottomRight:          complex(2, -2),
		}
		set.Init()
		b.StartTimer()
		set.Calculate()
	}
}

func BenchmarkComplexArea(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		set := mandelbrot.Area{
			HorizontalResolution: 1060,
			VerticalResolution:   730,
			MaxIterations:        3534,
			TopLeft:              complex(-1.401854499759, -0.000743603637),
			BottomRight:          complex(-1.399689899172, 0.000743603637),
		}
		set.Init()
		b.StartTimer()
		set.Calculate()
	}
}
