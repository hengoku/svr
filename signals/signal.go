package signals

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
	"math"
	"os"
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
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: s.xVals,
				YValues: s.yVals,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		panic(err)
	}

	f, err := os.Create("lab1.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := buffer.WriteTo(f); err != nil {
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
