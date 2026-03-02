package random

import (
	"math"

	"github.com/kelindar/noise"
)

type WienerNoise struct {
	Seed  uint32
	Sigma float64

	time  float64
	value float64
}

func NewWienerNoise(seed uint32, sigma float64) *WienerNoise {
	return &WienerNoise{
		Seed:  seed,
		Sigma: sigma,
	}
}

func (w *WienerNoise) Step(dt float64) float64 {
	w.time += dt

	// Brownian scaling
	stddev := w.Sigma * math.Sqrt(dt)

	n := float64(noise.White(w.Seed, uint(uint64(w.time*1e9))))
	w.value = stddev * n

	return w.value
}

func (w *WienerNoise) Value() float64 { return w.value }

type WhiteNoise struct {
	Seed  uint32
	Sigma float64

	time  uint64
	value float64
}

func NewWhiteNoise(seed uint32, sigma float64) *WienerNoise {
	return &WienerNoise{
		Seed:  seed,
		Sigma: sigma,
	}
}

func (w *WienerNoise) Next() float64 {
	w.time++

	n := float64(noise.White(w.Seed, uint(uint64(w.time*1e9))))
	w.value = w.Sigma * n

	return w.value
}

func (w *WhiteNoise) Value() float64 { return w.value }
