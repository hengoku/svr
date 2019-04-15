package signals

import (
	"github.com/mjibson/go-dsp/fft"
	"math"
	"testing"
)

var (
	sig = &Signal{
		WMax: 2100,
		HNum: 6,
		Generator: Generator{
			ABot:  0,
			ATop:  1,
			FiBot: 0,
			FiTop: 2 * math.Pi,
		},
	}
)

func BenchmarkSignal_DFTFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sig.Count(0, float64(i), 1)
		sig.GenerateHarmonics()

		sig.DFTFast()
	}
}

func BenchmarkSignal_DFTSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sig.Count(0, float64(i), 1)
		sig.GenerateHarmonics()

		sig.DFTSimple()
	}
}

func BenchmarkSignal_DFTImported(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sig.Count(0, float64(i), 1)
		sig.GenerateHarmonics()

		fft.FFTReal(sig.yVals)
	}
}
