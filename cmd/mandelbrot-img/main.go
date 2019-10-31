package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"time"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func main() {
	set := mandelbrot.NewAreaCentered(5000, 350, -0.7463, 0.1102, 0.005)
	set.Init()
	progress := set.Calculate()
	for p := range progress {
		fmt.Fprintf(os.Stderr, "%d\r", 100*p/len(set.Points))
		<-time.After(time.Second * 1)
	}

	img := image.NewRGBA(image.Rect(0, 0, set.HorizontalResolution, set.VerticalResolution))

	for x := 0; x < set.HorizontalResolution; x++ {
		for y := 0; y < set.VerticalResolution; y++ {
			point := set.GetPoint(x, y)
			intensity := 255 - (255 * point.Iterations() / set.MaxIterations)
			img.SetRGBA(x, y, color.RGBA{R: uint8(intensity), A: 255})
		}
	}
	err := jpeg.Encode(os.Stdout, img, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}
