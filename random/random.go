package random

type Random interface {
	Next() float64
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
