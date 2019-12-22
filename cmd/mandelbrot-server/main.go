package main

import (
	"context"
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"net/http"
	"runtime/trace"
	"strconv"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", "./page", "directory to serve")
)

func handler(w http.ResponseWriter, r *http.Request) {

	sheight := r.URL.Query().Get("height")
	swidth := r.URL.Query().Get("width")
	sa := r.URL.Query().Get("a")
	se := r.URL.Query().Get("e")
	sf := r.URL.Query().Get("f")
	log.Printf("Request %s, %s, %s, %s, %s", sheight, swidth, sa, se, sf)

	height, err := strconv.ParseFloat(sheight, 64)
	if err != nil {
		panic(err)
	}
	width, err := strconv.ParseFloat(swidth, 64)
	if err != nil {
		panic(err)
	}
	size := math.Max(height, width)

	a, err := strconv.ParseFloat(sa, 64)
	if err != nil {
		panic(err)
	}
	e, err := strconv.ParseFloat(se, 64)
	if err != nil {
		panic(err)
	}
	f, err := strconv.ParseFloat(sf, 64)
	if err != nil {
		panic(err)
	}

	log.Printf("w %f, h %f, s %f", width, height, size)

	pic := mandelbrot.NewPicture(
		complex(-2.1-3*e/(size*a), 1.5+3*f/(a*size)),
		3/a,
		int(size),
		32,
		1000,
	)
	pic.Init()
	img, err := Calculate(r.Context(), 6, pic)
	if err != nil {
		panic(err)
	}
	encodingError := png.Encode(w, img)
	if encodingError != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	log.Printf("listening on %q...", *listen)

	http.HandleFunc("/mandelbrot", handler)
	http.Handle("/", http.FileServer(http.Dir(*dir)))

	http.ListenAndServe(*listen, nil)
}

func Calculate(ctx context.Context, workers int, pic *mandelbrot.Picture) (*image.RGBA, error) {
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
