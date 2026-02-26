package geo

import (
	"github.com/SwatBros/fakeiot/random"
)

type Latitude struct {
	Min    float32
	Max    float32
	Random random.Random
}

func NewLatitude(random random.Random) *Latitude {
	return &Latitude{
		Min:    -90,
		Max:    90,
		Random: random,
	}
}

func NewLatitudeWithRange(min, max float32, random random.Random) *Latitude {
	return &Latitude{
		Min:    min,
		Max:    max,
		Random: random,
	}
}

func (l *Latitude) Next() float32 {
	span := l.Max - l.Min
	return l.Min + l.Random.Next()*span
}

type Longitude struct {
	Min    float32
	Max    float32
	Random random.Random
}

func NewLongitude(random random.Random) *Longitude {
	return &Longitude{
		Random: random,
	}
}

func NewLongitudeWithRange(min, max float32, random random.Random) *Longitude {
	return &Longitude{
		Min:    min,
		Max:    max,
		Random: random,
	}
}

func (l *Longitude) Next() float32 {
	span := l.Max - l.Min
	return l.Min + l.Random.Next()*span
}
