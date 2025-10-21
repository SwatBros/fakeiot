package fakeiot

import (
	"fmt"
	"time"
)

type Generator struct {
	World   World
	Sensors []Sensor
	Hooks   []Hook
}

func (g *Generator) GenerateSteps(steps uint) error {
	for i := uint(0); i < steps; i++ {
		g.World.Update()

		for _, hook := range g.Hooks {
			if hook.CheckCondition(g) {
				if err := hook.Run(g); err != nil {
					return err
				}
			}
		}

		for _, sensor := range g.Sensors {
			data := sensor.CollectData()
			if err := sensor.SendData(data); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) GenerateInterval(start, end time.Time, step time.Duration) error {
	ticks := uint(end.Sub(start) / step)
	return g.GenerateSteps(ticks)
}

func (g *Generator) RealTime(tick time.Duration) error {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for range ticker.C {
		g.World.Update()

		for _, hook := range g.Hooks {
			if hook.CheckCondition(g) {
				if err := hook.Run(g); err != nil {
					return err
				}
			}
		}

		for _, sensor := range g.Sensors {
			data := sensor.CollectData()
			if err := sensor.SendData(data); err != nil {
				return err
			}
		}
	}

	return nil
}

type World interface {
	// Update the world state
	Update()
}

type Clock interface {
	// Get the current time
	Now() time.Time
	// Tick the clock forward by one time step
	Tick()
	// TimeStep returns the duration of a time step
	TimeStep() time.Duration
}

type Sensor interface {
	// Collect data from the sensor
	CollectData() any
	// Send data
	SendData(data any) error
}

// NewSensor creates a new typesafe sensor with the given collection and sending functions.
func NewSensor[D any](c func() D, s func(D) error) Sensor {
	return &genericSensor[D]{collect: c, send: s}
}

type genericSensor[D any] struct {
	collect func() D
	send    func(D) error
}

func (gs *genericSensor[D]) CollectData() any {
	return gs.collect()
}

func (gs *genericSensor[D]) SendData(data any) error {
	d, ok := data.(D)
	if !ok {
		return fmt.Errorf("type mismatch: expected %T", *new(D))
	}
	return gs.send(d)
}
