package signals

import (
	"math/rand"
)

type Generator struct {
	ABot, ATop   float64
	FiBot, FiTop float64
}

func (g Generator) A() float64 {
	return g.ABot + rand.Float64()*(g.ATop-g.ABot)
}

func (g Generator) Fi() float64 {
	return g.FiBot + rand.Float64()*(g.FiTop-g.FiBot)
}
