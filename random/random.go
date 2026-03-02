package random

type Random interface {
	Next() float64
	Value() float64
}

type TimeDependantRandom interface {
	Next(float64) float64
	Value() float64
}

type MultiRandom struct {
	processes []Random
}

func NewMultiRandom(processes ...Random) *MultiRandom {
	return &MultiRandom{processes: processes}
}

func (m *MultiRandom) Next() float64 {
	var sum float64
	for _, p := range m.processes {
		sum += p.Next()
	}
	return sum
}

type TimeDependantMultiRandom struct {
	processes []TimeDependantRandom
}

func NewTimeDependantMultiRandom(processes ...TimeDependantRandom) *TimeDependantMultiRandom {
	return &TimeDependantMultiRandom{processes: processes}
}

func (m *TimeDependantMultiRandom) Next(dt float64) float64 {
	var sum float64
	for _, p := range m.processes {
		sum += p.Next(dt)
	}
	return sum
}
