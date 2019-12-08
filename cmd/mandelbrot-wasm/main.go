package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"runtime/trace"
	"syscall/js"
	"time"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func mandelbrotDraw(this js.Value, i []js.Value) interface{} {
	pic := mandelbrot.NewPicture(complex(-1.401854499759, -0.000743603637), 0.00021646*1024, 1024, 1, 1000)
	var err error
	pic.Init()
	result, err := Calculate(100, 6, pic)
	if err != nil {
		log.Panic(err)
	}

	buf := &bytes.Buffer{}

	encodingError := jpeg.Encode(buf, result, &jpeg.Options{Quality: 90})
	if encodingError != nil {
		panic(err)
	}
	return js.ValueOf(
		b64.StdEncoding.EncodeToString(buf.Bytes()),
	)
}

func registerCallbacks() {
	js.Global().Set("mandelbrotDraw", js.FuncOf(mandelbrotDraw))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
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
