// Package eventrouter provides first-match event routing for portwatch.
// Routes are evaluated in registration order; the first matching route
// handles the event. If no route matches, ErrNoMatch is returned.
package eventrouter
