package random

type Random interface {
	Next() float64
	Value() float64
}

type TimeDependantRandom interface {
	Step(float64) float64
	Value() float64
}

type MultiRandom struct {
	value     float64
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
	m.value = sum
	return sum
}

func (m *MultiRandom) Value() float64 {
	return m.value
}

type TimeDependantMultiRandom struct {
	value     float64
	processes []TimeDependantRandom
}

func NewTimeDependantMultiRandom(processes ...TimeDependantRandom) *TimeDependantMultiRandom {
	return &TimeDependantMultiRandom{processes: processes}
}

func (m *TimeDependantMultiRandom) Step(dt float64) float64 {
	var sum float64
	for _, p := range m.processes {
		sum += p.Step(dt)
	}
	m.value = sum
	return sum
}

func (m *TimeDependantMultiRandom) Value() float64 {
	return m.value
}
