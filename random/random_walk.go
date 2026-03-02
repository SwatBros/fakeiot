package random

import (
	"math"

	"github.com/kelindar/noise"
)

type RandomWalk struct {
	Seed  uint32
	Mu    float64
	Tau   float64
	Sigma float64

	time  float64
	value float64
}

func NewRandomWalk(seed uint32, mu, tau, sigma float64) *RandomWalk {
	return &RandomWalk{
		Seed:  seed,
		Mu:    mu,
		Tau:   tau,
		Sigma: sigma,
		value: mu,
	}
}

func (rw *RandomWalk) Step(dt float64) float64 {
	rw.time += dt

	expTerm := math.Exp(-dt / rw.Tau)

	variance := (rw.Sigma * rw.Sigma) * (rw.Tau / 2.0) * (1.0 - expTerm*expTerm)
	stddev := math.Sqrt(variance)

	noiseTerm := float64(noise.White(rw.Seed, uint(uint64(rw.time*1e9))))

	rw.value = rw.Mu +
		(rw.value-rw.Mu)*expTerm +
		stddev*noiseTerm

	return rw.value
}

func (rw *RandomWalk) Value() float64 { return rw.value }

func (rw *RandomWalk) Reset() {
	rw.time = 0
	rw.value = rw.Mu
}
