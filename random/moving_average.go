package random

import "github.com/kelindar/noise"

type MovingAverage struct {
	Seed  uint32
	Mu    float64
	Theta []float64

	t uint
	e []float64
}

func NewMovingAverage(seed uint32, mu float64, theta []float64) *MovingAverage {
	return &MovingAverage{
		Seed:  seed,
		Mu:    mu,
		Theta: theta,
		t:     0,
		e:     make([]float64, len(theta)),
	}
}

func (a *MovingAverage) Next() float64 {
	a.t++

	var sum float64
	q := uint(len(a.Theta))
	for i := uint(0); i < q; i++ {
		sum += a.Theta[i] * a.e[(a.t-i)%q]
	}

	et := float64(noise.White(a.Seed, a.t))
	next := a.Mu + sum + et

	a.e = append(a.e, et)
	a.e = a.e[1:]

	return next
}
