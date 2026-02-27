package random

import "github.com/kelindar/noise"

type AutoRegression struct {
	Seed uint32
	Phi  []float64

	t uint
	x []float64
}

func NewAutoRegression(seed uint32, phi []float64) *AutoRegression {
	return &AutoRegression{
		Seed: seed,
		Phi:  phi,
		t:    0,
		x:    make([]float64, len(phi)),
	}
}

func (a *AutoRegression) Next() float64 {
	a.t++

	var sum float64
	p := uint(len(a.Phi))
	for i := uint(0); i < p; i++ {
		sum += a.Phi[i] * a.x[(a.t-i)%p]
	}

	next := sum + float64(noise.White(a.Seed, a.t))

	a.x = append(a.x, next)
	a.x = a.x[1:]

	return next
}
