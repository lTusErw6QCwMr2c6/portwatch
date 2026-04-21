// Package eventpipeline provides a composable, ordered pipeline for processing
// alert.Event values through a series of named stages.
package eventpipeline

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// Stage is a single processing step that transforms a slice of events.
type Stage struct {
	Name    string
	Process func([]alert.Event) []alert.Event
}

// Pipeline executes an ordered sequence of stages against an event slice.
type Pipeline struct {
	stages []Stage
}

// New returns an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Register appends a named stage to the pipeline.
// Returns an error if the name is empty or already registered.
func (p *Pipeline) Register(s Stage) error {
	if s.Name == "" {
		return fmt.Errorf("eventpipeline: stage name must not be empty")
	}
	for _, existing := range p.stages {
		if existing.Name == s.Name {
			return fmt.Errorf("eventpipeline: stage %q already registered", s.Name)
		}
	}
	if s.Process == nil {
		return fmt.Errorf("eventpipeline: stage %q has nil Process func", s.Name)
	}
	p.stages = append(p.stages, s)
	return nil
}

// Stages returns the names of all registered stages in order.
func (p *Pipeline) Stages() []string {
	names := make([]string, len(p.stages))
	for i, s := range p.stages {
		names[i] = s.Name
	}
	return names
}

// Run executes each stage in registration order, passing the output of one
// stage as the input to the next. Returns the final event slice.
func (p *Pipeline) Run(events []alert.Event) []alert.Event {
	current := events
	for _, s := range p.stages {
		if len(current) == 0 {
			break
		}
		current = s.Process(current)
	}
	return current
}

// Len returns the number of registered stages.
func (p *Pipeline) Len() int {
	return len(p.stages)
}
