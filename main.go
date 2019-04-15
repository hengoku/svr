package main

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/wcharczuk/go-chart"
	"math"
	"time"
	"wrg/rts/lab/draws"
	"wrg/rts/lab/signals"
)

func main() {
	s1 := &signals.Signal{
		WMax: 2100,
		HNum: 6,
		Generator: signals.Generator{
			ABot:  0,
			ATop:  1,
			FiBot: 0,
			FiTop: 2 * math.Pi,
		},
	}

	s2 := &signals.Signal{
		WMax: 2100,
		HNum: 6,
		Generator: signals.Generator{
			ABot:  0,
			ATop:  1,
			FiBot: 0,
			FiTop: 2 * math.Pi,
		},
	}

	s1.Count(0, 1024, 1)
	s2.Count(0, 1024, 1)

	xVals, yVals, err := s1.Correlation(s2, 100)
	if err != nil {
		panic(err)
	}

	draws.DrawWith("correlation.png", chart.Series(chart.ContinuousSeries{
		XValues: xVals,
		YValues: yVals,
	}))

	xVals, yVals = s1.AutoCorrelation(100)
	draws.DrawWith("auto.png", chart.Series(chart.ContinuousSeries{
		XValues: xVals,
		YValues: yVals,
	}))

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

func CollectBenchmarks(s *signals.Signal, n int) {
	bMap := NewBenchMap(n)

	for i := 0; i < n; i++ {
		t := time.Now()
		s.DFTSimple()
		bMap.simpleY[i] = float64(time.Since(t).Nanoseconds())
		bMap.simpleX[i] = float64(i)

		t = time.Now()
		s.DFTFast()
		bMap.fastY[i] = float64(time.Since(t).Nanoseconds())
		bMap.fastX[i] = float64(i)

		t = time.Now()
		fft.FFTReal(s.YVals())
		bMap.libY[i] = float64(time.Since(t).Nanoseconds())
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

	if err := draws.DrawWith("bench.png",
		chart.ContinuousSeries{XValues: bMap.simpleX, YValues: bMap.simpleY},
		chart.ContinuousSeries{XValues: bMap.fastX, YValues: bMap.fastY},
		chart.ContinuousSeries{XValues: bMap.libX, YValues: bMap.libY},
	); err != nil {
		panic(err)
	}

}
