package random

import "github.com/kelindar/noise"

type WhiteNoise struct {
	Seed uint32

	t int
}

func (wn *WhiteNoise) Next() float32 {
	wn.t++
	return noise.White(wn.Seed, wn.t)
}
