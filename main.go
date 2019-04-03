package main

import (
	"github.com/wcharczuk/go-chart"
	"math"
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

	s1.GenerateHarmonics()
	s2.GenerateHarmonics()
	s1.Count(0, 512, 1)
	s2.Count(0, 512, 1)
	// s2.ChangeNTo(1024)

	s1.Draw()

	xVals, yVals, err := s1.Correlation(s2, 24)
	if err != nil {
		panic(err)
	}

	if err := draws.DrawWith(chart.ContinuousSeries{XValues: xVals, YValues: yVals}, "correlation.png"); err != nil {
		panic(err)
	}

	xVals, yVals = s1.AutoCorrelation(12)
	if err := draws.DrawWith(chart.ContinuousSeries{XValues: xVals, YValues: yVals}, "s1_auto.png"); err != nil {
		panic(err)
	}

	xVals, yVals = s2.AutoCorrelation(12)
	if err := draws.DrawWith(chart.ContinuousSeries{XValues: xVals, YValues: yVals}, "s2_auto.png"); err != nil {
		panic(err)
	}

}
