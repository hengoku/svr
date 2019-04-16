package main

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/wcharczuk/go-chart"
	"math"
	"svr/draws"
	"svr/signals"
	"time"
)

func main() {
	s := &signals.Signal{
		WMax: 2100,
		HNum: 6,
		Generator: signals.Generator{
			ABot:  0,
			ATop:  1,
			FiBot: 0,
			FiTop: 2 * math.Pi,
		},
	}

	s.GenerateHarmonics()
	CollectBenchmarks(100, 1000, 100, s)
}

type BenchMap struct {
	simpleX []float64
	simpleY []float64

	fastX []float64
	fastY []float64

	libX []float64
	libY []float64
}

func NewBenchMap(n int) *BenchMap {
	return &BenchMap{
		make([]float64, n),
		make([]float64, n),
		make([]float64, n),
		make([]float64, n),
		make([]float64, n),
		make([]float64, n),
	}
}

func DrawFourier() {

}

func CollectBenchmarks(from, to, step int, s *signals.Signal) {
	n := 11
	bMap := NewBenchMap(to / step)
	j := 0
	for i := from; i <= to; i += step {

		s := &signals.Signal{
			WMax: 2100,
			HNum: 6,
			Generator: signals.Generator{
				ABot:  0,
				ATop:  1,
				FiBot: 0,
				FiTop: 2 * math.Pi,
			},
		}

		s.GenerateHarmonics()
		s.Count(0, float64(i), 1)

		t := time.Now()
		fft.FFTReal(s.YVals())
		te := time.Now()
		bMap.fastY[j] = float64(te.Sub(t).Nanoseconds())
		bMap.fastX[j] = float64(i)

		t = time.Now()
		s.DFTSimple()
		te = time.Now()
		bMap.simpleY[j] = float64(te.Sub(t).Nanoseconds())
		bMap.simpleX[j] = float64(i)

		// t = time.Now()
		// fft.FFTReal(s.YVals())
		// te = time.Now()
		// bMap.libY[j] = float64(te.Sub(t).Nanoseconds())
		// bMap.libX[j] = float64(j)
		j++
	}

	chartX := make([]float64, n)
	chartY := make([]float64, n)
	s.Count(0, float64(n), 1)
	c := fft.FFTReal(s.YVals())
	for i := 0; i < n; i++ {
		chartX[i] = float64(i)
		chartY[i] = real(c[i])*real(c[i]) + imag(c[i])*imag(c[i])
	}

	if err := draws.DrawWith("fourier.png",
		chart.ContinuousSeries{XValues: chartX, YValues: chartY},
	); err != nil {
		panic(err)
	}

	if err := draws.DrawWith("bench.png",
		chart.ContinuousSeries{XValues: bMap.simpleX, YValues: bMap.simpleY, Name: "DFT"},
		chart.ContinuousSeries{XValues: bMap.fastX, YValues: bMap.fastY, Name: "FFT"},
	); err != nil {
		panic(err)
	}

}
