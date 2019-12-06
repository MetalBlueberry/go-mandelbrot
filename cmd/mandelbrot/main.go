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
	"runtime/trace"
	"time"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func main() {
	top := flag.Float64("top", 1.5, "Top mandelbrot position")
	left := flag.Float64("left", -2.1, "Left mandelbrot position")
	areaSize := flag.Float64("areaSize", 3, "From the TopLeft, the size of the complex area")

	imageSize := flag.Int("imageSize", 1920, "Size of the squared image generated in pixels")
	divisions := flag.Int("divisions", 50, "Number of divisions to split the work over multiple routines")
	maxIterations := flag.Int("maxIterations", 100, "Maximum number of iterations per point")

	workers := flag.Int("workers", runtime.NumCPU(), "Maximum number of iterations per point")
	out := flag.String("out", "mandelbrot.jpg", "output file, it can be png or jpg")
	timeout := flag.Int64("timeout", 20, "Maximum number of seconds to compute, if reached. the program will exit")

	flag.Parse()

	log.Printf("Start")

	pic := mandelbrot.NewPicture(complex(*left, *top), *areaSize, *imageSize, *divisions, *maxIterations)
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

func Calculate(timeout int64, workers int, pic *mandelbrot.Picture) (*image.RGBA, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer ctxCancel()
	ctx, task := trace.NewTask(ctx, "Calculate")
	defer task.End()
	doneIndex := pic.CalculateAsync(ctx, workers)
	img := image.NewRGBA(image.Rect(0, 0, pic.HorizontalResolution(), pic.VerticalResolution()))

	for {
		select {
		case <-ctx.Done():
			log.Print("CANCEL")
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
