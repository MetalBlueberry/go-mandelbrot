package mandelbrot_test

import (
	"context"
	"testing"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func benchmarkComplexPictureWorkers(b *testing.B, workers int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		pic := &mandelbrot.Picture{
			TopLeft:               complex(-1.401854499759, -0.000743603637),
			MaxIterations:         3534,
			ChunkSize:             0.00021646,
			HorizontalImageChunks: 106,
			VerticalImageChunks:   73,
			ChunkImageSize:        10,
		}
		ctx := context.Background()
		done := make(chan int)
		pic.Init()
		b.StartTimer()

		go pic.Calculate(ctx, done, workers)
		for range done {

		}
	}
}

func BenchmarkComplexPictureWorkers1(b *testing.B)  { benchmarkComplexPictureWorkers(b, 1) }
func BenchmarkComplexPictureWorkers2(b *testing.B)  { benchmarkComplexPictureWorkers(b, 2) }
func BenchmarkComplexPictureWorkers3(b *testing.B)  { benchmarkComplexPictureWorkers(b, 3) }
func BenchmarkComplexPictureWorkers4(b *testing.B)  { benchmarkComplexPictureWorkers(b, 4) }
func BenchmarkComplexPictureWorkers5(b *testing.B)  { benchmarkComplexPictureWorkers(b, 5) }
func BenchmarkComplexPictureWorkers6(b *testing.B)  { benchmarkComplexPictureWorkers(b, 6) }
func BenchmarkComplexPictureWorkers7(b *testing.B)  { benchmarkComplexPictureWorkers(b, 7) }
func BenchmarkComplexPictureWorkers8(b *testing.B)  { benchmarkComplexPictureWorkers(b, 8) }
func BenchmarkComplexPictureWorkers9(b *testing.B)  { benchmarkComplexPictureWorkers(b, 9) }
func BenchmarkComplexPictureWorkers10(b *testing.B) { benchmarkComplexPictureWorkers(b, 10) }
func BenchmarkComplexPictureWorkers11(b *testing.B) { benchmarkComplexPictureWorkers(b, 11) }
func BenchmarkComplexPictureWorkers12(b *testing.B) { benchmarkComplexPictureWorkers(b, 12) }

func benchmarkComplexPictureChunks(b *testing.B, HorizontalImageChunks, VerticalImageChunks, ChunkImageSize, workers int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		pic := &mandelbrot.Picture{
			TopLeft:               complex(-1.401854499759, -0.000743603637),
			MaxIterations:         3534,
			ChunkSize:             0.00021646 * float64(ChunkImageSize),
			HorizontalImageChunks: HorizontalImageChunks,
			VerticalImageChunks:   VerticalImageChunks,
			ChunkImageSize:        ChunkImageSize,
		}
		ctx := context.Background()
		done := make(chan int)
		pic.Init()
		b.StartTimer()

		go pic.Calculate(ctx, done, workers)
		for range done {

		}
	}
}

func BenchmarkComplexPictureChunks1024x1024x1w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 1)
}
func BenchmarkComplexPictureChunks352x352x3w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 352, 352, 3, 1)
}
func BenchmarkComplexPictureChunks96x96x11w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 96, 96, 11, 1)
}
func BenchmarkComplexPictureChunks32x32x33w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 33, 1)
}
func BenchmarkComplexPictureChunks16x16x66w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 66, 1)
}
func BenchmarkComplexPictureChunks8x8x132w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 132, 1)
}
func BenchmarkComplexPictureChunks4x4x264w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 264, 1)
}
func BenchmarkComplexPictureChunks2x2x528w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 528, 1)
}
func BenchmarkComplexPictureChunks1x1x1024w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 1)
}

func BenchmarkComplexPictureChunks1024x1024x1w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 2)
}
func BenchmarkComplexPictureChunks352x352x3w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 352, 352, 3, 2)
}
func BenchmarkComplexPictureChunks96x96x11w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 96, 96, 11, 2)
}
func BenchmarkComplexPictureChunks32x32x33w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 33, 2)
}
func BenchmarkComplexPictureChunks16x16x66w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 66, 2)
}
func BenchmarkComplexPictureChunks8x8x132w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 132, 2)
}
func BenchmarkComplexPictureChunks4x4x264w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 264, 2)
}
func BenchmarkComplexPictureChunks2x2x528w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 528, 2)
}
func BenchmarkComplexPictureChunks1x1x1024w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 2)
}

func BenchmarkComplexPictureChunks1024x1024x1w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 3)
}
func BenchmarkComplexPictureChunks352x352x3w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 352, 352, 3, 3)
}
func BenchmarkComplexPictureChunks96x96x11w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 96, 96, 11, 3)
}
func BenchmarkComplexPictureChunks32x32x33w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 33, 3)
}
func BenchmarkComplexPictureChunks16x16x66w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 66, 3)
}
func BenchmarkComplexPictureChunks8x8x132w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 132, 3)
}
func BenchmarkComplexPictureChunks4x4x264w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 264, 3)
}
func BenchmarkComplexPictureChunks2x2x528w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 528, 3)
}
func BenchmarkComplexPictureChunks1x1x1024w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 3)
}

func BenchmarkComplexPictureChunks1024x1024x1w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 4)
}
func BenchmarkComplexPictureChunks352x352x3w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 352, 352, 3, 4)
}
func BenchmarkComplexPictureChunks96x96x11w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 96, 96, 11, 4)
}
func BenchmarkComplexPictureChunks32x32x33w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 33, 4)
}
func BenchmarkComplexPictureChunks16x16x66w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 66, 4)
}
func BenchmarkComplexPictureChunks8x8x132w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 132, 4)
}
func BenchmarkComplexPictureChunks4x4x264w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 264, 4)
}
func BenchmarkComplexPictureChunks2x2x528w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 528, 4)
}
func BenchmarkComplexPictureChunks1x1x1024w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 4)
}

func BenchmarkComplexPictureChunks1024x1024x1w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 4)
}
func BenchmarkComplexPictureChunks352x352x3w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 352, 352, 3, 4)
}
func BenchmarkComplexPictureChunks96x96x11w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 96, 96, 11, 4)
}
func BenchmarkComplexPictureChunks32x32x33w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 33, 4)
}
func BenchmarkComplexPictureChunks16x16x66w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 66, 4)
}
func BenchmarkComplexPictureChunks8x8x132w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 132, 4)
}
func BenchmarkComplexPictureChunks4x4x264w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 264, 4)
}
func BenchmarkComplexPictureChunks2x2x528w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 528, 4)
}
func BenchmarkComplexPictureChunks1x1x1024w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 4)
}

func BenchmarkComplexPictureChunks1024x1024x1w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 4)
}
func BenchmarkComplexPictureChunks352x352x3w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 352, 352, 3, 4)
}
func BenchmarkComplexPictureChunks96x96x11w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 96, 96, 11, 4)
}
func BenchmarkComplexPictureChunks32x32x33w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 33, 4)
}
func BenchmarkComplexPictureChunks16x16x66w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 66, 4)
}
func BenchmarkComplexPictureChunks8x8x132w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 132, 4)
}
func BenchmarkComplexPictureChunks4x4x264w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 264, 4)
}
func BenchmarkComplexPictureChunks2x2x528w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 528, 4)
}
func BenchmarkComplexPictureChunks1x1x1024w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 4)
}
