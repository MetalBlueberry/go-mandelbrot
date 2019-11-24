package mandelbrot

import (
	"context"
	"log"
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

func NewPicture(topLeft complex128, chunkSize float64, imageSize int, divisions int, maxIterations int) *Picture {
	if imageSize%divisions != 0 {
		log.Printf("WARNING: ImageSize %d can't be divided in %d divisions, The final image will be smaller", imageSize, divisions)
	}
	chunkImageSize := imageSize / divisions
	return &Picture{
		TopLeft:               topLeft,
		MaxIterations:         maxIterations,
		ChunkSize:             chunkSize / float64(divisions),
		HorizontalImageChunks: divisions,
		VerticalImageChunks:   divisions,
		ChunkImageSize:        chunkImageSize,
	}
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

func (p *Picture) Calculate(ctx context.Context, workerCount int, doneIndex chan<- int) {
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)

	next := workQueue(ctx, len(p.areas))

	for worker := 0; worker < workerCount; worker++ {
		go doWork(wg, p.areas, next, doneIndex)
	}

	wg.Wait()

	close(doneIndex)
}

func (p *Picture) CalculateAsync(ctx context.Context, workerCount int) <-chan int {
	doneIndex := make(chan int)
	go p.Calculate(ctx, workerCount, doneIndex)
	return doneIndex
}

func workQueue(ctx context.Context, workCount int) <-chan int {
	next := make(chan int)
	go func() {
		for i := 0; i < workCount; i++ {
			select {
			case <-ctx.Done():
				return
			case next <- i:

			}
		}
		close(next)
	}()
	return next
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
