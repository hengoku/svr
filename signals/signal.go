package signals

import (
	"github.com/wcharczuk/go-chart"
	"math"
	"math/cmplx"
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

func (s *Signal) YVals() []float64 {
	return s.yVals
}

func (s *Signal) ChangeNTo(nNew int) {
	if nNew <= len(s.yVals) {
		return
	}

	dt := s.xVals[1] - s.xVals[0]
	for i := len(s.yVals) - 1; i <= nNew; i++ {
		s.xVals = append(s.xVals, s.xVals[len(s.xVals)-1]+dt)
		s.yVals = append(s.yVals, 0)
	}
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

func (s *Signal) Count(fromT, toT, timeShift float64) {
	s.xVals, s.yVals = []float64{}, []float64{}

	// count time shift for each iteration
	s.GenerateHarmonics()

	for t := fromT; t < toT; t += timeShift {
		var sum float64
		for i := 0; i < len(s.harmonics); i++ {
			sum += s.harmonics[i].Count(t)
		}
		s.xVals = append(s.xVals, t)
		s.yVals = append(s.yVals, sum)
	}
}

func (s *Signal) Draw() {
	if err := draws.DrawWith("lab1.png", chart.ContinuousSeries{XValues: s.xVals, YValues: s.yVals}); err != nil {
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

	if err := draws.DrawWith("bench.png", chart.ContinuousSeries{XValues: xVals, YValues: yVals}); err != nil {
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
	var xVals []float64
	var yVals []float64

	for tau := 0; tau <= maxTau; tau++ {
		result := 0.0
		for i := 0; i < len(s2.yVals); i++ {
			result += (s.CountAt(float64(i)) - s.ExpectedValue()) * (s2.CountAt(s2.xVals[i]+float64(tau)) - s2.ExpectedValue())
		}

		xVals = append(xVals, float64(tau))
		yVals = append(yVals, result/float64(len(s.xVals)-1))
	}

	return xVals, yVals, nil
}

func (s *Signal) AutoCorrelation(maxTau int) ([]float64, []float64) {
	var xVals []float64
	var yVals []float64

	for tau := 0; tau <= maxTau; tau++ {
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

func (s *Signal) DFTSimple() []complex128 {
	results := make([]complex128, len(s.yVals))
	arg := -2.0 * math.Pi / float64(len(s.yVals))
	for k := 0; k < len(s.yVals); k++ {
		var res complex128
		for n := 0; n < len(s.yVals); n++ {
			res += complex(s.yVals[n]*math.Cos(arg*float64(n)*float64(k)), s.yVals[n]*math.Sin(arg*float64(n)*float64(k)))
		}
		results[k] = res
	}

	return results
}

func (s *Signal) DFTFast() []complex128 {
	n := len(s.yVals)
	x := toComplex(s.yVals)

	j := 0
	for i := 0; i < n; i++ {
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
		m := n / 2
		for {
			if j < m {
				break
			}
			j = j - m
			m = m / 2
			if m < 2 {
				break
			}
		}
		j = j + m
	}
	kmax := 1
	for {
		if kmax >= n {
			return x
		}
		istep := kmax * 2
		for k := 0; k < kmax; k++ {
			theta := complex(0.0, -1.0*math.Pi*float64(k)/float64(kmax))
			for i := k; i < n; i += istep {
				j := i + kmax
				temp := x[j] * cmplx.Exp(theta)
				x[j] = x[i] - temp
				x[i] = x[i] + temp
			}
		}
		kmax = istep
	}
}
