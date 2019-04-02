package signals

import (
	"errors"
	"github.com/wcharczuk/go-chart"
	"math"
	"reflect"
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

func (s *Signal) CountAt(t float64) float64 {
	r := 0.0
	for _, h := range s.harmonics {
		r += h.Count(t)
	}

	return r
}

func (s *Signal) Correlation(s2 *Signal, maxTau int) ([]float64, []float64, error) {
	if !reflect.DeepEqual(s.xVals, s2.xVals) {
		return nil, nil, errors.New("time scales are not equal")
	}

	var xVals []float64
	var yVals []float64

	for tau := 1; tau <= maxTau; tau++ {
		result := 0.0
		for i := 0; i < len(s.yVals); i++ {
			result += (s.yVals[i] - s.ExpectedValue()) * (s2.CountAt(s2.xVals[i]+float64(tau)) - s2.ExpectedValue())
		}

		xVals = append(xVals, float64(tau))
		yVals = append(yVals, result/float64(len(s.xVals)-1))
	}

	return xVals, yVals, nil
}

func (s *Signal) AutoCorrelation(maxTau int) ([]float64, []float64) {
	var xVals []float64
	var yVals []float64

	for tau := 1; tau <= maxTau; tau++ {
		res := 0.0
		for i := 0; i < len(s.xVals); i++ {
			res += s.yVals[i] * s.CountAt(s.xVals[i]+float64(tau))
		}

		res /= float64(maxTau)

		xVals = append(xVals, float64(tau))
		yVals = append(yVals, res)

	}

	return xVals, yVals
}
