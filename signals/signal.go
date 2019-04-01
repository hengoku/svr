package signals

import (
	"github.com/wcharczuk/go-chart"
	"math"
	"time"
	"wrg/rts/lab/draws"
)

type Signal struct {
	WMax float64
	HNum int

	Generator Generator

	harmonics []Harmonic

	xVals []float64
	yVals []float64
}

func (s *Signal) GenerateHarmonics() {
	s.harmonics = make([]Harmonic, s.HNum)
	for i := 0; i < s.HNum; i++ {
		s.harmonics[i] = Harmonic{
			A:  s.Generator.A(),
			W:  s.WMax * float64((i+1)/s.HNum),
			Fi: s.Generator.Fi(),
		}
	}
}

func (s *Signal) Count(fromT, toT, nDiscrete float64) {
	s.xVals, s.yVals = []float64{}, []float64{}

	// count time shift for each iteration
	tShift := (toT - fromT) / nDiscrete

	s.GenerateHarmonics()
	for t := fromT; t <= toT; t += tShift {
		var sum float64
		for i := 0; i < len(s.harmonics); i++ {
			sum += s.harmonics[i].Count(t)
		}
		s.xVals = append(s.xVals, t)
		s.yVals = append(s.yVals, sum)

	}
}

func (s *Signal) Draw() {
	if err := draws.DrawWith(chart.ContinuousSeries{XValues: s.xVals, YValues: s.xVals}, "lab1.png"); err != nil {
		panic(err)
	}
}

func (s *Signal) ExpectedValue() float64 {
	// find similar
	sim := make(map[float64]float64)
	for _, v := range s.xVals {
		sim[v]++
	}

	mathExp := 0.0
	for k, v := range sim {
		mathExp += k * v / float64(len(s.yVals))
	}

	return mathExp
}

func (s *Signal) Dispersion() float64 {
	// find similar
	sim := make(map[float64]float64)
	for _, v := range s.xVals {
		sim[v]++
	}

	mathExp := 0.0
	for k, v := range sim {
		mathExp += k * k * v / float64(len(s.yVals))
	}

	return mathExp - math.Pow(s.ExpectedValue(), 2)
}

func (s *Signal) Bench() {
	xVals, yVals := []float64{}, []float64{}

	for i := 1.0; i <= s.WMax; i++ {
		xVals = append(xVals, i)
		tFrom := time.Now()
		s.Count(0, 1, i)
		yVals = append(yVals, float64(time.Since(tFrom).Nanoseconds()))
	}

	if err := draws.DrawWith(chart.ContinuousSeries{XValues: xVals, YValues: yVals}, "bench.png"); err != nil {
		panic(err)
	}
}
