package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/metalblueberry/mandelbrot/mandelbrot"
)

func main() {
	set := mandelbrot.Area{
		HorizontalResolution: 1060,
		VerticalResolution:   730,
		MaxIterations:        3534,
		TopLeft:              complex(-1.401854499759, -0.000743603637),
		BottomRight:          complex(-1.399689899172, 0.000743603637),
	}
	set.Init()
	set.Calculate()

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
