// Package cascade implements a multi-stage event processing chain.
// Each stage can transform, filter, or enrich events before passing
// them to the next stage in the pipeline.
package cascade

import "github.com/user/portwatch/internal/alert"

// Stage is a single processing step in the cascade.
type Stage interface {
	Process(events []alert.Event) []alert.Event
}

// StageFunc is a function adapter for Stage.
type StageFunc func([]alert.Event) []alert.Event

func (f StageFunc) Process(events []alert.Event) []alert.Event {
	return f(events)
}

// Cascade runs events through an ordered list of stages.
type Cascade struct {
	stages []Stage
}

// New creates a new Cascade with the given stages.
func New(stages ...Stage) *Cascade {
	return &Cascade{stages: stages}
}

// Add appends a stage to the cascade.
func (c *Cascade) Add(s Stage) {
	c.stages = append(c.stages, s)
}

// Run processes events through all stages in order.
// If any stage returns an empty slice, processing stops early.
func (c *Cascade) Run(events []alert.Event) []alert.Event {
	current := events
	for _, s := range c.stages {
		if len(current) == 0 {
			return current
		}
		current = s.Process(current)
	}
	return current
}

// Len returns the number of stages.
func (c *Cascade) Len() int {
	return len(c.stages)
}
