package random

import "github.com/kelindar/noise"

type Ar struct {
	Seed uint32
	P    uint
	Phi  []float32

	t uint
	x []float32
}

func NewAr(seed uint32, p uint, phi []float32) *Ar {
	if len(phi) != int(p) {
		panic("phi length must match p")
	}

	return &Ar{
		Seed: seed,
		P:    p,
		Phi:  phi,
		t:    0,
		x:    make([]float32, p),
	}
}

func (a *Ar) Next() float32 {
	a.t++

	var sum float32
	for i := uint(0); i < a.P; i++ {
		sum += a.Phi[i] * a.x[(a.t-i)%uint(a.P)]
	}

	next := sum + noise.White(a.Seed, a.t)

	a.x = append(a.x, next)
	a.x = a.x[1:]

	return next
}
