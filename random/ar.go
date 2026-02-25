package random

import "github.com/kelindar/noise"

type Ar struct {
	Seed uint32
	Phi  []float32

	t uint
	x []float32
}

func NewAr(seed uint32, phi []float32) *Ar {
	return &Ar{
		Seed: seed,
		Phi:  phi,
		t:    0,
		x:    make([]float32, len(phi)),
	}
}

func (a *Ar) Next() float32 {
	a.t++

	var sum float32
	p := uint(len(a.Phi))
	for i := uint(0); i < p; i++ {
		sum += a.Phi[i] * a.x[(a.t-i)%p]
	}

	next := sum + noise.White(a.Seed, a.t)

	a.x = append(a.x, next)
	a.x = a.x[1:]

	return next
}
