package random

import (
	"math"

	"github.com/kelindar/noise"
)

type RandomWalk struct {
	Seed uint32

	Mu    float64 // long-term mean
	Theta float64 // mean reversion strength (0 = pure random walk)
	Sigma float64 // noise scale

	t     uint
	value float64
}

func NewRandomWalk(seed uint32, mu, theta, sigma float64) *RandomWalk {
	return &RandomWalk{
		Seed:  seed,
		Mu:    mu,
		Theta: theta,
		Sigma: sigma,
		value: mu,
	}
}

// NewRWWithTimescale creates an OU process with a physical time constant.
// mean: long term mean
// tau:  time constant in seconds (how long memory lasts)
// sigma: continuous volatility (units per sqrt(second))
// dt:   step duration in seconds
func NewRandomWalkWithTimescale(seed uint32, mean, tau, sigma, dt float64) *RandomWalk {
	theta := dt / tau
	stepSigma := sigma * math.Sqrt(dt)

	return NewRandomWalk(seed, mean, theta, stepSigma)
}

func (rw *RandomWalk) Next() float64 {
	rw.t++

	noiseTerm := float64(noise.White(rw.Seed, rw.t))

	rw.value += rw.Theta*(rw.Mu-rw.value) + rw.Sigma*noiseTerm

	return rw.value
}
