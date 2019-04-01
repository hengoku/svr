package signals

import (
	"math"
)

type Harmonic struct {
	A, W, Fi float64
}

func (h Harmonic) Count(t float64) float64 {
	return h.A * math.Sin(h.W*t+h.Fi)
}
