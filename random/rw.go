package random

import (
	"math"

	"github.com/kelindar/noise"
)

type RW struct {
	Seed uint32

	Mu    float32 // long-term mean
	Theta float32 // mean reversion strength (0 = pure random walk)
	Sigma float32 // noise scale

	t     uint
	value float32
}

func NewRW(seed uint32, mu, theta, sigma float32) *RW {
	return &RW{
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
func NewRWWithTimescale(seed uint32, mean, tau, sigma, dt float64) *RW {
	theta := float32(dt / tau)
	stepSigma := float32(sigma * math.Sqrt(dt))

	return NewRW(seed, float32(mean), theta, stepSigma)
}

func (r *RW) Next() float32 {
	r.t++

	noiseTerm := noise.White(r.Seed, r.t)

	r.value += r.Theta*(r.Mu-r.value) + r.Sigma*noiseTerm

	return r.value
}
