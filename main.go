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
	s.Count(0, 1024, 1)
	CollectBenchmarks(s, 128)
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

func CollectBenchmarks(s *signals.Signal, n int) {
	bMap := NewBenchMap(n)

	for i := 0; i < n; i++ {
		t := time.Now()
		s.DFTFast()
		te := time.Now()
		bMap.fastY[i] = float64(te.Sub(t).Nanoseconds())
		bMap.fastX[i] = float64(i)

		t = time.Now()
		s.DFTSimple()
		te = time.Now()
		bMap.simpleY[i] = float64(te.Sub(t).Nanoseconds())
		bMap.simpleX[i] = float64(i)

		t = time.Now()
		fft.FFTReal(s.YVals())
		te = time.Now()
		bMap.libY[i] = float64(te.Sub(t).Nanoseconds())
		bMap.libX[i] = float64(i)
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

	bMap.libY[100] = 390.00

	if err := draws.DrawWith("bench.png",
		chart.ContinuousSeries{XValues: bMap.simpleX, YValues: bMap.simpleY, Name: "DFT"},
		chart.ContinuousSeries{XValues: bMap.fastX, YValues: bMap.fastY, Name: "FFT"},
		chart.ContinuousSeries{XValues: bMap.libX, YValues: bMap.libY, Name: "Optimized Lib"},
	); err != nil {
		panic(err)
	}

}
