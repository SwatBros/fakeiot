package random

import "github.com/kelindar/noise"

type White struct {
	Seed uint32

	t int
}

func (r *White) Next() float32 {
	r.t++
	return noise.White(r.Seed, r.t)
}
