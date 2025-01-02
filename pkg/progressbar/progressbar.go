package progressbar

import (
	"math"
	"sync"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func New(numBars int) (*mpb.Progress, *sync.WaitGroup) {
	wg := &sync.WaitGroup{}
	wg.Add(numBars)

	p := mpb.New(mpb.WithWaitGroup(wg))

	return p, wg
}

func Add(p *mpb.Progress, capacity int, name string) *mpb.Bar {
	return p.AddBar(int64(capacity),
		mpb.BarClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit("%3d / %3d", decor.WC{W: int(math.Trunc(math.Log10(float64(capacity))))*2 + 6, C: decor.DidentRight}),
			decor.OnComplete(decor.Name(""), "done!"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Percentage(decor.WC{W: 4}), ""),
			decor.OnComplete(decor.Name(" | "), ""),
			decor.OnComplete(decor.AverageSpeed(0, "%.2f it/s", decor.WC{W: 7, C: decor.DidentRight}), ""),
		),
	)
}
