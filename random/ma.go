package random

import "github.com/kelindar/noise"

type Ma struct {
	Seed  uint32
	Q     uint
	Mu    float32
	Theta []float32

	t uint
	e []float32
}

func NewMa(seed uint32, q uint, mu float32, theta []float32) *Ma {
	if len(theta) != int(q) {
		panic("phi length must match p")
	}

	return &Ma{
		Seed:  seed,
		Q:     q,
		Mu:    mu,
		Theta: theta,
		t:     0,
		e:     make([]float32, q),
	}
}

func (a *Ma) Next() float32 {
	a.t++

	var sum float32
	for i := uint(0); i < a.Q; i++ {
		sum += a.Theta[i] * a.e[(a.t-i)%uint(a.Q)]
	}

	et := noise.White(a.Seed, a.t)
	next := a.Mu + sum + et

	a.e = append(a.e, et)
	a.e = a.e[1:]

	return next
}
