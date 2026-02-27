package random

import "github.com/kelindar/noise"

type WhiteNoise struct {
	Seed uint32

	t int
}

func (wn *WhiteNoise) Next() float64 {
	wn.t++
	return float64(noise.White(wn.Seed, wn.t))
}
