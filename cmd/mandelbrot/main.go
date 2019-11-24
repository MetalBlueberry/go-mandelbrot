package main

import (
	"context"
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func main() {
	top := flag.Float64("top", 1.5, "Top mandelbrot position")
	left := flag.Float64("left", -2.1, "Left mandelbrot position")
	chunkSize := flag.Float64("chunkSize", 0.3, "Complex Chunk size")
	horizontalImageChunks := flag.Int("horizontalImageChunks", 10, "Count of horizontal image chunks")
	verticalImageChunks := flag.Int("verticalImageChunks", 10, "Count of vertical image chunks")
	chunkImageSize := flag.Int("chunkImageSize", 120, "Chunk dimensions for image size, image dimensions are horizontalImageChunks*ChunkImageSize x verticalImageChunks*ChunkImageSize")
	maxIterations := flag.Int("maxIterations", 100, "Maximun number of iterations per point")

	workers := flag.Int("workers", runtime.NumCPU(), "Maximun number of iterations per point")
	out := flag.String("out", "mandelbrot.jpg", "output file, it can be png or jpg")
	timeout := flag.Int64("timeout", 20, "Maximum number of seconds to compute, if reached. the program will exit")

	flag.Parse()

	if (*horizontalImageChunks)*(*verticalImageChunks)*(*chunkImageSize)*(*chunkImageSize) > 8294400 {
		log.Print("You are trying to generate a big image... your system my crash because you run out of memory, you have 5 seconds to ctrl+c to cancel.")
		time.Sleep(time.Second * 5)
	}

	log.Printf("Start")

	pic := mandelbrot.Picture{
		TopLeft:               complex(*left, *top),
		ChunkSize:             *chunkSize,
		MaxIterations:         *maxIterations,
		HorizontalImageChunks: *horizontalImageChunks,
		VerticalImageChunks:   *verticalImageChunks,
		ChunkImageSize:        *chunkImageSize,
	}
	pic.Init()

	log.Printf("Calculation started")

	img, err := Calculate(*timeout, *workers, pic)
	if err != nil {
		log.Printf("Calculation failed, image is not complete. cause: %s", err)
	}

	outFile, err := os.Create(*out)
	if err != nil {
		log.Fatalf("output file cannot be opened, cause: %s", err)
	}
	defer outFile.Close()

	switch filepath.Ext(*out) {
	case ".jpg", ".jpeg":
		encodingError := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
		if encodingError != nil {
			panic(err)
		}
	case ".png":
		encodingError := png.Encode(outFile, img)
		if encodingError != nil {
			panic(err)
		}
	}
}

func Calculate(timeout int64, workers int, pic mandelbrot.Picture) (*image.RGBA, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer ctxCancel()
	doneIndex := pic.CalculateAsync(ctx, workers)
	img := image.NewRGBA(image.Rect(0, 0, pic.HorizontalResolution(), pic.VerticalResolution()))

	for {
		select {
		case <-ctx.Done():
			return img, ctx.Err()
		case i, ok := <-doneIndex:
			if !ok {
				log.Print("Finished")
				return img, nil
			}
			log.Printf("Index %d done", i)
			offsetX, offsetY := pic.GetImageOffsetFor(i)
			paintAreaInImage(img, pic.GetArea(i), offsetX, offsetY)
		}
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
			img.Bounds().In()
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
