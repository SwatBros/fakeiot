package random

type Random interface {
	Next() float32
}

type MultiRandom struct {
	processes []Random
}

func NewMultiRandom(processes ...Random) *MultiRandom {
	return &MultiRandom{processes: processes}
}

func (m *MultiRandom) Next() float32 {
	var sum float32
	for _, p := range m.processes {
		sum += p.Next()
	}
	return sum
}
