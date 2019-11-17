package main

import (
	"context"
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func main() {
	top := flag.Float64("Top", 1.5, "Top mandelbrot position")
	left := flag.Float64("Left", -2.5, "Left mandelbrot position")
	chunkSize := flag.Float64("ChunkSize", 0.3333333, "Complex Chunk size")
	horizontalImageChunks := flag.Int("HorizontalImageChunks", 16, "Count of horizontal image chunks")
	verticalImageChunks := flag.Int("VerticalImageChunks", 9, "Count of vertical image chunks")
	chunkImageSize := flag.Int("ChunkImageSize", 120, "Chunk dimensions for image size, image dimensions are horizontalImageChunks*ChunkImageSize x verticalImageChunks*ChunkImageSize")
	maxIterations := flag.Int("MaxIterations", 1000, "Maximun number of iterations per point")

	workers := flag.Int("Workers", runtime.NumCPU(), "Maximun number of iterations per point")

	flag.Parse()
	log.Printf("Start")

	mandelPicture := mandelbrot.Picture{
		TopLeft:               complex(*left, *top),
		ChunkSize:             *chunkSize,
		MaxIterations:         *maxIterations,
		HorizontalImageChunks: *horizontalImageChunks,
		VerticalImageChunks:   *verticalImageChunks,
		ChunkImageSize:        *chunkImageSize,
	}
	mandelPicture.Init()

	img := image.NewRGBA(image.Rect(0, 0, mandelPicture.HorizontalResolution(), mandelPicture.VerticalResolution()))

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*20)
	doneIndex := make(chan int)

	go mandelPicture.Calculate(ctx, doneIndex, *workers)
	log.Printf("Calculation sent")

	for i := range doneIndex {
		log.Printf("Index %d done", i)
		offsetX, offsetY := mandelPicture.GetImageOffsetFor(i)
		paintAreaInImage(img, mandelPicture.GetArea(i), offsetX, offsetY)
	}

	ctxCancel()

	err := jpeg.Encode(os.Stdout, img, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

func paintAreaInImage(img *image.RGBA, area mandelbrot.Area, offsetX int, offsetY int) {
	for x := 0; x < area.HorizontalResolution; x++ {
		for y := 0; y < area.VerticalResolution; y++ {
			point := area.GetPoint(x, y)
			color := getColor(point, []color.RGBA{
				color.RGBA{
					A: 255,
				},
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
			})
			img.SetRGBA(offsetX+x, offsetY+y, color)
		}
	}
}

func getColor(point mandelbrot.Point, palette []color.RGBA) color.RGBA {
	index := point.Iterations() % len(palette)
	return palette[index]
}
