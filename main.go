package main

import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"math"
	"time"
	"wrg/rts/lab/draws"
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
}
