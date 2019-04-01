package main

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"math"
	"os"
	"time"
	"wrg/rts/lab/signals"
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
	s.Count(0, 1, 32)
	s.Draw()

	fmt.Printf("Expected value: %v\nDispersion: %v\n", s.ExpectedValue(), s.Dispersion())

	Bench(s)
	ExtraTask(s, 1, 300)
}

func ExtraTask(s *signals.Signal, fromN, toN float64) {
	xVals, yVals := []float64{}, []float64{}

	for n := fromN; n <= toN; n++ {
		s.Count(0, 1, n)
		xVals = append(xVals, n)
		yVals = append(yVals, s.ExpectedValue())
	}

	if err := DrawWith(chart.ContinuousSeries{XValues: xVals, YValues: yVals}, "extra_task.png"); err != nil {
		panic(err)
	}
}

func Bench(s *signals.Signal) {
	xVals, yVals := []float64{}, []float64{}

	for i := 1.0; i <= s.WMax; i++ {
		xVals = append(xVals, i)
		tFrom := time.Now()
		s.Count(0, 1, i)
		yVals = append(yVals, float64(time.Since(tFrom).Nanoseconds()))
	}

	if err := DrawWith(chart.ContinuousSeries{XValues: xVals, YValues: yVals}, "bench.png"); err != nil {
		panic(err)
	}
}

func DrawWith(series chart.Series, filename string) error {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},

		Series: []chart.Series{series},
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		panic(err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := buffer.WriteTo(f); err != nil {
		return err
	}

	return nil
}
