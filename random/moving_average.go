package random

import "github.com/kelindar/noise"

type MovingAverage struct {
	Seed  uint32
	Mu    float32
	Theta []float32

	t uint
	e []float32
}

func NewMovingAverage(seed uint32, mu float32, theta []float32) *MovingAverage {
	return &MovingAverage{
		Seed:  seed,
		Mu:    mu,
		Theta: theta,
		t:     0,
		e:     make([]float32, len(theta)),
	}
}

func (a *MovingAverage) Next() float32 {
	a.t++

	var sum float32
	q := uint(len(a.Theta))
	for i := uint(0); i < q; i++ {
		sum += a.Theta[i] * a.e[(a.t-i)%q]
	}

	et := noise.White(a.Seed, a.t)
	next := a.Mu + sum + et

	a.e = append(a.e, et)
	a.e = a.e[1:]

	return next
}
