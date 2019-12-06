package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"
	"testing"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

var result image.Image

func TestSimpleRenderTrace(t *testing.T) {
	// pic := &mandelbrot.Picture{
	// 	TopLeft:               complex(-1.401854499759, -0.000743603637),
	// 	MaxIterations:         1000,
	// 	ChunkSize:             0.00021646,
	// 	HorizontalImageChunks: 2,
	// 	VerticalImageChunks:   2,
	// 	ChunkImageSize:        512,
	// }
	pic := mandelbrot.NewPicture(complex(-1.401854499759, -0.000743603637), 0.00021646*1024, 1024, 32, 1000)
	var err error
	pic.Init()
	result, err = Calculate(100, 6, pic)
	if err != nil {
		log.Panic(err)
	}

	outFile, err := os.Create("test.jpg")
	if err != nil {
		log.Fatalf("output file cannot be opened, cause: %s", err)
	}
	defer outFile.Close()

	encodingError := jpeg.Encode(outFile, result, &jpeg.Options{Quality: 90})
	if encodingError != nil {
		panic(err)
	}
}
