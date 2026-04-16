package healthcheck

import (
	"fmt"
	"time"
)

// ProbeFunc is a function that checks a component and returns an error if unhealthy.
type ProbeFunc func() error

// Probe wraps a named ProbeFunc with a timeout.
type Probe struct {
	Name    string
	Fn      ProbeFunc
	Timeout time.Duration
}

// DefaultTimeout is used when no timeout is specified.
const DefaultTimeout = 3 * time.Second

// Run executes the probe and returns a Check result.
func (p Probe) Run() Check {
	timeout := p.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	type result struct {
		err error
	}
	ch := make(chan result, 1)
	go func() {
		ch <- result{err: p.Fn()}
	}()

	select {
	case r := <-ch:
		if r.err != nil {
			return Check{Name: p.Name, Status: StatusDown, Message: r.err.Error(), CheckedAt: time.Now()}
		}
		return Check{Name: p.Name, Status: StatusOK, Message: "ok", CheckedAt: time.Now()}
	case <-time.After(timeout):
		return Check{
			Name:      p.Name,
			Status:    StatusDegraded,
			Message:   fmt.Sprintf("probe timed out after %s", timeout),
			CheckedAt: time.Now(),
		}
	}
}

// RunAll executes all probes and records results into the monitor.
func RunAll(m *Monitor, probes []Probe) {
	for _, p := range probes {
		c := p.Run()
		m.Record(c.Name, c.Status, c.Message)
	}
}
