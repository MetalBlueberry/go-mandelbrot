package mandelbrot_test

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"testing"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

var pic *mandelbrot.Picture

func TestSimpleRenderTrace(t *testing.T) {
	pic = &mandelbrot.Picture{
		TopLeft:               complex(-1.401854499759, -0.000743603637),
		MaxIterations:         1000,
		ChunkSize:             0.00021646,
		HorizontalImageChunks: 4,
		VerticalImageChunks:   4,
		ChunkImageSize:        254,
	}
	ctx := context.Background()
	done := make(chan int)
	pic.Init()

	go pic.Calculate(ctx, 1, done)
	for range done {

	}
}

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

		go pic.Calculate(ctx, workers, done)
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

var saveImage = flag.Bool("saveImage", false, "save the result of running the benchmarks")

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
		pic.Init()
		b.StartTimer()

		done := pic.CalculateAsync(ctx, workers)
		for range done {

		}
		if *saveImage {
			b.StopTimer()
			saveImageFrom(pic, fmt.Sprintf("bench_%d_%d_%d_w_%d.jpeg", pic.HorizontalImageChunks, pic.VerticalImageChunks, pic.ChunkImageSize, workers))
		}
	}
}

func saveImageFrom(pic *mandelbrot.Picture, name string) {
	img := image.NewRGBA(image.Rect(0, 0, pic.HorizontalResolution(), pic.VerticalResolution()))
	for i := 0; i < pic.HorizontalImageChunks*pic.VerticalImageChunks; i++ {
		offsetX, offsetY := pic.GetImageOffsetFor(i)
		paintAreaInImage(img, pic.GetArea(i), offsetX, offsetY)
	}
	outFile, err := os.Create(name)
	if err != nil {
		log.Fatalf("output file cannot be opened, cause: %s", err)
	}
	defer outFile.Close()

	encodingError := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	if encodingError != nil {
		panic(err)
	}
}

func paintAreaInImage(img *image.RGBA, area mandelbrot.Area, offsetX int, offsetY int) {
	for x := 0; x < area.HorizontalResolution; x++ {
		for y := 0; y < area.VerticalResolution; y++ {
			point := area.GetPoint(x, y)
			color := getColor(point, []color.RGBA{
				color.RGBA{
					R: 255,
					A: 255,
				},
				color.RGBA{
					G: 255,
					A: 255,
				},
				color.RGBA{
					B: 255,
					A: 255,
				},
				color.RGBA{
					R: 255,
					G: 255,
					A: 255,
				},
				color.RGBA{
					G: 255,
					B: 255,
					A: 255,
				},
				color.RGBA{
					R: 255,
					B: 255,
					A: 255,
				},
				color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 255,
				},
			}, area.MaxIterations, color.RGBA{
				A: 255,
			})
			img.SetRGBA(offsetX+x, offsetY+y, color)
		}
	}
}
func getColor(point mandelbrot.Point, palette []color.RGBA, maxIterations int, maxIterationsColor color.RGBA) color.RGBA {
	if point.Iterations() == maxIterations {
		return maxIterationsColor
	}
	index := point.Iterations() % len(palette)
	return palette[index]
}
func BenchmarkComplexPictureChunks1024x1024x1w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 1)
}
func BenchmarkComplexPictureChunks512x512x2w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 1)
}
func BenchmarkComplexPictureChunks256x256x4w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 1)
}
func BenchmarkComplexPictureChunks128x128x8w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 1)
}
func BenchmarkComplexPictureChunks64x64x16w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 1)
}
func BenchmarkComplexPictureChunks32x32x32w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 1)
}
func BenchmarkComplexPictureChunks16x16x64w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 1)
}
func BenchmarkComplexPictureChunks8x8x128w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 1)
}
func BenchmarkComplexPictureChunks4x4x256w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 1)
}
func BenchmarkComplexPictureChunks2x2x512w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 1)
}
func BenchmarkComplexPictureChunks1x1x1024w1(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 1)
}

func BenchmarkComplexPictureChunks1024x1024x1w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 2)
}
func BenchmarkComplexPictureChunks512x512x2w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 2)
}
func BenchmarkComplexPictureChunks256x256x4w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 2)
}
func BenchmarkComplexPictureChunks128x128x8w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 2)
}
func BenchmarkComplexPictureChunks64x64x16w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 2)
}
func BenchmarkComplexPictureChunks32x32x32w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 2)
}
func BenchmarkComplexPictureChunks16x16x64w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 2)
}
func BenchmarkComplexPictureChunks8x8x128w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 2)
}
func BenchmarkComplexPictureChunks4x4x256w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 2)
}
func BenchmarkComplexPictureChunks2x2x512w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 2)
}
func BenchmarkComplexPictureChunks1x1x1024w2(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 2)
}

