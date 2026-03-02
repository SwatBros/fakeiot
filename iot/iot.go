package iot

import (
	"math"

	"github.com/SwatBros/fakeiot/random"
)

type Generator interface {
	Next(delta float64) float64
}

type TemperatureGenerator struct {
	random random.TimeDependantRandom

	dailySeasonality  *SineGenerator
	yearlySeasonality *SineGenerator
}

func NewTemperatureGenerator(random random.TimeDependantRandom, dailySeasonality, yearlySeasonality *SineGenerator) *TemperatureGenerator {
	return &TemperatureGenerator{
		random: random,
	}
}

// Advances the generator by the given delta time and return the next temperature value
//
// delta: The time elapsed since the last call to Next() in seconds
//
// Returns the next temperature value
func (tg *TemperatureGenerator) Next(delta float64) float64 {
	cycle := 0.
	if tg.dailySeasonality != nil {
		cycle = tg.dailySeasonality.Next(delta)
	}

	if tg.yearlySeasonality != nil {
		cycle += tg.yearlySeasonality.Next(delta)
	}

	return cycle + float64(tg.random.Step(delta))
}

type TemperatureGeneratorBuilder struct {
	random            random.TimeDependantRandom
	dailySeasonality  *SineGenerator
	yearlySeasonality *SineGenerator
}

func NewTemperatureGeneratorBuilder(random random.TimeDependantRandom) *TemperatureGeneratorBuilder {
	return &TemperatureGeneratorBuilder{
		random: random,
	}
}

func (tgb *TemperatureGeneratorBuilder) WithDailySeasonality(generator *SineGenerator) *TemperatureGeneratorBuilder {
	tgb.dailySeasonality = generator
	return tgb
}

func (tgb *TemperatureGeneratorBuilder) WithDefaultDailySeasonality(mean, amplitude float64) *TemperatureGeneratorBuilder {
	tgb.dailySeasonality = NewSine(amplitude, 86400.0, .3*math.Pi, mean)
	return tgb
}

func (tgb *TemperatureGeneratorBuilder) WithYearlySeasonality(generator *SineGenerator) *TemperatureGeneratorBuilder {
	tgb.yearlySeasonality = generator
	return tgb
}

func (tgb *TemperatureGeneratorBuilder) WithDefaultNorthernHemisphereYearlySeasonality(mean, amplitude float64) *TemperatureGeneratorBuilder {
	tgb.yearlySeasonality = NewSine(amplitude, 31536000.0, -.3*math.Pi, mean)
	return tgb
}

func (tgb *TemperatureGeneratorBuilder) Build() *TemperatureGenerator {
	return &TemperatureGenerator{
		random:            tgb.random,
		dailySeasonality:  tgb.dailySeasonality,
		yearlySeasonality: tgb.yearlySeasonality,
	}
}

type SineGenerator struct {
	Amplitude     float64
	AngularFreq   float64
	PhaseShift    float64
	VerticalShift float64

	t float64
}

func NewSine(amplitude, period, phaseShift, verticalShift float64) *SineGenerator {
	return &SineGenerator{Amplitude: amplitude, AngularFreq: 2 * math.Pi / period, PhaseShift: phaseShift, VerticalShift: verticalShift}
}

func (s *SineGenerator) Next(delta float64) float64 {
	s.t += delta
	return s.Amplitude*math.Sin(s.PhaseShift+s.t*s.AngularFreq) + s.VerticalShift
}
