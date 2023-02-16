package client

import (
	"math"
	"sort"

	"container/ring"
)

type slidingWindow struct {
	size        int
	currentSize int
	sum         float64
	samples     *ring.Ring
}

func newSlidingWindow(size int) *slidingWindow {
	if size < 1 {
		size = 1
	}
	w := &slidingWindow{
		size:    size,
		samples: ring.New(size),
	}
	for i := 0; i < w.size; i++ {
		w.samples.Value = math.NaN()
		w.samples = w.samples.Next()
	}
	return w
}

func (w *slidingWindow) add(sample float64) {
	w.samples = w.samples.Next()
	v := w.samples.Value.(float64)
	if !math.IsNaN(v) {
		w.sum -= v
	}
	if w.currentSize < w.size {
		w.currentSize++
	}
	w.samples.Value = float64(sample)
	w.sum += float64(sample)
}

func (w *slidingWindow) lastSample() float64 {
	return w.samples.Value.(float64)
}

func (w *slidingWindow) allSamples() []float64 {
	s := []float64{}
	r := w.samples
	for j := 0; j < w.size; j++ {
		v := r.Value.(float64)
		if !math.IsNaN(v) {
			s = append(s, v)
		}
		r = r.Prev()
	}
	return s
}

func mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func (w *slidingWindow) median() float64 {
	c := w.allSamples()
	sort.Float64s(c)
	l := len(c)
	if l == 0 {
		return math.NaN()
	} else if l%2 == 0 {
		return mean(c[l/2-1 : l/2+1])
	}
	return c[l/2]
}

func (w *slidingWindow) mean() float64 {
	return w.sum / float64(w.currentSize)
}