func BenchmarkComplexPictureChunks1024x1024x1w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 3)
}
func BenchmarkComplexPictureChunks512x512x2w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 3)
}
func BenchmarkComplexPictureChunks256x256x4w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 3)
}
func BenchmarkComplexPictureChunks128x128x8w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 3)
}
func BenchmarkComplexPictureChunks64x64x16w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 3)
}
func BenchmarkComplexPictureChunks32x32x32w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 3)
}
func BenchmarkComplexPictureChunks16x16x64w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 3)
}
func BenchmarkComplexPictureChunks8x8x128w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 3)
}
func BenchmarkComplexPictureChunks4x4x256w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 3)
}
func BenchmarkComplexPictureChunks2x2x512w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 3)
}
func BenchmarkComplexPictureChunks1x1x1024w3(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 3)
}

func BenchmarkComplexPictureChunks1024x1024x1w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 4)
}
func BenchmarkComplexPictureChunks512x512x2w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 4)
}
func BenchmarkComplexPictureChunks256x256x4w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 4)
}
func BenchmarkComplexPictureChunks128x128x8w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 4)
}
func BenchmarkComplexPictureChunks64x64x16w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 4)
}
func BenchmarkComplexPictureChunks32x32x32w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 4)
}
func BenchmarkComplexPictureChunks16x16x64w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 4)
}
func BenchmarkComplexPictureChunks8x8x128w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 4)
}
func BenchmarkComplexPictureChunks4x4x256w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 4)
}
func BenchmarkComplexPictureChunks2x2x512w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 4)
}
func BenchmarkComplexPictureChunks1x1x1024w4(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 4)
}

func BenchmarkComplexPictureChunks1024x1024x1w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 5)
}
func BenchmarkComplexPictureChunks512x512x2w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 5)
}
func BenchmarkComplexPictureChunks256x256x4w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 5)
}
func BenchmarkComplexPictureChunks128x128x8w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 5)
}
func BenchmarkComplexPictureChunks64x64x16w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 5)
}
func BenchmarkComplexPictureChunks32x32x32w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 5)
}
func BenchmarkComplexPictureChunks16x16x64w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 5)
}
func BenchmarkComplexPictureChunks8x8x128w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 5)
}
func BenchmarkComplexPictureChunks4x4x256w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 5)
}
func BenchmarkComplexPictureChunks2x2x512w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 5)
}
func BenchmarkComplexPictureChunks1x1x1024w5(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 5)
}

func BenchmarkComplexPictureChunks1024x1024x1w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 6)
}
func BenchmarkComplexPictureChunks512x512x2w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 6)
}
func BenchmarkComplexPictureChunks256x256x4w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 6)
}
func BenchmarkComplexPictureChunks128x128x8w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 6)
}
func BenchmarkComplexPictureChunks64x64x16w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 6)
}
func BenchmarkComplexPictureChunks32x32x32w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 6)
}
func BenchmarkComplexPictureChunks16x16x64w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 6)
}
func BenchmarkComplexPictureChunks8x8x128w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 6)
}
func BenchmarkComplexPictureChunks4x4x256w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 6)
}
func BenchmarkComplexPictureChunks2x2x512w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 6)
}
func BenchmarkComplexPictureChunks1x1x1024w6(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 6)
}

func BenchmarkComplexPictureChunks1024x1024x1w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 7)
}
func BenchmarkComplexPictureChunks512x512x2w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 7)
}
func BenchmarkComplexPictureChunks256x256x4w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 7)
}
func BenchmarkComplexPictureChunks128x128x8w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 7)
}
func BenchmarkComplexPictureChunks64x64x16w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 7)
}
func BenchmarkComplexPictureChunks32x32x32w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 7)
}
func BenchmarkComplexPictureChunks16x16x64w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 7)
}
func BenchmarkComplexPictureChunks8x8x128w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 7)
}
func BenchmarkComplexPictureChunks4x4x256w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 7)
}
func BenchmarkComplexPictureChunks2x2x512w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 7)
}
func BenchmarkComplexPictureChunks1x1x1024w7(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 7)
}
func BenchmarkComplexPictureChunks1024x1024x1w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1024, 1024, 1, 8)
}
func BenchmarkComplexPictureChunks512x512x2w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 512, 512, 2, 8)
}
func BenchmarkComplexPictureChunks256x256x4w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 256, 256, 4, 8)
}
func BenchmarkComplexPictureChunks128x128x8w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 128, 128, 8, 8)
}
func BenchmarkComplexPictureChunks64x64x16w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 64, 64, 16, 8)
}
func BenchmarkComplexPictureChunks32x32x32w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 32, 32, 32, 8)
}
func BenchmarkComplexPictureChunks16x16x64w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 16, 16, 64, 8)
}
func BenchmarkComplexPictureChunks8x8x128w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 8, 8, 128, 8)
}
func BenchmarkComplexPictureChunks4x4x256w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 4, 4, 256, 8)
}
func BenchmarkComplexPictureChunks2x2x512w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 2, 2, 512, 8)
}
func BenchmarkComplexPictureChunks1x1x1024w8(b *testing.B) {
	benchmarkComplexPictureChunks(b, 1, 1, 1024, 8)
}
