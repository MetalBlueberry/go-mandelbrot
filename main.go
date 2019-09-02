package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type MandelbrotPoint struct {
	ImagePoint    image.Point
	Point         complex128
	Z             complex128
	Iterations    int
	MaxIterations int
}

func (m *MandelbrotPoint) Calculate() {
	for !m.Diverges() && m.Iterations < m.MaxIterations {
		m.Iterations++
		m.Z = cmplx.Pow(m.Z, 2) + m.Point
	}
}

func (m *MandelbrotPoint) Diverges() bool {
	return cmplx.Abs(m.Z) > complex(2, 0)
}

func (m *MandelbrotPoint) Color() color.Color {
	return color.RGBA{
		R: 255 - uint8(255.0*float64(m.Iterations)/float64(m.MaxIterations)),
		A: 255,
	}
}

type Mandelbrot struct {
	Resolution  int
	TopLeft     complex128
	BottomRight complex128
	Image       image.Image
	Iterations  int
}

func NewMandelbrot(Resolution, Iterations int, x, y, area float64) *Mandelbrot {
	return &Mandelbrot{
		Resolution:  Resolution,
		Iterations:  Iterations,
		Image:       image.NewRGBA(image.Rect(0, 0, Resolution, Resolution)),
		TopLeft:     complex(x-area, y+area),
		BottomRight: complex(x+area, y-area),
	}
}

type Changeable interface {
	Set(x, y int, c color.Color)
}

func (m *Mandelbrot) Next() {
	totalPoints := float64(m.Resolution * m.Resolution)
	var points int32
	t := time.NewTicker(time.Millisecond * 200)
	go func() {
		for {
			_, ok := <-t.C
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "%0f\r", float64(points)/totalPoints)
		}
	}()
	waitgroup := sync.WaitGroup{}
	waitgroup.Add(m.Resolution)
	drawer := m.Image.(Changeable)
	for x := 0; x < m.Resolution; x++ {
		go func(x int) {
			for y := 0; y < m.Resolution; y++ {
				number := m.GetNumber(x, y)
				point := MandelbrotPoint{
					Point:         number,
					MaxIterations: m.Iterations,
				}
				point.Calculate()
				drawer.Set(x, y, point.Color())
				atomic.AddInt32(&points, 1)
			}
			waitgroup.Done()
		}(x)
	}
	waitgroup.Wait()
	t.Stop()
}

func (m *Mandelbrot) GetNumber(x, y int) complex128 {
	TopLeftReal := real(m.TopLeft)
	TopLeftImag := imag(m.TopLeft)
	BottomRightReal := real(m.BottomRight)
	BottomRightImag := imag(m.BottomRight)

	real := TopLeftReal + (float64(x)/float64(m.Resolution))*(BottomRightReal-TopLeftReal)
	imag := TopLeftImag + (float64(y)/float64(m.Resolution))*(BottomRightImag-TopLeftImag)
	return complex(real, imag)
}

func main() {
	set := NewMandelbrot(5000, 350, -0.7463, 0.1102, 0.005)
	set.Next()
	png.Encode(os.Stdout, set.Image)
}
