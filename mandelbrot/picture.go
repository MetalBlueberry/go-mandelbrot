package mandelbrot

import (
	"context"
	"sync"
)

type Picture struct {
	TopLeft       complex128
	ChunkSize     float64
	MaxIterations int

	HorizontalImageChunks int
	VerticalImageChunks   int
	ChunkImageSize        int

	areas []Area
}

func (p *Picture) Init() {
	p.areas = make([]Area, p.HorizontalImageChunks*p.VerticalImageChunks)
	for i := 0; i < len(p.areas); i++ {
		x, y := p.ForIndex(i)
		areaTopLeft := p.TopLeft + complex(p.ChunkSize*float64(x), -p.ChunkSize*float64(y))
		areaBottomRight := areaTopLeft + complex(p.ChunkSize, -p.ChunkSize)

		p.areas[i] = Area{
			TopLeft:              areaTopLeft,
			BottomRight:          areaBottomRight,
			HorizontalResolution: p.ChunkImageSize,
			VerticalResolution:   p.ChunkImageSize,
			MaxIterations:        p.MaxIterations,
		}
		p.areas[i].Init()
	}
}

func (p *Picture) Calculate(ctx context.Context, doneIndex chan<- int, workerCount int) {
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)

	next := make(chan int)
	go workQueue(ctx, next, len(p.areas))

	for worker := 0; worker < workerCount; worker++ {
		go doWork(wg, p.areas, next, doneIndex)
	}

	wg.Wait()

	close(doneIndex)
}

// workerQueue publish in a channel all the pending works to be done.
func workQueue(ctx context.Context, next chan<- int, workCount int) {
	for i := 0; i < workCount; i++ {
		select {
		case <-ctx.Done():
			return
		case next <- i:

		}
	}
	close(next)
}

func doWork(wg *sync.WaitGroup, areas []Area, next <-chan int, doneIndex chan<- int) {
	defer wg.Done()
	for i := range next {
		areas[i].Calculate()
		doneIndex <- i
	}
}

func (p *Picture) HorizontalResolution() int {
	return p.ChunkImageSize * p.HorizontalImageChunks
}

func (p *Picture) VerticalResolution() int {
	return p.ChunkImageSize * p.VerticalImageChunks
}

// IndexFor is an utility function to locate a x,y coordinate in the areas slice
func (p *Picture) IndexFor(x, y int) int {
	return x + y*p.HorizontalImageChunks
}

// ForIndex is an utility function to get the x,y values for a given index in the areas slice
func (p *Picture) ForIndex(i int) (x, y int) {
	y = i / p.HorizontalImageChunks
	x = i % p.HorizontalImageChunks
	return x, y
}

func (p *Picture) GetImageOffsetFor(index int) (width, hight int) {
	x, y := p.ForIndex(index)
	return x * p.ChunkImageSize, y * p.ChunkImageSize
}

func (p *Picture) GetArea(index int) Area {
	return p.areas[index]
}
