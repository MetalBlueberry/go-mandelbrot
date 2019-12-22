package mandelbrot

type Interlaced struct {
	TopLeft       complex128
	BottomRight   complex128
	MaxIterations int

	ImageWidth  int
	ImageHeight int

	areas []Area
}

func (p *Interlaced) Init() {
	p.areas = make([]Area, p.ImageHeight)
	for i := 0; i < len(p.areas); i++ {

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
